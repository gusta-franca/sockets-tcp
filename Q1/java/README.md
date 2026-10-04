# Servidor e Cliente TCP (Java)

Servidor e cliente TCP desenvolvidos em Java que permite que múltiplos clientes se conectem simultaneamente e realizem operações sobre um sistema de arquivos.

## Requisitos

- **Java 11 ou superior** (JDK)

## Variável de ambiente

Configure sua variável de ambiente:
cp .env.example .env

## Compilação

Dentro da pasta com os arquivos:

```bash
javac *.java
```

Isso gera os arquivos `.class` na mesma pasta.

## Execução

**1. Em um terminal, inicie o servidor:**

```bash
java Server
```

Ele fica escutando na porta **5000** e exibe `Servidor aguardando conexão ...`.

**2. Em outro terminal, inicie o cliente:**

```bash
java Client
```

Para testar vários clientes ao mesmo tempo, abra mais terminais e rode `java Client` em cada um.

### Conectando de outra máquina

No `Client.java`, altere o endereço do servidor:

```java
String serverHost = "127.0.0.1";   // troque pelo IP da máquina do servidor
int serverPort = 5000;
```

Depois recompile (`javac *.java`). O servidor e o cliente precisam usar a mesma porta.

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

A senha é digitada normalmente no cliente. O **cliente converte para SHA-512** antes de enviar, e o servidor compara apenas os hashes.

### Usuários cadastrados

| Usuário | Senha |
|---|---|
| `Carol` | `123mudar` |
| `Gustavo` | `123456` |
| `Maria` | `admin` |

Os nomes diferenciam maiúsculas de minúsculas (case sensitive).

### Uso

Antes, crie uma pasta de teste (o servidor cria `/tmp/<usuario>` no primeiro login):

```bash
mkdir -p /tmp/Carol/docs
touch /tmp/Carol/a.txt /tmp/Carol/b.txt