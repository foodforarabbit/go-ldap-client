package main

import (
	"flag"
	"log"
	"os"
  "encoding/json"
	"fmt"

	"github.com/foodforarabbit/go-ldap-client"
)

var base, bindDN, bindPassword, groupFilter, host, password, serverName, userFilter, username string
var port int
var useSSL bool
var skipTLS bool

var ldapConfiguration ldap.LDAPClient
var configuration Configuration

type UserConfiguration struct {
		Username     *string `json:"Username"`
		Password     *string `json:"Password"`
}

// configuration file
type Configuration struct {
    User    	  *UserConfiguration   `json:"User"`
}

type server struct{}

func main() {

  configurationFile := flag.String("c", "configuration.json", "configuration file path")
	flag.Parse()
	file, _ := os.Open(*configurationFile)

  // decode configuration json
	decoder := json.NewDecoder(file)
	ldapConfiguration = ldap.LDAPClient{}
	configuration = Configuration{}
	err := decoder.Decode(&ldapConfiguration)
	if err != nil {
	  fmt.Println("error1:", err)
	}

	file2, _ := os.Open(*configurationFile)
	decoder2 := json.NewDecoder(file2)
	err = decoder2.Decode(&configuration)
	if err != nil {
	  fmt.Println("error2:", err)
	}

	client := &ldapConfiguration
	defer client.Close()

	username = *(configuration.User.Username)
	password = *(configuration.User.Password)

	ok, user, err := client.Authenticate(username, password)
	if err != nil {
		log.Fatalf("Error authenticating user %s: %+v", username, err)
	}
	if !ok {
		log.Fatalf("Authenticating failed for user %s", username)
	}
	log.Printf("User: %+v", user)

	groups, err := client.GetGroupsOfUser(username)
	if err != nil {
		log.Fatalf("Error getting groups for user %s: %+v", username, err)
	}

	users, err := client.GetUsers()
	if err != nil {
		log.Fatalf("Error getting groups for user %s: %+v", username, err)
	}
	for _, user2 := range users {
		log.Printf("allUser: %+v", user2)
	}
	log.Printf("Groups: %+v", groups)
}
