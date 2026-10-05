# Servidor e Cliente TCP (Java)

Servidor e cliente TCP desenvolvidos em Java que permite que múltiplos clientes se conectem simultaneamente e realizem operações sobre um sistema de arquivos.

## Requisitos

- **Java 11 ou superior** (JDK)

## Variável de ambiente

Configure sua variável de ambiente:
cp .env.example .env

## Compilação
 
Na raiz do projeto:
 
```bash
javac -d out src/common/*.java src/server/*.java src/client/*.java
```
 
Os arquivos compilados vão para a pasta `out`, separados por pacote (`common`, `server` e `client`).
 
## Execução
 
**1. Em um terminal, inicie o servidor:**
 
```bash
java -cp out server.Server
```
 
**2. Em outro terminal, inicie o cliente:**
 
```bash
java -cp out client.Client
```
 
### Conectando de outra máquina
 
No `.env` da máquina do **cliente**, troque o endereço pelo IP da máquina do servidor:
 
```dotenv
SERVER_ADDR=127.0.0.1:5000
```

## Comandos

| Comando | Descrição | Resposta |
|---|---|---|
| `CONNECT usuario, senha` | Autentica o usuário. Note a **vírgula** e o **espaço** depois dela | `SUCCESS` ou `ERROR` |
| `PWD` | Mostra o diretório atual | caminho completo, com `/` |
| `CHDIR diretorio` | Muda o diretório atual | `SUCCESS` ou `ERROR` |
| `GETFILES` | Lista os arquivos do diretório atual | quantidade, depois um nome por mensagem |
| `GETDIRS` | Lista os subdiretórios do diretório atual | quantidade, depois um nome por mensagem |
| `EXIT` | Encerra a conexão | (sem resposta) |

Antes do `CONNECT`, qualquer comando (exceto `EXIT`) retorna `ERROR`.

### Usuários cadastrados

| Usuário | Senha |
|---|---|
| `Carol` | `123mudar` |
| `Gustavo` | `123456` |
| `Maria` | `admin` |

Os nomes diferenciam maiúsculas de minúsculas (case sensitive).