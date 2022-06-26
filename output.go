package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/fatih/color"
)

// confirm will prompt the user with the given toPrint string, and
// exit the program if N or n is input.
func confirm(p ...interface{}) {
	if yesEnabled {
		return
	}
	toPrint := fmt.Sprint(p...)
	toPrint = printer(color.FgYellow, "CONF", toPrint)
	fmt.Print(toPrint + " [Y/n]: ")
	var resp string
	fmt.Scanln(&resp)
	if strings.ToLower(strings.TrimSpace(resp)) == "n" {
		os.Exit(1)
	}
}

// ask will prompt the user with the given toPrint string, and
// return a boolean.
func ask(p ...interface{}) bool {
	if yesEnabled {
		return true
	}
	toPrint := fmt.Sprint(p...)
	toPrint = printer(color.FgBlue, "CONF", toPrint)
	fmt.Print(toPrint + " [Y/n]: ")
	var resp string
	fmt.Scanln(&resp)
	if strings.ToLower(strings.TrimSpace(resp)) == "n" {
		return false
	}
	return true
}

func pass(p ...interface{}) {
	toPrint := fmt.Sprintln(p...)
	printStr := printer(color.FgGreen, "PASS", toPrint)
	fmt.Printf(printStr)
}

func fail(p ...interface{}) {
	toPrint := fmt.Sprintln(p...)
	fmt.Printf(printer(color.FgRed, "FAIL", toPrint))
}

func warn(p ...interface{}) {
	toPrint := fmt.Sprintln(p...)
	fmt.Printf(printer(color.FgYellow, "WARN", toPrint))
}

func debug(p ...interface{}) {
	if debugEnabled {
		toPrint := fmt.Sprintln(p...)
		printStr := printer(color.FgMagenta, "DBUG", toPrint)
		fmt.Printf(printStr)
	}
}

func info(p ...interface{}) {
	if verboseEnabled {
		toPrint := fmt.Sprintln(p...)
		printStr := printer(color.FgCyan, "INFO", toPrint)
		fmt.Printf(printStr)
	}
}

func blue(head string, p ...interface{}) {
	toPrint := fmt.Sprintln(p...)
	printStr := printer(color.FgCyan, head, toPrint)
	fmt.Printf(printStr)
}

func red(head string, p ...interface{}) {
	toPrint := fmt.Sprintln(p...)
	fmt.Printf(printer(color.FgRed, head, toPrint))
}

func green(head string, p ...interface{}) {
	toPrint := fmt.Sprintln(p...)
	fmt.Printf(printer(color.FgGreen, head, toPrint))
}

func printer(colorChosen color.Attribute, messageType, toPrint string) string {
	printer := color.New(colorChosen, color.Bold)
	printStr := "["
	printStr += printer.Sprintf(messageType)
	printStr += fmt.Sprintf("] %s", toPrint)
	return printStr
}
