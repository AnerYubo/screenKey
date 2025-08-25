<template>
  <n-config-provider :locale="zhCN" :date-locale="dateZhCN" :theme="store.theme == 'dark' ? darkTheme: lightTheme">
    <n-dialog-provider>
      <n-message-provider>
        <n-config-provider :hljs="hljs">
          <div class="layout-wrapper">
            <my-menu />
            <n-scrollbar class="layout-content" trigger="hover">
              <router-view v-slot="{ Component, route }">
                <keep-alive :exclude="'About'">
                  <component :is="Component" :key="route.name" />
                </keep-alive>
              </router-view>
            </n-scrollbar>
          </div>
        </n-config-provider>
      </n-message-provider>
    </n-dialog-provider>
  </n-config-provider>
</template>

<script lang="ts" setup>
import myMenu from "./views/myMenu.vue";
import { zhCN, dateZhCN } from "naive-ui";
import hljs from "highlight.js/lib/core";
import javascript from "highlight.js/lib/languages/javascript";
import { init, setTheme, getConfig } from "./views/Setup/init";
import { onMounted } from "vue";
import { darkTheme, lightTheme } from 'naive-ui';
import { useStore } from "./views/stores/useStore";
const store = useStore();
hljs.registerLanguage("javascript", javascript);
// 阻止默认拖拽打开行为（尤其是图片、PDF 被系统打开）
window.addEventListener("dragover", (e) => {
  e.preventDefault();
});
window.addEventListener("drop", (e) => {
  e.preventDefault();
});
onMounted(() => {
  init();
  const General =  getConfig('General')
  setTheme(General.theme)
}); 
</script>

<style lang="scss" scoped>
.layout-wrapper {
  display: flex;
  flex-direction: column;
  height: 100vh;
  user-select: none;
  .layout-content {
    flex: 1;
    padding: 16px;
    overflow: hidden;

    :deep(.n-scrollbar-content) {
      min-height: 100%;
    }
  }
}
</style>
