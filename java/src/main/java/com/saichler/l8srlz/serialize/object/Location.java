package com.saichler.l8srlz.serialize.object;

/**
 * Wrapper for location pointer to mimic Go's pointer behavior
 */
public class Location {
    private int value;
    
    public Location() {
        this.value = 0;
    }
    
    public Location(int value) {
        this.value = value;
    }
    
    public int getValue() {
        return value;
    }
    
    public void setValue(int value) {
        this.value = value;
    }
    
    public void increment(int amount) {
        this.value += amount;
    }
    
    public void increment() {
        this.value++;
    }
    
    @Override
    public String toString() {
        return "Location{" + value + "}";
    }
}