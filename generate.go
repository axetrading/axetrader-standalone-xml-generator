package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"text/template"
)

type Configuration struct {
}

func main() {
	content, err := ioutil.ReadFile("./config.json")
	if err != nil {
		log.Fatal("Error when opening file: ", err)
	}

	config := Configuration{}
	err = json.Unmarshal(content, &config)
	if err != nil {
		log.Fatal("Error during Unmarshal(): ", err)
	}

	tmpl, err := template.New("standalone.xml.template").ParseFiles("standalone.xml.template")
	if err != nil {
		panic(err)
	}
	err = tmpl.Execute(os.Stdout, config)
	if err != nil {
		panic(err)
	}
}
