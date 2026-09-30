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

    public string Pwd() // retorna o caminho atual
    {
        return currentDir;
    }

    public bool Chdir(string path) // altera o caminho corrente para path
    {
        string newPath = Path.GetFullPath(Path.Combine(currentDir, path));
        string relativePath = Path.GetRelativePath(rootDir, newPath);
        // para analisar se o caminho está fora de root
        if (relativePath.StartsWith(".."))
        {
            return false;
        }

        if (!Directory.Exists(newPath))
        {
            return false;
        }

        currentDir = newPath;

        return true;
    }

    public string[] GetFiles() // retorna os arquivos do diretório corrente
    {
        string[] pathFileNames = Directory.GetFiles(currentDir);
        string[] fileNames = new string[pathFileNames.Length];
        for (int i = 0; i < pathFileNames.Length; i++)
        {
            fileNames[i] = Path.GetFileName(pathFileNames[i]);
        }
        return fileNames;
    }

    public string[] GetDirs() // retorna os diretórios do diretório corrente
    {
        string[] pathDirNames = Directory.GetDirectories(currentDir);
        string[] dirNames = new string[pathDirNames.Length];
        for (int i = 0; i < pathDirNames.Length; i++)
        {
            dirNames[i] = Path.GetFileName(pathDirNames[i]);
        }
        return dirNames;
    }
}