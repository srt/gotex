package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"text/template"
)

type GotexServer struct {
}

func (s GotexServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	body, _ := ioutil.ReadAll(r.Body)

	var model interface{}
	err := json.Unmarshal(body, &model)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotAcceptable)
		return
	}

	log.Printf("Received model %q", model)

	var t *template.Template
	t, err = template.ParseFiles("invoice.tmpl")
	if err != nil {
		panic(err)
	}

	err = t.Execute(w, model)
	if err != nil {
		panic(err)
	}

}

var configPath string
var config Config

func init() {
	flag.StringVar(&configPath, "c", "/etc/gotex.conf", "Config path")
}

func reloadConfig(c <-chan os.Signal) {
	for s := range c {
		log.Printf("Got %s signal: Reloading configuration", s)
		newConfig, err := ReadConfig(configPath)
		if err == nil {
			config = newConfig
		} else {
			log.Println(err)
		}
	}
}

func main() {
	flag.Parse()
	os.Exit(run())
}

func run() int {
	var err error

	config, err = ReadConfig(configPath)
	/*
		if err != nil {
			log.Fatal(err)
			return 1
		}
	*/

	http.Handle("/", http.StripPrefix("/", GotexServer{}))

	server := &http.Server{Addr: config.Addr}
	listener, err := net.Listen("tcp", server.Addr)
	if err != nil {
		log.Fatal(err)
		return 1
	}

	go server.Serve(listener)

	log.Printf("HTTP server started on %s", config.Addr)

	hupChannel := make(chan os.Signal, 1)
	signal.Notify(hupChannel, syscall.SIGHUP)
	go reloadConfig(hupChannel)

	killChannel := make(chan os.Signal, 1)
	signal.Notify(killChannel, os.Kill, os.Interrupt)

	<-killChannel
	log.Println("Exiting")
	listener.Close()

	return 0
}
