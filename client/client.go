package main

import (
	"fmt"
	"log"

	"github.com/vitorbass93/challenge1/client/gateway"
	"github.com/vitorbass93/challenge1/client/provider"
)

func main() {
	dolarProvider := provider.NewDolarProvider("http://localhost:8080/cotacao")
	fileGateway := gateway.NewFileGateway("cotacao.txt")
	bid, err := dolarProvider.GetDolar()
	if err != nil {
		log.Fatalf("dolarProvider.GetDolar() failed with err: %+v\n", err)
	}
	fileGateway.WriteFile(fmt.Sprintf("Dolar: %s", bid))
}
