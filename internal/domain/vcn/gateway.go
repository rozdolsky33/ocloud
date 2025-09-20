package vcn

type VcnsGateways struct {
	InternetGateway   string
	NatGateway        string
	ServiceGateway    string
	Drg               string
	LocalPeeringPeers []string
}
