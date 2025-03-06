package main

import (
	"fmt"
	"net/http"

	"github.com/Se623/calc-base-api/internal/agent"
	"github.com/Se623/calc-base-api/internal/lib"
	"github.com/Se623/calc-base-api/internal/orchestrator"
)

func main() {
	lib.InitLogger()

	lib.Sugar.Info("Initilized main program")

	for i := 0; i < lib.COMPUTING_AGENTS; i++ {
		go agent.Agent(i + 1)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/api/v1/calculate", orchestrator.Spliter)
	mux.HandleFunc("/internal/task", orchestrator.Distributor)
	mux.HandleFunc("localhost/api/v1/expressions", orchestrator.Displayer)
	lib.Sugar.Infof("Initilized orchestrator")

	lib.Sugar.Infof("Initilized server on port 8080")

	if err := http.ListenAndServe(":8080", mux); err != nil { // Запуск сервера
		fmt.Println(err)
		return
	}
}
