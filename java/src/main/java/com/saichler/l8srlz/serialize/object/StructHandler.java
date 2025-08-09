package com.saichler.l8srlz.serialize.object;

import com.google.protobuf.Message;

/**
 * Handler for struct/protobuf serialization
 * Translated from Go Struct.go
 */
public class StructHandler implements Complex {
    
    private final Int32Handler sizeHandler = new Int32Handler();
    private final StringHandler stringHandler = new StringHandler();
    
    @Override
    public Exception add(Object value, com.saichler.l8srlz.serialize.object.ByteBuffer data, Location location) {
        if (value == null) {
            sizeHandler.add(-1, data, location);
            return null;
        }
        
        if (!(value instanceof Message)) {
            return new Exception("Expected protobuf Message, got " + value.getClass());
        }
        
        Message message = (Message) value;
        String typeName = message.getClass().getSimpleName();
        
        try {
            // Marshal the protobuf message
            byte[] pbData = message.toByteArray();
            int size = pbData.length;
            
            // Calculate space needed: size + type name + data
            int typeNameBytes = typeName.getBytes("UTF-8").length;
            data.checkAndEnlarge(location.getValue(), 8 + typeNameBytes + size);
            
            if (size == 0) {
                sizeHandler.add(-2, data, location); // Empty message marker
            } else {
                sizeHandler.add(size, data, location);
            }
            
            // Write type name
            stringHandler.add(typeName, data, location);
            
            // Write protobuf data
            if (size > 0) {
                data.copyFrom(location.getValue(), pbData, 0, size);
                location.increment(size);
            }
            
            return null;
            
        } catch (Exception e) {
            return new Exception("Failed to marshal protobuf " + typeName + ": " + e.getMessage(), e);
        }
    }
    
    @Override
    public Result get(com.saichler.l8srlz.serialize.object.ByteBuffer data, Location location, IRegistry registry) {
        try {
            Integer sizeObj = (Integer) sizeHandler.get(data, location);
            int size = sizeObj;
            
            if (size == -1 || size == 0) {
                return new Result(null, null);
            }
            
            // Read type name
            String typeName = (String) stringHandler.get(data, location);
            
            // Get type info from registry
            if (registry == null) {
                return new Result(null, new Exception("Registry is required for protobuf deserialization"));
            }
            
            IInfo info;
            try {
                info = registry.info(typeName);
            } catch (Exception e) {
                return new Result(null, new Exception("Unknown protobuf type " + typeName + " in registry, please register it."));
            }
            
            // Create new instance
            Object instance;
            try {
                instance = info.newInstance();
            } catch (Exception e) {
                return new Result(null, new Exception("Error creating instance of " + typeName + ": " + e.getMessage()));
            }
            
            // If size is -2, it's an empty message
            if (size == -2) {
                return new Result(instance, null);
            }
            
            // Read and unmarshal protobuf data
            byte[] pbData = new byte[size];
            data.copyTo(location.getValue(), pbData, 0, size);
            location.increment(size);
            
            if (instance instanceof Message.Builder) {
                Message.Builder builder = (Message.Builder) instance;
                try {
                    builder.mergeFrom(pbData);
                    return new Result(builder.build(), null);
                } catch (Exception e) {
                    return new Result(null, new Exception("Failed to unmarshal protobuf " + typeName + ": " + e.getMessage()));
                }
            } else if (instance instanceof Message) {
                // Try to use parseFrom if available
                try {
                    Message message = (Message) instance;
                    Message.Builder builder = message.newBuilderForType();
                    builder.mergeFrom(pbData);
                    return new Result(builder.build(), null);
                } catch (Exception e) {
                    return new Result(null, new Exception("Failed to unmarshal protobuf " + typeName + ": " + e.getMessage()));
                }
            } else {
                return new Result(null, new Exception("Instance is not a protobuf Message: " + instance.getClass()));
            }
            
        } catch (Exception e) {
            return new Result(null, e);
        }
    }
}