package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/ini.v1"
)

type Message struct {
	Number  string `json:"number"`
	Message string `json:"message"`
}

type Group struct {
	ID   string `json:"id"`
	Name string `json:"name,omitempty"`
}

var serverURL string

func loadConfig() {
	executable, err := os.Executable()
	if err != nil {
		fmt.Println("Erro ao obter o caminho do executável:", err)
		os.Exit(1)
	}
	executablePath := filepath.Dir(executable)
	configPath := filepath.Join(executablePath, "sendwhats.conf")

	cfg, err := ini.Load(configPath)
	if err != nil {
		fmt.Println("Erro ao carregar arquivo de configuração:", err)
		os.Exit(1)
	}
	serverURL = cfg.Section("").Key("site").String()
	if serverURL == "" {
		fmt.Println("Endereço do servidor não especificado no arquivo de configuração.")
		os.Exit(1)
	}
}

func sendMessage(phone, title, message string) {
	url := serverURL + "/api/message"

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

func getGroups() {
	url := serverURL + "/api/groups"

	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Erro ao fazer a requisição:", err)
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Erro ao ler a resposta:", err)
		return
	}

	var groups []Group
	err = json.Unmarshal(body, &groups)
	if err != nil {
		fmt.Println("Erro ao decodificar resposta JSON:", err)
		return
	}

	fmt.Println("Lista de Grupos:")
	for _, group := range groups {
		if group.Name != "" {
			fmt.Printf("ID: %s - Nome: %s\n", group.ID, group.Name)
		} else {
			fmt.Printf("ID: %s - Nome: (sem nome)\n", group.ID)
		}
	}
}

func printUsage() {
	fmt.Println("Uso:")
	fmt.Println("  ./sendwhats -groups                     Lista os grupos e contatos")
	fmt.Println("  ./sendwhats <phone/id> <title> <message>   Envia uma mensagem")
}

func main() {
	// Carrega a configuração
	loadConfig()

	// Define a flag para listar grupos
	listGroups := flag.Bool("groups", false, "Listar grupos")
	flag.Parse()

	if *listGroups {
		getGroups()
		return
	}

	// Parâmetros para enviar mensagem
	if len(flag.Args()) < 3 {
		printUsage()
		return
	}

	phone := flag.Arg(0)
	title := flag.Arg(1)
	message := strings.Join(flag.Args()[2:], " ")

	sendMessage(phone, title, message)
}
