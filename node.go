package opts

import (
	"fmt"
	"strings"
)

// Node represent a node and its direct dependencies and conflicts
type Node struct {
	Name      string
	Deps      []*Node
	Conflicts []*Node
}

// NewNode create a new node with the name provided and empty Deps and Conflicts *Node arrays
func NewNode(name string) *Node {
	var deps []*Node
	var conflicts []*Node

	return &Node{
		Name:      name,
		Deps:      deps,
		Conflicts: conflicts,
	}
}

// addDep will add a node dependency to a node if the dependency does not exist
func (n *Node) addDep(dep *Node) {

	// check if the node exist otherwise add new node
	for _, v := range n.Deps {
		if v.Name == dep.Name {
			return
		}
	}

	n.Deps = append(n.Deps, dep)

}

// addConflict will add a node conflict to a node if the conflict does not exist
func (n *Node) addConflict(newNode *Node) {

	for _, v := range n.Conflicts {
		if v.Name == newNode.Name {
			return
		}
	}

	n.Conflicts = append(n.Conflicts, newNode)
}

// String is an implementation of the Stringer interface and returns a text containing the
// node name and a printable list of all dependencies and conflicts.
func (n *Node) String() string {
	return fmt.Sprintf("%s-> deps:%s conflicts:%s",
		n.Name,
		n.printAllNodes(n.Deps...),
		n.printAllNodes(n.Conflicts...))
}

func (n *Node) printAllNodes(nodes ...*Node) string {

	var s string
	for _, v := range nodes {
		if strings.TrimSpace(s) == "" {
			s += fmt.Sprintf("'%s'", v.Name)
		} else {
			s += fmt.Sprintf(",'%s'", v.Name)
		}
	}
	return s
}

// isCoherent return true if a node is coherent otherwise false
// by checking whether a node exists in both the Deps and the Conflicts array
func (n *Node) isCoherent() bool {

	for _, d := range n.Deps {
		for _, c := range n.Conflicts {
			if d.Name == c.Name {
				return false
			}
		}
	}

	return true
}
