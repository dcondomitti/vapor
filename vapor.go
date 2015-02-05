package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"text/template"
)

type Configuration struct {
	Token string
}

func loadConfiguration() *Configuration {
	c := new(Configuration)
	filename := "configuration.txt"
	token, err := ioutil.ReadFile(filename)

	if err != nil {
		panic(err)
	}

	c.Token = string(token)
	return c
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
	cfg := loadConfiguration()
	mac_address := r.URL.Path[len("/config/host/"):]
	ip := strings.Split(r.RemoteAddr, ":")[0]
	cloud_init := CloudInit{cfg.Token, ip, mac_address}

	log.Printf("request from %s@%s", ip, mac_address)
	generateCloudConfig(cloud_init, w)
}

func main() {
	http.HandleFunc("/config/host/", viewHandler)
	http.ListenAndServe(":8080", nil)
}
