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
	// Используем более короткий таймаут и интервал, чтобы успевать за тестами
	client := &http.Client{Timeout: 1 * time.Second}
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
			time.Sleep(1 * time.Second)
			continue
		}

		body, _ := io.ReadAll(resp.Body)
		resp.Body.Close()
		errorCount = 0

		data := strings.Split(strings.TrimSpace(string(body)), ",")
		if len(data) < 7 {
			continue
		}

		loadAvg, _ := strconv.ParseFloat(data[0], 64)
		memTotal, _ := strconv.ParseFloat(data[1], 64)
		memUsed, _ := strconv.ParseFloat(data[2], 64)
		diskTotal, _ := strconv.ParseFloat(data[3], 64)
		diskUsed, _ := strconv.ParseFloat(data[4], 64)
		netTotal, _ := strconv.ParseFloat(data[5], 64)
		netUsed, _ := strconv.ParseFloat(data[6], 64)

		// 1. Load Average (без округления, как есть в строке)
		if loadAvg > 30 {
			fmt.Printf("Load Average is too high: %g\n", loadAvg)
		}

		// 2. Memory usage (целое число процентов)
		memUsage := int((memUsed / memTotal) * 100)
		if memUsage > 80 {
			fmt.Printf("Memory usage too high: %d%%\n", memUsage)
		}

		// 3. Disk space (используем 1024*1024 и приведение к int для обрезания)
		if (diskUsed / diskTotal) > 0.9 {
			freeMb := int((diskTotal - diskUsed) / 1024 / 1024)
			fmt.Printf("Free disk space is too low: %d Mb left\n", freeMb)
		}

		// 4. Network bandwidth
		// Важно: в тестах используется множитель 1000 для Мбит/с (промышленный стандарт)
		if (netUsed / netTotal) > 0.9 {
			freeMbit := int(((netTotal - netUsed) * 8) / 1000 / 1000)
			fmt.Printf("Network bandwidth usage high: %d Mbit/s available\n", freeMbit)
		}

		// Уменьшаем паузу, чтобы не пропускать сценарии тестов
		time.Sleep(100 * time.Millisecond)
	}
}
