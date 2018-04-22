package opts

import "fmt"

// Opts defines the structure that toggles on/off options
type Opts struct {
	ruleSet *RuleSet
	options map[string]bool
	allOff  bool
}

// New creates an Opts with a RuleSet
func New(rs *RuleSet) *Opts {
	opts := Opts{ruleSet: rs,
		allOff:  true,
		options: make(map[string]bool)}

	opts.refreshOptions()
	return &opts
}

func (r *Opts) Toggle(opt1 string) {

	if !r.ruleSet.IsCoherent() {
		return
	}

	r.refreshOptions()
	isTurnedOn := r.isTurnedOn(opt1)
	node, err := r.getNode(opt1)

	if err != nil {
		return
	}

	// this means that the option is on so we need to turn it off
	if isTurnedOn {
		r.turnOffRecursively(opt1)
	}

	// this means that the option is off but there are conflicts
	// so just turn on all dependent options first
	// and then turn off all conflict recursively
	if !isTurnedOn {
		for _, o := range node.Deps {
			r.turnOn(o.Name)
		}

		for _, o := range node.Conflicts {
			r.turnOffRecursively(o.Name)
		}
	}

}

// StringSlice returns a list of options that are turned on
func (r *Opts) StringSlice() (out []string) {
	for k, n := range r.options {
		if n {
			out = append(out, k)
		}
	}
	return
}

// refreshOptions iterates through the RuleSet nodes and adds them into the Opts options
// map in case they do not exist
func (r *Opts) refreshOptions() {
	for _, n := range r.ruleSet.nodes {
		if _, ok := r.options[n.Name]; !ok {
			r.options[n.Name] = false
		}
	}

	r.allOff = true

	for _, n := range r.options {
		if n {
			r.allOff = false
			break
		}
	}

}

// getNodesDependentOn filters all nodes that depend on a node with name defined by the opt param
// and the list will not include a node that has the name with the opt param
func (r *Opts) getNodesDependentOn(opt string, includeSef bool) (out []*Node) {

	for _, n := range r.ruleSet.nodes {
		if n.Name != opt && !includeSef {
			for _, d := range n.Deps {
				if d.Name == opt {
					out = append(out, n)
				}
			}
		}
		if includeSef {
			for _, d := range n.Deps {
				if d.Name == opt {
					out = append(out, n)
				}
			}
		}

	}

	return
}

func (r *Opts) getNode(opt string) (*Node, error) {
	for _, n := range r.ruleSet.nodes {
		if n.Name == opt {
			return n, nil
		}
	}

	return nil, fmt.Errorf("could not find node '%s'", opt)
}

func (r *Opts) turnOn(opt string) error {
	if n, ok := r.options[opt]; ok {
		if !n {
			r.options[opt] = true
		}
		return nil
	}

	return fmt.Errorf("could not find option '%s' to turn it on", opt)
}

func (r *Opts) turnOff(opt string) error {
	if n, ok := r.options[opt]; ok {
		if n {
			r.options[opt] = false
		}
		return nil
	}

	return fmt.Errorf("could not find option '%s' to turn it on", opt)
}

func (r *Opts) turnOffRecursively(opt string) {
	if !r.isTurnedOn(opt) {
		return
	}

	node, err := r.getNode(opt)

	if err != nil {
		return
	}

	if err := r.turnOff(node.Name); err != nil {
		return
	}

	dependentNodes := r.getNodesDependentOn(opt, false)

	for _, n := range dependentNodes {
		if r.isTurnedOn(n.Name) {
			r.turnOffRecursively(n.Name)
		}
	}
}

func (r *Opts) isTurnedOn(opt string) bool {
	if n, ok := r.options[opt]; ok {
		return n
	}

	return false
}
