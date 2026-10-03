/*
    Data de criação: 30/09/2026
    Aluna: Caroline Marques Lau
    Arquivo do cliente que se conecta a um servidor de arquivos
*/

using System.Net.Sockets;
using System.Security.Cryptography;
using System.Text;
using Serilog;

class Client
{
    static void Main()
    {
        Log.Logger = new LoggerConfiguration().WriteTo.Console().CreateLogger();

        try
        {
            Log.Information("Conectando ao servidor...");
            using TcpClient client = new();
            client.Connect("127.0.0.1", 5000);
            Log.Information("Conectado ao servidor");
            NetworkStream stream = client.GetStream();
            Utf8Protocol protocol = new Utf8Protocol(stream);
            string response;
            while (true)
            {
                //Autenticação
                Console.Write("Usuário: ");
                string username = Console.ReadLine()!;
                Console.Write("Senha: ");
                string password = ReadPassword();
                string passwordHash = HashPassword(password);

                string connectMessage = $"CONNECT {username}, {passwordHash}";

                Log.Information("Enviando comando CONNECT para o servidor");
                protocol.WriteString(connectMessage);

                response = protocol.ReadString();

                Log.Information("Resposta do servidor: {Resposta}", response);
                if (response != "SUCCESS")
                {
                    Log.Error("Falha na autenticação");
                    continue;
                }
                else
                {
                    Log.Information("Autenticação realizada com sucesso");
                    break;
                }
            }

            // comandos após a autenticação
            while (true)
            {
                Console.Write("> ");
                string command = Console.ReadLine()!;

                Log.Information("Enviando comando: {Comando}", command);

                protocol.WriteString(command);

                if (command == "EXIT")
                {
                    Log.Information("Encerrando conexão");
                    break;
                }

                if (command == "GETFILES" || command == "GETDIRS") // pois recebe mais de uma mensagem
                {
                    int count = int.Parse(protocol.ReadString());
                    for (int i = 0; i < count; i++)
                    {
                        string item = protocol.ReadString();
                        Console.WriteLine(item);
                    }
                }
                else
                {
                    response = protocol.ReadString();
                    Log.Information("Resposta recebida: {Resposta}", response);
                }
            }
            Log.Information("Cliente encerrado");
        }
        catch (Exception ex)
        {
            Log.Error(ex, "Erro no cliente");
        }
        finally
        {
            Log.CloseAndFlush();
        }

    }

    static string HashPassword(string password) // gera o hash da senha
    {
        byte[] bytes = Encoding.UTF8.GetBytes(password);
        byte[] hash = SHA512.HashData(bytes);

        return Convert.ToHexString(hash).ToLower();
    }

    static string ReadPassword() // função feita por IA para esconder a senha
    {
        string password = "";
        while (true)
        {
            ConsoleKeyInfo key = Console.ReadKey(intercept: true);
            if (key.Key == ConsoleKey.Enter)
            {
                Console.WriteLine();
                break;
            }
            if (key.Key == ConsoleKey.Backspace)
            {
                if (password.Length > 0)
                {
                    password = password[..^1];
                }
                continue;
            }
            if (!char.IsControl(key.KeyChar))
            {
                password += key.KeyChar;
            }
        }
        return password;
    }
}