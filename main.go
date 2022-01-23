package main

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/alecthomas/kong"
	"github.com/sauercrowd/autoalias/store"
)

var CLI struct {
	Store struct {
		Command []string `arg: ""`
	} `cmd:"" help:"store a new command that was executed"`
	Render struct {
	} `cmd:"" help:"Generate a file that can be sourced by a shell"`
	History struct {
	} `cmd:"" help:"Print out the history of commands"`
}

func main() {
	ctx := kong.Parse(&CLI)
	store, err := store.New()
	if err != nil {
		log.Fatal(err)
	}
	switch ctx.Command() {
	case "store <command>":
		cmd := strings.Split(CLI.Store.Command[0], " ")
		if err := store.AddHistory(cmd); err != nil {
			log.Println(err)
			return
		}

		// generate new aliases since this is the only chance it could've changed
		if err := store.GenerateNewAlises(time.Hour*24*7, 20); err != nil {
			log.Println(err)
		}
	case "render":
		aliases, err := store.GetAliases()
		if err != nil {
			log.Fatal(err)
		}
		for _, alias := range aliases {
			parsedCmd := strings.Join(alias.Command, " ")
			fmt.Printf("alias %s='%s'\n", alias.Name, parsedCmd)
		}
	case "history":
		results, err := store.GetHistory()
		if err != nil {
			log.Fatal(err)
		}
		for _, result := range results {
			parsedCmd := strings.Join(result.Command, " ")
			fmt.Println(result.Executed.Format(time.RFC3339), parsedCmd)
		}
	}

}
