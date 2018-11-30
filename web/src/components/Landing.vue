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
      <b-jumbotron
        header="Oops!"
        lead="This release is not intended to be executable, come back soon.">
        <p>This pre-alpha build introduces some of the concepts to participants before the official release.</p>
        <p>We are working hard to foster an open and collaborative tool
        and welcome your feedback to help us develop the best tool.</p>

        <b-btn
          variant="primary"
          target="_blank"
          href="https://bitbucket.org/openbankingteam/conformance-suite/src/develop/README.md"
        >More Info</b-btn>
      </b-jumbotron>
      <!-- <router-link to="/wizard" class="go-to-config">Start a validation run</router-link> -->
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
