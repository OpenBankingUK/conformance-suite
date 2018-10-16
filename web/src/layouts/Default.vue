<template>
  <div
    v-if="isLoading"
    class="loading h100">
    <div class="logo">
      <img :src="logoDark">
    </div>
    <spinner />
  </div>
  <a-layout
    v-else
    class="h100">
    <side-bar />
    <a-layout>
      <nav-bar />
      <a-layout-content :style="{ margin: '24px 16px 0' }">
        <div :style="{ padding: '24px', background: '#fff' }">
          <slot />
        </div>
      </a-layout-content>
      <a-layout-footer class="centered">
        Functional Conformance Suite © Open Banking 2018
      </a-layout-footer>
    </a-layout>
  </a-layout>
</template>

<script>
import NavBar from '../components/NavBar.vue';
import SideBar from '../components/SideBar.vue';
import Spinner from '../components/Spinner.vue';
import logoDark from '../assets/images/logo-dark.svg';

export default {
  data() {
    return {
      logoDark,
    };
  },
  components: {
    NavBar,
    SideBar,
    Spinner,
  },
  computed: {
    isLoading() {
      return this.$store.getters['user/isLoading'];
    },
  },
};
</script>

<style>
.loading {
  display: flex;
  justify-content: center;
  align-items: center;
  flex-direction: column;
}
.logo {
  text-align: center;
  width: 200px;
  margin-bottom: 10px;
}
</style>
