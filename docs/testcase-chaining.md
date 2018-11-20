# Test case parameter chaining

Test case chaining is ability to use part of the results of one test case as input to another test case. Within the test case model, a rule has one or more test sequences, and each test sequence has one or more test cases. Parameter chaining occurs between two test cases within the same test sequence. The source test case is run before the receiving test case. Any part of a response that can be selected for comparison with a 'Matches' clause, can be used to extract data from the source test case, and make it available to the receiving test case.
Data is transferred between test cases using the Context. Test cases have access to four levels of context scope: 

1. test case scope - only visible to the test case
2. test sequence scope - visible to any test case within the test sequence
3. rule scope - visible to any test sequence within a rule, and all the test sequence's test cases
4. global/manifiest scope - visible to all testing running - discovery information would typically be accessible via global scope

For test case chaining, the test sequence scope is involved as this exists for the duration of a test sequence and therefore enables parameters to be available after a test case has completed and before a subsequent testcase has started.

In order to move parameter between test cases two directives have been introduced,






```go
// Use cases

/*
As a developer I want to perform a test where I load some json which defines a manifest, rule and testcases
I want the rule to manage the execution of the test that includes two test cases
I want the testcases to communicate paramaters between themselves using a context
I want the results of one test case being used as input to the other testcase
I want to use json pattern matching to extract the first returned AccountId from the first testcase and
use that value as the accountid parameter for the second testcase
*/

func TestChainedTestCases(t *testing.T) {
 manifest, err := loadManifest("testdata/passAccountId.json")
require.NoError(t, err)
assert.Equal(t, manifest.Name, "Basic Swagger 2.0 test run")

defer gock.Off()
gock.New("http://myaspsp").Get("/accounts").Reply(200).BodyString(string(getAccountResponse))

for _, rule := range manifest.Rules { // Iterate over Rules
  rule.Executor = executor{}
    for _, testcases := range rule.Tests {
      ctx := NewContext()
      ctx.Put("{AccountId}", "1231231")
      for _, testcase := range testcases {
        fmt.Println("\n==============Dumping testcase =-->")
        testcase.Dump()
        ctx.Dump()

        myReq, _ := testcase.Prepare(ctx, nil)      // Apply inputs, context - results on http object and context
        resp, err := rule.Execute(myReq, &testcase) // execute the testcase
        testcase.Validate(resp, ctx)

        _, _, _ = resp, err, myReq
        }
        _ = ctx
        }
    // rule.ProcessTestCases() // what does this even mean?
    // take the first testcase array,
    // is it a single testcase or multiple?
    // if single ... process/run
    // if multiple ... check that all requirements to run are met
    //- like available context variables
    //- like if an input variable is specificed - its used
    //- like all input variables are available
    //- like if an expects variable is specificed its used
  }
  fmt.Println("\n==============END =-->")
}

type executor struct {
}

func (e *executor) ExecuteTestCase(r *http.Request, t *TestCase, ctx *Context) (*http.Response, error) {
  return nil, nil
}

/*
As a developer I'd like many types of matches in my toolkit
- some matches simple match the response for success for failure
- some matches match a value and put it in the context - for access by other testcases



As a developer I'd like my testcase which uses accountid between testcases to be run against ozone bank

    			  "expect": {
                        "status-code": 200,
                        "matches": [{
                            "description": "A json match on response body",
                            "json": "Data.Account.Accountid",
                            "value": "@AccountId"   // store result in context - what if its not there? Match fails!!! with appropriate message
                        }]
                    }
                    "input": {
                        "method": "GET",
                        "endpoint": "/accounts/{AccountId}",
                        "contextGet": [{ // get a variable from the context and put it somewhere in the request
                            "name": "AccountId", // index into context
                            "purpose": "supplies the account number to be queried",
							"replaceInEndpoint": "{AccountId}" // replace strategy, also replace with jsonexpression - SJSON friend
							"replaceInBody": "{AccountId}" // replace strategy, also replace with jsonexpression - SJSON friend - after body constructed!!! (not get)
							"replaceInHeader/put header?"
							"replaceWithJson":"Data.OBresponse1.Field[3]" // build a response object and but this value in there
							"rawJsonBody":"{}" // raw json to but in the request body - simply use this as the body rather than building
                        }]
                    },
*/

// Context Tests
//
// Put things into a context
// Put and get strings from a context
// Put and get numbers from a context
// Put ang get structures from a context (deep vs shallow)
// Read context from a configuration file
//

/*

Test cases to handle looping through paginated output ... how?

Expects {
	pagecount indicator ... 1 of 6
	// current page
	// total pages
	// logic while current page < total page and not error
	// must be able to report currently page and total pages in any errors
	// pagereader:Testcase
	// firstpage:Testcase
	// loopconstruct:Rule? ... so testcases dont know and run unmodified
	// output feeds to input of same testcase until condition met
	// status - match, fail, loop

	Or ... changing parameter in context referred to by input section
	update context parameter in expects section
	repeat keyword in section
	repeatUntil {
		maxtimes: 10
		Match[] - condition
	}

}


///////////// TESTCASE GENERATION
// iterate over swagger
Get Endpoint/permission combinations
.. so figure out all positive permission permutations
end up with a list of structs that contains a permission + and method/endpoint tied to a testcase

for each permissioned endpoint, identify which permission set it can be satisfied from.
Figure out a set of rules for declarative testcase permission annotations + specifiy default if makes it clearer


OpenAPI - what does it do, what does is provide, is it upto date and reliable, can it be used to enrich our
our core model?

*/
```
