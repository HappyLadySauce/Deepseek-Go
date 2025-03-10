<script setup lang="ts">
import { computed } from 'vue'
import { useRoute } from 'vue-router'
import { ElContainer, ElHeader, ElMain, ElFooter, ElAside } from 'element-plus'
import Header from './components/Header.vue'
import Aside from './components/Aside.vue'
import { ref } from 'vue'
import { useRouter } from 'vue-router'

const route = useRoute()
const isCollapse = ref(false)
const router = useRouter()

// 判断当前是否在认证页面
const isAuthPage = computed(() => {
  return route.path.startsWith('/auth')
})

// 在组件中监听菜单点击
const handleMenuClick = (index: string) => {
  switch(index) {
    case '1':
      router.push('/chat')
      break
    case '2':
      router.push('/history')
      break
    case '3':
      router.push('/settings')
      break
  }
}

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
        <el-aside :width="isCollapse ? '64px' : '260px'" class="app-aside">
          <Aside />
        </el-aside>

        <!-- 主内容 -->
        <el-container>
          <el-main class="app-main">
            <router-view></router-view>
          </el-main>
        </el-container>
      </el-container>
    </el-container>
  </template>
</template>

<style>
html, body {
  margin: 0;
  padding: 0;
  height: 100vh;
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, 'Helvetica Neue', Arial, sans-serif;
}

#app {
  height: 100vh;
}

.app-container {
  height: 100vh;
}

.app-header {
  padding: 0;
  background: #fff;
  border-bottom: 1px solid #e6e6e6;
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  z-index: 1000;
}

.main-container {
  margin-top: 60px;
  height: calc(100vh - 60px);
}

.app-aside {
  background: #fff;
  transition: width 0.3s;
  overflow: hidden;
}

.app-main {
  background: #f5f7fa;
  padding: 20px;
}

/* 滚动条样式 */
::-webkit-scrollbar {
  width: 6px;
  height: 6px;
}

::-webkit-scrollbar-track {
  background: #f1f1f1;
  border-radius: 3px;
}

::-webkit-scrollbar-thumb {
  background: #c1c1c1;
  border-radius: 3px;
}

::-webkit-scrollbar-thumb:hover {
  background: #a8a8a8;
}

/* Element Plus 样式覆盖 */
.el-button {
  border-radius: 6px;
}

.el-input__inner {
  border-radius: 6px !important;
}
</style>
