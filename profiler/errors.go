package profiler

import "fmt"

var (
	NewProfilerErr = fmt.Errorf("error in creating Profiler")
	FileErr = fmt.Errorf("error in currFile Handling")
	ProfilingErr = fmt.Errorf("error in profiling")

	ProfileAnalyseErr = fmt.Errorf("error in analysing profile")
)
