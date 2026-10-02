/*
    Data de criação: 30/09/2026
    Aluna: Caroline Marques Lau
    Arquivo responsável por centralizar a lógica do gerenciador de arquivos
*/

class FileSystemManager
{
    private readonly string rootDir;
    private string currentDir;

    public FileSystemManager()
    {
        rootDir = Environment.GetEnvironmentVariable("SERVER_ROOT")
            ?? throw new Exception("SERVER_ROOT não foi definida.");

        currentDir = rootDir;
        if (!Directory.Exists(rootDir))
        {
            throw new DirectoryNotFoundException("O caminho base do servidor não existe.");
        }
    }

    public void AddFile(string filename, byte[]? file, uint? fileSize) // cria um arquivo no diretório atual
    {
        if (file == null)
            throw new InvalidDataException("O arquivo não pode ser nulo.");

        if (file.Length != fileSize)
            throw new InvalidDataException(
                "O tamanho real do arquivo não corresponde ao tamanho informado.");

        string path = Path.Combine(currentDir, filename);

        if (File.Exists(path))
            throw new IOException("O arquivo já existe.");

        try
        {
            File.WriteAllBytes(path, file);
        }
        catch (UnauthorizedAccessException ex)
        {
            throw new IOException("Sem permissão para criar o arquivo.", ex);
        }
        catch (IOException ex)
        {
            throw new IOException("Erro ao criar o arquivo.", ex);
        }
    }

    public void DeleteFile(string filename) // deleta um arquivo existente
    {

        string path = Path.Combine(currentDir, filename);
        if (!File.Exists(path))
            throw new FileNotFoundException(
                "O arquivo não existe.",
                filename);
        try
        {
            File.Delete(path);
        }
        catch (UnauthorizedAccessException ex)
        {
            throw new IOException(
                "Sem permissão para deletar o arquivo.",
                ex);
        }
        catch (IOException ex)
        {
            throw new IOException(
                "Erro ao deletar o arquivo.",
                ex);
        }
    }

    public string[] GetFilesList() // lista os arquivos do diretório
    {
        string[] pathFileNames = Directory.GetFiles(currentDir);
        string[] fileNames = new string[pathFileNames.Length];
        for (int i = 0; i < pathFileNames.Length; i++)
        {
            fileNames[i] = Path.GetFileName(pathFileNames[i]);
        }
        return fileNames;
    }

    public byte[]? GetFile(string filename) // envia os bytes do arquivo para permitir o download
    {
        string path = Path.Combine(currentDir, filename);
        if (!File.Exists(path)) return null;

        byte[] fileBytes = File.ReadAllBytes(path);
        return fileBytes;
    }
}