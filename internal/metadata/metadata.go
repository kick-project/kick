package metadata

type Global struct {
	File   string
	DBFile string
}

func (g *Global) Load() error {
	return nil
}

func (g *Global) ToDB() error {
	return nil
}

type Master struct {
}

func (g *Master) Load() error {
	return nil
}

func (g *Master) ToDB() error {
	return nil
}

type Template struct {
}
