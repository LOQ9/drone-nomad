package nomad

// Client is used to create a configuration of a scheduler
type Client struct {
	// Address is the url to connect to the cluster manager
	Address string
	// Region is the place where the service is going to be deployed
	Region string
	// Namespace allow jobs and their associated objects to be segmented from each other and other users of the cluster
	Namespace string
	// Token is the authorization code for nomad
	Token string
	// TLSConfig contains client tls configurations
	TLSConfig *ClientTLSConfig
}

// ClientTLSConfig contains client tls configurations
type ClientTLSConfig struct {
	// CACert is the path to a PEM-encoded CA cert file to use to verify the
	// Nomad server SSL certificate.
	CACert string

	// CAPath is the path to a directory of PEM-encoded CA cert files to verify
	// the Nomad server SSL certificate.
	CAPath string

	// CACertPem is the PEM-encoded CA cert to use to verify the Nomad server
	// SSL certificate.
	CACertPEM []byte

	// ClientCert is the path to the certificate for Nomad communication
	ClientCert string

	// ClientCertPEM is the PEM-encoded certificate for Nomad communication
	ClientCertPEM []byte

	// ClientKey is the path to the private key for Nomad communication
	ClientKey string

	// ClientKeyPEM is the PEM-encoded private key for Nomad communication
	ClientKeyPEM []byte

	// TLSServerName, if set, is used to set the SNI host when connecting via
	// TLS.
	TLSServerName string

	// Insecure enables or disables SSL verification
	Insecure bool
}
