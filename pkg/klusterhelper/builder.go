package klusterhelper

import (
	"fmt"
	"os"
	"path/filepath"

	"sigs.k8s.io/controller-runtime/pkg/client"
)

type KubeResource interface {
	client.Object
	validate() error
	marshal() ([]byte, error)
}

type Annotatable[K any] interface {
	MergeAnnotations(annotations map[string]string) K
	WithAnnotations(annotations map[string]string) K
}

type Builder struct {
	apps []*App
}

func NewBuilder() *Builder {
	return &Builder{}
}

func (b *Builder) AddApp(app *App) *Builder {
	b.apps = append(b.apps, app)
	return b
}

func (b *Builder) Validate() *Builder {
	var errs []string
	for _, c := range b.apps {
		if err := c.Validate(); err != nil {
			errs = append(errs, fmt.Sprintf("chart %v/%v validation failed: %v", c.name, c.namespace, err))
		}
	}
	if len(errs) > 0 {
		for _, e := range errs {
			fmt.Println(e)
		}
		panic("validation failed")
	}
	return b
}

func createFolder(folder string, removeExisting bool) error {
	if removeExisting {
		os.RemoveAll(folder)
	}
	if err := os.MkdirAll(filepath.Join(folder), 0755); err != nil {
		return fmt.Errorf("failed to create folder: %w", err)
	}
	return nil
}

func (b *Builder) Build(folder string, removeExisting bool) error {
	// don't remove root; only apps
	err := createFolder(folder, false)
	if err != nil {
		return err
	}

	for _, c := range b.apps {
		path := folder
		if c.subfolder != "" {
			path = filepath.Join(folder, c.subfolder)
			err := createFolder(path, removeExisting)
			if err != nil {
				return err
			}
		}
		// write ks.yaml
		ksData, err := c.ks.marshal()
		if err != nil {
			return fmt.Errorf("failed to marshal kustomization: %w", err)
		}
		ksPath := filepath.Join(path, "ks.yaml")
		if err := os.WriteFile(ksPath, ksData, 0644); err != nil {
			return fmt.Errorf("failed to write kustomization file: %w", err)
		}

		// separate folder with kustomization.yaml and associated resources

		path = filepath.Join(path, c.ks.subpath)
		err = createFolder(path, removeExisting)
		if err != nil {
			return err
		}

		// write individual resources
		var data []byte
		for _, obj := range c.resources {
			data = append(data, []byte("---\n")...)
			d, err := marshalCleanYAML(obj)
			if err != nil {
				return fmt.Errorf("failed to marshal object: %w", err)
			}
			data = append(data, d...)
		}
		if err := os.WriteFile(filepath.Join(path, fmt.Sprintf("%s.yaml", c.name)), data, 0644); err != nil {
			return fmt.Errorf("failed to write file: %w", err)
		}

		// write additional files
		for _, f := range c.files {
			if err := os.WriteFile(filepath.Join(path, f.Name()), f.Content(), 0644); err != nil {
				return fmt.Errorf("failed to write file %s: %w", f.Name(), err)
			}
		}

		// write kustomization.yaml
		c.kust.AddResources([]string{fmt.Sprintf("%s.yaml", c.name)})
		kust, err := c.kust.marshal()
		if err != nil {
			return fmt.Errorf("failed to marshal kustomization: %w", err)
		}
		if err := os.WriteFile(filepath.Join(path, "kustomization.yaml"), kust, 0644); err != nil {
			return fmt.Errorf("failed to write kustomization file: %w", err)
		}
	}

	return nil
}
