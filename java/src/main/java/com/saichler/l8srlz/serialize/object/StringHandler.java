package com.saichler.l8srlz.serialize.object;

import java.nio.charset.StandardCharsets;

/**
 * Handler for string serialization
 * Translated from Go String.go with performance optimizations
 */
public class StringHandler implements Primitive {
    
    private final Int32Handler sizeHandler = new Int32Handler();
    
    @Override
    public void add(Object value, com.saichler.l8srlz.serialize.object.ByteBuffer data, Location location) {
        String str = (String) value;
        
        // Convert string to UTF-8 bytes (Java default for network protocols)
        byte[] strBytes = str.getBytes(StandardCharsets.UTF_8);
        int length = strBytes.length;
        
        // Write length first
        sizeHandler.add(length, data, location);
        
        // Write string bytes
        data.checkAndEnlarge(location.getValue(), length);
        if (length > 0) {
            data.copyFrom(location.getValue(), strBytes, 0, length);
            location.increment(length);
        }
    }
    
    @Override
    public Object get(com.saichler.l8srlz.serialize.object.ByteBuffer data, Location location) {
        // Read length first
        Integer lengthObj = (Integer) sizeHandler.get(data, location);
        int length = lengthObj;
        
        if (length == 0) {
            return "";
        }
        
        // Read string bytes
        byte[] strBytes = new byte[length];
        data.copyTo(location.getValue(), strBytes, 0, length);
        location.increment(length);
        
        // Convert from UTF-8 bytes to string
        return new String(strBytes, StandardCharsets.UTF_8);
    }
}