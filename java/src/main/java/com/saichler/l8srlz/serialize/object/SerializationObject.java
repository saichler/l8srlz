package com.saichler.l8srlz.serialize.object;

import java.util.*;
import java.util.concurrent.ConcurrentHashMap;
import java.lang.reflect.*;
import java.util.Base64;

/**
 * Main serialization object that handles encoding and decoding of Java objects
 * Translated from Go Object.go
 */
public class SerializationObject {
    
    // Static registries for primitive and complex type handlers
    private static final Map<Class<?>, Primitive> primitives = new ConcurrentHashMap<>();
    private static final Map<Class<?>, Complex> complex = new ConcurrentHashMap<>();
    
    // Type handlers for size and string encoding
    private static final Int32Handler sizeObjectType = new Int32Handler();
    private static final StringHandler stringObjectType = new StringHandler();
    
    // Instance fields - package-private for access by handlers
    ByteBuffer data;
    Location location;
    IRegistry registry;
    
    static {
        initializeHandlers();
    }
    
    private static void initializeHandlers() {
        // Register primitive type handlers
        primitives.put(Integer.class, new IntHandler());
        primitives.put(int.class, new IntHandler());
        primitives.put(Long.class, new Int64Handler());
        primitives.put(long.class, new Int64Handler());
        primitives.put(Float.class, new Float32Handler());
        primitives.put(float.class, new Float32Handler());
        primitives.put(Double.class, new Float64Handler());
        primitives.put(double.class, new Float64Handler());
        primitives.put(String.class, new StringHandler());
        primitives.put(Boolean.class, new BoolHandler());
        primitives.put(boolean.class, new BoolHandler());
        
        // Register complex type handlers
        complex.put(List.class, new SliceHandler());
        complex.put(ArrayList.class, new SliceHandler());
        complex.put(Map.class, new MapHandler());
        complex.put(HashMap.class, new MapHandler());
        // Object/Struct handler will be added for protobuf messages
    }
    
    /**
     * Create a new encoding object
     */
    public static SerializationObject newEncode() {
        SerializationObject obj = new SerializationObject();
        obj.data = new ByteBuffer();
        obj.location = new Location();
        return obj;
    }
    
    /**
     * Create a new decoding object
     */
    public static SerializationObject newDecode(byte[] data, int location, IRegistry registry) {
        SerializationObject obj = new SerializationObject();
        obj.data = new ByteBuffer(data);
        obj.location = new Location(location);
        obj.registry = registry;
        return obj;
    }
    
    /**
     * Get the serialized data
     */
    public byte[] getData() {
        return data.copyRange(0, location.getValue());
    }
    
    /**
     * Get current location
     */
    public int getLocation() {
        return location.getValue();
    }
    
    /**
     * Add an object to the serialization buffer
     */
    public Exception add(Object value) throws Exception {
        if (value == null) {
            return addNull();
        }
        
        Class<?> valueClass = value.getClass();
        
        // Handle primitive types first (type switch optimization from performance analysis)
        if (valueClass == Integer.class || valueClass == int.class) {
            addKind(TypeKind.INT32);
            primitives.get(Integer.class).add(value, data, location);
            return null;
        } else if (valueClass == String.class) {
            addKind(TypeKind.STRING);
            primitives.get(String.class).add(value, data, location);
            return null;
        } else if (valueClass == Boolean.class || valueClass == boolean.class) {
            addKind(TypeKind.BOOL);
            primitives.get(Boolean.class).add(value, data, location);
            return null;
        } else if (valueClass == Long.class || valueClass == long.class) {
            addKind(TypeKind.INT64);
            primitives.get(Long.class).add(value, data, location);
            return null;
        } else if (valueClass == Double.class || valueClass == double.class) {
            addKind(TypeKind.FLOAT64);
            primitives.get(Double.class).add(value, data, location);
            return null;
        } else if (valueClass == Float.class || valueClass == float.class) {
            addKind(TypeKind.FLOAT32);
            primitives.get(Float.class).add(value, data, location);
            return null;
        }
        
        // Handle byte array specially
        if (value instanceof byte[]) {
            addKind(TypeKind.SLICE);
            Complex handler = complex.get(List.class);
            return handler.add(value, data, location);
        }
        
        // Handle collections
        if (value instanceof List) {
            addKind(TypeKind.SLICE);
            Complex handler = complex.get(List.class);
            return handler.add(value, data, location);
        } else if (value instanceof Map) {
            addKind(TypeKind.MAP);
            Complex handler = complex.get(Map.class);
            return handler.add(value, data, location);
        }
        
        // Handle protobuf messages and other objects
        if (com.google.protobuf.Message.class.isAssignableFrom(valueClass)) {
            addKind(TypeKind.PTR);
            Complex handler = new StructHandler();
            return handler.add(value, data, location);
        }
        
        // Fallback for reflection-based handling
        return addReflection(value);
    }
    
    /**
     * Get an object from the serialization buffer
     */
    public Object get() throws Exception {
        TypeKind kind = getKind();
        
        // Handle primitives by kind
        Primitive primitiveHandler = getPrimitiveHandler(kind);
        if (primitiveHandler != null) {
            return primitiveHandler.get(data, location);
        }
        
        Complex complexHandler = getComplexHandler(kind);
        if (complexHandler != null) {
            Complex.Result result = complexHandler.get(data, location, registry);
            if (result.hasError()) {
                throw result.getError();
            }
            return result.getValue();
        }
        
        throw new Exception("Did not find any handler for kind " + kind);
    }
    
    /**
     * Get Base64 representation of the data
     */
    public String toBase64() {
        return Base64.getEncoder().encodeToString(getData());
    }
    
    /**
     * Decode from Base64 string
     */
    public static byte[] fromBase64(String base64String) throws Exception {
        try {
            return Base64.getDecoder().decode(base64String);
        } catch (IllegalArgumentException e) {
            throw new Exception("Invalid Base64 string", e);
        }
    }
    
    /**
     * Utility method to serialize any object to bytes
     */
    public static byte[] dataOf(Object element) throws Exception {
        if (element == null) {
            return null;
        }
        SerializationObject obj = newEncode();
        Exception err = obj.add(element);
        if (err != null) {
            throw err;
        }
        return obj.getData();
    }
    
    /**
     * Utility method to deserialize bytes to object
     */
    public static Object elementOf(byte[] data, IRegistry registry) throws Exception {
        if (data == null) {
            return null;
        }
        SerializationObject obj = newDecode(data, 0, registry);
        return obj.get();
    }
    
    // Private helper methods
    
    private Exception addNull() {
        addKind(TypeKind.PTR);
        sizeObjectType.add(-1, data, location);
        return null;
    }
    
    private Exception addReflection(Object value) {
        // Fallback reflection-based serialization
        Class<?> valueClass = value.getClass();
        
        if (valueClass.isArray()) {
            addKind(TypeKind.SLICE);
            Complex handler = complex.get(List.class);
            return handler.add(value, data, location);
        }
        
        // Default to struct handling for other objects
        addKind(TypeKind.PTR);
        Complex handler = new StructHandler();
        return handler.add(value, data, location);
    }
    
    private void addKind(TypeKind kind) {
        sizeObjectType.add(kind.ordinal(), data, location);
    }
    
    private TypeKind getKind() {
        Integer kindInt = (Integer) sizeObjectType.get(data, location);
        return TypeKind.values()[kindInt];
    }
    
    private Primitive getPrimitiveHandler(TypeKind kind) {
        switch (kind) {
            case INT32: return primitives.get(Integer.class);
            case INT64: return primitives.get(Long.class);
            case FLOAT32: return primitives.get(Float.class);
            case FLOAT64: return primitives.get(Double.class);
            case STRING: return primitives.get(String.class);
            case BOOL: return primitives.get(Boolean.class);
            case INT: return primitives.get(Integer.class);
            default: return null;
        }
    }
    
    private Class<?> getClassForKind(TypeKind kind) {
        switch (kind) {
            case INT32: return Integer.class;
            case INT64: return Long.class;
            case FLOAT32: return Float.class;
            case FLOAT64: return Double.class;
            case STRING: return String.class;
            case BOOL: return Boolean.class;
            default: return null;
        }
    }
    
    private Complex getComplexHandler(TypeKind kind) {
        switch (kind) {
            case SLICE: return complex.get(List.class);
            case MAP: return complex.get(Map.class);
            case PTR: return new StructHandler();
            default: return null;
        }
    }
    
    /**
     * Type kinds that mirror Go's reflect.Kind
     */
    public enum TypeKind {
        INVALID,
        BOOL,
        INT,
        INT8,
        INT16,
        INT32,
        INT64,
        UINT,
        UINT8,
        UINT16,
        UINT32,
        UINT64,
        UINTPTR,
        FLOAT32,
        FLOAT64,
        COMPLEX64,
        COMPLEX128,
        ARRAY,
        CHAN,
        FUNC,
        INTERFACE,
        MAP,
        PTR,
        SLICE,
        STRING,
        STRUCT,
        UNSAFEPOINTER
    }
}