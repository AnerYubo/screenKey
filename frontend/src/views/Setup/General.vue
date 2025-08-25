<template>
  <div class="form" ref="containerRef">
    <n-form :model="form" label-placement="left" size="small">
      <n-form-item label="主题" path="theme">
        <div>
          <n-space vertical>
            <n-select
              v-model:value="form.theme"
              :options="themeOptions"
              placeholder="选择主题"
              style="width: 200px"
            />
          </n-space>
          <n-text
            depth="3"
            style="font-size: 12px; display: block; margin-top: 4px"
          >
            选择应用的外观主题，跟随系统时会自动切换明暗模式【建议使用亮色/暗色并无美化】。
          </n-text>
        </div>
      </n-form-item>
    </n-form>

    <n-alert
      v-if="isModified"
      type="warning"
      class="modified-floating"
      closable
    >
      当前设置有未保存的更改
    </n-alert>

    <n-space justify="center" style="margin-top: 16px">
      <n-button type="primary" @click="handleSubmit" :disabled="!isModified">
        保存
      </n-button>
      <n-button @click="handleReset" :disabled="!isModified">
        重置
      </n-button>
    </n-space>
  </div>
</template>

<script setup lang="ts">
import { onMounted, ref, computed } from "vue";
import { useMessage } from "naive-ui";
import { getConfig, setTheme } from "./init";

const themeOptions = [
  { label: "亮色", value: "light" },
  { label: "暗色", value: "dark" },
  { label: "跟随系统", value: "auto" },
];

const message = useMessage();
const containerRef = ref<HTMLElement | null>(null);

const form = ref<{ theme: "light" | "dark" | "auto" }>({
  theme: "auto",
});
const oldForm = ref({ ...form.value });

// 判断是否有修改
const isModified = computed(() => {
  return JSON.stringify(form.value) !== JSON.stringify(oldForm.value);
});

const handleSubmit = () => {
  localStorage.setItem("General", JSON.stringify(form.value));
  setTheme(form.value.theme); // ✅ 保存后立即应用主题
  message.success("已保存");
  oldForm.value = { ...form.value }; // 更新基准
};

const handleReset = () => {
  const configData = getConfig("General");
  form.value = configData;
  oldForm.value = { ...configData }; // 同步重置基准
  setTheme(form.value.theme);
};

onMounted(() => {
  handleReset();
});
</script>

<style scoped>
.form {
  position: relative;
  width: 100%;
  min-height: 100%;
  box-sizing: border-box;
  padding: 20px;
  display: flex;
  flex-direction: column;
  justify-content: space-between;
}
.modified-floating {
  position: fixed;
  bottom: 20px;
  right: 20px;
  width: auto;
  max-width: 300px;
  z-index: 999;
  font-size: 13px;
}
</style>
