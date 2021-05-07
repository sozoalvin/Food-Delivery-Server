package main

import (
	"database/sql"
	"fmt"
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

func insertUsersDB(username string, pw []byte, fname string, lname string, role string) {

	stmt, err := db.Prepare("insert into usersDB(un, password, firstName, lastName, userType) values (?,?,?,?,?)")

	check(err)
	defer stmt.Close()

	r, err := stmt.Exec(username, string(pw), fname, lname, role)

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

		foodname, price := retrieveFoodNameAndPrice(pid)
		FoodName = foodname
		UnitPrice = price
		totalCost += UnitPrice * float64(qty)

		cartDisplay = append(cartDisplay, cartFullData{pid, FoodName, qty, UnitPrice, totalCost, queryUserRole(un)})

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

		confirmationDisplay = append(confirmationDisplay, checkoutParseData{time, transID, sysID, foodName, qty, totalCost})
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

	var username string
	var transactionsDisplay []TransactionsParseData
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

		transactionsDisplay = append(transactionsDisplay, TransactionsParseData{username, time, transID, foodName, qty, totalCost})
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

func checkoutConfirm(u string, pi int) {

	generatedSysQueueID, err := generateSysQueueID()

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
		foodName, unitPrice = retrieveFoodNameAndPrice(pid)
		totalCost = unitPrice * float64(qty)

		insertIntoCheckoutDB(u, generatedID, generatedSysQueueID, totalCost, pi, foodName, qty)
		insertCheckoutDispaly(u, generatedID, generatedSysQueueID, totalCost, foodName, qty)
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
