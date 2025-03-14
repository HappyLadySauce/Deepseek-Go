<script setup lang="ts">
import { RouterView } from 'vue-router'
import { useUserStore } from './stores/user'
import { storeToRefs } from 'pinia'
import { onMounted, watch } from 'vue'
import { currentTheme } from './utils/theme'

// 用户存储
const userStore = useUserStore()
const { isAuthenticated } = storeToRefs(userStore)

// 检查本地存储中的令牌
onMounted(() => {
  userStore.checkAuth()
})

// 监听主题变化，设置Element Plus的主题
watch(currentTheme, (newTheme) => {
  if (newTheme === 'dark') {
    document.documentElement.setAttribute('class', 'dark')
  } else {
    document.documentElement.removeAttribute('class')
  }
}, { immediate: true })
</script>

<template>
  <div class="app-container">
    <RouterView />
  </div>
</template>

<style>
/* 全局CSS变量(默认亮色主题) */
:root {
  --page-bg: #f5f7fa;
  --card-bg: #ffffff;
  --header-bg: #ffffff;
  --text-color: #303133;
  --text-light: #606266;
  --text-muted: #909399;
  --menu-text-color: #303133;
  --border-color: #EBEEF5;
  --hover-color: #f5f7fa;
  --active-color: #ecf5ff;
  --primary-color: #409EFF;
  --secondary-color: #67C23A;
  --message-bg: #f5f7fa;
  --user-bg: #ecf5ff;
  --assistant-bg: #f5f7fa;
  --card-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.1);
  --message-shadow: 0 2px 4px rgba(0, 0, 0, 0.05);
  --tag-bg: #ecf5ff;
  --item-bg: #ffffff;
  --scrollbar-color: #C0C4CC;
  --scrollbar-track: #EBEEF5;
  --heading-color: #303133;
}

/* 暗色主题变量 */
.dark-theme {
  --page-bg: #1e1e1e;
  --card-bg: #252526;
  --header-bg: #1e1e1e;
  --text-color: #e1e1e1;
  --text-light: #b0b0b0;
  --text-muted: #8e8e8e;
  --menu-text-color: #e1e1e1;
  --border-color: #3e3e3e;
  --hover-color: #2a2d2e;
  --active-color: #094771;
  --primary-color: #409EFF;
  --secondary-color: #67C23A;
  --message-bg: #2d2d2d;
  --user-bg: #094771;
  --assistant-bg: #2d2d2d;
  --card-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.4);
  --message-shadow: 0 2px 4px rgba(0, 0, 0, 0.2);
  --tag-bg: #094771;
  --item-bg: #1e1e1e;
  --scrollbar-color: #4e4e4e;
  --scrollbar-track: #2d2d2d;
  --heading-color: #e1e1e1;
}

/* 全局样式 */
* {
  margin: 0;
  padding: 0;
  box-sizing: border-box;
}

body {
  font-family: 'Helvetica Neue', Helvetica, 'PingFang SC', 'Hiragino Sans GB',
    'Microsoft YaHei', '微软雅黑', Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  background-color: var(--page-bg);
  color: var(--text-color);
  transition: background-color 0.3s, color 0.3s;
}

.app-container {
  min-height: 100vh;
  display: flex;
  flex-direction: column;
}

/* 通用动画 */
.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.3s ease;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}

/* 响应式样式 */
@media (max-width: 768px) {
  .hidden-mobile {
    display: none;
  }
}
</style>
