
# Weather Info Server

This Go application provides a simple HTTP server that returns the IP address, location, and weather information based on the client's IP address. It uses the `ipinfo` and `openweathermap` APIs to retrieve the necessary information.

## Features

- Retrieves the client's IP address.
- Gets the geographical location based on the IP address.
- Fetches the current weather information for the detected location.
- Returns a greeting message including the client's name, temperature in Celsius, and location.

## Prerequisites

- Go 1.16 or later
- `ipinfo` API key
- `openweathermap` API key

## Installation

1. Clone the repository:

   ```sh
   git clone https://github.com/yourusername/weather-info-server.git
   cd weather-info-server
   ```

2. Install dependencies:

   ```sh
   go mod tidy
   ```

3. Create a `.env` file in the root directory of your project and add your API keys:

   ```sh
   API_KEY=your_ipinfo_api_key
   WEATHER_API_KEY=your_openweathermap_api_key
   ```

## Usage

1. Run the server:

   ```sh
   go run main.go
   ```

2. The server will start on port 8080. You can make a GET request to the server with a query parameter `visitor_name`:

   ```sh
   curl "http://localhost:8080/?visitor_name=Mark"
   ```

3. The server will respond with a JSON object containing the client's IP address, location, and a greeting message:

   ```json
   {
       "client_ip": "8.8.8.8",
       "location": "New York",
       "greeting": "Hello, Mark! The temperature is 11 degrees Celsius in New York."
   }
   ```

## Project Structure

- `main.go`: The main server file.
- `.env`: Environment file containing the API keys.

## Environment Variables

- `API_KEY`: Your `ipinfo` API key.
- `WEATHER_API_KEY`: Your `openweathermap` API key.

## Dependencies

- `github.com/briandowns/openweathermap`: OpenWeatherMap API client for Go.
- `github.com/ipinfo/go/v2/ipinfo`: IPinfo API client for Go.
- `github.com/joho/godotenv`: Library to load environment variables from a `.env` file.

## Code Overview

### Main Functions

- `getLocationAndIP(r *http.Request) string`: Retrieves the client's location based on the IP address.
- `getWeatherInfo(loc string) owm.Main`: Fetches weather information for the given location.
- `getIPAddress(r *http.Request) string`: Extracts the client's IP address from the request.
- `handler(w http.ResponseWriter, r *http.Request)`: Main handler function that processes the request and returns the response.
- `main()`: Initializes the server and loads environment variables.

