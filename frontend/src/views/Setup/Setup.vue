<template>
  <div class="setup-container">
    <n-tabs
      key="lineleft"
      addable
      type="line"
      animated
      placement="left"
    >
      <n-tab-pane name="General" tab="通用"> <General /> </n-tab-pane>
      <n-tab-pane name="Authenticator" tab="Authenticator"> <Authenticator /> </n-tab-pane>
      <n-tab-pane name="About" tab="关于"> <About /> </n-tab-pane>
    </n-tabs>
  </div>
</template>

<script setup lang="ts">
import { ref } from "vue";
import About from "../About/About.vue";
import Authenticator from "./Authenticator.vue";
import General from "./General.vue";
</script>

<style scoped lang="scss">
.setup-container {
  height: calc(100vh - 50px); // 整个容器撑满
  display: flex;
  flex-direction: column;
  overflow: hidden;
  *{
    box-sizing: border-box;
  }
}

:deep(.n-tabs) {
  flex: 1;
  display: flex;
  flex-direction: row;
  overflow: hidden; // 关键：父容器不能滚动
}

:deep(.n-tabs-nav.n-tabs-nav--left) {
  height: 100%;
  overflow-y: auto; // 左侧独立滚动
  flex-shrink: 0;
  border-right: 1px solid #eee;
}

:deep(.n-tabs-content) {
  flex: 1;
  display: flex;
  flex-direction: column;
  overflow: hidden;
}
:deep(.n-tabs.n-tabs--left .n-tab-pane ){
    padding: 0;
}
:deep(.n-tab-pane) {
  flex: 1;
  overflow-y: auto;
  padding: 0;
  background: var(--color-bg-main);
  box-sizing: border-box;
  height: 100%  ;
  /* 滚动条样式仿 Naive */
  scrollbar-width: thin;
  scrollbar-color: rgba(0,0,0,0.25) transparent;

  &::-webkit-scrollbar {
    width: 6px;
    height: 6px;
  }

  &::-webkit-scrollbar-thumb {
    background: rgba(0,0,0,0.25);
    border-radius: 3px;
    transition: background 0.3s;
  }

  &::-webkit-scrollbar-thumb:hover {
    background: rgba(0,0,0,0.4);
  }

  &::-webkit-scrollbar-track {
    background: transparent;
  }
}

</style>
