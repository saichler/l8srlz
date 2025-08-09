package com.saichler.l8srlz.serialize.object;

/**
 * Interface for complex type serializers
 */
public interface Complex {
    
    /**
     * Add a complex value to the buffer
     * @param value the value to serialize
     * @param data the byte buffer
     * @param location current location in the buffer (will be updated)
     * @return null on success, exception on error
     */
    Exception add(Object value, ByteBuffer data, Location location);
    
    /**
     * Get a complex value from the buffer
     * @param data the byte buffer
     * @param location current location in the buffer (will be updated)
     * @param registry type registry for creating instances
     * @return the deserialized value and any error
     */
    Result get(ByteBuffer data, Location location, IRegistry registry);
    
    /**
     * Result wrapper for complex deserialization
     */
    class Result {
        private final Object value;
        private final Exception error;
        
        public Result(Object value, Exception error) {
            this.value = value;
            this.error = error;
        }
        
        public Object getValue() { return value; }
        public Exception getError() { return error; }
        public boolean hasError() { return error != null; }
    }
}