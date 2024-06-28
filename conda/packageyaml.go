package conda

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/robocorp/rcc/cloud"
	"github.com/robocorp/rcc/pathlib"
	"gopkg.in/yaml.v2"
)

type (
	packageDependencies struct {
		CondaForge []string `yaml:"conda-forge,omitempty"`
		Pypi       []string `yaml:"pypi,omitempty"`
	}
	internalPackage struct {
		Dependencies *packageDependencies `yaml:"dependencies"`
		PostInstall  []string             `yaml:"post-install,omitempty"`
	}
)

func (it *internalPackage) AsEnvironment() *Environment {
	result := &Environment{
		Channels:    []string{"conda-forge"},
		PostInstall: []string{},
	}
	seenScripts := make(map[string]bool)
	result.PostInstall = addItem(seenScripts, it.PostInstall, result.PostInstall)
	pushConda(result, it.condaDependencies())
	pushPip(result, it.pipDependencies())
	result.pipPromote()
	return result
}

func fixPipDependency(dependency *Dependency) *Dependency {
	if dependency != nil {
		if dependency.Qualifier == "=" {
			dependency.Original = fmt.Sprintf("%s==%s", dependency.Name, dependency.Versions)
			dependency.Qualifier = "=="
		}
	}
	return dependency
}

func (it *internalPackage) pipDependencies() []*Dependency {
	result := make([]*Dependency, 0, len(it.Dependencies.Pypi))
	for _, item := range it.Dependencies.Pypi {
		dependency := AsDependency(item)
		if dependency != nil {
			result = append(result, fixPipDependency(dependency))
		}
	}
	return result
}

func (it *internalPackage) condaDependencies() []*Dependency {
	result := make([]*Dependency, 0, len(it.Dependencies.CondaForge))
	for _, item := range it.Dependencies.CondaForge {
		dependency := AsDependency(item)
		if dependency != nil {
			result = append(result, dependency)
		}
	}
	return result
}

func packageYamlFrom(content []byte) (*Environment, error) {
	result := new(internalPackage)
	err := yaml.Unmarshal(content, result)
	if err != nil {
		return nil, err
	}
	return result.AsEnvironment(), nil
}

func ReadPackageCondaYaml(filename string) (*Environment, error) {
	basename := strings.ToLower(filepath.Base(filename))
	if basename == "package.yaml" {
		environment, err := ReadPackageYaml(filename)
		if err == nil {
			return environment, nil
		}
	}
	return readCondaYaml(filename)
}

func ReadPackageYaml(filename string) (*Environment, error) {
	var content []byte
	var err error

	if pathlib.IsFile(filename) {
		content, err = os.ReadFile(filename)
	} else {
		content, err = cloud.ReadFile(filename)
	}
	if err != nil {
		return nil, fmt.Errorf("%q: %w", filename, err)
	}
	return packageYamlFrom(content)
}
