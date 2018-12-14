<template>
  <div
    id="app"
    class="d-flex flex-column app">
    <b-navbar
      type="dark"
      variant="dark">
      <b-navbar-brand to="/">Functional Conformance Suite</b-navbar-brand>
    </b-navbar>
    <div class="d-flex flex-column flex-fill">
      <router-view/>
    </div>
    <Footer/>
  </div>
</template>

<style lang="scss">
@import 'assets/css/fonts.css';

.app {
  font-family: 'Metropolis', 'Arial', sans-serif;
}

/* div {
  border: 1px solid #F6F6F6;
} */
</style>

<script>
import Footer from './components/Footer.vue';
import Wizard from './components/Wizard.vue';

export default {
  name: 'app',
  components: {
    Wizard,
    Footer,
  },
};
</script>
