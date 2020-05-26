package main

import (
	"fmt"
	"os/user"
)

func main() {
	printBanner()
	printGreeting()
	//repl.Start(os.Stdin, os.Stdout)
}

func printBanner() {
	banner := `
   __  ____                 
  /  |/  (_)__________  ___ 
 / /|_/ / / __/ __/ _ \/ _ \
/_/  /_/_/\__/_/  \___/_//_/

`

	fmt.Print(banner)
}

func printGreeting() {
	currentUser, err := user.Current()
	if err != nil {
		panic(err)
	}

	fmt.Printf("Hello %s! Welcome to the Micron language console!\n", currentUser.Username)
	fmt.Printf("Feel free to type in the code below\n\n")
}
