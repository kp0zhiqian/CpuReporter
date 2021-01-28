package main

import (
	"os"
	"bufio"
	"strings"
	"time"
	"fmt"
	"encoding/json"
	"strconv"
	"github.com/shirou/gopsutil/cpu"
)


type CpusTask struct {
	Usage float64
	Cpu string
	Task []string
}

type Timeline struct {
	Second int
	Total float64
	Data []CpusTask

}

var timeline Timeline


func getTask(result *Timeline) {

	
	file, err := os.Open("/proc/sched_debug")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	var currentCPU int
	var taskList []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "cpu#") {
			currentCPU, _ = strconv.Atoi(strings.Split(strings.Split(line, ",")[0], "#")[1])
			if currentCPU != 0 {
				result.Data[currentCPU-1].Task = taskList
				taskList = taskList[:0:0]
			}
		} else if strings.HasPrefix(line, " S") || strings.HasPrefix(line, " I") || strings.HasPrefix(line, ">R") {
			taskList = append(taskList, line)
		}
	}
	result.Data[currentCPU].Task = taskList
}


func main() {

	second := 0

	f, err := os.Create("cpu_data.json")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	enc := json.NewEncoder(f)
	for {
		var result = timeline
		result.Second = second
		percentAll, _ := cpu.Percent(time.Second*1, false)
		result.Total, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", percentAll[0] ), 64)
		percentPerCpu, _ := cpu.Percent(time.Second*1, true)
		for i := 0; i < len(percentPerCpu); i++ {
			var everySecond CpusTask
			everySecond.Usage, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", percentPerCpu[i]), 64)
			everySecond.Cpu = strconv.Itoa(i)
			result.Data = append(result.Data, everySecond)
		}
		getTask(&result)
		enc.Encode(result)
		second += 1
	}
	
}