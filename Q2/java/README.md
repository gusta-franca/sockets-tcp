# Servidor e Cliente TCP de Arquivos (Java)

Servidor e cliente TCP desenvolvidos em Java que permitem que múltiplos clientes se conectem simultaneamente e gerenciem um conjunto de arquivos remotos, usando um protocolo binário compatível com os servidores e clientes em C# e Go.

## Requisitos

- **Java 17 ou superior** (JDK)

## Variável de ambiente

Configure sua variável de ambiente:

```bash
cp .env.example .env
```

| Variável | Descrição | Padrão |
|---|---|---|
| `SERVER_ADDR` | Endereço do servidor no formato `host:porta` | `127.0.0.1:5000` |
| `BASE_DIR` | Pasta onde o servidor guarda os arquivos | `./storage` |
| `DOWNLOAD_DIR` | Pasta onde o cliente grava os arquivos baixados | `./downloads` |


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

O servidor escuta na porta de `SERVER_ADDR` em todas as interfaces de rede. O cliente usa o host e a porta do mesmo campo para se conectar.

### Conectando de outra máquina

No `.env` da máquina do **cliente**, troque o endereço pelo IP da máquina do servidor:

```dotenv
SERVER_ADDR=192.168.0.10:5000
```

## Comandos

O nome do arquivo é digitado na mesma linha do comando, no formato `COMANDO <arquivo>` (ex.: `ADDFILE ./docs/relatorio.pdf`). O `GETFILESLIST` é digitado sozinho.

| Comando | Descrição | Resultado |
|---|---|---|
| `ADDFILE` | Envia um arquivo ao servidor. | `Arquivo enviado com sucesso!` ou erro |
| `DELETE` | Remove um arquivo do servidor | `Arquivo deletado com sucesso!` ou erro |
| `GETFILESLIST` | Lista os arquivos do servidor | quantidade, depois um nome por linha |
| `GETFILE` | Baixa um arquivo do servidor para `DOWNLOAD_DIR` | `Arquivo ... baixado com sucesso!` ou erro |
| `EXIT` | Encerra a conexão | (sem resposta) |

Comandos não são sensíveis a maiúsculas e minúsculas, mas os nomes de arquivo são.

## Protocolo

Todos os campos numéricos são enviados em **Big Endian**, e os arquivos trafegam **byte a byte**.

### Requisição

| Campo | Tamanho |
|---|---:|
| Message Type (`0x01`) | 1 byte |
| Command Identifier | 1 byte |
| Filename Size | 1 byte |
| Filename | `Filename Size` bytes |
| File Size | 4 bytes* |
| File | `File Size` bytes* |

\* Presentes apenas no `ADDFILE`.

### Resposta

| Campo | Tamanho |
|---|---:|
| Message Type (`0x02`) | 1 byte |
| Command Identifier | 1 byte |
| Status Code | 1 byte |
| Número de arquivos | 2 bytes* |
| Nomes: tamanho (1 byte) + nome, repetido | variável* |
| File Size | 4 bytes** |
| File | `File Size` bytes** |

\* Apenas no `GETFILESLIST` com sucesso. \** Apenas no `GETFILE` com sucesso. Em caso de erro, a resposta contém somente o cabeçalho.

### Códigos

| Comando | Valor | Status | Valor |
|---|---:|---|---:|
| `ADDFILE` | `0x01` | `SUCCESS` | `0x01` |
| `DELETE` | `0x02` | `ERROR` | `0x02` |
| `GETFILESLIST` | `0x03` | | |
| `GETFILE` | `0x04` | | |

## Logs

O servidor registra cada conexão, desconexão e comando executado no console e no arquivo `server.log`, usando `java.util.logging`.