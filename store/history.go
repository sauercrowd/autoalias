package store

import (
	"time"
)

func (s *Store) AddHistory(command []string) error {
	commandStr, err := ArrayToJsonArray(command)
	if err != nil {
		return err
	}
	_, err = s.historyDb.Exec("INSERT INTO history (command_array, used) VALUES (?, ?)", commandStr, time.Now().Format(time.RFC3339))
	return err
}

type Result struct {
	Command  []string
	Executed time.Time
}

func (s *Store) GetHistory() ([]Result, error) {
	results := make([]Result, 0)
	rows, err := s.historyDb.Query("SELECT command_array, used FROM history")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var commandJson string
		var d string
		if err := rows.Scan(&commandJson, &d); err != nil {
			return nil, err
		}
		var r Result
		if r.Command, err = JsonArrayToArray(commandJson); err != nil {
			return nil, err
		}
		r.Executed, err = time.Parse(time.RFC3339, d)
		if err != nil {
			return nil, err
		}
		results = append(results, r)
	}
	return results, nil
}
