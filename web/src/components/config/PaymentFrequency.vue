<template>
  <b-container>
    <b-row>
      <b-col lg="12">
        <b-form-group
          :label="label"
          :description="description"
          label-size="sm">
          <b-input-group
            prepend="Frequency_1"
            size="sm">
            <b-form-input
              id="payment_frequency"
              v-model="payment_frequency"
              :state="payment_frequency_state"
              type="text"
            />
          </b-input-group>
          <small class="text-muted">{{ help_text }}</small>
        </b-form-group>
      </b-col>
    </b-row>
  </b-container>
</template>

<script>
import { frequency1Pattern, hasValue, isFrequency1 } from '../../store/modules/config/frequency';

export default {
  name: 'PaymentFrequency',
  props: {
    compatibilityLabel: {
      type: Boolean,
      default: false,
    },
  },
  computed: {
    label() {
      return this.compatibilityLabel ? 'Legacy/v3 Standing Order Frequency' : 'Frequency';
    },
    description() {
      return this.compatibilityLabel
        ? 'Legacy scalar Frequency_1 value used by v3 standing-order payloads.'
        : 'Legacy scalar Frequency_1 value used by v3 standing-order payloads.';
    },
    help_text() {
      return `Must match Frequency_1, for example EvryDay, IntrvlDay:02, IntrvlWkDay:01:03, or QtrDay:ENGLISH. Pattern: ${frequency1Pattern}`;
    },
    payment_frequency: {
      get() {
        return this.$store.state.config.configuration.payment_frequency;
      },
      set(value) {
        this.$store.commit('config/SET_PAYMENT_FREQUENCY', value);
      },
    },
    payment_frequency_state() {
      if (!hasValue(this.payment_frequency)) {
        return null;
      }
      return isFrequency1(this.payment_frequency);
    },
  },
};
</script>

<style scoped>
</style>
