// olb.go
package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"golang.org/x/net/context"
)

var ctx context.Context
var cli *client.Client

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

func handler(w http.ResponseWriter, r *http.Request) {

	// split the path on /
	// strips off first /

	urlpath := strings.Split(r.URL.Path[1:], "/")

	if stringInSlice("service", urlpath) {
		if stringInSlice("list", urlpath) {
			resp := servicelist()
			fmt.Fprintf(w, resp)
		} else if stringInSlice("inspect", urlpath) {
			serviceinspect(urlpath[2])
		}
	} else {
		http.Error(w, "Not a valid command call", http.StatusBadRequest)
	}

}

func servicelist() (response string) {

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

func serviceinspect(servicename string) {
	srvlist, err := cli.ServiceList(ctx, types.ServiceListOptions{})
	if err != nil {
		fmt.Print("Not a stack manager error")
		return
	}

	var found bool = false
	var id string = ""
	var name string = ""

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
		fmt.Println("Service not found Error")
		return
	}

	srvinspect, _, err := cli.ServiceInspectWithRaw(ctx, id, types.ServiceInspectOptions{})
	if err != nil {
		fmt.Print("error")
	}

	fmt.Println(srvinspect.Endpoint.Ports[0].PublishedPort)

	ndlist, err := cli.NodeList(ctx, types.NodeListOptions{})

	if err != nil {
		fmt.Print("error")
	}

	for _, node := range ndlist {
		fmt.Println(node.Description.Hostname)
	}

}

func main() {
	// init global docker client
	ctx = context.Background()
	var err error
	cli, err = client.NewEnvClient()
	if err != nil {
		panic(err)
	}

	// init webserver
	http.HandleFunc("/", handler)
	http.ListenAndServe(":6666", nil)
}
