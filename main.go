package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type process struct {
	PID    int        `json:"pid"`
	PPID   int        `json:"ppid"`
	CMD    string     `json:"cmd"`
	Childs []*process `json:"children,omitempty"`
}

func main() {
	// read csv
	pwd, _ := os.Getwd()
	FilePath := filepath.Join(pwd, "process.csv")
	file, err := os.OpenFile(FilePath, os.O_RDONLY, 0777)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		fmt.Println("Error reading CSV data:", err)
		return
	}
	res := ConvertCsvToJson(records)
	fmt.Println(res)
}

func ConvertCsvToJson(records [][]string) string {
	// make processes map
	processes := make(map[int]*process)
	for _, record := range records {
		pid, ppid := 0, 0
		fmt.Sscanf(record[0], "%d", &pid)
		fmt.Sscanf(record[1], "%d", &ppid)

		p := &process{
			PID:  pid,
			PPID: ppid,
			CMD:  record[2],
		}
		processes[pid] = p
		// Check if there is a parent node
		if parent, ok := processes[ppid]; ok {
			parent.Childs = append(parent.Childs, p)
		}
	}

	// find roots node
	var roots []*process
	for _, p := range processes {
		if p.PPID == 0 {
			roots = append(roots, p)
		}
	}

	// Marshal to json
	// can use MarshalIndent to format json
	result, err := json.Marshal(roots)
	if err != nil {
		fmt.Println("Error marshaling JSON:", err)
		return ""
	}
	return string(result)
}
