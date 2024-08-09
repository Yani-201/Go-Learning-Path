package main

import (
	"task-manager-api-Auth/router"
)

func main() {
	router := router.SetupRouter()
	if err := router.Run("localhost:8080"); err != nil {
		panic(err)
	}
}
