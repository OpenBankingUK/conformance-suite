//go:generate mockery -name Generator -inpkg
package generation

import (
	"bitbucket.org/openbankingteam/conformance-suite/pkg/schema"
	"github.com/sirupsen/logrus"

	"bitbucket.org/openbankingteam/conformance-suite/pkg/discovery"
	"bitbucket.org/openbankingteam/conformance-suite/pkg/manifest"
	"bitbucket.org/openbankingteam/conformance-suite/pkg/model"
	"bitbucket.org/openbankingteam/conformance-suite/pkg/names"
	"bitbucket.org/openbankingteam/conformance-suite/pkg/permissions"
)

// SpecificationTestCases - test cases generated for a specification
type SpecificationTestCases struct {
	Specification discovery.ModelAPISpecification `json:"apiSpecification"`
	TestCases     []model.TestCase                `json:"testCases"`
}

type GeneratorConfig struct {
	ClientID              string
	Aud                   string
	ResponseType          string
	Scope                 string
	AuthorizationEndpoint string
	RedirectURL           string
	ResourceIDs           model.ResourceIDs
}

// Generator - generates test cases from discovery model
type Generator interface {
	GenerateManifestTests(log *logrus.Entry, config GeneratorConfig, discovery discovery.ModelDiscovery,
		ctx *model.Context, conditional []discovery.ConditionalAPIProperties) (SpecRun, manifest.Scripts, map[string][]manifest.RequiredTokens)
}

// NewGenerator - returns implementation of Generator interface
func NewGenerator() Generator {
	return generator{
		resolver: permissions.Resolver,
	}
}

// generator - implements Generator interface
type generator struct {
	resolver func(groups []permissions.Group) permissions.CodeSetResultSet
}

// shouldIgnoreDiscoveryItem - determine if we should process a `SchemaVersion`. Currently only the following are supported:
// * `Account and Transaction API Specification`
// * `Confirmation of Funds API Specification`
//
// All else returns `true`.
func shouldIgnoreDiscoveryItem(apiSpecification discovery.ModelAPISpecification) bool {
	shouldIgnore := true

	supportedSchemaVersions := []string{
		// `Account and Transaction API Specification
		"https://raw.githubusercontent.com/OpenBankingUK/read-write-api-specs/v3.1.0/dist/account-info-swagger.json",
		// `Confirmation of Funds API Specification`
		"https://raw.githubusercontent.com/OpenBankingUK/read-write-api-specs/v3.1.0/dist/confirmation-funds-swagger.json",
	}
	for _, supportedSchemaVersion := range supportedSchemaVersions {
		if apiSpecification.SchemaVersion == supportedSchemaVersion {
			return false
		}
	}

	return shouldIgnore
}

// Work in progress to integrate Manifest Test
func (g generator) GenerateManifestTests(log *logrus.Entry, config GeneratorConfig, discovery discovery.ModelDiscovery,
	ctx *model.Context, conditionalProperties []discovery.ConditionalAPIProperties) (SpecRun, manifest.Scripts, map[string][]manifest.RequiredTokens) {
	log = log.WithField("module", "GenerateManifestTests")
	for k, item := range discovery.DiscoveryItems {
		spectype, err := manifest.GetSpecType(item.APISpecification.SchemaVersion)
		if err != nil {
			logrus.Warnf("Cannot get spec type from schema version: " + item.APISpecification.SchemaVersion)
			log.Warnf("specification %s not found", item.APISpecification.SchemaVersion)
			continue
		}

		log.Debugf("Generating testcases for %s API", spectype)
		discovery.DiscoveryItems[k].APISpecification.SpecType = spectype
	}

	specTestCases := []SpecificationTestCases{}
	scrSlice := []model.SpecConsentRequirements{}
	var filteredScripts manifest.Scripts
	tokens := map[string][]manifest.RequiredTokens{}

	for _, item := range discovery.DiscoveryItems {
		validator, err := schema.NewSwaggerOBSpecValidator(item.APISpecification.Name, item.APISpecification.Version)
		if err != nil {
			log.WithError(err).Warnf("manifest testcase generation failed for %s", item.APISpecification.SchemaVersion)
			validator = schema.NewNullValidator()
		}
		log.WithFields(logrus.Fields{"name": item.APISpecification.Name, "version": item.APISpecification.Version}).
			Info("swagger spec validator created")

		//scripts, _, err := manifest.LoadGenerationResources(specType, item.APISpecification.Manifest)

		params := manifest.GenerationParameters{
			//	Scripts:      scripts,
			Spec:         item.APISpecification,
			Baseurl:      item.ResourceBaseURI,
			Ctx:          ctx,
			Endpoints:    item.Endpoints,
			ManifestPath: item.APISpecification.Manifest,
			Validator:    validator,
			Conditional:  conditionalProperties,
		}
		tcs, fsc, err := manifest.GenerateTestCases(&params)

		filteredScripts = fsc
		if err != nil {
			log.Warnf("manifest testcase generation failed for %s", item.APISpecification.SchemaVersion)
			continue
		}

		spectype := item.APISpecification.SpecType
		requiredSpecTokens, err := manifest.GetRequiredTokensFromTests(tcs, spectype)
		if err != nil {
			log.Warnf("failed to retrieve required spec tokens from test for spec %s", spectype)
			continue
		}
		logrus.Debugf("%s required spec tokens: %+v", spectype, requiredSpecTokens)
		specreq, err := getSpecConsentsFromRequiredTokens(requiredSpecTokens, item.APISpecification.Name)
		if err != nil {
			log.Warnf("failed to retrieve spec consents from required spec tokens for spec %s", spectype)
			continue
		}
		scrSlice = append(scrSlice, specreq)
		if spectype == "payments" || spectype == "cbpii" { //
			// three sets of test case. all, UI, consent (Non-ui)
			tcs = getUITests(tcs)
		}
		stc := SpecificationTestCases{Specification: item.APISpecification, TestCases: tcs}
		logrus.Debugf("%d test cases generated for %s", len(tcs), item.APISpecification.Name)
		specTestCases = append(specTestCases, stc)
		tokens[spectype] = requiredSpecTokens
	}

	for _, item := range scrSlice {
		logrus.Tracef("%#v", item)
	}
	for _, v := range tokens {
		logrus.Tracef("%#v", v)
	}
	return SpecRun{specTestCases, scrSlice}, filteredScripts, tokens
}

// taks all the payment testscases
// returns two sets
// set 1) - payment tests that show in the UI and execution when runtests is called
// set 2) - payment consent tests that need to be authorised before runtests can happen
func getUITests(tcs []model.TestCase) []model.TestCase {

	uiTests := []model.TestCase{}
	consentJobs := manifest.GetConsentJobs()

	for _, test := range tcs {
		_, exists := consentJobs.Get(test.ID)
		if exists {
			logrus.Tracef("skipping job %s", test.ID)
			continue
		}
		uiTests = append(uiTests, test)
	}

	return uiTests
}

// Packages up Required tokens into a SpecConsentRequirements structure
func getSpecConsentsFromRequiredTokens(rt []manifest.RequiredTokens, apiName string) (model.SpecConsentRequirements, error) {
	npa := []model.NamedPermission{}
	for _, v := range rt {
		np := model.NamedPermission{}
		np.Name = v.Name
		np.CodeSet = permissions.CodeSetResult{}
		np.CodeSet.TestIds = append(np.CodeSet.TestIds, permissions.StringSliceToTestID(v.IDs)...)
		np.CodeSet.CodeSet = append(np.CodeSet.CodeSet, permissions.StringSliceToCodeSet(v.Perms)...)
		npa = append(npa, np)
	}
	specConsentReq := model.SpecConsentRequirements{Identifier: apiName, NamedPermissions: npa}
	return specConsentReq, nil
}

// consentRequirements calls resolver to get list of permission sets required to run all test cases
func (g generator) consentRequirements(specTestCases []SpecificationTestCases) []model.SpecConsentRequirements {
	nameGenerator := names.NewSequentialPrefixedName("to")
	specConsentRequirements := []model.SpecConsentRequirements{}
	for _, spec := range specTestCases {
		var groups []permissions.Group
		for _, tc := range spec.TestCases {
			g := model.NewDefaultPermissionGroup(tc)
			groups = append(groups, g)
		}
		resultSet := g.resolver(groups)
		consentRequirements := model.NewSpecConsentRequirements(nameGenerator, resultSet, spec.Specification.Name)
		specConsentRequirements = append(specConsentRequirements, consentRequirements)
	}
	return specConsentRequirements
}

// SpecRun represents all specs and their test and a list of tokens
// required to run those tests
type SpecRun struct {
	SpecTestCases           []SpecificationTestCases        `json:"specCases"`
	SpecConsentRequirements []model.SpecConsentRequirements `json:"specTokens"`
}
