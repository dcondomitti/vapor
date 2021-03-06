package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"text/template"
)

type Configuration struct {
	Token          string
	HostnameFormat string
}

var cfg Configuration

func loadConfiguration(c *Configuration) {
	token := os.Getenv("ETCD_DISCOVERY_TOKEN")
	hostname_format := os.Getenv("HOSTNAME_FORMAT")

	if token == "" {
		panic("ETCD_DISCOVERY_TOKEN not present")
	} else if hostname_format == "" {
		panic("HOSTNAME_FORMAT not present")
	} else {
		c.Token = token
		c.HostnameFormat = hostname_format
	}
}

type CloudInit struct {
	Token      string
	IPAddress  string
	MacAddress string
	Hostname   string
	EtcdId     string
}

type Host struct {
	IpAddress  string
	MacAddress string
	Hostname   string
	EtcdId     string
}

func NewHost(ip_address string, mac_address string) *Host {
	mac_address_parts := strings.Split(mac_address, ":")
	suffix := mac_address_parts[len(mac_address_parts)-2] + mac_address_parts[len(mac_address_parts)-1]
	hostname := fmt.Sprintf(cfg.HostnameFormat, suffix)
	etcd_id := strings.Split(hostname, ".")[0]

	h := Host{ip_address, mac_address, hostname, etcd_id}
	return &h
}

func generateCloudConfig(c CloudInit, w http.ResponseWriter) {
	tmpl, err := template.ParseFiles("cloud-config.template")
	if err != nil {
		panic(err)
	}

	err = tmpl.Execute(w, c)
	if err != nil {
		panic(err)
	}
}

func viewHandler(w http.ResponseWriter, r *http.Request) {
	ip := strings.Split(r.RemoteAddr, ":")[0]
	mac_address := r.URL.Path[len("/config/host/"):]
	h := NewHost(ip, mac_address)

	cloud_init := CloudInit{cfg.Token, h.IpAddress, h.MacAddress, h.Hostname, h.EtcdId}

	log.Printf("request from %s@%s", ip, mac_address)
	generateCloudConfig(cloud_init, w)
}

func main() {
	loadConfiguration(&cfg)
	http.HandleFunc("/config/host/", viewHandler)
	http.ListenAndServe(":8080", nil)
}
