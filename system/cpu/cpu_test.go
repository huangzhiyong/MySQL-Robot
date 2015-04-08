package cpu

import (
	"fmt"
	"testing"
)

func TestCPUTimes(t *testing.T) {
	v, err := CPUTimes(false)
	if err != nil {
		t.Errorf("error: %v", err)
	}

	if len(v) == 0 {
		t.Error("could not get CPUs:", err)
	}

	nodata := CPUTimesStat{}
	for _, vv := range v {
		if vv == nodata {
			t.Errorf("could not get CPU stat: %v", vv)
		}
		fmt.Println(vv)
	}
}

func TestCPUTimesAll(t *testing.T) {
	v, err := CPUTimes(true)
	if err != nil {
		t.Errorf("error: %v", err)
	}

	if len(v) == 0 {
		t.Errorf("could not get CPUs:", err)
	}

	nodata := CPUTimesStat{}
	for _, vv := range v {
		if vv == nodata {
			t.Errorf("could not get CPU stat: %v", vv)
		}
		fmt.Println(vv)
	}
}

func TestCPUInfo(t *testing.T) {
	v, err := NewCPUInfo()
	if err != nil {
		t.Error(err)
	}
	if len(v.Processors) == 0 {
		t.Errorf("could not get CPU infos:", err)
	}
	for _, vv := range v.Processors {
		fmt.Println(vv)
	}
	fmt.Println(v.NumCore(false))
	fmt.Println(v.NumCPU())
}
