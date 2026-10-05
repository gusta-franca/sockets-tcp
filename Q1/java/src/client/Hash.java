/*
    Última atualização: 04/10/2026
    Aluna: Maria Eduarda Bambini
    Arquivo do cliente responsável por converter a senha digitada em hash SHA-512.
*/

package client;
 
import java.nio.charset.StandardCharsets;
import java.security.MessageDigest;
import java.security.NoSuchAlgorithmException;

/**
 * Hash: calcula o SHA-512 de um texto e devolve o resultado em hexadecimal minúsculo
 * (128 caracteres). O cliente usa isso para nunca enviar a senha em texto puro.
 */
public class Hash {

    public static String sha512(String text) {
        try {
            MessageDigest digest = MessageDigest.getInstance("SHA-512");

            byte[] textBytes = text.getBytes(StandardCharsets.UTF_8);
            byte[] hashBytes = digest.digest(textBytes);

            StringBuilder hex = new StringBuilder();

            for (int i = 0; i < hashBytes.length; i++) {
                byte currentByte = hashBytes[i];

                // %02x = número em hexadecimal, com 2 dígitos, completando com zero
                String twoChars = String.format("%02x", currentByte);

                hex.append(twoChars);
            }

            String result = hex.toString();
            return result;

        } catch (NoSuchAlgorithmException e) {
            throw new RuntimeException("ERROR", e);
        }
    }
}