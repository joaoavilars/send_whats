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

// Estrutura para mensagem EvolutionAPI
type EvolutionMessage struct {
	Number      string `json:"number"`
	Text        string `json:"text"`
	InstanceName string `json:"instanceName"`
}

// Estrutura de resposta EvolutionAPI para grupos
type EvolutionGroup struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type EvolutionGroupsResponse struct {
	Groups []EvolutionGroup `json:"groups"`
}

// Configuração da aplicação
type Config struct {
	APIType          string // "whatsapp-web" ou "evolutionapi"
	ServerURL        string // URL do servidor WhatsApp-Web.js
	EvolutionURL     string // URL base da EvolutionAPI
	EvolutionAPIKey  string // API Key da EvolutionAPI
	EvolutionInstance string // Nome da instância EvolutionAPI
}

var config Config

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

	// Lê o tipo de API (padrão: whatsapp-web)
	config.APIType = cfg.Section("").Key("api_type").MustString("whatsapp-web")
	
	// Configurações para WhatsApp-Web.js
	config.ServerURL = cfg.Section("").Key("site").String()
	
	// Configurações para EvolutionAPI
	config.EvolutionURL = cfg.Section("").Key("evolution_url").String()
	config.EvolutionAPIKey = cfg.Section("").Key("evolution_apikey").String()
	config.EvolutionInstance = cfg.Section("").Key("evolution_instance").String()

	// Validação baseada no tipo de API
	if config.APIType == "whatsapp-web" {
		if config.ServerURL == "" {
			fmt.Println("Endereço do servidor (site) não especificado no arquivo de configuração.")
			os.Exit(1)
		}
	} else if config.APIType == "evolutionapi" {
		if config.EvolutionURL == "" {
			fmt.Println("URL da EvolutionAPI (evolution_url) não especificada no arquivo de configuração.")
			os.Exit(1)
		}
		if config.EvolutionAPIKey == "" {
			fmt.Println("API Key da EvolutionAPI (evolution_apikey) não especificada no arquivo de configuração.")
			os.Exit(1)
		}
		if config.EvolutionInstance == "" {
			fmt.Println("Nome da instância EvolutionAPI (evolution_instance) não especificado no arquivo de configuração.")
			os.Exit(1)
		}
	} else {
		fmt.Printf("Tipo de API inválido: %s. Use 'whatsapp-web' ou 'evolutionapi'.\n", config.APIType)
		os.Exit(1)
	}
}

// Envia mensagem via WhatsApp-Web.js
func sendMessageWhatsAppWeb(phone, title, message string) {
	url := config.ServerURL + "/api/message"

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
	
	// Lê o corpo da resposta para mostrar erros, se houver
	body, _ := ioutil.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		fmt.Printf("Erro na resposta: %s\n", string(body))
	}
}

// Envia mensagem via EvolutionAPI (suporta contatos e grupos)
func sendMessageEvolutionAPI(phone, title, message string) {
	// Remove barra final da URL se existir
	baseURL := strings.TrimSuffix(config.EvolutionURL, "/")
	url := baseURL + "/message/sendText/" + config.EvolutionInstance

	// Substitui \n por quebras de linha reais
	message = strings.ReplaceAll(message, "\\n", "\n")

	// Monta a mensagem completa com título
	fullMessage := fmt.Sprintf("%s\n\n%s", title, message)

	// Prepara o payload da EvolutionAPI
	payload := map[string]interface{}{
		"number": phone,
		"text":   fullMessage,
	}

	// Converte para JSON
	jsonBody, err := json.Marshal(payload)
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
	req.Header.Set("apikey", config.EvolutionAPIKey)

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
	
	// Lê o corpo da resposta para mostrar erros ou sucesso
	body, _ := ioutil.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		fmt.Printf("Erro na resposta: %s\n", string(body))
	} else {
		// Mostra resposta de sucesso se disponível
		if len(body) > 0 {
			fmt.Printf("Resposta: %s\n", string(body))
		}
	}
}

// Função principal de envio que roteia para a API correta
func sendMessage(phone, title, message string) {
	if config.APIType == "evolutionapi" {
		sendMessageEvolutionAPI(phone, title, message)
	} else {
		sendMessageWhatsAppWeb(phone, title, message)
	}
}

// Lista grupos via WhatsApp-Web.js
func getGroupsWhatsAppWeb() {
	url := config.ServerURL + "/api/groups"

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

// Lista grupos via EvolutionAPI
func getGroupsEvolutionAPI() {
	// Remove barra final da URL se existir
	baseURL := strings.TrimSuffix(config.EvolutionURL, "/")
	url := baseURL + "/group/fetchAllGroups/" + config.EvolutionInstance

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("Erro ao criar requisição HTTP:", err)
		return
	}
	req.Header.Set("apikey", config.EvolutionAPIKey)

	client := &http.Client{}
	resp, err := client.Do(req)
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

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Erro ao buscar grupos. Status: %s\n", resp.Status)
		fmt.Printf("Resposta: %s\n", string(body))
		return
	}

	// Tenta decodificar como array direto ou objeto com grupos
	var groups []EvolutionGroup
	err = json.Unmarshal(body, &groups)
	if err != nil {
		// Tenta como objeto com campo groups
		var response EvolutionGroupsResponse
		err2 := json.Unmarshal(body, &response)
		if err2 != nil {
			fmt.Println("Erro ao decodificar resposta JSON:", err)
			fmt.Printf("Resposta recebida: %s\n", string(body))
			return
		}
		groups = response.Groups
	}

	fmt.Println("Lista de Grupos:")
	if len(groups) == 0 {
		fmt.Println("Nenhum grupo encontrado.")
		return
	}
	
	for _, group := range groups {
		if group.Name != "" {
			fmt.Printf("ID: %s - Nome: %s\n", group.ID, group.Name)
		} else {
			fmt.Printf("ID: %s - Nome: (sem nome)\n", group.ID)
		}
	}
}

// Função principal que roteia para a API correta
func getGroups() {
	if config.APIType == "evolutionapi" {
		getGroupsEvolutionAPI()
	} else {
		getGroupsWhatsAppWeb()
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
