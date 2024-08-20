package main

import (
	"io"
	"net/http"
	"os"
	"regexp"
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
				for _, link := range links {
					rs, err := http.Get(link)
					if err != nil {
						t.Fatalf("\t\tfailed to get link %s: %s", link, err)
					}
					rs.Body.Close()
					if rs.StatusCode != http.StatusOK {
						t.Fatalf("\t\tlink %s returned status %d", link, rs.StatusCode)
					}
				}
			}
		}
	}
}
