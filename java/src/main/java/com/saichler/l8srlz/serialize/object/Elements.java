package com.saichler.l8srlz.serialize.object;

import java.util.*;
import java.lang.reflect.Array;

/**
 * Container for serialized elements with query support
 * Translated from Go Elements.go (simplified without GSQL dependencies)
 */
public class Elements {
    
    private List<Element> elements;
    private boolean notification = false;
    private boolean replicasRequest = false;
    
    /**
     * Create new Elements container with single value
     */
    public static Elements create(Exception error, Object value) {
        if (value != null && value instanceof Runnable) {
            throw new RuntimeException("value is a function, this is probably a mistake");
        }
        
        Elements result = new Elements();
        result.elements = new ArrayList<>();
        
        if (value == null) {
            Element element = new Element();
            element.error = error;
            result.elements.add(element);
            return result;
        }
        
        // Handle collections
        if (value instanceof Collection) {
            Collection<?> collection = (Collection<?>) value;
            int index = 0;
            for (Object item : collection) {
                result.add(item, index++, null);
            }
        } else if (value instanceof Map) {
            Map<?, ?> map = (Map<?, ?>) value;
            for (Map.Entry<?, ?> entry : map.entrySet()) {
                result.add(entry.getValue(), entry.getKey(), null);
            }
        } else if (value.getClass().isArray()) {
            int length = Array.getLength(value);
            for (int i = 0; i < length; i++) {
                Object item = Array.get(value, i);
                result.add(item, i, null);
            }
        } else {
            // Single element
            Element element = new Element();
            element.element = value;
            element.error = error;
            result.elements.add(element);
        }
        
        return result;
    }
    
    /**
     * Create Elements with error
     */
    public static Elements createError(String errorMessage) {
        return create(new Exception(errorMessage), null);
    }
    
    /**
     * Create notification Elements
     */
    public static Elements createNotify(Object value) {
        Elements elems = create(null, value);
        elems.notification = true;
        return elems;
    }
    
    /**
     * Create replicas request Elements
     */
    public static Elements createReplicasRequest(Elements source) {
        Elements copy = new Elements();
        copy.elements = new ArrayList<>(source.elements);
        copy.notification = source.notification;
        copy.replicasRequest = true;
        return copy;
    }
    
    /**
     * Add element with key and error
     */
    public void add(Object element, Object key, Exception error) {
        if (elements == null) {
            elements = new ArrayList<>();
        }
        
        Element elem = new Element();
        elem.element = element;
        elem.key = key;
        elem.error = error;
        elements.add(elem);
    }
    
    /**
     * Get all elements as list
     */
    public List<Object> getElements() {
        if (elements == null) {
            return new ArrayList<>();
        }
        
        List<Object> result = new ArrayList<>(elements.size());
        for (Element elem : elements) {
            result.add(elem.element);
        }
        return result;
    }
    
    /**
     * Get first element
     */
    public Object getElement() {
        if (elements == null || elements.isEmpty()) {
            return null;
        }
        return elements.get(0).element;
    }
    
    /**
     * Get all keys as list
     */
    public List<Object> getKeys() {
        if (elements == null) {
            return new ArrayList<>();
        }
        
        List<Object> result = new ArrayList<>(elements.size());
        for (Element elem : elements) {
            result.add(elem.key);
        }
        return result;
    }
    
    /**
     * Get first key
     */
    public Object getKey() {
        if (elements == null || elements.isEmpty()) {
            return null;
        }
        return elements.get(0).key;
    }
    
    /**
     * Get all errors as list
     */
    public List<Exception> getErrors() {
        if (elements == null) {
            return new ArrayList<>();
        }
        
        List<Exception> result = new ArrayList<>(elements.size());
        for (Element elem : elements) {
            result.add(elem.error);
        }
        return result;
    }
    
    /**
     * Get first error
     */
    public Exception getError() {
        if (elements == null || elements.isEmpty()) {
            return null;
        }
        return elements.get(0).error;
    }
    
    /**
     * Serialize Elements to bytes
     */
    public byte[] serialize() throws Exception {
        SerializationObject obj = SerializationObject.newEncode();
        
        // Serialize number of elements
        int size = (elements != null) ? elements.size() : 0;
        Exception err = obj.add(size);
        if (err != null) {
            throw err;
        }
        
        if (elements != null) {
            for (Element elem : elements) {
                // Serialize element
                err = obj.add(elem.element);
                if (err != null) {
                    throw err;
                }
                
                // Serialize key
                err = obj.add(elem.key);
                if (err != null) {
                    throw err;
                }
                
                // Serialize error message
                String errorMsg = (elem.error != null) ? elem.error.getMessage() : "";
                err = obj.add(errorMsg);
                if (err != null) {
                    throw err;
                }
            }
        }
        
        // Note: pquery serialization skipped as it requires GSQL dependencies
        err = obj.add((Object) null); // Placeholder for pquery
        if (err != null) {
            throw err;
        }
        
        return obj.getData();
    }
    
    /**
     * Deserialize Elements from bytes
     */
    public void deserialize(byte[] data, IRegistry registry) throws Exception {
        SerializationObject obj = SerializationObject.newDecode(data, 0, registry);
        
        // Read number of elements
        Object sizeObj = obj.get();
        int size = (Integer) sizeObj;
        
        elements = new ArrayList<>(size);
        
        for (int i = 0; i < size; i++) {
            Element elem = new Element();
            
            // Read element
            elem.element = obj.get();
            
            // Read key
            elem.key = obj.get();
            
            // Read error message
            Object errorMsgObj = obj.get();
            String errorMsg = (String) errorMsgObj;
            if (errorMsg != null && !errorMsg.isEmpty()) {
                elem.error = new Exception(errorMsg);
            }
            
            elements.add(elem);
        }
        
        // Read pquery (placeholder)
        obj.get(); // Skip pquery for now
    }
    
    /**
     * Check if this is a notification
     */
    public boolean isNotification() {
        return notification;
    }
    
    /**
     * Check if this is a replicas request
     */
    public boolean isReplicasRequest() {
        return replicasRequest;
    }
    
    /**
     * Append another Elements to this one
     */
    public void append(Elements other) {
        if (other.elements != null) {
            if (this.elements == null) {
                this.elements = new ArrayList<>();
            }
            this.elements.addAll(other.elements);
        }
    }
    
    /**
     * Get number of elements
     */
    public int size() {
        return (elements != null) ? elements.size() : 0;
    }
    
    /**
     * Check if empty
     */
    public boolean isEmpty() {
        return elements == null || elements.isEmpty();
    }
    
    /**
     * Internal Element class
     */
    private static class Element {
        Object element;
        Object key;
        Exception error;
    }
}