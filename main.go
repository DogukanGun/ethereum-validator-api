package main

import (
	"ethereum-validator-api/utils"
	"github.com/gin-gonic/gin"
	"log"
)

func main() {
	utils.InitializeENV(".env")
	router := gin.Default()
	err := utils.SetupEndpoints(router)
	if err != nil {
		log.Fatalf("Error is :%v", err)
	}
}
