package schemaprops

import (
	"bytes"
	"encoding/json"
	"log"
	"sort"
	"strconv"
	"strings"
)

type Collector struct {
	level      int
	currentApi int
	path       []string
	Apis       []PropertyOutput
}

type PropertyOutput struct {
	Api       string     `json:"api,omitempty"`
	Version   string     `json:"version,omitempty"`
	Endpoints []Endpoint `json:"endpoints,omitempty"`
	endpoints map[string]map[string]int
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

type PathRegex struct {
	Regex   string
	Method  string
	Name    string
	Mapping string
}

var subPathx = "[a-zA-Z0-9_{}-]+" // url sub path regex

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
	return c
}

func (c *Collector) SetCollectorAPIDetails(api, version string) {
	p := PropertyOutput{Api: api, Version: version}
	p.endpoints = make(map[string]map[string]int, 0)
	c.Apis = append(c.Apis, p)
	c.currentApi = len(c.Apis) - 1
}

func (c *Collector) GetProperties() map[string]map[string]int {
	return c.Apis[c.currentApi].endpoints
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
func sortPathStrings(m map[string]string) []string {
	keyslice := make([]string, 0)
	for k, _ := range m {
		if len(k) != 0 {
			keyslice = append(keyslice, k)
		}
	}
	sort.Strings(keyslice)
	return keyslice
}

func (c *Collector) CollectProperties(method, endpoint, body string, code int) {
	if len(c.Apis) == 0 {
		return
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
	for _, v := range keyslice {
		pathmap[v] = 0
	}

	shortname := c.stripName(endpoint)
	c.Apis[c.currentApi].endpoints[method+" "+shortname+" "+strconv.Itoa(code)] = pathmap

	return
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

func (c *Collector) OutputJSON() string {
	apis := c.Apis
	for i, api := range apis {
		endpoints := sortEndpoints(api.endpoints)
		for _, k := range endpoints {
			v := c.Apis[i].endpoints[k]
			method, path, code := c.parseEndpoint(k)
			path = pathToSwagger(path)
			endp := Endpoint{Method: method, Path: path}
			// find if endpoint and code already exists - if so use that endpoint
			response, endpoint := c.findEndpointResponse(api.Endpoints, endp, code)
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
			}
			endpoint.Responses = append(endpoint.Responses, response)
			c.Apis[i].Endpoints = append(c.Apis[i].Endpoints, endpoint)
		}
	}

	// Convert structs to JSON.
	jsondata, err := json.MarshalIndent(c.Apis, "", " ")
	if err != nil {
		log.Fatal(err)
	}

	return "{ \n\"responseFields\": " + string(jsondata) + "\n}"
}

func (c *Collector) parseEndpoint(ep string) (string, string, string) {
	split := strings.Split(ep, " ")
	if len(split) != 3 {
		return "", "", ""
	}
	return split[0], split[1], split[2]
}

func (c *Collector) findEndpointResponse(endpoints []Endpoint, endpoint Endpoint, code string) (Response, Endpoint) {
	for ek, ep := range endpoints {
		if endpoint.Method == ep.Method && endpoint.Path == ep.Path {
			for rk, resp := range ep.Responses {
				if resp.Code == code {
					return endpoints[ek].Responses[rk], endpoints[ek]
				}
			}
		}
	}
	return Response{}, Endpoint{}
}
