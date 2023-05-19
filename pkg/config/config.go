package config

import (
	"flag"
	"net/url"
	"os"
	"strings"
)

type Config struct {
	serverUrl    string
	clientId     string
	clientSecret string
	username     string
	password     string
	scope        string
	env          string
}

func Get() *Config {
	conf := &Config{}
	flag.StringVar(&conf.clientId, "clientId", os.Getenv("CLIENT_ID"), "Application Client ID")
	flag.StringVar(&conf.clientSecret, "consumerSecret", os.Getenv("CLIENT_SECRET"), "Application Client Secret")
	flag.StringVar(&conf.username, "username", os.Getenv("USERNAME"), "Accela Account Username")
	flag.StringVar(&conf.password, "password", os.Getenv("PASSWORD"), "Accela Account Password")
	flag.StringVar(&conf.scope, "scopes", os.Getenv("SCOPE"), "Scopes")
	flag.StringVar(&conf.env, "env", os.Getenv("ENVIRONMENT"), "Environment")

	flag.StringVar(&conf.serverUrl, "serverUrl", os.Getenv("SERVER_URL"), "Server Url")

	flag.Parse()

	return conf
}

func (c *Config) GetClientIDAndSecret() (string, string) {
	return c.clientId, c.clientSecret
}

func (c *Config) GetUsernameAndPassword() (string, string) {
	return c.username, c.password
}

func (c *Config) GetScopes() string {
	return c.scope
}

func (c *Config) GetEnv() string {
	return c.env
}

func (c *Config) GetServerURL() string {
	c.serverUrl = strings.TrimSuffix(c.serverUrl, "/")
	u, _ := url.Parse(c.serverUrl)
	if u.Scheme == "" {
		return "https://" + c.serverUrl
	} else {
		return c.serverUrl
	}
}

func (c *Config) GetServerHost() string {
	c.serverUrl = strings.TrimSuffix(c.serverUrl, "/")
	u, _ := url.Parse(c.serverUrl)
	if u.Scheme == "" {
		return u.Host
	} else {
		return u.Host
	}
}
