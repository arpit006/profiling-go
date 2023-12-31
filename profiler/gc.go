package profiler

import (
	"fmt"
	"runtime"
	"time"
)

func printGCStats(prefixStr string) {
	var stats runtime.MemStats
	runtime.ReadMemStats(&stats)

	fmt.Printf("[GC] [%s] Memory Stats:\n", prefixStr)
	fmt.Printf("[GC] [%s] Allocated Memory: %d bytes\n", prefixStr, stats.Alloc)
	fmt.Printf("[GC] [%s] Total Allocated Memory: %d bytes\n", prefixStr, stats.TotalAlloc)
	fmt.Printf("[GC] [%s] Frees count from heap: %d bytes\n", prefixStr, stats.Frees)
	fmt.Printf("[GC] [%s] Next GC target size: %d \n", prefixStr, stats.Frees)
	fmt.Printf("[GC] [%s] Last GC run at: %v time\n", prefixStr, time.Unix(0, int64(stats.LastGC)))
	fmt.Printf("[GC] [%s] Total GC pause duration: %v time\n", prefixStr, time.Duration(int64(stats.LastGC)).Minutes())
	fmt.Printf("[GC] [%s] Heap Memory In Use: %d bytes\n", prefixStr, stats.HeapAlloc)
	fmt.Printf("[GC] [%s] Heap Memory Idle: %d bytes\n", prefixStr, stats.HeapIdle)
	fmt.Printf("[GC] [%s] Heap Memory Released: %d bytes\n", prefixStr, stats.HeapReleased)
	fmt.Printf("[GC] [%s] Number of Garbage Collections: %d\n", prefixStr, stats.NumGC)
	fmt.Printf("[GC] [%s] Last GC Pause Time: %s\n", prefixStr, time.Duration(stats.PauseNs[(stats.NumGC+255)%256]))
	if stats.NumGC > 0 {
		fmt.Printf("[GC] [%s] GC Pause Times: %s\n", prefixStr, time.Duration(stats.PauseTotalNs/uint64(stats.NumGC)))
	}
}
