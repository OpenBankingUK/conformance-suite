<template>
  <div>
    <!-- <p>Lorem ipsum...</p> -->
    <!-- Leaving for next implementation (discovery endpoint) -->
    <!-- <a-input-search -->
    <!--   placeholder="Enter discovery URL" -->
    <!--   @search="onSearch" -->
    <!--   enterButton="Search" -->
    <!--   size="large" /> -->
    <div class="centered">
      <router-link
        to="/config"
        class="go-to-config">Start a validation run</router-link>
    </div>
  </div>
</template>

<script>
export default {
  methods: {
    // leaving for next implementation (discovery endpoint)
    // onSearch(value) {
    //   if (value) {
    //     this.$router.push('/config', value);
    //   }
    // },
  },
};
</script>
