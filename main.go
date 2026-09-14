package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/rivo/tview"
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
	data, err := os.ReadFile(inventoryFile)
	if err != nil {
		if os.IsNotExist(err) {
			return
		}
		log.Fatal("error reading inventory file - ", err)
	}
	if err := json.Unmarshal(data, &inventory); err != nil {
		log.Fatal("error parsing inventory file - ", err)
	}
}

func saveInventory() {
	data, err := json.MarshalIndent(inventory, "", "  ")
	if err != nil {
		log.Fatal("error saving inventory - ", err)
	}
	if err := os.WriteFile(inventoryFile, data, 0644); err != nil {
		log.Fatal("error writing inventory file - ", err)
	}
}

func deleteInventory(index int) bool {
	if index < 0 || index >= len(inventory) {
		return false
	}
	inventory = append(inventory[:index], inventory[index+1:]...)
	saveInventory()
	return true
}

func refreshInventory(inventoryList *tview.List) {
	inventoryList.Clear()
	for _, item := range inventory {
		inventoryList.AddItem(fmt.Sprintf("%s - %d items", item.Name, item.Stock), "", 0, nil)
	}
}

func main() {
	app := tview.NewApplication()
	loadInventory()

	pages := tview.NewPages()

	showModal := func(message string) {
		modal := tview.NewModal().
			SetText(message).
			AddButtons([]string{"OK"}).
			SetDoneFunc(func(buttonIndex int, buttonLabel string) {
				pages.RemovePage("modal")
			})
		pages.AddPage("modal", modal, true, true)
	}

	inventoryList := tview.NewList().ShowSecondaryText(false)
	inventoryList.SetBorder(true).SetTitle("Inventory Items")

	nameInput := tview.NewInputField().
		SetLabel("Name: ").
		SetFieldWidth(20)
	stockInput := tview.NewInputField().
		SetLabel("Stock: ").
		SetFieldWidth(20)

	addButton := tview.NewButton("Add Item").
		SetSelectedFunc(func() {
			if nameInput.GetText() == "" {
				showModal("Name cannot be empty")
				return
			}

			stock := 0
			if stockStr := stockInput.GetText(); stockStr != "" {
				var err error
				if stock, err = strconv.Atoi(stockStr); err != nil || stock <= 0 {
					showModal("Invalid stock value. Please enter a positive integer.")
					return
				}
			}

			inventory = append(inventory, Item{Name: nameInput.GetText(), Stock: stock})
			saveInventory()

			nameInput.SetText("")
			stockInput.SetText("")
			refreshInventory(inventoryList)
		})

	deleteButton := tview.NewButton("Delete Selected").
		SetSelectedFunc(func() {
			if !deleteInventory(inventoryList.GetCurrentItem()) {
				showModal("No item selected for deletion")
				return
			}
			refreshInventory(inventoryList)
		})

	quitButton := tview.NewButton("Quit").
		SetSelectedFunc(func() {
			app.Stop()
		})

	inputFlex := tview.NewFlex().
		SetDirection(tview.FlexColumn).
		AddItem(nameInput, 0, 1, false).
		AddItem(stockInput, 0, 1, false).
		AddItem(addButton, 0, 1, false)

	buttonFlex := tview.NewFlex().
		SetDirection(tview.FlexColumn).
		AddItem(deleteButton, 0, 1, false).
		AddItem(quitButton, 0, 1, false)

	flex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(inventoryList, 0, 1, true).
		AddItem(inputFlex, 3, 1, false).
		AddItem(buttonFlex, 3, 1, false)

	pages.AddPage("main", flex, true, true)
	app.SetRoot(pages, true)

	refreshInventory(inventoryList)

	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
