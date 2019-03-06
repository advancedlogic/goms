package rest

import "errors"

type Rest struct {
	port int
}

type Builder struct {
	*Rest
}

func NewBuilder() *Builder {
	return &Builder{
		Rest: &Rest{},
	}
}

func (sb *Builder) Port(port int) *Builder {
	sb.port = port
	return sb
}

func (sb *Builder) Build() (*Rest, error) {
	if sb.Rest != nil {
		return sb.Rest, nil
	}
	return nil, errors.New("")
}

func (r *Rest) Run(opts ...string) error {
	return nil
}
