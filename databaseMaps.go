package main

import (
	"strconv"
)

var MyFoodListMap = make(map[string]FoodInfo)
var MyFoodListMap2 = make(map[string]FoodInfo)
var foodNameAddresstoID = make(map[string]string)
var foodNameAddresstoname = make(map[string]string)

func CreateFoodListMap(ch chan string) {

	for _, v := range V {
		keyValue := v.FoodName + " - " + v.MerchantName + " - " + v.DetailedLocation
		MyFoodListMap[keyValue] = FoodInfo{v.FoodName, v.MerchantName, v.DetailedLocation, v.PostalCode, v.Price, v.OpeningPeriods}
	}
	ch <- "Food List Map Data Completed"
}

func FoodMerchantNameAddressProductID() {

	for _, v := range FoodMerchantNameAddress { //range through entire SLICE to populate map
		value1 := "pID" + strconv.Itoa(pid)
		foodNameAddresstoID[v] = value1
		foodNameAddresstoname[value1] = v
		pid++
	}
}
