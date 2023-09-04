package main

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"os"
	"reflect"

	"github.com/BurntSushi/toml"
)

// parseConfig takes the config content as a string and attempts to parse it
// into the conf struct based on the TOML spec.
func parseConfig(configContent string) {
	if configContent == "" {
		fail("Configuration is empty!")
		os.Exit(1)
	}
	md, err := toml.Decode(configContent, &conf)
	if err != nil {
		fail("Error decoding TOML: " + err.Error())
		os.Exit(1)
	}
	if verboseEnabled {
		for _, undecoded := range md.Undecoded() {
			warn("Undecoded scoring configuration key \"" + undecoded.String() + "\" will not be used.")
		}
	}

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
	}

	// Check if the config version matches ours.
	if conf.Version != version {
		warn("Scoring version does not match Aeacus version! Compatibility issues may occur.")
		info("Consider updating your config to include:")
		info("    version = '" + version + "'")
	}

	// Print warnings for impossible checks and undefined check types.
	for i, check := range conf.Check {
		if len(check.Pass) == 0 && len(check.PassOverride) == 0 {
			warn("Check " + fmt.Sprintf("%d", i+1) + " does not define any possible ways to pass!")
		}
		allConditions := append(append(append([]cond{}, check.Pass[:]...), check.Fail[:]...), check.PassOverride[:]...)
		for j, cond := range allConditions {
			if cond.Type == "" {
				warn("Check " + fmt.Sprintf("%d condition %d", i+1, j+1) + " does not have a check type!")
			}
		}
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
	for i, check := range conf.Check {
		green("CHCK", fmt.Sprintf("Check %d (%d points):", i+1, check.Points))
		fmt.Println("Message:", check.Message)
		for _, c := range check.Pass {
			fmt.Println("Pass Condition:")
			fmt.Print(c)
		}
		for _, c := range check.PassOverride {
			fmt.Println("PassOverride Condition:")
			fmt.Print(c)
		}
		for _, c := range check.Fail {
			fmt.Println("Fail Condition:")
			fmt.Print(c)
		}
	}
}

func obfuscateConfig() {
	if debugEnabled {
		debug("Obfuscating configuration...")
	}
	if err := obfuscateData(&conf.Password); err != nil {
		fail(err.Error())
	}
	for i, check := range conf.Check {
		if err := obfuscateData(&conf.Check[i].Message); err != nil {
			fail(err.Error())
		}
		if conf.Check[i].Hint != "" {
			if err := obfuscateData(&conf.Check[i].Hint); err != nil {
				fail(err.Error())
			}
		}
		for j := range check.Pass {
			if err := obfuscateCond(&conf.Check[i].Pass[j]); err != nil {
				fail(err.Error())
			}
		}
		for j := range check.PassOverride {
			if err := obfuscateCond(&conf.Check[i].PassOverride[j]); err != nil {
				fail(err.Error())
			}
		}
		for j := range check.Fail {
			if err := obfuscateCond(&conf.Check[i].Fail[j]); err != nil {
				fail(err.Error())
			}
		}
	}
}

// obfuscateCond is a convenience function to obfuscate all string fields of a
// struct using reflection. It assumes all struct fields are strings.
func obfuscateCond(c *cond) error {
	s := reflect.ValueOf(c).Elem()
	for i := 0; i < s.NumField(); i++ {
		if s.Type().Field(i).Name == "regex" {
			continue
		}
		datum := s.Field(i).String()
		if err := obfuscateData(&datum); err != nil {
			return err
		}
		s.Field(i).SetString(datum)
	}
	return nil
}

// deobfuscateCond is a convenience function to deobfuscate all string fields
// of a struct using reflection.
func deobfuscateCond(c *cond) error {
	s := reflect.ValueOf(c).Elem()
	for i := 0; i < s.NumField(); i++ {
		if s.Type().Field(i).Name == "regex" {
			continue
		}
		datum := s.Field(i).String()
		if err := deobfuscateData(&datum); err != nil {
			return err
		}
		s.Field(i).SetString(datum)
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
