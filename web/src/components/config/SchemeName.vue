<template>
  <div>
    <b-form-group
      label-for="scheme_name"
      label="SchemeName"
      description="OBExternalAccountIdentification4Code"
    >
      <b-form-select
        v-model="scheme_name_selector"
        :options="[
          'UK.OBIE.BBAN',
          'UK.OBIE.IBAN',
          'UK.OBIE.PAN' ,
          'UK.OBIE.Paym',
          'UK.OBIE.SortCodeAccountNumber',
          'Other'
        ]"
        :state="validSchemeName(scheme_name)"
        required
        @change="scheme_name_selector_change"
      />
    </b-form-group>
    <b-form-group
      v-if="custom_scheme_visible"
      label-for="scheme_name_other"
      label="Custom SchemeName"
      description="OBExternalAccountIdentification4Code"
    >
      <b-form-input
        v-model="scheme_name_other"
        :state="validSchemeName(scheme_name_other)"
        required
        @update="scheme_name_other_update"
      />
    </b-form-group>
  </div>
</template>

<script>
import * as _ from 'lodash';

const ACCOUNT_TYPE_INTERNATIONAL = 'International';
const ACCOUNT_TYPE_LOCAL = 'Local';
const ACCOUNT_TYPE_CBPII = 'CBPII';

export default {
  name: 'SchemeName',
  props: {
    creditorAccountType: {
      type: String,
      required: true,
    },
  },
  data() {
    let scheme_name_selector = null;
    let scheme_name_other = null;

    switch (this.creditorAccountType) {
      case ACCOUNT_TYPE_INTERNATIONAL:
        if (!this.isKnownSchemeName(this.$store.state.config.configuration.international_creditor_account.scheme_name) && this.isNotEmpty(this.$store.state.config.configuration.international_creditor_account.scheme_name)) {
          scheme_name_other = this.$store.state.config.configuration.international_creditor_account.scheme_name;
          scheme_name_selector = 'Other';
        } else {
          scheme_name_selector = this.$store.state.config.configuration.international_creditor_account.scheme_name;
        }
        break;
      case ACCOUNT_TYPE_LOCAL:
        if (!this.isKnownSchemeName(this.$store.state.config.configuration.creditor_account.scheme_name) && this.isNotEmpty(this.$store.state.config.configuration.creditor_account.scheme_name)) {
          scheme_name_other = this.$store.state.config.configuration.creditor_account.scheme_name;
          scheme_name_selector = 'Other';
        } else {
          scheme_name_selector = this.$store.state.config.configuration.creditor_account.scheme_name;
        }
        break;
      case ACCOUNT_TYPE_CBPII:
        if (!this.isKnownSchemeName(this.$store.state.config.configuration.cbpii_debtor_account.scheme_name) && this.isNotEmpty(this.$store.state.config.configuration.cbpii_debtor_account.scheme_name)) {
          scheme_name_other = this.$store.state.config.configuration.cbpii_debtor_account.scheme_name;
          scheme_name_selector = 'Other';
        } else {
          scheme_name_selector = this.$store.state.config.configuration.cbpii_debtor_account.scheme_name;
        }
        break;
      default:
        throw new Error('Unsupported Scheme Account Type');
    }

    return {
      scheme_name_other,
      scheme_name_selector,
      maxSchemeNameLength: 40,
    };
  },
  computed: {
    scheme_name: {
      get() {
        if (this.scheme_name_selector === 'Other' && this.isNotEmpty(this.scheme_name_other)) {
          return this.scheme_name_other;
        }

        return this.scheme_name_selector;
      },
    },
    custom_scheme_visible: {
      get() {
        return this.scheme_name_selector === 'Other';
      },
    },
  },
  methods: {
    isKnownSchemeName(schemeName) {
      return [
        'UK.OBIE.BBAN',
        'UK.OBIE.IBAN',
        'UK.OBIE.PAN',
        'UK.OBIE.Paym',
        'UK.OBIE.SortCodeAccountNumber',
      ].indexOf(schemeName) > -1;
    },
    scheme_name_other_update() {
      if (this.isNotEmpty(this.scheme_name_other) && this.creditorAccountType === ACCOUNT_TYPE_INTERNATIONAL) {
        this.$store.commit('config/SET_INTERNATIONAL_CREDITOR_ACCOUNT_NAME_SCHEME_NAME', this.scheme_name_other);
      }
      if (this.isNotEmpty(this.scheme_name_other) && this.creditorAccountType === ACCOUNT_TYPE_LOCAL) {
        this.$store.commit('config/SET_CREDITOR_ACCOUNT_NAME_SCHEME_NAME', this.scheme_name_other);
      }
      if (this.isNotEmpty(this.scheme_name_other) && this.creditorAccountType === ACCOUNT_TYPE_CBPII) {
        this.$store.commit('config/SET_CBPII_DEBTOR_ACCOUNT_SCHEME_NAME', this.scheme_name_other);
      }
    },
    scheme_name_selector_change() {
      if (this.scheme_name_selector !== 'Other' && this.creditorAccountType === ACCOUNT_TYPE_INTERNATIONAL) {
        this.$store.commit('config/SET_INTERNATIONAL_CREDITOR_ACCOUNT_NAME_SCHEME_NAME', this.scheme_name_selector);
      }
      if (this.scheme_name_selector !== 'Other' && this.creditorAccountType === ACCOUNT_TYPE_LOCAL) {
        this.$store.commit('config/SET_CREDITOR_ACCOUNT_NAME_SCHEME_NAME', this.scheme_name_selector);
      }
      if (this.scheme_name_selector !== 'Other' && this.creditorAccountType === ACCOUNT_TYPE_CBPII) {
        this.$store.commit('config/SET_CBPII_DEBTOR_ACCOUNT_SCHEME_NAME', this.scheme_name_selector);
      }
    },
    isNotEmpty(value) {
      return !_.isEmpty(value);
    },
    validSchemeName(value) {
      return this.isNotEmpty(value) && value.length <= this.maxSchemeNameLength;
    },
  },
};
</script>
