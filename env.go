package main

/*
Create env file
*/
func createEnvFile(rootDomain string, timezone string, useTraefikHub bool, traefikHubKey string) string {
	var env string = "PUID=1000\n" +
		"PGID=1000\n" +
		"TIMEZONE=" + timezone + "\n" +
		"ROOT_DOMAIN_NAME=" + rootDomain

	if useTraefikHub == true {
		env = env + "\nTRAEFIK_HUB_CODE=" + traefikHubKey
	}

	return env
}
