package bastion

import (
	"context"
	"fmt"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/rozdolsky33/ocloud/internal/app"
	"github.com/rozdolsky33/ocloud/internal/logger"
	ociadb "github.com/rozdolsky33/ocloud/internal/oci/database/autonomousdb"
	ocihwdb "github.com/rozdolsky33/ocloud/internal/oci/database/heatwavedb"
	adbSvc "github.com/rozdolsky33/ocloud/internal/services/database/autonomousdb"
	hwdbSvc "github.com/rozdolsky33/ocloud/internal/services/database/heatwavedb"
	bastionSvc "github.com/rozdolsky33/ocloud/internal/services/identity/bastion"
	"github.com/rozdolsky33/ocloud/internal/services/util"
)

// selectDatabaseType runs a TUI to choose between HeatWave and Autonomous Database.
func selectDatabaseType(ctx context.Context) (DatabaseType, error) {
	m := NewDatabaseTypeModel()
	p := tea.NewProgram(m, tea.WithContext(ctx))
	res, err := p.Run()
	if err != nil {
		return "", fmt.Errorf("database type selection TUI: %w", err)
	}
	out, ok := res.(DatabaseTypeModel)
	if !ok || out.Choice == "" {
		return "", ErrAborted
	}
	return out.Choice, nil
}

// connectDatabase runs the DB target flow. We can't always auto-verify reachability,
// so we surface that limitation to the user.
func connectDatabase(ctx context.Context, appCtx *app.ApplicationContext, svc *bastionSvc.Service,
	b bastionSvc.Bastion, sType SessionType) error {

	// Only Port-Forwarding is supported for databases
	if sType != TypePortForwarding {
		logger.Logger.Info("Only Port-Forwarding sessions are supported for database connections")
		return fmt.Errorf("only Port-Forwarding sessions are supported for database connections")
	}

	// Select a database type
	dbType, err := selectDatabaseType(ctx)
	if err != nil {
		return err
	}

	logger.Logger.Info("Selected database type", "type", dbType)

	switch dbType {
	case DatabaseHeatWave:
		return connectHeatWaveDatabase(ctx, appCtx, svc, b)
	case DatabaseAutonomous:
		return connectAutonomousDatabase(ctx, appCtx, svc, b)
	default:
		return fmt.Errorf("unknown database type: %s", dbType)
	}
}

// connectHeatWaveDatabase handles the HeatWave database connection flow.
func connectHeatWaveDatabase(ctx context.Context, appCtx *app.ApplicationContext, svc *bastionSvc.Service,
	b bastionSvc.Bastion) error {

	adapter, err := ocihwdb.NewAdapter(appCtx.Provider)
	if err != nil {
		return fmt.Errorf("error creating HeatWave database adapter: %w", err)
	}
	dbService := hwdbSvc.NewService(adapter, appCtx)

	dbs, _, _, err := dbService.FetchPaginatedHeatWaveDb(ctx, 1000, 0)
	if err != nil {
		return fmt.Errorf("list HeatWave databases: %w", err)
	}
	if len(dbs) == 0 {
		logger.Logger.Info("No HeatWave Databases found.")
		return nil
	}

	dm := NewHeatWaveDBListModelFancy(dbs)
	dp := tea.NewProgram(dm, tea.WithContext(ctx))
	dres, err := dp.Run()
	if err != nil {
		return fmt.Errorf("HeatWave DB selection TUI: %w", err)
	}
	chosen, ok := dres.(ResourceListModel)
	if !ok || chosen.Choice() == "" {
		return ErrAborted
	}

	var db hwdbSvc.HeatWaveDatabase
	for _, d := range dbs {
		if d.ID == chosen.Choice() {
			db = d
			break
		}
	}

	_, reason := svc.CanReach(ctx, b, db.VcnID, db.SubnetId)
	logger.Logger.Info("Reachability to HeatWave DB cannot be automatically verified", "reason", reason)
	logger.Logger.Info("Selected HeatWave database", "name", db.DisplayName, "id", db.ID)

	// Get SSH key pair
	pubKey, privKey, err := SelectSSHKeyPair(ctx)
	if err != nil {
		return err
	}

	// Get region
	region, regErr := appCtx.Provider.Region()
	if regErr != nil {
		return fmt.Errorf("get region: %w", regErr)
	}

	// Default port for MySQL/HeatWave
	defaultPort := 3306
	if db.Port != nil {
		defaultPort = *db.Port
	}

	port, err := util.PromptPort("Enter port to forward (local:target)", defaultPort)
	if err != nil {
		return fmt.Errorf("read port: %w", err)
	}

	// Create a port forwarding session
	sessID, err := svc.EnsurePortForwardSession(ctx, b.ID, db.IpAddress, port, pubKey)
	if err != nil {
		return fmt.Errorf("ensure port forward: %w", err)
	}

	// Build and spawn SSH tunnel
	sshTunnelArgs, err := bastionSvc.BuildPortForwardArgs(privKey, sessID, region, db.IpAddress, port, port)
	if err != nil {
		return fmt.Errorf("build args: %w", err)
	}

	logger.Logger.Info("Starting background tunnel", "args", sshTunnelArgs)
	pid, err := bastionSvc.SpawnDetached(sshTunnelArgs, "/tmp/ssh-tunnel.log")
	if err != nil {
		return fmt.Errorf("spawn detached: %w", err)
	}
	logger.Logger.V(logger.Debug).Info("spawned tunnel", "pid", pid)

	if err := bastionSvc.WaitForListen(port, 5*time.Second); err != nil {
		logger.Logger.Error(err, "warning")
	}

	logFile := fmt.Sprintf("~/.oci/.ocloud/ssh-tunnel-%d.log", port)
	logger.Logger.Info("SSH tunnel started in background", "logs", logFile, "local_port", port, "database", db.DisplayName)
	return nil
}

// connectAutonomousDatabase handles the Autonomous Database connection flow.
func connectAutonomousDatabase(ctx context.Context, appCtx *app.ApplicationContext, svc *bastionSvc.Service,
	b bastionSvc.Bastion) error {

	adapter, err := ociadb.NewAdapter(appCtx.Provider)
	if err != nil {
		return fmt.Errorf("error creating database adapter: %w", err)
	}
	dbService := adbSvc.NewService(adapter, appCtx)

	dbs, _, _, err := dbService.FetchPaginatedAutonomousDb(ctx, 1000, 0)
	if err != nil {
		return fmt.Errorf("list databases: %w", err)
	}
	if len(dbs) == 0 {
		logger.Logger.Info("No Autonomous Databases found.")
		return nil
	}

	dm := NewDBListModelFancy(dbs)
	dp := tea.NewProgram(dm, tea.WithContext(ctx))
	dres, err := dp.Run()
	if err != nil {
		return fmt.Errorf("DB selection TUI: %w", err)
	}
	chosen, ok := dres.(ResourceListModel)
	if !ok || chosen.Choice() == "" {
		return ErrAborted
	}

	var db adbSvc.AutonomousDatabase
	for _, d := range dbs {
		if d.ID == chosen.Choice() {
			db = d
			break
		}
	}

	_, reason := svc.CanReach(ctx, b, db.VcnID, db.SubnetId)
	logger.Logger.Info("Reachability to Autonomous DB cannot be automatically verified", "reason", reason)
	logger.Logger.Info("Selected Autonomous database", "name", db.Name, "id", db.ID)

	// Get SSH key pair
	pubKey, privKey, err := SelectSSHKeyPair(ctx)
	if err != nil {
		return err
	}

	// Get region
	region, regErr := appCtx.Provider.Region()
	if regErr != nil {
		return fmt.Errorf("get region: %w", regErr)
	}

	// Default port for Oracle Database
	defaultPort := 1521
	port, err := util.PromptPort("Enter port to forward (local:target)", defaultPort)
	if err != nil {
		return fmt.Errorf("read port: %w", err)
	}

	// Use private endpoint IP if available
	targetIP := db.PrivateEndpointIp
	if targetIP == "" {
		return fmt.Errorf("no private endpoint IP available for database %s", db.Name)
	}

	// Create a port forwarding session
	sessID, err := svc.EnsurePortForwardSession(ctx, b.ID, targetIP, port, pubKey)
	if err != nil {
		return fmt.Errorf("ensure port forward: %w", err)
	}

	// Build and spawn SSH tunnel
	sshTunnelArgs, err := bastionSvc.BuildPortForwardArgs(privKey, sessID, region, targetIP, port, port)
	if err != nil {
		return fmt.Errorf("build args: %w", err)
	}

	logger.Logger.Info("Starting background tunnel", "args", sshTunnelArgs)
	pid, err := bastionSvc.SpawnDetached(sshTunnelArgs, "/tmp/ssh-tunnel.log")
	if err != nil {
		return fmt.Errorf("spawn detached: %w", err)
	}
	logger.Logger.V(logger.Debug).Info("spawned tunnel", "pid", pid)

	if err := bastionSvc.WaitForListen(port, 5*time.Second); err != nil {
		logger.Logger.Error(err, "warning")
	}

	logFile := fmt.Sprintf("~/.oci/.ocloud/ssh-tunnel-%d.log", port)
	logger.Logger.Info("SSH tunnel started in background", "logs", logFile, "local_port", port, "database", db.Name)
	return nil
}
