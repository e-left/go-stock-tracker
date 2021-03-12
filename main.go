package main

import(
	"fmt"
	"log"
	"github.com/joho/godotenv"
	// "net/http"
)

func main() {
	fmt.Println("Golang crypto tracker")
	var envs map[string]string
	envs, err := godotenv.Read(".env")

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	api_key := envs["API_KEY"]
	fmt.Printf("Api key: %s\n", api_key)



}
