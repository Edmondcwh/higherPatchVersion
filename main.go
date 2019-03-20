package main

import (
	"context"
	"fmt"
	"sort"
	"github.com/coreos/go-semver/semver"
	"github.com/google/go-github/github"
	"bufio"
	"os"
	"log"
	"strings"
)


type Descend []*semver.Version

func (input Descend) Len() int {
	return len(input)
}

func (input Descend) Swap(i, j int) {
	input[i], input[j] = input[j], input[i]
}

func (input Descend) Less(i, j int) bool {
	return input[j].LessThan(*input[i])
}


// LatestVersions returns a sorted slice with the highest version as its first element and the highest version of the smaller minor versions in a descending order
func LatestVersions(releases []*semver.Version, minVersion *semver.Version) []*semver.Version {
	var versionSlice []*semver.Version
	sort.Sort(Descend(releases))
	for _, version :=  range releases {
		idx := len(versionSlice)
		if (minVersion.LessThan(*version)) {
			if (idx == 0) { // highest version as the versionSlice's first element
				versionSlice = append(versionSlice, version) 
			} else if (versionSlice[idx-1].Major == version.Major) { // make sure verison is the smaller minor versions of the highest version
				if (versionSlice[idx - 1].Minor > version.Minor) { 
					versionSlice = append(versionSlice, version)
				} else if(versionSlice[idx - 1].Minor == version.Minor && versionSlice[idx - 1].Patch < version.Patch) { // update the highest version of the minor version
					versionSlice = versionSlice[:idx-1]
					versionSlice = append(versionSlice, version)
				} 
			} else {
				break // end the loop when version is not the smaller minor versions of the highest version
			}
		} else {
			break // end the loop when version is smaller than minVersion
		}
	}
	return versionSlice
}
// Read file from given path
func readFile(path string) []string {
	
	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var content []string
	for scanner.Scan() {
		if (strings.Contains(scanner.Text(), "/")) {
			content = append(content, scanner.Text())
		}
	}
	return content
}	
// get respositories from the given path
func getRespositories(content []string) {
	client := github.NewClient(nil)
	ctx := context.Background()
	opt := &github.ListOptions{PerPage: 80}
	var verMap = make(map[string] string)
	var repMap = make(map[string] string)
	
	for _, con := range content { // divide the content of file into respository, owner and version
		fmt.Println(con)
		repAndVersion := strings.Split(con, ",") // repAndVersion[0] = respository, repAndVersion[1] = minVersion
		holder := strings.Split(repAndVersion[0], "/") // holder[0] = owner, holder[1] = name of respository
		verMap[holder[0]] = repAndVersion[1]
		repMap[holder[0]] = holder[1]
	}

	for k := range repMap { // get latest versions for each respository
		releases, _, err := client.Repositories.ListReleases(ctx, k, repMap[k], opt)

		if err != nil {
	    	fmt.Printf("error: %v\n", err)
		}
		minVersion := semver.New(verMap[k])
		allReleases := make([]*semver.Version, len(releases))
		for i, release := range releases {
			versionString := *release.TagName
			if versionString[0] == 'v' {
				versionString = versionString[1:]
			}
			allReleases[i] = semver.New(versionString)
		}
		versionSlice := LatestVersions(allReleases, minVersion)	

		fmt.Printf("latest versions of %s/%s: %s \n", k, repMap[k], versionSlice)
	}
}	
	// This is just an example structure of the code, if you implement this interface, the test cases in main_test.go are very easy to run}

// Here we implement the basics of communicating with github through the library as well as printing the version
// You will need to implement LatestVersions function as well as make this application support the file format outlined in the README
// Please use the format defined by the fmt.Printf line at the bottom, as we will define a passing coding challenge as one that outputs
// the correct information, including this line
func main() {
	// Github
	if len(os.Args) < 2 {
        fmt.Println("Missing parameter, provide file name!")
        return
    }
	content := readFile(os.Args[1])
	getRespositories(content)
}