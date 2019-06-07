package utils

import (
	"io/ioutil"
	"log"
)

func ReadBytesFromFile(filaName string) []byte {
	log.Println("Reading file: " + filaName)
	data, err := ioutil.ReadFile(filaName)
	if err != nil {
		log.Fatal(err)
	}

	return data
}
