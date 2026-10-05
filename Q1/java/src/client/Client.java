package client;
import java.io.*;
import java.net.*;
import java.util.Scanner;

import src.common.Config;
import src.common.Connection;

public class Client {
    public static void main(String[] tokens) {

        Config config = new Config(".env");

        String serverHost = config.getServerHost();
        int serverPort = config.getServerPort();
        
        try (Scanner reader = new Scanner(System.in);
             Socket socket = new Socket(serverHost, serverPort);
             Connection connection = new Connection(socket)) {

            while (true) {
                System.out.print("Mensagem: ");
                String buffer = reader.nextLine().trim();

                if (buffer.isEmpty()) {
                    continue;
                }

                String[] args = buffer.split(" ", 2);
                String command = args[0].toUpperCase();

                if (command.equals("CONNECT")) {

                    if (args.length != 2) {
                        System.out.println("Utilize o formato 'CONNECT <usuário>, <senha>'");
                        continue;
                    }

                    String[] values = args[1].split(" ");
                    if (values.length != 2) {
                        System.out.println("Utilize o formato 'CONNECT <usuário>, <senha>'");
                        continue;
                    }

                    String userWithComma = values[0];
                    String password = values[1];

                    if (userWithComma.endsWith(",") == false) {
                        System.out.println("Utilize o formato 'CONNECT <usuário>, <senha>'");
                        continue;
                    }

                    int userLastIndex = userWithComma.length() - 1;
                    String userName = userWithComma.substring(0, userLastIndex);
                    
                    String hash = Hash.sha512(password);

                    connection.send("CONNECT " + userName + ", " + hash);
                    System.out.println("Servidor: " + connection.receive());
                
                } else if (command.equals("EXIT")) {
                    connection.send("EXIT");
                    break;

                } else if (command.equals("GETFILES") || command.equals("GETDIRS")) {
                    connection.send(command);

                    String response = connection.receive();
                    if (response.equals("ERROR")) {
                        System.out.println("Server disse: ERROR.");
                    } else {
                        int count = Integer.parseInt(response);
                        System.out.println("Quantidade: " + count);

                        for (int i = 0; i < count; i++){
                            System.out.println(" - " + connection.receive());
                        }
                        
                    } 
                } else {
                    connection.send(buffer);
                    System.out.println("Server disse: " + connection.receive());
                }
            
            }

        } catch (UnknownHostException e) {
            System.out.println("Socket: " + e.getMessage());
        } catch (EOFException e) {
            System.out.println("EOF: " + e.getMessage());
        } catch (IOException e) {
            System.out.println("IO: " + e.getMessage());
        }
    }
}