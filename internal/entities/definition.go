package entities

// TODO: rename
type Definition struct {
	Sections []*Section `json:"sections"`
}

func (d *Definition) AllVariables() (vars []*Variable) {
	// TODO: refactor
	for _, sec := range d.Sections {
		for _, subsec := range sec.Sections {
			vars = append(vars, subsec.Variables...)
		}
	}

	return vars
}
