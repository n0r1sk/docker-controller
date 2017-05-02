// olb.go
package main

import (
	"fmt"
	"net/http"
	//"path"
	"strings"

	"golang.org/x/net/context"

	//"encoding/json"

	"github.com/docker/docker/api/types"
	//"github.com/docker/docker/api/types/swarm"
	"github.com/docker/docker/client"
)

var ctx context.Context
var cli *client.Client

func handler(w http.ResponseWriter, r *http.Request) {

	// TODO check if r.URL.Path is path!

	// split the path on /
	urlpath := strings.Split(r.URL.Path, "/")

	fmt.Println(urlpath[2])
	switch urlpath[1] {
	case "service":
		servicecallhandler(urlpath[2], urlpath[3])
	default:
		fmt.Fprintf(w, "not a valid call")
	}

}

func servicecallhandler(command string, arg string) {

	switch command {
	case "list":
		servicelist()
	case "inspect":
		serviceinspect(arg)
	default:
		fmt.Println("commant of list not known or implemented")
	}

}

func servicelist() {
	servicelist, err := cli.ServiceList(ctx, types.ServiceListOptions{})
	if err != nil {
		fmt.Print("error")
	}

	for _, service := range servicelist {
		fmt.Println(service.ID)
		fmt.Println(service.Spec.Name)
	}
}

func serviceinspect(servicename string) {
	srvlist, err := cli.ServiceList(ctx, types.ServiceListOptions{})
	if err != nil {
		fmt.Print("error")
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
		fmt.Println("service not found")
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
	http.ListenAndServe(":8080", nil)
}
