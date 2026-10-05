/*
    Última atualização: 04/10/2026
    Aluna: Maria Eduarda Bambini
    Arquivo compartilhado (cliente e servidor) com as constantes e estruturas do protocolo binário.
*/

package common;

/**
 * Protocol: códigos de mensagem, comandos e status, além dos cabeçalhos
 * de requisição e resposta definidos para a comunicação entre as linguagens.
 */
public final class Protocol {

    /** Comandos aceitos pelo servidor e seus códigos no protocolo. */
    public enum Command {
        ADDFILE(0x01),
        DELETE(0x02),
        GETFILESLIST(0x03),
        GETFILE(0x04);

        public final int code;

        Command(int code) {
            this.code = code;
        }

        /** Retorna o comando correspondente ao código ou null se desconhecido. */
        public static Command fromCode(int code) {
            for (Command command : values()) {
                if (command.code == code) {
                    return command;
                }
            }
            return null;
        }
    }

    public record RequestHeader(int messageType, int commandCode, String filename) {}

    public record AnswerHeader(int messageType, int commandCode, int statusCode) {}

    public static final int MESSAGE_REQUEST = 0x01;
    public static final int MESSAGE_ANSWER = 0x02;
    public static final int STATUS_SUCCESS = 0x01;
    public static final int STATUS_ERROR = 0x02;
    public static final int MAX_FILENAME_SIZE = 255;
    public static final long MAX_FILE_SIZE = 0xFFFFFFFFL;

    private Protocol() {}
}