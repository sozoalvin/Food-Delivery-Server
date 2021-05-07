package main

import (
	"fmt"
	"net/http"
	"time"

	uuid "github.com/satori/go.uuid"
)

// func getUser(w http.ResponseWriter, req *http.Request) *user {
func getUser(w http.ResponseWriter, req *http.Request) string {

	var username string

	c, err := req.Cookie("session") //requesting to check if the client has a cookie, named session

	if err != nil { //fired if no cookie named, session is present

		sID, err := uuid.NewV4() //if no cookie present, we give it one.
		//err handling
		if err != nil {
			fmt.Printf("Something went wrong: %s, err") //prints if there is error when have error generating UUID
		}
		c = &http.Cookie{
			Name:  "session",    //nametype of cookie
			Value: sID.String(), // Returns canonical string representation of UUID:
		}
	} // if cookie not present; all above codes will run.
	http.SetCookie(w, c) // http.SetCookie is required to 'set

	rows, err := db.Query(`SELECT un FROM sessionsDB where cValue =?`, c.Value)

	check(err)
	defer rows.Close()

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

func alreadyLoggedIn(w http.ResponseWriter, req *http.Request) bool { //checks if a user is logged or by returning equivalent bool value
	c, err := req.Cookie("session")
	if err != nil {
		return false
	}

	username := dbQuery(c.Value)
	if username == "" {
		//check if username exists on other sessions.
		return false
	} else {
		return true
	}

}

func cleanSessions() {
	fmt.Println("(before) db session cleaned") // for demonstration purposes
	showSessions()                             // for demonstration purposes
	for k, v := range dbSessions {
		if time.Now().Sub(v.lastActivity) > (time.Second * 30) {
			delete(dbSessions, k)
		}
	}
	dbSessionsCleaned = time.Now()
	fmt.Println("(after) db session cleaned") // for demonstration purposes
	showSessions()                            // for demonstration purposes
}

func showSessions() {
	fmt.Println("********")
	for k, v := range dbSessions {
		fmt.Println(k, v.un)
	}
	fmt.Println("")
}
