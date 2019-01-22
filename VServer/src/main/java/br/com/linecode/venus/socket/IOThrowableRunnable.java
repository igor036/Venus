/*
 * Author: Igor joaquim dos Santos Lima
 * Github: https://github.com/igor036
 * 
 */
package br.com.linecode.venus.socket;

import java.io.IOException;

public interface IOThrowableRunnable extends Runnable     {
    
    /**
     * method that throws an IO exception 
     */
    void consumer() throws IOException;
    
    @Override
    @SuppressWarnings("squid:S1148")
    default void run() {
        try {
            consumer();
        } catch (IOException e) {
            e.printStackTrace();
        }
    }
}
