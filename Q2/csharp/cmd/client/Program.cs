/*
    Data de criação: 02/10/2026
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
        DotNetEnv.Env.Load();
        string rootDir = Environment.GetEnvironmentVariable("CLIENT_ROOT")
            ?? throw new Exception("CLIENT_ROOT não foi definida.");
        Log.Logger = new LoggerConfiguration().WriteTo.Console().CreateLogger();
        try
        {
            Log.Information("Conectando ao servidor...");
            using TcpClient client = new();
            client.Connect("127.0.0.1", 5000);
            Log.Information("Conectado ao servidor");
            NetworkStream stream = client.GetStream();
            Protocol protocol = new(stream);
            try
            {
                while (true)
                {
                    Console.Write("Digite o comando desejado: ");
                    string command = Console.ReadLine()!;
                    if (command.ToUpper() == "EXIT")
                    {
                        break;
                    }
                    Request request = new()
                    {
                        MessageType = MessageType.Request
                    };
                    if (Enum.TryParse(command, true, out Command commandEnum))
                    {
                        request.CommandIdentifier = commandEnum;
                        Log.Information("Comando válido: {commandEnum}", commandEnum);
                        string fileNameString = "";
                        if (commandEnum == Command.GETFILESLIST)
                        {
                            request.FileNameSize = 0;
                            request.FileName = new byte[1];
                        }
                        while (true && commandEnum != Command.GETFILESLIST)
                        {
                            Console.Write("Digite o nome do arquivo: ");
                            fileNameString = Console.ReadLine()!;
                            byte[] fileNameByte = Encoding.UTF8.GetBytes(fileNameString);
                            if (fileNameByte.Length <= byte.MaxValue)
                            {
                                request.FileNameSize = (byte)fileNameByte.Length;
                                request.FileName = fileNameByte;
                                break;
                            }
                            Log.Information("O nome do arquivo ultrapassa o limite de 255 bytes");
                        }
                        switch (commandEnum)
                        {
                            case Command.ADDFILE:
                                string filePath = Path.Combine(rootDir, fileNameString);
                                if (File.Exists(filePath))
                                {
                                    request.File = File.ReadAllBytes(filePath);
                                    request.FileSize = (uint)request.File.Length;
                                }
                                else
                                {
                                    Log.Information("O arquivo informado não existe e não pode ser enviado ao servidor");
                                    break;
                                }
                                protocol.SendRequest(request);
                                Answer answerAdd = protocol.ReadAnswer();
                                if (answerAdd.StatusCode == StatusCode.SUCCESS)
                                {
                                    Log.Information("Arquivo enviado com sucesso!");
                                }
                                else
                                {
                                    Log.Information("Não foi possível enviar o arquivo");
                                }
                                break;
                            case Command.DELETE:
                                protocol.SendRequest(request);
                                Answer answerDelete = protocol.ReadAnswer();
                                if (answerDelete.StatusCode == StatusCode.SUCCESS)
                                {
                                    Log.Information("Arquivo deletado com sucesso!");
                                }
                                else
                                {
                                    Log.Information("Não foi possível deletar o arquivo");
                                }
                                break;
                            case Command.GETFILE:
                                protocol.SendRequest(request);
                                Answer answerGetFile = protocol.ReadAnswer();
                                if (answerGetFile.StatusCode == StatusCode.SUCCESS)
                                {
                                    if (answerGetFile.FileSize == null || answerGetFile.FileSize == 0)
                                    {
                                        Log.Information("O tamanho do arquivo enviado é nulo ou igual a 0");
                                        break;
                                    }
                                    if (answerGetFile.File == null || answerGetFile.File.Length == 0)
                                    {
                                        Log.Information("O campo do arquivo veio vazio");
                                    }
                                    byte[] file = answerGetFile.File!;
                                    uint fileSize = answerGetFile.FileSize.Value;
                                    if (file.Length != fileSize)
                                    {
                                        Log.Information(
                                            "O tamanho informado ({FileSize}) não bate com o tamanho recebido ({ActualSize})",
                                            fileSize,
                                            file.Length
                                        );
                                    }
                                    string path = Path.Combine(rootDir, fileNameString);

                                    File.WriteAllBytes(path, file);

                                    Log.Information(
                                        "Arquivo {FileName} baixado com sucesso!",
                                        fileNameString
                                    );
                                }
                                else
                                {
                                    Log.Information("Não foi possível baixar o arquivo");
                                }
                                break;
                            case Command.GETFILESLIST:
                                protocol.SendRequest(request);
                                Answer answerList = protocol.ReadAnswer();
                                if (answerList.StatusCode == StatusCode.SUCCESS)
                                {
                                    if (answerList.FileNames == null || answerList.FileNames.Length == 0)
                                    {
                                        Log.Information("A lista de arquivos veio vazia");
                                        break;
                                    }
                                    string[] fileNames = answerList.FileNames;
                                    Log.Information("Lista de arquivos obtida com sucesso!");
                                    Log.Information("Arquivos disponíveis:");
                                    foreach (string fileName in answerList.FileNames)
                                    {
                                        Log.Information("- {FileName}", fileName);
                                    }
                                }
                                else
                                {
                                    Log.Information("Não foi possível obter a lista de arquivos");
                                }
                                break;
                            default:
                                break;
                        }
                    }
                    else
                    {
                        Log.Information("Comando inválido!");
                    }
                }

            }
            catch (EndOfStreamException)
            {

                Log.Error("O servidor foi desconectado.");
            }
            catch (IOException ex)
            {
                Log.Error(ex, "Erro de comunicação com o servidor.");
            }
        }
        catch (Exception ex)
        {
            Log.Error(ex, "Erro inesperado.");
        }
    }
}