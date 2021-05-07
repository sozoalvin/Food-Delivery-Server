package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"sync"
	"text/template"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var MyFoodListDB = InitMyFoodList()
var TransID int = 0
var QueueID int = 500
var pid int = 0

var searchResult []string
var selectedProduct string

var searchResult2 []searchResultFormat
var mutex sync.Mutex

var db *sql.DB

type user struct {
	Username        string
	Fname           string
	Lname           string
	DeliveryAddress string
	PostalCode      string
	MobileNumber    string
}

type searchResultFormat struct {
	FoodName string
	PID      string
}

type cartDisplay struct {
	FoodName   string
	Quantity   string
	UnitPrice  float64
	TotalPrice float64
}

var cartDisplayList []cartDisplay

type FoodInfo struct {
	FoodName         string
	MerchantName     string
	DetailedLocation string
	PostalCode       int
	Price            float64
	OpeningPeriods   OpeningPeriods
}

type OpeningPeriods map[string][]string

type cartData struct {
	PID      string
	Quantity string
}

type cartDisplayData struct {
	FoodName  string
	Quantity  string
	UnitPrice float64
	TotalCost float64
}

type cartFullData struct {
	PID       string
	FoodName  string
	Quantity  int
	UnitPrice float64
	TotalCost float64
	UserRole  string
}

type systemQueueParseData struct {
	SysID    string
	Pi       int
	Time     string
	DriverID string
}

type TransactionsParseData struct {
	Username  string
	Time      string
	TransID   string
	FoodName  string
	Quantity  int
	TotalCost float64
}

type checkoutParseData struct {
	Time      string
	TransID   string
	SysID     string
	FoodName  string
	Quantity  int
	TotalCost float64
}

type session struct {
	un           string
	lastActivity time.Time
}

var tpl *template.Template

var dbSessions = map[string]session{}
var dbSessionsCleaned time.Time

func init() {
	tpl = template.Must(template.ParseGlob("templates/*"))
	dbSessionsCleaned = time.Now()
}

func main() {

	ch := make(chan string) //create a channel called c
	var dbErr error

	currentTime := time.Now()

	go CreateFoodList(ch) //newResult is a slice that is being returned byCreateFoodList function
	fmt.Println("\nSystem Message :", <-ch)
	go CreateFoodListMap(ch)
	go MyFoodListDB.PreInsertTrie(FoodMerchantNameAddress, ch) //populates Trie Data for Food LIst
	fmt.Println("System Message :", <-ch)
	fmt.Println("System Message :", <-ch)
	myPostalCodesDB := InitPostalCode()   //creates PostalCode BST DB
	myPostalCodesDB.PreInsertPostalCode() //preinset POSTAL Code DB
	FoodMerchantNameAddressProductID()
	fmt.Println("System Message : System is Ready", currentTime.Format("2006-01-02 15:04:05"))

	db, dbErr = sql.Open("mysql", "insert_SQL Data here")

	check(dbErr)
	defer db.Close()
	dbErr = db.Ping()
	check(dbErr)

	http.HandleFunc("/", index)
	http.HandleFunc("/signup", signup)
	http.HandleFunc("/login", login)
	http.HandleFunc("/logout", logout)
	http.HandleFunc("/userprofile", userprofile)
	http.HandleFunc("/assets/rs2.png", rs2)
	http.HandleFunc("/searchresult", searchresult)
	http.HandleFunc("/yourcart", yourcart)
	http.HandleFunc("/checkout_processing", checkout_processing)
	http.HandleFunc("/checkout", checkout)
	http.HandleFunc("/allsystemorders", allsystemorders)
	http.HandleFunc("/alltransactions", alltransactions)
	http.HandleFunc("/login_redirect", login_redirect)
	http.HandleFunc("/dispatchdriver", dispatchdriver)
	http.HandleFunc("/clearcart", clearcart)
	http.HandleFunc("/allthefoodisgone", allthefoodisgone)
	http.Handle("/favicon.ico", http.NotFoundHandler()) //NotFoundHandler returns a simple request handler that replies to each request with a “404 page not found” reply.
	http.ListenAndServe(":80", nil)                     //launches HTTP server

} // close main functioin

func check(err error) {
	if err != nil {
		fmt.Println(err)
	}
} //end function check
