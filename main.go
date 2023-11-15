package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
)

type repoInfo struct {
	Name     string
	Archived bool
	Url      string
}

func main() {
	fmt.Fprintln(os.Stderr, "Checking repos...")
	getJSON("freecodecamp")
	readFile()
}

func getJSON(org string) {
	var repos []repoInfo
	// pagination for github api
	for page := 1; ; page++ {
		fmt.Println("Fetching page", page)
		url := fmt.Sprintf("https://api.github.com/orgs/%s/repos?page=%d", org, page)
		cmd := exec.Command(
			"curl",
			"-H", "Authorization: Bearer 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",
			"-H", "Accept: application/vnd.github.v3+json",
			"-s", url,
		)
		output, err := cmd.Output()
		if err != nil {
			fmt.Println("Error executing curl command:", err)
			break
		}
		var pageRepos []repoInfo
		err = json.Unmarshal(output, &pageRepos)
		if err != nil {
			fmt.Println("Error unmarshaling JSON:", err)
			break
		}

		if len(pageRepos) == 0 {
			fmt.Println("repo fetch completed")
			break
		}
		repos = append(repos, pageRepos...)
	}
	writeToFile("repositories.json", repos)
}

func writeToFile(filename string, data interface{}) {
	file, err := os.Create(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")

	err = encoder.Encode(data)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Data written to %s\n", filename)
}

func readFile() {
	_, err := os.Stat("repositories.json")
	if os.IsNotExist(err) {
		fmt.Println("File 'repositories.json' not found.")
		return
	} else if err != nil {
		log.Fatal(err)
	}
	f, err := os.Open("repositories.json")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	b, err := io.ReadAll(f)
	if err != nil {
		log.Fatal(err)
	}
	var ri []repoInfo
	err = json.Unmarshal(b, &ri)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Found %d repos\n", len(ri))
	for _, v := range ri {
		fmt.Printf("name: %s archived: %t url: %s\n", v.Name, v.Archived, v.Url)
	}
}
