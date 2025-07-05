package backup

import (
    "crypto/aes"
    "crypto/cipher"
    "crypto/rand"
    "crypto/sha256"
    "encoding/base64"
    "fmt"
    "io"
    "os"
)

// encryption formatt
type EncryptionConfig struct {
    
	Password string
    Salt     []byte
}

// AES-GCM encrpption with password based derivation
func EncryptFile(filePath string, config *EncryptionConfig) (string, error) {
    
    data, err := os.ReadFile(filePath)
    
	if err != nil {
        	
		return "", fmt.Errorf("failed to read file %s: %w", filePath, err)
    
	}

    // derive key from password+salt
    key := deriveKey(config.Password, config.Salt)

    // AES cipher from key
    block, err := aes.NewCipher(key)
    if err != nil {
    
		return "", fmt.Errorf("failed to create cipher: %w", err)
    
	}

    // GCM mode from AES block
    gcm, err := cipher.NewGCM(block)
    if err != nil {
    
		return "", fmt.Errorf("failed to create GCM: %w", err)
    
	}

    //generate nonce
    nonce := make([]byte, gcm.NonceSize())
    if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
    
		return "", fmt.Errorf("failed to generate nonce: %w", err)
    
	}

    //encrypt data with nonce
    ciphertext := gcm.Seal(nonce, nonce, data, nil)

    //encode to base64
    encoded := base64.StdEncoding.EncodeToString(ciphertext)
    return encoded, nil
}

// // decryption methods - will be used later
// func DecryptFile(encryptedData string, config *EncryptionConfig) ([]byte, error) {
//     ///decode from base64
//     ciphertext, err := base64.StdEncoding.DecodeString(encryptedData)
//     if err != nil {
    
// 		return nil, fmt.Errorf("failed to decode base64: %w", err)
    
// 	}

//     //derive key from password
//     key := deriveKey(config.Password, config.Salt)

//     //create AES cipher
//     block, err := aes.NewCipher(key)
//     if err != nil {
    
// 		return nil, fmt.Errorf("failed to create cipher: %w", err)
    
// 	}

//     //create GCM mode
//     gcm, err := cipher.NewGCM(block)
//     if err != nil {
    
// 		return nil, fmt.Errorf("failed to create GCM: %w", err)
    
// 	}

//     //extract nonce
//     nonceSize := gcm.NonceSize()
//     if len(ciphertext) < nonceSize {
    
// 		return nil, fmt.Errorf("ciphertext too short")
    
// 	}

//     nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]

//     //decrypt data
//     plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
//     if err != nil {
    
// 		return nil, fmt.Errorf("failed to decrypt: %w", err)
    
// 	}

//     return plaintext, nil
// }



// derive a 32-byte key from password and salt using SHA-256
func deriveKey(password string, salt []byte) []byte {
    hash := sha256.New()
    hash.Write([]byte(password))
    hash.Write(salt)
    return hash.Sum(nil)
}

//create a random salt for key derivation
func GenerateSalt() ([]byte, error) {
    salt := make([]byte, 32)
    _, err := rand.Read(salt)
    return salt, err
}
