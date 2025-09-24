package klusterhelper

import (
	"fmt"
)

type ImageSpec struct {
	Repository string `json:"repository,omitempty"`
	Tag        string `json:"tag"`
}

func (i *ImageSpec) String() string {
	return fmt.Sprintf("%s:%s", i.Repository, i.Tag)
}
