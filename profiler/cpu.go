package profiler

import (
	"fmt"
	"os"
	"os/exec"
	"runtime/pprof"
)

type cpuProfiler struct {
	appName     string
	fileName    string
	host        string
	portNo      int
	getFilePath func() string
	file        *os.File
}

func NewCPUProfiler(appName, host string, port int) Profiler {
	fileName := fmt.Sprintf("%s%s", CPU_PROF, PROF_EXT)
	if appName != "" {
		fileName = fmt.Sprintf("%s-%s%s", CPU_PROF, appName, PROF_EXT)
	}
	return &cpuProfiler{
		appName:  appName,
		fileName: fileName,
		portNo:   port,
		host:     host,
		getFilePath: func() string {
			return fmt.Sprintf("%s%s%s", ROOT_DIR, FWD_SLASH, fileName)
		},
	}
}

func (cpu *cpuProfiler) Start() error {
	f, err := os.Create(cpu.getFilePath())
	if err != nil {
		fmt.Printf("[error] error in creating file. error is: [%s]\n", err)
		return fmt.Errorf("error in creating file. error is: [%s]\n", err)
	}

	cpu.file = f

	err = pprof.StartCPUProfile(cpu.file)
	if err != nil {
		return fmt.Errorf("error in starting CPU profiler. error is: [%s]\n", err)
	}

	return nil
}

func (cpu *cpuProfiler) Analyse() error {
	pprof.StopCPUProfile()

	err := cpu.file.Close()
	if err != nil {
		fmt.Printf("[error] in closing the file: [%s]. error is: [%s]", cpu.getFilePath(), err)
		return err
	}

	cmd := exec.Command("pprof", "-http", fmt.Sprintf("%s:%d", cpu.host, cpu.portNo), fmt.Sprintf("%s", cpu.getFilePath()))

	// set any env variables if required
	cmd.Env = append(os.Environ(), "ENV=local")

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err = cmd.Run()
	if err != nil {
		fmt.Printf("[error] in analysing the cpu-profiler for: [%s]. error is [%s]\n", cpu.appName, err)
		return err
	}
	return nil
}
