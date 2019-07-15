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
      this.openPopup(url, 'PSU Consent', 1074 * 0.75, 800 * 0.75);
      targetElement.innerHTML = 'PSU Consent (Started)'; // eslint-disable-line
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
