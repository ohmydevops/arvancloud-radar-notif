package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/gen2brain/beeep"
	"github.com/ohmydevops/arvancloud-radar-notif/radar"
)

// Max consecutive errors to consider outage
const maxISPError = 3

var (
	ISPErrorCounts = make(map[radar.ISP]int)
	erroredISPs    = make(map[radar.ISP]bool)
	mu             sync.Mutex // protects ISPErrorCounts & erroredISPs
)

func main() {
	cfg, err := parseFlags()
	if err != nil {
		if errors.Is(err, ErrEmptyService) {
			flag.Usage()
			os.Exit(1)
		}
		fmt.Println("❌", err)
		flag.Usage()
		os.Exit(1)
	}

	if cfg.ShowServices {
		printServices()
		os.Exit(0)
	}

	fmt.Println("📡 Arvan Cloud Radar Monitor")

	fmt.Printf("✅ Monitoring service: %s\n", cfg.Service)

	waitUntilNextMinute()

	for {
		fmt.Printf("⏰ %s\n", time.Now().Format("15:04:05"))

		var wg sync.WaitGroup

		for _, isp := range radar.AllISPs {
			wg.Add(1)
			go func(isp radar.ISP) {
				defer wg.Done()
				checkISP(isp, radar.Service(cfg.Service))
			}(isp)
		}

		wg.Wait()
		time.Sleep(1 * time.Minute)
	}
}

// checkISP handles checking & notification for a single ISP
func checkISP(isp radar.ISP, service radar.Service) {
	stats, err := radar.CheckISPServiceStatistics(isp, service)
	if err != nil {
		fmt.Printf("[%s] ⚠️ %v\n", isp, err)
		return
	}

	mu.Lock()
	defer mu.Unlock()

	if stats.IsAccessibleNow() {
		if erroredISPs[isp] {
			notifyRestored(service, isp)
		}
		erroredISPs[isp] = false
		ISPErrorCounts[isp] = 0
	} else {
		ISPErrorCounts[isp]++
		if ISPErrorCounts[isp] > maxISPError && !erroredISPs[isp] {
			notifyOutage(service, isp)
			erroredISPs[isp] = true
		}
	}
}

// notifyOutage sends notification when service becomes unreachable
func notifyOutage(service radar.Service, isp radar.ISP) {
	msg := fmt.Sprintf("%s unreachable from %s", service, isp)
	if err := beeep.Notify("🔴 Internet Outage", msg, "./icon.png"); err != nil {
		fmt.Printf("[%s] ❌ Notification error: %v", isp, err)
	}
	fmt.Printf("[%s] 🔴 %s outage detected", isp, service)
}

// notifyRestored sends notification when service is reachable again
func notifyRestored(service radar.Service, isp radar.ISP) {
	msg := fmt.Sprintf("%s is reachable again from %s", service, isp)
	if err := beeep.Notify("🟢 Internet Restored", msg, "./icon.png"); err != nil {
		log.Printf("[%s] ❌ Notification error: %v", isp, err)
	}
	log.Printf("[%s] 🟢 %s restored", isp, service)
}

// printServices prints available services
func printServices() {
	fmt.Println("Available services:")
	for _, s := range radar.AllServices {
		fmt.Printf("  - %s\n", s)
	}

}

// waitUntilNextMinute sleeps until next full minute
func waitUntilNextMinute() {
	time.Sleep(time.Until(time.Now().Truncate(time.Minute).Add(time.Minute)))
}
