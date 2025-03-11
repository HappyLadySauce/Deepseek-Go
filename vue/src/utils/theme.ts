import { ref } from 'vue'

// 当前主题，默认跟随系统
const prefersDark = window.matchMedia('(prefers-color-scheme: dark)').matches;
export const currentTheme = ref(prefersDark ? 'dark' : 'light');

// CSS变量映射
const themes = {
  light: {
    '--page-bg': '#f5f7fa',
    '--card-bg': '#ffffff',
    '--header-bg': '#ffffff',
    '--text-color': '#303133',
    '--text-light': '#606266',
    '--text-muted': '#909399',
    '--menu-text-color': '#303133',
    '--border-color': '#EBEEF5',
    '--hover-color': '#f5f7fa',
    '--active-color': '#ecf5ff',
    '--primary-color': '#409EFF',
    '--secondary-color': '#67C23A',
    '--message-bg': '#f5f7fa',
    '--user-bg': '#ecf5ff',
    '--assistant-bg': '#f5f7fa',
    '--card-shadow': '0 2px 12px 0 rgba(0, 0, 0, 0.1)',
    '--message-shadow': '0 2px 4px rgba(0, 0, 0, 0.05)',
    '--tag-bg': '#ecf5ff',
    '--item-bg': '#ffffff',
    '--scrollbar-color': '#C0C4CC',
    '--scrollbar-track': '#EBEEF5',
    '--heading-color': '#303133'
  },
  dark: {
    '--page-bg': '#1e1e1e',
    '--card-bg': '#252526',
    '--header-bg': '#1e1e1e',
    '--text-color': '#e1e1e1',
    '--text-light': '#b0b0b0',
    '--text-muted': '#8e8e8e',
    '--menu-text-color': '#e1e1e1',
    '--border-color': '#3e3e3e',
    '--hover-color': '#2a2d2e',
    '--active-color': '#094771',
    '--primary-color': '#409EFF',
    '--secondary-color': '#67C23A',
    '--message-bg': '#2d2d2d',
    '--user-bg': '#094771',
    '--assistant-bg': '#2d2d2d',
    '--card-shadow': '0 2px 12px 0 rgba(0, 0, 0, 0.4)',
    '--message-shadow': '0 2px 4px rgba(0, 0, 0, 0.2)',
    '--tag-bg': '#094771',
    '--item-bg': '#1e1e1e',
    '--scrollbar-color': '#4e4e4e',
    '--scrollbar-track': '#2d2d2d',
    '--heading-color': '#e1e1e1'
  }
};

// 应用主题
export function applyTheme(theme: 'light' | 'dark') {
  console.log('正在应用主题：', theme);
  
  // 更新类名，方便CSS选择器
  if (theme === 'dark') {
    document.documentElement.classList.add('dark-theme');
    document.documentElement.classList.remove('light-theme');
    document.body.classList.add('dark-theme');
    document.body.classList.remove('light-theme');
  } else {
    document.documentElement.classList.add('light-theme');
    document.documentElement.classList.remove('dark-theme');
    document.body.classList.add('light-theme');
    document.body.classList.remove('dark-theme');
  }
  
  // 存储用户主题偏好
  localStorage.setItem('theme', theme);
  console.log('主题已应用，当前主题：', document.body.classList.contains('dark-theme') ? 'dark' : 'light');
}

// 切换主题
export function toggleTheme() {
  console.log('切换主题，当前主题：', currentTheme.value);
  currentTheme.value = currentTheme.value === 'light' ? 'dark' : 'light';
  applyTheme(currentTheme.value === 'dark' ? 'dark' : 'light');
  console.log('主题已切换为：', currentTheme.value);
}

// 初始化主题
export function initTheme() {
  console.log('初始化主题...');
  // 首先检查是否有存储的主题偏好
  const savedTheme = localStorage.getItem('theme');
  
  if (savedTheme === 'dark' || savedTheme === 'light') {
    console.log('从localStorage中恢复主题：', savedTheme);
    currentTheme.value = savedTheme;
  } else {
    console.log('使用系统默认主题：', prefersDark ? 'dark' : 'light');
  }
  
  applyTheme(currentTheme.value === 'dark' ? 'dark' : 'light');
  
  // 监听系统主题变化
  window.matchMedia('(prefers-color-scheme: dark)').addEventListener('change', e => {
    console.log('系统主题变化：', e.matches ? 'dark' : 'light');
    if (!localStorage.getItem('theme')) {
      currentTheme.value = e.matches ? 'dark' : 'light';
      applyTheme(currentTheme.value === 'dark' ? 'dark' : 'light');
    } else {
      console.log('用户已设置主题偏好，忽略系统主题变化');
    }
  });
}

// 立即初始化主题
console.log('theme.ts 被加载');
initTheme(); 