package profiler

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"runtime/pprof"
	"time"
)

type asyncMemoryProfiler struct {
	appName            string
	host               string
	portNo             int
	fileNameWithoutExt string
	currFileName       string
	currFilePath       string
	currFile           *os.File
	dirPath            string
	startTime          time.Time
	analysisDuration   time.Duration
	sleepDuration      time.Duration
	iteration          int
	isRunning          bool
}

func NewAsyncMemProfiler(appName, host string, port int, analyseDuration, sleepDuration time.Duration) (Profiler, error) {
	fileNameWithoutExt := fmt.Sprintf("%s", MemProfFile)
	if appName != "" {
		fileNameWithoutExt = fmt.Sprintf("%s-%s", MemProfFile, appName)
	}
	currWorkingDir, err := getCurrWorkingDir()
	if err != nil {
		return nil, err
	}

	fileName := getFileName(fileNameWithoutExt)
	filePath := getFilePath(fileName)

	err = createDir(currWorkingDir)
	if err != nil {
		return nil, err
	}

	absolutePath := getAbsoluteFilePath(getCurrDirWritePath(currWorkingDir), filePath)

	return &asyncMemoryProfiler{
		appName:            appName,
		host:               host,
		portNo:             port,
		fileNameWithoutExt: fileNameWithoutExt,
		currFileName:       fileName,
		currFilePath:       absolutePath,
		currFile:           nil,
		dirPath:            currWorkingDir,
		startTime:          time.Now(),
		analysisDuration:   analyseDuration,
		sleepDuration:      sleepDuration,
		iteration:          0,
		isRunning: true,
	}, nil
}

func (a *asyncMemoryProfiler) Start() error {
	go func() {
		err := a.startInBackground()
		if err != nil {
			fmt.Printf("[error] in profiling memory async")
			return
		}
	}()
	return nil
}

func (a *asyncMemoryProfiler) startInBackground() error {
	f, err := os.Create(a.currFilePath)
	if err != nil {
		fmt.Printf("[error] error in creating currFile. error is: [%s]\n", err)
		return errors.Join(FileErr, fmt.Errorf("error in creating currFile. error is: [%s]\n", err))
	}
	a.currFile = f
	a.iteration = a.iteration + 1

	for a.shouldRunFurther() {
		time.Sleep(a.sleepDuration)
		err := a.rotate()
		if err != nil {
			return fmt.Errorf("error while rotating for memory profiler. error is: [%s]", err)
		}
	}
	return nil
}

func (a *asyncMemoryProfiler) rotate() error {
	err := pprof.WriteHeapProfile(a.currFile)
	if err != nil {
		return fmt.Errorf("error in writing heap profiler. error is: [%s]", err)
	}
	err = a.currFile.Close()
	if err != nil {
		return fmt.Errorf("error in closing profiling file. error is: [%s]", err)
	}

	// close old instances
	a.iteration = a.iteration + 1
	a.currFileName = getFileName(a.fileNameWithoutExt)
	wd, err := getCurrWorkingDir()
	if err != nil { return err }
	writeDir := getCurrDirWritePath(wd)
	a.currFilePath = getAbsoluteFilePath(writeDir, getFilePath(a.currFileName))

	// create new instances
	f, err := os.Create(a.currFilePath)
	if err != nil {
		fmt.Printf("[error] error in creating currFile. error is: [%s]\n", err)
		return errors.Join(FileErr, fmt.Errorf("error in creating currFile. error is: [%s]\n", err))
	}
	a.currFile = f
	fmt.Printf("******* rotating profiler [%s]\n", a.currFileName)
	return nil
}

func (a *asyncMemoryProfiler) Stop() error {
	err := pprof.WriteHeapProfile(a.currFile)
	if err != nil {
		return fmt.Errorf("error in writing heap profiler. error is: [%s]", err)
	}
	err = a.currFile.Close()
	if err != nil {
		return fmt.Errorf("error in closing profiling file. error is: [%s]", err)
	}
	a.isRunning = false
	return nil
}

func (a *asyncMemoryProfiler) Analyse() error {
	wd, err := getCurrWorkingDir()
	if err != nil { return err }
	go func() {
		err := a.analyseInBackground(wd)
		if err != nil {
			fmt.Printf("[error] in analysing the memory profiler in background")
			return
		}
	}()
	return nil
}

func (a *asyncMemoryProfiler) analyseInBackground(wd string) error {

	cmd := exec.Command("pprof", "-http", fmt.Sprintf("%s:%d", a.host, a.portNo), fmt.Sprintf("%s%s%s", getCurrDirWritePath(wd), FwdSlash, a.fileNameWithoutExt+"*"))

	// set any env variables if required
	cmd.Env = append(os.Environ(), "ENV=local")

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		fmt.Printf("[error] in analysing the memory-profiler for: [%s]. error is [%s]\n", a.appName, err)
		return errors.Join(ProfileAnalyseErr, err)
	}
	return nil
}

func (a *asyncMemoryProfiler) shouldRunFurther() bool {
	if time.Since(a.startTime).Milliseconds() > a.analysisDuration.Milliseconds() {
		return false
	}
	return true
}
