package main

import (
	"awesome/api"
	"awesome/config"
)

func init() {
	config.LoadConfiguration()
}

func main() {
	var service = api.NewAPIService()
	service.AddEndpoints()
	service.Run()
}
