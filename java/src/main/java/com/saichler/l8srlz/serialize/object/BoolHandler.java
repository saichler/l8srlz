package com.saichler.l8srlz.serialize.object;

/**
 * Handler for boolean serialization
 * Translated from Go Bool.go
 */
public class BoolHandler implements Primitive {
    
    @Override
    public void add(Object value, com.saichler.l8srlz.serialize.object.ByteBuffer data, Location location) {
        data.checkAndEnlarge(location.getValue(), 1);
        
        boolean boolValue;
        if (value instanceof Boolean) {
            boolValue = (Boolean) value;
        } else {
            // Handle primitive boolean
            boolValue = (boolean) value;
        }
        
        byte byteValue = boolValue ? (byte) 1 : (byte) 0;
        data.set(location.getValue(), byteValue);
        location.increment(1);
    }
    
    @Override
    public Object get(com.saichler.l8srlz.serialize.object.ByteBuffer data, Location location) {
        byte byteValue = data.get(location.getValue());
        location.increment(1);
        return byteValue == 1;
    }
}