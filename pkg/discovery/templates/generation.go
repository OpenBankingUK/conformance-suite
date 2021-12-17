package templates

import (
	"encoding/json"
	"io/ioutil"
	"regexp"
	"sort"
	"strings"

	"github.com/go-openapi/spec"
	"github.com/sirupsen/logrus"

	"bitbucket.org/openbankingteam/conformance-suite/pkg/discovery"
	"bitbucket.org/openbankingteam/conformance-suite/pkg/model"

	"github.com/go-openapi/loads"
)

// DraftDiscovery - Can be run to generate a "draft" generic discovery
// template by downloading configured specification swagger files, and
// writing out endpoint paths to the draft template.
//
// Not intended to be run in production.
func DraftDiscovery() error {
	template := newModel()
	for _, config := range specConfigs() {
		doc, err := loadSpec(config.SchemaVersion.String(), false)
		if err != nil {
			return err
		}
		specVersion := checkVersion(doc, config)
		checkName(doc, config)
		updateSpecVersion(&template, specVersion)

		item := newItem(config, specVersion)
		addEndpoints(&item, doc)
		template.DiscoveryItems = append(template.DiscoveryItems, item)
	}
	discModel := discovery.Model{
		DiscoveryModel: template,
	}
	r, err := json.MarshalIndent(discModel, "", "  ")
	if err != nil {
		return err
	}
	return ioutil.WriteFile("draft-"+discModel.DiscoveryModel.Name+".json", r, 0644)
}

func updateSpecVersion(template *discovery.ModelDiscovery, specVersion string) {
	template.Name = strings.Replace(template.Name, "[VERSION]", specVersion, 1)
	template.Description = strings.Replace(template.Description, "[VERSION]", specVersion, 1)
}

func addEndpoints(item *discovery.ModelDiscoveryItem, doc *loads.Document) {
	paths := doc.Spec().Paths.Paths
	keys := make([]string, 0, len(paths))
	for k := range paths {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, path := range keys {
		props := paths[path]
		for meth, op := range getOperations(&props) {
			if op != nil {
				ep := discovery.ModelEndpoint{
					Method: meth,
					Path:   path,
				}
				item.Endpoints = append(item.Endpoints, ep)
			}
		}
	}
}

func newModel() discovery.ModelDiscovery {
	template := discovery.ModelDiscovery{
		Name:             "ob-[VERSION]-generic",
		Description:      "An Open Banking UK generic discovery template for [VERSION] of Accounts and Payments.",
		DiscoveryVersion: "v0.4.0",
		TokenAcquisition: "psu",
		DiscoveryItems:   []discovery.ModelDiscoveryItem{},
	}
	return template
}

func specConfigs() []model.Specification {
	configs := []model.Specification{}
	for _, config := range model.Specifications() {
		if accountsOrPayments(config) {
			configs = append(configs, config)
		}
	}
	return configs
}

func accountsOrPayments(config model.Specification) bool {
	accountsOrPayments := strings.Contains(config.Identifier, "account") ||
		strings.Contains(config.Identifier, "payment")
	return accountsOrPayments
}

func newItem(spec model.Specification, specVersion string) discovery.ModelDiscoveryItem {
	item := discovery.ModelDiscoveryItem{
		APISpecification: discovery.ModelAPISpecification{
			Name:          spec.Name,
			URL:           spec.URL.String(),
			Version:       spec.Version,
			SchemaVersion: spec.SchemaVersion.String(),
		},
		OpenidConfigurationURI: "https://example.com/.well-known/openid-configuration",
		ResourceBaseURI:        "https://example.com:4501/open-banking/" + specVersion + "/",
		Endpoints:              []discovery.ModelEndpoint{},
	}
	return item
}

func checkName(doc *loads.Document, spec model.Specification) {
	title := doc.Spec().Info.Title
	if title != spec.Name {
		logrus.Println("Our configured spec.Name: " +
			spec.Name +
			" must match swagger title: " +
			title)
	}
}

func checkVersion(doc *loads.Document, spec model.Specification) string {
	version := doc.Spec().Info.Version
	patchRe := regexp.MustCompile(`\..$`)
	specVersion := patchRe.ReplaceAllString(version, "")
	if specVersion != spec.Version {
		logrus.Println("Our configured spec.Version: " +
			spec.Version +
			" must match swagger version: " +
			specVersion)
	}
	return specVersion
}

// loads specification via http or file
func loadSpec(spec string, print bool) (*loads.Document, error) {
	doc, err := loads.Spec(spec)
	if err != nil {
		return nil, err
	}
	if print {
		var jsondoc []byte
		jsondoc, err = json.MarshalIndent(doc.Spec(), "", "    ")
		if err != nil {
			return nil, err
		}

		logrus.Println(string(jsondoc))
	}
	return doc, err
}

// getOperations returns a mapping of HTTP Verb name to "spec operation name"
func getOperations(props *spec.PathItem) map[string]*spec.Operation {
	ops := map[string]*spec.Operation{
		"DELETE":  props.Delete,
		"GET":     props.Get,
		"HEAD":    props.Head,
		"OPTIONS": props.Options,
		"PATCH":   props.Patch,
		"POST":    props.Post,
		"PUT":     props.Put,
	}

	// Keep those != nil
	for key, op := range ops {
		if op == nil {
			delete(ops, key)
		}
	}
	return ops
}
