package apiserver

type Options struct {
	Addr string
}

func DefaultOptions() Options {
	return Options{
		Addr: ":8080",
	}
}

func (opt Options) WithAddr(addr string) Options {
	opt.Addr = addr
	return opt
}
