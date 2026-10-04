import java.nio.charset.StandardCharsets;
import java.security.MessageDigest;
import java.security.NoSuchAlgorithmException;

public class Hash {

    // Transforma um texto em hash SHA-512, escrito em hexadecimal minúsculo
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