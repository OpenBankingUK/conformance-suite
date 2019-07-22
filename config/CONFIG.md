# `CONFIG`

The admin configuration is available at: <https://ob19-admin.o3bank.co.uk/>.

## TPP - `conformance-suite`

<https://ob19-admin.o3bank.co.uk/perry/tpp/view_tpp?id=5d3055cd5bb9d65d12b9cc91>:

| TPP Name          | Organization Id    | Roles             | JWKS URL                                                                               |
|-------------------|--------------------|-------------------|----------------------------------------------------------------------------------------|
| conformance-suite | 0015800001041RbAAI | AISP, PISP, CBPII | https://keystore.openbankingtest.org.uk/0015800001041RbAAI/REfZKo7zN2IeE0X2RFGTb4.jwks |

*NB:* The _Certificate DN_ is not `C=GB,O=OpenBanking,OU=0015800001041RbAAI,CN=REfZKo7zN2IeE0X2RFGTb4` it is `CN=REfZKo7zN2IeE0X2RFGTb4,OU=0015800001041RbAAI,O=OpenBanking,C=GB`, see the _Certificate DN_ in the tables below.

## Software Statements - `conformance-suite`

<https://ob19-admin.o3bank.co.uk/perry/software-statement/view_software-statement?id=5d3065e25bb9d65d12b9cc97>:

| Software Statement Name | Software Statement ID  | TPP               | Roles             | Redirect Urls                                    | JWKS Uri                                                                               | Subject DN                                                         |
|-------------------------|------------------------|-------------------|-------------------|--------------------------------------------------|----------------------------------------------------------------------------------------|--------------------------------------------------------------------|
| conformance-suite       | REfZKo7zN2IeE0X2RFGTb4 | conformance-suite | AISP, PISP, CBPII | https://127.0.0.1:8443/conformancesuite/callback | https://keystore.openbankingtest.org.uk/0015800001041RbAAI/REfZKo7zN2IeE0X2RFGTb4.jwks | CN=REfZKo7zN2IeE0X2RFGTb4,OU=0015800001041RbAAI,O=OpenBanking,C=GB |

## Clients

### `conformance-suite_Open Banking - Strict, No CIBA, Not Headless`

<https://ob19-admin.o3bank.co.uk/perry/client/view_client?id=5d3066225bb9d65d12b9cc98>:

| Client Name                                                    | Client ID                            | Client Secret                        | TPP               | Software Statement Name | Software Statement Id  | Bank          | Bearer Token                 | Resource Server | OIDC SERVER CONFIG |   | Authorization End-point:                | Token Server Config                        | OIDC CLIENT CONFIG |   | Scopes:                                     | Redirect URIs                                    | Certificate DN                                                     | Token Endpoint Auth Method      | Response Types | ID token signed response alg | Request Object Signing Alg | Token Endpoint Signing Alg | JWKS URI (From Software Statement)                                                     |
|----------------------------------------------------------------|--------------------------------------|--------------------------------------|-------------------|-------------------------|------------------------|---------------|------------------------------|-----------------|--------------------|---|-----------------------------------------|--------------------------------------------|--------------------|---|---------------------------------------------|--------------------------------------------------|--------------------------------------------------------------------|---------------------------------|----------------|------------------------------|----------------------------|----------------------------|----------------------------------------------------------------------------------------|
| conformance-suite_Open Banking - Strict, No CIBA, Not Headless | bd9ea798-f3a3-494e-94b2-2bd69ac0245c | 07b051be-a568-4d34-b83c-36581d7ab86e | conformance-suite | conformance-suite       | REfZKo7zN2IeE0X2RFGTb4 | OpenBanking-1 | `<Bearer Token>`             |                 |                    |   | https://ob19-auth1-ui.o3bank.co.uk/auth | https://ob19-auth1.o3bank.co.uk:4201/token |                    |   | openid,payments,accounts,fundsconfirmations | https://127.0.0.1:8443/conformancesuite/callback | CN=REfZKo7zN2IeE0X2RFGTb4,OU=0015800001041RbAAI,O=OpenBanking,C=GB | private_key_jwt,tls_client_auth | code id_token  | PS256                        | PS256                      | PS256                      | https://keystore.openbankingtest.org.uk/0015800001041RbAAI/REfZKo7zN2IeE0X2RFGTb4.jwks |

### `conformance-suite_Open Banking - Permissive, CIBA, Headless`

<https://ob19-admin.o3bank.co.uk/perry/client/view_client?id=5d3180095bb9d65d12b9cc99>:

| Client Name                                                 | Client ID                            | Client Secret                        | TPP               | Software Statement Name | Software Statement Id  | Bank          | Bearer Token                 | Resource Server | OIDC SERVER CONFIG |   | Authorization End-point:                | Token Server Config                        | OIDC CLIENT CONFIG |   | Scopes:                                     | Redirect URIs                                    | Certificate DN                                                     | Token Endpoint Auth Method                                            | Response Types     | ID token signed response alg | Request Object Signing Alg | Token Endpoint Signing Alg | JWKS URI (From Software Statement)                                                     |
|-------------------------------------------------------------|--------------------------------------|--------------------------------------|-------------------|-------------------------|------------------------|---------------|------------------------------|-----------------|--------------------|---|-----------------------------------------|--------------------------------------------|--------------------|---|---------------------------------------------|--------------------------------------------------|--------------------------------------------------------------------|-----------------------------------------------------------------------|--------------------|------------------------------|----------------------------|----------------------------|----------------------------------------------------------------------------------------|
| conformance-suite_Open Banking - Permissive, CIBA, Headless | 081756dd-17f5-4543-a221-012e7ec8694e | b1491863-89ec-42f1-befe-97bbb9c51243 | conformance-suite | conformance-suite       | REfZKo7zN2IeE0X2RFGTb4 | OpenBanking-1 | `<Bearer Token>`             |                 |                    |   | https://ob19-auth1-ui.o3bank.co.uk/auth | https://ob19-auth1.o3bank.co.uk:4201/token |                    |   | openid,payments,accounts,fundsconfirmations | https://127.0.0.1:8443/conformancesuite/callback | CN=REfZKo7zN2IeE0X2RFGTb4,OU=0015800001041RbAAI,O=OpenBanking,C=GB | client_secret_basic,client_secret_jwt,private_key_jwt,tls_client_auth | code,code id_token | PS256                        | none,HS256,RS256,PS256     | none,HS256,RS256,PS256     | https://keystore.openbankingtest.org.uk/0015800001041RbAAI/REfZKo7zN2IeE0X2RFGTb4.jwks |
