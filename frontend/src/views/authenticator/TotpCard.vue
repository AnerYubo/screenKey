<template>
  <n-card
    size="small"
    class="totp-card"
    hoverable
    @click="copyCode"
    @contextmenu.prevent="onRightClick"
    :class="{ selected: selected }"
  >
    <div class="card-header">
      <div class="issuer">{{ data.issuer }}: {{ data.account }}</div>
      <div class="category">{{ data.category }}</div>
    </div>

    <div class="card-body">
      <span class="code" :class="{ 'flash-red': data.remainSeconds <= 5 }">
        {{ data.currentCode }}
      </span>

      <!-- 使用 CountdownCircle 组件 -->
      <CountdownCircle
        :remain-seconds="data.remainSeconds"
        :show-number="false"
      />
    </div>

    <div class="card-footer">
      <span>{{ data.remark }}</span>
    </div>
  </n-card>
</template>

<script setup lang="ts">
import { computed } from "vue";
import { NCard, useMessage } from "naive-ui";
import type { TotpEntry } from "./type";
import CountdownCircle from "@/components/CountdownCircle.vue";
const props = defineProps<{
  data: TotpEntry;
  selected?: boolean;
}>();

const emit = defineEmits<{
  (e: "update-selection", id: number): void;
}>();

const message = useMessage();

function copyCode() {
  navigator.clipboard.writeText(props.data.currentCode || "");
  message.success("验证码已复制", { duration: 1000 });
}

function onRightClick(e: MouseEvent) {
  emit("update-selection", props.data.id);
}
</script>

<style scoped lang="scss">
.totp-card {
  border-radius: 14px;
  padding: 10px;
  transition: all 0.2s ease-in-out;
  cursor: pointer;

  &:hover {
    transform: translateY(-2px);
    box-shadow:var(--box-shadow);
  }
  &.selected {
    box-shadow: 0 4px 12px rgba(24, 160, 88, 0.5);
    transform: translateY(-2px);
  }
  .card-header {
    display: flex;
    justify-content: space-between;
    margin-bottom: 6px;
    .issuer {
      font-weight: bold;
      font-size: 14px;
      color: var(--color-text);
      /* ✅ 文字过长省略号 */
      overflow: hidden;
      text-overflow: ellipsis;
      white-space: nowrap;
      max-width: 70%;
    }
    .category {
      font-size: 12px;
      background: #f0f0f0;
      padding: 2px 8px;
      border-radius: 8px;
      color: #666;
      /* ✅ 文字过长省略号 */
      overflow: hidden;
      text-overflow: ellipsis;
      white-space: nowrap;
      max-width: 20%;
    }
  }

  .card-footer {
    max-width: 90%;
    text-align: start;
    font-size: 12px;
    /* ✅ 文字过长省略号 */
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
  }

  .card-body {
    display: flex;
    align-items: center;
    margin-top: 8px;
    gap: 10px;

    .code {
      font-size: 18px;
      font-weight: bold;
      color: #18a058;
      letter-spacing: 2px;
    }

    .countdown {
      width: 28px;
      height: 28px;
    }

    .bg {
      fill: #eee;
    }
    .fg {
      transition: d 1s linear;
    }
  }
}

@keyframes flash-red {
  25%,
  100% {
    color: red;
    opacity: 1;
  }
  0%,
  50% {
    color: #ff6666;
    opacity: 0.5;
  }
}
.flash-red {
  animation: flash-red 1.5s ease-in-out infinite;
}
</style>
