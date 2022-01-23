package store

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"
)

type Alias struct {
	Name    string
	Command []string
}

func (s *Store) GetAliases() ([]Alias, error) {
	rows, err := s.aliasDb.Query("SELECT alias, command FROM aliases")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	aliases := make([]Alias, 0)
	for rows.Next() {
		var a Alias
		var jsonArr string
		rows.Scan(&a.Name, &jsonArr)
		a.Command, err = JsonArrayToArray(jsonArr)
		if err != nil {
			return nil, err
		}
		aliases = append(aliases, a)
	}
	return aliases, nil
}

func (s *Store) AddAlias(cmd []string, alias string) error {
	arr, err := ArrayToJsonArray(cmd)
	if err != nil {
		return err
	}
	_, err = s.aliasDb.Exec("INSERT INTO aliases(alias, command) VALUES(?, ?)", alias, arr)
	return err
}

func (s *Store) HasAliasWithCommand(jsonArr string) (bool, error) {
	rows, err := s.aliasDb.Query("SELECT * FROM aliases WHERE command=?", jsonArr)
	if err != nil {
		return false, err
	}
	defer rows.Close()
	for rows.Next() {
		return true, nil
	}
	return false, nil
}

func (s *Store) GenerateNewAlises(window time.Duration, amount int) error {
	rows, err := s.historyDb.Query(`SELECT command_array FROM history WHERE used >= ? GROUP BY command_array HAVING COUNT(used) >= ?`, time.Now().Add(-window).Format(time.RFC3339), amount)
	if err != nil {
		return err
	}
	defer rows.Close()

	aliases := make([]struct {
		cmd   []string
		alias string
	}, 0)

	foundRows := false
	for rows.Next() {
		var jsonArr string
		if err := rows.Scan(&jsonArr); err != nil {
			return err
		}
		ok, err := s.HasAliasWithCommand(jsonArr)
		if err != nil {
			return err
		}
		if ok {
			continue
		}
		cmd, err := JsonArrayToArray(jsonArr)
		if err != nil {
			return err
		}
		alias, err := generateAlias(cmd)
		if err != nil {
			continue // cause the only error can be a collission
		}
		foundRows = true
		aliases = append(aliases, struct {
			cmd   []string
			alias string
		}{
			alias: alias, cmd: cmd,
		})
		if err := s.AddAlias(cmd, alias); err != nil {
			return err
		}
	}

	if foundRows {
		fmt.Println("Hey! I created a new alias for you")
		for _, alias := range aliases {
			fmt.Printf("    %s=%s\n", alias.alias, strings.Join(alias.cmd, " "))
		}
	}
	return nil
}

func generateAlias(command []string) (string, error) {
	candidate := ""
	// first try just first letters
	for i, c := range command {
		c = strings.Trim(c, " 	")
		candidate += string(c[0])
		if i > 0 && AliasHasNoOverlap(candidate) {
			return candidate, nil
		}
	}

	// if we still weren't successfull add more from the last command
	for i, c := range command[len(command)-1] {
		if c == ' ' {
			continue
		}
		if i == 0 {
			continue
		}
		candidate += string(c)
		if AliasHasNoOverlap(candidate) {
			return candidate, nil
		}
	}

	return "", fmt.Errorf("Unable to create alias for command, to many collisions")

}

func AliasHasNoOverlap(alias string) bool {
	path, ok := os.LookupEnv("PATH")
	if !ok {
		return true
	}
	for _, dir := range strings.Split(path, ":") {
		files, err := ioutil.ReadDir(dir)
		if err != nil {
			continue
		}
		for _, f := range files {
			if f.Name() == alias {
				return false
			}
		}
	}
	return true
}
