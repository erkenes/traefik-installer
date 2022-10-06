package main

import (
	"fmt"
)

var appVersion = "1.0.0"

func main() {

	var rootDomain string
	var useTraefikHub = false
	var useCloudflare = false
	var timezone = "Europe/Berlin"
	var cloudflareEmail = ""
	var cloudflareApiKey = ""
	var traefikHubKey = ""
	var secureRootDomain bool
	var secureExtendedRootDomain bool
	var smtpUsername string
	var smtpPassword string
	var smtpHost string
	var smtpPort int
	var smtpSender string
	var smtpStartupAddress string
	var policyForRootDomain = "bypass"
	var isLocalEnvironment = true

	fmt.Println("#################")
	fmt.Println("")
	fmt.Println("traefik installer version: " + appVersion)
	fmt.Println("")
	fmt.Println("#################")

	fmt.Println("What is your Timezone?")
	fmt.Scan(&timezone)

	fmt.Println("Is this a local environment (not exposed to the internet)?")
	isLocalEnvironment = askForConfirmation()

	if isLocalEnvironment == false {
		fmt.Println("What is your root domain?")
		fmt.Scan(&rootDomain)

		fmt.Println("Do you want to use Cloudflare?")
		useCloudflare = askForConfirmation()

		if useCloudflare == true {
			fmt.Println("What is your Cloudflare email address?")
			fmt.Scan(&cloudflareEmail)

			fmt.Println("What is your Cloudflare Global API-Key?")
			fmt.Scan(&cloudflareApiKey)
		}

		fmt.Println("Do you want to use Traefik-Hub?")
		useTraefikHub = askForConfirmation()

		if useTraefikHub == true {
			fmt.Println("What is your Traefik-Hub Key? Create a new one here https://hub.traefik.io/agents/new")
			fmt.Scan(&traefikHubKey)
		}
	} else {
		fmt.Println("Your root domain is: local.dev")
		fmt.Println("We use a default cert file. You can find it in the 'certs' folder. Please install the RootCA.crt file.\nAll urls must have the tld '*.local.dev'")

		fmt.Println("\nCloudflare is disabled")
		fmt.Println("\nTraefik-Hub is disabled")

		rootDomain = "local.dev"
		useCloudflare = false
		useTraefikHub = false
	}

	fmt.Println("")
	fmt.Println(" --- Authelia Configuration --- ")

	fmt.Println("Should the root domain (" + rootDomain + ") be secured with authelia?")
	secureRootDomain = askForConfirmation()

	if secureRootDomain == true {
		fmt.Println("Should Two-Factor be used instead of One-Factor authentication for the root domain (" + rootDomain + ")?")
		secureExtendedRootDomain = askForConfirmation()

		if secureExtendedRootDomain == true {
			policyForRootDomain = "two_factor"
		} else {
			policyForRootDomain = "one_factor"
		}
	} else {
		policyForRootDomain = "bypass"
	}

	fmt.Println("")
	fmt.Println(" --- SMTP-Settings for Authelia --- ")

	fmt.Println("SMTP-Host:")
	fmt.Scan(&smtpHost)

	fmt.Println("SMTP-Port:")
	fmt.Scan(&smtpPort)

	fmt.Println("SMTP-Username:")
	fmt.Scan(&smtpUsername)

	fmt.Println("SMTP-Password:")
	fmt.Scan(&smtpPassword)

	fmt.Println("SMTP-Sender Address:")
	fmt.Scan(&smtpSender)

	fmt.Println("SMTP-Startup Check Address:")
	fmt.Scan(&smtpStartupAddress)

	traefikConfig := createTraefikFile(rootDomain, useCloudflare, cloudflareEmail, useTraefikHub)
	autheliaConfig := createAutheliaConfig(rootDomain, smtpUsername, smtpHost, smtpPort, smtpSender, smtpStartupAddress, policyForRootDomain)

	writeFile("AppData/traefik-proxy", "acme.json", []byte(""), 0600)
	writeFile("AppData/traefik-proxy", "traefik.yml", traefikConfig, 0644)
	writeFile("AppData/authelia", "configuration.yml", autheliaConfig, 0644)

	// Secrets
	writeFile("secrets", "authelia_notifier_smtp_password", []byte(smtpPassword), 0600)
	writeFile("secrets", "authelia_jwt_secret", []byte(randomString(32)), 0600)
	writeFile("secrets", "authelia_session_secret", []byte(randomString(32)), 0600)
	writeFile("secrets", "authelia_storage_encryption_key", []byte(randomString(32)), 0600)
	writeFile("secrets", "mysql_password", []byte(randomString(16)), 0600)
	writeFile("secrets", "mysql_root_password", []byte(randomString(16)), 0600)

	if useCloudflare == true {
		writeFile("secrets", "cf_email", []byte(cloudflareEmail), 0600)
		writeFile("secrets", "cf_api_key", []byte(cloudflareApiKey), 0600)
	}
	if isLocalEnvironment == true {
		traefikTlsDynamicConfig := createTraefikTlsFile()

		writeFile("AppData/traefik-proxy/rules", "local-tls.yml", traefikTlsDynamicConfig, 0644)
	}

	dockerComposeConfig := createDockerComposeFile(useTraefikHub, useCloudflare, isLocalEnvironment, secureRootDomain)
	writeFile(".", "docker-compose.yml", dockerComposeConfig, 0644)

	writeFile(".", ".env", []byte(createEnvFile(rootDomain, timezone, useTraefikHub, traefikHubKey)), 0700)
}
