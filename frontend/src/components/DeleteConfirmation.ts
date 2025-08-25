// src/hooks/useDeleteConfirm.ts
import { h } from "vue";

export interface DeleteConfirmOptions {
  count: number;
  dialog: any;
  message: any;
  onConfirm: () => Promise<void>;
  title?: string;
  content?: string;
}

/**
 * 显示删除确认对话框
 * @param options 配置项
 * @returns {Promise<boolean>} 用户是否点击了确认
 */
export function useDeleteConfirm(
  options: DeleteConfirmOptions
): Promise<boolean> {
  return new Promise((resolve) => {
    options.dialog.warning({
      title: () =>
        h("div", { style: "font-weight: bold; font-size: 16px;" }, [
          options.title || "确认删除 ",
          h(
            "span",
            { style: "color: red; font-weight: bold;" },
            String(options.count)
          ),
          " 条记录？",
        ]),
      content: () =>
        h("div", { style: "color: #555;" }, [
          h("p", null, options.content || "这些记录删除后将无法恢复，"),
          h("p", null, "请确认你是否要执行此操作。"),
        ]),
      positiveText: "删除",
      negativeText: "取消",
      async onPositiveClick() {
        try {
          await options.onConfirm();
          options.message.success(`✅ 成功删除 ${options.count} 条记录`);
          resolve(true);
        } catch (err) {
          options.message.error(
            `❌ 删除失败：${err instanceof Error ? err.message : String(err)}`
          );
          resolve(false);
        }
      },
      onNegativeClick() {
        resolve(false);
      },
    });
  });
}
