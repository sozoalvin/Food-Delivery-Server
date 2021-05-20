package main

import (
	"database/sql"
	"flag"
	"fmt"
	"net/http"
	"sync"
	"text/template"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"golang.org/x/crypto/acme/autocert"
)

var countpid int = 1

var MyFoodListDB = InitMyFoodList() //creates the trie required to power the
var TransID int = 0
var QueueID int = 500
var pid int = 0

var searchResult []string
var selectedProduct string

var searchResult2 []searchResultFormat
var mutex sync.Mutex

var db *sql.DB

var (
	domain        string
	productionFlg bool = false
)

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

type MerchantFoodInfo struct {
	FoodName string
	Price    string
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

	var dbErr error

	flag.StringVar(&domain, "domain", "", "domain name to request your certificate")
	flag.BoolVar(&productionFlg, "productionFlg", false, "if true, we start HTTPS server")
	flag.Parse()

	currentTime := time.Now()

	db, dbErr = sql.Open("mysql", "newuser1:password@tcp(127.0.0.1:54779)/my_db?charset=utf8")

	check(dbErr)
	defer db.Close()
	dbErr = db.Ping()
	check(dbErr)

	fmt.Println("System Message : System is Ready", currentTime.Format("2006-01-02 15:04:05"))

	router := mux.NewRouter()
	router.HandleFunc("/", index)
	router.HandleFunc("/signup", signup)
	router.HandleFunc("/login", login)
	router.HandleFunc("/logout", logout)
	router.HandleFunc("/userprofile", userprofile)
	router.HandleFunc("/assets/rs2.png", rs2)
	router.HandleFunc("/assets/rs3.png", rs3)
	router.HandleFunc("/assets/rs4.png", rs4)
	router.HandleFunc("/assets/rs5.png", rs5)
	router.HandleFunc("/assets/rs6.png", rs6)
	router.HandleFunc("/assets/rs7.png", rs7)
	router.HandleFunc("/assets/rs8.png", rs8)
	router.HandleFunc("/assets/login.png", loginbutton)
	router.HandleFunc("/assets/signup.png", signupbutton)
	router.HandleFunc("/searchresult", searchresult)
	router.HandleFunc("/yourcart", yourcart)
	router.HandleFunc("/checkout_processing", checkout_processing)
	router.HandleFunc("/checkout", checkout)
	router.HandleFunc("/allsystemorders", allsystemorders)
	router.HandleFunc("/alltransactions", alltransactions)
	router.HandleFunc("/login_redirect", login_redirect)
	router.HandleFunc("/dispatchdriver", dispatchdriver)
	router.HandleFunc("/clearcart", clearcart)
	router.HandleFunc("/allthefoodisgone", allthefoodisgone)
	router.HandleFunc("/api/v1/apivalidation", apiValidation).Methods("Get")
	router.HandleFunc("/api/v1/nameaddress", nameAddress).Methods("Get")
	router.HandleFunc("/api/v1/retrieveall", retrieveAll).Methods("Get")
	router.HandleFunc("/api/v1/additems", additems).Methods("Post")
	router.HandleFunc("/api/v1/merchantsetup", merchantsetup).Methods("Post")
	router.HandleFunc("/api/v1/edititems", editItems).Methods("Put")
	router.HandleFunc("/api/v1/deleteitems", deleteItems).Methods("Delete")
	router.Handle("/favicon.ico", http.NotFoundHandler()) //NotFoundHandler returns a simple request handler that replies to each request with a “404 page not found” reply.
	router.HandleFunc("/sample", sample)

	if productionFlg {
		certManager := autocert.Manager{
			Prompt:     autocert.AcceptTOS,
			HostPolicy: autocert.HostWhitelist(domain),
			Cache:      autocert.DirCache("certs"),
		}
		tlsConfig := certManager.TLSConfig()
		tlsConfig.GetCertificate = getSelfSignedOrLetsEncryptCert(&certManager)
		server := http.Server{
			Addr:      ":443",
			Handler:   router,
			TLSConfig: tlsConfig,
		}

		go http.ListenAndServe(":80", certManager.HTTPHandler(nil))
		fmt.Println("Server listening on", server.Addr)
		if err := server.ListenAndServeTLS("", ""); err != nil {
			fmt.Println(err)
		}
	} else {
		fmt.Println("ProductionFlg activated\nServer listening on :80")
		http.ListenAndServe(":80", router)
	}

} // close main functioin

func check(err error) {
	if err != nil {
		fmt.Println(err)
	}
} //end function check
