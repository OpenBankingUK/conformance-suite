<script>
import logo from '../assets/images/logo.svg';

export default {
  data() {
    return {
      logo,
    };
  },
};
</script>

<template>
  <a-layout-sider
    breakpoint="md"
  >
    <div class="logo">
      <img :src="logo">
      <span>Functional Comformance Suite</span>
    </div>
    <a-menu
      theme="dark"
      mode="inline"
    >
      <router-link
        to="/"
        tag="li"
        class="ant-menu-item">
        <a-icon type="home" />
        <span class="nav-text">Home</span>
      </router-link>
      <router-link
        to="/config"
        tag="li"
        class="ant-menu-item">
        <a-icon type="setting" />
        <span class="nav-text">Validation run</span>
      </router-link>
      <router-link
        to="/profile"
        tag="li"
        class="ant-menu-item">
        <a-icon type="user" />
        <span class="nav-text">Profile</span>
      </router-link>
    </a-menu>
  </a-layout-sider>
</template>

<style>
.logo {
  padding: 20px;
  text-align: center;
  text-transform: uppercase;
  color: rgba(255, 255, 255, 0.7);
}
@media (max-width: 768px){
  .logo {
    display: none;
  }
}
</style>
