package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	internal_flag "github.com/ohmydevops/arvancloud-radar-notif/internal/flag"
	"github.com/ohmydevops/arvancloud-radar-notif/internal/notification"
	"github.com/ohmydevops/arvancloud-radar-notif/radar"
)

// Max consecutive errors to consider outage
const maxISPError = 3

const notificationIconPath = "./icon.png"
const ProgramName = "📡 Arvan Cloud Radar Monitor"

var (
	ISPErrorCounts = make(map[radar.Datacenter]int)
	erroredISPs    = make(map[radar.Datacenter]bool)
	mu             sync.Mutex // protects ISPErrorCounts & erroredISPs
)

func main() {
	cfg, err := internal_flag.ParseFlags()
	if err != nil {
		fmt.Println("❌", err)
		flag.Usage()
		os.Exit(1)
	}

	if cfg.ShowServices {
		printServices()
		os.Exit(0)
	}

	// Create notification manager
	notifiers := []notification.Notifier{
		notification.NewConsoleNotifier(),
		notification.NewDesktopNotofier(ProgramName, notificationIconPath),
	}
	notifiersManager := notification.NewNotofiersManager(notifiers)

	fmt.Println(ProgramName)

	fmt.Printf("✅ Monitoring service: %s\n", cfg.Service)

	//waitUntilNextMinute()

	for {
		fmt.Printf("⏰ %s\n", time.Now().Format("15:04:05"))

		var wg sync.WaitGroup

		for _, isp := range radar.AllDatacenters {
			wg.Add(1)
			go func(isp radar.Datacenter) {
				defer wg.Done()
				checkISP(isp, radar.Service(cfg.Service), notifiersManager)
			}(isp)
		}

		wg.Wait()
		time.Sleep(10 * time.Second)
	}
}

// checkISP handles checking & notification for a single ISP
func checkISP(datacenter radar.Datacenter, service radar.Service, notifiersManager *notification.NotifiersManager) {
	stats, err := radar.CheckDatacenterServiceStatistics(datacenter, service)
	if err != nil {
		fmt.Printf("[%s] ⚠️ %v\n", datacenter, err)
		return
	}

	mu.Lock()
	defer mu.Unlock()

	if stats.IsAccessibleNow() {
		if erroredISPs[datacenter] {
			title := "🟢 Internet Restored"
			msg := fmt.Sprintf("%s is reachable again from %s", service, datacenter)
			if err := notifiersManager.Notify(title, msg); err != nil {
				log.Printf("[%s] ❌ Notification error: %v", datacenter, err)
			}
		}
		erroredISPs[datacenter] = false
		ISPErrorCounts[datacenter] = 0
	} else {
		ISPErrorCounts[datacenter]++
		if ISPErrorCounts[datacenter] > maxISPError && !erroredISPs[datacenter] {
			title := "🔴 Internet Outage"
			msg := fmt.Sprintf("%s unreachable from %s", service, datacenter)
			if err := notifiersManager.Notify(title, msg); err != nil {
				fmt.Printf("[%s] ❌ Notification error: %v", datacenter, err)
			}
			erroredISPs[datacenter] = true
		}
	}
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
