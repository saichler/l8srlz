package com.saichler.l8srlz.serialize.object;

import java.lang.reflect.Array;
import java.util.*;

/**
 * Handler for slice/array serialization
 * Translated from Go Slice.go
 */
public class SliceHandler implements Complex {
    
    private final Int32Handler sizeHandler = new Int32Handler();
    
    @Override
    public Exception add(Object value, com.saichler.l8srlz.serialize.object.ByteBuffer data, Location location) {
        if (value == null) {
            sizeHandler.add(-1, data, location);
            return null;
        }
        
        // Handle byte arrays specially
        if (value instanceof byte[]) {
            byte[] byteArray = (byte[]) value;
            if (byteArray.length == 0) {
                sizeHandler.add(-1, data, location);
                return null;
            }
            
            sizeHandler.add(byteArray.length, data, location);
            data.checkAndEnlarge(location.getValue(), 1);
            data.set(location.getValue(), (byte) 1); // Mark as byte array
            location.increment(1);
            
            data.checkAndEnlarge(location.getValue(), byteArray.length);
            data.copyFrom(location.getValue(), byteArray, 0, byteArray.length);
            location.increment(byteArray.length);
            return null;
        }
        
        // Handle Lists and Collections
        Collection<?> collection = null;
        int size = 0;
        
        if (value instanceof Collection) {
            collection = (Collection<?>) value;
            size = collection.size();
        } else if (value.getClass().isArray()) {
            size = Array.getLength(value);
        } else {
            return new Exception("Unsupported slice type: " + value.getClass());
        }
        
        if (size == 0) {
            sizeHandler.add(-1, data, location);
            return null;
        }
        
        sizeHandler.add(size, data, location);
        data.checkAndEnlarge(location.getValue(), 1);
        data.set(location.getValue(), (byte) 0); // Mark as object array
        location.increment(1);
        
        // Serialize each element
        SerializationObject obj = new SerializationObject();
        obj.data = data;
        obj.location = location;
        
        try {
            if (collection != null) {
                for (Object element : collection) {
                    Exception err = obj.add(element);
                    if (err != null) {
                        return err;
                    }
                }
            } else {
                // Handle arrays
                for (int i = 0; i < size; i++) {
                    Object element = Array.get(value, i);
                    Exception err = obj.add(element);
                    if (err != null) {
                        return err;
                    }
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
            
            byte arrayType = data.get(location.getValue());
            location.increment(1);
            
            if (arrayType == 1) {
                // Byte array
                byte[] result = new byte[size];
                data.copyTo(location.getValue(), result, 0, size);
                location.increment(size);
                return new Result(result, null);
            } else {
                // Object array - deserialize as List
                List<Object> elements = new ArrayList<>(size);
                
                SerializationObject obj = new SerializationObject();
                obj.data = data;
                obj.location = location;
                obj.registry = registry;
                
                for (int i = 0; i < size; i++) {
                    Object element = obj.get();
                    elements.add(element);
                }
                
                return new Result(elements, null);
            }
            
        } catch (Exception e) {
            return new Result(null, e);
        }
    }
}