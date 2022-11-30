package main

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
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
	var outputBuffer bytes.Buffer
	err = tmpl.Execute(&outputBuffer, config)
	if err != nil {
		panic(err)
	}
	outputFilename := "standalone.xml"
	if !isValid(outputBuffer.String()) {
		outputFilename += ".invalid"
		log.Fatalf("invalid xml, saving to %s for inspection\n", outputFilename)
	} else {
		log.Printf("writing xml to %s\n", outputFilename)
	}
	if err := os.WriteFile(outputFilename, outputBuffer.Bytes(), 0644); err != nil {
		log.Panicf("error writing to %s: %s", outputFilename, err)
	}
}

func isValid(s string) bool {
	return xml.Unmarshal([]byte(s), new(interface{})) == nil
}
