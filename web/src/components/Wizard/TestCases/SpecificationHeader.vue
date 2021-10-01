<template>
  <div>
    <h6>API Specification</h6>
    <b-table
      :items="[ apiSpecificationWithConsentUrls ]"
      :fields="tableFields"
      small
      fixed
      stacked
    >
      <template
        slot="name"
        slot-scope="data">
        <a
          :href="data.item.url"
          target="_blank">
          {{ data.item.name }}
        </a>
      </template>
      <template
        slot="schema"
        slot-scope="data"
      >
        <a
          :href="data.item.schemaVersion"
          target="_blank">
          {{ data.item.version }}
        </a>
      </template>
      <template
        slot="consentUrls"
        slot-scope="data">
        <template v-for="(url, index) in data.value">
          <a
            :key="url"
            :title="url"
            class="psu-consent-link"
            href="#"
            @click="startPsuConsent(url, $event.target)">
            PSU Consent
          </a>
          <span :key="'s' + index">{{ acquired(tokenName(url)) }}</span>
          <br :key="index">
        </template>
      </template>
    </b-table>
  </div>
</template>

<script>
import { createNamespacedHelpers } from 'vuex';
import axios from 'axios';

const {
  mapGetters,
  mapState,
} = createNamespacedHelpers('testcases');

export default {
  name: 'SpecificationHeader',
  props: {
    // Example value for `apiSpecification`.
    // {
    //     "name": "Account and Transaction API Specification",
    //     "url": "https://openbanking.atlassian.net/wiki/spaces/DZ/pages/642090641/Account+and+Transaction+API+Specification+-+v3.0",
    //     "version": "v3.0",
    //     "schemaVersion": "https://raw.githubusercontent.com/OpenBankingUK/read-write-api-specs/v3.0.0/dist/account-info-swagger.json"
    // }
    apiSpecification: {
      type: Object,
      required: true,
    },
  },
  computed: {
    ...mapState([
      'consentUrls',
    ]),
    specConsentUrls() {
      return this.consentUrls[this.apiSpecification.name];
      // Uncomment below and comment line above to test before backend consent URL changes finished:
      // return [`http://example.com/${this.apiSpecification.name}/1`, `http://example.com/${this.apiSpecification.name}/2`];
    },
    apiSpecificationWithConsentUrls() {
      const consentUrls = this.specConsentUrls;
      return Object.assign({ consentUrls }, this.apiSpecification);
    },
    hasConsentUrls() {
      return this.specConsentUrls && this.specConsentUrls.length > 0;
    },
    /**
     * Fields to display in API Specification Table.
     * See documentation: https://bootstrap-vue.js.org/docs/components/table#fields-column-definitions-
     */
    tableFields() {
      const fields = {
        name: {
          tdClass: 'table-data-breakable api-specification-table',
        },
        schema: {
          tdClass: 'table-data-breakable api-specification-table',
        },
      };
      if (this.hasConsentUrls) {
        Object.assign(fields, {
          consentUrls: {
            tdClass: 'table-data-breakable api-specification-table',
          },
        });
      }
      return fields;
    },
  },
  methods: {
    ...mapGetters([
      'tokenAcquired',
    ]),
    acquired(tokenName) {
      if (this.tokenAcquired()(tokenName)) {
        return ' Acquired';
      }
      return '';
    },
    tokenName(url) {
      const u = new URL(url);
      const state = u.searchParams.get('state');
      if (state) {
        // returns state value, e.g. 'Token0001' or 'Token0002'
        return state;
      }
      return null;
    },
    startPsuConsent(url, targetElement) {
      // this.openPopup(url, 'PSU Consent', 1074 * 0.75, 800 * 0.75);
      // targetElement.innerHTML = shortURL; // eslint-disable-line

      const shortUrl = this.getShortUrl(url);
      shortUrl.then(function(res) {
        console.log(res);
      }, function(err) {
        console.log(err);
      });

      // const nonce = this.getNonce('eyJhbGciOiJQUzI1NiIsImtpZCI6ImZhb0JPSGI1bV9EZFBqZ1JCTG9mX2hQMU52dyJ9.eyJzdWIiOiJ2cnAtMS1mNjlhOWM4ZS04MjFhLTQxZjktYTdkMy05NGNlYThhYzNjNjEiLCJvcGVuYmFua2luZ19pbnRlbnRfaWQiOiJ2cnAtMS1mNjlhOWM4ZS04MjFhLTQxZjktYTdkMy05NGNlYThhYzNjNjEiLCJwc3VfaWRlbnRpZmllcnMiOnsidXNlcklkIjoiNzAwMDAxMDAwMDAwMDAwMDAwMDAwMDAyIiwiY29tcGFueUlkIjoiMTIzNDUifSwiaXNzIjoiaHR0cHM6Ly9vYjE5LWF1dGgxLXVpLm8zYmFuay5jby51ayIsImF1ZCI6ImFlOTRmMzNkLTBlNTctNDVjYS1hMmZhLWQzMTEwMjVjODViNSIsImlhdCI6MTYzMTYxMTgwNywiZXhwIjoxNjMxNjE1NDA3LCJub25jZSI6ImMxYTA2MjFkLWFjMjQtNDM2My05ZmI0LTcxYmIxMTA5YzY4ZiIsImF1dGhfdGltZSI6MTYzMTYxMTgwNywiYXpwIjoiYWU5NGYzM2QtMGU1Ny00NWNhLWEyZmEtZDMxMTAyNWM4NWI1IiwicmVmcmVzaF90b2tlbl9leHBpcmVzX2F0IjoxNjMyMDAwNjA3LCJjX2hhc2giOiJDMERVODBvMkFydVNDeTVqeUwtWXFRIiwic19oYXNoIjoiY00xdVR1OE5zdWk4MGVZVTVXb0hMUSIsImFjciI6InVybjpvcGVuYmFua2luZzpwc2QyOnNjYSJ9.tfpKCJlN9ZE99BHEF0UGekISnvFEe26OmGNSlz2S1SmWXmkTuKD1RYQeK3u2FryY7_d9sWWdWdBwtFd_NU_N8mB1UfvV_7cd33jTvBU_MgiaRSRx0wIFCu7ILiuJQaUWQ7cW79rSvRZlGtFRcDdrZHze7wT076uPBzK08Aq8rjNTMcxMR4cb56MKrzMpCtEawthKjHwlk8liiPTK8ELq9Pb67XINv8hJAXBKDxPW0aTFaOyQQV4DIHFSjzVh1PufbphvGklb5viG5_gLICO5qWDWjC2v-REvMCxYMjeoE4Wa7yhEhgg1q9m5WCgGYCvQYfu8uakqFhgLs2_AoyV6wQ');
      // nonce.then(function(res) {
      //   console.log(res);
      //   const callback = this.getCallback(nonce);
      //   console.log(callback);
      // }, function(err) {
      //   console.log(err);
      // });
      const jwt = 'eyJhbGciOiJQUzI1NiIsImtpZCI6ImZhb0JPSGI1bV9EZFBqZ1JCTG9mX2hQMU52dyJ9.eyJzdWIiOiJ2cnAtMS1mNjlhOWM4ZS04MjFhLTQxZjktYTdkMy05NGNlYThhYzNjNjEiLCJvcGVuYmFua2luZ19pbnRlbnRfaWQiOiJ2cnAtMS1mNjlhOWM4ZS04MjFhLTQxZjktYTdkMy05NGNlYThhYzNjNjEiLCJwc3VfaWRlbnRpZmllcnMiOnsidXNlcklkIjoiNzAwMDAxMDAwMDAwMDAwMDAwMDAwMDAyIiwiY29tcGFueUlkIjoiMTIzNDUifSwiaXNzIjoiaHR0cHM6Ly9vYjE5LWF1dGgxLXVpLm8zYmFuay5jby51ayIsImF1ZCI6ImFlOTRmMzNkLTBlNTctNDVjYS1hMmZhLWQzMTEwMjVjODViNSIsImlhdCI6MTYzMTYxMTgwNywiZXhwIjoxNjMxNjE1NDA3LCJub25jZSI6ImMxYTA2MjFkLWFjMjQtNDM2My05ZmI0LTcxYmIxMTA5YzY4ZiIsImF1dGhfdGltZSI6MTYzMTYxMTgwNywiYXpwIjoiYWU5NGYzM2QtMGU1Ny00NWNhLWEyZmEtZDMxMTAyNWM4NWI1IiwicmVmcmVzaF90b2tlbl9leHBpcmVzX2F0IjoxNjMyMDAwNjA3LCJjX2hhc2giOiJDMERVODBvMkFydVNDeTVqeUwtWXFRIiwic19oYXNoIjoiY00xdVR1OE5zdWk4MGVZVTVXb0hMUSIsImFjciI6InVybjpvcGVuYmFua2luZzpwc2QyOnNjYSJ9.tfpKCJlN9ZE99BHEF0UGekISnvFEe26OmGNSlz2S1SmWXmkTuKD1RYQeK3u2FryY7_d9sWWdWdBwtFd_NU_N8mB1UfvV_7cd33jTvBU_MgiaRSRx0wIFCu7ILiuJQaUWQ7cW79rSvRZlGtFRcDdrZHze7wT076uPBzK08Aq8rjNTMcxMR4cb56MKrzMpCtEawthKjHwlk8liiPTK8ELq9Pb67XINv8hJAXBKDxPW0aTFaOyQQV4DIHFSjzVh1PufbphvGklb5viG5_gLICO5qWDWjC2v-REvMCxYMjeoE4Wa7yhEhgg1q9m5WCgGYCvQYfu8uakqFhgLs2_AoyV6wQ'
      const nonce = this.getNonce(jwt);
      nonce.then((res) => {
        console.log(res);
        const callback = this.getCallback(res);
        callback.then((res2) => {
          console.log(res2);
        })
      }, function(err) {
        console.log(err);
      });
    },
    async getShortUrl(url) {
      const req = await axios.post('https://ok3hqencc6.execute-api.eu-west-1.amazonaws.com/staging', {'url': url})
      if (req.status === 200) {
        return req.data
      }
      return false
      // .then(function (response) {
      //   console.log(response);
      //   })
      //   .catch(function (error) {
      //     console.log(error);
      //     });
    },
    async getNonce(id_token) {
      const params = new URLSearchParams([['id_token', id_token]]);
      const req = await axios.get('https://ok3hqencc6.execute-api.eu-west-1.amazonaws.com/staging', { params });
      if (req.status === 200) {
        return req.data
      }
      return false
      // const params = new URLSearchParams([['id_token', id_token]]);
      // const res = await axios.get('https://ok3hqencc6.execute-api.eu-west-1.amazonaws.com/staging', { params });
      // if (res.status === 200) {
      //   return res.data
      // }
      // return false
    },
    getNonce2(id_token) {
      const params = new URLSearchParams([['id_token', id_token]]);
      axios.get('https://ok3hqencc6.execute-api.eu-west-1.amazonaws.com/staging', { params })
      .then(function (response) {
        console.log(response);
        })
        .catch(function (error) {
          console.log(error);
          });
    },
    async getCallback(nonce) {
      const params = new URLSearchParams([['nonce', nonce]]);
      const req = await axios.get('https://ok3hqencc6.execute-api.eu-west-1.amazonaws.com/staging', { params });
      if (req.status === 200) {
        return req.data
      }
      return false
    },
    openPopup(url, title, w, h) {
      // Open a window popup via Javascript and supports single/dual displays
      // Credit: https://stackoverflow.com/a/16861050/225885
      // Fixes dual-screen position                         Most browsers      Firefox
      const dualScreenLeft = window.screenLeft !== undefined ? window.screenLeft : window.screenX;
      const dualScreenTop = window.screenTop !== undefined ? window.screenTop : window.screenY;

      const clientWidth = document.documentElement.clientWidth ? document.documentElement.clientWidth : window.screen.width;
      const width = window.innerWidth ? window.innerWidth : clientWidth;
      const clientHeight = document.documentElement.clientHeight ? document.documentElement.clientHeight : window.screen.height;
      const height = (window.innerHeight ? window.innerHeight : clientHeight) + 15;

      const systemZoom = width / window.screen.availWidth;
      const left = (width - w) / 2 / (systemZoom + dualScreenLeft);
      const top = (height - h) / 2 / (systemZoom + dualScreenTop);

      const wLocation = `width=${w / systemZoom}, height=${h / systemZoom}, top=${top}, left=${left}`;
      const wParams = `scrollbars=yes, ${wLocation}`;
      const newWindow = window.open(url, title, wParams);

      if (newWindow != null && window.focus) {
        newWindow.focus();
      }
    },
  },
};
</script>

<style scoped>
/*
 * Don't remove the `/deep/` here.
 *
 * This rule ensures values such as 'https://openbanking.atlassian.net/wiki/spaces/DZ/pages/642090641/Account+and+Transaction+API+Specification+-+v3.0'
 * in the table are broken into separate lines.
 */
.test-case /deep/ .table-data-breakable {
  word-break: break-all;
}

.test-case /deep/ .api-specification-table {
  grid-template-columns: 20% auto !important;
}
</style>
