package main

import (
	"fmt"
	"gopkg.in/yaml.v3"
)

type AutheliaServerStruct struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}
type AutheliaLogStruct struct {
	Level      string `yaml:"level"`
	FilePath   string `yaml:"file_path"`
	KeepStdout bool   `yaml:"keep_stdout"`
}
type AutheliaNtpStruct struct {
	Address             string `yaml:"address"`
	Version             int    `yaml:"version"`
	MaxDesync           string `yaml:"max_desync"`
	DisableStartupCheck bool   `yaml:"disable_startup_check"`
	DisableFailure      bool   `yaml:"disable_failure"`
}
type AutheliaTotpStruct struct {
	Disable    bool   `yaml:"disable"`
	Issuer     string `yaml:"issuer"`
	Algorithm  string `yaml:"algorithm"`
	Digits     int    `yaml:"digits"`
	Period     int    `yaml:"period"`
	Skew       int    `yaml:"skew"`
	SecretSize int    `yaml:"secret_size"`
}
type AutheliaWebauthnStruct struct {
	Disable                         bool   `yaml:"disable"`
	DisplayName                     string `yaml:"display_name"`
	AttestationConveyancePreference string `yaml:"attestation_conveyance_preference"`
	UserVerification                string `yaml:"user_verification"`
	Timeout                         string `yaml:"timeout"`
}

type AutheliaAuthBackendPasswordStruct struct {
	Disable bool `yaml:"disable"`
}
type AutheliaAuthBackendFilePasswordStruct struct {
	Algorithm   string `yaml:"algorithm"`
	Iterations  int    `yaml:"iterations"`
	KeyLength   int    `yaml:"key_length"`
	SaltLength  int    `yaml:"salt_length"`
	Parallelism int    `yaml:"parallelism"`
	Memory      int    `yaml:"memory"`
}
type AutheliaAuthBackendFileStruct struct {
	Path     string                                `yaml:"path"`
	Password AutheliaAuthBackendFilePasswordStruct `yaml:"password"`
}
type AutheliaAuthenticationBackendStruct struct {
	PasswordReset AutheliaAuthBackendPasswordStruct `yaml:"password_reset"`
	File          AutheliaAuthBackendFileStruct     `yaml:"file"`
}

type AutheliaAccessControlRulesStruct struct {
	Domain  []string `yaml:"domain"`
	Policy  string   `yaml:"policy"`
	Subject []string `yaml:"subject"`
}
type AutheliaAccessControlStruct struct {
	DefaultPolicy string                             `yaml:"default_policy"`
	Rules         []AutheliaAccessControlRulesStruct `yaml:"rules"`
}

type AutheliaSessionRedisStruct struct {
	Host                     string `yaml:"host"`
	Port                     int    `yaml:"port"`
	DatabaseIndex            int    `yaml:"database_index"`
	MaximumActiveConnections int    `yaml:"maximum_active_connections"`
	MinimumIdleConnections   int    `yaml:"minimum_idle_connections"`
}
type AutheliaSessionStruct struct {
	Redis      AutheliaSessionRedisStruct `yaml:"redis"`
	Name       string                     `yaml:"name"`
	Expiration int                        `yaml:"expiration"`
	Inactivity int                        `yaml:"inactivity"`
	Domain     string                     `yaml:"domain"`
}

type AutheliaRegulationStruct struct {
	MaxRetries int    `yaml:"max_retries"`
	FindTime   string `yaml:"find_time"`
	BanTime    string `yaml:"ban_time"`
}

type AutheliaPasswordPolicyStandardStruct struct {
	Enabled          bool `yaml:"enabled"`
	MinLength        int  `yaml:"min_length"`
	MaxLength        int  `yaml:"max_length"`
	RequireUppercase bool `yaml:"require_uppercase"`
	RequireLowercase bool `yaml:"require_lowercase"`
	RequireNumber    bool `yaml:"require_number"`
	RequireSpecial   bool `yaml:"require_special"`
}
type AutheliaPasswordPolicyStruct struct {
	Standard AutheliaPasswordPolicyStandardStruct `yaml:"standard"`
}

type AutheliaStorageMysqlStruct struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Database string `yaml:"database"`
	Username string `yaml:"username"`
}
type AutheliaStorageStruct struct {
	Mysql AutheliaStorageMysqlStruct `yaml:"mysql"`
}

type AutheliaNotifierSmtpStruct struct {
	Username            string `yaml:"username"`
	Host                string `yaml:"host"`
	Port                int    `yaml:"port"`
	Sender              string `yaml:"sender"`
	Subject             string `yaml:"subject"`
	StartupCheckAddress string `yaml:"startup_check_address"`
}
type AutheliaNotifierStruct struct {
	DisableStartupCheck bool                       `yaml:"disable_startup_check"`
	Smtp                AutheliaNotifierSmtpStruct `yaml:"smtp"`
}

type AutheliaConfigStruct struct {
	Server                AutheliaServerStruct                `yaml:"server"`
	Log                   AutheliaLogStruct                   `yaml:"log"`
	DefaultRedirectionUrl string                              `yaml:"default_redirection_url"`
	Ntp                   AutheliaNtpStruct                   `yaml:"ntp"`
	Totp                  AutheliaTotpStruct                  `yaml:"totp"`
	WebAuthn              AutheliaWebauthnStruct              `yaml:"webauthn"`
	AuthenticationBackend AutheliaAuthenticationBackendStruct `yaml:"authentication_backend"`
	AccessControl         AutheliaAccessControlStruct         `yaml:"access_control"`
	Session               AutheliaSessionStruct               `yaml:"session"`
	Regulation            AutheliaRegulationStruct            `yaml:"regulation"`
	PasswordPolicy        AutheliaPasswordPolicyStruct        `yaml:"password_policy"`
	Storage               AutheliaStorageStruct               `yaml:"storage"`
	Notifier              AutheliaNotifierStruct              `yaml:"notifier"`
}

/*
Create an authelia configuration
*/
func createAutheliaConfig() {
	policyRootDomain := "one_factor"

	if configPolicyForRootDomain == "two_factor" {
		policyRootDomain = "two_factor"
	}
	if configPolicyForRootDomain == "bypass" {
		policyRootDomain = "bypass"
	}
	autheliaConfig := AutheliaConfigStruct{
		Server: AutheliaServerStruct{
			Host: "0.0.0.0",
			Port: 9091,
		},
		Log: AutheliaLogStruct{
			Level:      "warn",
			FilePath:   "/config/authelia.log",
			KeepStdout: true,
		},
		DefaultRedirectionUrl: "https://authelia." + configRootDomain,
		Ntp: AutheliaNtpStruct{
			Address:             "time.cloudflare.com:123",
			Version:             3,
			MaxDesync:           "3s",
			DisableStartupCheck: false,
			DisableFailure:      false,
		},
		Totp: AutheliaTotpStruct{
			Disable:    false,
			Issuer:     "authelia.com",
			Algorithm:  "sha1",
			Digits:     6,
			Period:     30,
			Skew:       1,
			SecretSize: 32,
		},
		WebAuthn: AutheliaWebauthnStruct{
			Disable:                         false,
			DisplayName:                     "Authelia Traefik",
			AttestationConveyancePreference: "indirect",
			UserVerification:                "preferred",
			Timeout:                         "60s",
		},
		AuthenticationBackend: AutheliaAuthenticationBackendStruct{
			PasswordReset: AutheliaAuthBackendPasswordStruct{
				Disable: false,
			},
			File: AutheliaAuthBackendFileStruct{
				Path: "/config/users_database.yml",
				Password: AutheliaAuthBackendFilePasswordStruct{
					Algorithm:   "argon2id",
					Iterations:  3,
					KeyLength:   32,
					SaltLength:  16,
					Parallelism: 4,
					Memory:      1024,
				},
			},
		},
		AccessControl: AutheliaAccessControlStruct{
			DefaultPolicy: "deny",
			Rules: []AutheliaAccessControlRulesStruct{
				{
					Domain: []string{"authelia." + configRootDomain},
					Policy: "bypass",
				},
				{
					Domain: []string{configRootDomain},
					Policy: policyRootDomain,
				},
				{
					Domain:  []string{"traefik." + configRootDomain},
					Policy:  "two_factor",
					Subject: []string{"group:admins", "group:traefik"},
				},
			},
		},
		Session: AutheliaSessionStruct{
			Redis: AutheliaSessionRedisStruct{
				Host:                     "redis",
				Port:                     6379,
				DatabaseIndex:            0,
				MaximumActiveConnections: 8,
				MinimumIdleConnections:   0,
			},
			Name:       "authelia_session",
			Expiration: 3600,
			Inactivity: 300,
			Domain:     configRootDomain,
		},
		Regulation: AutheliaRegulationStruct{
			MaxRetries: 3,
			FindTime:   "2m",
			BanTime:    "5m",
		},
		PasswordPolicy: AutheliaPasswordPolicyStruct{
			Standard: AutheliaPasswordPolicyStandardStruct{
				Enabled:          true,
				MinLength:        8,
				MaxLength:        0,
				RequireUppercase: true,
				RequireLowercase: true,
				RequireNumber:    true,
				RequireSpecial:   false,
			},
		},
		Storage: AutheliaStorageStruct{
			Mysql: AutheliaStorageMysqlStruct{
				Host:     "mariadb",
				Port:     3306,
				Database: "authelia",
				Username: "authelia",
			},
		},
		Notifier: AutheliaNotifierStruct{
			DisableStartupCheck: false,
			Smtp: AutheliaNotifierSmtpStruct{
				Username:            configSmtpUsername,
				Host:                configSmtpHost,
				Port:                configSmtpPort,
				Sender:              "Authelia <" + configSmtpSender + ">",
				Subject:             "[Authelia] {title}",
				StartupCheckAddress: configSmtpStartupAddress,
			},
		},
	}

	yamlData, err := yaml.Marshal(&autheliaConfig)

	if err != nil {
		fmt.Printf("Error while Marshaling. %v", err)
	}

	writeFile("AppData/authelia", "configuration.yml", yamlData, 0644)

	// Secrets
	writeFile("secrets", "authelia_notifier_smtp_password", []byte(configSmtpPassword), 0600)
	writeFile("secrets", "authelia_jwt_secret", []byte(randomString(32)), 0600)
	writeFile("secrets", "authelia_session_secret", []byte(randomString(32)), 0600)
	writeFile("secrets", "authelia_storage_encryption_key", []byte(randomString(32)), 0600)
	writeFile("secrets", "mysql_password", []byte(randomString(16)), 0600)
	writeFile("secrets", "mysql_root_password", []byte(randomString(16)), 0600)
}
