<template>
  <div class="my-menu">
    <div
      class="menu-item"
      :class="{ active: isActive('/PortKnocking') }"
      @click="$router.push('/PortKnocking')"
    >
      端口敲门
    </div>
    <div
      class="menu-item"
      :class="{ active: isActive('/') }"
      @click="$router.push('/')"
    >
      Authenticator
    </div>
    <div
      style="display: flex; align-items: center; justify-content: end; flex: 1"
    >
      <n-popover trigger="hover" :delay="200" placement="bottom">
        <template #trigger>
          <n-icon
            size="20"
            style="cursor: pointer;"
            :color=" isTop ? '#18a058' : ''"
            @click="toggleTop"
          >
            <svg
              xmlns="http://www.w3.org/2000/svg"
              viewBox="0 0 16 16"
            >
              <g fill="none">
                <path
                  d="M10.059 2.445a1.5 1.5 0 0 0-2.386.354l-2.02 3.79l-2.811.937a.5.5 0 0 0-.196.828L4.793 10.5l-2.647 2.647L2 14l.853-.146L5.5 11.207l2.146 2.147a.5.5 0 0 0 .828-.196l.937-2.81l3.779-2.024a1.5 1.5 0 0 0 .354-2.38L10.06 2.444z"
                  fill="currentColor"
                ></path>
              </g>
            </svg>
          </n-icon>
        </template>
        <span>{{ isTop ? '取消置顶' : '置顶' }}</span>
      </n-popover>

      <div
        class="menu-item "
        :class="{ active: isActive('/Setup') }"
        @click="$router.push('/Setup')"
      >
      
        设置
        <n-icon><EllipsisVertical /></n-icon>
      </div>
    </div>
  </div>
</template>

<script lang="ts" setup>
import { ref } from "vue";
import { useRoute } from "vue-router";
// wails runtime API
import { WindowSetAlwaysOnTop } from "../../wailsjs/runtime/runtime";
import { EllipsisVertical } from "@vicons/ionicons5";

const route = useRoute();
const isTop = ref(false); // 是否置顶

function isActive(path: string) {
  return route.path === path;
}

async function toggleTop() {
  isTop.value = !isTop.value;
  await WindowSetAlwaysOnTop(isTop.value);
}
</script>

<style lang="scss" scoped>
.my-menu {
  width: 100%;
  display: flex;
  padding: 10px;
  user-select: none;
  background-color: var(--color-bg-main);
  border-bottom: 1px #c7c7c7 solid;

  .menu-item {
    cursor: pointer;
    font-weight: 500;
    box-sizing: border-box;
    padding: 2px 10px;
    transition: all 0.3s ease;
    display: flex;
    align-items: center;
    &:hover {
      --color: rgb(54, 158, 100);
      color: var(--color);
    }

    &.active {
      --color: rgb(54, 158, 100);
      color: var(--color);
      border-radius: 2px;
      font-weight: bolder;
    }
  }
}
</style>
