## Sobre o programa

Servidor TCP desenvolvido em Go que permite múltiplos clientes se conectem simultaneamente a um servidor e realizem operações sobre um sistema de arquivos.

Os comandos executados em uma sessão serão registrados em `logs/server.log`.

## Requisitos

- Go
- Makefile

## Configuração

Crie uma cópia de `.env.example`, a renomeie para `.env` e altere as variáveis conforme necessário.

## Como executar

Compile o servidor e o cliente:

```bash
make build
```

Inicie o servidor (terminal 1):

```bash
make run-server
```

Inicie o cliente (terminal 2):

```bash
make run-client
```

## Comandos

| Comando | Descrição |
|---|---|
| `CONNECT user,password` | Autentica um usuário |
| `PWD` | Mostra o diretório atual |
| `CHDIR <caminho>` | Altera o diretório atual |
| `GETFILES` | Lista os arquivos |
| `GETDIRS` | Lista os diretórios |
| `HELP` | Mostra os comandos disponíveis |
| `EXIT` | Encerra a conexão |

## Protocolo

A comunicação utiliza **TCP** e **UTF-8**.

Cada mensagem possui 4 bytes iniciais indicando seu tamanho, seguidos pela mensagem codificada em UTF-8:

```text
[4 bytes: tamanho da mensagem][mensagem UTF-8]
```

O tamanho da mensagem é representado em **Big-endian**.

Para comandos que retornam múltiplos resultados, como `GETFILES` e `GETDIRS`, o servidor envia primeiro a quantidade de itens e depois cada item individualmente.
