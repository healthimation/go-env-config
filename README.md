# go-env-config
This library is a shim to allow configuration from environment variables to satisfy the interfaces for config using [go-consul-client](http://www.github.com/divideandconquer/go-consul-client).  The constructs are specifically tailored to allow minimal changes to services already using go-consul-client, but allow us to run on Heroku-like platforms, as well as keeping it easy to migrate back to an environment using Consul.

It is expected that something external to the services is managing the available environment variables, and the apps are just consumers.

## Usage

### Application
There is no standalone usage for this library.

### Library
You can import this library into you golang application and then use it to access your environment configuration:


To fetch configuration:
```golang

import "github.com/healthimation/go-env-config/client"


func main() {

  conf = client.NewEnvLoader()

	//fetch data from the cache
	myConfigString := conf.MustGetString("ENV_VAR")
	myConfigBool := conf.MustGetBool("ENV_VAR")
	myConfigInt := conf.MustGetInt("ENV_VAR")
	myConfigDuration := conf.MustGetDuration("ENV_VAR")

	...
}
```


The balancer package returns the same types as the Consul balancer package, but just determines its information from the environment.
```golang

import "github.com/healthimation/go-env-config/balancer"


func main() {

  serviceName := "authentication"

  // use PGHOST:PGPORT env vars to find the DB, and append "_URL" to service name lookups
  // pass in different env var names for db host and port if required
  balancer = client.NewEnvBalancer("PGHOST", "PGPORT", "_URL")

  // setup db
	dbLoc, err := balancer.FindService(fmt.Sprintf("%s-db", serviceName))
	if err != nil {
		...
	}

  // use the db
	db, err := data.NewDb("postgres", fmt.Sprintf("postgres://%s:%s@%s:%d/authentication?sslmode=disable", dbUser, dbPass, dbLoc.URL, dbLoc.Port))
	...

  // find a service
  loc, err := lb.FindService("secure-code") // looks up SECURE_CODE_URL in environment
	if err != nil {
	   ...
  }

}
```
