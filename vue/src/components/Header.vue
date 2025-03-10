# 创建一个现代化的顶部导航栏
<script setup lang="ts">
import { ref } from 'vue'
import { Search } from '@element-plus/icons-vue'

const activeIndex = ref('1')
</script>

<template>
  <div class="header-container">
    <div class="logo">
      <img src="../assets/logo.png" alt="Logo" class="logo-img" />
      <span class="logo-text">DeepSeek Chat</span>
    </div>
    
    <el-menu
      :default-active="activeIndex"
      class="nav-menu"
      mode="horizontal"
      :ellipsis="false"
    >
      <el-menu-item index="1">首页</el-menu-item>
      <el-menu-item index="2">对话</el-menu-item>
      <el-menu-item index="3">历史记录</el-menu-item>
      <el-menu-item index="4">设置</el-menu-item>
      
      <div class="flex-grow" />
      
      <el-input
        class="search-input"
        placeholder="搜索对话..."
        :prefix-icon="Search"
      />
      
      <el-dropdown class="avatar-dropdown">
        <el-avatar :size="32" src="https://cube.elemecdn.com/0/88/03b0d39583f48206768a7534e55bcpng.png" />
        <template #dropdown>
          <el-dropdown-menu>
            <el-dropdown-item>个人中心</el-dropdown-item>
            <el-dropdown-item>账号设置</el-dropdown-item>
            <el-dropdown-item divided>退出登录</el-dropdown-item>
          </el-dropdown-menu>
        </template>
      </el-dropdown>
    </el-menu>
  </div>
</template>

<style scoped>
.header-container {
  height: 60px;
  display: flex;
  align-items: center;
  padding: 0 20px;
  background: #fff;
  box-shadow: 0 2px 8px rgba(0,0,0,0.1);
}

.logo {
  display: flex;
  align-items: center;
  margin-right: 40px;
}

.logo-img {
  height: 32px;
  margin-right: 10px;
}

.logo-text {
  font-size: 20px;
  font-weight: 600;
  color: #409EFF;
}

.nav-menu {
  flex: 1;
  border: none;
}

.flex-grow {
  flex-grow: 1;
}

.search-input {
  width: 200px;
  margin-right: 20px;
}

.avatar-dropdown {
  margin-left: 20px;
  cursor: pointer;
}

:deep(.el-menu--horizontal) {
  border-bottom: none;
}

:deep(.el-menu-item) {
  font-size: 15px;
}

:deep(.el-menu-item.is-active) {
  font-weight: 600;
}
</style>
