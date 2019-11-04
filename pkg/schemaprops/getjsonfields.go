package schemaprops

import (
	"bytes"
	"encoding/json"
	"fmt"
	"sort"
	"strconv"

	"github.com/sirupsen/logrus"
)

type Collector struct {
	level     int
	path      []string
	endpoints map[string]map[string]int
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

func (c *Collector) DumpProperties() {
	logrus.Debug("Dump Properties===============")
	endpoints := sortEndpoints(c.endpoints)
	for _, k := range endpoints {
		logrus.Debugf("%s", k)
		fmt.Printf("%s\n", k)
		v := c.endpoints[k]
		sortedv := sortPaths(v)
		for _, x := range sortedv {
			logrus.Debugf("%s", x)
			fmt.Printf("%s\n", x)
		}
	}
	logrus.Debug("End Dump Properties===============")
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
	requestPaths := make(map[string]int, 20)
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
	c.endpoints[method+" "+endpoint+" "+strconv.Itoa(code)] = pathmap // Store path under method/endpoint/response code key

	return pathmap
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
	c.level++
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
		c.level--
		return
	}
	for k := range r {
		c.path[c.level] = k
		m[c.makePath(c.level)] = 0
		c.expand(r[k], m)
	}
	c.level--
}
