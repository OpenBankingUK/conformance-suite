<template>
  <b-card
    :id="name"
    :title="title"
    style="max-width: 20rem;"
    class="discovery-card mb-2"
    @click="selectDiscovery()"
  >
    <div class="card-image">
      <b-card-img
        :src="imgSrc"
        :alt="name"
      />
    </div>
    <p class="card-text">{{ text }}</p>
  </b-card>
</template>

<script>
import { mapActions } from 'vuex';

export default {
  name: 'DiscoveryTemplateCard',
  props: {
    discoveryModel: {
      type: Object,
      private: true,
      default() {
        return null;
      },
    },
    image: {
      type: String,
      private: true,
      default() {
        return null;
      },
    },
  },
  computed: {
    title() {
      return '';
    },
    name() {
      return this.discoveryModel.name;
    },
    discoveryVersion() {
      return this.discoveryModel.discoveryVersion;
    },
    text() {
      return this.discoveryModel.description;
    },
    imgSrc() {
      return this.image;
    },
  },
  methods: {
    ...mapActions('config', [
      'setDiscoveryModel',
    ]),
    selectDiscovery() {
      this.setDiscoveryModel(JSON.stringify({ discoveryModel: this.discoveryModel }));
      // route to discovery configuration
      this.$router.push('discovery-config');
    },
  },
};
</script>

<style>
div.card-image {
  min-height: 140px;
  display: flex;
  align-items: center;
  justify-content: center;
}
.card-image > .card-img {
  width: 50%;
}
.discovery-card > .card-body {
  padding-top: 0.5rem;
}
</style>
