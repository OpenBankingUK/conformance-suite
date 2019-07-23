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
    <div v-if="serverError">
      <code class="serverError">{{ serverError }}</code>
      <br>
    </div>
    <p>The following response was received from the backend when processing the callback message.</p>
    <div v-if="serverResponse">
      <code class="response">{{ JSON.stringify(serverResponse, null, 2) }}</code>
    </div>
    <br>
    <b-button @click="closeWindow()">
      Close Window
    </b-button>
  </div>

</template>

<script>
import URI from 'urijs';
import api from '../../api';

export default {
  name: 'ConformanceSuiteCallback',
  props: {
    autoCloseOnSuccess: {
      type: Boolean,
      required: false,
      default: false,
    },
  },
  data() {
    return {
      serverResponse: null,
      uri: null,
      serverError: null,
    };
  },
  computed: {
    hasFragment() {
      return api.hasFragment(this.uri);
    },
    hasQuery() {
      return api.hasQuery(this.uri);
    },
    params() {
      return api.consentParams(this.uri);
    },
    isError() {
      return api.isConsentError(this.uri);
    },
    postUrl() {
      return api.consentCallbackEndpoint(this.uri);
    },
  },
  async created() {
    this.uri = new URI(this.$route.fullPath);
    await this.doPost();
    if (this.serverResponse.error == null && !this.isError && this.autoCloseOnSuccess) {
      try {
        window.close();
      } catch (e) {
        // Can ignore this exception as it's just an attempt to close window
      }
    }
  },
  methods: {
    async doPost() {
      try {
        const url = this.postUrl;
        const data = this.params;

        const result = await api.post(url, data);
        const { status } = result;

        if (status !== 200) {
          this.serverError = 'Error processing callback - expected HTTP 200 OK';
        }
        this.serverResponse = await result.json();
      } catch (err) {
        this.serverResponse = err;
      }
    },
    closeWindow() {
      try {
        window.close();
      } catch (e) {
        // Can ignore this exception as it's just an attempt to close window
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
