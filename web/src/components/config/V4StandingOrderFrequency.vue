<template>
  <b-container>
    <b-row>
      <b-col lg="12">
        <b-form-group
          :label="label"
          label-size="sm"
          description="OBFrequency6 defaults used by v4 standing-order success payloads. Use either a legacy Frequency_1 Type or a structured OBFrequency6Code Type with optional PointInTime or CountPerPeriod. Type-only structured values remain allowed." />
      </b-col>
      <b-col lg="12">
        <b-input-group
          prepend="Legacy Frequency_1 Type"
          size="sm">
          <b-form-input
            id="v4_standing_order_frequency_legacy_type"
            v-model="legacy_type"
            :state="legacy_type_state"
            type="text"
          />
        </b-input-group>
        <small class="text-muted">
          Legacy v3-style frequency values such as EvryDay or IntrvlWkDay:01:03 are mutually exclusive with structured v4 fields.
        </small>
      </b-col>
      <b-col
        class="mt-2"
        lg="4">
        <b-input-group
          prepend="OBFrequency6Code Type"
          size="sm">
          <b-form-select
            id="v4_standing_order_frequency_type"
            v-model="structured_type"
            :options="type_options"
            :state="structured_type_state"
          />
        </b-input-group>
      </b-col>
      <b-col
        class="mt-2"
        lg="4">
        <b-input-group
          prepend="PointInTime"
          append="e.g. '-1', '1', '01'"
          size="sm">
          <b-form-input
            id="v4_standing_order_frequency_point_in_time"
            v-model="point_in_time"
            :state="point_in_time_state"
            type="text"
          />
        </b-input-group>
      </b-col>
      <b-col
        class="mt-2"
        lg="4">
        <b-input-group
          prepend="CountPerPeriod"
          append="Int32"
          size="sm">
          <b-form-input
            id="v4_standing_order_frequency_count_per_period"
            v-model.number="count_per_period"
            :state="count_per_period_state"
            min="1"
            type="number"
          />
        </b-input-group>
      </b-col>
      <b-col
        v-if="legacy_type_has_structured_fields"
        lg="12">
        <small class="text-danger">Legacy Frequency_1 Type cannot be combined with PointInTime or CountPerPeriod.</small>
      </b-col>
      <b-col
        v-if="!mutually_exclusive"
        lg="12">
        <small class="text-danger">PointInTime and CountPerPeriod are mutually exclusive.</small>
      </b-col>
    </b-row>
  </b-container>
</template>

<script>
import {
  hasValue,
  isCountPerPeriod,
  isFrequency1,
  isOBFrequency6Code,
  isPointInTime,
  obFrequency6Codes,
} from '../../store/modules/config/frequency';

const defaultStructuredType = 'WEEK';

export default {
  name: 'V4StandingOrderFrequency',
  props: {
    compatibilityLabel: {
      type: Boolean,
      default: false,
    },
  },
  computed: {
    label() {
      return this.compatibilityLabel ? 'V4 Standing Order Frequency' : 'Standing Order Frequency';
    },
    frequency() {
      return this.$store.state.config.configuration.v4_standing_order_frequency || {};
    },
    type_options() {
      return [{ value: null, text: '-- Please select an option --', disabled: true }]
        .concat(obFrequency6Codes.map(value => ({ value, text: value })));
    },
    is_structured_type() {
      return isOBFrequency6Code(this.frequency.Type);
    },
    is_legacy_type() {
      return isFrequency1(this.frequency.Type) && !this.is_structured_type;
    },
    legacy_type: {
      get() {
        return hasValue(this.frequency.Type) && !this.is_structured_type ? this.frequency.Type : '';
      },
      set(value) {
        if (hasValue(value)) {
          this.commitFrequency({ Type: value });
          return;
        }
        if (this.is_legacy_type || (hasValue(this.frequency.Type) && !this.is_structured_type)) {
          this.commitFrequency({ Type: '' });
        }
      },
    },
    structured_type: {
      get() {
        return this.is_structured_type ? this.frequency.Type : null;
      },
      set(value) {
        if (hasValue(value)) {
          const frequency = this.is_structured_type ? { ...this.frequency, Type: value } : { Type: value };
          this.commitFrequency(frequency);
          return;
        }
        this.commitFrequency({ Type: '' });
      },
    },
    point_in_time: {
      get() {
        return this.is_structured_type ? this.frequency.PointInTime || '' : '';
      },
      set(value) {
        const frequency = { Type: this.structured_type || defaultStructuredType };
        if (hasValue(value)) {
          frequency.PointInTime = value;
        }
        this.commitFrequency(frequency);
      },
    },
    count_per_period: {
      get() {
        if (!this.is_structured_type || this.frequency.CountPerPeriod === undefined) {
          return '';
        }
        return this.frequency.CountPerPeriod;
      },
      set(value) {
        const frequency = { Type: this.structured_type || defaultStructuredType };
        if (hasValue(value)) {
          frequency.CountPerPeriod = value;
        }
        this.commitFrequency(frequency);
      },
    },
    legacy_type_state() {
      if (!hasValue(this.legacy_type)) {
        return null;
      }
      return isFrequency1(this.legacy_type) && !this.legacy_type_has_structured_fields;
    },
    structured_type_state() {
      if (!hasValue(this.structured_type)) {
        return null;
      }
      return isOBFrequency6Code(this.structured_type);
    },
    has_point_in_time() {
      return hasValue(this.frequency.PointInTime);
    },
    has_count_per_period() {
      return hasValue(this.frequency.CountPerPeriod);
    },
    legacy_type_has_structured_fields() {
      return this.is_legacy_type && (this.has_point_in_time || this.has_count_per_period);
    },
    mutually_exclusive() {
      return !(this.is_structured_type && this.has_point_in_time && this.has_count_per_period);
    },
    point_in_time_state() {
      if (!this.is_structured_type || !this.has_point_in_time) {
        return null;
      }
      return isPointInTime(this.frequency.PointInTime) && this.mutually_exclusive;
    },
    count_per_period_state() {
      if (!this.is_structured_type || !this.has_count_per_period) {
        return null;
      }
      return isCountPerPeriod(this.frequency.CountPerPeriod) && this.mutually_exclusive;
    },
  },
  methods: {
    commitFrequency(frequency) {
      this.$store.commit('config/SET_V4_STANDING_ORDER_FREQUENCY', frequency);
    },
  },
};
</script>

<style scoped>
</style>
