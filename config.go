package main

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"strconv"
)

type configuration struct {
	RootDomain       string `yaml:"root_domain,omitempty"`
	Timezone         string `yaml:"timezone,omitempty"`
	LocalEnvironment bool   `yaml:"local_environment,omitempty"`
	Security         struct {
		RootDomainPolicy string `yaml:"root_domain_policy,omitempty"`
	} `yaml:"security,omitempty"`
	Smtp struct {
		Username       string `yaml:"username,omitempty"`
		Password       string `yaml:"password,omitempty"`
		Host           string `yaml:"host,omitempty"`
		Port           int    `yaml:"port,omitempty"`
		Sender         string `yaml:"sender,omitempty"`
		StartupAddress string `yaml:"startup_address,omitempty"`
	} `yaml:"smtp,omitempty"`
	TraefikHub struct {
		Enabled bool   `yaml:"enabled,omitempty"`
		Key     string `yaml:"key,omitempty"`
	} `yaml:"traefik_hub,omitempty"`
	Cloudflare struct {
		Enabled bool   `yaml:"enabled,omitempty"`
		Email   string `yaml:"email,omitempty"`
		ApiKey  string `yaml:"api_key,omitempty"`
	} `yaml:"cloudflare,omitempty"`
}

/*
Write config to file
*/
func (c *configuration) writeConfig() {
	yamlData, err := yaml.Marshal(&c)

	if err != nil {
		fmt.Printf("Error while Marshaling. %v", err)
	}

	writeFile("AppData", "traefik-installer.yaml", yamlData, 0644)
}

/*
Read config file
*/
func (c *configuration) readConfig() *configuration {
	content := readFile("AppData", "traefik-installer.yaml")

	err := yaml.Unmarshal(content, c)

	if err != nil {
		return c
	}

	return c
}

/*
list configuration with println
*/
func (c *configuration) toText() {
	fmt.Println("Root-Domain: " + c.RootDomain)

	fmt.Println("Timezone: " + c.Timezone)

	if c.LocalEnvironment {
		fmt.Println("Is local environment: yes")
	} else {
		fmt.Println("Is local environment: no")
	}

	fmt.Println("Root-Domain-Policy: " + c.Security.RootDomainPolicy)
	fmt.Println("Smtp-Username: " + c.Smtp.Username)
	fmt.Println("Smtp-Password: " + c.Smtp.Password)
	fmt.Println("Smtp-Host: " + c.Smtp.Host)
	fmt.Println("Smtp-Port: " + strconv.Itoa(c.Smtp.Port))
	fmt.Println("Smtp-Sender: " + c.Smtp.Sender)
	fmt.Println("Smtp-Startup Address: " + c.Smtp.StartupAddress)

	if c.TraefikHub.Enabled {
		fmt.Println("Traefik-Hub enabled: yes")
		fmt.Println("Traefik-Hub Key: " + c.TraefikHub.Key)
	} else {
		fmt.Println("Traefik-Hub enabled: no")
	}

	if c.Cloudflare.Enabled {
		fmt.Println("Cloudflare enabled: yes")
		fmt.Println("Cloudflare Email: " + c.Cloudflare.Email)
		fmt.Println("Cloudflare Api-Key: " + c.Cloudflare.ApiKey)
	} else {
		fmt.Println("Cloudflare enabled: no")
	}
}
