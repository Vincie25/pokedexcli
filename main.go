package main

import (
    "github.com/Vincie25/pokedexcli/internal/pokecache"
    "bufio"
    "fmt"
    "os"
    "strings"
    "encoding/json"
    "io"
    "net/http"
    "time"
)

type cliCommand struct {
    name        string
    description string
    callback    func(*config, []string) error
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

type exploreResponse struct {
	PokemonEncounters []struct {
		Pokemon struct {
			Name string `json:"name"`
		} `json:"pokemon"`
	} `json:"pokemon_encounters"`
}

type config struct {
    Next     *string
    Previous *string
    cache    *pokecache.Cache
}

var commands map[string]cliCommand

func commandExit(cfg *config, args []string) error {
    fmt.Println("Closing the Pokedex... Goodbye!")
    os.Exit(0)
    return nil
}

func commandHelp(cfg *config, args []string) error {
    fmt.Println("Welcome to the Pokedex!")
    fmt.Print("Usage:\n\n")
    for _, cmd := range commands {
        fmt.Printf("%s: %s\n", cmd.name, cmd.description)
    }
    return nil
}

func commandMap(cfg *config, args []string) error {
    url := "https://pokeapi.co/api/v2/location-area"
    if cfg.Next != nil {
        url = *cfg.Next
    }

    var pokeRes pokeResponse

    if val, ok := cfg.cache.Get(url); ok {
        // Cache hit
        json.Unmarshal(val, &pokeRes)
    } else {
        // Cache miss
        res, err := http.Get(url)
        if err != nil {
            return err
        }
        defer res.Body.Close()
        body, err := io.ReadAll(res.Body)
        if err != nil {
            return err
        }
        cfg.cache.Add(url, body)
        json.Unmarshal(body, &pokeRes)
    }

    cfg.Next = pokeRes.Next
    cfg.Previous = pokeRes.Previous
    for _, location := range pokeRes.Results {
        fmt.Println(location.Name)
    }
    return nil
}

func commandMapb(cfg *config, args []string) error {
    if cfg.Previous == nil {
        fmt.Println("you're on the first page")
        return nil
    }

    url := *cfg.Previous
    var pokeRes pokeResponse

    if val, ok := cfg.cache.Get(url); ok {
        // Cache hit
        json.Unmarshal(val, &pokeRes)
    } else {
        // Cache miss
        res, err := http.Get(url)
        if err != nil {
            return err
        }
        defer res.Body.Close()
        body, err := io.ReadAll(res.Body)
        if err != nil {
            return err
        }
        cfg.cache.Add(url, body)
        json.Unmarshal(body, &pokeRes)
    }

    cfg.Next = pokeRes.Next
    cfg.Previous = pokeRes.Previous
    for _, location := range pokeRes.Results {
        fmt.Println(location.Name)
    }
    return nil
}

func commandExplore(cfg *config, args []string) error {
	if len(args) == 0 {
		fmt.Println("please provide a location area name")
		return nil
	}

	name := args[0]
	url := "https://pokeapi.co/api/v2/location-area/" + name
	fmt.Printf("Exploring %s...\n", name)

	var exploreRes exploreResponse
	if val, ok := cfg.cache.Get(url); ok {
		json.Unmarshal(val, &exploreRes)
	} else {
		res, err := http.Get(url)
		if err != nil {
			return err
		}
		defer res.Body.Close()
		body, err := io.ReadAll(res.Body)
		if err != nil {
			return err
		}
		cfg.cache.Add(url, body)
		json.Unmarshal(body, &exploreRes)
	}

	fmt.Println("Found Pokemon:")
	for _, encounter := range exploreRes.PokemonEncounters {
		fmt.Printf(" - %s\n", encounter.Pokemon.Name)
	}
	return nil
}

func main() {
    cfg := &config{
        cache: pokecache.NewCache(5 * time.Second),
    }
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
        "explore": {
			name:        "explore",
			description: "Explore a location area",
			callback:    commandExplore,
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
            err := cmd.callback(cfg, words[1:])
            if err != nil {
                fmt.Println(err)
            }
        } else {
            fmt.Println("Unknown command")
        }
    }
}
