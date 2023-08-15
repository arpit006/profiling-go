package profiler

import (
	"fmt"
	"os"
	"os/exec"
	"runtime/pprof"
)

type memoryProfiler struct {
	appName     string
	fileName    string
	host        string
	portNo      int
	getFilePath func() string
	file        *os.File
}

func NewMemoryProfiler(appName, host string, port int) Profiler {
	fileName := fmt.Sprintf("%s%s", MEM_PROF, PROF_EXT)
	if appName != "" {
		fileName = fmt.Sprintf("%s-%s%s", MEM_PROF, appName, PROF_EXT)
	}
	return &memoryProfiler{
		appName:  appName,
		fileName: fileName,
		portNo:   port,
		host:     host,
		getFilePath: func() string {
			return fmt.Sprintf("%s%s%s", ROOT_DIR, FWD_SLASH, fileName)
		},
	}
}

func (mem *memoryProfiler) Start() error {
	f, err := os.Create(mem.getFilePath())
	if err != nil {
		fmt.Printf("[error] error in creating file. error is: [%s]\n", err)
		return fmt.Errorf("error in creating file. error is: [%s]\n", err)
	}
	mem.file = f
	return nil
}

func (mem *memoryProfiler) Analyse() error {
	err := pprof.WriteHeapProfile(mem.file)
	if err != nil {
		return fmt.Errorf("error in starting CPU profiler. error is: [%s]\n", err)
	}

	err = mem.file.Close()
	if err != nil {
		fmt.Printf("[error] in closing the file: [%s]. error is: [%s]", mem.getFilePath(), err)
		return err
	}

	cmd := exec.Command("pprof", "-http", fmt.Sprintf("%s:%d", mem.host, mem.portNo), fmt.Sprintf("%s", mem.getFilePath()))

	// set any env variables if required
	cmd.Env = append(os.Environ(), "ENV=local")

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err = cmd.Run()
	if err != nil {
		fmt.Printf("[error] in analysing the memory-profiler for: [%s]. error is [%s]\n", mem.appName, err)
		return err
	}
	return nil
}
