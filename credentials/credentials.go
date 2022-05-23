package credentials

import (
	"fmt"
	"os"
)

type EnvPrefix string
type Source string

const (
	Prod         EnvPrefix = "PROD"
	Intermediate EnvPrefix = "INTERMEDIATE"
	Dev          EnvPrefix = "DEV"
	Staging      EnvPrefix = "STAGING"
)

const (
	TremSource  Source = "TREM"
	CoisaSource Source = "COISA"
)

type credential struct {
	host     string
	port     string
	user     string
	password string
	tremDb   string
	coisaDb  string
}

type envNames struct {
	host     string
	port     string
	user     string
	password string
	tremDb   string
	coisaDb  string
}

func NewCredentials(prefix EnvPrefix) ICredentials {
	envs := envNames{
		host:     fmt.Sprintf("%s_HOST", prefix),
		port:     fmt.Sprintf("%s_PORT", prefix),
		user:     fmt.Sprintf("%s_USER", prefix),
		password: fmt.Sprintf("%s_PASS", prefix),
		tremDb:   fmt.Sprintf("%s_DB_TREM", prefix),
		coisaDb:  fmt.Sprintf("%s_DB_COISA", prefix),
	}
	return getCredentials(envs)
}

func getCredentials(envs envNames) credential {
	return credential{
		host:     os.Getenv(envs.host),
		port:     os.Getenv(envs.port),
		user:     os.Getenv(envs.user),
		password: os.Getenv(envs.password),
		tremDb:   os.Getenv(envs.tremDb),
		coisaDb:  os.Getenv(envs.coisaDb),
	}
}

func (c credential) GetCredentialsBySource(source Source) ConnectionParams {
	credentials := c.basicCredentials()
	credentials.Database = c.databaseParam(source)
	return credentials
}

func (c credential) databaseParam(source Source) string {
	switch source {
	case CoisaSource:
		return c.coisaDb
	case TremSource:
		return c.tremDb
	}
	return ""
}

func (c credential) basicCredentials() ConnectionParams {
	return ConnectionParams{
		Host:     c.host,
		Port:     c.port,
		User:     c.user,
		Password: c.password,
	}
}
