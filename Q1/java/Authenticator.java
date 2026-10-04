import java.util.HashMap;
import java.util.Map;

public class Authenticator {
        private final Map<String, String> users = new HashMap<>();

        public Authenticator() {
            /* senha: 123mudar */
            users.put("Carol", "01c80f0a14035a8977ad7cfe0fd91e2f9c6def942f2447c4744e0a8b4c639bcef281ad9e2b97efe20accbf8b52ab7d7bafa0956f34d6d9e574ed12fa1a3ec56d");
            /* senha: 123456 */
            users.put("Gustavo", "ba3253876aed6bc22d4a6ff53d8406c6ad864195ed144ab5c87621b6c233b548baeae6956df346ec8c17f5ea10f35ee3cbc514797ed7ddd3145464e2a0bab413");
            /* senha: admin */
            users.put("Maria", "c7ad44cbad762a5da0a452f9e854fdc1e0e7a52a38015f23f3eab1d80b931dd472634dfac71cd34ebc35d16ab7fb8a90c81f975113d6c7538dc69dd8de9077ec");
        }

        public boolean authenticate(String user, String passwordHash) {
            String correctHash = users.get(user);
            
            if (correctHash == null) {
                return false;
            }

            return correctHash.equalsIgnoreCase(passwordHash);
        }
}