# Projeto: Servidor TCP

Alunos: Caroline Marques Lau, Gustavo Martins França e Maria Eduarda Bambini

Disciplina: Sistemas Distribuídos

## Questão 1

Mensagens em `String UTF`.

| Comando | Função |
|---|---|
| `CONNECT user, password` | Autentica usuário usando senha em SHA-512 |
| `PWD` | Retorna diretório atual |
| `CHDIR path` | Altera diretório |
| `GETFILES` | Lista arquivos |
| `GETDIRS` | Lista diretórios |
| `EXIT` | Encerra conexão |

Respostas: `SUCCESS` ou `ERROR`.
## Questão 2 

Esta aplicação implementa um servidor de arquivos remoto multiusuário utilizando comunicação TCP sob um protocolo binário customizado.

## Especificação do protocolo

### 1. Comandos suportados
- `0x01` - **ADDFILE**: Adiciona um arquivo no servidor.
- `0x02` - **DELETE**: Remove um arquivo do servidor.
- `0x03` - **GETFILESLIST**: Lista os nomes dos arquivos no servidor.
- `0x04` - **GETFILE**: Realiza o download de um arquivo.

### 2. Formato da requisição (Cliente -> Servidor)
- `1 byte`: Message Type (`0x01`)
- `1 byte`: Command Identifier (`0x01` a `0x04`)
- `1 byte`: Filename Size ($N$)
- `$N$ bytes`: Filename (0 a 255 bytes)

#### Campos adicionais de carga útil (Payload)
- **ADDFILE**: `4 bytes` (Tamanho do Arquivo em Big Endian) + `1 a 2^32 bytes` (Dados do arquivo).

### 3. Formato da resposta (Servidor -> Cliente)
- `1 byte`: Message Type (`0x02`)
- `1 byte`: Command Identifier (`0x01` a `0x04`)
- `1 byte`: Status Code (`1` = SUCCESS, `2` = ERROR)

#### Campos adicionais de carga útil (Payload)
- **GETFILESLIST**: `2 bytes` (Quantidade de arquivos em Big Endian) + Repetição de [`1 byte` (Tamanho do nome) + `N bytes` (Nome do arquivo)].
- **GETFILE**: `4 bytes` (Tamanho do arquivo em Big Endian) + `1 a 2^32 bytes` (Dados do arquivo).

### 4. Requisitos de implementação
- **Endianness**: Todos os inteiros de múltiplos bytes devem ser serializados/deserializados em **Big Endian**.
- **Download**: Os downloads efetuados pelo cliente são salvos automaticamente na pasta `Downloads/`.
- **Streaming**: Envio e recebimento de dados realizados byte a byte.
- **Logging**: O servidor registra eventos e erros usando bibliotecas nativas/padrão da linguagem.

## Estrutura

    sockets-tcp/
    ├── README.md
    ├── .gitignore
    │
    ├── Q1/                            # Questão 1: Protocolo de Texto (UTF-8)
    │   ├── c#/
    │   │   ├── README.md
    │   │   └── cmd/
    │   │       ├── client/            # Código do cliente em C#
    │   │       └── server/            # Código do servidor em C#
    │   ├── go/
    │   │   ├── README.md
    │   │   └── cmd/
    │   │       ├── client/            # Código do cliente em Go
    │   │       └── server/            # Código do servidor em Go
    │   └── java/
    │       ├── README.md
    │       └── cmd/
    │           ├── client/            # Código do cliente em Java
    │           └── server/            # Código do servidor em Java
    │
    └── Q2/                            # Questão 2: Protocolo Binário TCP
        ├── c#/
        │   ├── README.md
        │   └── cmd/
        │       ├── client/            # Código do cliente em C#
        │       └── server/            # Código do servidor em C#
        ├── go/
        │   ├── README.md
        │   └── cmd/
        │       ├── client/            # Código do cliente em Go
        │       └── server/            # Código do servidor em Go
        └── java/
            ├── README.md
            └── cmd/
                ├── client/            # Código do cliente em Java
                └── server/            # Código do servidor em Java