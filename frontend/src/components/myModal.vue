<template>
  <n-modal
    v-model:show="localShow"
    preset="card"
    :title="props.title"
    :mask-closable="true"
    style="
      max-width: 450px;
      border-radius: 12px;
      box-shadow: 0 8px 24px rgba(0, 0, 0, 0.1);
      margin: 20px auto;
      user-select: none;
    "
    @close="closeModal"
  >
    <slot></slot>
  </n-modal>
</template>

<script setup lang="ts">
import { ref, watch } from "vue";

interface Props {
  show: boolean;
  title: string;
}

const props = defineProps<Props>();
const emit = defineEmits(["update:show"]);

const localShow = ref(props.show);

// 监听 props.show 变化，同步给 localShow
watch(() => props.show, (val) => {
  localShow.value = val;
});

// 监听 localShow 变化，关闭时通知父组件
watch(localShow, (val) => {
  if (!val) {
    emit("update:show", false);
  }
});

const closeModal = () => {
  localShow.value = false;
};
</script>
