package cmd

import "testing"

func TestEncryption(t *testing.T) {
	// Testing encryptConfig
	plainText := "Test string."
	if encrypted, err := encryptConfig(plainText); err != nil {
		t.Errorf("encryptConfig returned an error: %s", err.Error())
	} else {
		if decrypted, err := decryptConfig(encrypted); decrypted != plainText {
			t.Errorf("decryptConfig(encryptConfig('%s')) == %s, should be '%s'", plainText, decrypted, plainText)
		} else if err != nil {
			t.Errorf("decryptConfig returned an error: %s", err.Error())
		}
	}

	// Testing encryptString
	password := "Password1!"
	encrypted := encryptString(password, plainText)
	if decrypted := decryptString(password, encrypted); decrypted != plainText {
		t.Errorf("decryptConfig(encryptConfig('%s')) == %s, should be '%s'", plainText, decrypted, plainText)
	}
}
