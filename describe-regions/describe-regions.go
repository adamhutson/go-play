package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

var profiles = []string{}
var profile string

func main() {
	loadProfiles()
	profileFlag := flag.String("profile", "default", availableProfiles())
	flag.Parse()
	profile = *profileFlag
	if !validProfile(profile) {
		log.Fatal(availableProfiles())
	}
	describeRegions()
}

func loadProfiles() {
	home, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	file, err := os.Open(filepath.Join(home, ".aws", "credentials"))
	if err != nil {
		panic(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "[") {
			profile := strings.Replace(line, "[", "", 1)
			profile = strings.Replace(profile, "]", "", 1)
			profiles = append(profiles, profile)
		}
		if err := scanner.Err(); err != nil {
			log.Fatal(err)
		}
	}
}

func availableProfiles() string {
	validProfiles := strings.Join(profiles, ",")
	return fmt.Sprintf("Available profiles are: %v", validProfiles)
}

func validProfile(profileToValidate string) bool {
	for _, p := range profiles {
		if strings.ToLower(p) == strings.ToLower(profileToValidate) {
			return true
		}
	}
	return false
}

func describeRegions() {
	sess, err := session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
		Profile:           profile,
	})

	svc := ec2.New(sess)
	regions, err := svc.DescribeRegions(nil)
	if err != nil {
		panic(err)
	}
	bs, err := json.Marshal(regions)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(bs))
}
