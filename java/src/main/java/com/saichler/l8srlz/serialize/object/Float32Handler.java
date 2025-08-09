package com.saichler.l8srlz.serialize.object;

import java.nio.ByteBuffer;
import java.nio.ByteOrder;

/**
 * Handler for float32 serialization
 * Translated from Go Float32.go
 */
public class Float32Handler implements Primitive {
    
    @Override
    public void add(Object value, com.saichler.l8srlz.serialize.object.ByteBuffer data, Location location) {
        data.checkAndEnlarge(location.getValue(), 4);
        
        float floatValue;
        if (value instanceof Float) {
            floatValue = (Float) value;
        } else {
            floatValue = ((Number) value).floatValue();
        }
        
        // Convert float to IEEE 754 bits and use big-endian encoding
        int bits = Float.floatToIntBits(floatValue);
        ByteBuffer bb = ByteBuffer.allocate(4).order(ByteOrder.BIG_ENDIAN);
        bb.putInt(bits);
        byte[] bytes = bb.array();
        
        data.copyFrom(location.getValue(), bytes, 0, 4);
        location.increment(4);
    }
    
    @Override
    public Object get(com.saichler.l8srlz.serialize.object.ByteBuffer data, Location location) {
        ByteBuffer bb = ByteBuffer.allocate(4).order(ByteOrder.BIG_ENDIAN);
        byte[] bytes = new byte[4];
        data.copyTo(location.getValue(), bytes, 0, 4);
        bb.put(bytes);
        bb.rewind();
        
        int bits = bb.getInt();
        float result = Float.intBitsToFloat(bits);
        location.increment(4);
        return result;
    }
}