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
}
