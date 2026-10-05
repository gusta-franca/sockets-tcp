/*
    Última atualização: 04/10/2026
    Aluna: Maria Eduarda Bambini
    Arquivo do cliente responsável por ler os comandos do usuário, enviá-los ao servidor e exibir as respostas.
*/

package client;

import common.Config;
import common.Connection;
import common.Protocol;
import java.io.BufferedInputStream;
import java.io.BufferedOutputStream;
import java.io.EOFException;
import java.io.IOException;
import java.io.InputStream;
import java.io.OutputStream;
import java.net.Socket;
import java.net.UnknownHostException;
import java.nio.charset.StandardCharsets;
import java.nio.file.Files;
import java.nio.file.InvalidPathException;
import java.nio.file.Path;
import java.util.Scanner;

/**
 * Client: cliente interativo. Lê um comando por vez do teclado, envia ao servidor
 * pelo protocolo binário e mostra a resposta. Os downloads são gravados em DOWNLOAD_DIR.
 */
public class Client {

    public static void main(String[] args) {

        Config config = new Config(".env");

        try {
            String serverHost = config.getServerHost();
            int serverPort = config.getServerPort();
            Path downloadDir = Path.of(config.getDownloadDir()).toAbsolutePath().normalize();
            run(serverHost, serverPort, downloadDir);
        } catch (UnknownHostException e) {
            System.out.println("Socket: " + e.getMessage());
        } catch (EOFException e) {
            System.out.println("O servidor foi desconectado.");
        } catch (IOException e) {
            System.out.println("IO: " + e.getMessage());
        } catch (IllegalArgumentException e) {
            System.out.println("Configuração inválida: " + e.getMessage());
        }
    }

    /** Conecta ao servidor e processa os comandos digitados até o usuário sair. */
    private static void run(String serverHost, int serverPort, Path downloadDir) throws IOException {
        try (Scanner reader = new Scanner(System.in);
             Socket socket = new Socket(serverHost, serverPort);
             Connection connection = new Connection(socket)) {

            Files.createDirectories(downloadDir);
            System.out.println("Conectado ao servidor");

            while (true) {
                System.out.print("> ");
                if (!reader.hasNextLine()) {
                    break;
                }
                String buffer = reader.nextLine().trim();

                if (buffer.isEmpty()) {
                    continue;
                }
                if (buffer.equalsIgnoreCase("EXIT")) {
                    break;
                }

                Protocol.Command command = parseCommand(buffer);
                if (command == null) {
                    System.out.println("Comando inválido!");
                    continue;
                }

                switch (command) {
                    case ADDFILE -> addFile(connection, readFilename(reader));
                    case DELETE -> delete(connection, readFilename(reader));
                    case GETFILESLIST -> getFilesList(connection);
                    case GETFILE -> getFile(connection, downloadDir, readFilename(reader));
                }
            }
        }
    }

    private static Protocol.Command parseCommand(String input) {
        try {
            return Protocol.Command.valueOf(input.toUpperCase());
        } catch (IllegalArgumentException e) {
            return null;
        }
    }

    /** Solicita o nome do arquivo até que tenha de 1 a 255 bytes em UTF-8. */
    private static String readFilename(Scanner reader) {
        while (true) {
            System.out.print("Digite o nome do arquivo: ");
            String name = reader.nextLine().trim();
            int size = name.getBytes(StandardCharsets.UTF_8).length;
            if (size >= 1 && size <= Protocol.MAX_FILENAME_SIZE) {
                return name;
            }
            System.out.println("O nome do arquivo deve ter entre 1 e 255 bytes");
        }
    }

    /** Envia o arquivo indicado pelo caminho; ao servidor vai apenas o nome do arquivo. */
    private static void addFile(Connection connection, String filePath) throws IOException {
        Path path = Path.of(filePath);
        if (!Files.isRegularFile(path)) {
            System.out.println("O arquivo informado não existe e não pode ser enviado ao servidor");
            return;
        }
        long size = Files.size(path);
        if (size > Protocol.MAX_FILE_SIZE) {
            System.out.println("O arquivo ultrapassa o limite de 4 GB");
            return;
        }

        connection.writeRequestHeader(Protocol.Command.ADDFILE, path.getFileName().toString());
        connection.writeUInt32(size);
        try (InputStream in = new BufferedInputStream(Files.newInputStream(path))) {
            connection.writeFrom(in, size);
        }
        connection.flush();

        boolean ok = readAnswer(connection).statusCode() == Protocol.STATUS_SUCCESS;
        System.out.println(ok ? "Arquivo enviado com sucesso!" : "Não foi possível enviar o arquivo");
    }

    private static void delete(Connection connection, String filename) throws IOException {
        connection.writeRequestHeader(Protocol.Command.DELETE, filename);
        connection.flush();

        boolean ok = readAnswer(connection).statusCode() == Protocol.STATUS_SUCCESS;
        System.out.println(ok ? "Arquivo deletado com sucesso!" : "Não foi possível deletar o arquivo");
    }

    private static void getFilesList(Connection connection) throws IOException {
        connection.writeRequestHeader(Protocol.Command.GETFILESLIST, "");
        connection.flush();

        if (readAnswer(connection).statusCode() != Protocol.STATUS_SUCCESS) {
            System.out.println("Não foi possível obter a lista de arquivos");
            return;
        }

        int count = connection.readUInt16();
        System.out.println("Quantidade: " + count);
        for (int i = 0; i < count; i++) {
            int nameSize = connection.readUnsignedByte();
            System.out.println(" - " + connection.readString(nameSize));
        }
    }

    private static void getFile(Connection connection, Path downloadDir, String filename) throws IOException {
        connection.writeRequestHeader(Protocol.Command.GETFILE, filename);
        connection.flush();

        if (readAnswer(connection).statusCode() != Protocol.STATUS_SUCCESS) {
            System.out.println("Não foi possível baixar o arquivo");
            return;
        }

        long size = connection.readUInt32();
        Path target = localTarget(downloadDir, filename);
        if (target == null) {
            connection.discard(size);
            System.out.println("Nome de arquivo inválido para gravação local");
            return;
        }

        try (OutputStream out = new BufferedOutputStream(Files.newOutputStream(target))) {
            connection.readFileTo(out, size);
        }
        System.out.println("Arquivo " + filename + " baixado com sucesso!");
    }

    /** Define o destino do download dentro da pasta padrão, ignorando qualquer diretório no nome. */
    private static Path localTarget(Path downloadDir, String filename) {
        try {
            Path name = Path.of(filename).getFileName();
            return name == null ? null : downloadDir.resolve(name);
        } catch (InvalidPathException e) {
            return null;
        }
    }

    /** Lê o cabeçalho da resposta e confirma que é uma mensagem do tipo resposta. */
    private static Protocol.AnswerHeader readAnswer(Connection connection) throws IOException {
        Protocol.AnswerHeader header = connection.readAnswerHeader();
        if (header.messageType() != Protocol.MESSAGE_ANSWER) {
            throw new IOException("Tipo de mensagem inesperado: " + header.messageType());
        }
        return header;
    }
}