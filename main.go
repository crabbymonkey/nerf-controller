package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"math"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

var isRunning = false
var appCode = ""
var hopper = 0

func logMessage(message string) {
	currentTime := time.Now()
	fmt.Println("{" + currentTime.Format(time.RFC1123) + "} " + message)
}

func getLifeLeftOfAccessToken(accessToken *AccessToken) time.Duration {
	return (time.Duration(accessToken.LifeTime) * time.Second) - time.Since(accessToken.TimeBorn)
}

func handleFireing() {
	for {
		if hopper > 0 {
			fire()
			hopper--
		}
		time.Sleep(time.Duration(randomValue(10, 90)) * time.Second)
	}
}

func listenAndHandleDonations(accessToken *AccessToken) {
	lastDonationID := getLastDonationID(accessToken)
	remainder := 0.00
	for {
		if !isRunning {
			break
		}

		refreshAccessToken(accessToken)
		// timeLeftOnAccessToken := getLifeLeftOfAccessToken(accessToken)
		// logMessage("Time remaining on access_token: " + (timeLeftOnAccessToken.Round(time.Second)).String())
		// if timeLeftOnAccessToken < time.Minute {
		// 	refreshAccessToken(accessToken)
		// 	if accessToken == nil {
		// 		log.Fatalln("error refreshing access token")
		// 		break
		// 	}
		// }
		logMessage("checking for donations")
		donations := getStreamlabsDonations(accessToken, nil, lastDonationID)
		var donation Donation
		for _, donation = range *donations {
			fmt.Printf("%+v\n", donation)
			intpart, div := math.Modf(donation.Amount)
			hopper += int(intpart)
			remainder += div
			if remainder > 1 {
				remainder--
				hopper++
			}
			logMessage("Hopper: " + strconv.Itoa(hopper) + ", Remainder: " + fmt.Sprintf("%f", remainder))
		}
		if len(*donations) != 0 {
			*lastDonationID = donation.DonationID //TODO: This shouldent be in the loop
		}
		time.Sleep(5 * time.Second)
	}
}

func getStreamlabsDonations(accessToken *AccessToken, numDonations *int, afterDonationID *int) *[]Donation {
	urlParams := url.Values{}
	urlParams.Add("access_token", accessToken.Val)
	urlParams.Add("currency", "USD")

	if numDonations != nil {
		urlParams.Add("limit", strconv.Itoa(*numDonations))
	}
	if afterDonationID != nil {
		urlParams.Add("after", strconv.Itoa(*afterDonationID))
	}
	finalURL := "https://streamlabs.com/api/v1.0/donations?" + urlParams.Encode()
	logMessage("request URL: " + finalURL)
	resp, err := http.Get(finalURL)
	if err != nil {
		log.Fatalln(err)
		return nil
	}
	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		log.Fatalln("bad response getting access_token: " + string(body))
		return nil
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
		return nil
	}

	var rawDonations DonationData
	err = json.Unmarshal(data, &rawDonations)
	if err != nil {
		log.Fatalln(err)
		return nil
	}

	return &rawDonations.Data
}

func getLastDonationID(accessToken *AccessToken) *int {
	if !isRunning {
		return nil
	}
	numDonations := 1
	donations := getStreamlabsDonations(accessToken, &numDonations, nil)
	if len((*donations)) != 1 {
		return nil
	}
	fmt.Printf("Last Donation:\n%+v\n", &(*donations)[0])
	return &(*donations)[0].DonationID
}

//Compile templates on start
var templates = template.Must(template.ParseFiles(
	"templates/notFound.html",
	"templates/header.html",
	"templates/footer.html",
	"templates/index.html"))

//GrantType is the type of the request, refresh or get a new token
type GrantType string

const (
	//RefreshToken used to refresh the token
	RefreshToken GrantType = "refresh_token"
	//AuthorizationCode used to get a new code
	AuthorizationCode GrantType = "authorization_code"
)

func makeAccesTokenRequest(accessToken *AccessToken, grantType GrantType) {
	app := getStreamlabsAppInfo()

	logMessage("making request for access_code using grant_type " + string(grantType))

	var grantTypeLabel string
	var grantTypeVal string
	if grantType == RefreshToken {
		grantTypeLabel = "refresh_token"
		grantTypeVal = accessToken.RefreshToken

	} else if grantType == AuthorizationCode {
		grantTypeLabel = "code"
		grantTypeVal = appCode
	}

	message := map[string]interface{}{
		"grant_type":    grantType,
		"client_id":     app.ClientID,
		"client_secret": app.ClientSecret,
		"redirect_uri":  app.RedirectURI,
		grantTypeLabel:  grantTypeVal,
	}

	bytesRepresentation, err := json.Marshal(message)
	if err != nil {
		log.Fatalln("error marshaling request body: " + err.Error())
	}

	resp, err := http.Post("https://streamlabs.com/api/v1.0/token",
		"application/json", bytes.NewBuffer(bytesRepresentation))
	if err != nil {
		log.Fatalln(err)
		accessToken = nil
		return
	}
	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		log.Fatalln("bad response getting access_token: " + string(body))
		accessToken = nil
		return
	}

	accessToken.TimeBorn = time.Now()

	err = json.NewDecoder(resp.Body).Decode(&accessToken)
	if err != nil {
		log.Fatalln("error decoding response for access token: " + err.Error())
		accessToken = nil
		return
	}
}

func refreshAccessToken(accessToken *AccessToken) {
	makeAccesTokenRequest(accessToken, RefreshToken)
}

func getAccessToken() *AccessToken {
	var accessToken AccessToken
	makeAccesTokenRequest(&accessToken, AuthorizationCode)

	return &accessToken
}

func fire() {
	logMessage("FIRE!!!")
}

func saveToken(token string) {
	logMessage("saving token: " + token)
}

//Display the named template
func display(w http.ResponseWriter, tmpl string, data interface{}) {
	templates.ExecuteTemplate(w, tmpl, data)
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	if isRunning {
		requestURL, err := url.Parse("/live")
		if err != nil {
			log.Fatal(err)
		}
		requestQuery := requestURL.Query()
		requestQuery.Set("code", appCode)
		requestURL.RawQuery = requestQuery.Encode()

		logMessage(requestURL.String())
		http.Redirect(w, r, requestURL.String(), http.StatusSeeOther)
	}

	data := Page{
		PageTitle:  "Home",
		Running:    isRunning,
		HopperSize: &hopper,
	}
	display(w, "index", data)
}

func liveHandler(w http.ResponseWriter, r *http.Request) {
	appCode = r.URL.Query().Get("code")
	if !isRunning && appCode != "" {
		accessToken := getAccessToken()
		if accessToken == nil {
			log.Fatalln("error getting the access token")
		}
		isRunning = true

		data := Page{
			PageTitle: "Live",
			Running:   isRunning,
		}
		display(w, "index", data)

		go listenAndHandleDonations(accessToken)
		go handleFireing()
	}
}

func randomPageHandler(w http.ResponseWriter, r *http.Request) {
	// If empty show the home page
	// If static page show the static package
	// Else show the 404 page
	if r.URL.Path == "/" {
		homeHandler(w, r)
	} else if r.URL.Path == "/fire" {
		fireHandler(w, r)
	} else if r.URL.Path == "/activate" {
		activateHandler(w, r)
	} else if r.URL.Path == "/live" {
		liveHandler(w, r)
	} else if r.URL.Path == "/stop" {
		stopHandler(w, r)
	} else if strings.HasSuffix(r.URL.Path[1:], ".html") {
		http.ServeFile(w, r, "static/html/"+r.URL.Path[1:])
	} else {
		fmt.Println("Sorry but it seems this page does not exist...")
		errorHandler(w, r, http.StatusNotFound)
	}
}

func fireHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		fire()
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func stopHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		isRunning = false
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func getStreamlabsAppInfo() StreamLabsApp {
	var app StreamLabsApp
	// pwd, _ := os.Getwd()
	file, err := os.Open("StreamLabsAPI.json")
	if err != nil {
		logMessage("ERROR OPENING FILE: " + err.Error())
	}
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&app)
	if err != nil {
		logMessage("ERROR DECODING JSON: " + err.Error())
	}
	return app
}

func activateHandler(w http.ResponseWriter, r *http.Request) {
	app := getStreamlabsAppInfo()

	requestURL, err := url.Parse("https://streamlabs.com/api/v1.0/authorize")
	if err != nil {
		log.Fatal("error makinking redirection url: " + err.Error())
	}
	requestQuery := requestURL.Query()
	requestQuery.Set("response_type", "code")
	requestQuery.Set("client_id", app.ClientID)
	requestQuery.Set("redirect_uri", app.RedirectURI)
	requestQuery.Set("scope", "donations.read")
	requestURL.RawQuery = requestQuery.Encode()

	logMessage(requestURL.String())
	http.Redirect(w, r, requestURL.String(), http.StatusSeeOther)
}

func errorHandler(w http.ResponseWriter, r *http.Request, status int) {
	w.WriteHeader(status)
	if status == http.StatusNotFound {
		display(w, "404", &Page{PageTitle: "404", Running: isRunning})
	}
}

// Gets a random value from the low to high values. This will include the low and high values.
func randomValue(low int, high int) int {
	var scaledInt = high - low + 1 // The +1 is to offset the values so it can be the high value.
	return rand.Intn(scaledInt) + low
}

func getPort() string {
	if value, ok := os.LookupEnv("PORT"); ok {
		return ":" + value
	}
	return ":8080"
}

func main() {
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.HandleFunc("/", randomPageHandler)
	var port = getPort()
	fmt.Println("Now listening to port " + port)
	log.Fatal(http.ListenAndServe(port, nil))
}

//A Page structure
type Page struct {
	PageTitle  string
	Running    bool
	HopperSize *int
}

//StreamLabsApp Object that holds the StreamLabs App info that comes from a JSON
type StreamLabsApp struct {
	ClientID     string
	ClientSecret string
	RedirectURI  string
}

//AccessToken stores the value and timers
type AccessToken struct {
	Val          string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	LifeTime     int    `json:"expires_in"`
	TokenType    string `json:"token_type"`
	TimeBorn     time.Time
}

//DonationData object of the raw data returned from StreamLabs when getting donations
type DonationData struct {
	Data []Donation `json:"data"`
}

//Donation object returned from StreamLabs
type Donation struct {
	DonationID int     `json:"donation_id"`
	CreatedAt  int     `json:"created_at"`
	Currency   string  `json:"currency"`
	Amount     float64 `json:"amount,string"`
	Name       string  `json:"name"`
	Message    string  `json:"message"`
	Email      string  `json:"email"`
}