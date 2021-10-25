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
          <span :key="'s' + index">
          <span><b-badge variant="primary">{{ acquired(tokenName(url)) }}</b-badge> </span>
          <small> <a @click="openPopup(localCallbackUrls[tokenName(url)])"
                     v-show="mobileConsent && localCallbackUrls[tokenName(url)] != null && !acquired(tokenName(url))"
                     href="#"
                     title="Use when popup blocker cause the callback handling to stop."> Open local callback url</a></small>
          </span>
          <br :key="index">
        </template>
      </template>
    </b-table>
    <div>
      <b-modal v-model="modalShow">
        <div style="text-align: center"><img :src=qrCodeUrl alt="QR code"></div>
      </b-modal>
    </div>
  </div>
</template>

<script>
import {createNamespacedHelpers, mapGetters} from 'vuex';
import axios from 'axios';

const {
  mapState,
} = createNamespacedHelpers('testcases');

export default {
  name: 'SpecificationHeader',
  data: function () {
    return {
      modalShow: false,
      qrCodeUrl: undefined,
      shortConsentUrls: {},
      localCallbackUrls: {},
    }
  },
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
    ...mapGetters('config', [
      'tokenAcquisition',
      'callbackProxyUrl',
    ]),
    ...mapState([
      'consentUrls',
    ]),
    mobileConsent() {
      return this.tokenAcquisition === 'mobile';
    },
    specConsentUrls() {
      return this.consentUrls[this.apiSpecification.name];
      // Uncomment below and comment line above to test before backend consent URL changes finished:
      // return [`http://example.com/${this.apiSpecification.name}/1`, `http://example.com/${this.apiSpecification.name}/2`];
    },
    apiSpecificationWithConsentUrls() {
      const consentUrls = this.specConsentUrls;
      return Object.assign({consentUrls}, this.apiSpecification);
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
    ...mapGetters('testcases', [
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

      return state ? state : null;
    },
    getCallbackProxyUrl() {
      return this.callbackProxyUrl;
    },
    getUrlShortenerUrl() {
      //seems redundant but provides clear contract and allows replacing without changing code below
      return this.callbackProxyUrl;
    },
    getIdTokenParserUrl() {
      //seems redundant but provides clear contract and allows replacing without changing code below
      return this.callbackProxyUrl;
    },
    getNonceExtractorUrl() {
      //seems redundant but provides clear contract and allows replacing without changing code below
      return this.callbackProxyUrl;
    },
    showQrCode(consentUrl) {
      const tokenName = this.tokenName(consentUrl);
      this.modalShow = true;
    },
    async startPsuConsent(consentUrl, targetElement) {
      const tokenName = this.tokenName(consentUrl);

      if (!this.mobileConsent) {
        this.openPopup(consentUrl, 'PSU Consent', 1074 * 0.75, 800 * 0.75);
        return;
      }

      /**
       * Utility to poll for callback consent
       *
       * @param {function} fn
       * @param {function} validate
       * @param {?number} interval
       * @param {?number} maxAttempts
       */
      async function createPoller(fn, validate, interval = 1000, maxAttempts = null) {
        let attempts = 0;

        const executePoll = async (resolve, reject) => {
          const result = await fn();
          attempts++;

          if (validate(result)) {
            return resolve(result);
          }

          if (maxAttempts != null && attempts === maxAttempts) {
            return reject(new Error(`'Max polling attempts (${maxAttempts}) exceeded!`));
          }

          setTimeout(executePoll, interval, resolve, reject);
        };

        return new Promise(executePoll);
      }

      console.debug(consentUrl)
      const idToken = (new URL(consentUrl)).searchParams.get('request')
      console.debug('idToken', idToken)

      const nonce = await this.getNonceFromToken(idToken);
      console.debug('nonce', nonce)

      const shortUrl = await this.getShortUrl(consentUrl);
      this.shortConsentUrls[tokenName] = shortUrl;
      this.modalShow = true;
      console.debug('Short url', shortUrl);

      //TODO: Replace with an internal service
      const qrCodeUrl = new URL("http://api.qrserver.com/v1/create-qr-code/?size=400x400");
      qrCodeUrl.searchParams.set('data', shortUrl.toString());
      this.qrCodeUrl = qrCodeUrl.toString();

      const pollForCallback = createPoller(
          async () => this.getCallbackPayload(nonce),
          payload => payload && payload.code, // assumes that if we have a code parameter returned, the rest is ok as well
          2000,
          1000
      )
          .then(cbPayload => {
            console.debug(cbPayload);
            const localCallbackUrl = new URL(window.location);
            localCallbackUrl.pathname = '/conformancesuite/callback';
            for (const cbField in cbPayload) {
              if (cbPayload.hasOwnProperty(cbField)) {
                localCallbackUrl.searchParams.set(cbField, cbPayload[cbField]);
              }
            }
            console.debug(localCallbackUrl.toString())

            // We currently reuse existing callback handling logic, routing etc.
            // so we build an URL from received callback params and "redirecting"
            this.localCallbackUrls[tokenName] = localCallbackUrl.toString();
            this.modalShow = false;

            this.openPopup(localCallbackUrl.toString());
          })
          .catch(err => console.error(err));
    },

    /**
     * Executes a requests to internal URL shortening service and returns shortened url
     * The reason we do shortening is to generate easy to scan and not excessively large QR codes
     * @param {string} url
     */
    async getShortUrl(url) {
      const req = await axios.post(this.getUrlShortenerUrl(), {'url': url})

      return req && req.status === 200 ? req.data : false;
    },

    /**
     * Nonce is used as a matching key between consent request and consent callback
     * We currently use an external service to process the JWT and extract nonce string
     *
     * TODO: Replace with either
     * - pure JS implementation
     * - a local call to FCS server
     *
     * @param {string} id_token
     * @return {Promise<string>}
     */
    async getNonceFromToken(id_token) {
      const params = new URLSearchParams([['id_token', id_token]]);
      const req = await axios.get(this.getNonceExtractorUrl(), {params});

      return req.status === 200 ? req.data : undefined;
    },

    getParsedJwt(id_token) {
      /**
       * TODO: Replace with either
       * - pure JS implementation
       * - a local call to FCS server
       *
       * This will parse the id_token JWT to extract the nonce that is used to match callback contents to consent
       *
       * @type {URLSearchParams}
       */
      const params = new URLSearchParams([['id_token', id_token]]);
      axios.get(this.getIdTokenParserUrl(), {params})
          .then(function (response) {
            console.log(response);
          })
          .catch(function (error) {
            console.log(error);
          });
    },

    /**
     * Executes a requests to internal URL shortening service and returns shortened url
     * @param {string} nonce
     */
    async getCallbackPayload(nonce) {
      const params = new URLSearchParams([['nonce', nonce]]);
      try {
        const req = await axios.get(this.getCallbackProxyUrl(), {params});
        if (req.status === 200) {
          return req.data
        }
      } catch (e) {
        console.log(e)
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
