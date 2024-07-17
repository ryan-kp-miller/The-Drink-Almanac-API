package main

import "os"

func main() {
	port, ok := os.LookupEnv("API_PORT")
	if !ok {
		port = "8000"
	}
	Start(port)
}
