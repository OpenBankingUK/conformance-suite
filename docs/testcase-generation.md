# Test Case Generation

We have introduced a new generation package into the application in order to explore options for generating test cases.

The test case generation is initially targeted at open api/swagger endpoint test cases that result from an analysis of the discovery input files.

Test case generation stems from information collected via the discovery model and typically submitted via the application UI.

The process of test case generation currently begins with an examination of the endpoint list submitted in the discovery configuration file. The endpoints are the mandatory, conditional and optional endpoints that have been implemented by the ASPSP.

The first pass at testcase generation, takes this endpoint list and generated a test case object for each of the endpoints. In additional, any resource ID's specified in the discovery configuration file are used to fill in the resource identifieds, for example {AccountId} in "GET /accounts/{AccountId}".
