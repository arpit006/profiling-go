package profiler

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime/pprof"
)

type memoryProfiler struct {
	appName     string
	fileName    string
	host        string
	portNo      int
	getFilePath func() string
	getAbsoluteFilePath func() string
	file        *os.File
}

func NewMemoryProfiler(appName, host string, port int) (Profiler, error) {
	fileName := fmt.Sprintf("%s%s", MemProfFile, ProfFileExt)
	if appName != "" {
		fileName = fmt.Sprintf("%s-%s%s", MemProfFile, appName, ProfFileExt)
	}

	// Get the absolute path of the directory containing the executable
	wd, err := os.Getwd()
	if err != nil {
		fmt.Printf("[error] in getting current working directory. error is: [%s]\n", err)
		return nil, errors.Join(NewProfilerErr, fmt.Errorf("[error] in getting current working directory. error is: [%s]", err))
	}

	filePath := fmt.Sprintf("%s%s%s", RootDir, FwdSlash, fileName)
	absolutePath := filepath.Join(wd, filePath)

	err = os.MkdirAll(filepath.Join(wd, RootDir), os.ModePerm)
	if err != nil {
		fmt.Printf("[error] in creating a dir: [%s]. error is: [%s]\n", RootDir, err)
		return nil, errors.Join(NewProfilerErr, fmt.Errorf("[error] in creating a dir: [%s]. error is: [%s]\n", RootDir, err))
	}
	return &memoryProfiler{
		appName:  appName,
		fileName: fileName,
		portNo:   port,
		host:     host,
		getFilePath: func() string {
			return filePath
		},
		getAbsoluteFilePath: func() string {
			return absolutePath
		},
	}, nil
}

func (mem *memoryProfiler) Start() error {
	f, err := os.Create(mem.getAbsoluteFilePath())
	if err != nil {
		fmt.Printf("[error] error in creating file. error is: [%s]\n", err)
		return errors.Join(FileErr, fmt.Errorf("error in creating file. error is: [%s]\n", err))
	}
	mem.file = f
	return nil
}

func (mem *memoryProfiler) Stop() error {
	err := pprof.WriteHeapProfile(mem.file)
	if err != nil {
		return fmt.Errorf("error in starting CPU profiler. error is: [%s]\n", err)
	}

	err = mem.file.Close()
	if err != nil {
		fmt.Printf("[error] in closing the file: [%s]. error is: [%s]", mem.getFilePath(), err)
		return errors.Join(FileErr, err)
	}
	return nil
}

func (mem *memoryProfiler) Analyse() error {
	cmd := exec.Command("pprof", "-http", fmt.Sprintf("%s:%d", mem.host, mem.portNo), fmt.Sprintf("%s", mem.getFilePath()))

	// set any env variables if required
	cmd.Env = append(os.Environ(), "ENV=local")

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		fmt.Printf("[error] in analysing the memory-profiler for: [%s]. error is [%s]\n", mem.appName, err)
		return errors.Join(ProfileAnalyseErr, err)
	}
	return nil
}
