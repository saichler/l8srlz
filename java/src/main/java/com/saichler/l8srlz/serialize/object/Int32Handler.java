package com.saichler.l8srlz.serialize.object;

import java.nio.ByteBuffer;
import java.nio.ByteOrder;

/**
 * Handler for int32 serialization using optimized binary encoding
 * Translated from Go Int32.go with performance optimizations
 */
public class Int32Handler implements Primitive {
    
    @Override
    public void add(Object value, com.saichler.l8srlz.serialize.object.ByteBuffer data, Location location) {
        data.checkAndEnlarge(location.getValue(), 4);
        
        int intValue;
        if (value instanceof Integer) {
            intValue = (Integer) value;
        } else {
            // Handle primitive int or other numeric types
            intValue = ((Number) value).intValue();
        }
        
        // Use Java's ByteBuffer for efficient big-endian encoding (performance optimization)
        ByteBuffer bb = ByteBuffer.allocate(4).order(ByteOrder.BIG_ENDIAN);
        bb.putInt(intValue);
        byte[] bytes = bb.array();
        
        data.copyFrom(location.getValue(), bytes, 0, 4);
        location.increment(4);
    }
    
    @Override
    public Object get(com.saichler.l8srlz.serialize.object.ByteBuffer data, Location location) {
        // Use Java's ByteBuffer for efficient big-endian decoding
        ByteBuffer bb = ByteBuffer.allocate(4).order(ByteOrder.BIG_ENDIAN);
        byte[] bytes = new byte[4];
        data.copyTo(location.getValue(), bytes, 0, 4);
        bb.put(bytes);
        bb.rewind();
        
        int result = bb.getInt();
        location.increment(4);
        return result;
    }
}