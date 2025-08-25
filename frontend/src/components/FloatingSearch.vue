<template>
  <div class="header-wrapper">
    <!-- 如果搜索框可见，替换 header -->
    <template v-if="visible">
      <transition name="fade">
        <div class="search-header">
          <n-input
            ref="searchInput"
            v-model:value="query"
            size="medium"
            placeholder="输入关键字搜索 (ESC 关闭)"
            clearable
            autofocus
            @keydown.esc="closeSearch"
            @input="emitSearch"
          >
            <template #prefix>
              <n-icon><SearchOutline /></n-icon>
            </template>
          </n-input>
        </div>
      </transition>
    </template>
    <!-- 默认 header -->
    <div v-else class="default-header">
      <slot> </slot>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, nextTick, onMounted, onBeforeUnmount } from "vue";
import { NInput, NIcon } from "naive-ui";
import { SearchOutline } from "@vicons/ionicons5";

const visible = ref(false);
const query = ref("");
const searchInput = ref();

const emit = defineEmits<{
  (e: "update:query", value: string): void;
}>();

function openSearch() {
  visible.value = true;
  nextTick(() => {
    searchInput.value?.focus();
  });
}

function closeSearch() {
  visible.value = false;
  query.value = "";
  emit("update:query", "");
}

function emitSearch() {
  emit("update:query", query.value);
}

function handleKeydown(e: KeyboardEvent) {
  // Ctrl+K / Cmd+K / /
  if ((e.ctrlKey || e.metaKey) && e.key.toLowerCase() === "k") {
    e.preventDefault();
    openSearch();
  } else if (e.key === "/") {
    e.preventDefault();
    openSearch();
  }
}

onMounted(() => window.addEventListener("keydown", handleKeydown));
onBeforeUnmount(() => window.removeEventListener("keydown", handleKeydown));
</script>

<style scoped lang="scss">
.header-wrapper {
  display: flex;
  align-items: center;
  width: 100%;
}

.search-header {
  flex: 1;
}

.default-header {
  flex: 1;
}

.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.25s ease;
}
.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}
</style>
