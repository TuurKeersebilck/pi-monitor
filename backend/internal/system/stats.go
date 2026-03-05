package system

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"syscall"
)

type Stats struct {
	CPU  CPUStats  `json:"cpu"`
	RAM  RAMStats  `json:"ram"`
	Disk DiskStats `json:"disk"`
	Temp TempStats `json:"temp"`
}

type TempStats struct {
	CPUTempC float64 `json:"cpu_temp_c"`
}

type CPUStats struct {
	UsagePercent float64 `json:"usage_percent"`
}

type RAMStats struct {
	TotalMB     uint64  `json:"total_mb"`
	UsedMB      uint64  `json:"used_mb"`
	UsagePercent float64 `json:"usage_percent"`
}

type DiskStats struct {
	TotalGB      float64 `json:"total_gb"`
	UsedGB       float64 `json:"used_gb"`
	UsagePercent float64 `json:"usage_percent"`
}

func GetStats() (*Stats, error) {
	cpu, err := getCPU()
	if err != nil {
		return nil, fmt.Errorf("cpu: %w", err)
	}

	ram, err := getRAM()
	if err != nil {
		return nil, fmt.Errorf("ram: %w", err)
	}

	disk, err := getDisk("/")
	if err != nil {
		return nil, fmt.Errorf("disk: %w", err)
	}

	temp := getTemp()

	return &Stats{CPU: cpu, RAM: ram, Disk: disk, Temp: temp}, nil
}

// prevCPUStat holds the last CPU snapshot for delta-based measurement across polls.
var prevCPUStat []uint64

// getCPU computes CPU usage as the delta between the current and previous /proc/stat
// snapshot. On the first call it returns 0% and stores the baseline.
func getCPU() (CPUStats, error) {
	cur, err := readCPUStat()
	if err != nil {
		return CPUStats{}, err
	}

	if prevCPUStat == nil {
		prevCPUStat = cur
		return CPUStats{UsagePercent: 0}, nil
	}

	totalDiff := float64(sum(cur) - sum(prevCPUStat))
	idleDiff := float64(cur[3] - prevCPUStat[3])
	prevCPUStat = cur

	if totalDiff == 0 {
		return CPUStats{UsagePercent: 0}, nil
	}

	usage := (1 - idleDiff/totalDiff) * 100
	return CPUStats{UsagePercent: roundTo(usage, 1)}, nil
}

func readCPUStat() ([]uint64, error) {
	f, err := os.Open("/proc/stat")
	if err != nil {
		return nil, err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "cpu ") {
			fields := strings.Fields(line)[1:]
			var vals []uint64
			for _, f := range fields {
				v, _ := strconv.ParseUint(f, 10, 64)
				vals = append(vals, v)
			}
			return vals, nil
		}
	}
	return nil, fmt.Errorf("cpu line not found in /proc/stat")
}

func getRAM() (RAMStats, error) {
	f, err := os.Open("/proc/meminfo")
	if err != nil {
		return RAMStats{}, err
	}
	defer f.Close()

	info := make(map[string]uint64)
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		fields := strings.Fields(scanner.Text())
		if len(fields) >= 2 {
			key := strings.TrimSuffix(fields[0], ":")
			val, _ := strconv.ParseUint(fields[1], 10, 64)
			info[key] = val // values are in kB
		}
	}

	totalKB := info["MemTotal"]
	availKB := info["MemAvailable"]
	usedKB := totalKB - availKB

	totalMB := totalKB / 1024
	usedMB := usedKB / 1024

	var pct float64
	if totalMB > 0 {
		pct = roundTo(float64(usedMB)/float64(totalMB)*100, 1)
	}

	return RAMStats{
		TotalMB:      totalMB,
		UsedMB:       usedMB,
		UsagePercent: pct,
	}, nil
}

func getDisk(path string) (DiskStats, error) {
	var stat syscall.Statfs_t
	if err := syscall.Statfs(path, &stat); err != nil {
		return DiskStats{}, err
	}

	total := float64(stat.Blocks) * float64(stat.Bsize)
	free := float64(stat.Bfree) * float64(stat.Bsize)
	used := total - free

	toGB := func(b float64) float64 { return roundTo(b/1024/1024/1024, 1) }

	var pct float64
	if total > 0 {
		pct = roundTo(used/total*100, 1)
	}

	return DiskStats{
		TotalGB:      toGB(total),
		UsedGB:       toGB(used),
		UsagePercent: pct,
	}, nil
}

func getTemp() TempStats {
	data, err := os.ReadFile("/sys/class/thermal/thermal_zone0/temp")
	if err != nil {
		return TempStats{}
	}
	raw, err := strconv.ParseInt(strings.TrimSpace(string(data)), 10, 64)
	if err != nil {
		return TempStats{}
	}
	return TempStats{CPUTempC: roundTo(float64(raw)/1000, 1)}
}

func sum(vals []uint64) uint64 {
	var s uint64
	for _, v := range vals {
		s += v
	}
	return s
}

func roundTo(val float64, decimals int) float64 {
	pow := float64(1)
	for i := 0; i < decimals; i++ {
		pow *= 10
	}
	return float64(int(val*pow+0.5)) / pow
}
