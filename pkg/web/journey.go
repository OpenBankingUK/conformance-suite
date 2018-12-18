package web

import (
	"bitbucket.org/openbankingteam/conformance-suite/pkg/discovery"
	"github.com/pkg/errors"
)

// Journey represents all possible steps for a user test conformance web journey
type Journey interface {
	SetDiscoveryModel(discoveryModel *discovery.Model) (discovery.ValidationFailures, error)
}

type journey struct {
	validator           discovery.Validator
	validDiscoveryModel *discovery.Model
}

var journeyInstance Journey

// NewWebJourney creates an instance for a user journey, assumes one user only no concurrency
// so a singleton is returned
func NewWebJourney(validator discovery.Validator) Journey {
	if journeyInstance == nil {
		journeyInstance = &journey{
			validator: validator,
		}
	}
	return journeyInstance
}

func (wj *journey) SetDiscoveryModel(discoveryModel *discovery.Model) (discovery.ValidationFailures, error) {
	failures, err := wj.validator.Validate(discoveryModel)
	if err != nil {
		return nil, errors.Wrap(err, "error setting discovery model")
	}

	if !failures.Empty() {
		return failures, nil
	}

	wj.validDiscoveryModel = discoveryModel

	return discovery.NoValidationFailures, nil
}
