# cpuReporter

## Overview
This is a CPU usage statistics collector of common use. The program will record the total CPU usage and per CPU usage **every 1s** by default during this program running, and render a charts report.

## Usage
You can simply run this program in background.
```shell
./cpuReport &
```
And when you want to end the recording, you could:
```shell
pkill cpuReport
```

## Options
- `-f` : The output filename you want to use, please note that the filename must end with `.html`. The default value is `cpu_report.html`
- `-t` : The page title of the html output file, default value is `CPU Report`
- `-i` : The interval to collect the cpu usage

## TODO
1. Add hover label to heatmap
2. Add support customized interval
