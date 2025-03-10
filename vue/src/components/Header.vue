# 创建一个现代化的顶部导航栏
<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { Search, User, Setting, SwitchButton } from '@element-plus/icons-vue'

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
        <span class="product-name">轻量应用服务器</span>
      </div>
      <div class="region-selector">
        <span class="region-text">成都</span>
        <el-icon class="region-icon"><SwitchButton /></el-icon>
      </div>
    </div>

    <div class="header-right">
      <div class="search-section">
        <el-input
          placeholder="搜索资源、产品或文档"
          :prefix-icon="Search"
          class="search-input"
        />
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
  background: #1e2f40;
  color: #fff;
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
  background: rgba(255, 255, 255, 0.2);
}

.product-name {
  font-size: 16px;
  font-weight: 500;
  color: #fff;
}

.region-selector {
  display: flex;
  align-items: center;
  gap: 8px;
  cursor: pointer;
  padding: 4px 8px;
  border-radius: 4px;
  transition: all 0.3s;
}

.region-selector:hover {
  background: rgba(255, 255, 255, 0.1);
}

.region-text {
  font-size: 14px;
}

.region-icon {
  font-size: 14px;
  color: rgba(255, 255, 255, 0.7);
}

.header-right {
  display: flex;
  align-items: center;
  gap: 24px;
}

.search-section {
  width: 300px;
}

.search-input :deep(.el-input__wrapper) {
  background: rgba(255, 255, 255, 0.1);
  box-shadow: none;
  border: 1px solid rgba(255, 255, 255, 0.1);
}

.search-input :deep(.el-input__inner) {
  color: #fff;
  height: 32px;
}

.search-input :deep(.el-input__inner::placeholder) {
  color: rgba(255, 255, 255, 0.5);
}

.search-input :deep(.el-input__prefix-icon) {
  color: rgba(255, 255, 255, 0.5);
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
  background: rgba(255, 255, 255, 0.1);
}

.username {
  font-size: 14px;
  color: #fff;
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
