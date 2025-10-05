package vcn

// Gateways represent a summarized view of gateways associated with a VCN.
// It is used by the OCI gateway adapter to return a concise, human-readable
// snapshot for presentation layers.
type Gateways struct {
	InternetGateway   string   // e.g., "igw-prod (present)" or "—"
	NatGateway        string   // e.g., "nat-prod (present)" or "—"
	ServiceGateway    string   // e.g., "sgw-prod (ObjectStorage, OSN)" or "—"
	Drg               string   // e.g., "drg-core (attached)" or "—"
	LocalPeeringPeers []string // e.g., ["lpg-a → vcn-b", "lpg-c → vcn-d"]
}
