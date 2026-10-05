/*
    Última atualização: 04/10/2026
    Aluna: Maria Eduarda Bambini
    Arquivo do servidor responsável por aceitar as conexões e criar uma thread para cada cliente.
*/

package server;

import common.Config;
import common.Connection;
import common.Protocol;

import java.io.EOFException;
import java.io.IOException;
import java.net.ServerSocket;
import java.net.Socket;
import java.nio.file.Files;
import java.nio.file.Path;
import java.util.logging.ConsoleHandler;
import java.util.logging.FileHandler;
import java.util.logging.Handler;
import java.util.logging.Level;
import java.util.logging.Logger;
import java.util.logging.SimpleFormatter;

/**
 * Server: escuta na porta configurada no .env e, a cada cliente que conecta,
 * cria uma ClientThread para atendê-lo sem bloquear os demais.
 */
public class Server {

    static {
        System.setProperty("java.util.logging.SimpleFormatter.format", "%1$tF %1$tT [%4$s] %5$s%6$s%n");
    }

    private static final Logger LOGGER = Logger.getLogger(Server.class.getName());

    public static void main(String[] args) {
        Config config = new Config(".env");
        int serverPort = config.getServerPort();
        Path baseDir = Path.of(config.getBaseDir()).toAbsolutePath().normalize();

        try (ServerSocket listenSocket = new ServerSocket(serverPort)) {
            configureLogging();
            Files.createDirectories(baseDir);
            LOGGER.log(Level.INFO, "Servidor aguardando conexão na porta {0}, pasta {1}",
                    new Object[] {String.valueOf(serverPort), baseDir});

            while (true) {
                Socket clientSocket = listenSocket.accept();

                ClientThread clientThread = new ClientThread(clientSocket, baseDir);
                clientThread.start();
            }

        } catch (IOException e) {
            LOGGER.log(Level.SEVERE, "Listen socket: " + e.getMessage(), e);
        }
    }

    /** Envia os logs para o console e para o arquivo server.log. */
    private static void configureLogging() throws IOException {
        Logger root = Logger.getLogger("");
        for (Handler handler : root.getHandlers()) {
            root.removeHandler(handler);
        }
        root.setLevel(Level.INFO);

        Handler console = new ConsoleHandler();
        console.setFormatter(new SimpleFormatter());
        root.addHandler(console);

        Handler file = new FileHandler("server.log", true);
        file.setFormatter(new SimpleFormatter());
        root.addHandler(file);
    }
}

/**
 * ClientThread: thread responsável pela comunicação com um cliente. Lê o cabeçalho
 * de cada requisição, entrega para a Session e encerra quando a Session pede
 * ou quando o cliente desconecta.
 */
class ClientThread extends Thread {

    private static final Logger LOGGER = Logger.getLogger(ClientThread.class.getName());

    private final Connection connection;
    private final Session session;
    private final String remote;

    public ClientThread(Socket clientSocket, Path baseDir) throws IOException {
        this.remote = String.valueOf(clientSocket.getRemoteSocketAddress());
        this.connection = new Connection(clientSocket);
        this.session = new Session(connection, baseDir, remote);
    }

    @Override
    public void run() {
        try (connection) {
            LOGGER.log(Level.INFO, "Cliente conectado: {0}", remote);
            while (true) {
                Protocol.RequestHeader request = connection.readRequestHeader();

                boolean keepRunning = session.commandHandler(request);
                if (keepRunning == false) {
                    break;
                }
            }
        } catch (EOFException e) {
            LOGGER.log(Level.INFO, "Cliente desconectado: {0}", remote);
        } catch (IOException e) {
            LOGGER.log(Level.WARNING, "Erro de comunicação com " + remote, e);
        }
        LOGGER.log(Level.INFO, "Thread de comunicação finalizada: {0}", remote);
    }
}