package main

import (
	// Standard imports
	"fmt"
	"log"

	// Files, bytes, and crypto
	"bufio"
	"bytes"
	"io/ioutil"
	"compress/zlib"

    // Checks
	"os"
    "os/exec"
    "crypto/sha1"
    "encoding/hex"

	// Github/External
	"github.com/BurntSushi/toml"
	"github.com/urfave/cli"
)

//////////////////////////////////////////////////////////////////
//  .oooo.    .ooooo.   .oooo.    .ooooo.  oooo  oooo   .oooo.o //
// `P  )88b  d88' `88b `P  )88b  d88' `"Y8 `888  `888  d88  "8  //
//  .oP"888  888ooo888  .oP"888  888        888   888  `"Y88b.  //
// d8(  888  888    .o d8(  888  888   .o8  888   888  o.  )88b //
// `Y888""8o `Y8bod8P' `Y888""8o `Y8bod8P'  `V88V"V8P' 8""888P' //
//////////////////////////////////////////////////////////////////

func main() {
	configName := "/opt/aeacus/scoring.conf"
	dataName := "/opt/aeacus/scoring.dat"
	var scoreType string
	id := imageData{0, 0, []scoreItem{}}
	app := &cli.App{
		UseShortOptionHandling: true,
		EnableBashCompletion:   true,
		Name:                   "aeacus",
		Usage:                  "setup and score vulnerabilities in an image",
		Action: func(c *cli.Context) error {
			mc := metaConfig{c, configName, dataName, scoringChecks{}}
			readConfig(&mc, readFile(mc.ConfigName))
			if c.Bool("v") {
				printConfig(&mc)
			}
			return nil
		},
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:    "verbose",
				Aliases: []string{"v"},
				Usage:   "Print extra information",
			},
		},
		Commands: []*cli.Command{
			{
				Name:    "check",
				Aliases: []string{"c"},
				Usage:   "(default) Check that the scoring config is valid",
				Action: func(c *cli.Context) error {
					mc := metaConfig{c, configName, dataName, scoringChecks{}}
					readConfig(&mc, readFile(mc.ConfigName))
					if c.Bool("v") {
						printConfig(&mc)
					}
					return nil
				},
			},
			{
				Name:    "score",
				Aliases: []string{"s"},
				Usage:   "Score image with current config",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:        "type",
						Aliases:     []string{"t"},
						Usage:       "Type of score output (console, html)",
						Destination: &scoreType,
					},
				},
				Action: func(c *cli.Context) error {
					var config scoringChecks
					mc := metaConfig{c, configName, dataName, config}
					readConfig(&mc, readFile(mc.ConfigName))
					scoreImage(&mc, &id, "stdout")
					return nil
				},
			},
			{
				Name:    "decrypt",
				Aliases: []string{"d"},
				Usage:   "Decrypt scoring.dat and print to terminal",
				Action: func(c *cli.Context) error {
					mc := metaConfig{c, configName, dataName, scoringChecks{}}
					decryptedData := readData(&mc)
					fmt.Println(decryptedData)
					return nil
				},
			},
			{
				Name:    "release",
				Aliases: []string{"r"},
				Usage:   "Encrypt scoring.conf, clean up image for release",
				Action: func(c *cli.Context) error {
					var config scoringChecks
					mc := metaConfig{c, configName, dataName, config}
					readConfig(&mc, readFile(mc.ConfigName))
					releaseImage(&mc)
					return nil
				},
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

/////////////////////
// DATA STRUCTURES //
/////////////////////

type metaConfig struct {
	Cli        *cli.Context
	ConfigName string
	DataName   string
	Config     scoringChecks
}

type imageData struct {
	Score       int
	RunningTime int // change to time or smth idk
	Points      []scoreItem
}

type scoreItem struct {
	Message string
	Points  int
}

type scoringChecks struct {
	Name  string
	Title string
	User  string
	Check []check
}

type check struct {
	Message string
	Points  int
	Pass    []condition
	Fail    []condition
}

type condition struct {
	Type string
	Arg1 string
	Arg2 string
}

///////////////////////
// CONTROL FUNCTIONS //
///////////////////////

func scoreImage(mc *metaConfig, id *imageData, outputType string) {
	readConfig(mc, readData(mc))
	id.Score = 0
    for i, check := range mc.Config.Check {
		status := false
		for _, condition := range check.Pass {
			status = processCheck(mc, &check, condition.Type, condition.Arg1, condition.Arg2)
			if status {
				break
			}
		}
		for _, condition := range check.Fail {
			status = processCheck(mc, &check, condition.Type, condition.Arg1, condition.Arg2)
			if status {
				break
			}
		}
		if status {
			if mc.Cli.Bool("v") {
				fmt.Printf("[PASS] Check passed: ")
                fmt.Printf("%d :: %s - %d pts\n", i, check.Message, check.Points)
			}
			id.Score += check.Points
		} else {
			if mc.Cli.Bool("v") {
				fmt.Printf("[FAIL] Check failed: ")
                fmt.Printf("%d :: %s - %d pts\n", i, check.Message, check.Points)
			}
        }
	}
    if mc.Cli.Bool("v") {
		fmt.Printf("[INFO] Score: ")
	}
    fmt.Println(id.Score)
}

func processCheck(mc *metaConfig, check *check, checkType string, arg1 string, arg2 string) bool {
	switch checkType {
	case "Command":
        if check.Message == "" {
            check.Message = "Command \"" + arg1 + "\" passed"
        }
		return Command(arg1)
	case "CommandNot":
        if check.Message == "" {
            check.Message = "Command \"" + arg1 + "\" failed"
        }
		return !Command(arg1)
    case "FileExists":
        if check.Message == "" {
            check.Message = "File \"" + arg1 + "\" exists"
        }
        return FileExists(arg1)
    case "FileExistsNot":
        if check.Message == "" {
            check.Message = "File \"" + arg1 + "\" does not exist"
        }
        return !FileExists(arg1)
    case "FileEquals":
        if check.Message == "" {
            check.Message = "File \"" + arg1 + "\" matches hash"
        }
        return FileEquals(arg1, arg2)
    case "FileEqualsNot":
        if check.Message == "" {
            check.Message = "File \"" + arg1 + "\" doesn't match hash"
        }
        return !FileEquals(arg1, arg2)
	default:
		if mc.Cli.Bool("v") {
			fmt.Println("[ERROR] No check type " + checkType)
		}
	}
	return false
}

func releaseImage(mc *metaConfig) {
	cleanUp()
	writeConfig(mc)
	fmt.Println("release - put stuff on desktop, service, etc")
	// add self to services
	// set up notifications

}

////////////////////////////
// MISC/UTILITY FUNCTIONS //
////////////////////////////

func cleanUp() {
	fmt.Println("[INFO] Cleaning up the system...")
	// viminfo, scoring.conf, etc

	//    <Execute>rm -f /home/*/.local/share/recently-used.xbel</Execute>
	//    <Execute>echo Running installation commands</Execute>
	//    <Execute>rm -f /home/*/Desktop/*~</Execute>
	//    <Execute>rm -f /var/crash/*.crash</Execute>
	//    <Execute>rm -f /var/VMwareDnD/*</Execute>
}

func destroyImage() {
	// destroy the image if outside time range or time limit
	fmt.Println("destroying the system lol")
}

func outputHtml(outFile string, data *imageData) {
	// cat header
	fmt.Println("generating html, outputting to" + outFile)
	// report gen at
	// approximate image running time
	// team id?
	//points out of Point
	// if remote, public scoreboard
	// connection status
	// penalties
	// x out of x scored security issues fixed, for a gain of X points:
	// message - $points pts
	// cat footer

}

func readFile(fileName string) string {
	fileContent, err := ioutil.ReadFile(fileName)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return string(fileContent)
}

//////////////////////
// CONFIG FUNCTIONS //
//////////////////////

func readConfig(mc *metaConfig, configContent string) {
	if _, err := toml.Decode(configContent, &mc.Config); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func printConfig(mc *metaConfig) {
	fmt.Printf("[INFO] Configuration " + mc.ConfigName + " check passed!\n")
	fmt.Printf("Title: %s (%s)\n", mc.Config.Title, mc.Config.Name)
	fmt.Printf("User: %s\n", mc.Config.User)

	fmt.Println("Checks:")
	for _, check := range mc.Config.Check {
		fmt.Printf("\tCheck X (%d points):\n", check.Points)
		fmt.Printf("\t\tMessage: %s\n", check.Message)
		if check.Pass != nil {
			fmt.Printf("\t\tPassConditions:\n")
			for _, condition := range check.Pass {
				fmt.Printf("\t\t\t%s: %s, %s\n", condition.Type, condition.Arg1, condition.Arg2)
			}
		}
		if check.Fail != nil {
			fmt.Printf("\t\tFailConditions:\n")
			for _, condition := range check.Fail {
				fmt.Printf("\t\t\t%s: %s, %s\n", condition.Type, condition.Arg1, condition.Arg2)
			}
		}
	}
}

func writeConfig(mc *metaConfig) {
	if mc.Cli.Bool("v") {
		fmt.Println("[INFO] Reading configuration from " + mc.ConfigName + "...")
	}

	configFile, err := os.Open(mc.ConfigName)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer configFile.Close()

	info, _ := configFile.Stat()
	var size int64 = info.Size()
	configBuffer := make([]byte, size)
	buffer := bufio.NewReader(configFile)
	_, err = buffer.Read(configBuffer)

	if mc.Cli.Bool("v") {
		fmt.Println("[INFO] Encrypting configuration...")
	}

	encryptedBuffer := encryptConfig(configBuffer)

	if mc.Cli.Bool("v") {
		fmt.Println("[INFO] Writing data to " + mc.DataName + "...")
	}
	err = ioutil.WriteFile(mc.DataName, encryptedBuffer.Bytes(), info.Mode())
}

func readData(mc *metaConfig) string {
	if mc.Cli.Bool("v") {
		fmt.Println("[INFO] Decrypting data from" + mc.DataName)
	}

	dataFile, err := os.Open(mc.DataName)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer dataFile.Close()

	info, err := os.Stat(mc.ConfigName)
	if err != nil {
		fmt.Println(err)
	}
	var size int64 = info.Size()
	dataBuffer := make([]byte, size) // FIXME TODO fixed buffer size cuts off config
	buffer := bufio.NewReader(dataFile)
	_, err = buffer.Read(dataBuffer)

	return decryptData(dataBuffer)
}

/////////////////////////////
// CRYPTOGRAPHIC FUNCTIONS //
/////////////////////////////

func encryptConfig(configFile []byte) bytes.Buffer {
	// xor with defined byte string
	// zlib
	var encryptedFile bytes.Buffer
	writer := zlib.NewWriter(&encryptedFile)
	writer.Write(configFile)
	writer.Close()
	// aes with reversed byte string or something

	return encryptedFile
}

func decryptData(dataFile []byte) string {
	// aes with reversed byte string

	tmpBuff := bytes.NewReader(dataFile)
	reader, _ := zlib.NewReader(tmpBuff)
	reader.Read(dataFile)
	reader.Close()

	// xor with defined byte string

	return string(dataFile)
}

/////////////////////
// CHECK FUNCTIONS //
/////////////////////

func Command(commandGiven string) bool {
    cmd := exec.Command("sh", "-c", commandGiven)
    if err := cmd.Run(); err != nil {
        fmt.Println(err)
        if exitError, ok := err.(*exec.ExitError); ok {
            fmt.Println(exitError.ExitCode())
            return false
        }
    }
	return true
}

func FileExists(checkFile string) bool {
    if _, err := os.Stat(checkFile); os.IsNotExist(err) {
        return false
    }
    return true
}

func FileContains() {

}

func FileEquals(fileName string, fileHash string) bool {
	fileContent, err := ioutil.ReadFile(fileName)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
    hasher := sha1.New()
    hasher.Write(fileContent)
    hash := hex.EncodeToString(hasher.Sum(nil))
    if hash == fileHash {
        return true
    }
    return false
}

func PackageInstalled() {

}

func UserExists() {

}
