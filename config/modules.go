package config

type Module struct {
	Variables map[string]*Variable `json:"variables"`
}

func newModule() *Module {
	return &Module{
		Variables: make(map[string]*Variable),
	}
}
