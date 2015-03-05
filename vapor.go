package main

import (
	"log"
	"net/http"
	"os"
	"strings"
	"text/template"
)

type Configuration struct {
	Token string
}

var cfg Configuration

func loadConfiguration(c *Configuration) {
	token := os.Getenv("ETCD_DISCOVERY_TOKEN")

	if token == "" {
		panic("ETCD_DISCOVERY_TOKEN not present")
	} else {
		c.Token = token
	}
}

type CloudInit struct {
	Token      string
	IPAddress  string
	MacAddress string
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
	mac_address := r.URL.Path[len("/config/host/"):]
	ip := strings.Split(r.RemoteAddr, ":")[0]
	cloud_init := CloudInit{cfg.Token, ip, mac_address}

	log.Printf("request from %s@%s", ip, mac_address)
	generateCloudConfig(cloud_init, w)
}

func main() {
	loadConfiguration(&cfg)
	http.HandleFunc("/config/host/", viewHandler)
	http.ListenAndServe(":8080", nil)
}
