<template>
  <b-container>
    <b-row>
      <b-col lg="12">
        <b-form-group
          label="Frequency"
          label-size="sm" />
      </b-col>
      <b-col lg="4">
        <b-input-group
          prepend="Schedule Code"
          size="sm">
          <b-form-select
            id="schedule_code_selected"
            v-model="schedule_code_selected"
            :options="options"
            label="Schedule Code"
            required
            @change="on_schedule_code_change"
          />
        </b-input-group>
      </b-col>
      <b-col lg="8">
        <b-input-group
          v-if="should_display_schedule_code_value_input"
          :append="schedule_code_regex"
          prepend="Schedule Code Value"
          size="sm"
        >
          <b-form-input
            id="schedule_code_value"
            v-model="schedule_code_value"
            :state="valid_schedule_code_value"
            label="Schedule Code Value"
            required
            type="text"
            @update="on_schedule_code_value_update"
          />
        </b-input-group>
      </b-col>
    </b-row>
  </b-container>
</template>

<script>
import * as _ from 'lodash';

const isNotEmpty = value => !_.isEmpty(value);

const validator = {
  regex: /^(EvryDay)$|^(EvryWorkgDay)$|^(IntrvlWkDay:0[1-9]:0[1-7])$|^(WkInMnthDay:0[1-5]:0[1-7])$|^(IntrvlMnthDay:(0[1-6]|12|24):(-0[1-5]|0[1-9]|[12][0-9]|3[01]))$|^(QtrDay:(ENGLISH|SCOTTISH|RECEIVED))$/,
  frequencies: {
    EvryDay: /^$/,
    EvryWorkgDay: /^$/,
    IntrvlWkDay: /^0[1-9]:0[1-7]$/,
    WkInMnthDay: /^0[1-5]:0[1-7]$/,
    IntrvlMnthDay: /^(0[1-6]|12|24):(-0[1-5]|0[1-9]|[12][0-9]|3[01])$/,
    QtrDay: /^(ENGLISH|SCOTTISH|RECEIVED)$/,
  },
};

export default {
  name: 'PaymentFrequency',
  data() {
    const options = [{ value: null, text: '-- Please select an option --', disabled: true }].concat(
      _.map(validator.frequencies, (value, key) => ({ value: key, text: key, disabled: false })),
    );

    // Don't do `this.payment_frequency`, it won't work.
    const { payment_frequency = null } = this.$store.state.config.configuration;
    // => IntrvlWkDay:01:03

    if (_.isEmpty(payment_frequency)) {
      return {
        schedule_code_selected: null,
        schedule_code_value: null,
        options,
      };
    }

    const frequencies_with_no_input = /^(EvryDay)$|^(EvryWorkgDay)$/;
    const does_not_require_input = isNotEmpty(payment_frequency.match(frequencies_with_no_input));
    if (does_not_require_input) {
      return {
        schedule_code_selected: payment_frequency,
        schedule_code_value: null,
        options,
      };
    }

    const valid_frequency = isNotEmpty(payment_frequency.match(validator.regex));
    if (!valid_frequency) {
      return {
        schedule_code_selected: null,
        schedule_code_value: null,
        options,
      };
    }

    const schedule_code_selected = payment_frequency.substr(0, payment_frequency.indexOf(':'));
    // => IntrvlWkDay
    const schedule_code_value = payment_frequency.substr(payment_frequency.indexOf(':') + 1);
    // => 01:03

    return {
      schedule_code_selected,
      schedule_code_value,
      options,
    };
  },
  computed: {
    payment_frequency: {
      get() {
        return this.$store.state.config.configuration.payment_frequency;
      },
      set(value) {
        this.$store.commit('config/SET_PAYMENT_FREQUENCY', value);
      },
    },
    requires_input() {
      // If empty regex matches the `schedule_code_selected` we don't require further input from the user.
      const regex = validator.frequencies[this.schedule_code_selected];
      const requires_input = _.isEmpty(''.match(regex));
      return requires_input;
    },
    should_display_schedule_code_value_input() {
      // Nothing selected, i.e., `Select Schedule Code` so don't display an additional input area.
      if (_.isEmpty(this.schedule_code_selected)) {
        return false;
      }

      return this.requires_input;
    },
    schedule_code_regex() {
      if (_.isEmpty(this.schedule_code_selected)) {
        return false;
      }
      const regex = validator.frequencies[this.schedule_code_selected];

      return regex.toString();
    },
    valid_schedule_code_value() {
      if (_.isEmpty(this.schedule_code_selected)) {
        return false;
      }

      const regex = validator.frequencies[this.schedule_code_selected];
      const value = this.schedule_code_value || '';
      const valid = isNotEmpty(value.match(regex));

      return valid;
    },
  },
  methods: {
    on_schedule_code_change() {
      if (this.requires_input) {
        this.payment_frequency = null;
        return;
      }

      // If the selected payment frequency does not require input, it's valid so commit it to the store.
      const regex = validator.frequencies[this.schedule_code_selected];
      const value = this.schedule_code_selected || '';
      const valid = isNotEmpty(value.match(regex));
      if (!valid) {
        return;
      }

      this.payment_frequency = this.schedule_code_selected;
    },
    on_schedule_code_value_update() {
      // Only commit new payment frequency if the value is valid.
      if (!this.valid_schedule_code_value) {
        return;
      }

      // Final payment_frequency looks like `IntrvlWkDay:01:03` when input is required.
      const payment_frequency = [this.schedule_code_selected, this.schedule_code_value].join(':');
      const valid = isNotEmpty(payment_frequency.match(validator.regex));
      if (!valid) {
        return;
      }

      this.payment_frequency = payment_frequency;
    },
  },
};
</script>

<style scoped>
</style>
