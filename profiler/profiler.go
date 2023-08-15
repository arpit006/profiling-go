package profiler

type Profiler interface {
	Start() error
	Stop() error
	Analyse() error
}


