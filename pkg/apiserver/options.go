package apiserver

type Options struct {
	Addr       string
	ConfigPath string
}

func DefaultOptions() Options {
	return Options{
		Addr:       ":8080",
		ConfigPath: "./config.yaml",
	}
}

func (opt Options) WithAddr(addr string) Options {
	opt.Addr = addr
	return opt
}
func (opt Options) WithConfigPath(configPath string) Options {
	opt.ConfigPath = configPath
	return opt
}
