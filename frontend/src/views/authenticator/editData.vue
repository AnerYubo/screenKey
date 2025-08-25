<!-- src/views/authenticator/editData.vue -->
<template>
  <n-form
    ref="formRef"
    :model="form"
    :rules="rules"
    :label-width="100"
    style="text-align: start"
  >
    <n-form-item path="issuer" label="颁发者">
      <n-input v-model:value="form.issuer" />
    </n-form-item>
    <n-form-item path="account" label="账户">
      <n-input v-model:value="form.account" />
    </n-form-item>
    <n-form-item path="category" label="分类">
      <n-select
        v-model:value="form.category"
        filterable
        tag
        :options="categoriesOptions"
      />
    </n-form-item>    
    <n-form-item path="remark" label="备注">
      <n-input v-model:value="form.remark" />
    </n-form-item>

    <n-space justify="end" style="margin-top: 16px">
      <n-button quaternary @click="emit('close')" :disabled="loading">取消</n-button>
      <n-button type="primary" @click="handleConfirm" :loading="loading">保存</n-button>
    </n-space>
  </n-form>
</template>

<script setup lang="ts">
import { ref, reactive, watch, computed } from "vue";
import { useMessage } from "naive-ui";
import type { TotpEntry } from "./type";
import { UpdateTOTP } from "../../../wailsjs/go/app/TOTPHandler";

const props = defineProps<{ data: TotpEntry | null, categoriesList: any }>();
const emit = defineEmits(["success", "close"]);

const message = useMessage();
const formRef = ref();
const loading = ref(false);
const categoriesOptions = computed(() =>
  props.categoriesList.map((k:any) => ({ label: k, value: k }))
)
// 初始化表单数据（用 reactive 创建空模型）
const form = reactive<TotpEntry>({
  id: 0,
  issuer: "",
  account: "",
  remark: "",
  secret: "",
  otpauth: "",
  currentCode: "",
  category: '未分类',
  remainSeconds: 0,
});

// props.data 变化时同步
watch(
  () => props.data,
  (val) => {
    if (val) {
      Object.assign(form, val);
    }
  },
  { immediate: true }
);

const rules = {
  issuer: { required: true, message: "请输入颁发者", trigger: "blur" },
  account: { required: true, message: "请输入账户", trigger: "blur" },
};

async function handleConfirm() {
  if (loading.value) return;
  try {
    await formRef.value?.validate();

    loading.value = true;

    const payload: TotpEntry = {
      ...form,
      issuer: form.issuer.trim(),
      account: form.account.trim(),
      remark: form.remark.trim(),
      category: form.category.trim(),
    };

    const err = await UpdateTOTP(form.id, payload);
    if (err != null) throw new Error(err);

    emit("success", { ...payload });
    message.success("✅ 修改成功");
    emit("close");
  } catch (e) {
    message.error("❌ 数据不符合规范或保存失败");
  } finally {
    loading.value = false;
  }
}
</script>

<style scoped>
/* 可选美化 */
</style>
