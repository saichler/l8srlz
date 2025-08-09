package com.saichler.l8srlz.serialize.object;

import java.util.concurrent.ConcurrentHashMap;
import java.util.Map;
import java.lang.reflect.Method;
import com.google.protobuf.Message;

/**
 * Simple implementation of IRegistry for managing protobuf types
 */
public class SimpleRegistry implements IRegistry {
    
    private final Map<String, TypeInfo> types = new ConcurrentHashMap<>();
    
    @Override
    public void register(String typeName, Object instance) {
        if (instance == null) {
            throw new IllegalArgumentException("Instance cannot be null");
        }
        
        Class<?> clazz = instance.getClass();
        types.put(typeName, new TypeInfo(typeName, clazz, instance));
    }
    
    @Override
    public IInfo info(String typeName) throws Exception {
        TypeInfo info = types.get(typeName);
        if (info == null) {
            throw new Exception("Type not registered: " + typeName);
        }
        return info;
    }
    
    @Override
    public boolean isRegistered(String typeName) {
        return types.containsKey(typeName);
    }
    
    /**
     * Register a protobuf message using its class
     */
    public void register(Class<? extends Message> messageClass) {
        try {
            Message defaultInstance = getDefaultInstance(messageClass);
            String typeName = messageClass.getSimpleName();
            register(typeName, defaultInstance);
        } catch (Exception e) {
            throw new RuntimeException("Failed to register protobuf class " + messageClass, e);
        }
    }
    
    private Message getDefaultInstance(Class<? extends Message> messageClass) throws Exception {
        try {
            Method method = messageClass.getMethod("getDefaultInstance");
            return (Message) method.invoke(null);
        } catch (Exception e) {
            // Try to create using newBuilder
            Method newBuilder = messageClass.getMethod("newBuilder");
            Message.Builder builder = (Message.Builder) newBuilder.invoke(null);
            return builder.build();
        }
    }
    
    /**
     * Internal type information holder
     */
    private static class TypeInfo implements IInfo {
        private final String typeName;
        private final Class<?> typeClass;
        private final Object template;
        
        public TypeInfo(String typeName, Class<?> typeClass, Object template) {
            this.typeName = typeName;
            this.typeClass = typeClass;
            this.template = template;
        }
        
        @Override
        public Object newInstance() throws Exception {
            if (template instanceof Message) {
                Message message = (Message) template;
                return message.newBuilderForType();
            } else {
                return typeClass.newInstance();
            }
        }
        
        @Override
        public String getTypeName() {
            return typeName;
        }
        
        @Override
        public Class<?> getTypeClass() {
            return typeClass;
        }
    }
}