package main

import (
	"io"
	"net/http"
	"os"
	"regexp"
	"sync"
	"testing"
)

/*
This test makes sure all links in README.adoc are reachable (status 200)
*/
func TestLinksInReadmeValid(t *testing.T) {
	t.Log("Given README.adoc loaded")
	{
		file, err := os.Open("README.adoc")
		if err != nil {
			t.Fatalf("failed to open README.adoc: %s", err)
		}
		defer file.Close()
		t.Log("\tWhen links are obtained")
		{
			// read entire file to string at once
			content, err := io.ReadAll(file)
			if err != nil {
				t.Fatalf("failed to read README.adoc: %s", err)
			}
			results := regexp.MustCompile(`link:(https://github\.com/[^\[]+)\[([^]]+)]`).FindAllStringSubmatch(string(content), -1)
			links := make([]string, len(results))
			t.Logf("\t\tThen %d links are found", len(results))
			{
				if results == nil || len(results) == 0 {
					t.Fatal("\t\tno links found in README.adoc")
				}
				for i, result := range results {
					links[i] = result[1]
					t.Logf("\t\t\t- %s", links[i])
				}
			}
			t.Logf("\t\tThen all links are valid")
			{
				var wg sync.WaitGroup
				wg.Add(len(links))
				results := make(chan result)

				for _, link := range links {
					go func(l string) {
						rs, err := http.Get(l)
						if err != nil {
							results <- result{link: l, error: err}
						} else {
							results <- result{link: l, code: &rs.StatusCode}
						}
						rs.Body.Close()
						wg.Done()
					}(link)
				}
				go func() {
					wg.Wait()
					close(results)
				}()

				failed := false
				for r := range results {
					if r.error != nil {
						t.Logf("\t\t\tfailed to get link %s: %s", r.link, r.error)
						failed = true
					}
					if r.code != nil && *r.code != http.StatusOK {
						t.Logf("\t\t\tlink %s returned status %d", r.link, *r.code)
						failed = true
					}
				}
				if failed {
					t.Fatal("\t\t\tfailed to verify any of the links")
				}
			}
		}
	}
}

type result struct {
	link  string
	code  *int
	error error
}
