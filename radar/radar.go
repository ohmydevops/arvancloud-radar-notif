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

type ISP string

const (
	ISP_MCI         ISP = "mci"
	ISP_Irancell    ISP = "irancell"
	ISP_Tehran2     ISP = "tehran-2"
	ISP_Tehran3     ISP = "tehran-3"
	ISP_Hostiran    ISP = "hostiran"
	ISP_Parsonline  ISP = "parsonline"
	ISP_Afranet     ISP = "afranet"
	ISP_BertinaXrx  ISP = "bertina-xrx"
	ISP_BertinaThr  ISP = "bertina-thr"
	ISP_AjkAbrbaran ISP = "ajk-abrbaran"

	// Not found data ISPs
	// ISP_SindadThrFanava ISP = "sindad-thr-fanava"
	// ISP_SindadBuf   ISP = "sindad-buf"
	// ISP_SindadThr   ISP = "sindad-thr"
)

var AllISPs = []ISP{
	ISP_MCI,
	ISP_Irancell,
	ISP_Tehran2,
	ISP_Tehran3,
	ISP_Hostiran,
	ISP_Parsonline,
	ISP_Afranet,
	ISP_BertinaXrx,
	ISP_BertinaThr,
	ISP_AjkAbrbaran,
	// ISP_SindadThrFanava,
	// ISP_SindadBuf,
	// ISP_SindadThr,
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

// CheckISPServiceStatistics fetches the latest monitoring value for the given ISP and service.
func CheckISPServiceStatistics(isp ISP, service Service) (*ServiceStatistics, error) {
	client := getClient()

	body, err := fetchData(client, string(isp))
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
func fetchData(client *http.Client, isp string) ([]byte, error) {
	resp, err := client.Get(baseURL + isp)
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
