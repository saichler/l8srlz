package com.saichler.l8srlz.serialize.object;

import java.util.*;

/**
 * Handler for map serialization
 * Translated from Go Map.go
 */
public class MapHandler implements Complex {
    
    private final Int32Handler sizeHandler = new Int32Handler();
    
    @Override
    public Exception add(Object value, com.saichler.l8srlz.serialize.object.ByteBuffer data, Location location) {
        if (value == null) {
            sizeHandler.add(-1, data, location);
            return null;
        }
        
        if (!(value instanceof Map)) {
            return new Exception("Expected Map, got " + value.getClass());
        }
        
        Map<?, ?> map = (Map<?, ?>) value;
        if (map.isEmpty()) {
            sizeHandler.add(-1, data, location);
            return null;
        }
        
        sizeHandler.add(map.size(), data, location);
        
        // Serialize each key-value pair
        SerializationObject obj = new SerializationObject();
        obj.data = data;
        obj.location = location;
        
        try {
            for (Map.Entry<?, ?> entry : map.entrySet()) {
                // Serialize key
                Exception err = obj.add(entry.getKey());
                if (err != null) {
                    return err;
                }
                
                // Serialize value
                err = obj.add(entry.getValue());
                if (err != null) {
                    return err;
                }
            }
        } catch (Exception e) {
            return e;
        }
        
        return null;
    }
    
    @Override
    public Result get(com.saichler.l8srlz.serialize.object.ByteBuffer data, Location location, IRegistry registry) {
        try {
            Integer sizeObj = (Integer) sizeHandler.get(data, location);
            int size = sizeObj;
            
            if (size == -1 || size == 0) {
                return new Result(null, null);
            }
            
            Map<Object, Object> result = new HashMap<>(size);
            
            SerializationObject obj = new SerializationObject();
            obj.data = data;
            obj.location = location;
            obj.registry = registry;
            
            for (int i = 0; i < size; i++) {
                // Deserialize key
                Object key = obj.get();
                
                // Deserialize value
                Object val = obj.get();
                
                result.put(key, val);
            }
            
            return new Result(result, null);
            
        } catch (Exception e) {
            return new Result(null, e);
        }
    }
}