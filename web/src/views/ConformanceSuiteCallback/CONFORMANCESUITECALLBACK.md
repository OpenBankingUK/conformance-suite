# `CONFORMANCESUITECALLBACK.md`
[`./ConformanceSuiteCallback.vue`](./ConformanceSuiteCallback.vue) handles the redirect that happens was part of authorisation.

## redirects
The authorisation response can be contained in a query string or as a fragment.

### query string
```sh
https://client.example.com/cb?code=SplxlOBeZQQYbYS6WxSbIA
               &state=xyz
```

See:
* https://tools.ietf.org/html/rfc6749#section-4.1.2

### fragment
```sh
https://client.example.org/cb#
    code=SplxlOBeZQQYbYS6WxSbIA
    &id_token=eyJ0 ... NiJ9.eyJ1c ... I6IjIifX0.DeWt4Qu ... ZXso
    &state=af0ifjsldkjs
```

See:
* https://openid.net/specs/openid-connect-core-1_0.html#ImplicitCallback
* https://openid.net/specs/openid-connect-core-1_0.html#HybridCallback
* https://openid.net/specs/openid-connect-core-1_0.html#FragmentNotes

### reading
* https://medium.com/@darutk/diagrams-of-all-the-openid-connect-flows-6968e3990660
* https://openid.net/specs/openid-connect-core-1_0.html#AuthResponse
* https://openid.net/specs/openid-connect-core-1_0.html#ImplicitAuthResponse
* https://openid.net/specs/openid-connect-core-1_0.html#HybridAuthResponse
