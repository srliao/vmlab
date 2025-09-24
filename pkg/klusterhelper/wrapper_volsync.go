package klusterhelper

import (
	vs1alpha1 "github.com/backube/volsync/api/v1alpha1"
)

type ReplicationDestinationWrapper struct {
	*vs1alpha1.ReplicationDestination
}

type ReplicationSourceWrapper struct {
	*vs1alpha1.ReplicationSource
}

var _ KubeResource = &ReplicationDestinationWrapper{}
var _ KubeResource = &ReplicationSourceWrapper{}

func (r *ReplicationDestinationWrapper) validate() error { return nil }
func (r *ReplicationDestinationWrapper) marshal() ([]byte, error) {
	return marshalCleanYAML(r.ReplicationDestination)
}

func (r *ReplicationSourceWrapper) validate() error { return nil }
func (r *ReplicationSourceWrapper) marshal() ([]byte, error) {
	return marshalCleanYAML(r.ReplicationSource)
}
