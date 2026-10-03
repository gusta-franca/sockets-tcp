import java.io.File;
import java.io.IOException;
import java.nio.file.Files;
import java.nio.file.Path;
import java.util.ArrayList;
import java.util.List;

public class Session {
    private final Connection connection;
    private final Authenticator authenticator;
    
    private boolean auth = false;
    private String user;
    private Path baseDir;
    private Path currDir;

    public Session(Connection connection, Authenticator authenticator) {
        this.connection = connection;
        this.authenticator = authenticator;
    }

    public boolean commandHandler(String command) throws IOException {
        String[] tokens = command.split(" ", 2);
        String cmd = tokens[0].toUpperCase();

        String args = "";
        if (tokens.length == 2) {
            args = tokens[1].trim();
        }

        if (!auth && !cmd.equals("CONNECT") && !cmd.equals("EXIT")) {
            connection.send("ERROR");
            System.out.println("Servidor: É necessário estabelecer conexão antes de utilizar um comando.");
            return true;
        }

        switch (cmd) {
            case "CONNECT": 
                connect(args);
                return true;
            case "PWD":
                printWorkingDirectory();
                return true;
            case "CHDIR":
                changeDirectory(args);
                return true;
            case "GETFILES":
                sendItems(true);
                return true;
            case "GETDIRS":
                sendItems(false);
                return true;
            case "EXIT":
                return false;
            default:
                connection.send("ERROR");
                System.out.println("Servidor: Comando não encontrado.");
                return true;
        }
    }

    private void connect(String args) throws IOException {
        String[] values = args.split(" ");

        if (values.length != 2) {
            connection.send("ERROR");
            System.out.println("Servidor: Utilize o formato 'CONNECT <usuario> <senha>'");
            return;
        }

        String userNameWithComma = values[0];
        String passwordHash = values[1];

        if (userNameWithComma.endsWith(",") == false) {
            connection.send("ERROR");
            System.out.println("Servidor: Utilize vírgula após o nome de usuário.");
            return;
        }

        int userLastIndex = userNameWithComma.length() - 1;
        String userName = userNameWithComma.substring(0, userLastIndex);

        if (!authenticator.authenticate(userName, passwordHash)) {
            connection.send("ERROR");
            System.out.println("Servidor: Usuário ou senha inválidos.");
            return;
        }

        baseDir = Path.of("/tmp", userName);
        Files.createDirectories(baseDir);

        this.user = userName;
        this.currDir = baseDir;
        this.auth = true;
        connection.send("SUCCESS");
    }

    private void printWorkingDirectory() throws IOException {
        String path = currDir.toString();
        connection.send(path);
    }

    private void changeDirectory(String args) throws IOException {

        if (args.isEmpty()) {
            connection.send("ERROR");
            System.out.println("Servidor: Utilize o formato 'CHDIR <diretório>'");
            return;
        }

        Path newCurrDir = currDir.resolve(args); 

        if (newCurrDir.startsWith(baseDir) == false) {
            connection.send("ERROR");
            System.out.println("Diretórios válidos começam na raíz 'tmp'.");
            return;
        }

        if (Files.isDirectory(newCurrDir) == false) {
            connection.send("ERROR: diretório inválido.");
            return;
        }

        currDir = newCurrDir;
        connection.send("SUCCESS");
    }

    private void sendItems(boolean wantsFiles) throws IOException {
        List<String> filesNames = new ArrayList<>();
        List<String> directoriesNames = new ArrayList<>();

        File[] items = currDir.toFile().listFiles();

        if (items != null) {
            for (int i = 0; i < items.length; i++) {
                File item = items[i];

                if (item.isFile()) {
                    filesNames.add(item.getName());
                }

                if (item.isDirectory()) {
                    directoriesNames.add(item.getName());
                }
            }
        }

        List<String> namesToSend;
        if (wantsFiles) {
            namesToSend = filesNames;
        } else {
            namesToSend = directoriesNames;
        }

        connection.send(String.valueOf(namesToSend.size()));
        for (int i = 0; i < namesToSend.size(); i++) {
            connection.send(namesToSend.get(i));
        }
    }

}
