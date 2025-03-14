<script setup lang="ts">
import { computed, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { ElContainer, ElHeader, ElMain, ElFooter, ElAside } from 'element-plus'
import Header from './components/Header.vue'
import Aside from './components/Aside.vue'
import { ref, provide, watchEffect } from 'vue'
import { useRouter } from 'vue-router'
import { applyTheme, initTheme } from '@/utils/theme'

// 初始化主题
onMounted(() => {
  console.log('App.vue: 初始化主题')
  initTheme()
})

const route = useRoute()
// 将 isCollapse 状态提升到父组件
const isCollapse = ref(false)
const router = useRouter()

// 提供折叠状态给子组件
provide('isCollapse', isCollapse)

const toggleCollapse = () => {
  isCollapse.value = !isCollapse.value
}

// 提供切换折叠函数给子组件
provide('toggleCollapse', toggleCollapse)

// 判断当前是否在认证页面
const isAuthPage = computed(() => {
  return route.path.startsWith('/auth')
})

// 处理用户相关操作
const handleCommand = (command: string) => {
  switch(command) {
    case 'profile':
      router.push('/settings')
      break
    case 'logout':
      // 处理登出逻辑
      localStorage.removeItem('token')
      router.push('/auth/login')
      break
  }
}
</script>

<template>
  <!-- 根据路由判断显示认证页面还是主应用布局 -->
  <template v-if="isAuthPage">
    <router-view></router-view>
  </template>

  <!-- 主应用布局 -->
  <template v-else>
    <el-container class="app-container">
      <!-- 头部 -->
      <el-header class="app-header" height="60px">
        <Header />
      </el-header>

      <!-- 主容器 -->
      <el-container class="main-container">
        <!-- 侧边栏 -->
        <el-aside :width="isCollapse ? '64px' : '240px'" class="app-aside">
          <Aside />
        </el-aside>

        <!-- 主内容 -->
        <el-main class="app-main">
          <router-view></router-view>
        </el-main>
      </el-container>
    </el-container>
  </template>
</template>

<style>
/* App全局样式 - 使用CSS变量确保主题适用 */
.app-container {
  height: 100vh;
}

.app-header {
  padding: 0;
  background: var(--header-bg);
  border-bottom: 1px solid var(--border-color);
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  z-index: 1000;
  transition: background-color 0.3s ease;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}

.main-container {
  margin-top: 60px;
  height: calc(100vh - 60px);
}

.app-aside {
  background: var(--card-bg);
  transition: width 0.3s ease, background-color 0.3s ease;
  overflow: hidden;
  border-right: 1px solid var(--border-color);
  box-shadow: 2px 0 8px rgba(0, 0, 0, 0.1);
}

.app-main {
  background: var(--page-bg);
  padding: 20px;
  transition: background-color 0.3s ease;
  flex: 1;
  overflow: auto;
}

#app {
  height: 100vh;
  color: var(--text-color);
  background-color: var(--page-bg);
  transition: background-color 0.3s, color 0.3s;
}

/* Element Plus 样式覆盖 */
.el-button {
  border-radius: 6px;
}

.el-input__inner {
  border-radius: 6px !important;
}

/* 根据主题自动调整下拉菜单颜色 */
.el-dropdown-menu {
  background-color: var(--card-bg) !important;
  border-color: var(--border-color) !important;
  box-shadow: 0 3px 6px rgba(0, 0, 0, 0.1) !important;
}

.el-dropdown-menu__item {
  color: var(--text-color) !important;
}

.el-dropdown-menu__item:hover {
  background-color: var(--hover-color) !important;
}

/* 修复主内容区扩展问题 */
.el-container {
  height: 100%;
}
</style>
