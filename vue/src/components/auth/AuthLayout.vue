<script setup lang="ts">
</script>

<template>
  <div class="auth-container">
    <div class="auth-background">
      <div class="auth-content">
        <router-view v-slot="{ Component }">
          <transition name="fade" mode="out-in">
            <component :is="Component" />
          </transition>
        </router-view>
      </div>
      <div class="auth-footer">
        <p>© {{ new Date().getFullYear() }} DeepSeek. All rights reserved.</p>
        <div class="footer-links">
          <el-link type="info" href="#" class="footer-link">关于我们</el-link>
          <el-link type="info" href="#" class="footer-link">服务条款</el-link>
          <el-link type="info" href="#" class="footer-link">隐私政策</el-link>
          <el-link type="info" href="#" class="footer-link">帮助中心</el-link>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.auth-container {
  min-height: 100vh;
  display: flex;
  flex-direction: column;
  background: linear-gradient(135deg, #e0f2fe 0%, #f0f9ff 100%);
}

.auth-background {
  flex: 1;
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;
  padding: 20px;
}

.auth-content {
  width: 100%;
  max-width: 400px;
  margin-bottom: 40px;
}

.auth-footer {
  text-align: center;
  padding: 20px;
  color: #666;
}

.footer-links {
  display: flex;
  justify-content: center;
  gap: 20px;
  margin-top: 10px;
}

.footer-link {
  font-size: 14px;
}

/* 过渡动画 */
.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.3s ease;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}
</style> 