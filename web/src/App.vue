<template>
  <div
    id="app"
    class="d-flex flex-column app">
    <Header/>
    <div class="d-flex flex-column flex-fill">
      <router-view/>
    </div>
  </div>
</template>

<style lang="scss">
@import 'assets/css/fonts.css';

.app {
  font-family: 'Metropolis', 'Arial', sans-serif;
}

</style>

<script>
import Header from './components/Header.vue';
import Wizard from './components/Wizard.vue';

export default {
  name: 'app',
  components: {
    Wizard,
    Header,
  },
};
</script>
