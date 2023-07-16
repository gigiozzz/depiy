package bundles

import (
	"fmt"
	"io/ioutil"

	"gopkg.in/yaml.v3"
)

type ComponentType string

const (
	ManifestComponentType ComponentType = "MANIFEST"
	PluginComponentType   ComponentType = "PLUGIN"
)

type Component struct {
	Name string        `json:"name,omitempty"`
	Type ComponentType `json:"type,omitempty"`
	Spec interface{}   `yaml:"-"`
}

type Plugin struct {
	IngressName     string `yaml:"ingressName,omitempty"`
	IngressHost     string `yaml:"ingressHost,omitempty"`
	IngressPath     string `yaml:"ingressPath,omitempty"`
	Repository      string `yaml:"repository,omitempty"`
	Tag             string `yaml:"tag,omitempty"`
	Digest          string `yaml:"digest,omitempty"`
	HealthCheckPath string `yaml:"healthCheckPath,omitempty"`
	Port            int    `yaml:"port,omitempty"`
}

type Manifest struct {
	FilePath string `yaml:"filePath,omitempty"`
}

type BundleDescriptor struct {
	Version      string      `yaml:"version"`
	Name         string      `yaml:"name"`
	Descriptor   string      `yaml:"descriptor"`
	Dependencies []string    `yaml:"dependencies"`
	Components   []Component `yaml:"components"`
}

func (s *Component) UnmarshalYAML(n *yaml.Node) error {
	type S Component
	type T struct {
		*S   `yaml:",inline"`
		Spec yaml.Node `yaml:"spec"`
	}

	obj := &T{S: (*S)(s)}
	if err := n.Decode(obj); err != nil {
		return err
	}

	switch s.Type {
	case ManifestComponentType:
		s.Spec = new(Manifest)
	case PluginComponentType:
		s.Spec = new(Plugin)
	default:
		panic(fmt.Sprintf("kind unknown %s", s.Type))
	}
	return obj.Spec.Decode(s.Spec)
}

func (s *Component) GetIfIsPlugin() (bool, *Plugin) {
	plugin, isPlugin := s.Spec.(*Plugin)
	return isPlugin, plugin
}

func (s *Component) GetIfIsManifest() (bool, *Manifest) {
	manifest, isManifest := s.Spec.(*Manifest)
	return isManifest, manifest
}

func ReadBundleDescriptor(fileDescriptorPath string) (*BundleDescriptor, error) {
	yfile, err := ioutil.ReadFile(fileDescriptorPath)
	if err != nil {
		return nil, err
	}

	data := &BundleDescriptor{}
	err = yaml.Unmarshal(yfile, &data)
	return data, err
}
