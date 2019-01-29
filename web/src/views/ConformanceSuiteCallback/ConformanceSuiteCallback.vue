<!--
See:
* `web/src/views/CONFORMANCESUITECALLBACK.md`
* https://openid.net/specs/openid-connect-core-1_0.html#FragmentNotes
-->
<template>
  <div class="d-flex flex-column flex-fill p-3">
    <h1>ConformanceSuiteCallback</h1>
    <hr>
    <h2>Redirect</h2>
    <div>
      <b>hasQuery:</b>
      <code class="has-query">{{ hasQuery }}</code>
    </div>
    <div>
      <b>hasFragment:</b>
      <code class="has-fragment">{{ hasFragment }}</code>
    </div>
    <div>
      <b>params:</b>
      <code class="params">{{ params }}</code>
    </div>
    <div>
      <b>isError:</b>
      <code class="is-error">{{ isError }}</code>
    </div>
    <hr>
    <h2>Response</h2>
    <div>
      <code class="response">{{ response }}</code>
    </div>
  </div>
</template>

<script>
import URI from 'urijs';
import * as _ from 'lodash';
import api from '../../api/apiUtil';

export default {
  name: 'ConformanceSuiteCallback',
  data() {
    return {
      response: null,
      uri: null,
    };
  },
  computed: {
    hasFragment() {
      // http://localhost:8080/conformancesuite/callback#code=a052c795-742d-415a-843f-8a4939d740d1&scope=openid%20accounts&id_token=eyJ0eXAiOiJKV1QiLCJraWQiOiJGb2w3SXBkS2VMWm16S3RDRWdpMUxEaFNJek09IiwiYWxnIjoiRVMyNTYifQ.eyJzdWIiOiJtYmFuYSIsImF1ZGl0VHJhY2tpbmdJZCI6IjY5YzZkZmUzLWM4MDEtNGRkMi05Mjc1LTRjNWVhNzdjZWY1NS0xMDMzMDgyIiwiaXNzIjoiaHR0cHM6Ly9tYXRscy5hcy5hc3BzcC5vYi5mb3JnZXJvY2suZmluYW5jaWFsL29hdXRoMi9vcGVuYmFua2luZyIsInRva2VuTmFtZSI6ImlkX3Rva2VuIiwibm9uY2UiOiI1YTZiMGQ3ODMyYTlmYjRmODBmMTE3MGEiLCJhY3IiOiJ1cm46b3BlbmJhbmtpbmc6cHNkMjpzY2EiLCJhdWQiOiI1NGY2NDMwOS00MzNkLTQ2MTAtOTVkMi02M2QyZjUyNTM0MTIiLCJjX2hhc2giOiIxbGt1SEFuaVJDZlZNS2xEc0pxTTNBIiwib3BlbmJhbmtpbmdfaW50ZW50X2lkIjoiQTY5MDA3Nzc1LTcwZGQtNGIyMi1iZmM1LTlkNTI0YTkxZjk4MCIsInNfaGFzaCI6ImZ0OWRrQTdTWXdlb2hlZXpjOGFHeEEiLCJhenAiOiI1NGY2NDMwOS00MzNkLTQ2MTAtOTVkMi02M2QyZjUyNTM0MTIiLCJhdXRoX3RpbWUiOjE1Mzk5NDM3NzUsInJlYWxtIjoiL29wZW5iYW5raW5nIiwiZXhwIjoxNTQwMDMwMTgxLCJ0b2tlblR5cGUiOiJKV1RUb2tlbiIsImlhdCI6MTUzOTk0Mzc4MX0.8bm69KPVQIuvcTlC-p0FGcplTV1LnmtacHybV2PTb2uEgMgrL3JNA0jpT2OYO73r3zPC41mNQlMDvVOUn78osQ&state=5a6b0d7832a9fb4f80f1170a
      return !_.isEmpty(this.uri.fragment());
    },
    hasQuery() {
      // http://localhost:8080/conformancesuite/callback?code=a052c795-742d-415a-843f-8a4939d740d1&scope=openid%20accounts&id_token=eyJ0eXAiOiJKV1QiLCJraWQiOiJGb2w3SXBkS2VMWm16S3RDRWdpMUxEaFNJek09IiwiYWxnIjoiRVMyNTYifQ.eyJzdWIiOiJtYmFuYSIsImF1ZGl0VHJhY2tpbmdJZCI6IjY5YzZkZmUzLWM4MDEtNGRkMi05Mjc1LTRjNWVhNzdjZWY1NS0xMDMzMDgyIiwiaXNzIjoiaHR0cHM6Ly9tYXRscy5hcy5hc3BzcC5vYi5mb3JnZXJvY2suZmluYW5jaWFsL29hdXRoMi9vcGVuYmFua2luZyIsInRva2VuTmFtZSI6ImlkX3Rva2VuIiwibm9uY2UiOiI1YTZiMGQ3ODMyYTlmYjRmODBmMTE3MGEiLCJhY3IiOiJ1cm46b3BlbmJhbmtpbmc6cHNkMjpzY2EiLCJhdWQiOiI1NGY2NDMwOS00MzNkLTQ2MTAtOTVkMi02M2QyZjUyNTM0MTIiLCJjX2hhc2giOiIxbGt1SEFuaVJDZlZNS2xEc0pxTTNBIiwib3BlbmJhbmtpbmdfaW50ZW50X2lkIjoiQTY5MDA3Nzc1LTcwZGQtNGIyMi1iZmM1LTlkNTI0YTkxZjk4MCIsInNfaGFzaCI6ImZ0OWRrQTdTWXdlb2hlZXpjOGFHeEEiLCJhenAiOiI1NGY2NDMwOS00MzNkLTQ2MTAtOTVkMi02M2QyZjUyNTM0MTIiLCJhdXRoX3RpbWUiOjE1Mzk5NDM3NzUsInJlYWxtIjoiL29wZW5iYW5raW5nIiwiZXhwIjoxNTQwMDMwMTgxLCJ0b2tlblR5cGUiOiJKV1RUb2tlbiIsImlhdCI6MTUzOTk0Mzc4MX0.8bm69KPVQIuvcTlC-p0FGcplTV1LnmtacHybV2PTb2uEgMgrL3JNA0jpT2OYO73r3zPC41mNQlMDvVOUn78osQ&state=5a6b0d7832a9fb4f80f1170a
      return !_.isEmpty(this.uri.query());
    },
    params() {
      if (this.hasFragment) {
        return URI.parseQuery(this.uri.fragment());
      }
      return URI.parseQuery(this.uri.query());
    },
    isError() {
      // http://localhost:8080/conformancesuite/callback?error_description=JWT%20invalid.%20Expiration%20time%20incorrect.&state=5a6b0d7832a9fb4f80f1170a&error=invalid_request
      if (!_.isEmpty(_.get(this, 'params.error'))) {
        return true;
      }
      return false;
    },
    postUrl() {
      if (this.isError) {
        return '/api/redirect/error';
      }
      if (this.hasFragment) {
        return '/api/redirect/fragment/ok';
      }
      if (this.hasQuery) {
        return '/api/redirect/query/ok';
      }

      throw new Error(`invalid state, uri: ${this.uri.toString()}, isError: ${this.isError}, hasFragment: ${this.hasFragment}, hasQuery: ${this.hasQuery}`);
    },
  },
  async created() {
    this.uri = new URI(this.$route.fullPath);
    await this.doPost();
  },
  methods: {
    async doPost() {
      try {
        const url = this.postUrl;
        const data = this.params;

        const result = await api.post(url, data);
        this.response = await result.json();
      } catch (err) {
        this.response = err;
      }
    },
  },
};
</script>

<style scoped>
code::before {
  content: " ";
}
</style>
