package manifest

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/sirupsen/logrus"

	"bitbucket.org/openbankingteam/conformance-suite/pkg/model"
)

// TestPlan species a list of scripts, asserts and other entities required to run a set of test
type TestPlan struct {
	Scripts    Scripts
	References References
}

// Scripts -
type Scripts struct {
	Scripts []Script `json:"scripts,omitempty"`
}

// Script represents a highlevel test definition
type Script struct {
	Description      string            `json:"description,omitempty"`
	Detail           string            `json:"detail,omitempty"`
	ID               string            `json:"id,omitempty"`
	RefURI           string            `json:"refURI,omitempty"`
	Parameters       map[string]string `json:"parameters,omitempty"`
	Headers          map[string]string `json:"headers,omitempty"`
	Resource         string            `json:"resource,omitempty"`
	Asserts          []string          `json:"asserts,omitempty"`
	Method           string            `json:"method,omitempty"`
	URI              string            `json:"uri,omitempty"`
	URIImplemenation string            `json:"uri_implemenation,omitempty"`
	SchemaCheck      bool              `json:"schemaCheck,omitempty"`
}

// References - reference collection
type References struct {
	References map[string]Reference `json:"references,omitempty"`
}

// Reference is an item referred to by the test script list an assert of token reqirement
type Reference struct {
	Expect      model.Expect `json:"expect,omitempty"`
	Permissions []string     `json:"permissions,omitempty"`
}

// AccountData stores account number to be used in the test scripts
type AccountData struct {
	Ais           map[string]string `json:"ais,omitempty"`
	AisConsentIds []string          `json:"ais.ConsetnAccoutId,omitempty"`
	Pis           PisData           `json:"pis,omitempty"`
}

// PisData contains information about PIS accounts required for the test scrips
type PisData struct {
	Currency        string            `json:"Currency,omitempty"`
	DebtorAccount   map[string]string `json:"DebtorAccount,omitempty"`
	MADebtorAccount map[string]string `json:"MADebtorAccount,omitempty"`
}

func loadScripts(filename string) (Scripts, error) {
	plan, err := ioutil.ReadFile(filename)
	if err != nil {
		return Scripts{}, err
	}
	var m Scripts
	err = json.Unmarshal(plan, &m)
	if err != nil {
		return Scripts{}, err
	}
	return m, nil
}

func loadReferences(filename string) (References, error) {
	plan, err := ioutil.ReadFile(filename)
	if err != nil {
		return References{}, err
	}
	var m References
	err = json.Unmarshal(plan, &m)
	if err != nil {
		return References{}, err
	}
	return m, nil
}

func loadTestPlan(filename string) (TestPlan, error) {
	plan, err := ioutil.ReadFile(filename)
	if err != nil {
		return TestPlan{}, err
	}
	var m TestPlan
	err = json.Unmarshal(plan, &m)
	if err != nil {
		return TestPlan{}, err
	}
	return m, nil
}

func testCaseBuilder(s Script, refs map[string]Reference, ctx *model.Context, consents []string) (model.TestCase, error) {
	tc := model.TestCase{}
	tc.ID = s.ID
	tc.Name = s.Description
	tc.Input = buildInputSection(s)
	tc.Purpose = s.Detail
	tc.Context = model.Context{}

	tc.Context.PutContext(ctx)

	for _, a := range s.Asserts {
		ref, exists := refs[a]
		if !exists {
			msg := fmt.Sprintf("assertion %s do not exist in reference data", a)
			logrus.Error(msg)
			return tc, errors.New(msg)
		}
		tc.Expect = ref.Expect
		tc.Expect.SchemaValidation = s.SchemaCheck
	}

	tc.ProcessReplacementFields(ctx, false)
	return tc, nil
}

func getAccountConsent(refs References, vx string) []string {
	ref := refs.References[vx]
	return ref.Permissions
}

func buildInputSection(s Script) model.Input {
	i := model.Input{}
	i.Method = strings.ToUpper(s.Method)
	i.Endpoint = s.URI
	for k, v := range s.Headers {
		i.Headers[k] = v
	}
	return i
}
