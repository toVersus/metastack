package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/servers"
	"github.com/gophercloud/gophercloud/pagination"
	"github.com/kelseyhightower/envconfig"
)

type AuthConfig struct {
	RegionName         string `split_words:"true"`
	AuthUrl            string `split_words:"true"`
	Username           string
	Password           string
	ProjectId          string `split_words:"true"`
	UserDomainName     string `split_words:"true"`
	IdentityApiVersion string `split_words:"true"`
	Interface          string
}

func main() {
	flag.Parse()
	args := flag.Args()
	if len(args) > 1 {
		log.Fatal("Too many arguments. Please specify single server name")
	} else if len(args) == 0 {
		log.Fatal("Too few arguments. Please specify single server name")
	}
	serverName := args[0]

	var config AuthConfig
	err := envconfig.Process("OS", &config)
	if err != nil {
		log.Fatal(err)
	}

	authOpts := gophercloud.AuthOptions{
		IdentityEndpoint: config.AuthUrl,
		Username:         config.Username,
		Password:         config.Password,
		DomainName:       config.UserDomainName,
		Scope: &gophercloud.AuthScope{
			ProjectID: config.ProjectId,
		},
	}
	provider, err := openstack.AuthenticatedClient(authOpts)
	if err != nil {
		log.Fatalf("Failed to initialize OpenStack provider with auth options (%#v): %s", authOpts, err)
	}

	client, err := openstack.NewComputeV2(provider, gophercloud.EndpointOpts{
		Region: config.RegionName,
	})
	if err != nil {
		log.Fatalf("Failed to initialize Openstack client: %s", err)
	}

	lsOpts := servers.ListOpts{Name: serverName}
	if err := servers.List(client, lsOpts).EachPage(func(page pagination.Page) (bool, error) {
		serverList, err := servers.ExtractServers(page)
		if err != nil {
			return false, err
		}
		if len(serverList) == 0 {
			return false, fmt.Errorf("server not found: %s", serverName)
		}

		// Only extract metadata of first matched server
		result, err := servers.Metadata(client, serverList[0].ID).Extract()
		if err != nil {
			return false, err
		}

		bytes, err := json.Marshal(result)
		if err != nil {
			return false, fmt.Errorf("failed to encode metadata to JSON: %s", err)
		}
		fmt.Println(string(bytes))

		return true, nil
	}); err != nil {
		log.Fatal(err)
	}
}
