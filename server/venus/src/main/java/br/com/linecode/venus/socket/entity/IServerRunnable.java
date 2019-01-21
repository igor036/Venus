/*
 * Author: Igor joaquim dos Santos Lima
 * Github: https://github.com/igor036
 * 
 */
package br.com.linecode.venus.socket.entity;

import java.io.IOException;

public interface IServerRunnable extends Runnable {
	
	/**
	 * method that throws an IO exception 
	 */
	void consumer() throws IOException;
	
	
	/**
	 * function that will be run in all server lifetime
	 */
	@Override
	default void run() {
		while(true) {
			try {
				consumer();
			} catch (IOException e) {
				e.printStackTrace();
				System.exit(1);
			}
		}
	}
}
