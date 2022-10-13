package main

import (
	"fmt"
)

var appVersion = "1.0.0"

var configRootDomain string
var configUseTraefikHub = false
var configUseCloudflare = false
var configTimezone = "Europe/Berlin"
var configCloudflareEmail = ""
var configCloudflareApiKey = ""
var configTraefikHubKey = ""
var configSecureRootDomain bool
var configSecureExtendedRootDomain bool
var configSmtpUsername string
var configSmtpPassword string
var configSmtpHost string
var configSmtpPort int
var configSmtpSender string
var configSmtpStartupAddress string
var configPolicyForRootDomain = "bypass"
var configIsLocalEnvironment = true

func main() {
	showStartMessage()

	printSectionHeader("Read Configuration")
	var c configuration
	c.readConfig()

	configTimezone = getTextInput("What is your "+colorYellow+"Timezone"+colorReset+"?", true, c.Timezone)
	configIsLocalEnvironment = getConfirmInput("Is this a "+colorYellow+"local environment"+colorReset+" (not exposed to the internet)?", c.LocalEnvironment)

	if configIsLocalEnvironment == false {
		configRootDomain = getTextInput("What is your "+colorYellow+"root domain"+colorReset+"?", true, c.RootDomain)
		configUseCloudflare = getConfirmInput("Do you want to "+colorYellow+"use Cloudflare"+colorReset+"?", c.Cloudflare.Enabled)

		if configUseCloudflare == true {
			configCloudflareEmail = getTextInput("What is your "+colorYellow+"Cloudflare email address"+colorReset+"?", true, c.Cloudflare.Email)
			configCloudflareApiKey = getTextInput("What is your "+colorYellow+"Cloudflare Global API-Key"+colorReset+"?", true, c.Cloudflare.ApiKey)
		}

		configUseTraefikHub = getConfirmInput("Do you want to "+colorYellow+"use Traefik-Hub"+colorReset+"?", c.TraefikHub.Enabled)

		if configUseTraefikHub == true {
			configTraefikHubKey = getTextInput("What is your "+colorYellow+"Traefik-Hub Key"+colorReset+"? Create a new one here "+colorCyan+"https://hub.traefik.io/agents/new"+colorReset, true, c.TraefikHub.Key)
		}
	} else {
		fmt.Println("Your root domain is: " + colorPurple + "local.dev" + colorReset)
		fmt.Println("We use a default cert file. You can find it in the '" + colorBlue + "certs" + colorReset + "' folder. Please install the RootCA.crt file.\nAll urls must have the tld '" + colorBlue + "*.local.dev" + colorReset + "'")

		fmt.Println(colorGreen + "\nCloudflare is disabled" + colorReset)
		fmt.Println(colorGreen + "\nTraefik-Hub is disabled" + colorReset)

		configRootDomain = "local.dev"
		configUseCloudflare = false
		configUseTraefikHub = false
	}

	printSectionHeader("Authelia Configuration")

	securedRootDomain := false
	securedRootDomainExtended := false
	if c.Security.RootDomainPolicy != "bypass" {
		securedRootDomain = true

		if c.Security.RootDomainPolicy == "two_factor" {
			securedRootDomainExtended = true
		}
	}
	configSecureRootDomain = getConfirmInput("Should the root domain ("+configRootDomain+") be "+colorYellow+"secured with authelia"+colorReset+"?", securedRootDomain)

	if configSecureRootDomain == true {
		configSecureExtendedRootDomain = getConfirmInput("Should "+colorYellow+"Two-Factor"+colorReset+" be used instead of One-Factor authentication for the root domain ("+configRootDomain+")?", securedRootDomainExtended)

		if configSecureExtendedRootDomain == true {
			configPolicyForRootDomain = "two_factor"
		} else {
			configPolicyForRootDomain = "one_factor"
		}
	} else {
		configPolicyForRootDomain = "bypass"
	}

	// SMTP-Settings
	printSectionHeader("SMTP-Settings for Authelia")
	configSmtpHost = getTextInput(colorYellow+"SMTP-Host"+colorReset+":", true, c.Smtp.Host)
	configSmtpPort = getNumberInput(colorYellow+"SMTP-Port"+colorReset+":", true, c.Smtp.Port)
	configSmtpUsername = getTextInput(colorYellow+"SMTP-Username"+colorReset+":", true, c.Smtp.Username)
	configSmtpPassword = getPasswordInput(colorYellow+"SMTP-Password"+colorReset+":", true, c.Smtp.Password)
	configSmtpSender = getEmailAddressInput(colorYellow+"SMTP-Sender Address"+colorReset+":", true, c.Smtp.Sender)
	configSmtpStartupAddress = getEmailAddressInput(colorYellow+"SMTP-Startup Check Address"+colorReset+":", true, c.Smtp.StartupAddress)

	createTraefikFile()
	createAutheliaConfig()
	createDockerComposeFile()
	writeFile(".", ".env", []byte(createEnvFile(configRootDomain, configTimezone, configUseTraefikHub, configTraefikHubKey)), 0700)

	// Update config
	c.RootDomain = configRootDomain
	c.Timezone = configTimezone
	c.LocalEnvironment = configIsLocalEnvironment
	c.Security.RootDomainPolicy = configPolicyForRootDomain
	c.Smtp.Username = configSmtpUsername
	c.Smtp.Password = configSmtpPassword
	c.Smtp.Host = configSmtpHost
	c.Smtp.Port = configSmtpPort
	c.Smtp.Sender = configSmtpSender
	c.Smtp.StartupAddress = configSmtpStartupAddress
	c.TraefikHub.Enabled = configUseTraefikHub
	c.TraefikHub.Key = configTraefikHubKey
	c.Cloudflare.Enabled = configUseCloudflare
	c.Cloudflare.Email = configCloudflareEmail
	c.Cloudflare.ApiKey = configCloudflareApiKey
	c.writeConfig()
}

func showStartMessage() {
	fmt.Println("#################")
	fmt.Println("")
	fmt.Println("traefik installer version: " + colorCyan + appVersion + colorReset)
	fmt.Println("")
	fmt.Println("#################")
}
