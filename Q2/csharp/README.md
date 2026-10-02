## Sobre o programa

Aplicação cliente-servidor desenvolvida em C# para transferência e gerenciamento de arquivos por meio de uma conexão TCP.

## Dependências

- .NET
- DotNetEnv
- Serilog
- Serilog.Sinks.Console

Para restaurar as dependências:

```bash
dotnet restore
```



## Configuração

O servidor e o cliente possuem arquivos `.env` separados.

### Servidor

Crie um arquivo `.env` no projeto do servidor:

```env
SERVER_ROOT=/caminho/do/diretorio
```

`SERVER_ROOT` define o diretório utilizado pelo servidor para armazenar e gerenciar os arquivos.

### Cliente

Crie um arquivo `.env` no projeto do cliente:

```env
CLIENT_ROOT=/caminho/do/diretorio
```

`CLIENT_ROOT` define o diretório utilizado pelo cliente para armazenar e buscar arquivos.


## Como rodar

Inicie primeiro o servidor e, em seguida, execute o cliente:

```bash
dotnet run
```

O cliente se conecta ao servidor em:

```text
127.0.0.1:5000
```

Para compilar:

```bash
dotnet build
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

### Answer

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