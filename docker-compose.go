package main

import (
	"fmt"
	"github.com/imdario/mergo"
	"gopkg.in/yaml.v3"
)

/*
Create the docker compose file
*/
func createDockerComposeFile() {

	type IpamConfigSubnetStruct struct {
		Subnet string `yaml:"subnet"`
	}
	type IpamStruct struct {
		Config []IpamConfigSubnetStruct `yaml:"config"`
	}
	type NetworkStruct struct {
		Name   string     `yaml:"name"`
		Driver string     `yaml:"driver"`
		Ipam   IpamStruct `yaml:"ipam"`
	}
	type NetworksStruct struct {
		TraefikInternal NetworkStruct `yaml:"traefik-internal"`
		TraefikProxy    NetworkStruct `yaml:"traefik-proxy"`
		SocketProxy     NetworkStruct `yaml:"socket-proxy"`
	}

	type SecretStruct struct {
		File string `yaml:"file"`
	}
	type SecretsStruct struct {
		CfEmail                      SecretStruct `yaml:"cf_email,omitempty"`
		CfApiKey                     SecretStruct `yaml:"cf_api_key,omitempty"`
		AutheliaJwtSecret            SecretStruct `yaml:"authelia_jwt_secret"`
		AutheliaSessionSecret        SecretStruct `yaml:"authelia_session_secret"`
		AutheliaNotifierSmtpPassword SecretStruct `yaml:"authelia_notifier_smtp_password"`
		AutheliaStorageEncryptionKey SecretStruct `yaml:"authelia_storage_encryption_key"`
		MysqlRootPassword            SecretStruct `yaml:"mysql_root_password"`
		MysqlPassword                SecretStruct `yaml:"mysql_password"`
	}

	type NetworkIpV4AddressStruct struct {
		Ipv4Address string `yaml:"ipv4_address,omitempty"`
	}
	type AllNetworksStruct struct {
		TraefikInternal NetworkIpV4AddressStruct `yaml:"traefik-internal,omitempty"`
		TraefikProxy    NetworkIpV4AddressStruct `yaml:"traefik-proxy,omitempty"`
		SocketProxy     NetworkIpV4AddressStruct `yaml:"socket-proxy,omitempty"`
	}

	type HealthcheckStruct struct {
		Test        []string `yaml:"test,flow"`
		Interval    string   `yaml:"interval,omitempty"`
		Retries     int      `yaml:"retries,omitempty"`
		Timeout     string   `yaml:"timeout,omitempty"`
		StartPeriod string   `yaml:"start_period,omitempty"`
	}

	type ServiceStruct struct {
		ContainerName string            `yaml:"container_name,omitempty"`
		Image         string            `yaml:"image"`
		Restart       string            `yaml:"restart"`
		Command       []string          `yaml:"command,omitempty"`
		DependsOn     []string          `yaml:"depends_on,omitempty"`
		Networks      AllNetworksStruct `yaml:"networks"`
		SecurityOpt   []string          `yaml:"security_opt,omitempty"`
		Healthcheck   HealthcheckStruct `yaml:"healthcheck,omitempty"`
		Ports         []string          `yaml:"ports,omitempty"`
		Volumes       []string          `yaml:"volumes,omitempty"`
		Environment   []string          `yaml:"environment,omitempty"`
		Secrets       []string          `yaml:"secrets,omitempty"`
		Labels        []string          `yaml:"labels,omitempty"`
	}
	type ServicesStruct struct {
		TraefikProxy ServiceStruct `yaml:"traefik-proxy"`
		TraefikHub   ServiceStruct `yaml:"traefik-hub,omitempty"`
		SocketProxy  ServiceStruct `yaml:"socket-proxy"`
		MariaDb      ServiceStruct `yaml:"mariadb"`
		Authelia     ServiceStruct `yaml:"authelia"`
		Redis        ServiceStruct `yaml:"redis"`
	}

	type DockerComposeFile struct {
		Version  string         `yaml:"version"`
		Networks NetworksStruct `yaml:"networks"`
		Secrets  SecretsStruct  `yaml:"secrets"`
		Services ServicesStruct `yaml:"services"`
	}

	// network traefik-internal
	subnetTraefikInternal := "192.169.1.0/24"
	traefikProxyNetTraefikInternal := "192.169.1.11"
	autheliaNetTraefikInternal := "192.169.1.12"
	mariaDbNetTraefikInternal := "192.169.1.13"
	redisNetTraefikInternal := "192.169.1.14"
	traefikHubNetTraefikInternal := "192.169.1.15"

	// network socket-proxy
	subnetSocketProxy := "192.169.2.0/24"
	traefikProxyNetSocketProxy := "192.169.2.11"
	socketProxyNetSocketProxy := "192.169.2.12"

	// network traefik-proxy
	subnetTraefikProxy := "192.169.3.0/24"
	traefikProxyNetTraefikProxy := "192.169.3.11"
	autheliaNetTraefikProxy := "192.169.3.12"

	traefikProxyPorts := []string{
		"0.0.0.0:80:80",
		"0.0.0.0:443:443",
		//"127.0.0.1:8080:8080",
	}
	traefikProxyEnvironment := []string{
		"ROOT_DOMAIN_NAME=${ROOT_DOMAIN_NAME}",
	}
	var traefikProxySecrets []string
	traefikProxyVolumes := []string{
		"./AppData/traefik-proxy:/etc/traefik",
	}
	traefikProxyLabels := []string{
		"traefik.enable=true",
		"traefik.http.routers.http-catchall.entrypoints=http",
		"traefik.http.routers.http-catchall.rule=HostRegexp(`{host:.+}`)",
		"traefik.http.routers.http-catchall.middlewares=middlewares-https-redirectscheme@file",
		"traefik.http.routers.traefik-rtr.tls=true",
		"traefik.http.routers.traefik-rtr.entrypoints=https",
		"traefik.http.routers.traefik-rtr.rule=Host(`traefik.${ROOT_DOMAIN_NAME}`)",
		"traefik.http.routers.traefik-rtr.service=api@internal",
		"traefik.http.routers.ping.rule=Host(`traefik.${ROOT_DOMAIN_NAME}`) && Path(`/ping`)",
		"traefik.http.routers.ping.tls=true",
		"traefik.http.routers.ping.service=ping@internal",
	}
	traefikHubService := ServiceStruct{}
	secrets := SecretsStruct{
		AutheliaJwtSecret: SecretStruct{
			File: "./secrets/authelia_jwt_secret",
		},
		AutheliaSessionSecret: SecretStruct{
			File: "./secrets/authelia_session_secret",
		},
		AutheliaNotifierSmtpPassword: SecretStruct{
			File: "./secrets/authelia_notifier_smtp_password",
		},
		AutheliaStorageEncryptionKey: SecretStruct{
			File: "./secrets/authelia_storage_encryption_key",
		},
		MysqlRootPassword: SecretStruct{
			File: "./secrets/mysql_root_password",
		},
		MysqlPassword: SecretStruct{
			File: "./secrets/mysql_password",
		},
	}

	if configUseTraefikHub == true {
		traefikHubServiceDef := ServiceStruct{
			ContainerName: "traefik-hub-agent",
			Image:         "ghcr.io/traefik/hub-agent-traefik:v0.7.2",
			Restart:       "always",
			Command: []string{
				"run",
				"--hub.token=${TRAEFIK_HUB_CODE}",
				"--auth-server.advertise-url=http://traefik-hub-agent",
				"--traefik.host=traefik-proxy",
				"--traefik.tls.insecure=true",
			},
			Volumes: []string{
				"/var/run/docker.sock:/var/run/docker.sock",
			},
			DependsOn: []string{
				"traefik-proxy",
			},
			Networks: AllNetworksStruct{
				TraefikInternal: NetworkIpV4AddressStruct{
					traefikHubNetTraefikInternal,
				},
			},
		}

		mergo.Merge(&traefikHubService, traefikHubServiceDef)

		traefikProxyPorts = append(traefikProxyPorts, "0.0.0.0:9900:9900", "0.0.0.0:9901:9901")
	}

	if configUseCloudflare == true {
		cloudflareSecrets := SecretsStruct{
			CfEmail: SecretStruct{
				File: "./secrets/cf_email",
			},
			CfApiKey: SecretStruct{
				File: "./secrets/cf_api_key",
			},
		}

		traefikProxyEnvironment = append(traefikProxyEnvironment, "CF_API_EMAIL_FILE=/run/secrets/cf_email", "CF_API_KEY_FILE=/run/secrets/cf_api_key")
		traefikProxySecrets = append(traefikProxySecrets, "cf_email", "cf_api_key")

		mergo.Merge(&secrets, cloudflareSecrets)
	}

	if configIsLocalEnvironment == true {
		traefikProxyVolumes = append(traefikProxyVolumes, "./certs:/certs:ro")
	}

	if configSecureRootDomain == true {
		traefikProxyLabels = append(traefikProxyLabels, "traefik.http.routers.traefik-rtr.middlewares=chain-authelia@file")
	}

	dockerComposeConfig := DockerComposeFile{
		Version: "3.9",
		Networks: NetworksStruct{
			TraefikInternal: NetworkStruct{
				Name:   "traefik-internal",
				Driver: "bridge",
				Ipam: IpamStruct{
					Config: []IpamConfigSubnetStruct{
						{
							subnetTraefikInternal,
						},
					},
				},
			},
			SocketProxy: NetworkStruct{
				Name:   "socket-proxy",
				Driver: "bridge",
				Ipam: IpamStruct{
					Config: []IpamConfigSubnetStruct{
						{
							subnetSocketProxy,
						},
					},
				},
			},
			TraefikProxy: NetworkStruct{
				Name:   "traefik-proxy",
				Driver: "bridge",
				Ipam: IpamStruct{
					Config: []IpamConfigSubnetStruct{
						{
							subnetTraefikProxy,
						},
					},
				},
			},
		},
		Secrets: secrets,
		Services: ServicesStruct{
			TraefikProxy: ServiceStruct{
				ContainerName: "traefik-proxy",
				Image:         "traefik:2.8",
				Restart:       "always",
				Networks: AllNetworksStruct{
					TraefikInternal: NetworkIpV4AddressStruct{
						Ipv4Address: traefikProxyNetTraefikInternal,
					},
					SocketProxy: NetworkIpV4AddressStruct{
						Ipv4Address: traefikProxyNetSocketProxy,
					},
					TraefikProxy: NetworkIpV4AddressStruct{
						Ipv4Address: traefikProxyNetTraefikProxy,
					},
				},
				SecurityOpt: []string{
					"no-new-privileges:true",
				},
				Healthcheck: HealthcheckStruct{
					Test: []string{
						"CMD",
						"traefik",
						"healthcheck",
						"--ping",
					},
					Interval: "5s",
					Retries:  3,
				},
				Ports:       traefikProxyPorts,
				Volumes:     traefikProxyVolumes,
				Environment: traefikProxyEnvironment,
				Labels:      traefikProxyLabels,
				Secrets:     traefikProxySecrets,
			},
			SocketProxy: ServiceStruct{
				ContainerName: "socket-proxy",
				Image:         "tecnativa/docker-socket-proxy:0.1.1",
				Restart:       "always",
				Networks: AllNetworksStruct{
					SocketProxy: NetworkIpV4AddressStruct{
						socketProxyNetSocketProxy,
					},
				},
				SecurityOpt: []string{
					"no-new-privileges:true",
				},
				Volumes: []string{
					"/var/run/docker.sock:/var/run/docker.sock",
				},
				Environment: []string{
					"LOG_LEVEL=info",
					"EVENTS=1",
					"PING=1",
					"VERSION=1",
					"AUTH=0",
					"SECRETS=0",
					"POST=0",
					"BUILD=0",
					"COMMIT=0",
					"CONFIGS=0",
					"CONTAINERS=1",
					"DISTRIBUTION=0",
					"EXEC=0",
					"IMAGES=0",
					"INFO=0",
					"NETWORKS=0",
					"NODES=0",
					"PLUGINS=0",
					"SERVICES=0",
					"SESSION=0",
					"SWARM=0",
					"SYSTEM=0",
					"TASKS=0",
					"VOLUMES=0",
				},
				Healthcheck: HealthcheckStruct{
					Test: []string{
						"wget --spider http://localhost:2375/version || exit 1",
					},
					Interval:    "29s",
					Timeout:     "5s",
					Retries:     3,
					StartPeriod: "21s",
				},
			},
			MariaDb: ServiceStruct{
				Image: "linuxserver/mariadb:10.5.17",
				SecurityOpt: []string{
					"no-new-privileges:true",
				},
				Restart: "always",
				Ports: []string{
					"127.0.0.1:3306:3306",
				},
				Volumes: []string{
					"./AppData/MariaDB/data:/config",
				},
				Networks: AllNetworksStruct{
					TraefikInternal: NetworkIpV4AddressStruct{
						mariaDbNetTraefikInternal,
					},
				},
				Environment: []string{
					"TZ=${TIMEZONE}",
					"PUID=${PUID}",
					"PGID=${PGID}",
					"MYSQL_DATABASE=authelia",
					"MYSQL_USER=authelia",
					"FILE__MYSQL_PASSWORD=/run/secrets/mysql_password",
					"FILE__MYSQL_ROOT_PASSWORD=/run/secrets/mysql_root_password",
				},
				Secrets: []string{
					"mysql_root_password",
					"mysql_password",
				},
			},
			Authelia: ServiceStruct{
				Image:   "authelia/authelia:4.36.9",
				Restart: "always",
				Networks: AllNetworksStruct{
					TraefikProxy: NetworkIpV4AddressStruct{
						autheliaNetTraefikProxy,
					},
					TraefikInternal: NetworkIpV4AddressStruct{
						autheliaNetTraefikInternal,
					},
				},
				Volumes: []string{
					"./AppData/authelia:/config",
				},
				Environment: []string{
					"TZ=${TIMEZONE}",
					"AUTHELIA_JWT_SECRET_FILE=/run/secrets/authelia_jwt_secret",
					"AUTHELIA_SESSION_SECRET_FILE=/run/secrets/authelia_session_secret",
					"AUTHELIA_STORAGE_MYSQL_PASSWORD_FILE=/run/secrets/mysql_password",
					"AUTHELIA_NOTIFIER_SMTP_PASSWORD_FILE=/run/secrets/authelia_notifier_smtp_password",
					"AUTHELIA_STORAGE_ENCRYPTION_KEY_FILE=/run/secrets/authelia_storage_encryption_key",
				},
				Secrets: []string{
					"authelia_jwt_secret",
					"authelia_session_secret",
					"mysql_password",
					"authelia_notifier_smtp_password",
					"authelia_storage_encryption_key",
				},
				Labels: []string{
					"traefik.enable=true",
					"traefik.http.routers.authelia-rtr.entrypoints=https",
					"traefik.http.routers.authelia-rtr.rule=HostHeader(`authelia.${ROOT_DOMAIN_NAME}`)",
					"traefik.http.routers.authelia-rtr.tls=true",
					"traefik.http.routers.authelia-rtr.middlewares=chain-authelia@file",
					"traefik.http.routers.authelia-rtr.service=authelia-svc",
					"traefik.http.services.authelia-svc.loadbalancer.server.port=9091"},
			},
			Redis: ServiceStruct{
				Image: "redis:7.0.5-alpine",
				SecurityOpt: []string{
					"no-new-privileges:true",
				},
				Restart: "always",
				Networks: AllNetworksStruct{
					TraefikInternal: NetworkIpV4AddressStruct{
						redisNetTraefikInternal,
					},
				},
				Healthcheck: HealthcheckStruct{
					Test: []string{
						"CMD",
						"redis-cli",
						"--raw",
						"incr",
						"ping",
					},
				},
				Command: []string{
					"redis-server",
					"--appendonly",
					"yes",
					"--maxmemory",
					"5124mb",
					"--maxmemory-policy",
					"allkeys-lru",
				},
			},
			TraefikHub: traefikHubService,
		},
	}

	yamlData, err := yaml.Marshal(&dockerComposeConfig)

	if err != nil {
		fmt.Printf("Error while Marshaling. %v", err)
	}

	writeFile(".", "docker-compose.yml", yamlData, 0644)
}
