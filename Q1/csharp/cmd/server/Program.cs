/*
    Data de criação: 30/09/2026
    Aluna: Caroline Marques Lau
    Arquivo do servidor responsável por processar as mensagens dos clientes
*/

using System.Net;
using System.Net.Sockets;
using Serilog;
using Serilog.Context;

class Server
{
    static void Main()
    {
        DotNetEnv.Env.Load();
        Log.Logger = new LoggerConfiguration()
        .Enrich.FromLogContext()
        // para diferenciar os clientes
        .WriteTo.Console(outputTemplate: "[{Timestamp:HH:mm:ss} {Level:u3}] [{Cliente}] {Message:lj}{NewLine}{Exception}")
        .CreateLogger();

        // escuta na porta 5000  em todas as interfaces de rede da máquina
        TcpListener server = new(IPAddress.Any, 5000);
        server.Start(); // começa a escutar conexões TCP na porta
        Log.Information("Servidor iniciado na porta 5000");
        AuthenticationManager authenticationManager = new AuthenticationManager();

        while (true)
        {
            TcpClient client = server.AcceptTcpClient();
            Task.Run(() => HandleClient(client, authenticationManager));
        }
    }

    static void HandleClient(TcpClient client, AuthenticationManager authenticationManager)
    {
        using (LogContext.PushProperty("Cliente", client.Client.RemoteEndPoint))
        {
            try
            {
                Log.Information("Cliente conectado: {Cliente}", client.Client.RemoteEndPoint);

                using (client) // quando o bloco termina, o recurso é descartado
                using (NetworkStream stream = client.GetStream())
                {
                    FileSystemManager fs = new FileSystemManager();
                    Utf8Protocol protocol = new Utf8Protocol(stream);
                    bool authenticated = false;
                    // bloco de autenticação do cliente
                    while (true)
                    {
                        string message;

                        try
                        {
                            message = protocol.ReadString();
                        }
                        catch (EndOfStreamException)
                        {
                            Log.Information("Cliente desconectou durante autenticação");
                            break;
                        }

                        if (message.StartsWith("CONNECT "))
                        {
                            string credentials = message["CONNECT ".Length..];
                            string[] parts = credentials.Split(',');
                            if (parts.Length != 2)
                            {
                                protocol.WriteString("Comando inválido");
                                continue;
                            }

                            string username = parts[0].Trim();
                            string passwordHash = parts[1].Trim();

                            authenticated = authenticationManager.Authenticate(
                                username,
                                passwordHash
                            );

                            if (authenticated)
                            {
                                protocol.WriteString("SUCCESS");
                                Log.Information("Cliente logado com sucesso!");
                                break;
                            }
                            else
                            {
                                protocol.WriteString("ERROR");
                                Log.Information("Senha ou usuários inválidos");
                            }
                        }
                        else if (message == "EXIT")
                        {
                            break;
                        }
                        else
                        {
                            protocol.WriteString("Realize a sua autenticação");
                            Log.Information("Comando de autenticação inválido");
                        }
                    }

                    // caso esteja autenticado, recebe comandos
                    try
                    {
                        while (authenticated)
                        {
                            string message = protocol.ReadString();
                            Log.Information("Mensagem recebida: {Mensagem}", message);
                            if (message == "PWD")
                            {
                                string path = fs.Pwd();
                                protocol.WriteString(path);

                            }
                            else if (message.StartsWith("CHDIR "))
                            {
                                string path = message[6..];
                                bool success = fs.Chdir(path);
                                protocol.WriteString(success ? "SUCCESS" : "ERROR");
                            }
                            else if (message == "GETFILES")
                            {
                                string[] files = fs.GetFiles();
                                protocol.WriteString(files.Length.ToString());
                                foreach (string file in files)
                                {
                                    protocol.WriteString(file);
                                }
                            }
                            else if (message == "GETDIRS")
                            {
                                string[] dirs = fs.GetDirs();
                                protocol.WriteString(dirs.Length.ToString());
                                foreach (string dir in dirs)
                                {
                                    protocol.WriteString(dir);
                                }
                            }
                            else if (message == "EXIT")
                            {
                                authenticated = false;
                            }
                            else if (message == "HELP")
                            {
                                protocol.WriteString(
                                    "\nComandos disponíveis:\n" +
                                    "PWD - Mostra o diretório atual.\n" +
                                    "CHDIR <caminho> - Muda o diretório atual.\n" +
                                    "GETFILES - Lista os arquivos do diretório atual.\n" +
                                    "GETDIRS - Lista os diretórios do diretório atual.\n" +
                                    "EXIT - Encerra a conexão com o servidor.\n" +
                                    "HELP - Mostra esta lista de comandos."
                                );
                            }
                            else
                            {
                                protocol.WriteString("ERROR");
                            }
                        }
                        Log.Information("Cliente desconectado!");
                    }
                    catch (EndOfStreamException)
                    {
                        Log.Information("Cliente desconectou inesperadamente");
                    }
                }
            }
            catch (Exception ex)
            {
                Log.Error(ex, "Erro na conexão do cliente");
            }
        }
    }
}