package utils

import (
	"fmt"
	"os"
	"runtime"
	"time"

	"github.com/shirou/gopsutil/process"
)

func MeasureExecutionStats(task func()) {
	// Get process details
	pid := int32(os.Getpid())
	proc, _ := process.NewProcess(pid)

	// Capture memory stats at the start
	var memStart runtime.MemStats
	runtime.ReadMemStats(&memStart)
	startTime := time.Now()

	// Execute the given task
	task()

	// Wait for a brief interval to measure CPU usage accurately
	time.Sleep(1 * time.Second)

	// Capture memory stats at the end
	var memEnd runtime.MemStats
	runtime.ReadMemStats(&memEnd)
	endTime := time.Now()

	// Measure CPU at the end
	cpuUsage, _ := proc.CPUPercent()

	// Results
	fmt.Printf("Elapsed time: %v\n", endTime.Sub(startTime))
	fmt.Printf("CPU usage by process: %.2f%%\n", cpuUsage)
	fmt.Printf("Memory allocated during computation: %v MB\n", (memEnd.TotalAlloc-memStart.TotalAlloc)/(1024*1024))
	fmt.Printf("Memory in use at the end: %v MB\n", memEnd.HeapAlloc/(1024*1024))
}
