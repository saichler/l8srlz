package com.saichler.l8srlz.serialize.object;

/**
 * Interface for type information used for creating instances
 */
public interface IInfo {
    
    /**
     * Create a new instance of the registered type
     * @return new instance of the type
     * @throws Exception if instantiation fails
     */
    Object newInstance() throws Exception;
    
    /**
     * Get the type name
     * @return the name of the type
     */
    String getTypeName();
    
    /**
     * Get the class of the type
     * @return the Class object for the type
     */
    Class<?> getTypeClass();
}