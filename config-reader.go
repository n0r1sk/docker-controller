/*
Copyright 2017 Mario Kleinsasser and Bernhard Rausch

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	//"fmt"
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v2"
)

type Tcfg struct {
	General struct {
		Api_key string
	}
}

func ReadConfigfile() (config Tcfg, ok bool) {
	cfgdata, err := ioutil.ReadFile("/config/docker-controller.yml")

	if err != nil {
		log.Panic("Cannot open config file from /config/docker-controller.yml")
	}

	t := Tcfg{}

	err = yaml.Unmarshal([]byte(cfgdata), &t)
	if err != nil {
		log.Panic("Cannot map yml config file to interface, possible syntax error")
		log.Panic(err)
	}

	return t, true
}
