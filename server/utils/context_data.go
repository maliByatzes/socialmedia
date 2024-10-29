package utils

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"strings"
	"time"
)

type IPContext struct {
	IP        string `json:"ip"`
	IsPrivate bool   `json:"is_private"`

	Country     string `json:"country"`
	City        string `json:"city"`
	Region      string `json:"region"`
	Timezone    string `json:"timezone"`
	Coordinates struct {
		Latitude  float64 `json:"latitude"`
		Longitude float64 `json:"longitude"`
	} `json:"coordinates"`

	Browser    string `json:"browser"`
	Platform   string `json:"platform"`
	OS         string `json:"os"`
	Device     string `json:"device"`
	DeviceType string `json:"device_type"`

	ISP string `json:"isp"`
	ASN string `json:"asn"`
}

type IPGeolocation struct {
	Status     string  `json:"status"`
	Country    string  `json:"country"`
	City       string  `json:"city"`
	RegionName string  `json:"regionName"`
	Timezone   string  `json:"timezone"`
	ISP        string  `json:"isp"`
	AS         string  `json:"as"`
	Lat        float64 `json:"lat"`
	Lon        float64 `json:"lon"`
}

func GetCurrentContextData(ipAddress string, r *http.Request) (*IPContext, error) {
	context := &IPContext{
		IP: ipAddress,
	}

	ip := net.ParseIP(ipAddress)
	if ip == nil {
		return nil, fmt.Errorf("invalid IP address.")
	}
	context.IsPrivate = isPrivateIP(ip)

	if context.IsPrivate {
		if r != nil {
			parseUserAgent(r.UserAgent(), context)
		}
		return context, nil
	}

	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	resp, err := client.Get(fmt.Sprintf("http://ip-api.com/json/%s", ipAddress))
	if err != nil {
		return nil, fmt.Errorf("error making geolocation request: %v", err)
	}
	defer resp.Body.Close()

	var geoData IPGeolocation
	if err := json.NewDecoder(resp.Body).Decode(&geoData); err != nil {
		return nil, fmt.Errorf("error decoding geolocation response: %v", err)
	}

	context.Country = geoData.Country
	context.City = geoData.City
	context.Region = geoData.RegionName
	context.Timezone = geoData.Timezone
	context.Coordinates.Latitude = geoData.Lat
	context.Coordinates.Longitude = geoData.Lon
	context.ISP = geoData.ISP
	context.ASN = geoData.AS

	if r != nil {
		parseUserAgent(r.UserAgent(), context)
	}

	return context, nil
}

func isPrivateIP(ip net.IP) bool {
	privateNetworks := []struct {
		network string
		mask    string
	}{
		{"10.0.0.0", "255.0.0.0"},
		{"172.16.0.0", "255.240.0.0"},
		{"192.168.0.0", "255.255.0.0"},
	}

	for _, network := range privateNetworks {
		_, subnet, _ := net.ParseCIDR(fmt.Sprintf("%s/%s", network.network, network.mask))
		if subnet.Contains(ip) {
			return true
		}
	}

	return false
}

func parseUserAgent(userAgent string, context *IPContext) {
	ua := strings.ToLower(userAgent)

	switch {
	case strings.Contains(ua, "chrome"):
		context.Browser = "Chrome"
	case strings.Contains(ua, "firefox"):
		context.Browser = "Firefox"
	case strings.Contains(ua, "safari"):
		context.Browser = "Safari"
	case strings.Contains(ua, "edge"):
		context.Browser = "Edge"
	default:
		context.Browser = "Unknown"
	}

	switch {
	case strings.Contains(ua, "windows"):
		context.OS = "Windows"
	case strings.Contains(ua, "mac"):
		context.OS = "macOS"
	case strings.Contains(ua, "linux"):
		context.OS = "Linux"
	case strings.Contains(ua, "android"):
		context.OS = "Android"
	case strings.Contains(ua, "iphone") || strings.Contains(ua, "ipad"):
		context.OS = "iOS"
	default:
		context.OS = "Unknown"
	}

	switch {
	case strings.Contains(ua, "mobile"):
		context.Device = "Mobile"
	case strings.Contains(ua, "tablet"):
		context.Device = "Tablet"
	default:
		context.Device = "Desktop"
	}
}
