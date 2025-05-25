package crypto

import (
	"bytes"
	"testing"
)

func TestGenerateSalt(t *testing.T) {
	salt1, err := GenerateSalt()
	if err != nil {
		t.Fatalf("Failed to generate salt: %v", err)
	}

	if len(salt1) != saltLength {
		t.Fatalf("Expected salt length %d, got %d", saltLength, len(salt1))
	}

	// Generate another salt to ensure they're different
	salt2, err := GenerateSalt()
	if err != nil {
		t.Fatalf("Failed to generate second salt: %v", err)
	}

	if bytes.Equal(salt1, salt2) {
		t.Fatal("Generated salts should be different")
	}
}

func TestDeriveKey(t *testing.T) {
	password := "test-password"
	salt := []byte("test-salt-32-bytes-long-exactly!!")

	key1 := DeriveKey(password, salt)
	if len(key1) != keyLength {
		t.Fatalf("Expected key length %d, got %d", keyLength, len(key1))
	}

	// Same password and salt should produce same key
	key2 := DeriveKey(password, salt)
	if !bytes.Equal(key1, key2) {
		t.Fatal("Same password and salt should produce same key")
	}

	// Different password should produce different key
	key3 := DeriveKey("different-password", salt)
	if bytes.Equal(key1, key3) {
		t.Fatal("Different passwords should produce different keys")
	}

	// Different salt should produce different key
	differentSalt := []byte("different-salt-32-bytes-long-ex!!")
	key4 := DeriveKey(password, differentSalt)
	if bytes.Equal(key1, key4) {
		t.Fatal("Different salts should produce different keys")
	}
}

func TestEncryptDecrypt(t *testing.T) {
	key := []byte("test-key-32-bytes-long-exactly!!")
	plaintext := []byte("Hello, World! This is a test message.")

	// Test encryption
	ciphertext, err := Encrypt(plaintext, key)
	if err != nil {
		t.Fatalf("Failed to encrypt: %v", err)
	}

	if len(ciphertext) <= len(plaintext) {
		t.Fatal("Ciphertext should be longer than plaintext (includes nonce)")
	}

	// Test decryption
	decrypted, err := Decrypt(ciphertext, key)
	if err != nil {
		t.Fatalf("Failed to decrypt: %v", err)
	}

	if !bytes.Equal(plaintext, decrypted) {
		t.Fatalf("Decrypted text doesn't match original. Expected: %s, Got: %s",
			string(plaintext), string(decrypted))
	}
}

func TestEncryptDecryptDifferentKeys(t *testing.T) {
	plaintext := []byte("Secret message")
	key1 := []byte("key1-32-bytes-long-exactly-here!")
	key2 := []byte("key2-32-bytes-long-exactly-here!")

	// Encrypt with key1
	ciphertext, err := Encrypt(plaintext, key1)
	if err != nil {
		t.Fatalf("Failed to encrypt: %v", err)
	}

	// Try to decrypt with key2 (should fail)
	_, err = Decrypt(ciphertext, key2)
	if err == nil {
		t.Fatal("Decryption with wrong key should fail")
	}
}

func TestEncryptDecryptEmptyData(t *testing.T) {
	key := []byte("test-key-32-bytes-long-exactly!!")
	plaintext := []byte("")

	ciphertext, err := Encrypt(plaintext, key)
	if err != nil {
		t.Fatalf("Failed to encrypt empty data: %v", err)
	}

	decrypted, err := Decrypt(ciphertext, key)
	if err != nil {
		t.Fatalf("Failed to decrypt empty data: %v", err)
	}

	if !bytes.Equal(plaintext, decrypted) {
		t.Fatal("Decrypted empty data doesn't match original")
	}
}

func TestDecryptInvalidData(t *testing.T) {
	key := []byte("test-key-32-bytes-long-exactly!!")

	// Test with data too short
	_, err := Decrypt([]byte("short"), key)
	if err == nil {
		t.Fatal("Decryption of too-short data should fail")
	}

	// Test with corrupted data
	validCiphertext, _ := Encrypt([]byte("test"), key)
	corruptedCiphertext := make([]byte, len(validCiphertext))
	copy(corruptedCiphertext, validCiphertext)
	corruptedCiphertext[len(corruptedCiphertext)-1] ^= 0x01 // Flip last bit

	_, err = Decrypt(corruptedCiphertext, key)
	if err == nil {
		t.Fatal("Decryption of corrupted data should fail")
	}
}
