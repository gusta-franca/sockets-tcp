/*
    Última atualização: 04/10/2026
    Aluna: Maria Eduarda Bambini
    Arquivo compartilhado (cliente e servidor) responsável por enviar e receber mensagens pelo socket.
*/

package common;
 
import java.io.*;
import java.net.Socket;
import java.nio.charset.StandardCharsets;

/**
 * TCPConnection: encapsula um Socket TCP e os streams de leitura/escrita UTF.
 * Usada tanto pelo servidor quanto pelo cliente.
 */
public class Connection implements Closeable {

    private Socket socket;
    private DataInputStream in;
    private DataOutputStream out;

    public Connection(Socket socket) throws IOException {
        this.socket = socket;
        this.in = new DataInputStream(socket.getInputStream());
        this.out = new DataOutputStream(socket.getOutputStream());
    }

    /** Envia uma mensagem UTF. */
    public void send(String msg) throws IOException {
        byte[] data = msg.getBytes(StandardCharsets.UTF_8);
        out.writeInt(data.length);
        out.write(data);
        out.flush();
    }

    /** Aguarda e retorna uma mensagem UTF. */
    public String receive() throws IOException {
        int length = in.readInt();

        if (length < 0 || length > 10000) {
            throw new IOException("Mensagem muito longa: " + length);
        }

        byte[] data = new byte[length];
        in.readFully(data);

        return new String(data, StandardCharsets.UTF_8);
    }

    /** Fecha streams e socket. */
    @Override
    public void close() {
        try { in.close(); }  catch (IOException e) { System.err.println("IOE: " + e); }
        try { out.close(); } catch (IOException e) { System.err.println("IOE: " + e); }
        try { socket.close(); } catch (IOException e) { System.err.println("IOE: " + e); }
    }
}