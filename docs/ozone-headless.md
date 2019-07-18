# Ozone Headless Mechanism

Ozone headless consent flow allows the automatic provision of PSU consent without a PSU interaction. This facility is useful in automated test environments and it allows a non-interactive way of getting an access token.

The Ozone headless consent flow relies on having an Ozone clientId that has been configured to allow a headless interaction. The clientid is configured by the Ozone model bank administrator for this mode of operation.

Once an appropriately configured clientid and associated certificates have been obtained, the following steps take place:

## Client Credentials Grant

Run as normal

## Post Account/Payment Consents

Run as normal

## Generate PSU consent URL

The PSU consent URL is generated in the normal manner. The PSU is required to go to the generated url in their web browser then logon to the ASPSP portal in order to provide their consent for the requested permissions.

For Ozone headless flow. The consent URL is called directly from the functional conformance suite. When the suite calls Ozone, Ozone replies with a `302` `HTTP_FOUND` (MOVED_TEMPORARILY) redirect. The redirect response contains a  `Location` header which is in the form redirect-url/?code=Exchange_code.

The exchange code is retrieved from the location header url parameter `code`. And the normal exchange code sequence is performed in order to get an access token.

### Example Calling Ozone PSU ConsentURL

```httptrace
---------------------- REQUEST LOG -----------------------
GET  /auth?client_id=72b79ddd-4674-43bb-96c2-992f79cd6e62&redirect_uri=https%3A%2F%2F127.0.0.1%3A8443%2Fconformancesuite%2Fcallback&request=eyJhbGciOiJub25lIn0.eyJhdWQiOiJodHRwczovL21vZGVsb2JhbmthdXRoMjAxOC5vM2JhbmsuY28udWs6NDEwMSIsImNsYWltcyI6eyJpZF90b2tlbiI6eyJvcGVuYmFua2luZ19pbnRlbnRfaWQiOnsiZXNzZW50aWFsIjp0cnVlLCJ2YWx1ZSI6ImFhYy1mNGJmY2VlNi1hNmZkLTRhNTktOWUxYS1mOTdhNTczMTNiYzgifX19LCJpc3MiOiI3MmI3OWRkZC00Njc0LTQzYmItOTZjMi05OTJmNzljZDZlNjIiLCJyZWRpcmVjdF91cmkiOiJodHRwczovLzEyNy4wLjAuMTo4NDQzL2NvbmZvcm1hbmNlc3VpdGUvY2FsbGJhY2siLCJzY29wZSI6Im9wZW5pZCBhY2NvdW50cyJ9.&response_type=code&scope=openid+accounts&state=  HTTP/1.1
HOST   : modelobankauth2018.o3bank.co.uk:4101
HEADERS:
               User-Agent: go-resty/1.10.3 (https://github.com/go-resty/resty)
BODY   :
***** NO CONTENT *****
----------------------------------------------------------
```

After calling the PSU Consent url - the `Location` header of the `302` redirection response is captured and provided in the following call, as the code parameter.

### Example Location header returned in HTTP 302 response

```logtrace
Location Header: "https://127.0.0.1:8443/conformancesuite/callback#code=d037ebc3-121a-4df7-88f6-5a67322bf8eb&id_token=eyJhbGciOiJQUzI1NiIsImtpZCI6IjVVWXFjdGNOblZkSTl0VXRlYXA0dFNtV241NCJ9.eyJzdWIiOiJzZHAtMS0wYWZhOWE4Yi1kODI4LTQyOWEtOGUxNS01ZTlkMmRjNjUxNzAiLCJvcGVuYmFua2luZ19pbnRlbnRfaWQiOiJzZHAtMS0wYWZhOWE4Yi1kODI4LTQyOWEtOGUxNS01ZTlkMmRjNjUxNzAiLCJpc3MiOiJodHRwczovL21vZGVsb2JhbmthdXRoMjAxOC5vM2JhbmsuY28udWs6NDEwMSIsImF1ZCI6IjcyYjc5ZGRkLTQ2NzQtNDNiYi05NmMyLTk5MmY3OWNkNmU2MiIsImlhdCI6MTU2MzM3NDQwNiwiZXhwIjoxNTYzMzc4MDA2LCJjX2hhc2giOiJFREIyb3M3alVDSFVRai1OZzViaGF3Iiwic19oYXNoIjoiZVVuZzcxWU9fc0lFVnRDdXJyWTNYZyIsImFjciI6InVybjpvcGVuYmFua2luZzpwc2QyOnNjYSJ9.vJ31O1YvdJ5D8CIKmSWoAwhFO5f0_TD7LsngRjdsfcZpWWw4xdlEu4sVj3PZfgt1op3revo3HwOu6Xk9ICdsCD8QbSx5Jz5d59-xxVoIx_exgID2oe5-KiogXhgoklveeaLFt9dh-rn4ONyDluExeaHWOG0Rexxv7x-MYt439-xpR_nE0zs58QKenzGn1IWdJc0JV7z9BAbT6NFlOaIEaxRnva2-JyjL3Pdtm2ySGjB41f7gdfMlV1kLxWhBsMlaJiAlqgtbHSKJ-eLYQVERR4t1P5Z2twfrWytI97wvbtzAcVUN_duD2VtSZg-rk2XVsS283mj9y8NCK2rIOuGYCA&state=accountToken0001" 
```

The `code` url query parameter is extracted from the location header above. It has the value  `d037ebc3-121a-4df7-88f6-5a67322bf8eb`

### Example Authorization_code Exchange

```httptrace

---------------------- REQUEST LOG -----------------------
POST  /token  HTTP/1.1
HOST   : modelobank2018.o3bank.co.uk:4201
HEADERS:
                   Accept: */*
            Authorization: Basic NzJiNzlkZGQtNDY3NC00M2JiLTk2YzItOTkyZjc5Y2Q2ZTYyOjM1ZjBkZTliLWFhYTYtNGJmZi1hZDg0LTNlNjUyODU3NDcyMw==
             Content-Type: application/x-www-form-urlencoded
               User-Agent: go-resty/1.10.3 (https://github.com/go-resty/resty)
BODY   :
code=d037ebc3-121a-4df7-88f6-5a67322bf8eb&grant_type=authorization_code&redirect_uri=https%3A%2F%2F127.0.0.1%3A8443%2Fconformancesuite%2Fcallback&scope=accounts
----------------------------------------------------------

```

#### Example Response to Authorization code exchange providing access token

```httptrace

---------------------- RESPONSE LOG -----------------------
STATUS        : 200 OK
RECEIVED AT   : 2019-07-17T15:18:01.995716787+01:00
RESPONSE TIME : 93.960257ms
HEADERS:
               Connection: keep-alive
           Content-Length: 1014
             Content-Type: application/json; charset=utf-8
                     Date: Wed, 17 Jul 2019 14:18:01 GMT
                     Etag: W/"3f6-/+5GQGHyCqs1CqHqbjbWzAnmb2k"
                   Server: nginx/1.14.1
             X-Powered-By: Express
BODY   :
{
   "access_token": "7e602be4-973f-4950-b9b3-c3a8d2adf922",
   "token_type": "Bearer",
   "expires_in": 3600,
   "scope": "openid accounts",
   "id_token": "eyJhbGciOiJQUzI1NiIsImtpZCI6IjVVWXFjdGNOblZkSTl0VXRlYXA0dFNtV241NCJ9.eyJzdWIiOiJhYWMtZjRiZmNlZTYtYTZmZC00YTU5LTllMWEtZjk3YTU3MzEzYmM4Iiwib3BlbmJhbmtpbmdfaW50ZW50X2lkIjoiYWFjLWY0YmZjZWU2LWE2ZmQtNGE1OS05ZTFhLWY5N2E1NzMxM2JjOCIsImlzcyI6Imh0dHBzOi8vbW9kZWxvYmFua2F1dGgyMDE4Lm8zYmFuay5jby51azo0MTAxIiwiYXVkIjoiNzJiNzlkZGQtNDY3NC00M2JiLTk2YzItOTkyZjc5Y2Q2ZTYyIiwiaWF0IjoxNTYzMzczMDgyLCJleHAiOjE1NjMzNzY2ODIsImNfaGFzaCI6IkJwM0xWdlpFTFR5MlJRTk5wWGIxcmciLCJzX2hhc2giOiI0N0RFUXBqOEhCU2EtX1RJbVctNUpBIiwiYWNyIjoidXJuOm9wZW5iYW5raW5nOnBzZDI6c2NhIn0.mQ6PRcWE67dbziZ3gH5QUQDcSCuiZgdcWjYbT1L2lZzb4JQa5WsAkOODt5ZOc30eIwz3FCiN2tyULYmdufdkMfkQwsO7TFlgI8m2Ovd2kTUilYsy_1KJnjDRZFVDHiVfTbGcsmMwUFG7_5TzP6b-RbUhjMZPS-6pr7R8Qc5ziNkE2QBif6xGcBk7pdP3TZEeykJYuR2D0aDnQ8eaaS8rAZen8tifCUaS18YCUIYj5d15CODP4oprHxzpXZx_4Ym1NweCO3i1vm7woNggJN8S-Gi-5kE5dSMhqWij8WVsGImakAh36d-YEzyh8Oa3i9yEBLff1CARPqx1Q0d4EvnvfQ"
}
````
