<template>
  <a-collapse
    class="endpoint"
    :bordered="false"
    :data-endpoint="endpoint"
    accordion
  >
    <a-collapse-panel :disabled="!value.failed_calls">
      <endpoint-header
        slot="header"
        :data="value"
        :endpoint="endpoint"
      />
      <failures
        v-if="value.failures && value.failures.length"
        :failures="value.failures"
      />
    </a-collapse-panel>
  </a-collapse>
</template>

<script>
import Failures from './Failures';
import EndpointHeader from './EndpointHeader';

export default {
  components: {
    Failures,
    EndpointHeader,
  },
  props: {
    value: {
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
.endpoint .ant-collapse-item-disabled .ant-collapse-header {
  cursor: default;
  color: #000000;
}
</style>
