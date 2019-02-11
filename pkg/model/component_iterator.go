package model

import (
	"fmt"
	"reflect"
)

/* The componentIterator executes a component whilst passing a variable number of paramters
to that component. Parameters are string arrays, but may also be two dimensional string arrays.

[]string or
[][]string

The iterator allows multiple parameter arrays which are substituted into the context that the
component uses at execution time.
One use case is multiple token acquisition

*/

// ComponentIterator -
type ComponentIterator struct {
	component  *Component             // component to iterate over
	parameters map[string]interface{} // parameters to vary across each component iteration
	len        int
}

// NewComponentIterator -
func NewComponentIterator(comp *Component, params map[string]interface{}) *ComponentIterator {
	ci := ComponentIterator{component: comp, parameters: params}
	return &ci
}

// CheckParameters -
func (c *ComponentIterator) CheckParameters() {
	var maxlen int
	for k, v := range c.parameters {
		switch reflect.TypeOf(v).Kind() {
		case reflect.Slice:
			var arrayOfArrays [][]string
			var array []string
			var ok bool
			var ln int

			arrayOfArrays, ok = v.([][]string)
			if ok {
				ln = len(arrayOfArrays)
			}
			if !ok {
				array, ok = v.([]string)
				if ok {
					ln = len(array)
				}
			}

			if ln > maxlen {
				maxlen = ln
			}
			fmt.Println(k, v)
		}
	}
	c.len = maxlen
}

// Iterate - using the configured component and supplied parameters, iterate over the component
// using the next parameters from the parameter list each time
func (c *ComponentIterator) Iterate(ctx Context) {
	// copy context
	// pass context to component and execute
	// wait on channel for component completion - waiting for posted context in channel to indicate completion
	// get context results from all contexts by looking at component result parameters
	// execute test cases asynchronously
	// error reporting
	// - component name/id
	// - iterator variations applied
	// - component parameters
	// - error message
	// - collected at listening point
	// - timeout value

	// potential modification to testcase - input context/output context(existing/empty?) -- allows safe context updating
}

func putStringParamInContext(key string, arry [][]string, ctx *Context) {
	ctx.Put(key, arry)

}

func putArrayParamInContext(arry [][]string) {

}
