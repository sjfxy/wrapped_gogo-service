package main

import (
	"fmt"
	"github.com/cloudfoundry-community/go-cfenv"
	"os"
	"sin.com/gogo-service/service"
)

func main() {
	port := os.Getenv("PORT")
	if len(port) == 0 {
		port = "3000"
	}
	appEnv, err := cfenv.Current()
	//fmt.Println(err)

	if err != nil {
		fmt.Println("CF Environment not detected.")
	}
	//fmt.Println(appEnv)
	//server :=
	server := service.NewServer(appEnv)
	server.Run(":" + port)
}
