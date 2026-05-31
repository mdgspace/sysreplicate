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

	"github.com/mdgspace/sysreplicate/internal/domain"
	"golang.org/x/crypto/argon2"
)

type EncryptionConfig struct {
	Key []byte
}

func EncryptFile(filePath string, config *EncryptionConfig) (string, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to read file %s: %w", filePath, err)
	}

	block, err := aes.NewCipher(config.Key)
	if err != nil {
		return "", fmt.Errorf("failed to create cipher: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("failed to create GCM: %w", err)
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", fmt.Errorf("failed to generate nonce: %w", err)
	}

	ciphertext := gcm.Seal(nonce, nonce, data, nil)
	encoded := base64.StdEncoding.EncodeToString(ciphertext)
	return encoded, nil
}

const (
	argonTime    uint32 = 3
	argonMemory  uint32 = 64 * 1024
	argonThreads uint8  = 4
	saltLen             = 32
)

func DeriveKey(passphrase string, params *domain.KeyDerivationParams) []byte {
	return argon2.IDKey([]byte(passphrase), params.Salt, params.TimeCost, params.MemoryCost, params.Threads, 32)
}

func NewKeyDerivationParams() (*domain.KeyDerivationParams, error) {
	salt := make([]byte, saltLen)
	if _, err := rand.Read(salt); err != nil {
		return nil, fmt.Errorf("failed to generate salt: %w", err)
	}
	return &domain.KeyDerivationParams{
		Salt:       salt,
		TimeCost:   argonTime,
		MemoryCost: argonMemory,
		Threads:    argonThreads,
	}, nil
}

func HashPayload(data []byte) (string, error) {
	h := sha256.New()
	if _, err := h.Write(data); err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", h.Sum(nil)), nil
}

func VerifyHash(data []byte, expected string) bool {
	return fmt.Sprintf("%x", sha256.Sum256(data)) == expected
}
