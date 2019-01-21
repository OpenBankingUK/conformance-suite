<template>
  <b-navbar>
    <b-navbar-brand to="/">Conformance Suite</b-navbar-brand>
    <notification-bell
      v-b-popover.click.blur.bottomleft.html="notificationText"
      id="notifyPopover"
      ref="notify-bell"
      :size="25"
      :count="count"
      :animated="true"
      title="Notifications"
      tabindex="0"
      iconColor="#fff"
    />
  </b-navbar>

</template>

<script>
import NotificationBell from 'vue-notification-bell';

export default {
  name: 'TheHeader',
  components: {
    NotificationBell,
  },
  data() {
    return {
      unreadNotifications: 0,
      notifications: [],
    };
  },
  computed: {
    count() {
      return this.unreadNotifications;
    },
    notificationText() {
      if (this.notifications.length === 0) {
        return 'There are no notifications';
      }

      let result = '<ul>';
      for (let i = 0; i < this.notifications.length; i += 1) {
        result += `<li>${this.notifications[i]}</li>`;
      }
      return `${result}</ul>`;
    },
  },
  methods: {
    pushNotification(message) {
      this.notifications.push(message);
      this.unreadNotifications += 1;
    },
    shown() {
      this.unreadNotifications = 0;
    },
  },

};
</script>

<style scoped>
.navbar {
  background: linear-gradient(90deg, #6180c3, #6180c3);
  box-shadow: 0 6px 25px 0 rgba(38, 50, 56, 0.2);
}

.navbar-brand {
  color: #ffffff;
  font-size: 1rem;
}
</style>
