package main

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

func TestObfuscation(t *testing.T) {
	// Testing obfuscateData
	plainText := "I am data!"
	cipherText := "I am data!"
	if err := obfuscateData(&cipherText); err != nil {
		t.Errorf("obfuscateData returned an error: %s", err.Error())
	} else {
		if err := deobfuscateData(&cipherText); cipherText != plainText {
			t.Errorf("deobfuscateData(obfuscateData(%s)) == %s, should be '%s'", plainText, cipherText, plainText)
		} else if err != nil {
			t.Errorf("deobufscateData returned an error: %s", err.Error())
		}
	}
}
