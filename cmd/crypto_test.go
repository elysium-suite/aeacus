package cmd

import "testing"

func TestEncryption(t *testing.T) {
	// Testing encryptConfig
	plainText := "Test string."
	encrypted := encryptConfig(plainText)
	if decrypted := decryptConfig(encrypted); decrypted != plainText {
		t.Errorf("decryptConfig(encryptConfig('%s')) == %s, should be '%s'", plainText, decrypted, plainText)
	}

	// Testing encryptString
	password := "Password1!"
	encrypted = encryptString(password, plainText)
	if decrypted := decryptString(password, encrypted); decrypted != plainText {
		t.Errorf("decryptConfig(encryptConfig('%s')) == %s, should be '%s'", plainText, decrypted, plainText)
	}
}
