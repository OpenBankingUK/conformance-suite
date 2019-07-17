# Manifest Functions

The Context Model consists of a series name/value pairs that are used to convey variables to testcases and between test cases.

The Matching model consists of a number of match types that allow test case responses to be validated, and for elements of the test case response to be captured and put into the Context. The type of match is determined by which fields are present on the match object.

Context usage currently accomodates the following

|Variable Name | Value     |
|--------------|-----------|
|varA          | literal text - no replacements|
|varB          | reference to another context variable - so gets replaced with the the value of another variable|
|varC          | mixture of literals and references|

## Adding Context Functions to the context Model

Adding functions to the context model would suggest the following expansion to the context usage

|Variable Name | Value     |
|--------------|-----------|
|varA          | literal text - no replacements|
|varB          | reference to another context variable - so gets replaced with the the value of another variable|
|varC          | mixture of literals and references|
|varD          | context function|
|varE          | mixture of literals and functions|
|VarF          | mixture of literals, references and functions|

## Adding context functions to Manifests

Three situations arise when adding function into manifiests

- function result is single usage - so need not be stored for future use
- function result needs to be used in more than one place - so the result need to passed through the context
- function result is some form of auto increment variable (i.e. the function has state)

Additionally, the results of a function call may vary depending on time so the placement of functions may at some stage become time sensistive.

## Passing function results via Manifest 'KeepContextOnSuccess'

The KeepContextOnSuccess section of the manfiest is used to create JSON match types and pass testcase response values into the Context using the ContextPut mechanism of the Expect object. So its essentially an area for applying matches to the test case http response.

## Proposed Solution Manifest Functions

An approach to this would be to manipulate the context variables in the `parameters` section of the manifest, where context interactions already occur. Setting a function result value to a context variable in the `parameters` section removes the need to create additional match types, determine match types by field contents and add function checking into the context variable replacement logic, resulting in a simplfied focused solution.

The solution effectively narrows the context function calls to test case creation for the manfest. Limiting any possible side effects and reducing complexity.
