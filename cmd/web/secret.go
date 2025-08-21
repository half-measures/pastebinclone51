package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

type Config struct {
	SQLSecret string `json:"sqlsecret"`
}

func getSQLSecret() string {
	file, err := ioutil.ReadFile("config.json")
	if err != nil {
		log.Fatalf("Failed to read config file for SQL - %v", err)
		os.Exit(1)

	}
	//unmarshal json
	var config Config
	//defer configfile.Close()

	err = json.Unmarshal(file, &config)
	if err != nil {
		//log error
		log.Fatalf("Failed to unmarshal JSON: %v", err)
		os.Exit(1)
	}
	//pass now in config.SQLSecret var
	fmt.Printf("SQL password loaded %s\n", config.SQLSecret)
	//sqlsecret := bytevalue
	return config.SQLSecret

}
