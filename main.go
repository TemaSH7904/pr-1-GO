package main

import (
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func main() {
	client := &http.Client{Timeout: 2 * time.Second}
	errorCount := 0
	url := "http://srv.msk01.gigacorp.local/_stats"

	for {
		resp, err := client.Get(url)
		if err != nil || resp.StatusCode != http.StatusOK {
			errorCount++
			if errorCount >= 3 {
				fmt.Println("Unable to fetch server statistic")
			}
			if resp != nil {
				resp.Body.Close()
			}
			time.Sleep(5 * time.Second)
			continue
		}

		body, _ := io.ReadAll(resp.Body)
		resp.Body.Close()

		data := strings.Split(strings.TrimSpace(string(body)), ",")
		if len(data) != 7 {
			errorCount++
			continue
		}

		errorCount = 0

		loadAvg, _ := strconv.ParseFloat(data[0], 64)
		memTotal, _ := strconv.ParseFloat(data[1], 64)
		memUsed, _ := strconv.ParseFloat(data[2], 64)
		diskTotal, _ := strconv.ParseFloat(data[3], 64)
		diskUsed, _ := strconv.ParseFloat(data[4], 64)
		netTotal, _ := strconv.ParseFloat(data[5], 64)
		netUsed, _ := strconv.ParseFloat(data[6], 64)

		if loadAvg > 30 {
			fmt.Printf("Load Average is too high: %v\n", loadAvg)
		}

		memUsage := (memUsed / memTotal) * 100
		if memUsage > 80 {
			fmt.Printf("Memory usage too high: %.0f%%\n", memUsage)
		}

		diskUsage := (diskUsed / diskTotal) * 100
		if diskUsage > 90 {
			freeMb := (diskTotal - diskUsed) / 1024 / 1024
			fmt.Printf("Free disk space is too low: %.0f Mb left\n", freeMb)
		}

		netUsage := (netUsed / netTotal) * 100
		if netUsage > 90 {
			freeMbit := ((netTotal - netUsed) * 8) / 1024 / 1024
			fmt.Printf("Network bandwidth usage high: %.0f Mbit/s available\n", freeMbit)
		}

		time.Sleep(5 * time.Second)
	}
}
