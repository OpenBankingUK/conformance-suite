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

<script>
import Header from './components/Header.vue';

export default {
  name: 'app',
  components: {
    Header,
  },
};
</script>

<style>
@import "~@/assets/css/app.css";
@import "~@/assets/css/fonts.css";

.app {
  font-family: 'Metropolis', 'Arial', sans-serif;
}
</style>
