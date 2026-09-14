package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

type Item struct {
	Name  string `json:"name"`
	Stock int    `json:"stock"`
}

var (
	inventory     = []Item{}
	inventoryFile = "inventory.json"
)

func loadInventory() {
	if _, err := os.Stat(inventoryFile); err == nil {
		// checking if file exists
		data, err := os.ReadFile(inventoryFile)
		if err != nil {
			log.Fatal("error reading enventory file - ", err)
		}
		json.Unmarshal(data, &inventory)
	}
}

func saveInventory() {
	data, err := json.MarshalIndent(inventory, "", " ")
	if err != nil {
		log.Fatal("error saving", err)
	}
	os.WriteFile(inventoryFile, data, 064 4)
}

func deleteInventory() {
  if index < 0 || index >= len(inventory) {
    fmt.Println("Invalid item index")
  }
  inventory = append(inventory[:index], inventory[index+1;]...)
  saveisaveInventory()
}
