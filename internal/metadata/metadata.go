package metadata

type Global struct {
	File string
}

func (g *Global) Load(config string) error {
	return nil
}

func (g *Global) ToDB() error {
	return nil
}

type Master struct {
}

func (m *Master) Load(config string) error {
	return nil
}

type Template struct {
}

func (t *Template) Load(config string) error {
	return nil
}
