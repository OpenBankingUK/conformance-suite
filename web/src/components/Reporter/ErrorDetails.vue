<template>
  <a-list
    :dataSource="errors"
    itemLayout="horizontal"
  >
    <a-list-item
      slot="renderItem"
      slot-scope="item, index">
      <a-list-item-meta>
        <h3 slot="title">{{ item.message }}</h3>
        <strong slot="description">Path: /{{ item.path.join('/') }}</strong>
        <a-tag
          slot="avatar"
          :color="type === 'errors' ? 'red' : 'orange'">
          {{ item.code }}
        </a-tag>
      </a-list-item-meta>
    </a-list-item>
  </a-list>
</template>

<script>
export default {
  props: {
    errors: {
      type: Array,
      default: () => null,
    },
    type: {
      type: String,
      default: null,
    },
  },
};
</script>
