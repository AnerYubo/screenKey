import { useStore } from "../stores/useStore";
let mediaQuery: MediaQueryList | null = null;
let mediaListener: ((e: MediaQueryListEvent) => void) | null = null;
//初始化本地配置
const config: any = {
    Authenticator: {
        exportCount: 10, //每次导出数量
    },
    General:{
        theme:"light", //主题 light dark auto
    },
}

const init = () =>{
    var localData = localStorage.getItem("Authenticator");
    if(localData == null || localData == undefined){
        localStorage.setItem("Authenticator",JSON.stringify(config.Authenticator));
    }
    var localData = localStorage.getItem("General");
    if(localData == null || localData == undefined){
        localStorage.setItem("General",JSON.stringify(config.General));
    }    

}

const getConfig = (name :string) =>{
    const localData = localStorage.getItem(name);
    if(localData){
        return JSON.parse(localData);
    }else{
        return config[name];
    }
}

// 设置主题
const setTheme = (theme: "light" | "dark" | "auto") => {
  const store = useStore();

  let finalTheme = theme;

  // 先移除旧的监听（避免重复）
  if (mediaQuery && mediaListener) {
    mediaQuery.removeEventListener("change", mediaListener);
    mediaQuery = null;
    mediaListener = null;
  }

  if (theme === "auto") {
    mediaQuery = window.matchMedia("(prefers-color-scheme: dark)"); // ✅ 保存对象
    finalTheme = mediaQuery.matches ? "dark" : "light";

    // 监听系统主题变化
    mediaListener = (e: MediaQueryListEvent) => {
      const newTheme = e.matches ? "dark" : "light";
      applyTheme(newTheme);
    };

    // ✅ 这里 TS 就不会提示 mediaQuery 可能为 null 了
    mediaQuery.addEventListener("change", mediaListener);
  }

  // 设置根元素属性（方便 CSS 变量生效）
  document.documentElement.setAttribute("data-theme", finalTheme);

  // 返回给 naive-ui 使用
  store.theme = finalTheme;
  return finalTheme;
};


// 应用主题到 DOM & 返回 naive-ui 需要的 theme
const applyTheme = (theme: "light" | "dark") => {
const store = useStore();

  document.documentElement.setAttribute("data-theme", theme);
store.theme = theme;
};
export { init, getConfig, setTheme, config };