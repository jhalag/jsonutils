# JSON Helper Utilities

This repository provides helper utilities for working with JSON in Go, with a focus on secure data handling through encryption and base64 encoding.

## Overview

As I keep writing boilerplate for basic everyday JSON marshaling, I decided to split it out into a reusable module. This library offers convenient tools for handling JSON data while maintaining security through encryption capabilities and transparent base64 encoding.

The library currently provides two main types:
- **Base64 String**: Automatically encodes/decodes strings as base64 when marshaling/unmarshaling JSON
- **Encrypted String**: Encrypts sensitive fields in your structs when marshaling to JSON, while still being able to unmarshal them back correctly

## Features

- **Secure JSON Encryption**: Encrypt sensitive fields in your structs when marshaling to JSON
- **Automatic Base64 Encoding**: Transparently encode/decode strings as base64 during JSON operations
- **Automatic Decryption**: Transparent decryption when unmarshaling JSON data
- **Backward Compatibility**: Can handle both encrypted and plaintext data during unmarshaling
- **Easy Integration**: Simple API that integrates seamlessly with Go's standard `json` package

## Usage Examples

### Base64 String Usage

```go
package main

import (
    "encoding/json"
    "fmt"
    
    base64 "github.com/jhalag/jsonutils/base64"
)

type User struct {
    Name  string
    Email string
    Data  base64.String // This field will be automatically base64 encoded
}

func main() {
    user := User{
        Name:  "John Doe",
        Email: "john@example.com",
        Data:  "data that should be base64 encoded",
    }
    
    // Marshal to JSON (data will be base64 encoded)
    jsonData, err := json.MarshalIndent(user, "", "  ")
    if err != nil {
        panic(err)
    }
    
    fmt.Println(string(jsonData))
    // Output will contain the base64 encoded data
    
    // Unmarshal back from JSON
    var decodedUser User
    err = json.Unmarshal(jsonData, &decodedUser)
    if err != nil {
        panic(err)
    }
    
    fmt.Println("Decoded data:", decodedUser.Data)
}
```

### Encrypted String Usage

```go
package main

import (
    "encoding/json"
    "fmt"
    
    encryption "github.com/jhalag/jsonutils/encryption"
)

type User struct {
    Name     string
    Email    string
    Password encryption.String // This field will be encrypted
}

func main() {
    // Set the encryption key (required before any marshaling/unmarshaling)
    encryption.SetMarshalEncryptionKey("my-secret-key")
    
    // Create a user with sensitive data
    user := User{
        Name:     "John Doe",
        Email:    "john@example.com",
        Password: encryption.String{Value: "secret-password"},
    }
    
    // Marshal to JSON (password will be encrypted)
    jsonData, err := json.MarshalIndent(user, "", "  ")
    if err != nil {
        panic(err)
    }
    
    fmt.Println(string(jsonData))
    // Output will contain the encrypted password, not the plaintext
    
    // Unmarshal back from JSON
    var decodedUser User
    err = json.Unmarshal(jsonData, &decodedUser)
    if err != nil {
        panic(err)
    }
    
    fmt.Println("Decrypted password:", decodedUser.Password.Value)
}
```

## How It Works

### Base64 String
1. **Encoding**: When marshaling a struct containing `base64.String` fields, the library automatically base64-encodes the values
2. **Storage**: Base64-encoded data is stored as strings in JSON format
3. **Decoding**: When unmarshaling, the library automatically decodes base64 data back to original strings

### Encrypted String
1. **Encryption**: When marshaling a struct containing `encryption.String` fields, the library automatically encrypts the values using AES-256-GCM encryption
2. **Storage**: Encrypted data is stored as base64-encoded strings in JSON format
3. **Decryption**: When unmarshaling, the library automatically decrypts encrypted fields and handles plaintext fields transparently

## Security Notes

- The encryption key must be set before any marshaling/unmarshaling operations with encrypted fields
- This implementation uses AES-256-GCM with a randomly generated nonce for each encryption
- The library supports both encrypted and plaintext data during unmarshaling for backward compatibility

## Installation

```bash
go get github.com/jhalag/jsonutils
```

## License

MIT License - see the LICENSE file for details.