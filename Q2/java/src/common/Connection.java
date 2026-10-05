/*
    Última atualização: 04/10/2026
    Aluna: Maria Eduarda Bambini
    Arquivo compartilhado (cliente e servidor) responsável por enviar e receber o protocolo binário pelo socket.
*/

package common;

import java.io.BufferedInputStream;
import java.io.BufferedOutputStream;
import java.io.Closeable;
import java.io.EOFException;
import java.io.IOException;
import java.io.InputStream;
import java.io.OutputStream;
import java.net.Socket;
import java.nio.charset.StandardCharsets;

/**
 * Connection: encapsula um Socket TCP e os streams de leitura/escrita em bytes.
 * Todos os números são enviados em big endian e os arquivos trafegam byte a byte.
 */
public class Connection implements Closeable {

    private final Socket socket;
    private final InputStream in;
    private final OutputStream out;

    public Connection(Socket socket) throws IOException {
        this.socket = socket;
        this.in = new BufferedInputStream(socket.getInputStream());
        this.out = new BufferedOutputStream(socket.getOutputStream());
    }

    /** Lê um byte sem sinal (0-255); lança EOFException se a conexão fechar. */
    public int readUnsignedByte() throws IOException {
        int value = in.read();
        if (value < 0) {
            throw new EOFException("Conexão encerrada");
        }
        return value;
    }

    /** Lê um inteiro sem sinal de 2 bytes. */
    public int readUInt16() throws IOException {
        return (readUnsignedByte() << 8) | readUnsignedByte();
    }

    /** Lê um inteiro sem sinal de 4 bytes. */
    public long readUInt32() throws IOException {
        long value = 0;
        for (int i = 0; i < 4; i++) {
            value = (value << 8) | readUnsignedByte();
        }
        return value;
    }

    /** Lê {@code size} bytes e os interpreta como texto UTF-8. */
    public String readString(int size) throws IOException {
        byte[] data = new byte[size];
        for (int i = 0; i < size; i++) {
            data[i] = (byte) readUnsignedByte();
        }
        return new String(data, StandardCharsets.UTF_8);
    }

    /** Recebe {@code size} bytes, um a um, gravando-os em {@code dest}. */
    public void readFileTo(OutputStream dest, long size) throws IOException {
        for (long i = 0; i < size; i++) {
            dest.write(readUnsignedByte());
        }
    }

    /** Descarta {@code size} bytes do fluxo de entrada. */
    public void discard(long size) throws IOException {
        for (long i = 0; i < size; i++) {
            readUnsignedByte();
        }
    }

    /** Lê o cabeçalho comum das requisições. */
    public Protocol.RequestHeader readRequestHeader() throws IOException {
        int messageType = readUnsignedByte();
        int commandCode = readUnsignedByte();
        int filenameSize = readUnsignedByte();
        return new Protocol.RequestHeader(messageType, commandCode, readString(filenameSize));
    }

    /** Lê o cabeçalho comum das respostas. */
    public Protocol.AnswerHeader readAnswerHeader() throws IOException {
        return new Protocol.AnswerHeader(readUnsignedByte(), readUnsignedByte(), readUnsignedByte());
    }

    public void writeByte(int value) throws IOException {
        out.write(value);
    }

    /** Envia um inteiro sem sinal de 2 bytes. */
    public void writeUInt16(int value) throws IOException {
        out.write((value >> 8) & 0xFF);
        out.write(value & 0xFF);
    }

    /** Envia um inteiro sem sinal de 4 bytes. */
    public void writeUInt32(long value) throws IOException {
        for (int shift = 24; shift >= 0; shift -= 8) {
            out.write((int) ((value >> shift) & 0xFF));
        }
    }

    /** Envia os bytes de {@code data}, um a um. */
    public void writeBytes(byte[] data) throws IOException {
        for (byte b : data) {
            out.write(b);
        }
    }

    /** Envia {@code size} bytes lidos de {@code src}, um a um. */
    public void writeFrom(InputStream src, long size) throws IOException {
        for (long i = 0; i < size; i++) {
            int value = src.read();
            if (value < 0) {
                throw new EOFException("Arquivo menor do que o tamanho informado");
            }
            out.write(value);
        }
    }

    /** Envia o cabeçalho comum das requisições. */
    public void writeRequestHeader(Protocol.Command command, String filename) throws IOException {
        byte[] name = filename.getBytes(StandardCharsets.UTF_8);
        if (name.length > Protocol.MAX_FILENAME_SIZE) {
            throw new IllegalArgumentException("O nome do arquivo ultrapassa o limite de 255 bytes");
        }
        writeByte(Protocol.MESSAGE_REQUEST);
        writeByte(command.code);
        writeByte(name.length);
        writeBytes(name);
    }

    /** Envia o cabeçalho comum das respostas. */
    public void writeAnswerHeader(int commandCode, int statusCode) throws IOException {
        writeByte(Protocol.MESSAGE_ANSWER);
        writeByte(commandCode);
        writeByte(statusCode);
    }

    /** Força o envio dos bytes pendentes no buffer. */
    public void flush() throws IOException {
        out.flush();
    }

    /** Fecha streams e socket. */
    @Override
    public void close() {
        try { in.close(); }  catch (IOException e) { System.err.println("IOE: " + e); }
        try { out.close(); } catch (IOException e) { System.err.println("IOE: " + e); }
        try { socket.close(); } catch (IOException e) { System.err.println("IOE: " + e); }
    }
}