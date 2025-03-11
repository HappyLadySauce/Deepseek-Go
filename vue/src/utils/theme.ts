import { ref } from 'vue'

// 定义主题类型
export type ThemeMode = 'light' | 'dark'

// 创建响应式主题状态
export const currentTheme = ref<ThemeMode>('dark')

// 主题配置
export const themeConfig = {
  dark: {
    backgroundColor: '#0d1014',
    textColor: '#ffffff',
    primaryColor: '#3f85ed',
    secondaryColor: '#a2c5f9',
    asideBackground: '#14181e',
    headerBackground: '#14181e',
    borderColor: '#303c4b',
    hoverColor: 'rgba(63, 133, 237, 0.1)',
    menuTextColor: '#c7dcfb',
    menuActiveTextColor: '#a2c5f9',
    menuActiveBackground: 'rgba(63, 133, 237, 0.2)',
    shadowColor: 'rgba(0, 0, 0, 0.3)',
    accentColor: '#5f7ca5'
  },
  light: {
    backgroundColor: '#f0f9ff',
    textColor: '#333333',
    primaryColor: '#1890ff',
    secondaryColor: '#52c41a',
    asideBackground: '#ffffff',
    headerBackground: '#ffffff',
    borderColor: '#e6e6e6',
    hoverColor: '#f5f7fa',
    menuTextColor: '#333333',
    menuActiveTextColor: '#1890ff',
    menuActiveBackground: 'rgba(24, 144, 255, 0.1)',
    shadowColor: 'rgba(0, 0, 0, 0.1)',
    accentColor: '#61dafb'
  }
}

// 切换主题
export const toggleTheme = () => {
  currentTheme.value = currentTheme.value === 'dark' ? 'light' : 'dark'
  localStorage.setItem('theme', currentTheme.value)
  applyTheme()
}

// 应用主题
export const applyTheme = () => {
  // 从 localStorage 获取已保存的主题设置
  const savedTheme = localStorage.getItem('theme') as ThemeMode
  if (savedTheme && (savedTheme === 'light' || savedTheme === 'dark')) {
    currentTheme.value = savedTheme
  }
  
  const theme = themeConfig[currentTheme.value]
  
  // 更新 CSS 变量
  document.documentElement.style.setProperty('--bg-color', theme.backgroundColor)
  document.documentElement.style.setProperty('--text-color', theme.textColor)
  document.documentElement.style.setProperty('--primary-color', theme.primaryColor)
  document.documentElement.style.setProperty('--secondary-color', theme.secondaryColor)
  document.documentElement.style.setProperty('--aside-bg', theme.asideBackground)
  document.documentElement.style.setProperty('--header-bg', theme.headerBackground)
  document.documentElement.style.setProperty('--border-color', theme.borderColor)
  document.documentElement.style.setProperty('--hover-color', theme.hoverColor)
  document.documentElement.style.setProperty('--menu-text-color', theme.menuTextColor)
  document.documentElement.style.setProperty('--menu-active-text-color', theme.menuActiveTextColor)
  document.documentElement.style.setProperty('--menu-active-bg', theme.menuActiveBackground)
  document.documentElement.style.setProperty('--shadow-color', theme.shadowColor)
  document.documentElement.style.setProperty('--accent-color', theme.accentColor)
}

// 初始化主题
applyTheme() 