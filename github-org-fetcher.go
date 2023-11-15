package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"
	"sync"
	"time"
)

type repoInfo struct {
	Name     string
	Archived bool
	Url      string
}

func main() {
	var org string
	flag.StringVar(&org, "o", "", "Organization name")

	flag.Parse()

	if org == "" {
		args := flag.Args()
		if len(args) > 0 {
			org = args[0]
		} else {
			fmt.Fprintln(os.Stderr, "error: organization name is required.")
			os.Exit(1)
		}
	}

	var wg sync.WaitGroup
	stopCh := make(chan struct{})
	defer close(stopCh)

	wg.Add(1)
	go func() {
		defer wg.Done()
		printAnimatedDots(stopCh)
	}()
	fmt.Fprintln(os.Stderr, "checking repos for organization: ", org)
	getJSON(org)
	readFile()
	fmt.Fprintln(os.Stderr, "repositories.json created")
}

func printAnimatedDots(stopCh <-chan struct{}) {
	dots := []string{".", "..", "...", "...."}
	i := 0
	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			fmt.Print("\rFetching" + dots[i])
			i = (i + 1) % len(dots)
		case <-stopCh:
			return
		}
	}
}

func getJSON(org string) {
	var repos []repoInfo
	// pagination for github api
	for page := 1; ; page++ {
		url := fmt.Sprintf("https://api.github.com/orgs/%s/repos?page=%d", org, page)
		cmd := exec.Command(
			"curl",
			"-H", "Authorization: Bearer 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",
			"-H", "Accept: application/vnd.github.v3+json",
			"-s", url,
		)
		output, err := cmd.Output()
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error executing curl command:", err)
			break
		}
		var pageRepos []repoInfo
		err = json.Unmarshal(output, &pageRepos)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error unmarshaling JSON:", err)
			break
		}
		if len(pageRepos) == 0 {
			fmt.Fprintln(os.Stderr, "repo fetch completed")
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
		fmt.Fprintf(os.Stderr, "error while encoding data: %s\n", err)
	}
	fmt.Fprintf(os.Stderr, "Data written to %s\n", filename)
}

func readFile() {
	_, err := os.Stat("repositories.json")
	if os.IsNotExist(err) {
		fmt.Fprintln(os.Stderr, "File 'repositories.json' not found.")
		return
	} else if err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", err)
	}
	f, err := os.Open("repositories.json")
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", err)
	}
	defer f.Close()
	b, err := io.ReadAll(f)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", err)
	}
	var ri []repoInfo
	err = json.Unmarshal(b, &ri)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", err)
	}
	fmt.Fprintf(os.Stderr, "Found %d repos\n", len(ri))

	// format output
	fmt.Printf("%-30s%-10s%s\n", "Name", "Archived", "URL")
	fmt.Println(strings.Repeat("-", 60))

	for _, v := range ri {
		fmt.Printf("%-30s%-10t%s\n", v.Name, v.Archived, v.Url)
	}
}
