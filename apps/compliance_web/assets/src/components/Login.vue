<template>
  <a-row
    type="flex"
    align="middle"
    justify="center"
    class="h100">
    <a-col
      :xs="{ span: 24 }"
      :md="{ span: 12 }"
      :xll="{ span: 8 }">
      <div class="centered">
        <img :src="logoDark">
      </div>
      <a-divider>Functional Conformance Suite</a-divider>
      <a-card
        :bordered="false"
        class="login">
        <p>
          Test ASPSPs APIs for functional conformance with
          the Open Banking API Read/Write specifications.
        </p>
        <div class="centered">
          <a-button
            type="primary"
            size="large"
            icon="google"
            @click="signIn()">Sign in with Google</a-button>
        </div>
      </a-card>
    </a-col>
  </a-row>
</template>

<script>
import { mapActions } from 'vuex';
import logoDark from '../../static/images/logo-dark.svg';

export default {
  name: 'login',
  data() {
    return {
      logoDark,
    };
  },
  methods: {
    ...mapActions('user', ['signIn']),
  },
};
</script>
