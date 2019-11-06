package schemaprops

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"sort"
	"strconv"
	"strings"

	"github.com/sirupsen/logrus"
)

type Collector struct {
	level     int
	path      []string
	endpoints map[string]map[string]int
}

type PropertyOutput struct {
	Api       string     `json:"api,omitempty"`
	Version   string     `json:"version,omitempty"`
	Endpoints []Endpoint `json:"endpoints,omitempty"`
}

type Endpoint struct {
	Method    string     `json:"method,omitempty"`
	Path      string     `json:"path,omitempty"`
	Responses []Response `json:"responses,omitempty"`
}

type Response struct {
	Code   string   `json:"code,omitempty"`
	Fields []string `json:"fields,omitempty"`
}

var (
	collector *Collector
)

func GetPropertyCollector() *Collector {
	if collector == nil {
		collector = MakeCollector()
	}
	return collector
}

type PropertyCollector interface {
	CollectProperties(string, string, string, int) map[string]int
	GetProperties() map[string]map[string]int
	DumpProperties()
}

func MakeCollector() *Collector {
	c := &Collector{}
	c.path = make([]string, 20)
	c.endpoints = make(map[string]map[string]int, 0)
	return c
}

func (c *Collector) GetProperties() map[string]map[string]int {
	return c.endpoints
}

func sortEndpoints(m map[string]map[string]int) []string {
	keyslice := make([]string, 0)
	for k, _ := range m {
		if len(k) != 0 {
			keyslice = append(keyslice, k)
		}
	}
	sort.Strings(keyslice)
	return keyslice
}

func sortPaths(m map[string]int) []string {
	keyslice := make([]string, 0)
	for k, _ := range m {
		if len(k) != 0 {
			keyslice = append(keyslice, k)
		}
	}
	sort.Strings(keyslice)
	return keyslice
}

func (c *Collector) CollectProperties(method, endpoint, body string, code int) map[string]int {
	if strings.Contains(endpoint, "PsuDummyURL") {
		return map[string]int{}
	}
	requestPaths := make(map[string]int, 20)
	c.path = make([]string, 20)
	var anyJson map[string]interface{}
	json.Unmarshal([]byte(body), &anyJson)
	for k, _ := range anyJson {
		c.path[c.level] = k
		requestPaths[c.makePath(c.level)] = 0
		mapInterface, ok := anyJson[k].(map[string]interface{})
		if ok {
			c.expand(mapInterface, requestPaths)
		}
	}

	keyslice := make([]string, 0)
	for k, _ := range requestPaths {
		if len(k) != 0 {
			keyslice = append(keyslice, k)
		}
	}
	sort.Strings(keyslice)

	pathmap := make(map[string]int, 0)
	logrus.Debugf("Paths for endpoint: %s %s %d", method, endpoint, code)
	for _, v := range keyslice {
		pathmap[v] = 0
	}

	shortname := c.stripName(endpoint)

	c.endpoints[method+" "+shortname+" "+strconv.Itoa(code)] = pathmap // Store path under method/endpoint/response code key

	return pathmap
}

func (c *Collector) stripName(endpoint string) string {
	result := strings.Split(endpoint, "/open-banking/")
	len := len(result)
	return "/open-banking/" + result[len-1]
}

func (c *Collector) makePath(level int) string {
	var b bytes.Buffer
	for i := 0; i <= level; i++ {
		if len(c.path[i]) > 0 {
			if i > 0 {
				b.WriteString(".")
			}
			b.WriteString(c.path[i])
		}
	}
	return b.String()
}

func (c *Collector) expand(i interface{}, m map[string]int) {
	r, ok := i.(map[string]interface{})
	if !ok {
		switch i.(type) {
		case []interface{}:
			x := i.([]interface{})
			for _, v := range x {
				c.expand(v, m)
			}
		case string:
		default:
		}
		return
	}
	c.level++
	for k := range r {
		c.path[c.level] = k
		m[c.makePath(c.level)] = 0
		c.expand(r[k], m)
	}
	c.level--
}

func (c *Collector) DumpProperties() {
	if logrus.GetLevel() == logrus.TraceLevel {
		logrus.Debug("Dump Properties===============")
		endpoints := sortEndpoints(c.endpoints)
		for _, k := range endpoints {
			logrus.Debugf("%s", k)
			v := c.endpoints[k]
			sortedv := sortPaths(v)
			for _, x := range sortedv {
				logrus.Debugf("%s", x)
			}
		}
		logrus.Debug("End Dump Properties===============")
	}
}

func (c *Collector) OutputJSON(props PropertyOutput) string {
	logrus.Debug("---------------------")
	logrus.Debug("OutputJSON")
	props.Endpoints = make([]Endpoint, 0)

	endpoints := sortEndpoints(c.endpoints)
	for _, k := range endpoints {
		logrus.Debugf("%s", k)

		v := c.endpoints[k]
		method, path, code := c.parseEndpoint(k)
		logrus.Debugf("endpoint: %s, %s, %s", method, path, code)
		endp := Endpoint{Method: method, Path: path}
		// find if endpoint and code already exists - if so use that endpoint

		// Find if endpoint already exists
		// - if so,
		response, endpoint := c.findEndpointResponse(&props, endp, code)
		if response.Code == "" {
			response = Response{Code: code}
			response.Fields = make([]string, 0)
		}
		if endpoint.Method == "" {
			endpoint = Endpoint{}
			endpoint.Method = method
			endpoint.Path = path
			endpoint.Responses = make([]Response, 0)
		}

		sortedv := sortPaths(v)
		for _, x := range sortedv {
			response.Fields = append(response.Fields, x)
			logrus.Debugf("%s", x)
		}

		endpoint.Responses = append(endpoint.Responses, response)

		props.Endpoints = append(props.Endpoints, endpoint)

	}

	// Convert structs to JSON.
	jsondata, err := json.MarshalIndent(props, "", " ")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("\n%s\n", jsondata)

	return ""
}

func (c *Collector) parseEndpoint(ep string) (string, string, string) {
	split := strings.Split(ep, " ")
	if len(split) != 3 {
		return "", "", ""
	}
	return split[0], split[1], split[2]
}

func (c *Collector) findEndpointResponse(props *PropertyOutput, endpoint Endpoint, code string) (Response, Endpoint) {
	for ek, ep := range props.Endpoints {
		if endpoint.Method == ep.Method && endpoint.Path == ep.Path {
			for rk, resp := range ep.Responses {
				if resp.Code == code {
					return props.Endpoints[ek].Responses[rk], props.Endpoints[ek]
				}
			}
		}
	}
	return Response{}, Endpoint{}
}
