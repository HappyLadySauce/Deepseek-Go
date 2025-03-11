# 创建一个现代化的顶部导航栏
<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { ChatDotSquare, User, Setting, Sunny, MoonNight } from '@element-plus/icons-vue'
import { toggleTheme, currentTheme } from '@/utils/theme'

const router = useRouter()
const activeIndex = ref('1')

const handleCommand = (command: string) => {
  switch(command) {
    case 'profile':
      router.push('/settings')
      break
    case 'logout':
      localStorage.removeItem('token')
      router.push('/auth/login')
      break
  }
}
</script>

<template>
  <div class="header-container">
    <div class="header-left">
      <div class="logo-section">
        <img src="../assets/logo.png" alt="Logo" class="logo-img" />
        <span class="divider"></span>
        <span class="product-name">DeepSeek-Go —— 集成 AI 的智能运维监控平台</span>
      </div>
      <div class="region-selector">
        <span class="region-text">DeepSeek Chat</span>
        <el-icon><ChatDotSquare /></el-icon>
      </div>
    </div>

    <div class="header-right">
      <div class="theme-toggle">
        <el-button 
          type="text" 
          @click="toggleTheme" 
          class="theme-button"
          :title="currentTheme === 'dark' ? '切换到亮色模式' : '切换到暗色模式'"
        >
          <el-icon>
            <component :is="currentTheme === 'dark' ? Sunny : MoonNight" />
          </el-icon>
        </el-button>
      </div>
      
      <div class="user-section">
        <el-dropdown @command="handleCommand" trigger="click">
          <div class="user-info">
            <el-avatar 
              :size="32" 
              src="https://cube.elemecdn.com/0/88/03b0d39583f48206768a7534e55bcpng.png" 
              class="user-avatar"
            />
            <span class="username">管理员</span>
          </div>
          <template #dropdown>
            <el-dropdown-menu>
              <el-dropdown-item command="profile">
                <el-icon><User /></el-icon>
                <span>个人中心</span>
              </el-dropdown-item>
              <el-dropdown-item command="settings">
                <el-icon><Setting /></el-icon>
                <span>账号设置</span>
              </el-dropdown-item>
              <el-dropdown-item divided command="logout">退出登录</el-dropdown-item>
            </el-dropdown-menu>
          </template>
        </el-dropdown>
      </div>
    </div>
  </div>
</template>

<style scoped>
.header-container {
  height: 100%;
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 0 16px;
  background: var(--header-bg);
  color: var(--text-color);
  transition: background-color 0.3s ease;
}

.header-left {
  display: flex;
  align-items: center;
  gap: 24px;
}

.logo-section {
  display: flex;
  align-items: center;
  gap: 12px;
}

.logo-img {
  height: 28px;
  width: auto;
}

.divider {
  width: 1px;
  height: 20px;
  background: var(--border-color);
}

.product-name {
  font-size: 16px;
  font-weight: 500;
  color: var(--menu-text-color);
  background: linear-gradient(to right, var(--primary-color), var(--secondary-color));
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
}

.region-selector {
  display: flex;
  align-items: center;
  gap: 8px;
  cursor: pointer;
  padding: 4px 8px;
  border-radius: 4px;
  transition: all 0.3s;
  color: var(--menu-text-color);
}

.region-selector:hover {
  background: var(--hover-color);
  color: var(--secondary-color);
}

.region-text {
  font-size: 14px;
}

.header-right {
  display: flex;
  align-items: center;
  gap: 24px;
}

.theme-toggle {
  display: flex;
  align-items: center;
  justify-content: center;
}

.theme-button {
  padding: 8px;
  border-radius: 50%;
  transition: all 0.3s ease;
  color: var(--menu-text-color) !important;
  border: none !important;
}

.theme-button:hover {
  background-color: var(--hover-color) !important;
  color: var(--primary-color) !important;
  transform: rotate(30deg);
}

.theme-button .el-icon {
  font-size: 18px;
}

.user-section {
  display: flex;
  align-items: center;
}

.user-info {
  display: flex;
  align-items: center;
  gap: 8px;
  cursor: pointer;
  padding: 4px 8px;
  border-radius: 4px;
  transition: all 0.3s;
}

.user-info:hover {
  background: var(--hover-color);
}

.username {
  font-size: 14px;
  color: var(--menu-text-color);
}

:deep(.el-dropdown-menu__item) {
  display: flex;
  align-items: center;
  gap: 8px;
}

:deep(.el-dropdown-menu__item i) {
  margin-right: 0;
}
</style>
