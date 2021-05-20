package main

import (
	"database/sql"
	"fmt"
	"strconv"
	"time"
)

func dbQuery(cValue string) string {

	rows, err := db.Query(`SELECT un FROM sessionsDB where cValue =?`, cValue)

	check(err)
	defer rows.Close()
	var username string
	for rows.Next() {
		err = rows.Scan(&username)
		check(err)
	}

	if username != "" { //means the user is actually logged in
		updateUserLastTiming(username)
		return username

	} else {
		return ""
	}

}

func updateUserLastTiming(username string) {

	var timeNow string = time.Now().Format("2006-01-02 15:04:05")
	stmt, err := db.Prepare("UPDATE sessionsDB SET lastact=? WHERE un=?")

	check(err)
	defer stmt.Close()

	r, err := stmt.Exec(timeNow, username)
	check(err)

	_, err = r.RowsAffected()

	check(err)

}

func insertSessionsDB(username string, cValue string) {

	var timeNow string = time.Now().Format("2006-01-02 15:04:05")
	stmt, err := db.Prepare("insert into sessionsDB(un, lastact, cValue) values (?,?,?)")
	check(err)
	defer stmt.Close()
	r, err := stmt.Exec(username, timeNow, cValue)
	check(err)
	_, err = r.RowsAffected()

	check(err)

}

func updateSessionsDB(username string, cValue string) {

	var timeNow string = time.Now().Format("2006-01-02 15:04:05")

	stmt, err := db.Prepare("UPDATE sessionsDB SET cValue=? , lastAct=? WHERE un=?")
	check(err)
	defer stmt.Close()
	r, err := stmt.Exec(cValue, timeNow, username)
	check(err)
	n, err := r.RowsAffected()
	_ = n
	check(err)

}

func insertUsersDB(username string, pw []byte, fname, lname, role, uapikey string) {

	stmt, err := db.Prepare("insert into usersDB(un, password, firstName, lastName, userType, apiKey) values (?,?,?,?,?,?)")
	check(err)
	defer stmt.Close()
	r, err := stmt.Exec(username, string(pw), fname, lname, role, uapikey)
	check(err)

	n, err := r.RowsAffected()
	_ = n
	check(err)

}

func queryUsernameUsersDB(un string) bool {

	var e int

	err := db.QueryRow(`SELECT 1 FROM usersDB where un =?`, un).Scan(&e)
	if err == sql.ErrNoRows { //no username found in usersDB
		return false
	} else {
		return true
	}

}

func queryPasswordUsersDB(un string) string {

	var password string

	rows, err := db.Query(`SELECT password FROM usersDB where un =?`, un)
	check(err)
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&password)
		check(err)
	}
	return password
}

//updates user's last search information into user databaes for easy query and marketing
func updateUserLastSearch(searchText, username string) {

	stmt, err := db.Prepare("UPDATE usersDB SET lastSearch=? WHERE un=?")
	check(err)
	defer stmt.Close()

	r, err := stmt.Exec(searchText, username)
	check(err)

	n, err := r.RowsAffected()
	_ = n
	check(err)

}

//inserts user's search & time logs into overall database to find out winning keywords
func insertUserSearchLogs(searchText, username string) {

	var timeNow string = time.Now().Format("2006-01-02 15:04:05")
	stmt, err := db.Prepare("insert into searchLogsDB(un, searches, time) values (?,?,?)")
	check(err)
	defer stmt.Close()
	r, err := stmt.Exec(username, searchText, timeNow)
	check(err)
	n, err := r.RowsAffected()
	_ = n
	check(err)

}

func queryUserLastSearchTerm(un string) string {

	var lastSearchTerm string
	rows, err := db.Query(`SELECT lastSearch FROM usersDB where un =?`, un)
	check(err)
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&lastSearchTerm)
		check(err)
	}

	return lastSearchTerm
}

//inserts user's cart item into the overall database to store what they have added into cart.
func insertItemIntoCart(username, pid string, qty int) {

	var timeNow string = time.Now().Format("2006-01-02 15:04:05")
	stmt, err := db.Prepare("insert into cartDB(un, pid, qty, time) values (?,?,?,?)")
	check(err)
	defer stmt.Close()
	r, err := stmt.Exec(username, pid, qty, timeNow)
	check(err)
	n, err := r.RowsAffected()
	_ = n
	check(err)

}

func queryUserRole(un string) string {

	var userRole string

	rows, err := db.Query(`SELECT userType FROM usersDB where un =?`, un)

	check(err)
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&userRole)
		check(err)
	}

	return userRole
}

//queries all the items in the cart for a particulate username
func queryCartItems(un string) []cartFullData {

	var cartDisplay []cartFullData
	var pid string
	var qty int
	var FoodName string
	var UnitPrice float64

	rows, err := db.Query(`SELECT pid, qty FROM cartDB where un =?`, un)

	check(err)
	defer rows.Close()

	for rows.Next() {
		var totalCost float64
		err = rows.Scan(&pid, &qty)
		check(err)

		// foodname, price := retrieveFoodNameAndPrice(pid)
		foodname, price := retrieveFoodNameAndPrice2(pid)
		FoodName = foodname
		UnitPrice = price
		totalCost += UnitPrice * float64(qty)
		cartDisplay = append(cartDisplay, cartFullData{pid, FoodName, qty, UnitPrice, totalCost, queryUserRole(un)})
	}
	return cartDisplay
}

func queryCartItems2(un string) []cartFullData {

	var cartDisplay []cartFullData
	var pid string
	var qty int
	var FoodName string
	var UnitPrice float64

	rows, err := db.Query(`SELECT pid, qty FROM cartDB where un =?`, un)

	check(err)
	defer rows.Close()

	for rows.Next() {
		var totalCost float64
		err = rows.Scan(&pid, &qty)
		check(err)

		foodname, price := retrieveFoodNameAndPrice2(pid)
		FoodName = foodname
		UnitPrice = price
		totalCost = UnitPrice * float64(qty)
		cartDisplay = append(cartDisplay, cartFullData{pid, toTitle(FoodName), qty, UnitPrice, totalCost, queryUserRole(un)})
	}
	return cartDisplay
}

func emptyCartPing(un string) bool {

	var e int

	err := db.QueryRow(`SELECT 1 FROM cartDB where un =?`, un).Scan(&e)

	if err == sql.ErrNoRows { //no username found in usersDB
		return false
	} else {
		return true
	}

}

//deletes a pid value from a user's cart when the pid's quantity value (rq) == 0
func deleteUserCartItem(un, rp string) {

	stmt, err := db.Prepare("DELETE FROM cartDB where un=? AND pid=?")

	check(err)
	defer stmt.Close()

	r, err := stmt.Exec(un, rp)
	check(err)
	n, err := r.RowsAffected()
	_ = n
	check(err)

}

//update's user's cart pid (food name) quantity
func updateUserCartItem(un, rp string, rq int) {

	stmt, err := db.Prepare("UPDATE cartDB SET qty=? WHERE un=? AND pid=?")
	check(err)
	defer stmt.Close()

	r, err := stmt.Exec(rq, un, rp)
	check(err)

	n, err := r.RowsAffected()
	_ = n
	check(err)

}

func insertCheckoutItemsCheckoutDB(un, generatedID, generatedSysQueueID string, totalCost float64, pi int, foodName string, qty int) {

	var timeNow string = time.Now().Format("2006-01-02 15:04:05")

	stmt, err := db.Prepare("insert into checkoutDB(un, transID, sysID, totalCost, time, pi, foodName, qty) values (?,?,?,?,?,?,?,?)")

	check(err)
	defer stmt.Close()

	r, err := stmt.Exec(un, generatedID, generatedSysQueueID, totalCost, timeNow, pi, foodName, qty)

	check(err)

	n, err := r.RowsAffected()
	_ = n
	check(err)

}

func deletePreviousConfirmation(un string) {

	stmt, err := db.Prepare("DELETE FROM confirmationCheckoutDB where un=?")

	check(err)
	defer stmt.Close()

	r, err := stmt.Exec(un)
	check(err)

	n, err := r.RowsAffected()
	_ = n
	check(err)

}

func clearCart(un string) {

	stmt, err := db.Prepare("DELETE FROM cartDB where un=?")

	check(err)
	defer stmt.Close()

	r, err := stmt.Exec(un)
	check(err)

	n, err := r.RowsAffected()
	_ = n
	check(err)

}

//queries all the items in the confirmationDB.
func queryCheckoutConfirmationItems(un string) []checkoutParseData {

	var confirmationDisplay []checkoutParseData
	var time string
	var transID string
	var sysID string
	var foodName string
	var qty int
	var totalCost float64

	rows, err := db.Query(`SELECT transID, sysID, totalCost, time, foodName, qty FROM confirmationCheckoutDB where un =?`, un)

	check(err)
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&transID, &sysID, &totalCost, &time, &foodName, &qty)
		check(err)

		confirmationDisplay = append(confirmationDisplay, checkoutParseData{time, transID, sysID, toTitle(foodName), qty, totalCost})
	}

	return confirmationDisplay
}

func queryFname(un string) string {

	var firstName string

	rows, err := db.Query(`SELECT firstName FROM usersDB where un =?`, un)

	check(err)
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&firstName)
		check(err)
	}

	return firstName
}

func queryUsername(un string) bool {

	var e int

	err := db.QueryRow(`SELECT 1 FROM cartDB where un =?`, un).Scan(&e)

	if err == sql.ErrNoRows { //no username found in usersDB
		return false
	} else {
		return true
	}

}

func queryAllTransactions() []TransactionsParseData {

	var transactionsDisplay []TransactionsParseData
	var username string
	var time string
	var transID string
	var foodName string
	var qty int
	var totalCost float64

	rows, err := db.Query(`SELECT un, transID, totalCost, time, foodName, qty FROM checkoutDB ORDER BY(time)desc`)

	check(err)
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&username, &transID, &totalCost, &time, &foodName, &qty)
		check(err)

		transactionsDisplay = append(transactionsDisplay, TransactionsParseData{username, time, transID, toTitle(foodName), qty, totalCost})
	}

	return transactionsDisplay
}

//inserts systemID into database
func insertCheckoutItemsSysIDDB(username string, sysID string, pi int) {

	var timeNow string = time.Now().Format("2006-01-02 15:04:05")

	stmt, err := db.Prepare("insert into sysIDDB(un, sysID, pi, time, driverID) values (?,?,?,?,?)")

	check(err)
	defer stmt.Close()

	r, err := stmt.Exec(username, sysID, pi, timeNow, "")

	check(err)

	n, err := r.RowsAffected()
	_ = n
	check(err)

}

func querySystemIDDrivers() []systemQueueParseData {

	var SystemIDDriverDisplay []systemQueueParseData
	var time string
	var sysID string
	var pi int
	var driverID string

	rows, err := db.Query(`SELECT sysID, pi, time, driverID FROM sysIDDB ORDER BY pi desc, time asc`)

	check(err)
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&sysID, &pi, &time, &driverID)
		check(err)

		SystemIDDriverDisplay = append(SystemIDDriverDisplay, systemQueueParseData{sysID, pi, time, driverID})
	}

	return SystemIDDriverDisplay
}

func updateSysIDwDriverID(sysID, driverID string) {

	var timeNow string = time.Now().Format("2006-01-02 15:04:05")

	stmt, err := db.Prepare("UPDATE sysIDDB SET driverID=?, driverDispatchTime=? WHERE sysID=?")
	check(err)
	defer stmt.Close()

	r, err := stmt.Exec(driverID, timeNow, sysID)
	check(err)

	n, err := r.RowsAffected()
	_ = n
	check(err)

}

//deletes a entire sessiosROW from a user's cart when the user logs out.
func deleteSessionsDB(cValue string) {

	stmt, err := db.Prepare("DELETE FROM sessionsDB where cValue=?")

	check(err)
	defer stmt.Close()

	r, err := stmt.Exec(cValue)
	check(err)

	n, err := r.RowsAffected()
	_ = n
	check(err)

}

func updateUserPassword(username string, pn []byte) {

	stmt, err := db.Prepare("UPDATE usersDB SET password=? WHERE un=?")
	check(err)
	defer stmt.Close()

	r, err := stmt.Exec(string(pn), username)
	check(err)

	n, err := r.RowsAffected()
	_ = n
	check(err)

	fmt.Printf("\nusername: %v - Passwordchanged\n", username)

}

// deletes an existing session if the user is already logged in
func removeDuplicateSessionsDB(un string) {

	stmt, err := db.Prepare("DELETE FROM sessionsDB where un=?")

	check(err)
	defer stmt.Close()

	r, err := stmt.Exec(un)
	check(err)
	n, err := r.RowsAffected()
	_ = n
	check(err)

}

func moveCartDBtoCheckoutDB(u, sysID string) {

	stmt, err := db.Prepare("INSERT into checkoutDB(un, pid, qty) SELECT (un, pid, qty) from cartDB")

	check(err)
	defer stmt.Close()

	r, err := stmt.Exec()

	check(err)

	n, err := r.RowsAffected()
	_ = n
	check(err)

}

func checkoutConfirm(u, s string, pi int) {

	// generatedSysQueueID, err := generateSysQueueID()

	var pid string
	var qty int
	var foodName string
	var unitPrice float64
	var totalCost float64

	rows, err := db.Query(`SELECT pid, qty FROM cartDB where un =?`, u)

	check(err)
	defer rows.Close()

	for rows.Next() {
		generatedID, err := generateTransactionID()
		err = rows.Scan(&pid, &qty)

		check(err)
		foodName, unitPrice = retrieveFoodNameAndPrice2(pid)
		totalCost = unitPrice * float64(qty)

		insertIntoCheckoutDB(u, generatedID, s, totalCost, pi, foodName, qty)
		insertCheckoutDispaly(u, generatedID, s, totalCost, foodName, qty)
	}

}

func insertIntoCheckoutDB(u, transID, sysID string, totalCost float64, pi int, foodName string, qty int) {

	var timeNow string = time.Now().Format("2006-01-02 15:04:05")

	stmt, err := db.Prepare("insert into checkoutDB(un, transID, sysID, totalCost, time, pi, foodName, qty) values (?,?,?,?,?,?,?,?)")

	check(err)
	defer stmt.Close()

	r, err := stmt.Exec(u, transID, sysID, totalCost, timeNow, pi, foodName, qty)

	check(err)

	n, err := r.RowsAffected()
	_ = n
	check(err)
}

//inserts items into confrimatinCheckoutDB
func insertCheckoutDispaly(un, generatedID, generatedSysQueueID string, totalCost float64, foodName string, qty int) {

	var timeNow string = time.Now().Format("2006-01-02 15:04:05")

	stmt, err := db.Prepare("insert into confirmationCheckoutDB(un, transID, sysID, totalCost, time, foodName, qty) values (?,?,?,?,?,?,?)")

	check(err)
	defer stmt.Close()

	r, err := stmt.Exec(un, generatedID, generatedSysQueueID, totalCost, timeNow, foodName, qty)

	check(err)

	n, err := r.RowsAffected()
	_ = n
	check(err)

}

func retrieveUserApiKey(un string) string {

	var apiKey string

	rows, err := db.Query(`SELECT apiKey FROM usersDB where un =?`, un)

	check(err)
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&apiKey)
		check(err)
	}
	return apiKey
}

func validateAPIkey(s string) bool {

	var e int

	err := db.QueryRow(`SELECT 1 FROM usersDB where apiKey =?`, s).Scan(&e)

	if err == sql.ErrNoRows { //no username found in usersDB
		return false
	} else {
		return true
	}

}

func updateApiKeyUsersDB(s1, s2 string) {

	stmt, err := db.Prepare("UPDATE usersDB SET apiKey=? WHERE un=?")

	check(err)
	defer stmt.Close()

	r, err := stmt.Exec(s2, s1)
	check(err)

	_, err = r.RowsAffected()

	check(err)

}

func retriAPIUsername(s string) string {

	rows, err := db.Query(`SELECT un FROM usersDB where apiKey =?`, s)

	check(err)
	defer rows.Close()
	var username string

	for rows.Next() {
		err = rows.Scan(&username)
		check(err)
	}

	return username
}

func retriMercInfoAPI(s string) (string, string) {

	rows, err := db.Query(`SELECT merchantName, detailedLocation FROM merchantsDB where username =? LIMIT 1`, s)

	check(err)
	defer rows.Close()
	var merchantName string
	var detailedLocation string

	for rows.Next() {
		err = rows.Scan(&merchantName, &detailedLocation)
		check(err)
	}

	return toTitle(merchantName), toTitle(detailedLocation)
}

func checkifFoodExists(s1, s2 string) bool {

	var e int
	err := db.QueryRow(`SELECT 1 FROM foodListDB where username =? and foodName =?  `, s1, s2).Scan(&e)
	if err == sql.ErrNoRows { //no username found in usersDB
		return false
	} else {
		return true
	}

}

func retriMerchFooditems(s string) []MerchantFoodInfo {

	var result []MerchantFoodInfo
	var merchantName string
	var price string

	rows, err := db.Query(`SELECT foodName, price FROM foodListDB where username =?`, s)

	check(err)
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&merchantName, &price)
		check(err)
		result = append(result, MerchantFoodInfo{merchantName, price})
	}

	return result
}

// username, data.OldFoodname, data.NewFoodname, data.Price
func updatefoodListDB(s1, s2, s3, s4 string) {

	stmt, err := db.Prepare("UPDATE foodListDB SET lastact=? WHERE un=?")

	check(err)
	defer stmt.Close()

	r, err := stmt.Exec()
	check(err)

	_, err = r.RowsAffected()

	check(err)

}

func retrieveFoodNameAndPriceDB(s1, s2 string) string {

	var price string
	rows, err := db.Query(`SELECT price FROM foodListDB where username =? and foodName=? `, s1, s2)
	check(err)
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&price)
		check(err)
	}

	return price
}

func updatePriceofItem(s1, s2, s3 string) {

	stmt, err := db.Prepare(`UPDATE foodListDB SET price=? WHERE foodName=? and username=?`)

	check(err)
	defer stmt.Close()
	r, err := stmt.Exec(s2, s1, s3)
	check(err)
	_, err = r.RowsAffected()

	check(err)

}

func updateNameofItem(s1, s2, s3 string) {

	stmt, err := db.Prepare(`UPDATE foodListDB SET foodName=? WHERE foodName=? and username=?`)

	check(err)
	defer stmt.Close()
	r, err := stmt.Exec(s2, s1, s3)
	check(err)
	_, err = r.RowsAffected()
	check(err)

}

func updateNamePriceofItem(s1, s2, s3, s4 string) {

	stmt, err := db.Prepare(`UPDATE foodListDB SET foodName=?, price=? WHERE foodName=? and username=?`)

	check(err)
	defer stmt.Close()
	r, err := stmt.Exec(s2, s3, s1, s4)
	check(err)
	_, err = r.RowsAffected()

	check(err)

}

var FoodMerchantBrandNames []string //global variable

func createFoodList() {

	// var result []string
	var FoodName string
	var BrandName string

	rows, err := db.Query(`SELECT foodName, merchantName FROM foodListDB`)
	check(err)
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&FoodName, &BrandName)
		check(err)
		FoodMerchantBrandNames = append(FoodMerchantBrandNames, FoodName+" "+"-"+" "+BrandName)
	}

	// fmt.Printf("\n%+v\n", result)

} //end createFoodList

func updatepid(s string, i int) {

	result := strconv.Itoa(countpid)
	var pid string
	pid = "KVPID" + result

	stmt, err := db.Prepare(`update foodListDB SET pid=? where index1=? and username=?`)
	check(err)
	defer stmt.Close()
	r, err := stmt.Exec(pid, i, s)
	check(err)
	_, err = r.RowsAffected()
	check(err)
	countpid++
	// fmt.Println(pid, countpid)

}

func retrivePIDvalue(s1, s2 string) string {

	var pid string

	rows, err := db.Query(`SELECT pid from foodListDB where foodName=? and merchantName=?`, s1, s2)
	check(err)
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&pid)
		check(err)
	}
	return pid
}

func retrieveFoodNameAndPrice2(s string) (string, float64) {

	var foodName string
	var unitPrice float64

	rows, err := db.Query(`SELECT foodName, price from foodListDB where pid=?`, s)
	check(err)
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&foodName, &unitPrice)
		check(err)
	}
	return foodName, unitPrice
}

func addNewFoodItems(s1, s2, s3 string) {
	merchantInfo := retrieveMerchantDetailedInformation(s1)
	amendFoodListDB("add", s1, s2, s3, merchantInfo)
}

func retrieveMerchantDetailedInformation(s string) []string {

	var merchantName string
	var detailedLocation string
	var postalCode string
	var monWH string
	var tuesWH string
	var wedWH string
	var thursWH string
	var friWH string
	var satWH string
	var sunWH string
	var phWH string
	var cot string
	var result []string

	rows, err := db.Query(`SELECT merchantName, detailedLocation, postalCode, monWH, tuesWH, wedWH, thursWH, friWH, satWH, sunWH, phWH, cot from merchantsDB where username=? LIMIT 1`, s)
	check(err)
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&merchantName, &detailedLocation, &postalCode, &monWH, &tuesWH, &wedWH, &thursWH, &friWH, &satWH, &sunWH, &phWH, &cot)
		check(err)
	}

	result = append(result, merchantName, detailedLocation, postalCode, monWH, tuesWH, wedWH, thursWH, friWH, satWH, sunWH, phWH, cot)
	return result
}

func amendFoodListDB(fnType, s1, s2, s3 string, s4 []string) { //s2 is username

	mutex.Lock()
	{
		upid := retriveLastPIDvalue()
		pid := processPIDvalue(upid)

		if fnType == "add" {
			addMerchantFoodInformation(pid, s1, s2, s3, s4)
		} else {
			deleteMerchantFoodInformation(s2, s3)
		}

	}
	mutex.Unlock()
}

// pid, username, foodname, pricename, merchant informatioin
func addMerchantFoodInformation(pid, s1, s2, s3 string, s4 []string) {

	stmt, err := db.Prepare("insert into foodListDB(username, pid, foodName, price, merchantName, detailedLocation, postalCode, monWH, tuesWH, wedWH, thursWH, friWH, satWH, sunWH, phWH, cot) values (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)")
	check(err)
	defer stmt.Close()
	r, err := stmt.Exec(s1, pid, s2, s3, s4[0], s4[1], s4[2], s4[3], s4[4], s4[5], s4[6], s4[7], s4[8], s4[9], s4[10], s4[11])
	check(err)

	n, err := r.RowsAffected()
	_ = n
	check(err)

}

func retriveLastPIDvalue() string {

	var pid string

	rows, err := db.Query(`SELECT pid from foodListDB ORDER BY id DESC LIMIT 1`)
	check(err)
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&pid)
		check(err)
	}
	return pid
}

func deleteMerchantFoodInformation(s1, s2 string) {

	stmt, err := db.Prepare("DELETE FROM foodListDB where username=? AND foodName=?")

	check(err)
	defer stmt.Close()

	r, err := stmt.Exec(s1, s2)
	check(err)
	n, err := r.RowsAffected()
	_ = n
	check(err)

}

func insertMerchantInformationDB(s1, s2, s3, s4, s5, s6, s7, s8, s9, s10, s11, s12, s13 string) {

	stmt, err := db.Prepare("insert into merchantsDB(username, merchantName, detailedLocation, postalCode, monWH, tuesWH, wedWH, thursWH, friWH, satWH, sunWH, phWH, cot) values (?,?,?,?,?,?,?,?,?,?,?,?,?)")
	check(err)
	defer stmt.Close()
	r, err := stmt.Exec(s1, s2, s3, s4, s5, s6, s7, s8, s9, s10, s11, s12, s13)
	check(err)

	n, err := r.RowsAffected()
	_ = n
	check(err)

}
