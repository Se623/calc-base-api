package main

import (
	"fmt"
	"os"

	"github.com/Se623/calc-base-api/internal/calc"
	"go.uber.org/zap"
)

func main() {
	logger, _ := zap.NewProduction()
	defer logger.Sync()
	sugar := logger.Sugar()
	sugar.Info("Main: Main program is started")
	entry := make(chan string)
	go calc.Orchestrator(entry, sugar)
	expr := ""
	for {
		fmt.Fscan(os.Stdin, &expr)
		if expr == "!show" {
			fmt.Println("nuh-uh")
		} else {
			entry <- expr
		}
	}
}
