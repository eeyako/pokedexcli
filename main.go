package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/chelloiaco/pokedexcli/internal"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Pokedex > ")
		scanner.Scan()
		text := strings.Split(scanner.Text(), " ")
		result, ok := internal.GetCommands()[text[0]]
		if ok {
			if len(text) < 2 {
				result.Callback("")
			} else {
				result.Callback(text[1])
			}
		} else {
			fmt.Printf("\nInput not recognized - '%s'\nType 'help' for more info\n", scanner.Text())
		}
		fmt.Println()
	}
}
