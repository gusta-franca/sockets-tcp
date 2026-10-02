/*
    Data de criação: 30/09/2026
    Aluna: Caroline Marques Lau
    Arquivo do protocolo para troca de mensagens
*/

using System.Buffers.Binary;
using System.Net.Sockets;
using System.Text;

public class Protocol(NetworkStream stream)
{
    private readonly NetworkStream stream = stream;
    public void SendRequest(Request request)
    {
        if (request.MessageType != MessageType.Request)
        {
            throw new InvalidDataException("Tipo de mensagem inválido!");
        }
        byte[] fileName = request.FileName;
        // o nome não pode passar de 255 bytes
        if (fileName.Length > 255) throw new InvalidDataException("O nome do arquivo não pode passar de 255 bytes");
        // é o comando AddFile?
        bool hasFile = request.CommandIdentifier == Command.ADDFILE;
        // se for o campo File não pode ser nulo
        if (hasFile && request.File == null) throw new InvalidDataException("O campo para o arquivo veio nulo"); ;
        // armazenando o tamanho do File para ver se bate com o do Request
        uint fileSize = hasFile ? (uint)request.File!.Length : 0;

        if (hasFile && request.FileSize != fileSize)
            throw new InvalidDataException("O tamanho real do arquivo e o tamanho enviado não batem");
        // tamanho padrão do Request
        int packetSize = 3 + fileName.Length;
        // se for AddFile o Request é maior
        if (hasFile)
            packetSize += 4 + request.File!.Length;
        // Montando o pacote a ser enviado
        byte[] packet = new byte[packetSize];
        int offset = 0;
        packet[offset++] = (byte)request.MessageType;
        packet[offset++] = (byte)request.CommandIdentifier;
        packet[offset++] = (byte)fileName.Length;
        Array.Copy(fileName, 0, packet, 3, fileName.Length);
        offset += fileName.Length;
        if (hasFile)
        {
            BinaryPrimitives.WriteUInt32BigEndian(packet.AsSpan(offset, 4), fileSize); // pega 4 bytes começando de offset
            offset += 4;
            Array.Copy(request.File!, 0, packet, offset, request.File!.Length);
        }
        try
        {
            stream.Write(packet);
        }
        catch (IOException ex)
        {
            throw new IOException("Erro ao enviar o pacote pelo stream.", ex);
        }
    }

    public Request ReadRequest()
    {
        byte messageType = ReadByte();

        if (messageType != (byte)MessageType.Request)
            throw new InvalidDataException("Tipo de mensagem inválido!");

        byte commandByte = ReadByte();

        if (!Enum.IsDefined(typeof(Command), commandByte))
            throw new InvalidDataException("Comando inválido!");

        Command command = (Command)commandByte;

        byte fileNameSize = ReadByte();

        byte[] fileName = new byte[fileNameSize];
        stream.ReadExactly(fileName);

        Request request = new()
        {
            MessageType = (MessageType)messageType,
            CommandIdentifier = command,
            FileName = fileName
        };

        if (command == Command.ADDFILE)
        {
            byte[] fileSize = new byte[4];
            stream.ReadExactly(fileSize);
            uint size = BinaryPrimitives.ReadUInt32BigEndian(fileSize);

            const uint MaxFileSize = 100 * 1024 * 1024; // sugestão da IA

            if (size > MaxFileSize)
                throw new InvalidDataException("Arquivo muito grande!");

            request.FileSize = size;
            byte[] file = new byte[size];
            stream.ReadExactly(file);
            request.File = file;
        }

        return request;
    }

    private byte ReadByte()
    {
        byte[] buffer = new byte[1];
        stream.ReadExactly(buffer);
        return buffer[0];
    }

    public Answer ReadAnswer()
    {
        byte messageType = ReadByte();

        if (messageType != (byte)MessageType.Answer)
            throw new InvalidDataException("Tipo de mensagem inválido!");

        byte commandByte = ReadByte();

        if (!Enum.IsDefined(typeof(Command), commandByte))
            throw new InvalidDataException("Comando inválido!");

        Command command = (Command)commandByte;
        byte statusCode = ReadByte();

        if (!Enum.IsDefined(typeof(StatusCode), statusCode))
            throw new InvalidDataException("Status inválido!");

        Answer answer = new()
        {
            MessageType = (MessageType)messageType,
            CommandIdentifier = command,
            StatusCode = (StatusCode)statusCode
        };

        if (command == Command.GETFILESLIST)
        {
            byte[] fileAmount = new byte[2];
            stream.ReadExactly(fileAmount);
            ushort amount = BinaryPrimitives.ReadUInt16BigEndian(fileAmount);

            if (amount <= 0)
            {
                throw new InvalidDataException("A quantidade no comando GETFILESLIST precisa ser maior que 0");
            }
            string[] fileNames = new string[amount];
            for (int i = 0; i < amount; i++)
            {
                byte fileNameSize = ReadByte();
                byte[] fileName = new byte[fileNameSize];
                stream.ReadExactly(fileName);
                fileNames[i] = Encoding.UTF8.GetString(fileName);
            }
            answer.FileAmount = amount;
            answer.FileNames = fileNames;
        }
        else if (command == Command.GETFILE)
        {
            byte[] fileSize = new byte[4];
            stream.ReadExactly(fileSize);
            uint size = BinaryPrimitives.ReadUInt32BigEndian(fileSize);
            byte[] file = new byte[size];
            stream.ReadExactly(file);
            answer.FileSize = size;
            answer.File = file;
        }

        return answer;
    }

    public void SendAnswer(Answer answer)
    {
        int packetSize = 3;

        if (answer.MessageType != MessageType.Answer)
            throw new InvalidDataException("Tipo de mensagem inválido!");


        if (answer.CommandIdentifier == Command.GETFILESLIST)
        {
            if (answer.FileNames == null)
                throw new InvalidDataException("A lista de arquivos não pode ser nula");

            if (answer.FileNames.Length > ushort.MaxValue)
                throw new InvalidDataException("A quantidade de arquivos não pode passar de 65535");

            foreach (string fileName in answer.FileNames)
            {
                byte[] fileNameBytes = Encoding.UTF8.GetBytes(fileName);

                if (fileNameBytes.Length > 255)
                    throw new InvalidDataException(
                        "O nome do arquivo não pode passar de 255 bytes");

                packetSize += 1 + fileNameBytes.Length; // tamanho do nome mais o nome
            }

            packetSize += 2; // número de arquivos em 2 bytes
        }
        else if (answer.CommandIdentifier == Command.GETFILE)
        {
            if (answer.File == null)
                throw new InvalidDataException("O arquivo não pode ser nulo");

            packetSize += 4 + answer.File.Length; // tamaho do arquivo + arquivo
        }
        byte[] packet = new byte[packetSize];
        int offset = 0;
        packet[offset++] = (byte)answer.MessageType;
        packet[offset++] = (byte)answer.CommandIdentifier;
        packet[offset++] = (byte)answer.StatusCode;

        if (answer.CommandIdentifier == Command.GETFILESLIST)
        {
            ushort amount = (ushort)answer.FileNames!.Length;

            BinaryPrimitives.WriteUInt16BigEndian(
                packet.AsSpan(offset, 2),
                amount);

            offset += 2; // número de arquivos

            foreach (string fileName in answer.FileNames)
            {
                byte[] fileNameBytes = Encoding.UTF8.GetBytes(fileName);

                packet[offset++] = (byte)fileNameBytes.Length;

                Array.Copy(
                    fileNameBytes,
                    0,
                    packet,
                    offset,
                    fileNameBytes.Length);

                offset += fileNameBytes.Length;
            }
        }
        else if (answer.CommandIdentifier == Command.GETFILE)
        {
            uint fileSize = (uint)answer.File!.Length;

            BinaryPrimitives.WriteUInt32BigEndian(
                packet.AsSpan(offset, 4),
                fileSize);

            offset += 4;

            Array.Copy(
                answer.File,
                0,
                packet,
                offset,
                answer.File.Length);
        }
        try
        {
            stream.Write(packet);
        }
        catch (IOException ex)
        {
            throw new IOException("Erro ao enviar a resposta.", ex);
        }
    }
}