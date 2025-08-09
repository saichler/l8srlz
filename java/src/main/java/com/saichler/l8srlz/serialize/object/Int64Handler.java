package com.saichler.l8srlz.serialize.object;

import java.nio.ByteBuffer;
import java.nio.ByteOrder;

/**
 * Handler for int64 serialization using optimized binary encoding
 * Translated from Go Int64.go with performance optimizations
 */
public class Int64Handler implements Primitive {
    
    @Override
    public void add(Object value, com.saichler.l8srlz.serialize.object.ByteBuffer data, Location location) {
        data.checkAndEnlarge(location.getValue(), 8);
        
        long longValue;
        if (value instanceof Long) {
            longValue = (Long) value;
        } else {
            // Handle primitive long or other numeric types
            longValue = ((Number) value).longValue();
        }
        
        // Use Java's ByteBuffer for efficient big-endian encoding
        ByteBuffer bb = ByteBuffer.allocate(8).order(ByteOrder.BIG_ENDIAN);
        bb.putLong(longValue);
        byte[] bytes = bb.array();
        
        data.copyFrom(location.getValue(), bytes, 0, 8);
        location.increment(8);
    }
    
    @Override
    public Object get(com.saichler.l8srlz.serialize.object.ByteBuffer data, Location location) {
        // Use Java's ByteBuffer for efficient big-endian decoding
        ByteBuffer bb = ByteBuffer.allocate(8).order(ByteOrder.BIG_ENDIAN);
        byte[] bytes = new byte[8];
        data.copyTo(location.getValue(), bytes, 0, 8);
        bb.put(bytes);
        bb.rewind();
        
        long result = bb.getLong();
        location.increment(8);
        return result;
    }
}