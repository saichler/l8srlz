# object

Simple, deep, struct serializing method.

## Overview
Missing from Protocol Buffers is the Object notion. This is understandable as 
some of the languages are not object-oriented and are missing the Type Registry 
that exist in languages like Java, where you can instantiate a type just by its 
name using a ClassLoader.

This package is coming to overcome this challenge and introduce the object notion via 
utilizing the **Type Registry** inside the **shared** project.

## Prolog
This package will be used by the model introspector to generate **Delta Notifications** 
for models in a use-case of **Stateful Microservices**.


