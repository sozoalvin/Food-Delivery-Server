package main

import (
	"crypto/tls"
	"fmt"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	uuid "github.com/satori/go.uuid"
	"golang.org/x/crypto/acme/autocert"
)

var FoodMerchantNameAddress []string
var FoodMerchantNameAddress1 []string

func generateTransactionID() (string, error) {

	var generatedID string

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

func getSelfSignedOrLetsEncryptCert(certManager *autocert.Manager) func(hello *tls.ClientHelloInfo) (*tls.Certificate, error) {
	return func(hello *tls.ClientHelloInfo) (*tls.Certificate, error) {
		dirCache, ok := certManager.Cache.(autocert.DirCache)
		if !ok {
			dirCache = "certs"
		}

		keyFile := filepath.Join(string(dirCache), hello.ServerName+".key")
		crtFile := filepath.Join(string(dirCache), hello.ServerName+".crt")
		certificate, err := tls.LoadX509KeyPair(crtFile, keyFile)
		if err != nil {
			fmt.Printf("%s\nFalling back to Letsencrypt\n", err)
			return certManager.GetCertificate(hello)
		}
		fmt.Println("Loaded selfsigned certificate.")
		return &certificate, err
	}
}

func generateApiKey(s string) string {

	//time format in DDMMYY
	var timeNow string = time.Now().Format("010206")

	result := strings.Split(s, "")

	result1 := result[0:26]

	result1WithDate := append(result1, timeNow)

	userApiKey := strings.Join(result1WithDate, "")

	return userApiKey
}

func reGenerateApiKey() string {

	var timeNow string = time.Now().Format("010206")

	sID, err := uuid.NewV4()
	//err handling
	if err != nil {
		fmt.Printf("Something went wrong: %s", err)
	}

	result := strings.Split(sID.String(), "")

	result1 := result[0:26]

	result1WithDate := append(result1, timeNow)

	userApiKey := strings.Join(result1WithDate, "")

	return userApiKey

}

func retriFoodMerchantName(s string) (string, string) {

	slice := strings.Split(s, "")
	var i int = 0
	for slice[i] != "-" {
		i++
	}

	result := slice[0 : i-1]
	foodName := strings.Join(result, "")

	var j int = i
	for i < len(slice) {
		i++

	}
	result3 := slice[j+2 : i]
	merchantName := strings.Join(result3, "")

	return foodName, merchantName
}

func toTitle(s string) string {

	slice := strings.Split(s, "")
	var pos []int
	for i := 0; i < len(slice); i++ {

		if slice[i] == " " {
			pos = append(pos, i+1)
			i++
		}

	}
	pos = append(pos, 0)
	for i := 0; i < len(pos); i++ {

		slice[pos[i]] = strings.ToUpper(slice[pos[i]])

	}
	result := strings.Join(slice, "")

	return result
}

func processPIDvalue(s string) string {

	r1 := strings.Split(s, "")
	r2 := r1[5:]
	r3 := strings.Join(r2, "")
	r4, err := strconv.Atoi(r3)
	if err != nil {
		fmt.Println("There was an error converting PID string into int")
	}

	r5 := "KVPID" + strconv.Itoa(r4+1)

	return r5
}
