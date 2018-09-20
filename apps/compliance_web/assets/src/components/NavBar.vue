<template>
  <a-layout-header class="navbar">
    <a-dropdown
      v-if="signedIn"
      :trigger="['click']"
    >
      <a
        class="avatar ant-dropdown-link"
        href="#">
        <a-avatar
          :src="profile.avatar ? profile.avatar : ''"
          :icon="!profile.avatar ? 'user' : ''"
          size="large"
        />
      </a>
      <a-menu
        slot="overlay"
        class="dropdown user-menu">
        <a-menu-item key="0">
          <router-link to="/profile">Profile</router-link>
        </a-menu-item>
        <a-menu-divider />
        <a-menu-item key="1">
          <a @click="signOut()">Sign out</a>
        </a-menu-item>
      </a-menu>
    </a-dropdown>
  </a-layout-header>
</template>

<script>
import { mapState, mapActions } from 'vuex';

export default {
  computed: {
    ...mapState({
      signedIn: state => state.user.signedIn,
      profile: state => state.user.profile,
    }),
  },
  methods: {
    ...mapActions('user', ['signOut']),
  },
};
</script>

<style>
.navbar {
  background: #fff;
  display: flex;
  align-items: center;
  justify-content: flex-end;
  padding: 0 25px;
}
.ant-dropdown-link {
  display: flex;
  height: 63px;
  align-items: center;
}
.dropdown {
  width: 120px;
}
</style>
