// Copy of Bootstrap 4.2 Spinner from Bootstrap-Vue master branch:
// https://github.com/bootstrap-vue/bootstrap-vue/tree/dev/src/components/spinner
//
// Can be removed when next release of Bootstrap-Vue is made and we've upgraded.
import { mergeData } from 'vue-functional-data-merge';

// @vue/component
export default {
  name: 'BSpinner',
  functional: true,
  props: {
    type: {
      type: String,
      default: 'border', // SCSS currently supports 'border'
    },
    label: {
      type: String,
      default: null,
    },
    variant: {
      type: String,
      default: null,
    },
    small: {
      type: Boolean,
      default: false,
    },
    role: {
      type: String,
      default: 'status',
    },
    tag: {
      type: String,
      default: 'span',
    },
  },
  render(h, { props, data, slots }) {
    let label = h(false);
    const hasLabel = slots().label || props.label;
    if (hasLabel) {
      label = h('span', { staticClass: 'sr-only' }, hasLabel);
    }
    return h(
      props.tag,
      mergeData(data, {
        attrs: {
          role: hasLabel ? props.role || 'status' : null,
          'aria-hidden': hasLabel ? null : 'true',
        },
        class: {
          [`spinner-${props.type}`]: Boolean(props.type),
          [`spinner-${props.type}-sm`]: props.small,
          [`text-${props.variant}`]: Boolean(props.variant),
        },
      }),
      [label],
    );
  },
};
