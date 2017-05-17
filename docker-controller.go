// olb.go
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"golang.org/x/net/context"
)

var ctx context.Context
var cli *client.Client
var config Tcfg

type Message struct {
	Acode   int64
	Astring string
	Aslice  []string
}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func ResponseHelper(code int64, message string, payload []string) (resp string) {
	log.Print(message)
	var m Message
	m.Acode = code
	m.Astring = message
	m.Aslice = payload
	b, _ := json.Marshal(m)
	return string(b)
}

func handler(w http.ResponseWriter, r *http.Request) {

	// split the path on /
	// strips off first /

	urlpath := strings.Split(r.URL.Path[1:], "/")
	parameters := r.URL.Query()

	if len(parameters) < 1 {
		fmt.Fprintf(w, ResponseHelper(501, "No url parameters given.", nil))
		return
	}

	api_key, ok := parameters["api_key"]

	if !ok {
		fmt.Fprintf(w, ResponseHelper(501, "No api_key parameter given!", nil))
		return
	}

	if api_key[0] != config.General.Api_key {
		fmt.Fprintf(w, ResponseHelper(501, "api_keys between server and client not equal!", nil))
		return
	}

	if stringInSlice("service", urlpath) {
		if stringInSlice("list", urlpath) {
			resp := servicelist()
			fmt.Fprintf(w, resp)
		} else if stringInSlice("inspect", urlpath) {
			resp := serviceinspect(urlpath[2])
			fmt.Fprintf(w, resp)
		}
	} else {
		http.Error(w, "Not a valid command call", http.StatusBadRequest)
	}

}

func servicelist() string {

	var m Message

	servicelist, err := cli.ServiceList(ctx, types.ServiceListOptions{})
	if err != nil {
		m.Acode = 500
		b, _ := json.Marshal(m)
		return string(b)
	}

	for _, service := range servicelist {
		fmt.Println(service.ID)
		fmt.Println(service.Spec.Name)
	}

	m.Acode = 500
	b, _ := json.Marshal(m)
	return string(b)
}

func serviceinspect(servicename string) (response string) {
	srvlist, err := cli.ServiceList(ctx, types.ServiceListOptions{})
	if err != nil {
		log.Print("Not a stack manager error")
		log.Print(err)
		return
	}

	var found bool = false
	var id string = ""
	var name string = ""
	var publishedport = ""
	var m Message

	for _, service := range srvlist {
		if service.Spec.Name == servicename {
			found = true
			id = service.ID
			name = service.Spec.Name
		}
	}

	if found == true {
		fmt.Println(id)
		fmt.Println(name)
	} else {
		log.Print("Service not found")
		m.Acode = 500
		m.Astring = "Service not found"
		b, _ := json.Marshal(m)
		return string(b)
	}

	srvinspect, _, err := cli.ServiceInspectWithRaw(ctx, id, types.ServiceInspectOptions{})
	if err != nil {
		log.Print(err)
	}

	publishedport = fmt.Sprint(srvinspect.Endpoint.Ports[0].PublishedPort)

	ndlist, err := cli.NodeList(ctx, types.NodeListOptions{})

	if err != nil {
		log.Print(err)
	}

	// add nodes to message
	for _, node := range ndlist {
		nodename := node.Description.Hostname
		m.Aslice = append(m.Aslice, nodename+" "+publishedport)
	}

	m.Acode = 100
	b, _ := json.Marshal(m)
	return string(b)

}

func main() {

	// read config
	config, _ = ReadConfigfile()

	// init global docker client
	ctx = context.Background()
	var err error
	cli, err = client.NewEnvClient()
	if err != nil {
		log.Panic(err)
	}

	// init webserver
	http.HandleFunc("/", handler)
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Panic(err)
	}
}
