# SendWhats - Cliente de Envio de Mensagens WhatsApp

Aplicação em Go para envio de mensagens WhatsApp via linha de comando, com suporte para múltiplas APIs.

## Funcionalidades

- Envio de mensagens para contatos individuais
- Envio de mensagens para grupos
- Suporte para WhatsApp-Web.js (API customizada)
- Suporte para EvolutionAPI
- Listagem de grupos disponíveis

## Uso

### Enviar Mensagem

```bash
./sendwhats <numero_ou_id> <titulo> <mensagem>
```

**Exemplos:**
```bash
# Enviar para contato individual
./sendwhats 5511999999999 "Título" "Sua mensagem aqui"

# Enviar para grupo (use o ID do grupo)
./sendwhats 120363123456789012@g.us "Título" "Sua mensagem aqui"

# Mensagem com quebras de linha
./sendwhats 5511999999999 "Título" "Linha 1\\nLinha 2"
```

### Listar Grupos

```bash
./sendwhats -groups
```

## Configuração

O arquivo `sendwhats.conf` deve estar no mesmo diretório do executável.

### Configuração para WhatsApp-Web.js

```ini
api_type=whatsapp-web
site=http://192.168.10.32:4000
```

### Configuração para EvolutionAPI

```ini
api_type=evolutionapi
evolution_url=http://localhost:8080
evolution_apikey=sua_chave_api_aqui
evolution_instance=nome_da_instancia
```

**Parâmetros:**
- `api_type`: Tipo de API a ser usada (`whatsapp-web` ou `evolutionapi`)
- `site`: URL do servidor WhatsApp-Web.js (apenas para `whatsapp-web`)
- `evolution_url`: URL base da EvolutionAPI (apenas para `evolutionapi`)
- `evolution_apikey`: Chave de API da EvolutionAPI (apenas para `evolutionapi`)
- `evolution_instance`: Nome da instância configurada na EvolutionAPI (apenas para `evolutionapi`)

## Configuração da EvolutionAPI

### Instalação via Docker

A forma mais simples de instalar a EvolutionAPI é usando Docker:

```bash
docker run -d \
  --name evolution_api \
  -p 8080:8080 \
  -e AUTHENTICATION_API_KEY=sua_chave_secreta_aqui \
  -e CONFIG_SESSION_PHONE_CLIENT=Chrome \
  -e CONFIG_SESSION_PHONE_NAME=chrome \
  atendai/evolution-api:latest
```

**Parâmetros importantes:**
- `AUTHENTICATION_API_KEY`: Define a chave de autenticação da API (use esta chave no `evolution_apikey` do `sendwhats.conf`)
- `8080`: Porta padrão da EvolutionAPI (ajuste se necessário)

### Criar uma Instância

Após iniciar o container, você precisa criar uma instância do WhatsApp:

1. **Acesse a interface web da EvolutionAPI:**
   ```
   http://localhost:8080
   ```

2. **Ou use a API diretamente para criar uma instância:**
   ```bash
   curl -X POST http://localhost:8080/instance/create \
     -H "apikey: sua_chave_secreta_aqui" \
     -H "Content-Type: application/json" \
     -d '{
       "instanceName": "minha_instancia",
       "token": "token_opcional",
       "qrcode": true
     }'
   ```

3. **Escaneie o QR Code** que será gerado para conectar sua conta WhatsApp.

4. **Use o nome da instância** (`minha_instancia` no exemplo) no campo `evolution_instance` do `sendwhats.conf`.

### Verificar Status da Instância

```bash
curl -X GET http://localhost:8080/instance/fetchInstances \
  -H "apikey: sua_chave_secreta_aqui"
```

### Documentação Completa

Para mais informações sobre a EvolutionAPI, consulte:
- [Documentação Oficial](https://doc.evolution-api.com/)
- [GitHub da EvolutionAPI](https://github.com/EvolutionAPI/evolution-api)

## Formato de Números e IDs

### Contatos Individuais

- **WhatsApp-Web.js**: Aceita números em qualquer formato (ex: `5511999999999`, `+55 11 99999-9999`)
- **EvolutionAPI**: Recomenda-se usar formato internacional sem caracteres especiais (ex: `5511999999999`)

### Grupos

- **WhatsApp-Web.js**: Use o ID do grupo retornado pelo comando `./sendwhats -groups`
- **EvolutionAPI**: Use o ID do grupo no formato `{groupId}@g.us` (ex: `120363123456789012@g.us`)

Para obter os IDs dos grupos, use o comando:
```bash
./sendwhats -groups
```

## Compilação

```bash
go build -o sendwhats main.go
```

## Requisitos

- Go 1.21.1 ou superior
- Servidor WhatsApp-Web.js OU EvolutionAPI configurado e rodando

## Estrutura do Projeto

```
send_whats/
├── main.go              # Código principal
├── sendwhats.conf       # Arquivo de configuração
├── go.mod              # Dependências Go
├── go.sum              # Checksums das dependências
└── README.md           # Este arquivo
```

## Troubleshooting

### Erro: "Endereço do servidor não especificado"
- Verifique se o arquivo `sendwhats.conf` existe no mesmo diretório do executável
- Verifique se a chave `site` (para whatsapp-web) ou `evolution_url` (para evolutionapi) está configurada

### Erro: "API Key da EvolutionAPI não especificada"
- Verifique se `evolution_apikey` está configurado no `sendwhats.conf`
- A chave deve ser a mesma definida em `AUTHENTICATION_API_KEY` ao iniciar o Docker

### Erro: "Nome da instância EvolutionAPI não especificado"
- Verifique se `evolution_instance` está configurado no `sendwhats.conf`
- O nome deve corresponder a uma instância criada na EvolutionAPI

### Mensagem não é enviada
- Verifique se o servidor (WhatsApp-Web.js ou EvolutionAPI) está rodando
- Verifique se a instância da EvolutionAPI está conectada (status `open`)
- Verifique os logs do servidor para mais detalhes

### Grupos não aparecem
- Para EvolutionAPI, certifique-se de que a instância está conectada e sincronizada
- Algumas APIs podem demorar para sincronizar grupos após a conexão inicial

## Licença

Este projeto é fornecido como está, sem garantias.
