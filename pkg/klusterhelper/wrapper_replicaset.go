package klusterhelper

import (
	appsv1 "k8s.io/api/apps/v1"
)

type ReplicaSetWrapper struct {
	*appsv1.ReplicaSet
}

var _ KubeResource = &ReplicaSetWrapper{}

func (r *ReplicaSetWrapper) validate() error          { return nil }
func (r *ReplicaSetWrapper) marshal() ([]byte, error) { return marshalCleanYAML(r.ReplicaSet) }
