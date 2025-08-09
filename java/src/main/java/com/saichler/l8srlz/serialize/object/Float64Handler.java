package com.saichler.l8srlz.serialize.object;

import java.nio.ByteBuffer;
import java.nio.ByteOrder;

/**
 * Handler for float64 serialization
 * Translated from Go Float64.go with optimized binary encoding
 */
public class Float64Handler implements Primitive {
    
    @Override
    public void add(Object value, com.saichler.l8srlz.serialize.object.ByteBuffer data, Location location) {
        data.checkAndEnlarge(location.getValue(), 8);
        
        double doubleValue;
        if (value instanceof Double) {
            doubleValue = (Double) value;
        } else {
            doubleValue = ((Number) value).doubleValue();
        }
        
        // Convert double to IEEE 754 bits and use big-endian encoding
        long bits = Double.doubleToLongBits(doubleValue);
        ByteBuffer bb = ByteBuffer.allocate(8).order(ByteOrder.BIG_ENDIAN);
        bb.putLong(bits);
        byte[] bytes = bb.array();
        
        data.copyFrom(location.getValue(), bytes, 0, 8);
        location.increment(8);
    }
    
    @Override
    public Object get(com.saichler.l8srlz.serialize.object.ByteBuffer data, Location location) {
        ByteBuffer bb = ByteBuffer.allocate(8).order(ByteOrder.BIG_ENDIAN);
        byte[] bytes = new byte[8];
        data.copyTo(location.getValue(), bytes, 0, 8);
        bb.put(bytes);
        bb.rewind();
        
        long bits = bb.getLong();
        double result = Double.longBitsToDouble(bits);
        location.increment(8);
        return result;
    }
}