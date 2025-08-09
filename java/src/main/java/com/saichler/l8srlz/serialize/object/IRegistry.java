package com.saichler.l8srlz.serialize.object;

/**
 * Interface for type registry that manages protobuf message types
 */
public interface IRegistry {
    
    /**
     * Register a protobuf message type by name
     * @param typeName the name of the type
     * @param instance an instance of the type for reflection
     */
    void register(String typeName, Object instance);
    
    /**
     * Get type information for a given type name
     * @param typeName the name of the type
     * @return type info containing instantiation details
     * @throws Exception if type is not registered
     */
    IInfo info(String typeName) throws Exception;
    
    /**
     * Check if a type is registered
     * @param typeName the name of the type
     * @return true if registered, false otherwise
     */
    boolean isRegistered(String typeName);
}