<script setup lang="ts">
import { ElContainer, ElHeader, ElMain, ElFooter, ElAside } from 'element-plus'
import Header from './components/Header.vue'
import Aside from './components/Aside.vue'
import { ref } from 'vue'

const isCollapse = ref(false)
</script>

<template>
  <el-container class="app-container">
    <el-header class="app-header" height="60px">
      <Header />
    </el-header>
    
    <el-container class="main-container">
      <el-aside :width="isCollapse ? '64px' : '260px'" class="app-aside">
        <Aside />
      </el-aside>
      
      <el-container>
        <el-main class="app-main">
          <router-view></router-view>
        </el-main>
      </el-container>
    </el-container>
  </el-container>
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
