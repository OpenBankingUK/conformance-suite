<template>
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
    class="notification-bell"
  />
</template>

<script>
import NotificationBell from 'vue-notification-bell';
import { createNamespacedHelpers } from 'vuex';

const { mapGetters } = createNamespacedHelpers('status');

export default {
  name: 'TheNotificationBell',
  components: {
    NotificationBell,
  },
  data() {
    return {
      unreadNotifications: 0,
    };
  },
  computed: {
    ...mapGetters([
      'hasNotifications',
      'notifications',
    ]),
    count() {
      return this.notifications.length;
    },
    notificationText() {
      if (!this.hasNotifications) {
        return 'There are no notifications';
      }

      let result = '<ul>';
      this.notifications.forEach((n) => {
        let target = null;
        let url = null;
        if (n.extURL) {
          url = n.extURL;
          target = '_blank';
        } else if (n.infoURL) {
          url = n.infoURL;
          target = '_self';
        }
        const infoLink = url ? ` <a href="${url}" target="${target}">More info</a>` : '';

        result += `<li>${n.message}${infoLink}</li>`;
      });
      return `${result}</ul>`;
    },
  },
};
</script>

<style scoped>
.notification-bell {
  margin-right: 8px;
  outline : none;
}
</style>
