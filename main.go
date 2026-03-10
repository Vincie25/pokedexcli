package main

import (
    "bufio"
    "fmt"
    "os"
    "strings"
    "encoding/json"
    "io"
    "net/http"
)

type cliCommand struct {
    name        string
    description string
    callback    func(*config) error
}

type pokeResponse struct {
	Count    int     `json:"count"`
	Next     *string `json:"next"`
	Previous *string `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}

type config struct {
    Next     *string
    Previous *string
}

var commands map[string]cliCommand

func commandExit(cfg *config) error {
    fmt.Println("Closing the Pokedex... Goodbye!")
    os.Exit(0)
    return nil
}

func commandHelp(cfg *config) error {
    fmt.Println("Welcome to the Pokedex!")
    fmt.Println("Usage:\n")
    for _, cmd := range commands {
        fmt.Printf("%s: %s\n", cmd.name, cmd.description)
    }
    return nil
}

func commandMap(cfg *config) error {
	url := "https://pokeapi.co/api/v2/location-area"
	if cfg.Next != nil {
		url = *cfg.Next
	}

	res, err := http.Get(url)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}

	var pokeRes pokeResponse
	err = json.Unmarshal(body, &pokeRes)
	if err != nil {
		return err
	}

	cfg.Next = pokeRes.Next
	cfg.Previous = pokeRes.Previous

	for _, location := range pokeRes.Results {
		fmt.Println(location.Name)
	}

	return nil
}

func commandMapb(cfg *config) error {
	if cfg.Previous == nil {
		fmt.Println("you're on the first page")
		return nil
	}

	res, err := http.Get(*cfg.Previous)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}

	var pokeRes pokeResponse
	err = json.Unmarshal(body, &pokeRes)
	if err != nil {
		return err
	}

	cfg.Next = pokeRes.Next
	cfg.Previous = pokeRes.Previous

	for _, location := range pokeRes.Results {
		fmt.Println(location.Name)
	}

	return nil
}

func main() {
    cfg := &config{}
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
        "map": {
			name:        "map",
			description: "Displays next 20 location areas",
			callback:    commandMap,
		},
		"mapb": {
			name:        "mapb",
			description: "Displays previous 20 location areas",
            callback:    commandMapb,
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
            err := cmd.callback(cfg)
            if err != nil {
                fmt.Println(err)
            }
        } else {
            fmt.Println("Unknown command")
        }
    }
}
