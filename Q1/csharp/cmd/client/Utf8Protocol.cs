/*
    Data de criação: 30/09/2026
    Aluna: Caroline Marques Lau
    Arquivo do protocolo criado
*/

using System.Net.Sockets;
using System.Text;
using System.Buffers.Binary;

class Utf8Protocol(NetworkStream stream)
{
    private readonly NetworkStream stream = stream;

    /*
        Os primeiros quatro bytes sempre representam o tamanho da mensagem que vem em seguida
        O tamanho é enviado utilizando o padrão Big Endian
    */
    public string ReadString()
    {
        byte[] lengthBytes = new byte[4];
        stream.ReadExactly(lengthBytes);
        int length = BinaryPrimitives.ReadInt32BigEndian(lengthBytes);
        byte[] data = new byte[length];
        stream.ReadExactly(data);
        return Encoding.UTF8.GetString(data);
    }

    public void WriteString(string message)
    {
        byte[] data = Encoding.UTF8.GetBytes(message);
        int length = data.Length;
        byte[] lengthBytes = new byte[4];
        BinaryPrimitives.WriteInt32BigEndian(lengthBytes, length);
        stream.Write(lengthBytes);
        stream.Write(data);
    }
}