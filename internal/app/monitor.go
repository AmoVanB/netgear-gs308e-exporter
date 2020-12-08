package app

import (
	"crypto/md5"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/amovanb/netgear-gs308e-exporter/internal/pkg/config"
	"github.com/amovanb/netgear-gs308e-exporter/internal/pkg/metrics"
	"github.com/prometheus/client_golang/prometheus"
	"log"
	"math/big"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"
)

func MonitorSwitch(config config.SwitchConfig, interval int) error {
	cookie, err := login(config.Url, config.Password)
	if err != nil {
		return err
	}

	log.Printf("%s - cookie: %s", config.Url, cookie)

	ticker := time.NewTicker(time.Duration(interval) * time.Second)

	for range ticker.C {
		wg := sync.WaitGroup{}
		wg.Add(2)
		go func() {
			if err := updatePortStats(config, cookie); err != nil {
				log.Fatal(err)
			}
			wg.Done()
		}()

		go func() {
			if err := updatePortStatus(config, cookie); err != nil {
				log.Fatal(err)
			}
			wg.Done()
		}()
		wg.Wait()
	}

	return nil
}

// returns the cookie
func login(switchUrl, password string) (string, error) {
	client := http.Client{}

	// get login page
	resp, err := client.Get(fmt.Sprintf("%s/login.cgi", switchUrl))
	if err != nil {
		return "", err
	}

	// parse login page to get the random number
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return "", err
	}

	random, exists := doc.Find("#rand").Attr("value")
	if !exists {
		return "", fmt.Errorf("couldn't find the value of the 'rand' id")
	}

	log.Printf("%s - parsed random value: %s", switchUrl, random)

	// compute the hash to send (see Javascript code of the switch)
	postPwd := fmt.Sprintf("%x", md5.Sum([]byte(merge(password, random))))
	log.Printf("%s - posting password: %s\n", switchUrl, postPwd)

	// post the login
	resp, err = client.PostForm(fmt.Sprintf("%s/login.cgi", switchUrl), url.Values{
		"password": {postPwd},
	})
	if err != nil {
		return "", err
	}

	cookie := resp.Header.Get("Set-Cookie")
	if cookie == "" {
		// check for an error message
		// parse login page to get the random number
		doc, err = goquery.NewDocumentFromReader(resp.Body)
		if err != nil {
			return "", err
		}

		errMsg, exists := doc.Find("#err_msg").Attr("value")
		if exists {
			return "", fmt.Errorf("error logging in: %s", errMsg)
		} else {
			return "", fmt.Errorf("couldn't find 'Set-Cookie' header: %s", resp.Header)
		}
	}

	return cookie, nil
}

// translation of the Javascript function of the web interface
func merge(str1, str2 string) string {
	arr1 := strings.Split(str1, "")
	arr2 := strings.Split(str2, "")
	result := ""
	index1 := 0
	index2 := 0
	for (index1 < len(arr1)) || (index2 < len(arr2)) {
		if index1 < len(arr1) {
			result += arr1[index1]
			index1++
		}
		if index2 < len(arr2) {
			result += arr2[index2]
			index2++
		}
	}

	return result
}

func updatePortStats(config config.SwitchConfig, cookie string) error {
	client := http.Client{}
	urlObj, err := url.Parse(fmt.Sprintf("%s/portStats.htm", config.Url))
	if err != nil {
		return err
	}

	request := &http.Request{
		Method: http.MethodGet,
		URL:    urlObj,
		Header: http.Header{
			"Cookie": []string{cookie},
		},
	}

	resp, err := client.Do(request)
	if err != nil {
		return err
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return err
	}

	for i, node := range doc.Find(".portID").Nodes {
		docFromNode := goquery.NewDocumentFromNode(node)
		labels := prometheus.Labels{"switch": urlObj.Host, "port": strconv.Itoa(i + 1)}
		docFromNode.Find("input").Each(func(j int, selection *goquery.Selection) {
			val := new(big.Int)
			val.SetString(selection.AttrOr("value", "0"), 16)
			switch selection.AttrOr("name", "error") {
			case "rxPkt":
				metrics.ServiceMetricsVar.ReceivedBytes.With(labels).Set(float64(val.Uint64()))
			case "txpkt":
				metrics.ServiceMetricsVar.TransmittedBytes.With(labels).Set(float64(val.Uint64()))
			case "crcPkt":
				metrics.ServiceMetricsVar.CRCErrorPackets.With(labels).Set(float64(val.Uint64()))
			default:
				log.Printf("%s - error parsing input field: unknown input name", urlObj.Host)
			}
		})
	}

	return nil
}

func updatePortStatus(config config.SwitchConfig, cookie string) error {
	client := http.Client{}
	urlObj, err := url.Parse(fmt.Sprintf("%s/status.htm", config.Url))
	if err != nil {
		return err
	}

	request := &http.Request{
		Method: http.MethodGet,
		URL:    urlObj,
		Header: http.Header{
			"Cookie": []string{cookie},
		},
	}

	resp, err := client.Do(request)
	if err != nil {
		return err
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return err
	}

	for i, node := range doc.Find(".portID").Nodes {
		docFromNode := goquery.NewDocumentFromNode(node)
		labels := prometheus.Labels{"switch": urlObj.Host, "port": strconv.Itoa(i + 1)}
		docFromNode.Find("input").Each(func(j int, selection *goquery.Selection) {
			switch j {
			case 2:
				// the third input shows the Up/Down info
				if selection.AttrOr("value", "Down") == "Up" {
					metrics.ServiceMetricsVar.PortStatus.With(labels).Set(1)
				} else {
					metrics.ServiceMetricsVar.PortStatus.With(labels).Set(0)
				}
			case 4:
				// the fifth input shows the link speed
				speed, err := strconv.Atoi(strings.Replace(selection.AttrOr("value", "0M"), "M", "000000", 1))
				if err != nil {
					// 'No Speed' of a disconnected link, set speed to 0
					metrics.ServiceMetricsVar.PortSpeed.With(labels).Set(0)
				} else {
					metrics.ServiceMetricsVar.PortSpeed.With(labels).Set(float64(speed))
				}
			default:
				// we don't care about the other fields
			}
		})
	}

	return nil
}
