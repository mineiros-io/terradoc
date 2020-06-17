package config

// Module holds information of a HCL structure
// TODO: we currently only support variable blocks
type Module struct {
	Variables map[string]*Variable `json:"variables"`
}

func newModule() *Module {
	return &Module{
		Variables: make(map[string]*Variable),
	}
}
