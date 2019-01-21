/*
 * Author: Igor joaquim dos Santos Lima
 * Github: https://github.com/igor036
 * 
 */
package br.com.linecode.venus.socket;

import java.io.IOException;
import java.net.InetAddress;
import java.net.ServerSocket;
import java.net.Socket;
import java.util.HashMap;
import java.util.Scanner;

import br.com.linecode.venus.socket.entity.IServerRunnable;

public class Server {

	
	private static final int PORT = 180;
	
	private ServerSocket server;
	private HashMap<InetAddress, Socket> clients;
	
	public Server() {
		
		clients = new HashMap<InetAddress, Socket>();
		serverInit();
		
	}
	
	/**
	 * Start the server, accept clients and
	 * read data from clients. 
	 */
	private void serverInit() {
		try {
			
			server = new ServerSocket(PORT);
			System.out.printf("Listen port %d ...\n",PORT);
			
			//accept clients thread
			new Thread(((IServerRunnable)() -> {
				
				Socket client = server.accept();
				Scanner clientInput = new Scanner(client.getInputStream());
				
				System.out.printf("Client accept: %s\n",client.getInetAddress());
				clients.put(client.getInetAddress(), client);
				
				//read client thread
				new Thread(() ->  {
					while (clientInput.hasNext()) 
						System.out.printf("%s: %s\n", client.getInetAddress(), clientInput.nextLine());
				}).start();
				
			})).start();
			
		} catch (IOException e) {
			e.printStackTrace();
			System.exit(1);
		}
	}
}
