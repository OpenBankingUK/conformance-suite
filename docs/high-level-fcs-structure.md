# AppJourney Steps

1. **Initialization**
   - Initialize components: generator, validator, daemonController, etc.

2. **Set Configuration**
   - `SetConfig()`: Set journey configuration (certificates, endpoints, client details, etc.)
   - Put parameters into journey shared context

3. **Set Discovery Model**
   - `SetDiscoveryModel()`: Validate and set the discovery model
   - Process conditional API properties if applicable

4. **Validate Discovery Model**
   - Use `Validator` to check the discovery model
   - Record any validation failures

5. **Generate Test Cases**
   - `TestCases()`: Generate manifest tests based on the discovery model
   - For each discovery item:
     - Create a swagger spec validator
     - Generate test cases using `manifest.GenerateTestCases()`
     - Get required tokens for the spec
   - Determine API versions
   - Generate `SpecRun` with test cases and permissions

6. **TLS Validation**
   - For each discovery item:
     - Use `TLSValidator` to check TLS version of the resource base URI
     - Store TLS validation results in the context

7. **Token Acquisition**
   - Handle PSU Consent or Headless Token acquisition based on discovery model
   - For PSU Consent:
     - Get consent IDs and token map
     - Map tokens to test cases
     - Create token collector
   - For Headless Token:
     - Get headless consent
     - Map tokens to test cases
     - Set `allCollected` to true

8. **Dynamic Resource ID Handling** (if enabled)
   - Get dynamic resource IDs
   - Map dynamic IDs to relevant test cases

9. **Collect Tokens** (for PSU Consent)
   - `CollectToken()`: Exchange authorization code for access token
   - Store token in context
   - Get dynamic resource IDs if applicable
   - Mark token as collected in the collector

10. **Prepare for Test Run**
    - Ensure all tokens are collected
    - Apply dynamic resource IDs to test cases if applicable
    - Map tokens to test cases for each specification type

11. **Run Tests**
    - `RunTests()`: Create test case runner and execute test cases
      - Check if test cases have been generated
      - Verify all required tokens have been collected
      - If using dynamic resource IDs:
        - Apply dynamic resource IDs to relevant test cases
        - Set default account ID and statement ID in the journey context
      - Map tokens to test cases for each specification type
      - Create a run definition including:
        - Discovery model
        - SpecRun (generated test cases)
        - Signing certificate
        - Transport certificate
    - Create a new test case runner with the run definition
    - Set the journey context phase to "run"
    - Execute test cases:
      - For each specification in SpecTestCases:
        - For each test case in the specification:
          - Prepare the test case (e.g., replace placeholders, set up authentication)
          - Execute the test case
          - Collect and store the test results
      - Handle any errors during test execution
    - Update the daemon controller with test progress and results
    - Emit events for test execution progress and completion

12. **Monitor Results**
    - `Results()`: Provide access to daemon controller for monitoring test progress

13. **Clean Up**
    - `StopTestRun()`: Stop the test run if needed
    - `NewDaemonController()`: Reset daemon controller and events for new runs

Throughout the journey, various events are emitted, and the context is updated with relevant information. The journey also handles conditional properties based on the discovery model.