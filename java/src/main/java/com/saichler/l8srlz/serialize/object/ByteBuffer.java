package com.saichler.l8srlz.serialize.object;

import java.util.Arrays;

/**
 * Dynamic byte buffer that grows as needed
 */
public class ByteBuffer {
    private byte[] data;
    
    public ByteBuffer() {
        this.data = new byte[1024]; // Initial size like Go implementation
    }
    
    public ByteBuffer(byte[] data) {
        this.data = data;
    }
    
    /**
     * Ensure the buffer has enough space
     * @param location current location
     * @param need additional bytes needed
     */
    public void checkAndEnlarge(int location, int need) {
        if (location + need > data.length) {
            // Exponential growth with minimum threshold (from performance optimizations)
            int newCapacity = data.length * 2;
            if (newCapacity < location + need + 512) {
                newCapacity = location + need + 512;
            }
            data = Arrays.copyOf(data, newCapacity);
        }
    }
    
    /**
     * Get the underlying byte array
     * @return byte array
     */
    public byte[] getData() {
        return data;
    }
    
    /**
     * Set the underlying byte array
     * @param data new byte array
     */
    public void setData(byte[] data) {
        this.data = data;
    }
    
    /**
     * Get the length of the buffer
     * @return buffer length
     */
    public int length() {
        return data.length;
    }
    
    /**
     * Get byte at specific position
     * @param index position
     * @return byte value
     */
    public byte get(int index) {
        return data[index];
    }
    
    /**
     * Set byte at specific position
     * @param index position
     * @param value byte value
     */
    public void set(int index, byte value) {
        data[index] = value;
    }
    
    /**
     * Copy data into the buffer at specified location
     * @param location starting position
     * @param src source data
     * @param srcOffset source offset
     * @param length number of bytes to copy
     */
    public void copyFrom(int location, byte[] src, int srcOffset, int length) {
        System.arraycopy(src, srcOffset, data, location, length);
    }
    
    /**
     * Copy data from the buffer to destination
     * @param location starting position in buffer
     * @param dest destination array
     * @param destOffset destination offset
     * @param length number of bytes to copy
     */
    public void copyTo(int location, byte[] dest, int destOffset, int length) {
        System.arraycopy(data, location, dest, destOffset, length);
    }
    
    /**
     * Create a copy of a portion of the buffer
     * @param start start position
     * @param end end position
     * @return new byte array containing the copy
     */
    public byte[] copyRange(int start, int end) {
        return Arrays.copyOfRange(data, start, end);
    }
}