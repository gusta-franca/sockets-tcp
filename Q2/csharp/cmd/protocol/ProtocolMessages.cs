public enum Command : byte
{
    ADDFILE = 0x01,
    DELETE = 0x02,
    GETFILESLIST = 0x03,
    GETFILE = 0x04
}

public enum StatusCode : byte
{
    SUCCESS = 0x01,
    ERROR = 0x02
}

public struct Request
{
    public MessageType MessageType;
    public Command CommandIdentifier;
    public byte FileNameSize;
    public byte[] FileName;

    public uint? FileSize;
    public byte[]? File;
}

public enum MessageType : byte
{
    Request = 0x01,
    Answer = 0x02
}

public struct Answer
{
    public MessageType MessageType;
    public Command CommandIdentifier;
    public StatusCode StatusCode;
    public ushort? FileAmount;
    public string[]? FileNames;
    public uint? FileSize;
    public byte[]? File;
}