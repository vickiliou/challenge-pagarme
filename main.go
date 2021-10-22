package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/pagarme/marshals/labs/vicki/desafioGo/config"
	"github.com/pagarme/marshals/labs/vicki/desafioGo/logger"
	"github.com/pagarme/marshals/labs/vicki/desafioGo/transaction"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type txn struct {
	apiKey, endpoint string
	authType         []string
	paymentMethods   []paymentMethod
}

type paymentMethod struct {
	method []byte
}

func readFile(fileName string) []byte {
	jsonFile, err := ioutil.ReadFile(fileName)
	if err != nil {
		log.Printf("Error reading the file %s: %s", fileName, err.Error())
	}
	return jsonFile
}

func main() {
	endpoint, err := config.GetEnvVar(config.APIEndpointEnv)
	if err != nil {
		log.Fatal("Endpoint not found")
	}

	apiKey, err := config.GetEnvVar(config.APIKeyEnv)
	if err != nil {
		log.Fatal("API key not found")
	}

	creditCardTransaction := readFile("transaction/creditCardTransaction.json")
	boletoTransaction := readFile("transaction/boletoTransaction.json")

	t := txn{
		apiKey:   apiKey,
		endpoint: endpoint,
		authType: []string{
			"inBody",
			"pathParam",
			"basicAuth",
		},
		paymentMethods: []paymentMethod{
			{method: creditCardTransaction},
			{method: boletoTransaction},
		},
	}

	http.Handle("/metrics", promhttp.Handler())
	for _, dataAuth := range t.authType {
		for _, dataPaymentMethod := range t.paymentMethods {
			body, err := transaction.SendTx(t.endpoint, t.apiKey, dataAuth, dataPaymentMethod.method)
			if err != nil {
				fmt.Printf("Error with transaction %s - %s: %s", dataAuth, dataPaymentMethod.method, err.Error())
				continue
			}
			logger.InfoLog.Println(string(body))
		}
	}

	config.VersionHandler()
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err.Error())
	}
}
