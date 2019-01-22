/*
 * Author: Igor joaquim dos Santos Lima
 * Github: https://github.com/igor036
 * 
 */
package br.com.linecode.venus.socket;

import java.io.IOException;
import java.net.Socket;
import java.util.Scanner;

@SuppressWarnings("squid:S106")
public class Client {
    
    private Scanner scanner;
    private Socket client;

    @SuppressWarnings("squid:S3457")
    public Client(Socket client) throws IOException {
        
        this.client = client;
        this.scanner = new Scanner(client.getInputStream());
        
        new Thread(() ->  {
            while(scanner.hasNext()) {
                System.out.printf("%s: %s\n", client.getLocalAddress(), scanner.nextLine());
            }
            scanner.close();
        }).start();
    }
    
    /**
     * Close client connection and stop thread of read.
     */
    @SuppressWarnings("squid:S1148")
    public void closeConnection() {
        try {
            client.close();
        }catch (IOException e) {
            e.printStackTrace();
        }
    }
}
