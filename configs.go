package main

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"os"
	"reflect"

	"github.com/BurntSushi/toml"
	"github.com/safinsingh/aeaconf2"
	"github.com/safinsingh/aeaconf2/compat"
)

// parseConfig takes the config content as a string and attempts to parse it
// into the conf struct based on the TOML spec.
func parseConfig(configContent string) {
	if configContent == "" {
		fail("Configuration is empty!")
		os.Exit(1)
	}

	headerRaw, checksRaw, err := compat.SeparateConfig([]byte(configContent))
	if err != nil {
		fail("error separating config file: " + err.Error())
		os.Exit(1)
	}

	cfg := new(config)
	err = toml.Unmarshal(headerRaw, cfg)
	if err != nil {
		fail("error parsing config file header: " + err.Error())
		os.Exit(1)
	}

	ab := aeaconf2.DefaultAeaconfBuilder(checksRaw, funcRegistry).
		SetLineOffset(countLines(headerRaw)).
		SetMaxPoints(cfg.MaxPoints)

	cfg.Checks = ab.GetChecks()

	// If there's no remote, local must be enabled.
	if conf.Remote == "" {
		conf.Local = true
		if conf.DisableRemoteEncryption {
			fail("Remote encryption cannot be disabled if remote is not enabled!")
			os.Exit(1)
		}
	} else {
		if conf.Remote[len(conf.Remote)-1] == '/' {
			fail("Your remote URL must not end with a slash: try", conf.Remote[:len(conf.Remote)-1])
			os.Exit(1)
		}
		if conf.Name == "" {
			fail("Need image name in config if remote is enabled.")
			os.Exit(1)
		}

		if conf.Password == "" && conf.DisableRemoteEncryption == false {
			fail("Need password in config if remote is enabled.")
			os.Exit(1)
		}

		if conf.DisableRemoteEncryption && conf.Password != "" {
			warn("Remote encryption is disabled, but a password is still defined!")
		}

	}

	// Check if the config version matches ours.
	if conf.Version != version {
		warn("Scoring version does not match Aeacus version! Compatibility issues may occur.")
		info("Consider updating your config to include:")
		info("    version = '" + version + "'")
	}
}

// writeConfig writes the in-memory config to disk as the an encrypted
// configuration file.
func writeConfig() {
	buf := new(bytes.Buffer)
	if err := toml.NewEncoder(buf).Encode(conf); err != nil {
		fail(err.Error())
		os.Exit(1)
		return
	}

	dataPath := dirPath + scoringData
	encryptedConfig, err := encryptConfig(buf.String())
	if err != nil {
		fail("Encrypting config failed: " + err.Error())
		os.Exit(1)
	} else if verboseEnabled {
		info("Writing data to " + dataPath + "...")
	}

	writeFile(dataPath, encryptedConfig)
}

// ReadConfig parses the scoring configuration file.
func readConfig() {
	fileContent, err := readFile(dirPath + scoringConf)
	if err != nil {
		fail("Configuration file (" + dirPath + scoringConf + ") not found!")
		os.Exit(1)
	}
	parseConfig(fileContent)
	assignPoints()
	assignDescriptions()
	if verboseEnabled {
		printConfig()
	}
	obfuscateConfig()
}

// PrintConfig offers a printed representation of the config, as parsed
// by readData and parseConfig.
func printConfig() {
	pass("Configuration " + dirPath + scoringConf + " validity check passed!")
	blue("CONF", scoringConf)
	if conf.Version != "" {
		pass("Version:", conf.Version)
	}
	if conf.Title == "" {
		red("MISS", "Title:", "N/A")
	} else {
		pass("Title:", conf.Title)
	}
	if conf.Name == "" {
		red("MISS", "Name:", "N/A")
	} else {
		pass("Name:", conf.Name)
	}
	if conf.OS == "" {
		red("MISS", "OS:", "N/A")
	} else {
		pass("OS:", conf.OS)
	}
	if conf.User == "" {
		red("MISS", "User:", "N/A")
	} else {
		pass("User:", conf.User)
	}
	if conf.Remote != "" {
		pass("Remote:", conf.Remote)
	}
	if conf.DisableRemoteEncryption {
		pass("Remote Encryption:", "Disabled")
	} else {
		pass("Remote Encryption:", "Enabled")
	}
	if conf.Local {
		pass("Local:", conf.Local)
	}
	if conf.EndDate != "" {
		pass("End Date:", conf.EndDate)
	}
	for i, check := range conf.Checks {
		green("CHCK", fmt.Sprintf("Check %d: %s", i+1, check.Debug()))
	}
}

func obfuscateConfig() {
	if debugEnabled {
		debug("Obfuscating configuration...")
	}
	if err := obfuscateData(&conf.Password); err != nil {
		fail(err.Error())
	}
	for i := range conf.Checks {
		if err := obfuscateData(&conf.Checks[i].Message); err != nil {
			fail(err.Error())
		}
		if conf.Checks[i].Hint != "" {
			if err := obfuscateData(&conf.Checks[i].Hint); err != nil {
				fail(err.Error())
			}
		}
		if err := obfuscateCond(conf.Checks[i].Condition); err != nil {
			fail(err.Error())
		}
	}
}

// obfuscateCond is a convenience function to obfuscate all string fields of a
// struct using reflection. It assumes all struct fields are strings.

// ummmmmm
func obfuscateCond(c aeaconf2.Condition) error {
	s := reflect.ValueOf(c).Elem()
	t := s.Type()

	for i := 0; i < s.NumField(); i++ {
		field := s.Field(i)
		fieldType := t.Field(i).Type

		if fieldType.Kind() == reflect.String {
			datum := field.String()
			if err := obfuscateData(&datum); err != nil {
				return err
			}
			field.SetString(datum)
		} else if fieldType.Implements(reflect.TypeOf((*aeaconf2.Condition)(nil)).Elem()) {
			if err := obfuscateCond(field.Addr().Interface().(aeaconf2.Condition)); err != nil {
				return err
			}
		}
	}

	return nil
}

// deobfuscateCond is a convenience function to deobfuscate all string fields
// of a struct using reflection.
func deobfuscateCond(c aeaconf2.Condition) error {
	s := reflect.ValueOf(c).Elem()
	t := s.Type()

	for i := 0; i < s.NumField(); i++ {
		field := s.Field(i)
		fieldType := t.Field(i).Type

		if fieldType.Kind() == reflect.String {
			datum := field.String()
			if err := deobfuscateData(&datum); err != nil {
				return err
			}
			field.SetString(datum)
		} else if fieldType.Implements(reflect.TypeOf((*aeaconf2.Condition)(nil)).Elem()) {
			if err := deobfuscateCond(field.Addr().Interface().(aeaconf2.Condition)); err != nil {
				return err
			}
		}
	}

	return nil
}

func xor(key, plaintext string) string {
	ciphertext := make([]byte, len(plaintext))
	for i := 0; i < len(plaintext); i++ {
		ciphertext[i] = key[i%len(key)] ^ plaintext[i]
	}
	return string(ciphertext)
}

func hexEncode(inputString string) string {
	return hex.EncodeToString([]byte(inputString))
}

func hexDecode(inputString string) (string, error) {
	result, err := hex.DecodeString(inputString)
	if err != nil {
		return "", err
	}
	return string(result), nil
}
