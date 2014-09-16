package server

type Config struct {
	Port int
	Filename string
	Testmode bool
}

func NewConfig(portByFlag int, filename string) Config {

	if len(filename) == 0 {
		return Config{portByFlag, filename, false}
	} else {
		return Config{portByFlag, filename, true}
	}
}
