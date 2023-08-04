package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {

	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", "http://server-app:8080/cotacao", nil)
	if err != nil {
		panic(err)
	}
	response, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		panic(err)
	}

	log.Println("JSON formatado:")
	log.Println(string(body))

	var data map[string]string
	err = json.Unmarshal(body, &data)
	if err != nil {
		log.Println("Error decoding JSON:", err)
		panic(err)
		return
	}

	bidValue, ok := data["bid"]
	if !ok {
		log.Println("JSON does not contain the 'bid' field")
		panic(err)
		return
	}

	content := fmt.Sprintf("Dólar: %s", bidValue)
	file, err := os.Create("cotacao.txt")
	if err != nil {
		panic(err)
	}

	_, err = file.Write([]byte(content))
	if err != nil {
		panic(err)
	}

	log.Println("Arquivo criado com sucesso!")
	file.Close()
	//io.Copy(os.Stdout, response.Body)
}
