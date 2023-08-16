package profiler

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

func getCurrTimeAsStr() string {
	return time.Now().Format(DateFmt)
}

func getFileName(fileNameWithoutExt string) string {
	return fmt.Sprintf("%s-%s%s", fileNameWithoutExt, getCurrTimeAsStr(), ProfFileExt)
}

func getFilePath(fileName string) string {
	return fmt.Sprintf("%s%s", FwdSlash, fileName)
}

func getAbsoluteFilePath(workDir, filePath string) string {
	return filepath.Join(workDir, filePath)
}

func getCurrWorkingDir() (string, error) {
	// Get the absolute path of the directory containing the executable
	wd, err := os.Getwd()
	if err != nil {
		fmt.Printf("[error] in getting current working directory. error is: [%s]\n", err)
		return EmptyStr, errors.Join(NewProfilerErr, fmt.Errorf("[error] in getting current working directory. error is: [%s]", err))
	}
	fmt.Printf("[info] currentWorkingDir: [%s]\n", wd)
	return wd, nil
}

func getCurrDirWritePath(wd string) string {
	return filepath.Join(wd, RootDir)
}

func createDir(wd string) error {
	err := os.MkdirAll(getCurrDirWritePath(wd), os.ModePerm)
	if err != nil {
		fmt.Printf("[error] in creating a dir: [%s]. error is: [%s]\n", RootDir, err)
		return errors.Join(NewProfilerErr, fmt.Errorf("[error] in creating a dir: [%s]. error is: [%s]\n", RootDir, err))
	}
	fmt.Printf("[info] createdDir: [%s]\n", getCurrDirWritePath(wd))
	return nil
}
