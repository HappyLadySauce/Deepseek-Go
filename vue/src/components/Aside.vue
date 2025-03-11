<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { 
  ChatLineRound,
  Document,
  Setting,
  Monitor,
  Connection,
  Warning,
  Refresh
} from '@element-plus/icons-vue'

// 路由
const router = useRouter()
// 是否折叠
const isCollapse = ref(false)
// 当前选中的菜单
const activeMenu = ref('overview')

// 监控子菜单是否展开
const isMonitorSubMenuOpen = ref(false)

// 菜单点击事件
const handleSelect = (key: string) => {
  switch(key) {
    case 'overview':
      router.push('/overview')
      break
    case 'system-monitor':
      router.push('/monitor/system')
      break
    case 'network-monitor':
      router.push('/monitor/network')
      break
    case 'alarm-monitor':
      router.push('/monitor/alarm')
      break
    case 'settings':
      router.push('/settings')
      break
  }
}
</script>

<template>
  <div class="aside-container">
    <el-menu
      :default-active="activeMenu"
      class="aside-menu"
      :collapse="isCollapse"
      background-color="#1e2f40"
      text-color="#fff"
      active-text-color="#409EFF"
      @select="handleSelect"
    >
      <div class="menu-header">
        <span class="menu-title" v-if="!isCollapse">实例列表</span>
        <el-icon class="refresh-icon" title="刷新"><Refresh /></el-icon>
      </div>

      <el-menu-item index="overview">
        <el-icon><ElementPlus /></el-icon>
        <template #title>总体概览</template>
      </el-menu-item>

      <el-sub-menu index="monitor">
        <template #title>
          <el-icon><Odometer /></el-icon>
          <span>实时监控</span>
        </template>
        <el-menu-item index="system-monitor">
          <el-icon><Monitor /></el-icon>
          系统监控
        </el-menu-item>
        <el-menu-item index="network-monitor">
          <el-icon><Connection /></el-icon>
          网络监控
        </el-menu-item>
        <el-menu-item index="alarm-monitor">
          <el-icon><Warning /></el-icon>
          监控告警
        </el-menu-item>
      </el-sub-menu>

      <el-menu-item index="settings">
        <el-icon><Setting /></el-icon>
        <template #title>系统配置</template>
      </el-menu-item>
    </el-menu>

    <div class="collapse-trigger" @click="isCollapse = !isCollapse">
      <el-icon :class="{ 'is-collapse': isCollapse }">
        <svg viewBox="0 0 1024 1024" xmlns="http://www.w3.org/2000/svg">
          <path fill="currentColor" d="M831.872 340.864L512 652.672 192.128 340.864a30.592 30.592 0 00-42.752 0 29.12 29.12 0 000 41.6L489.664 714.24a32 32 0 0044.672 0l340.288-331.712a29.12 29.12 0 000-41.728 30.592 30.592 0 00-42.752 0z"/>
        </svg>
      </el-icon>
    </div>
  </div>
</template>

<style scoped>
.aside-container {
  height: 100%;
  display: flex;
  flex-direction: column;
  background: #1e2f40;
  position: relative;
}

.aside-menu {
  flex: 1;
  border-right: none;
}

.menu-header {
  height: 48px;
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 20px;
  color: #fff;
  border-bottom: 1px solid rgba(255, 255, 255, 0.1);
}

.menu-title {
  font-size: 14px;
  font-weight: 500;
}

.refresh-icon {
  font-size: 16px;
  cursor: pointer;
  padding: 4px;
  border-radius: 4px;
  transition: all 0.3s;
}

.refresh-icon:hover {
  background: rgba(255, 255, 255, 0.1);
}

.collapse-trigger {
  height: 40px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #fff;
  cursor: pointer;
  border-top: 1px solid rgba(255, 255, 255, 0.1);
  transition: all 0.3s;
}

.collapse-trigger:hover {
  background: rgba(255, 255, 255, 0.1);
}

.collapse-trigger .el-icon {
  transition: transform 0.3s;
  font-size: 16px;
}

.collapse-trigger .is-collapse {
  transform: rotate(180deg);
}

:deep(.el-menu) {
  border-right: none;
}

:deep(.el-menu-item) {
  height: 40px;
  line-height: 40px;
  
  &:hover {
    background-color: rgba(255, 255, 255, 0.1) !important;
  }
  
  &.is-active {
    background-color: rgba(64, 158, 255, 0.1) !important;
  }
}

:deep(.el-menu--collapse) {
  width: 64px;
}

:deep(.el-menu-item .el-icon) {
  font-size: 18px;
}

:deep(.el-sub-menu__title) {
  height: 40px !important;
  line-height: 40px !important;
}

:deep(.el-sub-menu .el-menu-item) {
  height: 40px !important;
  line-height: 40px !important;
  padding-left: 50px !important;
}
</style>
