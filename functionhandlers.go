package main

import (
	"fmt"
	"html"
	"net/http"
	"regexp"
	"strconv"
	"time"

	uuid "github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
)

func index(w http.ResponseWriter, req *http.Request) {

	u := getUser(w, req)

	if req.Method == http.MethodPost {
		searchKW := req.FormValue("searchtext")

		if u == "" { //no username present!
			http.Redirect(w, req, "/login_redirect", http.StatusSeeOther)
			return
		}

		sr := html.EscapeString(searchKW)

		updateUserLastSearch(sr, u)
		insertUserSearchLogs(sr, u)

		http.Redirect(w, req, "/searchresult", http.StatusSeeOther)
	}

	parseData := struct {
		U    string
		Data string
	}{
		u, "",
	}

	showSessions() // for demonstration purposes
	tpl.ExecuteTemplate(w, "index.gohtml", parseData)
}

func signup(w http.ResponseWriter, req *http.Request) {
	if alreadyLoggedIn(w, req) {
		http.Redirect(w, req, "/", http.StatusSeeOther)
		return
	}

	var me = make(map[string]string) //make map for error

	if req.Method == http.MethodPost {
		// get form values
		unUnsanitized := req.FormValue("username")
		p := req.FormValue("password")
		fUnsanitized := req.FormValue("firstname")
		lUnsanitized := req.FormValue("lastname")
		r := req.FormValue("role")

		un := html.EscapeString(unUnsanitized)
		f := html.EscapeString(fUnsanitized)
		l := html.EscapeString(lUnsanitized)

		//validateInputs perfom input sanitsation
		boolresult, mapresult := validateInputs(un, p, f, l, me)

		if boolresult == false {
			tpl.ExecuteTemplate(w, "signup.gohtml", mapresult)
			return
		}

		if queryUsername(un) {
			mapresult["Username"] = "Username already exists!"
			tpl.ExecuteTemplate(w, "signup.gohtml", mapresult)
			return
		}

		sID, err := uuid.NewV4()
		//err handling
		if err != nil {
			fmt.Printf("Something went wrong: %s, err")
		}

		c := &http.Cookie{
			Name:     "session",
			Value:    sID.String(),
			HttpOnly: true,
		}

		http.SetCookie(w, c)

		dbSessions[c.Value] = session{un, time.Now()} // i wil store your informtion with cookie value UUID

		insertSessionsDB(un, c.Value)

		// store user in dbUsers
		bs, err := bcrypt.GenerateFromPassword([]byte(p), bcrypt.MinCost)
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		insertUsersDB(un, bs, f, l, r)

		http.Redirect(w, req, "/", http.StatusSeeOther) //once logged in, redirect to where you want the user to be redirected to
		return
	}
	showSessions() // for demonstration purposes

	tpl.ExecuteTemplate(w, "signup.gohtml", nil)
}

func login(w http.ResponseWriter, req *http.Request) {
	if alreadyLoggedIn(w, req) {
		http.Redirect(w, req, "/", http.StatusSeeOther) //if alreadyLoggedIn == true -> returns them to see what they're supposed to see etc.
		return
	}

	var me = make(map[string]string) //make map for error
	rx := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

	// process form submission
	if req.Method == http.MethodPost {

		unUnsanitized := req.FormValue("username") //the 'name' of the field
		p := req.FormValue("password")             //the 'name' of the field
		un := html.EscapeString(unUnsanitized)

		if !rx.MatchString(un) || len(un) > 20 {
			me["Username1"] = "Username entered is not a valid email address."
			tpl.ExecuteTemplate(w, "login.gohtml", me)
			return
		}
		//if username is correct, then we check password
		dbPassword := queryPasswordUsersDB(un)

		// err := bcrypt.CompareHashAndPassword(u.Password, []byte(p))
		err := bcrypt.CompareHashAndPassword([]byte(dbPassword), []byte(p))
		if err != nil {

			me["Password"] = "Invaid Password Entered. Please try again"
			tpl.ExecuteTemplate(w, "login.gohtml", me)

			return
		}

		sID, err := uuid.NewV4()
		//err handling
		if err != nil {
			fmt.Printf("Something went wrong: %s, err")
		}

		c := &http.Cookie{
			Name:     "session",
			Value:    sID.String(),
			HttpOnly: true,
		}

		http.SetCookie(w, c)

		//check for duplicate sessions and kill it. this forces the other session to be logged out
		removeDuplicateSessionsDB(un)
		// updateSessionsDB(un, c.Value)
		insertSessionsDB(un, c.Value)
		// dbSessions[c.Value] = session{un, time.Now()}
		http.Redirect(w, req, "/", http.StatusSeeOther)
		return
	}
	showSessions() // for demonstration purposes
	tpl.ExecuteTemplate(w, "login.gohtml", nil)
}

func logout(w http.ResponseWriter, req *http.Request) {
	if !alreadyLoggedIn(w, req) { //if you are not logged in, there's nothing you need to do. whatever UUID cookie value, belongs to a non-logged in visitor
		http.Redirect(w, req, "/", http.StatusSeeOther)
		return
	}
	c, _ := req.Cookie("session")

	deleteSessionsDB(c.Value)

	c = &http.Cookie{
		Name:   "session",
		Value:  "",
		MaxAge: -1,
	}
	http.SetCookie(w, c)

	// clean up dbSessions
	if time.Now().Sub(dbSessionsCleaned) > (time.Second * 30) {
		go cleanSessions()
	}

	http.Redirect(w, req, "/", http.StatusSeeOther) //goes back to home page after logging out
}

func rs2(w http.ResponseWriter, req *http.Request) {
	http.ServeFile(w, req, "assets/rs2.png")
}

func searchresult(w http.ResponseWriter, req *http.Request) {

	u := getUser(w, req) //getUser function call

	var localSlice = []searchResultFormat{}

	lastSearchTerm := queryUserLastSearchTerm(u)

	localSearchResult := MyFoodListDB.GetSuggestion(lastSearchTerm, 50) // you will always append a global variable so you pass data this way.

	for _, v := range localSearchResult { //range through all available data in the slice

		valuepair := foodNameAddresstoID[v]
		localSlice = append(localSlice, searchResultFormat{v, valuepair}) //everytime a new item is added into cart, this gets appended
	}

	if req.Method == http.MethodPost {
		productid := req.FormValue("pid") //pid is also known as the productID
		quantity1 := req.FormValue("quantity")

		intQuantity, err := strconv.Atoi(quantity1)
		if err != nil || intQuantity <= 0 {
			intQuantity = 1 //if error, we default quantity to 1.
		}
		insertItemIntoCart(u, productid, intQuantity)

		http.Redirect(w, req, "/yourcart", http.StatusSeeOther)
	}

	showSessions() // for demonstration purposes

	parseData := struct {
		U    string
		Data []searchResultFormat
	}{
		u, localSlice,
	}

	tpl.ExecuteTemplate(w, "searchresult.gohtml", parseData)

}

func yourcart(w http.ResponseWriter, req *http.Request) {

	var localCartDisplay []cartFullData
	u := getUser(w, req) //getUser function call

	if u == "" {
		http.Redirect(w, req, "/login_redirect", http.StatusSeeOther)
		return
	}
	showSessions() // for demonstration purposes
	if req.Method == http.MethodPost {
		rb := req.FormValue("updatecart")
		if rb == "updatecart" {

			rp := req.FormValue("pid")
			rq := req.FormValue("quantity")

			if rq == "0" { //if quantity is 0, delete key (pid) from map - CartMapData
				deleteUserCartItem(u, rp)
				http.Redirect(w, req, "/yourcart", http.StatusSeeOther)
				return
			}

			irq, err := strconv.Atoi(rq) //conversion of string rq to integer

			if err != nil || irq < 0 {
				irq = 1 //if error, we default quantity to 1.
			}
			//updataes user cart items.
			updateUserCartItem(u, rp, irq)
			http.Redirect(w, req, "/yourcart", http.StatusSeeOther)

			return
		} else {

			pi, piErr := strconv.Atoi(req.FormValue("priorityindex"))
			if piErr != nil || pi < 0 {
				pi = 0 //if error, we default quantity to 1.
			}

			generatedSysQueueID, err := generateSysQueueID()

			if err != nil {
				http.Error(w, "Internal server error", http.StatusInternalServerError)
				return
			}
			deletePreviousConfirmation(u)

			checkoutConfirm(u, pi)

			insertCheckoutItemsSysIDDB(u, generatedSysQueueID, pi)

			http.Redirect(w, req, "/checkout_processing", http.StatusSeeOther)
			return
		}
	}

	localCartDisplay = queryCartItems(u)

	parseData := struct {
		U    string
		Data []cartFullData
	}{
		u, localCartDisplay,
	}

	if queryUserRole(u) == "Customer Service Officer" || queryUserRole(u) == "superuser#1" {

		tpl.ExecuteTemplate(w, "yourcart_admin.gohtml", parseData)
	} else {

		tpl.ExecuteTemplate(w, "yourcart.gohtml", parseData)
	}
}

func checkout_processing(w http.ResponseWriter, req *http.Request) {

	u := getUser(w, req) //getUser function call

	firstName := queryFname(u)

	customer := &user{
		Username: u,
		Fname:    firstName,
	}

	parseData := struct {
		U    string
		Data *user
	}{
		u, customer,
	}

	showSessions() // for demonstration purposes

	tpl.ExecuteTemplate(w, "checkout_processing.gohtml", parseData)

}

func checkout(w http.ResponseWriter, req *http.Request) {

	u := getUser(w, req) //getUser function call

	showSessions() // for demonstration purposes

	data := queryCheckoutConfirmationItems(u)

	fmt.Printf("\n%+v\n", data)

	if req.Method == http.MethodPost {
		//no username present!
		r := req.FormValue("homebutton")
		if r == "home" {
			http.Redirect(w, req, "/", http.StatusSeeOther)
		}

	}

	parseData := struct {
		U    string
		Data []checkoutParseData
	}{
		u, data,
	}

	tpl.ExecuteTemplate(w, "checkout.gohtml", parseData)
}

func allsystemorders(w http.ResponseWriter, req *http.Request) {

	u := getUser(w, req) //getUser function call

	showSessions() // for demonstration purposes

	if u == "" {
		//no username present!
		http.Redirect(w, req, "/allthefoodisgone", http.StatusSeeOther)
		return
	}

	if queryUserRole(u) == "Customer Service Officer" || queryUserRole(u) == "superuser#1" {

		tpl.ExecuteTemplate(w, "allsystemorders.gohtml", nil) //replace nil as data

	} else {

		tpl.ExecuteTemplate(w, "allthefoodisgone.gohtml", nil)
	}
}

func alltransactions(w http.ResponseWriter, req *http.Request) {

	u := getUser(w, req) //getUser function call

	showSessions() // for demonstration purposes

	if u == "" {
		//no username present!
		http.Redirect(w, req, "/allthefoodisgone", http.StatusSeeOther)
		return
	}

	data := queryAllTransactions()

	parseData := struct {
		U    string
		Data []TransactionsParseData
	}{
		u, data,
	}

	if queryUserRole(u) == "Customer Service Officer" || queryUserRole(u) == "superuser#1" {

		tpl.ExecuteTemplate(w, "alltransactions.gohtml", parseData) //please replace nil as data

	} else {

		tpl.ExecuteTemplate(w, "allthefoodisgone.gohtml", nil)
	}

}

func login_redirect(w http.ResponseWriter, req *http.Request) {

	showSessions() // for demonstration purposes

	u := getUser(w, req)

	parseData := struct {
		U    string
		Data string
	}{
		u, "",
	}

	tpl.ExecuteTemplate(w, "login_redirect.gohtml", parseData)

}

func clearcart(w http.ResponseWriter, req *http.Request) {

	u := getUser(w, req) //getUser function call

	if u == "" {
		//no username present!
		http.Redirect(w, req, "/login_redirect", http.StatusSeeOther)
		return
	}

	showSessions() // for demonstration purposes

	clearCart(u)

	http.Redirect(w, req, "/yourcart", http.StatusSeeOther)
	return

}

func validateInputs(un string, p string, f string, l string, me map[string]string) (bool, map[string]string) {

	rx := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

	if len(un) == 0 || len(un) > 40 {
		me["Username"] = "Username is not valid. Please Enter again"
	} else if !rx.MatchString(un) {
		me["Username"] = "Username is not valid email address. Please Enter again"
	}
	if len(p) == 0 {
		me["Password"] = "Password is not valid. Please Enter again"
	}
	if len(f) == 0 || len(f) > 45 {
		me["FirstName"] = "First Name is not valid. Please Enter again"
	}
	if len(l) == 0 || len(l) > 45 {
		me["LastName"] = "Last Name is not valid. Please Enter again"
	}
	if len(un) != 0 && len(p) != 0 && len(un) != 0 && len(l) != 0 && rx.MatchString(un) && len(un) < 40 && len(f) < 45 && len(l) < 45 {
		return true, me
	}
	return false, me
}

func dispatchdriver(w http.ResponseWriter, req *http.Request) {

	u := getUser(w, req) //getUser function call
	if u == "" {
		http.Redirect(w, req, "/allthefoodisgone", http.StatusSeeOther)
		return
	}

	if req.Method == http.MethodPost {
		rb := req.FormValue("updatedriver")
		if rb == "updatedriver" {

			rsq := req.FormValue("systemqueuenumber")   //request system queue number
			rdn_noescape := req.FormValue("drivername") //request assigned driver name
			rdn := html.EscapeString(rdn_noescape)
			if len(rdn) > 20 {
				rdn = ""
			}
			updateSysIDwDriverID(rsq, rdn)
			http.Redirect(w, req, "/dispatchdriver", http.StatusSeeOther)
			return
		}
	}

	showSessions() // for demonstration purposes

	data := querySystemIDDrivers()

	parseData := struct {
		U    string
		Data []systemQueueParseData
	}{
		u, data,
	}

	if queryUserRole(u) == "Dispatch Supervisor" || queryUserRole(u) == "superuser#1" {

		tpl.ExecuteTemplate(w, "dispatchdriver.gohtml", parseData)
	} else {
		tpl.ExecuteTemplate(w, "allthefoodisgone.gohtml", nil)

	}

}

func allthefoodisgone(w http.ResponseWriter, req *http.Request) {

	showSessions() // for demonstration purposes

	tpl.ExecuteTemplate(w, "allthefoodisgone.gohtml", nil)

}

func userprofile(w http.ResponseWriter, req *http.Request) {

	u := getUser(w, req) //getUser function call

	if u == "" {
		//no username present!
		http.Redirect(w, req, "/login_redirect", http.StatusSeeOther)
		return
	}

	// var u user

	var me = make(map[string]string) //make map for error

	parseData := struct {
		U    string
		Data map[string]string
	}{
		u, me,
	}

	if req.Method == http.MethodPost {

		// get form values
		rb := req.FormValue("updatePassword")
		if rb == "updatePassword" {
			po := req.FormValue("passwordOld")
			pn := req.FormValue("passwordNew")
			dbPassword := queryPasswordUsersDB(u)
			// err := bcrypt.CompareHashAndPassword(u.Password, []byte(p))
			err := bcrypt.CompareHashAndPassword([]byte(dbPassword), []byte(po))
			if err != nil {
				me["Password"] = "Your Old Password is wrong. Please try again"

				tpl.ExecuteTemplate(w, "userprofile.gohtml", parseData) //replace nil as data
				return

			} else {
				//update saved password in db.
				bs, err := bcrypt.GenerateFromPassword([]byte(pn), bcrypt.MinCost)
				if err != nil {
					http.Error(w, "Internal server error", http.StatusInternalServerError)
					return
				}
				updateUserPassword(u, bs)
				me["Password"] = "Password updated successfully."

				tpl.ExecuteTemplate(w, "userprofile.gohtml", parseData) //replace nil as data
				return

			}
		} //ends updatedPassword

		rc := req.FormValue("updateAddress")

		if rc == "updateAddress" {

			a1 := req.FormValue("fullAddress")
			a2 := req.FormValue("postalCode")

			me["Address"] = "Address updated successfully."

			_, _ = a1, a2

			fmt.Printf("\nusername: %v, updated their address\n", u)

			tpl.ExecuteTemplate(w, "userprofile.gohtml", parseData) //replace nil as data
			return

		}

		rm := req.FormValue("updateContact")

		_ = rm

		me["ContactNumber"] = "Your contact Number has been updated successfully"

		fmt.Printf("\nusername: %v, updated their contact number\n", u)

		tpl.ExecuteTemplate(w, "userprofile.gohtml", parseData) //replace nil as data
		return

	}
	showSessions() // for demonstration purposes
	tpl.ExecuteTemplate(w, "userprofile.gohtml", parseData)
}

func sample(w http.ResponseWriter, req *http.Request) {

	// u := getUser(w, req) //getUser function call
	u := "username"
	if u == "" {
		//no username present!
		http.Redirect(w, req, "/login_redirect", http.StatusSeeOther)
		return
	}

	showSessions() // for demonstration purposes

	data := queryCheckoutConfirmationItems(u)

	parseData := struct {
		U    string
		Data []checkoutParseData
	}{
		u, data,
	}

	if req.Method == http.MethodPost {
		//no username present!
		r := req.FormValue("homebutton")
		if r == "home" {
			http.Redirect(w, req, "/", http.StatusSeeOther)
		}

	}

	tpl.ExecuteTemplate(w, "sample.gohtml", parseData)

}
