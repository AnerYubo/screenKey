<template>
  <n-form
    ref="formRef"
    :model="form"
    :rules="rules"
    :label-width="100"
    style="text-align: start"
  >
    <n-form-item path="issuer" label="颁发者">
      <n-input v-model:value="form.issuer" placeholder="例如：GitHub" />
    </n-form-item>
    <n-form-item path="account" label="账户">
      <n-input
        v-model:value="form.account"
        placeholder="例如：example@mail.com"
      />
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
      <n-input v-model:value="form.remark" placeholder="例如：工作登录" />
    </n-form-item>
  </n-form>

  <n-space justify="end" style="margin-top: 16px">
    <n-button quaternary @click="emit('close')" :disabled="loading"
      >取消</n-button
    >
    <n-button type="primary" @click="handleSubmit" :loading="loading"
      >生成</n-button
    >
  </n-space>
</template>

<script setup lang="ts">
import { ref, reactive, computed } from "vue";
import { useMessage } from "naive-ui";
import { GenerateTOTP } from "../../../wailsjs/go/app/TOTPHandler";
const props = defineProps<{ categoriesList: any }>();

const emit = defineEmits(["success", "close"]);
const message = useMessage();

const formRef = ref();
const loading = ref(false);

const categoriesOptions = computed(() =>
  props.categoriesList.map((k:any) => ({ label: k, value: k }))
)
const form = reactive({
  issuer: "",
  account: "",
  remark: "",
  category: '未分类',
});

const rules = {
  issuer: { required: true, message: "请输入颁发者", trigger: "blur" },
  account: { required: true, message: "请输入账户", trigger: "blur" },
};

function resetForm() {
  form.issuer = "";
  form.account = "";
  form.remark = "";
}

const handleSubmit = () => {
  if (loading.value) return;

  formRef.value?.validate(async (errors: any) => {
    if (!errors) {
      loading.value = true;
      try {
        const res = await GenerateTOTP({
          issuer: form.issuer.trim(),
          account: form.account.trim(),
          remark: form.remark.trim(),
          category: form.remark.trim(),
        });
        emit("success", res);
        message.success("TOTP 生成成功");
        resetForm();
        emit("close");
      } catch (err) {
        message.error(
          "❌ 生成失败：" + (err instanceof Error ? err.message : "未知错误")
        );
      } finally {
        loading.value = false;
      }
    } else {
      message.warning("请填写完整的表单信息");
    }
  });
};
</script>
