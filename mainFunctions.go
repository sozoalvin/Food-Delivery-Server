package main

import (
	"fmt"
	"sort"
	"strconv"
)

var FoodMerchantNameAddress []string
var FoodMerchantNameAddress1 []string

func concatenateFoodList() { //return you a slice of strings
	fmt.Println(V)
}

func CreateFoodList(ch chan string) {
	// func CreateFoodList() []string {

	for _, v := range V {
		FoodMerchantNameAddress = append(FoodMerchantNameAddress, v.FoodName+" - "+v.MerchantName+" - "+v.DetailedLocation)
	}
	sort.Strings(FoodMerchantNameAddress)
	ch <- "Mandatory - Food List Data Generated"
}

func generateTransactionID() (string, error) {

	var generatedID string
	//for every element wihch is stored based on PID as key, should have unique generated ID

	mutex.Lock()
	{
		sTransID := strconv.Itoa(TransID)

		generatedID = "MC" + sTransID + "KV" // generated ID will always be unique

		TransID++
	}
	mutex.Unlock()

	return generatedID, nil

}

func generateSysQueueID() (string, error) {

	var generatedSysQueueID string

	mutex.Lock()
	{
		SQueueID := strconv.Itoa(QueueID)

		generatedSysQueueID = "OS" + SQueueID + "KV" // generated ID will always be unique

		QueueID++
	}
	mutex.Unlock()

	return generatedSysQueueID, nil
}

func retrieveFoodNameAndPrice(pid string) (string, float64) {

	foodName := foodNameAddresstoname[pid] //foodNameAddresstname = global MAP

	unitCost := MyFoodListMap[foodName].Price

	return foodName, unitCost

}
