/*
    Data de criação: 30/09/2026
    Aluna: Caroline Marques Lau
    Arquivo do servidor responsável por processar as mensagens dos clientes
*/

using System.Net;
using System.Net.Sockets;
using Serilog;
using Serilog.Context;
using System.Text;

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
        FileSystemManager fs = new();

        while (true)
        {
            TcpClient client = server.AcceptTcpClient();
            Task.Run(() => HandleClient(client, fs));
        }
    }

    static void HandleClient(TcpClient client, FileSystemManager fs)
    {
        using (LogContext.PushProperty("Cliente", client.Client.RemoteEndPoint))
        {
            Log.Information("Cliente conectado");
            using (client)
            using (NetworkStream stream = client.GetStream())
            {
                Protocol protocol = new(stream);
                while (true)
                {
                    Request request;
                    try
                    {
                        request = protocol.ReadRequest();
                        Command command = request.CommandIdentifier;
                        string fileName = Encoding.UTF8.GetString(request.FileName);
                        switch (command)
                        {
                            case Command.ADDFILE:
                                try
                                {
                                    fs.AddFile(fileName, request.File, request.FileSize);
                                    Answer answer = new()
                                    {
                                        MessageType = MessageType.Answer,
                                        CommandIdentifier = Command.ADDFILE,
                                        StatusCode = StatusCode.SUCCESS
                                    };

                                    protocol.SendAnswer(answer);

                                    Log.Information(
                                        "Arquivo {FileName} adicionado com sucesso",
                                        fileName
                                    );
                                }
                                catch (Exception ex)
                                {
                                    Answer answer = new()
                                    {
                                        MessageType = MessageType.Answer,
                                        CommandIdentifier = Command.ADDFILE,
                                        StatusCode = StatusCode.ERROR
                                    };

                                    protocol.SendAnswer(answer);

                                    Log.Error(
                                        ex,
                                        "Erro ao adicionar o arquivo {FileName}",
                                        fileName
                                    );
                                }
                                break;
                            case Command.DELETE:
                                try
                                {
                                    fs.DeleteFile(fileName);
                                    Answer answer = new()
                                    {
                                        MessageType = MessageType.Answer,
                                        CommandIdentifier = Command.DELETE,
                                        StatusCode = StatusCode.SUCCESS
                                    };

                                    protocol.SendAnswer(answer);

                                    Log.Information(
                                        "Arquivo {FileName} deletado com sucesso",
                                        fileName
                                    );
                                }
                                catch (Exception ex)
                                {
                                    Answer answer = new()
                                    {
                                        MessageType = MessageType.Answer,
                                        CommandIdentifier = Command.DELETE,
                                        StatusCode = StatusCode.ERROR
                                    };

                                    protocol.SendAnswer(answer);

                                    Log.Error(
                                        ex,
                                        "Erro ao deletar o arquivo {FileName}",
                                        fileName
                                    );
                                }
                                break;
                            case Command.GETFILE:
                                try
                                {
                                    byte[]? file = fs.GetFile(fileName);
                                    Answer answer = new()
                                    {
                                        MessageType = MessageType.Answer,
                                        CommandIdentifier = Command.GETFILE,
                                        StatusCode = StatusCode.SUCCESS,
                                        FileSize = (uint)file!.Length,
                                        File = file
                                    };

                                    protocol.SendAnswer(answer);

                                    Log.Information(
                                        "Arquivo {FileName} enviado com sucesso",
                                        fileName
                                    );
                                }
                                catch (Exception ex)
                                {
                                    Answer answer = new()
                                    {
                                        MessageType = MessageType.Answer,
                                        CommandIdentifier = Command.GETFILE,
                                        StatusCode = StatusCode.ERROR
                                    };

                                    protocol.SendAnswer(answer);

                                    Log.Error(
                                        ex,
                                        "Erro ao enviar o arquivo {FileName}",
                                        fileName
                                    );
                                }
                                break;
                            case Command.GETFILESLIST:
                                try
                                {
                                    string[] fileList = fs.GetFilesList();
                                    Answer answer = new()
                                    {
                                        MessageType = MessageType.Answer,
                                        CommandIdentifier = Command.GETFILESLIST,
                                        StatusCode = StatusCode.SUCCESS,
                                        FileAmount = (ushort)fileList.Length,
                                        FileNames = fileList
                                    };

                                    protocol.SendAnswer(answer);

                                    Log.Information(
                                        "Lista de arquivos enviada com sucesso"
                                    );
                                }
                                catch (Exception ex)
                                {
                                    Answer answer = new()
                                    {
                                        MessageType = MessageType.Answer,
                                        CommandIdentifier = Command.GETFILESLIST,
                                        StatusCode = StatusCode.ERROR
                                    };

                                    protocol.SendAnswer(answer);

                                    Log.Error(
                                        ex,
                                        "Erro ao enviar lista de arquivos"
                                    );
                                }
                                break;
                            default:
                                Log.Information("Comando inválido");
                                break;
                        }
                    }
                    catch (EndOfStreamException)
                    {
                        Log.Information("Cliente desconectou");
                        break;
                    }

                }
            }
        }
    }
}