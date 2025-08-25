<template>
  <div class="form" ref="containerRef">
    <n-form :model="form" label-placement="left" size="small">
      <n-form-item
        label="批量导出数量"
        path="exportCount"

      >
        <div>
          <n-input-number
            v-model:value="form.exportCount"
            :min="1"
            :max="100"
            placeholder="请输入数量"
            style="width: 200px"
          />
          <n-text
            depth="3"
            style="font-size: 12px; display: block; margin-top: 4px"
          >
            一次性导出多个密钥时使用。建议根据实际需要调整，过大可能导致界面卡顿，扫码识别失败。
          </n-text>
        </div>
      </n-form-item>


    </n-form>

    <n-alert v-if="isModified" type="warning" class="modified-floating" closable>
      当前设置有未保存的更改
    </n-alert>
      <n-space justify="center" style="margin-top: 16px">
        <n-button type="primary" @click="handleSubmit" :disabled="!isModified"
          >保存</n-button
        >
        <n-button @click="handleReset" :disabled="!isModified">重置</n-button>
      </n-space>    
  </div>
</template>

<script setup lang="ts">
import { onMounted, ref, computed } from "vue";
import { useMessage } from "naive-ui";
import { getConfig } from "./init";

const message = useMessage();
const containerRef = ref<HTMLElement | null>(null);
const form = ref({
  exportCount: 10,
});
const oldForm = ref({ ...form.value });

// 判断是否有修改
const isModified = computed(() => {
  return JSON.stringify(form.value) !== JSON.stringify(oldForm.value);
});

const handleSubmit = () => {
  localStorage.setItem("Authenticator", JSON.stringify(form.value));
  message.success(`已保存`);
  oldForm.value = { ...form.value }; // 更新基准
};

const handleReset = () => {
  const configData = getConfig("Authenticator");
  form.value = configData;
  oldForm.value = { ...configData }; // 同步重置基准
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
