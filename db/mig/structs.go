package mig // import "iris.arke.works/forum/db/mig"

import (
	"errors"
	"gopkg.in/yaml.v2"
	"strings"
)

type Unit struct {
	Name        string     `yaml:"-"`
	Description string     `yaml:"description"`
	DependsOn   []string   `yaml:"depends_on"`
	AlwaysExec  bool       `yaml:"always_exec"`
	Type        UnitType   `yaml:"type"`
	SQL         SQLSection `yaml:"sql"`

	executed bool
}

func (u Unit) DependsOnWithoutNothing() []string {
	var deps = u.DependsOn
	var retDeps = []string{}
	for _, v := range deps {
		if v != "nothing" {
			retDeps = append(retDeps, v)
		}
	}
	return retDeps
}

type SQLSection struct {
	Postgres string `yaml:"postgres"`
}

type UnitType string

const (
	UnitTypeMigration     UnitType = "migration"
	UnitTypeVirtualTarget          = "target"
)

var nothingUnit = Unit{
	Name:        "nothing",
	Description: "Use as dependency if a unit has no dependencies",
	DependsOn:   []string{},
	Type:        UnitTypeMigration,
	SQL: SQLSection{
		Postgres: "",
	},
	executed: true,
}

// loadUnitFile will load and parse a Unit file
// It accepts a basepath and a filename which are joined together. The filename should be relative
// to the basepath and acts as a unitname.
// The filename must end in either .mig or .trgt
func loadUnitFile(basepath, filename string) (*Unit, error) {
	var retUnit = &Unit{}

	if !strings.HasSuffix(filename, ".yaml") {
		return nil, errors.New("Unit file must have file extension .yaml")
	}

	box, err := boxConf.FindBox(basepath)
	if err != nil {
		return nil, err
	}

	dat, err := box.Bytes(filename)
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(dat, retUnit)
	if err != nil {
		return nil, err
	}

	if retUnit.Type == "" {
		retUnit.Type = UnitTypeMigration
	}

	retUnit.Name = strings.TrimSuffix(filename, ".yaml")

	if retUnit.Name == "" {
		return retUnit, errors.New("Could not determine name of unit")
	}

	return retUnit, nil
}
