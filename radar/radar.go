package radar

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"sync"
)

var (
	client     *http.Client
	clientOnce sync.Once
)

// Define types and constants

type Datacenter string

const (
	Datacenter_MCI         Datacenter = "mci"
	Datacenter_Irancell    Datacenter = "irancell"
	Datacenter_Tehran2     Datacenter = "tehran-2"
	Datacenter_Tehran3     Datacenter = "tehran-3"
	Datacenter_Hostiran    Datacenter = "hostiran"
	Datacenter_Parsonline  Datacenter = "parsonline"
	Datacenter_Afranet     Datacenter = "afranet"
	Datacenter_BertinaXrx  Datacenter = "bertina-xrx"
	Datacenter_BertinaThr  Datacenter = "bertina-thr"
	Datacenter_AjkAbrbaran Datacenter = "ajk-abrbaran"

	// Not found data DataCenters
	// DataCenter_SindadThrFanava Datacenter = "sindad-thr-fanava"
	// DataCenter_SindadBuf   Datacenter = "sindad-buf"
	// DataCenter_SindadThr   Datacenter = "sindad-thr"
)

var AllDatacenters = []Datacenter{
	Datacenter_MCI,
	Datacenter_Irancell,
	Datacenter_Tehran2,
	Datacenter_Tehran3,
	Datacenter_Hostiran,
	Datacenter_Parsonline,
	Datacenter_Afranet,
	Datacenter_BertinaXrx,
	Datacenter_BertinaThr,
	Datacenter_AjkAbrbaran,
	// Datacenter_SindadThrFanava,
	// Datacenter_SindadBuf,
	// Datacenter_SindadThr,
}

type Service string

const (
	Service_Google      Service = "google"
	Service_Github      Service = "github"
	Service_Wikipedia   Service = "wikipedia"
	Service_Playstation Service = "playstation"
	Service_Bing        Service = "bing"
	Service_Digikala    Service = "digikala"
	Service_Divar       Service = "divar"
	Service_Aparat      Service = "aparat"
)

var AllServices = []Service{
	Service_Google,
	Service_Github,
	Service_Wikipedia,
	Service_Playstation,
	Service_Bing,
	Service_Digikala,
	Service_Divar,
	Service_Aparat,
}

func ParseService(s string) (Service, bool) {
	for _, svc := range AllServices {
		if string(svc) == s {
			return svc, true
		}
	}
	return "", false
}

// Arvan Cloud API base URL
const baseURL = "https://radar.arvancloud.ir/api/v1/internet-monitoring?isp="

type ServiceStatistics struct {
	Service    Service
	Statistics []float64
}

func (ss *ServiceStatistics) IsAccessibleNow() bool {
	return ss.Statistics[len(ss.Statistics)-1] == 0
}

// CheckDatacenterServiceStatistics fetches the latest monitoring value for the given datacenter and service.
func CheckDatacenterServiceStatistics(datacenter Datacenter, service Service) (*ServiceStatistics, error) {
	client := getClient()

	body, err := fetchData(client, string(datacenter))
	if err != nil {
		return nil, err
	}

	return parseServiceStatsFromData(body, service)
}
func getClient() *http.Client {
	clientOnce.Do(func() {
		jar, _ := cookiejar.New(nil)
		client = &http.Client{Jar: jar}
	})
	return client
}

// fetchData performs the HTTP GET request and returns response body bytes.
func fetchData(client *http.Client, datacenter string) ([]byte, error) {
	resp, err := client.Get(baseURL + datacenter)
	if err != nil {
		return nil, fmt.Errorf("request error: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read error: %v", err)
	}

	return body, nil
}

// parseServiceStatsFromData parses and returns back the data for the service
func parseServiceStatsFromData(data []byte, service Service) (*ServiceStatistics, error) {
	var parsed map[string][]float64
	if err := json.Unmarshal(data, &parsed); err != nil {
		return nil, fmt.Errorf("JSON parse error: %v", err)
	}

	// Get values of the service
	values, ok := parsed[string(service)]
	if !ok || len(values) == 0 {
		return nil, fmt.Errorf("no data for service: %s", service)
	}

	return &ServiceStatistics{Service: service, Statistics: values}, nil
}
