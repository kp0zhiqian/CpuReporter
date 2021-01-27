package main

import (
	"os"
	"bufio"
	"strings"
	"strconv"
	"time"
	"fmt"
	"encoding/json"
	"github.com/shirou/gopsutil/cpu"
)

type CpusTask struct {
	Usage float64
	Cpu string
	Task []string
}

type Timeline struct {
	Second int
	Data []CpusTask
}

func getTask() []CpusTask {
	var Results []CpusTask
	file, err := os.Open("sched_debug")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	var cpu_i int64
	
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "cpu#") {
			var SingleCpu CpusTask
			cpuTag := strings.Split(strings.Split(line, ",")[0], "#")[1]
			cpu_i, _ = strconv.ParseInt(cpuTag, 10, 64)
			percentPerCpu, _ := cpu.Percent(time.Second*1, true)
			SingleCpu.Usage = float64(percentPerCpu[cpu_i])
			fmt.Println(SingleCpu.Usage)
			SingleCpu.Cpu = cpuTag
			Results = append(Results, SingleCpu)
		} else if strings.HasPrefix(line, " S") || strings.HasPrefix(line, " I") || strings.HasPrefix(line, ">R") {
			Results[cpu_i].Task = append(Results[cpu_i].Task, line)
		}
	}
	
	return Results
}


func main() {

	second := 0

	var result Timeline

	f, err := os.Create("cpu_data.json")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	enc := json.NewEncoder(f)
	for {
		cpuTaskNow := getTask()
		result.Second = second
		result.Data = cpuTaskNow
		time.Sleep(1 * time.Second)
		enc.Encode(result)
		second += 1
	}
	
}