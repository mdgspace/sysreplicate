package backup

import (
    "crypto/aes"
    "crypto/cipher"
    "crypto/rand"
    "encoding/base64"
    "fmt"
    "io"
    "os"
)

//simplified encryption config without password
type EncryptionConfig struct {
    Key []byte // Direct 32-byte key instead of password+salt
}

//AES-GCM encryption with direct key (no password derivation)
func EncryptFile(filePath string, config *EncryptionConfig) (string, error) {
    data, err := os.ReadFile(filePath)
    if err != nil {
        return "", fmt.Errorf("failed to read file %s: %w", filePath, err)
    }

    //use direct key (no password derivation)
    block, err := aes.NewCipher(config.Key)
    if err != nil {
        return "", fmt.Errorf("failed to create cipher: %w", err)
    }

    // GCM mode from AES block
    gcm, err := cipher.NewGCM(block)
    if err != nil {
        return "", fmt.Errorf("failed to create GCM: %w", err)
    }

    // Generate nonce
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

//generate a random 32-byte key for AES-256
func GenerateKey() ([]byte, error) {
    key := make([]byte, 32) // 32 bytes for AES-256
    _, err := rand.Read(key)
    return key, err
}
