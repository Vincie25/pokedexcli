package main

import (
    "bufio"
    "fmt"
    "os"
    "strings"
)

type cliCommand struct {
    name        string
    description string
    callback    func() error
}

var commands map[string]cliCommand

func commandExit() error {
    fmt.Println("Closing the Pokedex... Goodbye!")
    os.Exit(0)
    return nil
}

func commandHelp() error {
    fmt.Println("Welcome to the Pokedex!")
    fmt.Println("Usage:\n")
    for _, cmd := range commands {
        fmt.Printf("%s: %s\n", cmd.name, cmd.description)
    }
    return nil
}

func main() {
    commands = map[string]cliCommand{
        "exit": {
            name:        "exit",
            description: "Exit the Pokedex",
            callback:    commandExit,
        },
        "help": {
            name:        "help",
            description: "Displays a help message",
            callback:    commandHelp,
        },
    }
    scanner := bufio.NewScanner(os.Stdin)

    for {
        fmt.Print("Pokedex > ")
        scanner.Scan()
        input := scanner.Text()
        input = strings.TrimSpace(input)
        input = strings.ToLower(input)

        words := strings.Fields(input)
        if len(words) == 0 {
            continue
        }

        firstWord := words[0]
        cmd, ok := commands[firstWord]
        if ok {
            err := cmd.callback()
            if err != nil {
                fmt.Println(err)
            }
        } else {
            fmt.Println("Unknown command")
        }
    }
}
