package main

import (
	"encoding/json"
	"fmt"
	_ "io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	_ "strings"

	owm "github.com/briandowns/openweathermap"
	"github.com/ipinfo/go/v2/ipinfo"
	"github.com/joho/godotenv"
)


type Resp struct {
   Client_ip    string `json:"client_ip"` 
   Location    string `json:"location"`
   Greeting    string  `json:"greeting"`
}

type WeatherCondition struct {
    Code        int
    Main        string
    Description string
    Icon        string
}



func getLocationAndIP(r *http.Request) string {
    token := os.Getenv("API_KEY")

    if token == "" {
         log.Fatal("API key is empty")
    }

    

    client := ipinfo.NewClient(nil, nil, token)

    

    // const ip_address = "8.8.8.8"

    ip_address := getIPAddress(r)

    // fmt.Println(ip_address)

	info, err := client.GetIPInfo(net.ParseIP(ip_address))

	if err != nil {
		log.Fatal(err)
	}

    

    city := info.City

    fmt.Println(city)

    return city
    
}

func getWeatherInfo(loc string) owm.Main {
    apiKey := os.Getenv("WEATHER_API_KEY")

    w, err := owm.NewCurrent("C", "EN", apiKey) 
    if err != nil {
        log.Fatalln(err)
    }

    if loc == "" {
        fmt.Println("Invalid location parameter")
    }

    w.CurrentByName(loc)

    fmt.Println(w.Main.Temp)

    return w.Main
}


func getIPAddress(r *http.Request) string {
    ip, _, err := net.SplitHostPort(r.RemoteAddr)
	
    if err != nil {
        log.Println("Error getting IP address:", err)
        return ""
    }

    userIP := net.ParseIP(ip)
    if userIP == nil {
        log.Println("Error parsing IP address:", err)
        return ""
    }

    // fmt.Println(userIP)

    return userIP.String()

}

func handler(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("visitor_name")
	
    clientIP := getIPAddress(r)

    city := getLocationAndIP(r)

    weatherArr := getWeatherInfo(city)

    temp := int(weatherArr.Temp)

     greeting := fmt.Sprintf("Hello, %s! The temperature is %d degrees Celsius in %s.", name, temp, city)

    resp := Resp{
        Client_ip: clientIP,
        Location:  city,
        Greeting:  greeting,
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(resp)


}

func main() {
    envError := godotenv.Load()
    if envError != nil {
        log.Fatal("Error loading .env file")
    }

    http.HandleFunc("/", handler)

    fmt.Println("Starting server on :8080")

    err := http.ListenAndServe(":8080", nil)

    if err != nil {
        fmt.Println("Error starting server:", err)
    }
}