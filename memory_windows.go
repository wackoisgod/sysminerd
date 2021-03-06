// +build windows

package main

import (
	"time"
	"unsafe"
)

func (m *MemoryInputModule) GetMetrics() ([]Metric, error) {
	metrics := make([]Metric, 0, 50)
	now := time.Now()

	var memData memorystatusex
	memData.dwLength = dword(unsafe.Sizeof(memData))

	r1, _, err := _globalMemoryStatusEx.Call(uintptr(unsafe.Pointer(&memData)))

	if r1 != 1 {
		return nil, err
	}

	metrics = append(metrics, Metric{module: m.Name(), name: "CurrentLoad", value: float64(memData.dwMemoryLoad), timestamp: now})
	metrics = append(metrics, Metric{module: m.Name(), name: "Total", value: float64(memData.ullTotalPhys), timestamp: now})
	metrics = append(metrics, Metric{module: m.Name(), name: "Free", value: float64(memData.ullAvailPhys), timestamp: now})

	used := float64(memData.ullTotalPhys) - float64(memData.ullAvailPhys)

	metrics = append(metrics, Metric{module: m.Name(), name: "Used", value: used, timestamp: now})

	virTotal := float64(memData.ullTotalVirtual)
	virFree := float64(memData.ullAvailVirtual)
	virUsed := virTotal - virFree

	metrics = append(metrics, Metric{module: m.Name(), name: "TotalVirtual", value: virTotal, timestamp: now})
	metrics = append(metrics, Metric{module: m.Name(), name: "FreeVirtual", value: virFree, timestamp: now})
	metrics = append(metrics, Metric{module: m.Name(), name: "UsedVirtual", value: virUsed, timestamp: now})

	return metrics, nil
}
