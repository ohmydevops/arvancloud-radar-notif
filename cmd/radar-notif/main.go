package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
	"time"

	internal_flag "github.com/ohmydevops/arvancloud-radar-notif/internal/flag"
	"github.com/ohmydevops/arvancloud-radar-notif/internal/notification"
	"github.com/ohmydevops/arvancloud-radar-notif/radar"
)

const ProgramName = "📡 Arvan Cloud Radar Monitor"

// Max consecutive errors to consider outage
const maxConsecutiveErrorsForOutage = 3
const notificationIconPath = "./icon.png"

var (
	DatacenterErrorCounts = make(map[radar.Datacenter]int)
	erroredDatacenters    = make(map[radar.Datacenter]bool)
	mu                    sync.Mutex // protects DatacenterErrorCounts & erroredDatacenters
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

	fmt.Printf("✅ Monitoring service: %s\n\n", capitalizeFirst(cfg.Service))

	// performDelay(cfg.CheckDelay)

	for {
		fmt.Printf("⏰ %s\n", time.Now().Format("15:04:05"))

		var wg sync.WaitGroup

		for _, datacenter := range radar.AllDatacenters {
			wg.Add(1)
			go func(dc radar.Datacenter) {
				defer wg.Done()
				checkDatacenter(dc, radar.Service(cfg.Service), notifiersManager)
			}(datacenter)
		}

		wg.Wait()

		performDelay(cfg.CheckDelay)
	}
}

// checkDatacenter handles checking & notification for a single datacenter
func checkDatacenter(datacenter radar.Datacenter, service radar.Service, notifiersManager *notification.NotifiersManager) {

	// Retrieve connectivity statistics between the datacenter and the service
	stats, err := radar.CheckDatacenterServiceStatistics(datacenter, service)
	if err != nil {
		fmt.Printf("⚠️ Statistics: %v from [%s]\n", err, datacenter)
		return
	}

	// Prevents race condition
	mu.Lock()
	defer mu.Unlock()

	if stats.IsAccessibleNow() {
		if erroredDatacenters[datacenter] {
			title := "🟢 Internet Restored"
			msg := fmt.Sprintf("%s is reachable again from %s", capitalizeFirst(string(service)), datacenter)

			// Notify through notification mechanisms
			if err := notifiersManager.Notify(title, msg); err != nil {
				log.Printf("❌ Notification: %v from [%s]", err, datacenter)
			}
		}
		erroredDatacenters[datacenter] = false
		DatacenterErrorCounts[datacenter] = 0
	} else {
		DatacenterErrorCounts[datacenter]++
		if DatacenterErrorCounts[datacenter] >= maxConsecutiveErrorsForOutage && !erroredDatacenters[datacenter] {
			title := "🔴 Internet Outage"
			msg := fmt.Sprintf("%s is unreachable from %s", capitalizeFirst(string(service)), datacenter)

			// Notify through notification mechanisms
			if err := notifiersManager.Notify(title, msg); err != nil {
				fmt.Printf("❌ Notification: %v from [%s]", err, datacenter)
			}
			erroredDatacenters[datacenter] = true
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

// performDelay sleeps until next full minute
func performDelay(minutes int) {
	time.Sleep(time.Duration(minutes) * time.Minute)
}

// capitalizeFirst makes the first letter uppercase
func capitalizeFirst(s string) string {
	if len(s) == 0 {
		return s
	}
	return strings.ToUpper(s[:1]) + s[1:]
}
