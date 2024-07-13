package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
)

var hexToLetter map[string]string
var letterToHex map[string]string

func init() {
	hexToLetter = createHexToLetterMap()
	letterToHex = make(map[string]string)
	for hex, letter := range hexToLetter {
		letterToHex[letter] = hex
	}
}

func createHexToLetterMap() map[string]string {
	hexToLetter := make(map[string]string)
	letters := "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	
	index := 0
	for i := 0; i < 16; i++ {
		for j := 0; j < 16; j++ {
			hex := fmt.Sprintf("%02X", index)
			letter1 := string(letters[i])
			letter2 := string(letters[j])
			
			if strings.ContainsAny(hex, "ABCDEF") {
				letter2 = string(letters[(j+1)%26])
			}
			
			hexToLetter[hex] = letter1 + letter2
			index++
		}
	}
	
	return hexToLetter
}

func encode(r io.Reader, w io.Writer, lineLength int, uppercase bool) error {
	scanner := bufio.NewScanner(r)
	lineCount := 0

	for scanner.Scan() {
		line := scanner.Text()
		line = strings.TrimSpace(line)
		if len(line)%2 != 0 {
			return fmt.Errorf("invalid hex input: odd number of characters")
		}

		for i := 0; i < len(line); i += 2 {
			hex := strings.ToUpper(line[i : i+2])
			if letter, ok := hexToLetter[hex]; ok {
				if !uppercase {
					letter = strings.ToLower(letter)
				}
				fmt.Fprint(w, letter)
				lineCount += 2
				if lineLength > 0 && lineCount >= lineLength {
					fmt.Fprintln(w)
					lineCount = 0
				}
			} else {
				return fmt.Errorf("invalid hex value: %s", hex)
			}
		}
	}
	if lineCount > 0 && lineLength > 0 {
		fmt.Fprintln(w)
	}
	return scanner.Err()
}

func decode(r io.Reader, w io.Writer, uppercase bool) error {
	scanner := bufio.NewScanner(r)
	
	for scanner.Scan() {
		line := scanner.Text()
		line = strings.TrimSpace(line)
		line = strings.ToUpper(line)
		if len(line)%2 != 0 {
			return fmt.Errorf("invalid input: odd number of characters")
		}

		for i := 0; i < len(line); i += 2 {
			pair := line[i : i+2]
			if hex, ok := letterToHex[pair]; ok {
				if uppercase {
					hex = strings.ToUpper(hex)
				} else {
					hex = strings.ToLower(hex)
				}
				fmt.Fprint(w, hex)
			} else {
				return fmt.Errorf("invalid letter pair: %s", pair)
			}
		}
	}
	fmt.Fprintln(w)
	return scanner.Err()
}

func usage() {
	fmt.Fprintf(os.Stderr, "Usage: %s [OPTIONS]\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "Encode hexadecimal data to letter pairs or decode letter pairs back to hexadecimal.\n\n")
	fmt.Fprintf(os.Stderr, "Options:\n")
	flag.PrintDefaults()
}

func main() {
	decodeFlag := flag.Bool("d", false, "Decode mode")
	lineLengthFlag := flag.Int("l", 64, "Line length for encoding (0 for no line breaks)")
	uppercaseFlag := flag.Bool("u", false, "Use uppercase letters (default is lowercase)")

	flag.Usage = usage
	flag.Parse()

	stdinInfo, _ := os.Stdin.Stat()
	hasInput := (stdinInfo.Mode() & os.ModeCharDevice) == 0

	if flag.NFlag() == 0 && flag.NArg() == 0 && !hasInput {
		usage()
		os.Exit(1)
	}

	var err error
	if *decodeFlag {
		err = decode(os.Stdin, os.Stdout, *uppercaseFlag)
	} else {
		err = encode(os.Stdin, os.Stdout, *lineLengthFlag, *uppercaseFlag)
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
