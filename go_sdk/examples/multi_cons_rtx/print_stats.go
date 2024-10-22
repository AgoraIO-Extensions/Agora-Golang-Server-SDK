package main

import (
	"fmt"
	"runtime"
)

func bToMb(b uint64) uint64 {
	return b / 1024 / 1024
}

func printStats() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	numCPU := runtime.NumCPU()

	fmt.Printf("stats: cpus %d, alloc %v MiB, sys %v MiB, num gc %v\n",
		numCPU, bToMb(m.Alloc), bToMb(m.Sys), m.NumGC)
}
