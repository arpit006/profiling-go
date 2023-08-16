package profiler

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime/pprof"
)

type cpuProfiler struct {
	appName             string
	fileName            string
	host                string
	portNo              int
	getFilePath         func() string
	getAbsoluteFilePath func() string
	file                *os.File
}

func NewCPUProfiler(appName, host string, port int) (Profiler, error) {
	fileName := fmt.Sprintf("%s%s", CpuProfFile, ProfFileExt)
	if appName != "" {
		fileName = fmt.Sprintf("%s-%s%s", CpuProfFile, appName, ProfFileExt)
	}

	// Get the absolute path of the directory containing the executable
	wd, err := os.Getwd()
	if err != nil {
		fmt.Printf("[error] in getting current working directory. error is: [%s]", err)
		return nil, NewProfilerErr
	}

	filePath := fmt.Sprintf("%s%s%s", RootDir, FwdSlash, fileName)
	absolutePath := filepath.Join(wd, filePath)

	err = os.MkdirAll(filepath.Join(wd, RootDir), os.ModePerm)
	if err != nil {
		fmt.Printf("[error] in creating a dir: [%s]. error is: [%s]\n", RootDir, err)
		return nil, errors.Join(NewProfilerErr, fmt.Errorf("[error] in creating a dir: [%s]. error is: [%s]", RootDir, err))
	}

	return &cpuProfiler{
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

func (cpu *cpuProfiler) Start() error {
	f, err := os.Create(cpu.getAbsoluteFilePath())
	if err != nil {
		fmt.Printf("[error] error in creating currFile. error is: [%s]\n", err)
		return errors.Join(FileErr, fmt.Errorf("error in creating currFile. error is: [%s]\n", err))
	}

	cpu.file = f

	err = pprof.StartCPUProfile(cpu.file)
	if err != nil {
		fmt.Printf("error in starting CPU profiler. error is: [%s]\n", err)
		return errors.Join(ProfilingErr, err)
	}

	return nil
}

func (cpu *cpuProfiler) Stop() error {
	pprof.StopCPUProfile()

	err := cpu.file.Close()
	if err != nil {
		fmt.Printf("[error] in closing the currFile: [%s]. error is: [%s]", cpu.getFilePath(), err)
		return errors.Join(FileErr, err)
	}
	return nil
}

func (cpu *cpuProfiler) Analyse() error {
	cmd := exec.Command("pprof", "-http", fmt.Sprintf("%s:%d", cpu.host, cpu.portNo), fmt.Sprintf("%s", cpu.getFilePath()))

	// set any env variables if required
	cmd.Env = append(os.Environ(), "ENV=local")

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		fmt.Printf("[error] in analysing the cpu-profiler for: [%s]. error is [%s]\n", cpu.appName, err)
		return errors.Join(ProfileAnalyseErr, err)
	}
	return nil
}
