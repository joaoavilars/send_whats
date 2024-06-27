package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
)

type Message struct {
	Number  string `json:"number"`
	Message string `json:"message"`
}

func main() {
	url := "http://192.168.10.39:4000/api/message"
	phone := os.Args[1]
	title := os.Args[2]
	message := strings.Join(os.Args[3:], " ")

	// Substitui \n por quebras de linha reais
	message = strings.ReplaceAll(message, "\\n", "\n")

	msg := Message{
		Number:  phone,
		Message: fmt.Sprintf("%s\n\n%s", title, message),
	}

	// Converte a mensagem para JSON
	jsonBody, err := json.Marshal(msg)
	if err != nil {
		fmt.Println("Erro ao codificar mensagem JSON:", err)
		return
	}

	// Cria uma requisição HTTP POST
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		fmt.Println("Erro ao criar requisição HTTP:", err)
		return
	}
	req.Header.Set("Content-Type", "application/json")

	// Cliente HTTP para enviar a requisição
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Erro ao enviar requisição HTTP:", err)
		return
	}
	defer resp.Body.Close()

	// Verifica o status da resposta
	fmt.Println("Status da resposta:", resp.Status)
}
