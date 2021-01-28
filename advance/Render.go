package main

import  (
	"os"
	"bufio"
	"encoding/json"
	"io"
	"flag"
	"strings"
	"strconv"
	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/components"
	"github.com/go-echarts/go-echarts/v2/opts"
)

type PerCPU struct {
	Usage int `json:"Usage"`
	CPU string `json:"Cpu"`
	Task []string `json:"Task"`
}

type PerSecond struct {
	Second int `json:"Second"`
	Data []PerCPU  `json:"Data"`
	Total float64 `json:"Total"`
}

type SingleHmData struct {
	Second float64
	CPU float64
	Usage float64
	Task []string
}

var (
	timeline []PerSecond
	hmData = []SingleHmData{}
	seconds = []int{}
	cpuList = []string{}
	totalCpu []opts.LineData
) 

func generateCpulist() {
	for _, c := range timeline[0].Data {
		cpuList = append(cpuList, string(c.CPU))
	}
}

func getTotalCpuUsage() {
	for _, s := range timeline {
		totalCpu = append(totalCpu, opts.LineData{Value: s.Total})
	}
		
}

var totalCpuFormatter = `function (params) {
	return params.value + '%'
}`

func createTotalCpuChart() *charts.Line {
	// Create a new bar instance
	line := charts.NewLine()
	// Set some global options like Title/Legend/ToolTip or anything else
	line.SetGlobalOptions(
		charts.WithInitializationOpts(opts.Initialization{
			Width:  "1800px",
			Height: "700px",
		}),
		charts.WithTitleOpts(opts.Title{
			Title: "Total CPU Usage",
		}),
		charts.WithTooltipOpts(opts.Tooltip{
			Show: true,
			Formatter: opts.FuncOpts(totalCpuFormatter),
		}),
		charts.WithYAxisOpts(opts.YAxis{
			Name: "CPU Usage(%)",
		}),
		charts.WithXAxisOpts(opts.XAxis{
			Name: "Seconds",
		}),
		charts.WithDataZoomOpts(opts.DataZoom{
			Start:      0,
			End:        100,
			XAxisIndex: []int{0},
		}),
	)

	// Put data into instance
	line.SetXAxis(seconds).
		AddSeries("Total CPU Usage", totalCpu).
		SetSeriesOptions(
			charts.WithAreaStyleOpts(opts.AreaStyle{
				Opacity: 0.5,
			}),
		)
	return line
}

func generateHeatMapData() {
	for _, s := range timeline {
		for _, c := range s.Data {
			CPU_f,_ := strconv.ParseFloat(c.CPU, 64)
			var r SingleHmData
			r.Second = float64(s.Second)
			r.CPU = float64(CPU_f)
			r.Usage = float64(c.Usage)
			r.Task = c.Task
			hmData = append(hmData, r)
		}
		seconds = append(seconds, s.Second)
	}

}

func getCpuHeatMapData() []opts.HeatMapData {
	items := make([]opts.HeatMapData, 0)
	for _, singleData := range hmData {
		var tasklist string
		var valuelist [3]float64
		tasklist += "<table>"
		for _, task := range singleData.Task{
			task = strings.Join(strings.Fields(task), " ")
			tasklist += "<tr>"
			colums := strings.Split(task, " ")
			for _, col := range colums {
				tasklist += "<td>" + col + "</td>"
			}
			tasklist += "</tr>"
		}
		tasklist += "</table>"
		valuelist[0] = singleData.Second
		valuelist[1] = singleData.CPU
		valuelist[2] = singleData.Usage
		items = append(items, opts.HeatMapData{Name: tasklist, Value: valuelist})
	}
	return items
}


var getCpuTask = `function (params) {
		return params.name
	}`

//Create heatmap chart
func createPerCpuHeatMap(cpus []string) *charts.HeatMap {
	hm := charts.NewHeatMap()
	hm.SetGlobalOptions(
		charts.WithInitializationOpts(opts.Initialization{
			Width:  "1800px",
			Height: "900px",
		}),
		charts.WithTooltipOpts(opts.Tooltip{
			Show: true,
			Formatter: opts.FuncOpts(getCpuTask),
		}),
		charts.WithTitleOpts(opts.Title{
			Title: "Per CPU Usage HeatMap",
		}),
		charts.WithXAxisOpts(opts.XAxis{
			Name:      "Seconds",
			Type:      "category",
			SplitArea: &opts.SplitArea{Show: true},
		}),
		charts.WithYAxisOpts(opts.YAxis{
			Name:      "CPU",
			Type:      "category",
			Data:      cpus,
			SplitArea: &opts.SplitArea{Show: true},
		}),
		charts.WithVisualMapOpts(opts.VisualMap{
			Calculable: true,
			Min:        0,
			Max:        100,
			InRange: &opts.VisualMapInRange{
				Color: []string{"#50a3ba", "#eac736", "#d94e5d"},
			},
		}),
		charts.WithDataZoomOpts(opts.DataZoom{
			Type:       "slider",
			Start:      0,
			End:        100,
			XAxisIndex: []int{0},
		}),
	)

	hm.SetXAxis(seconds).
		AddSeries("cpus_heatmap", getCpuHeatMapData()).
		SetSeriesOptions(
			charts.WithLabelOpts(opts.Label{
				Show:     false,
				Position: "inside",
			}),
		)
	return hm

}


func main() {

	var outputFile string
	var pageTitle string
	var interval string
	flag.StringVar(&outputFile, "f", "cpu_report.html", "output file name, must end with .html, default:cpu_report.html")
	flag.StringVar(&pageTitle, "t", "CPU Report", "page title of the report, default: CPU Report")
	flag.StringVar(&interval, "i", "1", "the interval to record the cpu usage, default: 1s")
	flag.Parse()
	
	f, err := os.Open("cpu_data.json")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	s := bufio.NewReader(f)
	for {
		var nowData PerSecond
		line, err := s.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			panic(err)
		}
		json.Unmarshal([]byte(line), &nowData)
		timeline = append(timeline, nowData)
	}
	getTotalCpuUsage()
	generateCpulist()
	generateHeatMapData()
	page := components.NewPage()
	page.Initialization.PageTitle = pageTitle
	page.SetLayout(components.PageCenterLayout)
	page.AddCharts(
		createTotalCpuChart(),
		createPerCpuHeatMap(cpuList),
	)
	finalFile, err := os.Create(outputFile)
	if err != nil {
		panic(err)
	}
	page.Render(io.MultiWriter(finalFile))
	
}