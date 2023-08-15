PROFILING_PKG := github.com/google/pprof

setup-profiler:
	go install $(PROFILING_PKG)@latest
	brew install graphviz

