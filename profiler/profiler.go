package profiler

type Profiler interface {
	Start() error
	Analyse() error
}


