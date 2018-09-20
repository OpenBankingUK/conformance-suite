<template>
  <a-layout class="h100">
    <a-layout>
      <a-layout-content class="clean">
        <slot />
      </a-layout-content>
      <a-layout-footer class="centered">
        Functional Conformance Suite © Open Banking 2018
      </a-layout-footer>
    </a-layout>
  </a-layout>
</template>

<style>
.clean {
  display: flex;
  justify-content: center;
  align-items: center;
  padding: 15px;
}
</style>
