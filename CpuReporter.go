package main

import (
	"flag"
	"fmt"
	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/components"
	"github.com/go-echarts/go-echarts/v2/opts"
	"github.com/shirou/gopsutil/cpu"
	"io"
	"os"
	"strconv"
	"time"
)

var (
	// Store the time line array which will be the X axis of charts
	timeArray []string

	// Job running time checking
	second int

	// data used by heat map
	hmData = [][3]float64{
		//{cpu,second,usage}
	}

	// CPU list used by per cpu chart
	cpuList = []string{}

	// Seconds list use by heatmap X axis
	seconds = []string{}

	// opts.LineData list use by total CPU collection
	totalCpu []opts.LineData
)

// Generate random data for bar chart
func getTotalCpuUsage() {
	percentAll, _ := cpu.Percent(time.Second, false)
	totalCpu = append(totalCpu, opts.LineData{Value: percentAll[0]})
}

func createTotalCpuChart(seconds []string, usage []opts.LineData) *charts.Line {
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
		AddSeries("Category A", usage).
		SetSeriesOptions(
			charts.WithAreaStyleOpts(opts.AreaStyle{
				Opacity: 0.5,
			}),
		)
	return line
}

// Get the opts.HeatMapData that the heat map needs
func getCpuHeatMapData() []opts.HeatMapData {
	items := make([]opts.HeatMapData, 0)
	for _, singleData := range hmData {
		items = append(items, opts.HeatMapData{Value: singleData})
	}
	return items
}

//Create heatmap chart
func createPerCpuHeatMap(cpus []string) *charts.HeatMap {
	hm := charts.NewHeatMap()
	hm.SetGlobalOptions(
		charts.WithInitializationOpts(opts.Initialization{
			Width:  "1800px",
			Height: "700px",
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

// Get per cpu usage data
func getSepCpuUsage() {
	percentPerCpu, _ := cpu.Percent(time.Second*5, true)
	for i, cpuUse := range percentPerCpu {
		cpuUse, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", cpuUse), 64)
		var data = [3]float64{float64(second), float64(i), cpuUse}
		hmData = append(hmData, data)
	}
}

func main() {

	// Get filename and page title from command line
	var outputFile string
	var pageTitle string
	var interval string
	flag.StringVar(&outputFile, "f", "cpu_report.html", "output file name, must end with .html, default:cpu_report.html")
	flag.StringVar(&pageTitle, "t", "CPU Report", "page title of the report, default: CPU Report")
	flag.StringVar(&interval, "i", "1", "the interval to record the cpu usage, default: 1s")
	flag.Parse()
	interval_i, err := strconv.Atoi(interval)
	if err != nil {
        // handle error
        fmt.Println(err)
        os.Exit(2)
    }

	// Get CPU list
	cpuN, _ := cpu.Counts(true)
	for i := 0; i < cpuN; i++ {
		cpuList = append(cpuList, strconv.Itoa(i))
	}

	// Start keep rendering report
	for {
		second += interval_i
		timeArray = append(timeArray, strconv.Itoa(second))
		seconds = append(seconds, strconv.Itoa(second))
		//add cpu usage to global var
		getTotalCpuUsage()
		getSepCpuUsage()
		// Generate a cpu report page
		page := components.NewPage()
		page.Initialization.PageTitle = pageTitle
		page.SetLayout(components.PageCenterLayout)
		page.AddCharts(
			createTotalCpuChart(timeArray, totalCpu),
			createPerCpuHeatMap(cpuList),
		)
		f, err := os.Create(outputFile)
		if err != nil {
			panic(err)
		}
		page.Render(io.MultiWriter(f))

		time.Sleep(time.Second * time.Duration(interval_i))
	}

}
