package mig // import "iris.arke.works/forum/db/mig"

import (
	"fmt"
	"github.com/GeertJohan/go.rice"
	"github.com/restic/restic/src/restic/errors"
	"os"
)

var boxConf = rice.Config{
	LocateOrder: []rice.LocateMethod{
		rice.LocateEmbedded,
		rice.LocateAppended,
		rice.LocateFS,
		rice.LocateWorkingDirectory,
	},
}

// Graph represents a set of unit files with dependencies that form a Direct Acyclic Graph.
type Graph struct {
	nodes map[string]*Unit
}

// NewGraph creates a graph that only contains the "nothing" unit which is used to bootstrap
// units without dependencies.
//
// This is necessary because the graph has no concept of a node without dependency. This makes determining
// executable nodes much easier (search for nodes with executed dependencies instead of the former *and* searching
// for nodes without dependencies)
func NewGraph() *Graph {
	return &Graph{
		nodes: map[string]*Unit{
			"nothing": &nothingUnit,
		},
	}
}

// Load will use go.rice to load an embedded resources (or from LocalFS if not found) into the graph
func (g *Graph) Load(basepath string) error {
	box, err := boxConf.FindBox(basepath)
	if err != nil {
		return err
	}
	return box.Walk("", func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		unit, err := loadUnitFile(basepath, path)
		if err != nil {
			return err
		}
		g.nodes[unit.Name] = unit
		return nil
	})
}

// GetUnit returns the specified Unit as a struct
func (g *Graph) GetUnit(name string) (Unit, error) {
	if node, ok := g.nodes[name]; ok {
		return *node, nil
	}
	return Unit{}, errors.New("Node not found")
}

// ValidateNodes will check that the graph is properly formed but will not check for cycles.
//
// A properly formed graph contains nodes where their key in the internal map and their name match,
// that nodes always have a dependency if they have not been executed and that the dependencies of
// all nodes exist.
//
// If this function returns no error then it is guaranteed that the graph is directed and may be acyclic.
func (g *Graph) ValidateNodes() error {
	for name, node := range g.nodes {
		if name != node.Name {
			return fmt.Errorf("Node Name and Node Key are mismatched")
		}
		if !node.executed && (node.DependsOn == nil || len(node.DependsOn) == 0) {
			return fmt.Errorf("Node %s is unexecuted and has no dependencies", node.Name)
		}
		if node.executed {
			continue
		}
		for _, dependency := range node.DependsOn {
			if _, ok := g.nodes[dependency]; !ok {
				return fmt.Errorf("Node %s depends on Node %s which does not exist", node.Name, dependency)
			}
		}
	}
	return nil
}

// CanExecuteNode returns true if the specified node is executeable, if the key is not present or the node cannot
// be executed, it return false
func (g *Graph) CanExecuteNode(name string) bool {
	node, ok := g.nodes[name]
	if !ok {
		return false
	}
	return g.canExecuteNode(node)
}

func (g *Graph) canExecuteNode(node *Unit) bool {
	if node.executed {
		return false
	}
	for _, v := range node.DependsOn {
		subNode := g.nodes[v]
		if !subNode.executed {
			return false
		}
	}
	return true
}

// Size returns the total number of nodes on the graph
func (g *Graph) Size() int {
	return len(g.nodes)
}

// RemainingSize indicates how many nodes need to be executed still
func (g *Graph) RemainingSize() int {
	var counter = 0
	for _, v := range g.nodes {
		if !v.executed {
			counter++
		}
	}
	return counter
}

// IsStuck determines if the graph can continue to execute nodes.
// To check this, IsStick counts how many unexecuted nodes are on the graph and then
// checks if any nodes can be executed. If there are unexecuted nodes but we cannot execute
// any nodes, then we are stuck.
// To ensure this check does not trigger, you should select a target and use GetTargetSubgraph to
// obtain a graph only containing the target's dependencies, direct or indirect.
func (g *Graph) IsStuck() bool {
	if g.RemainingSize() > 0 {
		if len(g.GetAllRunnableNodes()) == 0 {
			return true
		}
	}
	return false
}

// GetAllRunnableNodes returns a list of all nodes that can be executed on the current graph or an empty slice.
func (g *Graph) GetAllRunnableNodes() []string {
	var foundNodes = []string{}
	for _, v := range g.nodes {
		if g.canExecuteNode(v) {
			foundNodes = append(foundNodes, v.Name)
		}
	}
	return foundNodes
}

// GetTargetSubgraph will take a target and create a graph that only contains nodes
// that are direct or indirect dependencies of that target. If a unit is not reachable from
// the current target, it will not be included in the subgraph.
func (g *Graph) GetTargetSubgraph(targetName string) (*Graph, error) {
	target, ok := g.nodes[targetName]
	if !ok {
		return nil, fmt.Errorf("Target %s does not exist", targetName)
	}
	if target.Type != UnitTypeVirtualTarget {
		return nil, fmt.Errorf("Node %s is a Migration not a target", targetName)
	}
	var newGraph = &Graph{nodes: map[string]*Unit{
		target.Name:      target,
		nothingUnit.Name: &nothingUnit,
	}}
	var searchSet = target.DependsOn
	for len(searchSet) > 0 {
		node, ok := g.nodes[searchSet[0]]
		if !ok {
			return nil, fmt.Errorf("Node %s wanted by target but does not exist", searchSet[0])
		}
		if len(searchSet) > 1 {
			searchSet = searchSet[1:]
		} else {
			searchSet = []string{}
		}
		newGraph.nodes[node.Name] = node
		searchSet = append(searchSet, node.DependsOnWithoutNothing()...)
	}
	return newGraph, newGraph.ValidateNodes()
}

// MarkNodesRun will mark a node as executed on the graph, allowing the graph to proceed the execution.
func (g *Graph) MarkNodesRun(names ...string) error {
	for _, name := range names {
		if _, ok := g.nodes[name]; ok {
			g.nodes[name].executed = true
		} else {
			return fmt.Errorf("Node %s not found on graph", name)
		}
	}
	return nil
}
