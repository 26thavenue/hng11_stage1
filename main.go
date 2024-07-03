package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"strings"

	owm "github.com/briandowns/openweathermap"
	"github.com/ipinfo/go/v2/ipinfo"
	"github.com/joho/godotenv"
)


type Resp struct {
   Client_ip    string `json:"client_ip"` 
   Location    string `json:"location"`
   Greeting    string  `json:"greeting"`
}

type Greeting struct {
	Message string `json:"message"`
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


    ip_address := getIPAddress(r)
    

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
    // Check for X-Forwarded-For header first
    ip := r.Header.Get("X-Forwarded-For")
    if ip != "" {
        
        ips := strings.Split(ip, ",")
        return strings.TrimSpace(ips[0])
    }

    
    ip = r.Header.Get("X-Real-IP")
    if ip != "" {
        return ip
    }

    ip, _, err := net.SplitHostPort(r.RemoteAddr)
    if err != nil {
        return ""
    }

    defer fmt.Println(ip)

    
    if net.ParseIP(ip).IsPrivate() {
        return getPublicIP()
    }

    return ip
}

func getPublicIP() string {
    resp, err := http.Get("https://api.ipify.org")
    if err != nil {
        return ""
    }
    defer resp.Body.Close()

    ip, err := io.ReadAll(resp.Body)
    if err != nil {
        return ""
    }

    fmt.Println(ip)
    return string(ip)
}

func Handler(w http.ResponseWriter, r *http.Request) {
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

func Greet(w http.ResponseWriter, r *http.Request){
    response := Greeting{
		Message: "Hello there",
	}

	
	w.Header().Set("Content-Type", "application/json")

	
	json.NewEncoder(w).Encode(response)
}

func main(){
    envError := godotenv.Load()
    if envError != nil {
        fmt.Printf("Error loading .env file")
    }

    log.Printf("API_KEY: %s", os.Getenv("API_KEY"))
    log.Printf("WEATHER_API_KEY: %s", os.Getenv("WEATHER_API_KEY"))

    http.HandleFunc("/api/hello", Handler)

    http.HandleFunc("/", Greet)

    const port = "8080"

    fmt.Println("Sever has started")

    log.Fatal(http.ListenAndServe("0.0.0.0:" + port, nil))

    
}
