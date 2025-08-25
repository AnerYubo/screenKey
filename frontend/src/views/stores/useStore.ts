import { defineStore } from "pinia";


// 创建 store
export const useStore = defineStore("store", {
  state: () => ({
    theme: "dark", //主题 light dark auto
  }),
  actions: {
    setTheme(theme: "light" | "dark" | "auto") {
      this.theme = theme;
    },
  },
});
