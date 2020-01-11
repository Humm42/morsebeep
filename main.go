package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"time"
)

func beepMilliseconds(time string) {
	err := exec.Command("beep", "-l", time).Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: beep: %s\n", os.Args[0], err)
		os.Exit(1)
	}
}

func beepText(str string) {
	fields := strings.Fields(str)
	var runeArg string
	for _, v := range fields {
		switch v {
		case "dit":
			runeArg += "."
		case "daw":
			runeArg += "-"
		case "dah":
			runeArg += "-"
		default:
			fmt.Fprintf(os.Stderr, "\nnot ‘dit’ or ‘daw’: ‘%s’\n",
				v)
			os.Exit(1)
		}
	}
	beepRunes(runeArg)
}

func beepRunes(str string) {
	if strings.TrimSpace(str) == "" {
		// BUG: sleeps too short if first input line
		// sleep 4
		time.Sleep(200 * time.Millisecond)
		return
	}
	for _, v := range str {
		switch v {
		case '.':
			// beep 1
			beepMilliseconds("50")
		case '-':
			// beep 3
			beepMilliseconds("150")
		case ' ':
		default:
			fmt.Fprintf(os.Stderr, "\nnot ‘.’ or ‘-’: %v\n", v)
			os.Exit(1)
		}
		// sleep 1
		time.Sleep(50 * time.Millisecond)
	}
	// sleep 2
	time.Sleep(100 * time.Millisecond)
}

func main() {
	var err error

	_, err = exec.LookPath("beep")
	if err != nil {
		fmt.Fprintln(os.Stderr, "executable “beep” not found")
		os.Exit(1)
	}

	var sFlag bool
	flag.BoolVar(&sFlag, "s", false, "accept dashes and dots")
	flag.Parse()
	args := flag.Args()

	var inputFile *os.File

	switch len(args) {
	case 0:
		inputFile = os.Stdin
	case 1:
		inputFile, err = os.Open(args[0])
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s: %s: %s\n", os.Args[0],
				args[0], err)
			os.Exit(1)
		}
	default:
		flag.Usage()
		os.Exit(1)
	}

	var inputBytes []byte
	inputBytes, err = ioutil.ReadAll(inputFile)
	if inputFile != os.Stdin {
		inputFile.Close()
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: %s\n", os.Args[0], err)
		os.Exit(1)
	}
	input := string(inputBytes)

	// remove trailing newline
	if len(input) > 0 && input[len(input)-1] == '\n' {
		input = input[:len(input)-1]
	}

	lines := strings.Split(input, "\n")

	if sFlag {
		for _, v := range lines {
			beepRunes(v)
		}
	} else {
		for _, v := range lines {
			beepText(v)
		}
	}
}
