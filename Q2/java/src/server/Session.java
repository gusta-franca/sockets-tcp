/*
    Última atualização: 04/10/2026
    Aluna: Maria Eduarda Bambini
    Arquivo do servidor responsável por executar os comandos de um cliente sobre a pasta de arquivos.
*/

package server;

import common.Connection;
import common.Protocol;

import java.io.BufferedInputStream;
import java.io.BufferedOutputStream;
import java.io.IOException;
import java.io.InputStream;
import java.io.OutputStream;
import java.nio.charset.StandardCharsets;
import java.nio.file.Files;
import java.nio.file.InvalidPathException;
import java.nio.file.Path;
import java.nio.file.StandardCopyOption;
import java.util.ArrayList;
import java.util.List;
import java.util.logging.Level;
import java.util.logging.Logger;
import java.util.stream.Stream;

/**
 * Session: executa os comandos ADDFILE, DELETE, GETFILESLIST e GETFILE recebidos
 * de um cliente, respondendo pelo protocolo binário e registrando cada ação no log.
 */
class Session {

    private static final Logger LOGGER = Logger.getLogger(Session.class.getName());

    private final Connection connection;
    private final Path baseDir;
    private final String remote;

    public Session(Connection connection, Path baseDir, String remote) {
        this.connection = connection;
        this.baseDir = baseDir;
        this.remote = remote;
    }

    /** Executa a requisição; retorna false se a sessão deve ser encerrada. */
    public boolean commandHandler(Protocol.RequestHeader header) throws IOException {
        Protocol.Command command = Protocol.Command.fromCode(header.commandCode());

        if (header.messageType() != Protocol.MESSAGE_REQUEST || command == null) {
            LOGGER.log(Level.WARNING, "Requisição inválida de {0}: tipo={1}, comando={2}",
                    new Object[] {remote, header.messageType(), header.commandCode()});
            reply(header.commandCode(), false);
            return false;
        }

        boolean success = switch (command) {
            case ADDFILE -> addFile(header.filename());
            case DELETE -> delete(header.filename());
            case GETFILESLIST -> getFilesList();
            case GETFILE -> getFile(header.filename());
        };

        LOGGER.log(Level.INFO, "{0} executou {1} [{2}]: {3}",
                new Object[] {remote, command, header.filename(), success ? "SUCCESS" : "ERROR"});
        return true;
    }

    private boolean addFile(String filename) throws IOException {
        long size = connection.readUInt32();
        Path target = resolve(filename);
        if (target == null) {
            connection.discard(size);
            return reply(Protocol.Command.ADDFILE.code, false);
        }

        Path temp = Files.createTempFile("upload-", ".tmp");
        boolean stored;
        try {
            try (OutputStream out = new BufferedOutputStream(Files.newOutputStream(temp))) {
                connection.readFileTo(out, size);
            }
            stored = moveIntoPlace(temp, target);
        } finally {
            Files.deleteIfExists(temp);
        }
        return reply(Protocol.Command.ADDFILE.code, stored);
    }

    /** Move o arquivo temporário para o destino final, substituindo um existente. */
    private boolean moveIntoPlace(Path temp, Path target) {
        try {
            Files.move(temp, target, StandardCopyOption.REPLACE_EXISTING);
            return true;
        } catch (IOException e) {
            LOGGER.log(Level.WARNING, "Erro ao gravar " + target, e);
            return false;
        }
    }

    private boolean delete(String filename) throws IOException {
        Path target = resolve(filename);
        boolean deleted = false;
        if (target != null) {
            try {
                deleted = Files.deleteIfExists(target);
            } catch (IOException e) {
                LOGGER.log(Level.WARNING, "Erro ao remover " + target, e);
            }
        }
        return reply(Protocol.Command.DELETE.code, deleted);
    }

    private boolean getFilesList() throws IOException {
        List<byte[]> names = new ArrayList<>();
        try (Stream<Path> entries = Files.list(baseDir)) {
            entries.filter(Files::isRegularFile)
                    .map(p -> p.getFileName().toString().getBytes(StandardCharsets.UTF_8))
                    .filter(n -> n.length >= 1 && n.length <= Protocol.MAX_FILENAME_SIZE)
                    .limit(0xFFFF)
                    .forEach(names::add);
        } catch (IOException e) {
            LOGGER.log(Level.WARNING, "Erro ao listar " + baseDir, e);
            return reply(Protocol.Command.GETFILESLIST.code, false);
        }

        connection.writeAnswerHeader(Protocol.Command.GETFILESLIST.code, Protocol.STATUS_SUCCESS);
        connection.writeUInt16(names.size());
        for (byte[] name : names) {
            connection.writeByte(name.length);
            connection.writeBytes(name);
        }
        connection.flush();
        return true;
    }

    private boolean getFile(String filename) throws IOException {
        Path target = resolve(filename);
        if (target == null || !Files.isRegularFile(target)) {
            return reply(Protocol.Command.GETFILE.code, false);
        }

        long size;
        InputStream in;
        try {
            size = Files.size(target);
            if (size > Protocol.MAX_FILE_SIZE) {
                return reply(Protocol.Command.GETFILE.code, false);
            }
            in = new BufferedInputStream(Files.newInputStream(target));
        } catch (IOException e) {
            LOGGER.log(Level.WARNING, "Erro ao ler " + target, e);
            return reply(Protocol.Command.GETFILE.code, false);
        }

        try (InputStream source = in) {
            connection.writeAnswerHeader(Protocol.Command.GETFILE.code, Protocol.STATUS_SUCCESS);
            connection.writeUInt32(size);
            connection.writeFrom(source, size);
            connection.flush();
        }
        return true;
    }

    /** Envia uma resposta apenas com o cabeçalho e retorna o próprio status. */
    private boolean reply(int commandCode, boolean success) throws IOException {
        connection.writeAnswerHeader(commandCode, success ? Protocol.STATUS_SUCCESS : Protocol.STATUS_ERROR);
        connection.flush();
        return success;
    }

    /** Resolve o nome dentro da pasta base; retorna null se for inválido ou escapar dela. */
    private Path resolve(String filename) {
        if (filename.isEmpty()) {
            return null;
        }
        try {
            Path target = baseDir.resolve(filename).normalize();
            if (!baseDir.equals(target.getParent())) {
                LOGGER.log(Level.WARNING, "Nome de arquivo rejeitado de {0}: {1}", new Object[] {remote, filename});
                return null;
            }
            return target;
        } catch (InvalidPathException e) {
            return null;
        }
    }
}