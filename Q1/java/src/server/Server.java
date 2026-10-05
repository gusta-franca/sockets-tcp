package src.common.server;
import java.io.*;
import java.net.*;
import java.nio.file.Path;
import src.common.Config;
import src.common.Connection;

public class Server {

    public static void main(String[] args) {
        Config config = new Config(".env");
        int serverPort = config.getServerPort();
        Path usersRoot = Path.of(config.getBaseDir()).toAbsolutePath().normalize();

        try (ServerSocket listenSocket = new ServerSocket(serverPort);) {
            Authenticator authenticator = new Authenticator();
            System.out.println("Servidor aguardando conexão ...");

            while (true) {
                Socket clientSocket = listenSocket.accept();

                ClientThread clientThread = new ClientThread(clientSocket, authenticator, usersRoot);
                clientThread.start();
            }

        } catch (IOException e) {
            System.out.println("Listen socket:" + e.getMessage());
        }
    }
}

class ClientThread extends Thread {

    private final Connection connection;
    private final Session session;

    public ClientThread(Socket clientSocket, Authenticator authenticator, Path usersRoot) throws IOException {
        this.connection = new Connection(clientSocket);
        this.session = new Session(connection, authenticator, usersRoot);
    }

    @Override
    public void run() {
        try (connection) {
            while (true) {
                String message = connection.receive();
                String command = message.trim();

                boolean keepRunning = session.commandHandler(command);
                if (keepRunning == false) {
                    break;
                }
            }
        } catch (EOFException e) {
            System.out.println("EOF: " + e.getMessage());
        } catch (IOException e) {
            System.out.println("IOE: " + e.getMessage());
        }
        System.out.println("Thread comunicação cliente finalizada.");
    }
}