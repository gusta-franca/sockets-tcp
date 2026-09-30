## Sobre o programa

Servidor TCP desenvolvido em C# que permite que múltiplos clientes se conectem simultaneamente e realizem operações sobre um sistema de arquivos.

## Requisitos

- .NET SDK instalado

## Dependências

O projeto utiliza:

- Serilog
- Serilog.Sinks.Console
- DotNetEnv

As dependências já estão definidas no arquivo `.csproj` e serão restauradas automaticamente pelo .NET.

## Configuração

Crie um arquivo `.env` na raiz do projeto:

```env
SERVER_ROOT=/caminho/para/o/diretorio
```

## Como executar

Primeiro, restaure as dependências:

```bash
dotnet restore
```

Inicie o servidor:

```bash
dotnet run
```

Em outro terminal, inicie o cliente:

```bash
dotnet run
```

## Comandos

| Comando | Descrição |
|---|---|
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