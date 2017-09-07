package util

type LogWriter interface {
	Out(s string)
}

type Debug struct {
	*log.Logger
}

func New(path, label string) *Debug {
	f, e := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	
	if e != nil {
		log.Panic(e)
	}

	flags := 0

	return &Debug{ 
		log.New(f, label, flags),
	}
}