package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"time"

	"github.com/gen2brain/beeep"
)

const baseURL = "https://radar.arvancloud.ir/api/v1/internet-monitoring?isp="

var (
	errorCounts = make(map[string]int)
	erroredISPs = make(map[string]bool)
)

var serviceIndicator string = "google"

var isps = []string{
	"sindad-buf",
	"sindad-thr-fanava",
	"sindad-thr",
	"bertina-xrx",
	"ajk-abrbaran",
	"tehran-3",
	"tehran-2",
	"bertina-thr",
	"hostiran",
	"parsonline",
	"afranet",
	"mci",
	"irancell",
}

func fetchData(client *http.Client, isp string) (float64, error) {
	url := baseURL + isp

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return 0, fmt.Errorf("❌ request creation error: %v", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		return 0, fmt.Errorf("❌ request error: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, fmt.Errorf("❌ error reading response: %v", err)
	}

	var data map[string][]float64
	if err := json.Unmarshal(body, &data); err != nil {
		return 0, fmt.Errorf("❌ JSON parse error: %v", err)
	}

	googleValues, ok := data[serviceIndicator]
	if !ok || googleValues == nil || len(googleValues) == 0 {
		return 0, fmt.Errorf("-")
	}

	return googleValues[len(googleValues)-1], nil
}

func checkStatus(isp string, value float64) {
	fmt.Printf("[%s] => Value: %.2f\n", isp, value)
	beeep.AppName = "Arvan Cloud Radar"
	if value != 0 {
		errorCounts[isp]++
		if errorCounts[isp] >= 3 && !erroredISPs[isp] {
			err := beeep.Notify("🔴 Internet Outage", fmt.Sprintf("Google unreachable from %s", isp), "./icon.png")
			if err != nil {
				fmt.Printf("[%s] ❌ Notification error: %v\n", isp, err)
			}
			erroredISPs[isp] = true
		}
	} else {
		if erroredISPs[isp] {
			err := beeep.Notify("🟢 Internet Outage fixed", fmt.Sprintf("Google reachable from %s", isp), "./icon.png")
			if err != nil {
				fmt.Printf("[%s] ❌ Notification error: %v\n", isp, err)
			}
			fmt.Printf("[%s] 🟢 Google is reachable again\n", isp)
		}
		errorCounts[isp] = 0
		erroredISPs[isp] = false
	}
}

func waitUntilNextMinute() {
	now := time.Now()
	delay := time.Until(now.Truncate(time.Minute).Add(time.Minute))
	time.Sleep(delay)
}

func main() {
	fmt.Println("Arvan Cloud Outage Notification started ...")
	jar, _ := cookiejar.New(nil)
	client := &http.Client{Jar: jar}

	waitUntilNextMinute()
	for {
		fmt.Printf("⏰ %s\n", time.Now().Format("15:04:05"))
		for _, isp := range isps {
			go func(isp string) {
				value, err := fetchData(client, isp)
				if err != nil {
					fmt.Printf("[%s] %v\n", isp, err)
					return
				}
				checkStatus(isp, value)
			}(isp)
		}
		time.Sleep(1 * time.Minute)
	}
}
