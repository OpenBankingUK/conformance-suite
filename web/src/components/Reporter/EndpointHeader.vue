<template>
  <div class="endpoint-header">
    <span>
      <a-badge
        v-if="data.failed_calls"
        :count="data.failed_calls"
        :offset="[-8, -2]">
        <a-tag color="red">Failures</a-tag>
      </a-badge>
      <a-tag
        v-else
        color="green">
        Passing
      </a-tag>
      <strong>{{ endpoint }}</strong>
    </span>
    <span>
      <a-tag>Total calls: {{ data.total_calls }}</a-tag>
    </span>
  </div>
</template>

<script>
export default {
  props: {
    data: {
      type: Object,
      default: null,
    },
    endpoint: {
      type: String,
      default: null,
    },
  },
};
</script>

<style>
.endpoint-header {
  display: flex;
  justify-content: space-between;
}
.endpoint-header strong {
  margin-left: 10px;
}
</style>
