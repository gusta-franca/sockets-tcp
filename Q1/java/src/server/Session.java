/*
    Última atualização: 04/10/2026
    Aluna: Maria Eduarda Bambini
    Arquivo do servidor responsável por processar as mensagens dos clientes.
*/

package server;
 
import common.Connection;
import java.io.File;
import java.io.IOException;
import java.nio.file.Files;
import java.nio.file.Path;
import java.util.ArrayList;
import java.util.List;

/**
 * Session: representa um cliente conectado. Guarda o estado dele (se está autenticado
 * e em qual diretório está) e executa cada comando recebido. Cada comando gera
 * exatamente uma resposta, exceto o EXIT.
 */
public class Session {
    private final Connection connection;
    private final Authenticator authenticator;
    
    private boolean auth = false;
    private String user;
    private Path baseDir;
    private Path currDir;
    private final Path usersRoot;

    public Session(Connection connection, Authenticator authenticator, Path usersRoot) {
        this.connection = connection;
        this.authenticator = authenticator;
        this.usersRoot = usersRoot;
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
                return true;
        }
    }

    private void connect(String args) throws IOException {
        int commaIndex = args.indexOf(",");

        if (commaIndex == -1) {
            connection.send("ERROR");
            return;
        }

        String userName = args.substring(0, commaIndex).trim();
        String passwordHash = args.substring(commaIndex + 1).trim();

        if (userName.isEmpty() || passwordHash.isEmpty()) {
            connection.send("ERROR");
            return;
        }

        boolean validLogin = authenticator.authenticate(userName, passwordHash);
        if (validLogin == false) {
            connection.send("ERROR");
            return;
        }

        baseDir = usersRoot.resolve(userName);
        Files.createDirectories(baseDir);

        this.user = userName;
        this.currDir = baseDir;
        this.auth = true;
        connection.send("SUCCESS");
    }

    private void printWorkingDirectory() throws IOException {
        Path relative = baseDir.relativize(currDir);
        String path = "/" + relative.toString();
        connection.send(path);
    }

    private void changeDirectory(String args) throws IOException {

        if (args.isEmpty()) {
            connection.send("ERROR");
            return;
        }

        Path newCurrDir = currDir.resolve(args); 

        if (args.startsWith("/")) {
            String withoutSlash = args.substring(1);
            newCurrDir = baseDir.resolve(withoutSlash).normalize();
        } else {
            newCurrDir = currDir.resolve(args).normalize();
        }

        if (newCurrDir.startsWith(baseDir) == false) {
            connection.send("ERROR");
            return;
        }

        if (Files.isDirectory(newCurrDir) == false) {
            connection.send("ERROR");
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
