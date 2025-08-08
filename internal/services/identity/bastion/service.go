package bastion

// Service provides methods for bastion operations
type Service struct{}

// NewService creates a new bastion service
func NewService() *Service {
	return &Service{}
}

// GetDummyBastions returns a list of dummy bastion options
func (s *Service) GetDummyBastions() []Bastion {
	return []Bastion{
		{ID: "ocid1.bastion.oc1.dummy.bastion1", Name: "bastion_1"},
		{ID: "ocid1.bastion.oc1.dummy.bastion2", Name: "basstion_1"},
		{ID: "ocid1.bastion.oc1.dummy.bastion3", Name: "bestion three"},
	}
}
