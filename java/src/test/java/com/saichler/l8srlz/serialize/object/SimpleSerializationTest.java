package com.saichler.l8srlz.serialize.object;

import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.DisplayName;
import static org.junit.jupiter.api.Assertions.*;

import java.util.*;

/**
 * Simplified test suite for Java serialization library
 * Tests basic functionality without complex protobuf dependencies
 */
public class SimpleSerializationTest {
    
    private SimpleRegistry registry;
    
    @BeforeEach
    void setUp() {
        registry = new SimpleRegistry();
    }
    
    private Object testType(Object value) throws Exception {
        // Serialize
        SerializationObject encoder = SerializationObject.newEncode();
        Exception err = encoder.add(value);
        if (err != null) {
            throw err;
        }
        
        byte[] data = encoder.getData();
        
        // Deserialize
        SerializationObject decoder = SerializationObject.newDecode(data, 0, registry);
        return decoder.get();
    }
    
    @Test
    @DisplayName("Test Primitive Types Serialization")
    void testPrimitiveTypes() throws Exception {
        // Test integers
        assertEquals(42, testType(42));
        assertEquals(-42, testType(-42));
        assertEquals(0, testType(0));
        assertEquals(Integer.MAX_VALUE, testType(Integer.MAX_VALUE));
        assertEquals(Integer.MIN_VALUE, testType(Integer.MIN_VALUE));
        
        // Test longs
        assertEquals(9223372036854775807L, testType(9223372036854775807L));
        assertEquals(-9223372036854775807L, testType(-9223372036854775807L));
        assertEquals(0L, testType(0L));
        
        // Test floats
        assertEquals(3.14159f, (Float) testType(3.14159f), 0.00001f);
        assertEquals(-3.14159f, (Float) testType(-3.14159f), 0.00001f);
        assertEquals(0.0f, (Float) testType(0.0f), 0.00001f);
        
        // Test doubles
        assertEquals(3.141592653589793, (Double) testType(3.141592653589793), 0.000000000000001);
        assertEquals(-3.141592653589793, (Double) testType(-3.141592653589793), 0.000000000000001);
        assertEquals(0.0, (Double) testType(0.0), 0.000000000000001);
        
        // Test booleans
        assertEquals(true, testType(true));
        assertEquals(false, testType(false));
        
        // Test strings
        assertEquals("", testType(""));
        assertEquals("Hello, World!", testType("Hello, World!"));
        assertEquals("Hello, 世界! 🌍", testType("Hello, 世界! 🌍"));
        
        // Test long string
        String longString = "a".repeat(1000); // Smaller for quick test
        assertEquals(longString, testType(longString));
    }
    
    @Test
    @DisplayName("Test Collections Serialization")
    void testCollections() throws Exception {
        // Test Lists
        List<Integer> intList = Arrays.asList(1, 2, 3, 4, 5);
        @SuppressWarnings("unchecked")
        List<Integer> resultIntList = (List<Integer>) testType(intList);
        assertEquals(intList.size(), resultIntList.size());
        for (int i = 0; i < intList.size(); i++) {
            assertEquals(intList.get(i), resultIntList.get(i));
        }
        
        // Test empty list (should return null based on Go behavior)
        List<Integer> emptyList = new ArrayList<>();
        assertNull(testType(emptyList));
        
        // Test string list
        List<String> stringList = Arrays.asList("a", "b", "c");
        @SuppressWarnings("unchecked")
        List<String> resultStringList = (List<String>) testType(stringList);
        assertEquals(stringList.size(), resultStringList.size());
        for (int i = 0; i < stringList.size(); i++) {
            assertEquals(stringList.get(i), resultStringList.get(i));
        }
        
        // Test byte array
        byte[] byteArray = {0x48, 0x65, 0x6c, 0x6c, 0x6f}; // "Hello"
        byte[] resultByteArray = (byte[]) testType(byteArray);
        assertArrayEquals(byteArray, resultByteArray);
        
        // Test empty byte array (should return null)
        byte[] emptyByteArray = new byte[0];
        assertNull(testType(emptyByteArray));
    }
    
    @Test
    @DisplayName("Test Maps Serialization")
    void testMaps() throws Exception {
        // Test String to Integer map
        Map<String, Integer> stringIntMap = new HashMap<>();
        stringIntMap.put("one", 1);
        stringIntMap.put("two", 2);
        stringIntMap.put("three", 3);
        
        @SuppressWarnings("unchecked")
        Map<String, Integer> resultStringIntMap = (Map<String, Integer>) testType(stringIntMap);
        assertEquals(stringIntMap.size(), resultStringIntMap.size());
        for (String key : stringIntMap.keySet()) {
            assertEquals(stringIntMap.get(key), resultStringIntMap.get(key));
        }
        
        // Test empty map (should return null based on Go behavior)
        Map<String, Integer> emptyMap = new HashMap<>();
        assertNull(testType(emptyMap));
        
        // Test Integer to String map
        Map<Integer, String> intStringMap = new HashMap<>();
        intStringMap.put(1, "one");
        intStringMap.put(2, "two");
        intStringMap.put(3, "three");
        
        @SuppressWarnings("unchecked")
        Map<Integer, String> resultIntStringMap = (Map<Integer, String>) testType(intStringMap);
        assertEquals(intStringMap.size(), resultIntStringMap.size());
        for (Integer key : intStringMap.keySet()) {
            assertEquals(intStringMap.get(key), resultIntStringMap.get(key));
        }
    }
    
    @Test
    @DisplayName("Test Null Values")
    void testNullValues() throws Exception {
        // Test null values
        assertNull(testType(null));
        
        // Test null collections
        List<Integer> nullList = null;
        assertNull(testType(nullList));
        
        Map<String, Integer> nullMap = null;
        assertNull(testType(nullMap));
    }
    
    @Test
    @DisplayName("Test Base64 Encoding")
    void testBase64() throws Exception {
        String testValue = "Hello, Base64!";
        
        // Serialize
        SerializationObject encoder = SerializationObject.newEncode();
        Exception err = encoder.add(testValue);
        assertNull(err);
        
        // Get Base64 string
        String base64Str = encoder.toBase64();
        assertNotNull(base64Str);
        assertFalse(base64Str.isEmpty());
        
        // Decode from Base64
        byte[] data = SerializationObject.fromBase64(base64Str);
        assertNotNull(data);
        assertTrue(data.length > 0);
        
        // Deserialize
        SerializationObject decoder = SerializationObject.newDecode(data, 0, registry);
        Object result = decoder.get();
        
        assertEquals(testValue, result);
    }
    
    @Test
    @DisplayName("Test Utility Functions")
    void testUtilityFunctions() throws Exception {
        Integer testValue = 12345;
        
        // Test dataOf
        byte[] data = SerializationObject.dataOf(testValue);
        assertNotNull(data);
        assertTrue(data.length > 0);
        
        // Test elementOf
        Object result = SerializationObject.elementOf(data, registry);
        assertEquals(testValue, result);
        
        // Test with null
        assertNull(SerializationObject.dataOf(null));
        assertNull(SerializationObject.elementOf(null, registry));
    }
    
    @Test
    @DisplayName("Test Large Data Serialization")
    void testLargeData() throws Exception {
        // Create large list
        List<Integer> largeList = new ArrayList<>();
        for (int i = 0; i < 1000; i++) { // Smaller than Go test for faster execution
            largeList.add(i);
        }
        
        // Serialize
        SerializationObject encoder = SerializationObject.newEncode();
        Exception err = encoder.add(largeList);
        assertNull(err);
        
        byte[] data = encoder.getData();
        assertTrue(data.length > 0);
        
        // Deserialize
        SerializationObject decoder = SerializationObject.newDecode(data, 0, registry);
        @SuppressWarnings("unchecked")
        List<Integer> result = (List<Integer>) decoder.get();
        
        assertNotNull(result);
        assertEquals(largeList.size(), result.size());
        
        // Spot check some values
        for (int i = 0; i < largeList.size(); i += 100) {
            assertEquals(largeList.get(i), result.get(i));
        }
    }
    
    @Test
    @DisplayName("Test Buffer Expansion")
    void testBufferExpansion() throws Exception {
        SerializationObject encoder = SerializationObject.newEncode();
        
        // Add many small items to force buffer expansion
        for (int i = 0; i < 100; i++) { // Smaller for quick test
            Exception err = encoder.add(i);
            assertNull(err, "Failed to add item " + i);
        }
        
        byte[] data = encoder.getData();
        assertTrue(data.length > 0);
        
        // Verify we can still serialize properly
        SerializationObject finalEncoder = SerializationObject.newEncode();
        Exception err = finalEncoder.add("test after expansion");
        assertNull(err);
    }
    
    @Test
    @DisplayName("Test Error Handling")
    void testErrorHandling() throws Exception {
        // Test invalid Base64
        assertThrows(Exception.class, () -> {
            SerializationObject.fromBase64("invalid-base64!");
        });
    }
    
    @Test
    @DisplayName("Test Elements Container")
    void testElements() throws Exception {
        // Test creating Elements from list
        List<Integer> list = Arrays.asList(1, 2, 3);
        Elements elements = Elements.create(null, list);
        
        assertFalse(elements.isEmpty());
        assertEquals(3, elements.size());
        
        List<Object> values = elements.getElements();
        assertEquals(3, values.size());
        assertEquals(1, values.get(0));
        assertEquals(2, values.get(1));
        assertEquals(3, values.get(2));
        
        // Test creating Elements from map
        Map<String, Integer> map = new HashMap<>();
        map.put("one", 1);
        map.put("two", 2);
        Elements mapElements = Elements.create(null, map);
        
        assertEquals(2, mapElements.size());
        List<Object> keys = mapElements.getKeys();
        assertTrue(keys.contains("one"));
        assertTrue(keys.contains("two"));
        
        // Test error Elements
        Elements errorElements = Elements.createError("test error");
        assertNotNull(errorElements.getError());
        assertEquals("test error", errorElements.getError().getMessage());
        
        // Test serialization/deserialization
        byte[] serialized = elements.serialize();
        assertNotNull(serialized);
        
        Elements deserialized = new Elements();
        deserialized.deserialize(serialized, registry);
        
        assertEquals(elements.size(), deserialized.size());
        assertEquals(elements.getElements(), deserialized.getElements());
    }
}