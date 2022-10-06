package main

import (
	"fmt"
	"gopkg.in/yaml.v3"
)

func createTraefikTlsFile() []byte {
	type CertsStruct struct {
		CertFile string `yaml:"certFile,omitempty"`
		KeyFile  string `yaml:"KeyFile,omitempty"`
	}
	type DefaultCertificateStruct struct {
		DefaultCertificate CertsStruct `yaml:"defaultCertificate,omitempty"`
	}
	type DefaultStruct struct {
		Default DefaultCertificateStruct `yaml:"default,omitempty"`
	}
	type StoresStruct struct {
		Stores DefaultStruct `yaml:"stores,omitempty"`
	}
	type TlsStruct struct {
		Tls StoresStruct `yaml:"tls,omitempty"`
	}
	tlsConfig := TlsStruct{
		Tls: StoresStruct{
			Stores: DefaultStruct{
				Default: DefaultCertificateStruct{
					DefaultCertificate: CertsStruct{
						CertFile: "/certs/local.crt",
						KeyFile:  "/certs/local.key",
					},
				},
			},
		},
	}

	yamlData, err := yaml.Marshal(&tlsConfig)

	if err != nil {
		fmt.Printf("Error while Marshaling. %v", err)
	}

	return yamlData
}
