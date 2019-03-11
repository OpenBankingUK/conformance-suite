<template>
  <b-card
    :id="idSelector"
    :title="title"
    style="max-width: 20rem;"
    class="discovery-card mb-2"
    @click="selectDiscovery()"
  >
    <div class="card-image">
      <b-card-img
        :src="image"
        :alt="name"/>
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
      default() {
        return null;
      },
    },
    image: {
      type: String,
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
    idSelector() {
      return this.name.replace(/\./g, '-');
    },
    text() {
      return this.discoveryModel.description;
    },
  },
  methods: {
    ...mapActions('config', ['setDiscoveryModel']),
    selectDiscovery() {
      this.setDiscoveryModel(JSON.stringify({ discoveryModel: this.discoveryModel }));
      // route to discovery configuration
      this.$router.push('discovery-config');
    },
  },
};
</script>

<style scoped>
div.card-image {
  min-height: 140px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-bottom: 1px solid #f6f6f6;
  margin-bottom: 1rem;
}
.card-image > .card-img {
  width: 50%;
}
.discovery-card > .card-body {
  padding-top: 0.5rem;
}
set .card {
  background: #fff;
  box-shadow: 0 20px 50px rgba(0, 0, 0, 0.1);
  border-radius: 10px;
  transition: 0.5s;
}
.card:hover {
  box-shadow: 0 30px 70px rgba(0, 0, 0, 0.2);
}
.card .box {
  position: absolute;
  top: 50%;
  left: 0;
  transform: translateY(-50%);
  text-align: center;
  padding: 20px;
  box-sizing: border-box;
  width: 100%;
}
.card .box .img {
  width: 120px;
  height: 120px;
  margin: 0 auto;
  border-radius: 50%;
  overflow: hidden;
}
.card .box .img img {
  width: 100%;
  height: 100%;
}
.card .box h2 {
  font-size: 20px;
  color: #262626;
  margin: 20px auto;
}
.card .box h2 span {
  font-size: 14px;
  background: #e91e63;
  color: #fff;
  display: inline-block;
  padding: 4px 10px;
  border-radius: 15px;
}
.card .box p {
  color: #262626;
}
.card .box span {
  display: inline-flex;
}
.card .box ul {
  margin: 0;
  padding: 0;
}
.card .box ul li {
  list-style: none;
  float: left;
}
.card .box ul li a {
  display: block;
  color: #aaa;
  margin: 0 10px;
  font-size: 20px;
  transition: 0.5s;
  text-align: center;
}
</style>
