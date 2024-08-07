package main

import (
	"fmt"
	"task-manager-api/router"
)

func main() {
	fmt.Println("Task Manger API")
	route := router.SetupRouter()
	route.Run("localhost:8080")
}