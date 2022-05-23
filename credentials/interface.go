package credentials

type ConnectionParams struct {
	Host     string
	Port     string
	User     string
	Password string
	Database string
}

type ICredentials interface {
	GetCredentialsBySource(source Source) ConnectionParams
}

type Databases struct {
	Prod         ConnectionParams
	Intermediate ConnectionParams
	Targets      []ConnectionParams
	Product      string
}
