package opts

import (
	"fmt"
)

// RuleSet defines a list of nodes with their dependencies and conflicts
// The concept of the data structure is based on predicate logic where
// AddDeb(a,b) means
// a->a, b
// b->b
// and then
// AddDeb(b,c) means
// a->a, b, c
// b->b, c
// c->c
// and then
// AddConflict(b,a) means
// a->a, b, c, ~b
// b->b, c, ~a
// c->c
// In case a node includes the same node in its Dep and Conflict arrays then the RuleSet is
// not coherent
type RuleSet struct {
	isCoherent bool
	nodes      []*Node
}

// NewRuleSet create a new RuleSet
func NewRuleSet() *RuleSet {
	var nodes []*Node
	return &RuleSet{
		isCoherent: true,
		nodes:      nodes,
	}
}

// AddDep is used to add a dependency between 2 options
func (r *RuleSet) AddDep(opt1, opt2 string) {

	opt1Node := r.addNode(opt1)

	if opt1 == opt2 {
		return
	}

	opt2Node := r.addNode(opt2)
	opt1Node.addDep(opt2Node)

	// now lets iterate through the rest of the nodes and add opt2 as a dependency where opt2 exists
	for _, v := range r.nodes {
		if v.Name != opt2 && v.Name != opt1 {
			for _, dep := range v.Deps {
				if dep.Name == opt1 {
					v.addDep(opt2Node)
				}
			}
		}
	}
}

func (r *RuleSet) addNode(opt string) *Node {
	for _, v := range r.nodes {
		if v.Name == opt {
			return v
		}
	}

	newNode := NewNode(opt)
	// we add own self as a dependency as well
	newNode.addDep(newNode)
	r.nodes = append(r.nodes, newNode)
	return newNode
}

// AddConflict defines a new conflict relationship between 2 options.  Every time this is called
// the data structure will re-evaluate its coherence
func (r *RuleSet) AddConflict(opt1, opt2 string) {
	r.addConflict(opt1, opt2)
	r.addConflict(opt2, opt1)
}

func (r *RuleSet) addConflict(opt1, opt2 string) {
	opt1Node := r.addNode(opt1)
	opt2Node := r.addNode(opt2)

	if opt1 == opt2 {
		return
	}

	opt1Node.addConflict(opt2Node)

	// now lets iterate through the rest of the nodes and add opt2 as a dependency where opt2 exists
	for _, v := range r.nodes {
		if v.Name != opt2 && v.Name != opt1 {
			for _, dep := range v.Deps {

				if dep.Name == opt1 {
					v.addConflict(opt2Node)
				}
			}
		}
	}
}

// IsCoherent returns true if the data structure is coherent, that is, that no option can depend,
// directly or indirectly, on another option and also be mutually exclusive with it.
func (r *RuleSet) IsCoherent() bool {

	if r.isCoherent {
		for _, a := range r.nodes {
			r.isCoherent = a.isCoherent()
			if !r.isCoherent {
				break
			}
		}
	}

	return r.isCoherent
}

// String is an implementation of the Stringer interface and prints all the RuleSets nodes
func (r *RuleSet) String() string {
	s := ""
	for _, n := range r.nodes {
		s += fmt.Sprintf("%s\n", n)
	}
	return s
}
