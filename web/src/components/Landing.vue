<template>
  <div class="centered">
    <b-jumbotron
      header="Oops!"
      lead="This release is not intended to be executable, come back soon."
    >
      <p>This pre-alpha build introduces some of the concepts to participants before the official release.</p>
      <p>
        We are working hard to foster an open and collaborative tool
        and welcome your feedback to help us develop the best tool.
      </p>
      <hr class="my-4">
      <div class="d-flex flex-row justify-content-around">
        <b-btn
          variant="primary"
          target="_blank"
          href="https://bitbucket.org/openbankingteam/conformance-suite/src/develop/README.md"
        >More Info</b-btn>
        <b-btn
          variant="success"
          to="/wizard/continue-or-start">Start a validation run</b-btn>
      </div>
    </b-jumbotron>
  </div>
</template>

<style scoped>
</style>

<script>
export default {
  name: 'Landing',
};
</script>
