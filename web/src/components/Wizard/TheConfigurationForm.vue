<template>
  <div class="p-3">
    <b-form>
      <!--
              maybe limit file selection to these file types:
              * .key: application/x-iwork-keynote-sffkey
      -->
      <ConfigurationFormFile
        id="signing_private"
        setter-method-name-suffix="SigningPrivate"
        label="Private Signing Certificate (.key):"
        validExtension=".key"
      />
      <!--
              maybe limit file selection to these file types:
              * .pem: application/x-x509-ca-cert
      -->
      <ConfigurationFormFile
        id="signing_public"
        setter-method-name-suffix="SigningPublic"
        label="Public Signing Certificate (.pem):"
        validExtension=".pem"
      />
      <!--
              maybe limit file selection to these file types:
              * .key: application/x-iwork-keynote-sffkey
      -->
      <ConfigurationFormFile
        id="transport_private"
        setter-method-name-suffix="TransportPrivate"
        label="Private Transport Certificate (.key):"
        validExtension=".key"
      />
      <!--
              maybe limit file selection to these file types:
              * .pem: application/x-x509-ca-cert
      -->
      <ConfigurationFormFile
        id="transport_public"
        setter-method-name-suffix="TransportPublic"
        label="Public Transport Certificate (.pem):"
        validExtension=".pem"
      />
    </b-form>
  </div>
</template>

<script>
import ConfigurationFormFile from './ConfigurationFormFile.vue';

export default {
  name: 'TheConfigurationForm',
  components: {
    ConfigurationFormFile,
  },
};
</script>

<style scoped>
</style>
