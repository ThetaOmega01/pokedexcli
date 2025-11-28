package cli

import (
	"fmt"
)

func commandPokedex(args ...string) error {
	if len(pokedex) == 0 {
		fmt.Println("Your Pokedex is empty. Catch some Pokemon first!")
		return nil
	}

	fmt.Println("Your Pokedex:")
	for name := range pokedex {
		fmt.Printf("  - %s\n", name)
	}

	return nil
}
