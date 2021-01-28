# Advance Cpu Reporter Usage

## Overview
This is a advanced cpu reporter whhich will record task in every cpu queue. This program includs two parts, **Watcher** and **Render**. 
### Watcher
Watcher is a recorder that will record the cpu queue data and cpu usage, then save all the data to a file.

### Render
Render will generate a html file that visualizes all the data. The data chart will present as Linemap and Heatmap.

## How to use
### Step.1
Build program
```shell
go build Watcher.go
go build Render.go
```

### Step.2
Start Watcher in background
```shell
./Watcher &
```

### Step.3
Do whatever you need to do

### Step.4
Kill the Watcher and run the Render
```shell
pkill Watcher
./Render
```

The **Render** has two options:
- `-f` : The output filename you want to use, please note that the filename must end with `.html`. The default value is `cpu_report.html`
- `-t` : The page title of the html output file, default value is `CPU Report`

### Check the chart
You will find the cpu_report.html file, and open it on any morder browser like Chrome.

> NOTE: Please keep in mind that if you run the Watcher for a very long time, the output data file of Watcher could be very large, and the Render output html file will also be very large. Still tring figure out a good idea to solve this problem, welcome to any good idea.