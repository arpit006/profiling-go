package profiler

import (
	"fmt"
	"runtime"
	"time"
)

func printGCStats() {
	var stats runtime.MemStats
	runtime.ReadMemStats(&stats)

	fmt.Println("[GC] Memory Stats:")
	fmt.Printf("[GC] Allocated Memory: %d bytes\n", stats.Alloc)
	fmt.Printf("[GC] Total Allocated Memory: %d bytes\n", stats.TotalAlloc)
	fmt.Printf("[GC] Heap Memory In Use: %d bytes\n", stats.HeapAlloc)
	fmt.Printf("[GC] Heap Memory Idle: %d bytes\n", stats.HeapIdle)
	fmt.Printf("[GC] Heap Memory Released: %d bytes\n", stats.HeapReleased)
	fmt.Printf("[GC] Number of Garbage Collections: %d\n", stats.NumGC)
	fmt.Printf("[GC] Last GC Pause Time: %s\n", time.Duration(stats.PauseNs[(stats.NumGC+255)%256]))
	fmt.Printf("[GC] GC Pause Times: %s\n", time.Duration(stats.PauseTotalNs/uint64(stats.NumGC)))
}
