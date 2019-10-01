<template>
  <b-form-group
    :id="group_id"
    :label="label"
    :label-for="id"
    :description="group_description">
    <b-form-input
      :id="id"
      :state="valid"
      v-model="model"
      type="text"
      required />
  </b-form-group>
</template>

<script>
import moment from 'moment';

const isISO8601 = (value) => {
  const date = moment(value, moment.ISO_8601);
  return date.isValid();
};

export default {
  name: 'DateTimeISO8601',
  props: {
    /**
     * ID to give to the HTML input element.
     */
    id: {
      type: String,
      required: true,
    },
    /**
     * The name of the key in `this.$store.state.config.configuration`, e.g.,
     * `this.$store.state.config.configuration['first_payment_date_time']`.
     */
    field: {
      type: String,
      required: true,
    },
    label: {
      type: String,
      required: true,
    },
    description: {
      type: String,
      required: true,
    },
    /**
     * The suffix of the mutation type.
     * The final type used will be, e.g., `config/SET_FIRST_PAYMENT_DATE_TIME`.
     */
    mutation: {
      type: String,
      required: true,
    },
  },
  data() {
    return {
    };
  },
  computed: {
    group_id() {
      return `${this.id}_group`;
    },
    group_description() {
      return `${this.description} formatted as ISO 8601 date (eg. 2006-01-02T15:04:05-07:00)`;
    },
    model: {
      get() {
        return this.$store.state.config.configuration[this.field];
      },
      set(value) {
        const type = `config/${this.mutation}`;
        this.$store.commit(type, value);
      },
    },
    valid() {
      return isISO8601(this.model);
    },
  },
  methods: {
  },
};
</script>

<style scoped>
</style>
