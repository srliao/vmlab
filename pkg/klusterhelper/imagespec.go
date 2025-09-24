package klusterhelper

import (
	"fmt"
)

type ImageSpec struct {
	Registry string `json:"registry,omitempty"`
	Name     string `json:"name"`
	Tag      string `json:"tag"`
}

func (i *ImageSpec) String() string {
	if i.Registry != "" {
		return fmt.Sprintf("%s/%s:%s", i.Registry, i.Name, i.Tag)
	}
	return fmt.Sprintf("%s:%s", i.Name, i.Tag)
}
