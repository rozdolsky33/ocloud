package oke

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"

	"github.com/rozdolsky33/ocloud/internal/services/util"
)

// kubeconfig model (unexported)
type kubeConfig struct {
	APIVersion     string         `yaml:"apiVersion"`
	Kind           string         `yaml:"kind"`
	Clusters       []namedCluster `yaml:"clusters"`
	Users          []namedUser    `yaml:"users"`
	Contexts       []namedContext `yaml:"contexts"`
	CurrentContext string         `yaml:"current-context"`
}

// namedCluster represents a named Kubernetes cluster configuration consisting of a name and its corresponding details.
type namedCluster struct {
	Name    string    `yaml:"name"`
	Cluster kcCluster `yaml:"cluster"`
}

// kcCluster represents a Kubernetes cluster configuration consisting of a server address and certificate details.
type kcCluster struct {
	Server                   string `yaml:"server"`
	CertificateAuthorityData string `yaml:"certificate-authority-data,omitempty"`
	InsecureSkipTLSVerify    bool   `yaml:"insecure-skip-tls-verify,omitempty"`
}

// namedUser represents a named Kubernetes user configuration consisting of a name and its corresponding details.
type namedUser struct {
	Name string `yaml:"name"`
	User kcUser `yaml:"user"`
}

// kcUser represents a Kubernetes user configuration.
type kcUser struct {
	Exec *kcExec `yaml:"exec,omitempty"`
}

// kcExec represents a Kubernetes exec configuration.
type kcExec struct {
	APIVersion         string   `yaml:"apiVersion"`
	Command            string   `yaml:"command"`
	Args               []string `yaml:"args"`
	Env                []any    `yaml:"env"`
	InteractiveMode    string   `yaml:"interactiveMode"`
	ProvideClusterInfo bool     `yaml:"provideClusterInfo"`
}

// namedContext represents a named Kubernetes context configuration consisting of a name and its corresponding details.
type namedContext struct {
	Name    string    `yaml:"name"`
	Context kcContext `yaml:"context"`
}

// kcContext represents Kubernetes context details including cluster, namespace, and user mappings.
type kcContext struct {
	Cluster   string `yaml:"cluster"`
	Namespace string `yaml:"namespace"`
	User      string `yaml:"user"`
}

// EnsureKubeconfigForOKE ensures kubeconfig entries for the given cluster/region/profile and local port.
// If entries already exist, it is a no-op.
func EnsureKubeconfigForOKE(cluster Cluster, region, profile string, localPort int) error {
	home, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("get home dir: %w", err)
	}
	kubeDir := filepath.Join(home, ".kube")
	cfgPath := filepath.Join(kubeDir, "config")

	if err := os.MkdirAll(kubeDir, 0o700); err != nil {
		return fmt.Errorf("ensure kube dir: %w", err)
	}

	var kc kubeConfig
	if b, err := os.ReadFile(cfgPath); err == nil {
		_ = yaml.Unmarshal(b, &kc)
	} else if !errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("read kubeconfig: %w", err)
	}
	if kc.APIVersion == "" {
		kc.APIVersion = "v1"
	}

	// First, check if any existing user exec config already targets this exact cluster id, region and profile.
	for _, u := range kc.Users {
		if u.User.Exec == nil {
			continue
		}
		if matchOKEExec(u.User.Exec, cluster.ID, region, profile) {
			// Found an existing matching entry, skip creation entirely.
			return nil
		}
	}

	suffix := shortID(cluster.ID)
	cName := "cluster-" + suffix
	uName := "user-" + suffix
	ctxName := "context-" + suffix

	// If all present by our naming, skip
	if hasNamed(kc.Users, func(n namedUser) bool { return n.Name == uName }) &&
		hasNamed(kc.Clusters, func(n namedCluster) bool { return n.Name == cName }) &&
		hasNamed(kc.Contexts, func(n namedContext) bool { return n.Name == ctxName }) {
		return nil
	}

	if util.PromptYesNo(fmt.Sprintf("Do you want to enter a custom kube context name for this cluster? (Default is '%s')", ctxName)) {
		if name, err := util.PromptString("Enter kube context name", ctxName); err == nil {
			name = strings.TrimSpace(name)
			if name != "" {
				ctxName = name
			}
		}
	}

	server := fmt.Sprintf("https://127.0.0.1:%d", localPort)

	kc.Clusters = upsertCluster(kc.Clusters, namedCluster{
		Name: cName,
		Cluster: kcCluster{
			Server:                server,
			InsecureSkipTLSVerify: true,
		},
	})

	kc.Users = upsertUser(kc.Users, namedUser{
		Name: uName,
		User: kcUser{Exec: &kcExec{
			APIVersion:         "client.authentication.k8s.io/v1beta1",
			Command:            "oci",
			Args:               []string{"ce", "cluster", "generate-token", "--cluster-id", cluster.ID, "--region", region, "--profile", profile, "--auth", "security_token"},
			Env:                []any{},
			InteractiveMode:    "",
			ProvideClusterInfo: false,
		}},
	})

	kc.Contexts = upsertContext(kc.Contexts, namedContext{
		Name: ctxName,
		Context: kcContext{
			Cluster:   cName,
			Namespace: "",
			User:      uName,
		},
	})

	if kc.CurrentContext == "" {
		kc.CurrentContext = ctxName
	}

	// write atomically to the target file.
	// Backup if exists first.
	if _, err := os.Stat(cfgPath); err == nil {
		if old, err := os.ReadFile(cfgPath); err == nil {
			bak := cfgPath + ".bak"
			_ = os.WriteFile(bak, old, 0o600)
		}
	}

	f, err := os.OpenFile(cfgPath, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0o600)
	if err != nil {
		return fmt.Errorf("open kubeconfig for write: %w", err)
	}
	defer f.Close()

	enc := yaml.NewEncoder(f)
	enc.SetIndent(2)
	if err := enc.Encode(&kc); err != nil {
		_ = enc.Close()
		return fmt.Errorf("marshal kubeconfig: %w", err)
	}
	if err := enc.Close(); err != nil {
		return fmt.Errorf("finalize kubeconfig write: %w", err)
	}
	return nil
}

// shortID returns a shortened version of the given cluster id.
func shortID(id string) string {
	// Try to take suffix after the last '.' or '/'
	s := id
	if idx := strings.LastIndex(s, "."); idx >= 0 {
		s = s[idx+1:]
	}
	if idx := strings.LastIndex(s, "/"); idx >= 0 {
		s = s[idx+1:]
	}
	if len(s) > 12 {
		s = s[len(s)-12:]
	}
	return s
}

// hasNamed returns true if the given array contains an element satisfying the given predicate.
func hasNamed[T any](arr []T, pred func(T) bool) bool {
	for _, v := range arr {
		if pred(v) {
			return true
		}
	}
	return false
}

// upsert* functions are used to upsert an element into an array of named elements.
func upsertCluster(arr []namedCluster, item namedCluster) []namedCluster {
	for i, v := range arr {
		if v.Name == item.Name {
			arr[i] = item
			return arr
		}
	}
	return append(arr, item)
}

// upsert* functions are used to upsert an element into an array of named elements.
func upsertUser(arr []namedUser, item namedUser) []namedUser {
	for i, v := range arr {
		if v.Name == item.Name {
			arr[i] = item
			return arr
		}
	}
	return append(arr, item)
}

// upsert* functions are used to upsert an element into an array of named elements.
func upsertContext(arr []namedContext, item namedContext) []namedContext {
	for i, v := range arr {
		if v.Name == item.Name {
			arr[i] = item
			return arr
		}
	}
	return append(arr, item)
}

// matchOKEExec returns true if the kcExec represents an OCI command generating a token for the given
// cluster id, region, and profile.
func matchOKEExec(exec *kcExec, clusterID, region, profile string) bool {
	if exec == nil {
		return false
	}
	if exec.Command != "oci" {
		return false
	}
	if !containsStr(exec.Args, "generate-token") {
		return false
	}
	flags := parseArgsToMap(exec.Args)
	if flags["--cluster-id"] != clusterID {
		return false
	}
	if flags["--region"] != region {
		return false
	}
	if flags["--profile"] != profile {
		return false
	}
	return true
}

// parseArgsToMap converts a slice of CLI args into a simple flag->value map.
// Supports both ["--flag", "value"] and ["--flag=value"] forms.
func parseArgsToMap(args []string) map[string]string {
	out := make(map[string]string)
	for i := 0; i < len(args); i++ {
		a := args[i]
		if strings.HasPrefix(a, "--") {
			if idx := strings.IndexByte(a, '='); idx > 0 {
				key := a[:idx]
				val := a[idx+1:]
				out[key] = val
				continue
			}
			// no equal sign, take next as value if present and not a flag
			if i+1 < len(args) && !strings.HasPrefix(args[i+1], "-") {
				out[a] = args[i+1]
				i++
			} else {
				out[a] = ""
			}
		}
	}
	return out
}

func containsStr(arr []string, needle string) bool {
	for _, v := range arr {
		if v == needle {
			return true
		}
	}
	return false
}

// KubeconfigExistsForOKE checks whether ~/.kube/config already contains an entry
// for the given OKE cluster identified by cluster ID, region and profile.
// It returns true if a matching user exec section is found, false if not.
func KubeconfigExistsForOKE(cluster Cluster, region, profile string) (bool, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return false, fmt.Errorf("get home dir: %w", err)
	}
	cfgPath := filepath.Join(home, ".kube", "config")

	var kc kubeConfig
	if b, err := os.ReadFile(cfgPath); err == nil {
		_ = yaml.Unmarshal(b, &kc)
	} else if errors.Is(err, os.ErrNotExist) {
		// kubeconfig file is not present at all -> no match
		return false, nil
	} else {
		return false, fmt.Errorf("read kubeconfig: %w", err)
	}

	for _, u := range kc.Users {
		if u.User.Exec == nil {
			continue
		}
		if matchOKEExec(u.User.Exec, cluster.ID, region, profile) {
			return true, nil
		}
	}
	return false, nil
}
