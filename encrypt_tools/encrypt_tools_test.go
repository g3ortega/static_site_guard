package encrypt_tools

import "testing"

func TestEncryptDecrypt(t *testing.T) {
	key := "0123456789abcdef" // 16-byte key
	original := "hello world"
	encrypted := Encrypt(key, original)
	if encrypted == original {
		t.Fatalf("expected encrypted text to differ from original")
	}
	decrypted := Decrypt(key, encrypted)
	if decrypted != original {
		t.Fatalf("expected decrypted text to be %q, got %q", original, decrypted)
	}
}
