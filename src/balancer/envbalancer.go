package balancer

import (
	"fmt"
	"net"
	"net/url"
	"os"
	"strconv"
	"strings"

	dns "github.com/divideandconquer/go-consul-client/src/balancer"
)

// DATABASE_URL	postgres://user:password@host:5432/dbname?sslmode=disable
// DNS balancer finds services through dns and balances load across them
// type DNS interface {
// 	FindService(serviceName string) (*ServiceLocation, error)
// }
//
// // ServiceLocation is a representation of where a service lives
// type ServiceLocation struct {
// 	URL  string
// 	Port int
// }

// envBalancer satisfies the DNS interface from go-consul-client.
// Specifically, this is so we can parse env vars and return them instead of
// performing consul lookups, and not need to change a ton of service code.
type envBalancer struct {
	dbHostVar        string
	dbPortVar        string
	serviceVarPrefix string
	serviceVarSuffix string
}

func NewEnvBalancer(dbHostVar, dbPortVar, serviceVarPrefix, serviceVarSuffix string) dns.DNS {
	return envBalancer{
		dbHostVar:        dbHostVar,
		dbPortVar:        dbPortVar,
		serviceVarPrefix: serviceVarPrefix,
		serviceVarSuffix: serviceVarSuffix,
	}
}

func (e envBalancer) FindService(serviceName string) (*dns.ServiceLocation, error) {
	// This is kind of lame.
	if strings.HasSuffix(serviceName, "-db") {
		return e.findDB()
	}

	return e.findService(serviceName)
}

func (e envBalancer) GetHttpUrl(serviceName string, useTLS bool) (url.URL, error) {
	result := url.URL{}
	loc, err := e.FindService(serviceName)
	if err != nil {
		return result, err
	}
	result.Host = loc.URL
	if loc.Port != 0 {
		result.Host = net.JoinHostPort(loc.URL, strconv.Itoa(loc.Port))
	}
	if useTLS {
		result.Scheme = "https"
	} else {
		result.Scheme = "http"
	}
	return result, nil
}

// example var: HMD_SECURE_CODE_URL=http://code-1547809825.internal:8080
func (e envBalancer) findService(serviceName string) (*dns.ServiceLocation, error) {
	service, err := getEnvNotBlank(e.buildServiceVar(serviceName))
	if err != nil {
		return nil, err
	}

	u, err := url.Parse(service)
	if err != nil {
		return nil, err
	}

	host, port, err := net.SplitHostPort(u.Host)
	if err != nil {
		return nil, err
	}

	p, err := strconv.Atoi(port)
	if err != nil {
		return nil, err
	}

	return &dns.ServiceLocation{URL: host, Port: p}, nil
}

func (e envBalancer) findDB() (*dns.ServiceLocation, error) {
	// We only care that the values are blank when we try to use them, since
	// services might not have a DB.
	if e.dbHostVar == "" {
		return nil, fmt.Errorf("db host var name cannot be blank")
	}

	if e.dbPortVar == "" {
		return nil, fmt.Errorf("db port var name cannot be blank")
	}

	dbHost, err := getEnvNotBlank(e.dbHostVar)
	if err != nil {
		return nil, err
	}

	dbPort, err := getEnvNotBlank(e.dbPortVar)
	if err != nil {
		return nil, err
	}
	p, err := strconv.Atoi(dbPort)
	if err != nil {
		return nil, err
	}

	return &dns.ServiceLocation{URL: dbHost, Port: p}, nil

}

// takes a service name like "secure-code" and turns it into an name like SECURE_CODE_URL
func (e envBalancer) buildServiceVar(serviceName string) string {
	// split string on any dashes to remove them and get a slice
	slugs := strings.Split(serviceName, "-")

	// prepend the service prefix if it's set
	if e.serviceVarPrefix != "" {
		slugs = append([]string{e.serviceVarPrefix}, slugs...)
	}

	// append the service suffix if it's set
	if e.serviceVarSuffix != "" {
		slugs = append(slugs, e.serviceVarSuffix)
	}

	// return upper case and joined with underscores
	return strings.ToUpper(strings.Join(slugs, "_"))
}

// // Takes a url and returns a service location
func (e envBalancer) urlStringToServiceLocation(serviceUrl string) (*dns.ServiceLocation, error) {
	u, err := url.Parse(serviceUrl)
	if err != nil {
		return nil, err
	}

	host, port, err := net.SplitHostPort(u.Host)
	if err != nil {
		return nil, err
	}

	p, err := strconv.Atoi(port)
	if err != nil {
		return nil, err
	}

	return &dns.ServiceLocation{URL: host, Port: p}, nil
}

func getEnvNotBlank(key string) (string, error) {
	v := os.Getenv(key)
	if v == "" {
		return "", fmt.Errorf("(%s) is blank or unset in the environment", key)
	}

	return v, nil
}
