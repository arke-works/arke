package mig // import "iris.arke.works/forum/db/mig"

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMigrations(t *testing.T) {
	assert := assert.New(t)

	graph := NewGraph()

	assert.Error(graph.Load("non-existent-basepath"))

	assert.NoError(graph.Load("arke"))

	assert.NoError(graph.ValidateNodes())

	oldLoad := loadUnitFile
	loadUnitFile = func(string, string) (*Unit, error) {
		return nil, errors.New("Test")
	}

	assert.Error(graph.Load("arke"))

	loadUnitFile = oldLoad
}

func TestGraph_GetUnit(t *testing.T) {
	assert := assert.New(t)

	graph := Graph{nodes: map[string]*Unit{
		"node1": {
			Name:        "node1",
			Description: "test node",
		},
	}}

	n, err := graph.GetUnit("node1")
	assert.EqualValues(n, Unit{
		Name:        "node1",
		Description: "test node",
	}, "Returning node should contain specified data")
	assert.NoError(err)

	n, err = graph.GetUnit("node-does-not-exist")
	assert.EqualValues(n, Unit{})
	assert.Error(err)
}

func TestGraph_ValidateNodes(t *testing.T) {
	assert := assert.New(t)

	assert.NotNil(NewGraph())

	graph := Graph{nodes: map[string]*Unit{
		"default": {
			Name:      "default1",
			DependsOn: []string{},
		},
	}}

	assert.Error(graph.ValidateNodes())

	graph.nodes["default"].Name = "default"
	graph.nodes["default"].DependsOn = []string{"nothing"}

	graph.nodes["nothing"] = &nothingUnit

	assert.NoError(graph.ValidateNodes())

	graph.nodes["default"].DependsOn = []string{
		"default2",
	}

	assert.Error(graph.ValidateNodes())

	graph.nodes["default"].DependsOn = nil

	assert.Error(graph.ValidateNodes())

	graph.nodes["default"].executed = true

	assert.NoError(graph.ValidateNodes())

}

func TestGraph_CanExecuteNode(t *testing.T) {
	assert := assert.New(t)

	graph := NewGraph()

	graph.nodes["default"] = &Unit{
		Name:      "default",
		DependsOn: []string{"default2"},
	}

	graph.nodes["default2"] = &Unit{
		Name:      "defaul2",
		DependsOn: []string{"nothing"},
	}

	assert.True(graph.CanExecuteNode("default2"))
	assert.False(graph.CanExecuteNode("default"))

	graph.nodes["default2"].executed = true

	assert.True(graph.CanExecuteNode("default"))

	graph.nodes["default"].executed = true

	assert.False(graph.CanExecuteNode("default"))

	assert.False(graph.CanExecuteNode("does-not-exist"))
}

func TestGraph_Size(t *testing.T) {
	if NewGraph().Size() != 1 {
		t.Log("New Graph must have size of 1")
		t.Fail()
	}
}

func TestGraph_RemainingSize(t *testing.T) {
	assert := assert.New(t)

	graph := NewGraph()
	graph.nodes["default"] = &Unit{
		Name:      "default",
		DependsOn: []string{"default2"},
	}
	graph.nodes["default2"] = &Unit{
		Name:      "default2",
		DependsOn: []string{"nothing"},
	}

	assert.Equal(2, graph.RemainingSize())

	graph.nodes["default2"].executed = true

	assert.Equal(1, graph.RemainingSize())

	graph.nodes["default"].executed = true

	assert.Equal(0, graph.RemainingSize())
}

func TestGraph_IsStuck(t *testing.T) {
	graph := NewGraph()

	graph.nodes["default"] = &Unit{
		Name:      "default",
		DependsOn: []string{"default2"},
	}

	graph.nodes["default2"] = &Unit{
		Name:      "default2",
		DependsOn: []string{"default", "nothing"},
	}

	if !graph.IsStuck() {
		t.Log("Graph is stuck but not indicated as stuck")
		t.Fail()
	}

	graph.nodes["default2"].DependsOn = []string{"nothing"}

	if graph.IsStuck() {
		t.Log("Graph is not stuck but indicated as stuck")
		t.Fail()
	}
}

func TestGraph_GetAllRunnableNodes(t *testing.T) {
	assert := assert.New(t)

	graph := NewGraph()

	graph.nodes["default"] = &Unit{
		Name:      "default",
		DependsOn: []string{"default2"},
	}

	graph.nodes["default2"] = &Unit{
		Name:      "default2",
		DependsOn: []string{"nothing"},
	}

	assert.EqualValues([]string{"default2"}, graph.GetAllRunnableNodes())
	assert.Equal(3, graph.Size())
}

func TestGraph_GetTargetSubgraph(t *testing.T) {
	assert := assert.New(t)

	graph := NewGraph()

	graph.nodes["default"] = &Unit{
		Name:      "default",
		DependsOn: []string{"default2"},
		Type:      UnitTypeVirtualTarget,
	}

	graph.nodes["default2"] = &Unit{
		Name:      "default2",
		DependsOn: []string{"default3", "default4", "default3"},
		Type:      UnitTypeVirtualTarget,
	}

	graph.nodes["default3"] = &Unit{
		Name:      "default3",
		DependsOn: []string{"nothing"},
		Type:      UnitTypeMigration,
	}

	graph.nodes["default4"] = &Unit{
		Name:      "default4",
		DependsOn: []string{"default3"},
		Type:      UnitTypeMigration,
	}

	graph.nodes["target-faulty-dep"] = &Unit{
		Name:      "target-faulty-dep",
		DependsOn: []string{"does-not-exist"},
		Type:      UnitTypeVirtualTarget,
	}

	subGraph, err := graph.GetTargetSubgraph("default2")

	assert.NoError(err)

	assert.NoError(subGraph.ValidateNodes())

	assert.Equal(4, subGraph.Size())

	assert.Equal(3, subGraph.RemainingSize())

	_, err = graph.GetTargetSubgraph("default_x")

	assert.Error(err)

	_, err = graph.GetTargetSubgraph("default3")

	assert.Error(err)

	_, err = graph.GetTargetSubgraph("target-faulty-dep")

	assert.Error(err)
}

func TestGraph_MarkNodesRun(t *testing.T) {
	assert := assert.New(t)

	graph := NewGraph()

	graph.nodes["default"] = &Unit{
		Name:      "default",
		DependsOn: []string{"default2"},
		Type:      UnitTypeVirtualTarget,
	}

	graph.nodes["default2"] = &Unit{
		Name:      "default2",
		DependsOn: []string{"default3", "default4", "default3"},
		Type:      UnitTypeVirtualTarget,
	}

	graph.nodes["default3"] = &Unit{
		Name:      "default3",
		DependsOn: []string{"nothing"},
		Type:      UnitTypeMigration,
	}

	graph.nodes["default4"] = &Unit{
		Name:      "default4",
		DependsOn: []string{"default3"},
		Type:      UnitTypeMigration,
	}

	assert.NoError(graph.MarkNodesRun("default4"))

	assert.Error(graph.MarkNodesRun("default4", "Nodefault"))

	assert.NoError(graph.MarkNodesRun())

	assert.True(graph.nodes["default4"].executed)
}
