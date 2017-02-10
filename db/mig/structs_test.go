package mig

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestLoadUnitFile(t *testing.T) {
	assert := require.New(t)

	_, err := loadUnitFile("this-is-not-a-path", "unit.xml")
	assert.Error(err)

	_, err = loadUnitFile("this-is-not-a-apth", "unit.yaml")
	assert.Error(err)

	_, err = loadUnitFile("arke", "does-not-exist-unit.yaml")
	assert.Error(err)

	_, err = loadUnitFile("mock-migs", "no-marshal.yaml")
	assert.Error(err)

	_, err = loadUnitFile("mock-migs", ".yaml")
	assert.Error(err)
}
