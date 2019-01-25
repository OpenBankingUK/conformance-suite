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
      if (this.notifications.length === 0) {
        return 'There are no notifications';
      }

      let result = '<ul>';
      for (let i = 0; i < this.notifications.length; i += 1) {
        const n = this.notifications[i];

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
      }
      return `${result}</ul>`;
    },
  },
};
</script>
