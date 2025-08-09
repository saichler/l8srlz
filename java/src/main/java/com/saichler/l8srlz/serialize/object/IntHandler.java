package com.saichler.l8srlz.serialize.object;

import java.nio.ByteBuffer;
import java.nio.ByteOrder;

/**
 * Handler for int serialization (maps to int32)
 * Translated from Go Int.go
 */
public class IntHandler implements Primitive {
    
    private final Int32Handler int32Handler = new Int32Handler();
    
    @Override
    public void add(Object value, com.saichler.l8srlz.serialize.object.ByteBuffer data, Location location) {
        // In Java, int is 32-bit, so we can delegate to Int32Handler
        int32Handler.add(value, data, location);
    }
    
    @Override
    public Object get(com.saichler.l8srlz.serialize.object.ByteBuffer data, Location location) {
        return int32Handler.get(data, location);
    }
}