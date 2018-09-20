<template>
  <div>
    <p>
      First name: {{ profile.first_name }}
    </p>
    <p>
      Last name: {{ profile.last_name }}
    </p>
    <p>
      Email: {{ profile.email }}
    </p>
  </div>
</template>

<script>
import { mapState } from 'vuex';

export default {
  name: 'profile',
  computed: {
    ...mapState({
      profile: state => state.user.profile,
    }),
  },
};
</script>
