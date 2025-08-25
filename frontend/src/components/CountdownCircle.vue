<template>
  <svg class="countdown" viewBox="0 0 36 36">
    <!-- 背景圆 -->
    <circle class="bg" cx="18" cy="18" r="16" />

    <!-- 动态扇形 -->
    <path
      class="fg"
      :fill="remainSeconds <= 5 ? '#d03050' : '#18a058'"
      :d="sectorPath"
    />

    <!-- 可选：在圆心显示剩余秒数 -->
    <text
      x="18"
      y="21"
      text-anchor="middle"
      font-size="10"
      fill="#fff"
      v-if="showNumber"
    >
      {{ remainSeconds }}
    </text>
  </svg>
</template>

<script setup lang="ts">
import { computed } from "vue";

const props = defineProps<{
  remainSeconds: number;   // 剩余秒数
  total?: number;          // 总时长，默认 30s
  showNumber?: boolean;    // 是否显示圆心数字
}>();

const total = props.total ?? 30;
const r = 16;
const cx = 18;
const cy = 18;

// ✅ 计算扇形 Path
const sectorPath = computed(() => {
  const remain = props.remainSeconds ?? 0;
  const percent = remain / total;

  if (percent <= 0) return "";  // ⬅️ 完全不画
  if (percent >= 1) {
    // ⬅️ 满圆
    return `
      M ${cx} ${cy}
      m 0 -${r}
      a ${r} ${r} 0 1 1 0 ${2 * r}
      a ${r} ${r} 0 1 1 0 -${2 * r}
      Z
    `;
  }

  const angle = percent * 2 * Math.PI - Math.PI / 2;
  const x = cx + r * Math.cos(angle);
  const y = cy + r * Math.sin(angle);
  const largeArc = percent > 0.5 ? 1 : 0;

  return `
    M ${cx} ${cy}
    L ${cx} ${cy - r}
    A ${r} ${r} 0 ${largeArc} 1 ${x} ${y}
    Z
  `;
});

</script>

<style scoped>
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
</style>
