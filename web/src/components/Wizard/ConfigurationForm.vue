<template>
  <div class="p-3">
    <b-form>
      <!--
              maybe limit file selection to these file types:
              * .key: application/x-iwork-keynote-sffkey
              -->
      <ConfigurationFormFile
        id="signing_private"
        setterMethodNameSuffix="SigningPrivate"
        label="Private Signing Certificate (.key):"
      />
      <!--
              maybe limit file selection to these file types:
              * .pem: application/x-x509-ca-cert
              -->
      <ConfigurationFormFile
        id="signing_public"
        setterMethodNameSuffix="SigningPublic"
        label="Public Signing Certificate (.pem):"
      />
      <!--
              maybe limit file selection to these file types:
              * .key: application/x-iwork-keynote-sffkey
              -->
      <ConfigurationFormFile
        id="transport_private"
        setterMethodNameSuffix="TransportPrivate"
        label="Private Transport Certificate (.key):"
      />
      <!--
              maybe limit file selection to these file types:
              * .pem: application/x-x509-ca-cert
              -->
      <ConfigurationFormFile
        id="transport_public"
        setterMethodNameSuffix="TransportPublic"
        label="Public Transport Certificate (.pem):"
      />
    </b-form>
  </div>
</template>

<script>

import ConfigurationFormFile from './ConfigurationFormFile.vue';

export default {
  name: 'ConfigurationForm',
  components: {
    ConfigurationFormFile,
  },
};
</script>

<style scoped>
</style>
