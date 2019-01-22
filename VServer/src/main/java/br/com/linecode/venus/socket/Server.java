/*
 * Author: Igor joaquim dos Santos Lima
 * Github: https://github.com/igor036
 * 
 * only static methods for server socket.
 */
package br.com.linecode.venus.socket;

import java.io.IOException;
import java.net.ServerSocket;
import java.net.Socket;
import java.util.LinkedList;
import java.util.List;

public class Server {
	
	private static final int PORT = 180;
	
	private static boolean running;
	private static ServerSocket server;  
	private static List<Client> clients;
	
	private Server() {}
	
	
	/**
	 * Start the server, accept clients and
	 * read data from clients. 
	 */
	@SuppressWarnings({
	    "squid:S2589","squid:S1148",
	    "squid:S106", "squid:S3457"
	})
	public static void start() {
		try {
			
		    running = true;
			server = new ServerSocket(PORT);
			clients = new LinkedList<>();
			
			System.out.printf("Lisen for connections in port %d...\n",PORT);
			
		    new Thread(((IOThrowableRunnable)() -> {
		        
		        while(running) {
		            
	                Socket client = server.accept();
	                clients.add(new Client(client));
	                
	                System.out.printf("New connection with %s\n",client.getLocalAddress());
	                
	            }
	        })).start();
			
	
		} catch (IOException e) {
			e.printStackTrace();
			System.exit(1);
		}
	}
	
	/**
	 * Stop server and close all connections.
	 */
	public static void stop() {
	    running = false;
	    while(!clients.isEmpty()) {
	        Client client = clients.remove(0);
	        client.closeConnection();
	    }
	}
}
