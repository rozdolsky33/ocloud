package vcn

type Gateways struct {
	InternetGateway   string
	NatGateway        string
	ServiceGateway    string
	Drg               string
	LocalPeeringPeers []string
}
