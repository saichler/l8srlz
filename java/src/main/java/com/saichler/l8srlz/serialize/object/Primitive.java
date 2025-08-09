package com.saichler.l8srlz.serialize.object;

/**
 * Interface for primitive type serializers
 */
public interface Primitive {
    
    /**
     * Add a primitive value to the buffer
     * @param value the value to serialize
     * @param data the byte buffer
     * @param location current location in the buffer (will be updated)
     */
    void add(Object value, ByteBuffer data, Location location);
    
    /**
     * Get a primitive value from the buffer
     * @param data the byte buffer
     * @param location current location in the buffer (will be updated)
     * @return the deserialized value
     */
    Object get(ByteBuffer data, Location location);
}