## Sobre o programa

Aplicação cliente-servidor desenvolvida em Go para transferência e gerenciamento de arquivos por meio de uma conexão TCP.

Os comandos executados em uma sessão serão registrados em `logs/server.log`.

## Dependências

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

## Protocolo

Todos os campos numéricos são enviados em **Big Endian**.

### Request

| Campo | Tamanho |
|---|---:|
| MessageType | 1 byte |
| CommandIdentifier | 1 byte |
| FileNameSize | 1 byte |
| FileName | `FileNameSize` bytes |
| FileSize | 4 bytes* |
| File | `FileSize` bytes* |

\* Presentes apenas no comando `ADDFILE`.

### Response

| Campo | Tamanho |
|---|---:|
| MessageType | 1 byte |
| CommandIdentifier | 1 byte |
| StatusCode | 1 byte |
| FileAmount | 2 bytes* |
| FileNames | variável* |
| FileSize | 4 bytes* |
| File | `FileSize` bytes* |

\* Dependem do comando.

### Comandos

| Comando | Valor |
|---|---:|
| ADDFILE | `0x01` |
| DELETE | `0x02` |
| GETFILESLIST | `0x03` |
| GETFILE | `0x04` |

### Status

| Status | Valor |
|---|---:|
| SUCCESS | `0x01` |
| ERROR | `0x02` |

### Tipos de mensagem

| Tipo | Valor |
|---|---:|
| Request | `0x01` |
| Answer | `0x02` |
