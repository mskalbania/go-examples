package which

import (
	"fmt"
	"os"
	"path/filepath"
)

//Which
/*
Based on example from "Mastering Go 4th ed".
Expects two args to be present - bin <executable_name>. Eg - bin git
Scans PATH for given binary.
Prints result
*/
func Which() {
	var foundLocations []string
	executableName := getBinaryName()
	fmt.Printf("Looking for binary '%s' in PATH\n", executableName)
	for _, path := range getPathPaths() {
		fmt.Printf("Searching in: %s\n", path)
		files, _ := os.ReadDir(path)
		for _, file := range files {
			info, _ := file.Info()
			if file.Name() == executableName && info.Mode()&0111 != 0 {
				foundLocations = append(foundLocations, filepath.Join(path, executableName))
			}
		}
	}
	if len(foundLocations) == 0 {
		fmt.Printf("\nBinary [%s] not found in PATH\n", executableName)
	} else {
		fmt.Printf("\nBinary [%s] found at:\n", executableName)
		for _, location := range foundLocations {
			fmt.Println(location)
		}
		fmt.Println()
	}
}

func getBinaryName() string {
	for i, arg := range os.Args {
		if arg == "bin" {
			return os.Args[i+1]
		}
	}
	fmt.Println("Provide executable name after bin flag example - 'bin git'")
	os.Exit(1)
	return ""
}

func getPathPaths() []string {
	path := os.Getenv("PATH")
	if path != "" {
		return filepath.SplitList(path)
	}
	fmt.Println("PATH env not found")
	os.Exit(1)
	return nil
}
