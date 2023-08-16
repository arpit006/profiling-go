package main

import (
	"fmt"
	"github.com/arpit006/profiling-go/profiler"
	"math/rand"
	"strconv"
	"time"
)

var (
	m = map[string]string{}
)

func main() {
	//analyseCpu()
	//analyseMem()

	analyseAsyncMem()
}

func analyseCpu() {
	cpu, _ := profiler.NewCPUProfiler("test", "localhost", 8001)

	err := cpu.Start()
	if err != nil {
		fmt.Printf("[error] in cpu profiling. error is: [%s]", err)
		return
	}
	work()

	cpu.Stop()
	err = cpu.Analyse()
	if err != nil {
		fmt.Printf("[error] in analysing. error is: [%s]", err)
	}
}

func analyseMem() {
	mem, _ := profiler.NewMemoryProfiler("test", "localhost", 8001)

	err := mem.Start()
	if err != nil {
		fmt.Printf("[error] in mem profiling. error is: [%s]", err)
		return
	}

	work()

	mem.Stop()
	err = mem.Analyse()
	if err != nil {
		fmt.Printf("[error] in analysing. error is: [%s]", err)
	}
}

func analyseAsyncMem() {
	mem, _ := profiler.NewAsyncMemProfiler("test", "localhost", 8001, 2 * time.Minute, 30 * time.Second)

	err := mem.Start()
	if err != nil {
		fmt.Printf("[error] in mem profiling. error is: [%s]", err)
		return
	}

	work()

	//mem.Stop()
	time.Sleep(5 * time.Minute)
}

func work() {
	i := 1
	for  {
		//time.Sleep(10 * time.Millisecond)
		//fmt.Printf("i: [%d]\n", i)
		r := rand.Intn(i)
		key := fmt.Sprintf("%d:%d", i, r)
		m[key] = strconv.Itoa(i)
		i = i + 1
	}
}
