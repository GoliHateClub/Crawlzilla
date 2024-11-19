package utils

import (
	"fmt"
	"os"
	"runtime"
	"time"

	"github.com/shirou/gopsutil/process"
)

func MeasureExecutionStats(task func()) string {
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

	return fmt.Sprintf("Elapsed time: %v\nCPU usage by process: %.2f%%\nMemory allocated during computation: %v MB\nMemory in use at the end: %v MB\n", endTime.Sub(startTime), cpuUsage, (memEnd.TotalAlloc-memStart.TotalAlloc)/(1024*1024), memEnd.HeapAlloc/(1024*1024))
}
