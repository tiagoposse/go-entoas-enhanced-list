package filter

import "github.com/ogen-go/ogen"

type Opt struct {
	Name string
	In string
}

func (o *Opt) Set(param *ogen.Parameter) {
	param.In = o.In
	param.Name = o.Name
}


type OptConfig func(*Opt)
func In(in string) OptConfig {
	return func(o *Opt) {
		o.In = in
	}
}

func Name(name string) OptConfig {
	return func(o *Opt) {
		o.Name = name
	}
}
