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

type Graph struct {
	nodes map[string]*Unit
}

func NewGraph() *Graph {
	return &Graph{
		nodes: map[string]*Unit{
			"nothing": &nothingUnit,
		},
	}
}

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

func (g *Graph) GetUnit(name string) (Unit, error) {
	if node, ok := g.nodes[name]; ok {
		return *node, nil
	}
	return Unit{}, errors.New("Node not found")
}

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

// Return the number of existing nodes in the graph
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

// If the graph contains unexecuted nodes but cannot execute any nodes then
// we are stuck. It means we have nodes with either a cyclic dependency
// or unresolvable dependencies, both of which are critical failures.
func (g *Graph) IsStuck() bool {
	if g.RemainingSize() > 0 {
		if len(g.GetAllRunnableNodes()) == 0 {
			return true
		}
	}
	return false
}

func (g *Graph) GetAllRunnableNodes() []string {
	var foundNodes = []string{}
	for _, v := range g.nodes {
		if g.canExecuteNode(v) {
			foundNodes = append(foundNodes, v.Name)
		}
	}
	return foundNodes
}

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
		if len(searchSet) > 1 {
			searchSet = searchSet[1:]
		} else {
			searchSet = []string{}
		}
		if !ok {
			return nil, fmt.Errorf("Node %s wanted by target but does not exist", searchSet[0])
		}
		newGraph.nodes[node.Name] = node
		searchSet = append(searchSet, node.DependsOnWithoutNothing()...)
	}
	return newGraph, newGraph.ValidateNodes()
}

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
