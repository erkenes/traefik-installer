package main

import (
	"fmt"
	"github.com/imdario/mergo"
	"gopkg.in/yaml.v3"
)

type GlobalStruct struct {
	CheckNewVersion    bool `yaml:"checkNewVersion"`
	SendAnonymousUsage bool `yaml:"sendAnonymousUsage"`
}
type EntryPointHttpStruct struct {
	Address string `yaml:"address"`
}
type HttpTlsDomainsStruct struct {
	Main string   `yaml:"main"`
	Sans []string `yaml:"sans"`
}
type TlsStruct struct {
	CertResolver string                 `yaml:"certResolver,omitempty"`
	Domains      []HttpTlsDomainsStruct `yaml:"domains"`
}
type HttpsHttpStruct struct {
	Tls TlsStruct
}
type ForwardedHeadersStruct struct {
	TrustedIPs []string `yaml:"trustedIPs,omitempty"`
}
type EntryPointHttpsStruct struct {
	Address          string                 `yaml:"address"`
	ForwardedHeaders ForwardedHeadersStruct `yaml:"forwardedHeaders,omitempty"`
	Http             HttpsHttpStruct        `yaml:"http"`
}
type EntryPointPing struct {
	Address string `yaml:"address"`
}
type EntryPoints struct {
	Http  EntryPointHttpStruct  `yaml:"http"`
	Https EntryPointHttpsStruct `yaml:"https"`
	Ping  EntryPointPing        `yaml:"ping"`
}
type ExperimentalStruct struct {
	Hub bool `yaml:"hub"`
}

type HubTlsStruct struct {
	Insecure bool `yaml:"insecure"`
}
type HubStruct struct {
	Tls HubTlsStruct `yaml:"tls"`
}

type PrometheusStruct struct {
	AddRoutersLabels bool `yaml:"addRoutersLabels"`
}
type MetricsStruct struct {
	Prometheus PrometheusStruct `yaml:"prometheus"`
}

type ApiStruct struct {
	Insecure  bool `yaml:"insecure"`
	Dashboard bool `yaml:"dashboard"`
}

type PingStruct struct {
	EntryPoint string `yaml:"entryPoint"`
}

type LogStruct struct {
	Level    string `yaml:"level"`
	FilePath string `yaml:"filePath"`
}

type AccessLogFiltersStruct struct {
	StatusCodes []string `yaml:"statusCodes"`
}
type AccessLogStruct struct {
	FilePath      string                 `yaml:"filePath"`
	BufferingSize int                    `yaml:"bufferingSize"`
	Filters       AccessLogFiltersStruct `yaml:"filters"`
}

type ProvidersDockerStruct struct {
	Endpoint         string `yaml:"endpoint"`
	ExposedByDefault bool   `yaml:"exposedByDefault"`
	Network          string `yaml:"network"`
	SwarmMode        bool   `yaml:"swarmMode"`
}
type ProvidersFileStruct struct {
	Directory string `yaml:"directory"`
	Watch     bool   `yaml:"watch"`
}
type ProvidersStruct struct {
	Docker ProvidersDockerStruct `yaml:"docker"`
	File   ProvidersFileStruct   `yaml:"file"`
}

type DnsChallengeStruct struct {
	Provider         string   `yaml:"provider,omitempty"`
	Resolvers        []string `yaml:"resolvers,omitempty"`
	DelayBeforeCheck string   `yaml:"delayBeforeCheck,omitempty"`
}
type CertAcmeStruct struct {
	Email        string             `yaml:"email,omitempty"`
	Storage      string             `yaml:"storage,omitempty"`
	DnsChallenge DnsChallengeStruct `yaml:"dnsChallenge,omitempty"`
}
type CertDnsCloudflareStruct struct {
	Acme CertAcmeStruct `yaml:"acme,omitempty"`
}
type CertificatesResolversStruct struct {
	DnsCloudflare CertDnsCloudflareStruct `yaml:"dns-cloudflare,omitempty"`
}

type TraefikConfig struct {
	Global                GlobalStruct                `yaml:"global"`
	EntryPoints           EntryPoints                 `yaml:"entryPoints"`
	Experimental          ExperimentalStruct          `yaml:"experimental,omitempty"`
	Hub                   HubStruct                   `yaml:"hub"`
	Metrics               MetricsStruct               `yaml:"metrics"`
	Api                   ApiStruct                   `yaml:"api"`
	Ping                  PingStruct                  `yaml:"ping"`
	Log                   LogStruct                   `yaml:"log"`
	AccessLog             AccessLogStruct             `yaml:"accessLog"`
	Providers             ProvidersStruct             `yaml:"providers"`
	CertificatesResolvers CertificatesResolversStruct `yaml:"certificatesResolvers,omitempty"`
}

/*
Create traefik config file
*/
func createTraefikFile() {
	entryPointsHttps := EntryPointHttpsStruct{
		Address: ":443",
		Http: HttpsHttpStruct{
			Tls: TlsStruct{
				Domains: []HttpTlsDomainsStruct{
					{
						Main: "traefik." + configRootDomain,
						Sans: []string{
							"*." + configRootDomain,
						},
					},
				},
			},
		},
	}
	certificatesResolvers := CertificatesResolversStruct{}

	if configUseCloudflare == true {
		entryPointsHttpsCloudflare := EntryPointHttpsStruct{
			Http: HttpsHttpStruct{
				Tls: TlsStruct{
					CertResolver: "dns-cloudflare",
				},
			},
			ForwardedHeaders: ForwardedHeadersStruct{
				TrustedIPs: []string{
					"173.245.48.0/20",
					"103.21.244.0/22",
					"103.22.200.0/22",
					"103.31.4.0/22",
					"141.101.64.0/18",
					"108.162.192.0/18",
					"190.93.240.0/20",
					"188.114.96.0/20",
					"197.234.240.0/22",
					"198.41.128.0/17",
					"162.158.0.0/15",
					"104.16.0.0/13",
					"104.24.0.0/14",
					"172.64.0.0/13",
					"131.0.72.0/22",
				},
			},
		}
		certificatesResolversCloudflare := CertificatesResolversStruct{
			DnsCloudflare: CertDnsCloudflareStruct{
				Acme: CertAcmeStruct{
					Email:   configCloudflareEmail,
					Storage: "/etc/traefik/acme.json",
					DnsChallenge: DnsChallengeStruct{
						Provider:         "cloudflare",
						Resolvers:        []string{"1.1.1.1:53", "1.0.0.1:54"},
						DelayBeforeCheck: "90",
					},
				},
			},
		}

		mergo.Merge(&entryPointsHttps, entryPointsHttpsCloudflare)
		mergo.Merge(&certificatesResolvers, certificatesResolversCloudflare)
	}

	traefikConfig := TraefikConfig{
		Global: GlobalStruct{
			CheckNewVersion:    true,
			SendAnonymousUsage: false,
		},
		EntryPoints: EntryPoints{
			Http: EntryPointHttpStruct{
				Address: ":80",
			},
			Https: entryPointsHttps,
			Ping: EntryPointPing{
				Address: ":8081",
			},
		},
		Experimental: ExperimentalStruct{
			Hub: configUseTraefikHub,
		},
		Hub: HubStruct{
			Tls: HubTlsStruct{
				Insecure: false,
			},
		},
		Metrics: MetricsStruct{
			Prometheus: PrometheusStruct{
				AddRoutersLabels: true,
			},
		},
		Api: ApiStruct{
			Insecure:  false,
			Dashboard: true,
		},
		Ping: PingStruct{
			EntryPoint: "ping",
		},
		Log: LogStruct{
			Level:    "WARN",
			FilePath: "/etc/traefik/traefik.log",
		},
		AccessLog: AccessLogStruct{
			FilePath:      "/etc/traefik/access.log",
			BufferingSize: 100,
			Filters: AccessLogFiltersStruct{
				StatusCodes: []string{"400-499"},
			},
		},
		Providers: ProvidersStruct{
			Docker: ProvidersDockerStruct{
				Endpoint:         "tcp://socket-proxy:2375",
				ExposedByDefault: false,
				Network:          "traefik-proxy",
				SwarmMode:        false,
			},
			File: ProvidersFileStruct{
				Directory: "/etc/traefik/rules",
				Watch:     true,
			},
		},
		CertificatesResolvers: certificatesResolvers,
	}

	yamlData, err := yaml.Marshal(&traefikConfig)

	if err != nil {
		fmt.Printf("Error while Marshaling. %v", err)
	}

	writeFile("AppData/traefik-proxy", "acme.json", []byte(""), 0600)
	writeFile("AppData/traefik-proxy", "traefik.yml", yamlData, 0644)

	if configUseCloudflare == true {
		writeFile("secrets", "cf_email", []byte(configCloudflareEmail), 0600)
		writeFile("secrets", "cf_api_key", []byte(configCloudflareApiKey), 0600)
	}
	if configIsLocalEnvironment == true {
		traefikTlsDynamicConfig := createTraefikTlsFile()

		writeFile("AppData/traefik-proxy/rules", "local-tls.yml", traefikTlsDynamicConfig, 0644)
	}
}
