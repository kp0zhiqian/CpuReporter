package main

import (
	"os"
	"bufio"
	"strings"
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

func getTask(cpuNum string) []string {
	var Results []string
	file, err := os.Open("/proc/sched_debug")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "cpu#"+cpuNum) {
			Results = append(Results, line)
		} else if strings.HasPrefix(line, " S") || strings.HasPrefix(line, " I") || strings.HasPrefix(line, ">R") {
			Results = append(Results, line)
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
		result.Second = second
		percentPerCpu, _ := cpu.Percent(time.Second*1, true)
		for i, usg := range percentPerCpu {
			var everySecond CpusTask
			currentCPU, _ := fmt.Printf("%d", i)
			everySecond.Usage = float64(usg)
			everySecond.Cpu = strconv.Itoa(currentCPU)
			everySecond.Task = getTask(string(i))
			result.Data = append(result.Data, everySecond)
		}
		fmt.Println(result)
		enc.Encode(result)
		second += 1
	}
	
}