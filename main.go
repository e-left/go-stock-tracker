package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/joho/godotenv"
)

const (
	BASE_URL = "https://pro-api.coinmarketcap.com/v1/cryptocurrency/quotes/latest"
)

type Res struct {
	Time         string  `json:"time"`
	BitcoinPrice float64 `json:"btcPrice"`
	EtherumPrice float64 `json:"ethPrice"`
}

func getCryptoPrices(apiKey string) Res {
	client := &http.Client{}
	request, err := http.NewRequest("GET", BASE_URL, nil)
	if err != nil {
		log.Fatal(err)
	}

	q := url.Values{}
	// q.Add("start", "1")
	// q.Add("limit", "5000")
	q.Add("convert", "USD")
	q.Add("symbol", "BTC,ETH")

	request.Header.Set("Accepts", "application/json")
	request.Header.Set("X-CMC_PRO_API_KEY", apiKey)
	request.URL.RawQuery = q.Encode()

	response, err := client.Do(request)

	if err != nil {
		log.Fatal("Error sending request to server")
	}

	// fmt.Println(response.Status)
	responseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal("Error reading response body")
	}

	resJSON := make(map[string]interface{})
	err = json.Unmarshal(responseBody, &resJSON)
	if err != nil {
		log.Fatal(err)
	}

	// fmt.Println(resJSON)
	btcPrice := resJSON["data"].(map[string]interface{})["BTC"].(map[string]interface{})["quote"].(map[string]interface{})["USD"].(map[string]interface{})["price"]
	ethPrice := resJSON["data"].(map[string]interface{})["ETH"].(map[string]interface{})["quote"].(map[string]interface{})["USD"].(map[string]interface{})["price"]

	now := time.Now()
	// fmt.Printf("%v:%v:%v, %v-%v-%v\n", now.Hour(), now.Minute(), now.Second(), now.Day(), now.Month(), now.Year())
	// fmt.Printf("Bitcoin price: %v\n", btcPrice)
	// fmt.Printf("Etherum price: %v\n", ethPrice)
	// fmt.Println(string(responseBody))

	sendBack := Res{now.String(), btcPrice.(float64), ethPrice.(float64)}

	return sendBack

}

func main() {
	fmt.Println("Golang crypto tracker")
	var envs map[string]string
	envs, err := godotenv.Read(".env")

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	api_key := envs["API_KEY"]
	// fmt.Printf("Api key: %s\n", api_key)

	http.Handle("/", http.FileServer(http.Dir("./static/")))
	http.HandleFunc("/prices", func(w http.ResponseWriter, r *http.Request) {
		log.Print("New request")
		sendBack := getCryptoPrices(api_key)
		json.NewEncoder(w).Encode(sendBack)
	})

	fmt.Println("Web server starting")
	log.Fatal(http.ListenAndServe(":8000", nil))

}
