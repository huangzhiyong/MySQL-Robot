package cpu

import (
	"encoding/json"
	"errors"
	common "github.com/huangzhiyong/MySQL-Robot/common"
	"runtime"
	"strconv"
	"strings"
)

const (
	STAT_FILE     = "/proc/stat"
	CPU_INFO_FILE = "/proc/cpuinfo"
	CPU_TICK      = 100
)

type CPUTimesStat struct {
	CPU       string  `json:"cpu"`
	User      float64 `json:"user"`
	Nice      float64 `json:"nice"`
	Sys       float64 `json:"sys"`
	Idle      float64 `json:"idle"`
	Iowait    float64 `json:"iowait"`
	Irq       float64 `json:"irq"`
	Softirq   float64 `json:"softirq"`
	Steal     float64 `json:"steal"`
	Guest     float64 `json:"guest"`
	GuestNice float64 `json:"guest_nice"`
}

func CPUCounts() int {
	return runtime.NumCPU()
}

func (c CPUTimesStat) String() string {
	v := []string{
		`"cpu":"` + c.CPU + `"`,
		`"user":"` + strconv.FormatFloat(c.User, 'f', 2, 64) + `"`,
		`"nice":"` + strconv.FormatFloat(c.Nice, 'f', 2, 64) + `"`,
		`"sys":"` + strconv.FormatFloat(c.Sys, 'f', 2, 64) + `"`,
		`"idle":"` + strconv.FormatFloat(c.Idle, 'f', 2, 64) + `"`,
		`"iowait":"` + strconv.FormatFloat(c.Iowait, 'f', 2, 64) + `"`,
		`"irq":"` + strconv.FormatFloat(c.Irq, 'f', 2, 64) + `"`,
		`"softirq":"` + strconv.FormatFloat(c.Softirq, 'f', 2, 64) + `"`,
		`"steal":"` + strconv.FormatFloat(c.Steal, 'f', 2, 64) + `"`,
		`"guest":"` + strconv.FormatFloat(c.Guest, 'f', 2, 64) + `"`,
		`"guest_nice":"` + strconv.FormatFloat(c.GuestNice, 'f', 2, 64) + `"`,
	}

	return `{` + strings.Join(v, ",") + `}`
}

func parseCPUStatLine(line string) (*CPUTimesStat, error) {
	fields := strings.Fields(line)

	if !strings.HasPrefix(fields[0], "cpu") {
		return nil, errors.New("not contain cpu")
	}

	cpu := fields[0]
	if cpu == "cpu" {
		cpu = "all"
	}
	user, err := strconv.ParseFloat(fields[1], 64)
	if err != nil {
		return nil, err
	}
	nice, err := strconv.ParseFloat(fields[2], 64)
	if err != nil {
		return nil, err
	}

	sys, err := strconv.ParseFloat(fields[3], 64)
	if err != nil {
		return nil, err
	}

	idle, err := strconv.ParseFloat(fields[4], 64)
	if err != nil {
		return nil, err
	}
	iowait, err := strconv.ParseFloat(fields[5], 64)
	if err != nil {
		return nil, err
	}
	irq, err := strconv.ParseFloat(fields[6], 64)
	if err != nil {
		return nil, err
	}
	softirq, err := strconv.ParseFloat(fields[7], 64)
	if err != nil {
		return nil, err
	}

	res := &CPUTimesStat{
		CPU:     cpu,
		User:    user / CPU_TICK,
		Nice:    nice / CPU_TICK,
		Sys:     sys / CPU_TICK,
		Idle:    idle / CPU_TICK,
		Iowait:  iowait / CPU_TICK,
		Irq:     irq / CPU_TICK,
		Softirq: softirq / CPU_TICK,
	}

	vlen := len(fields)
	if vlen >= 8 { // Linux kernel >= 2.6.11
		steal, err := strconv.ParseFloat(fields[8], 64)
		if err != nil {
			return nil, err
		}
		res.Steal = steal
	}
	if vlen >= 9 { // Linux kernel >= 2.6.24
		guest, err := strconv.ParseFloat(fields[9], 64)
		if err != nil {
			return nil, err
		}
		res.Guest = guest
	}
	if vlen >= 10 { // Linux kernel >= 3.2.0
		guest_nice, err := strconv.ParseFloat(fields[10], 64)
		if err != nil {
			return nil, err
		}
		res.GuestNice = guest_nice
	}

	return res, nil
}

func CPUTimes(percpu bool) ([]CPUTimesStat, error) {
	lines := make([]string, 0)
	if percpu {
		contents, _ := common.ReadLinesAll(STAT_FILE)
		for _, line := range contents {
			if !strings.HasPrefix(line, "cpu") {
				break
			}
			lines = append(lines, line)
		}
	} else {
		lines, _ = common.ReadLinesOffset(STAT_FILE, 0, 1)
	}

	res := make([]CPUTimesStat, 0, len(lines))

	for _, line := range lines {
		ct, err := parseCPUStatLine(line)
		if err != nil {
			continue
		}
		res = append(res, *ct)
	}
	return res, nil
}

type CPUInfo struct {
	Processors []CPUInfoStat `json:"processors"`
}

type CPUInfoStat struct {
	Processor       int32    `json:"processor"`
	VendorId        string   `json:"vendor_id"`
	CpuFamily       int32    `json:"cpu_family"`
	Model           int32    `json:"model"`
	ModelName       string   `json:"model_name"`
	Stepping        int32    `json:"stepping"`
	CpuMHz          float64  `json:"cpu_mhz"`
	CacheSize       int32    `json:"cache_size"`
	PhysicalId      int32    `json:"physical_id"`
	Silbings        int32    `json:"silbings"`
	CoreId          int32    `json:"core_id"`
	CpuCores        int32    `json:"cpu_cores"`
	Apicid          int32    `json:"apicid"`
	InitApicid      int32    `json:"init_apicid"`
	Fpu             string   `json:"fpu"`
	FpuException    string   `json:"fpu_exception"`
	CpuidLevel      int32    `json:"cpuid_level"`
	Wp              string   `json:"wp"`
	Flags           []string `json:"flags"`
	Bogomips        float64  `json:"bogomips"`
	ClflushSize     int32    `json:"clflush_size"`
	CacheAlignment  int32    `json:"cache_alignment"`
	AddressSizes    string   `json:"address_sizes"`
	PowerManagement string   `json:"power_management"`
}

func (c CPUInfo) NumCPU() int {
	return len(c.Processors)
}

func (c CPUInfo) NumCore(phycpu bool) int {
	core := make(map[string]bool)

	for _, p := range c.Processors {
		p_id := p.PhysicalId
		c_id := p.CoreId
		var key string
		if phycpu {
			key = strconv.FormatInt(int64(p_id), 10)
		} else {
			key = strconv.FormatInt(int64(p_id), 10) + "-" + strconv.FormatInt(int64(c_id), 10)
		}
		core[key] = true
	}

	return len(core)
}

func (c CPUInfoStat) string() string {
	s, _ := json.Marshal(c)
	return string(s)
}

func NewCPUInfo() (*CPUInfo, error) {
	contents, _ := common.ReadLinesAll(CPU_INFO_FILE)

	res := &CPUInfo{}

	var info CPUInfoStat
	for _, line := range contents {
		fields := strings.Split(line, ":")
		if len(fields) < 2 {
			if info.VendorId != "" {
				res.Processors = append(res.Processors, info)
			}
			continue
		}

		key := strings.TrimSpace(fields[0])
		value := strings.TrimSpace(fields[1])
		switch key {
		case "processor":
			info = CPUInfoStat{}
			v, err := strconv.ParseInt(value, 10, 64)
			if err != nil {
				return nil, err
			}
			info.Processor = int32(v)
		case "vendor_id":
			info.VendorId = value
		case "cpu family":
			v, err := strconv.ParseInt(value, 10, 64)
			if err != nil {
				return nil, err
			}
			info.CpuFamily = int32(v)
		case "model":
			v, err := strconv.ParseInt(value, 10, 64)
			if err != nil {
				return nil, err
			}
			info.Model = int32(v)
		case "model name":
			info.ModelName = value
		case "stepping":
			v, err := strconv.ParseInt(value, 10, 64)
			if err != nil {
				return nil, err
			}
			info.Stepping = int32(v)
		case "cpu MHz":
			v, err := strconv.ParseFloat(value, 64)
			if err != nil {
				return nil, err
			}
			info.CpuMHz = v
		case "cache size":
			v, err := strconv.ParseInt(strings.Replace(value, " KB", "", 1), 10, 64)
			if err != nil {
				return nil, err
			}
			info.CacheSize = int32(v)
		case "physical id":
			v, err := strconv.ParseInt(value, 10, 64)
			if err != nil {
				return nil, err
			}
			info.PhysicalId = int32(v)
		case "silbings":
			v, err := strconv.ParseInt(value, 10, 64)
			if err != nil {
				return nil, err
			}
			info.Silbings = int32(v)
		case "core id":
			v, err := strconv.ParseInt(value, 10, 64)
			if err != nil {
				return nil, err
			}
			info.CoreId = int32(v)
		case "cpu cores":
			v, err := strconv.ParseInt(value, 10, 64)
			if err != nil {
				return nil, err
			}
			info.CpuCores = int32(v)
		case "apicid":
			v, err := strconv.ParseInt(value, 10, 64)
			if err != nil {
				return nil, err
			}
			info.Apicid = int32(v)
		case "initial apicid":
			v, err := strconv.ParseInt(value, 10, 64)
			if err != nil {
				return nil, err
			}
			info.InitApicid = int32(v)
		case "fpu":
			info.Fpu = value
		case "fpu_exception":
			info.FpuException = value
		case "cpuid level":
			v, err := strconv.ParseInt(value, 10, 64)
			if err != nil {
				return nil, err
			}
			info.CpuidLevel = int32(v)
		case "wp":
			info.Wp = value
		case "flags":
			v := strings.Split(value, " ")
			info.Flags = v
		case "bogomips":
			v, err := strconv.ParseFloat(value, 64)
			if err != nil {
				return nil, err
			}
			info.Bogomips = v
		case "clflush size":
			v, err := strconv.ParseInt(value, 10, 64)
			if err != nil {
				return nil, err
			}
			info.ClflushSize = int32(v)
		case "cache_alignment":
			v, err := strconv.ParseInt(value, 10, 64)
			if err != nil {
				return nil, err
			}
			info.CacheAlignment = int32(v)
		case "address sizes":
			info.AddressSizes = value
		case "power management":
			info.PowerManagement = value
		}
	}
	return res, nil
}
