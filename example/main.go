package main

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/ganmacs/swim"
	"github.com/ganmacs/swim/logger"
)

func main() {
	sport := os.Getenv("PORT")
	port, _ := strconv.Atoi(sport)

	config := swim.DefaultConfig()
	config.BindAddr = "127.0.0.1"
	config.BindPort = port

	config.Logger = logger.NewLevelLogger(os.Stdout, logger.INFO)

	swim, err := swim.New(config)
	if err != nil {
		panic(err)
	}

	clusterSize, err := swim.Join("127.0.0.1:3000")
	if err != nil {
		panic(err)
	}

	fmt.Printf("cluster size is %d\n", clusterSize)

	time.Sleep(time.Second * 60)
}
