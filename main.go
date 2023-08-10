package main

import (
	"log"
	"encoding/json"
	"os"
    "fmt"
)

func main() {
	configFile, err := os.Open("config.json")
	if err != nil {
		log.Fatal("Error opening config file:", err)
	}
	defer configFile.Close()

	var configuration Configuration
	err = json.NewDecoder(configFile).Decode(&configuration)
	if err != nil {
		log.Fatal("Error reading config file:", err)
	}

	// to implement db
	// store, err := NewMysqlStore(config)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// if err := store.Init(); err != nil {
	// 	log.Fatal(err)
	// }

	server, err := NewAPIServer(":4000", configuration)
    if err != nil {
		fmt.Println("Error creating APIServer:", err)
		return
	}

	server.Run()
}