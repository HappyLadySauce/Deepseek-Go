<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted, watch } from 'vue'
import { useRouter } from 'vue-router'
import { useUserStore } from '../stores/user'
import { useChatStore } from '../stores/chat'
import { ElMessageBox } from 'element-plus'
import { storeToRefs } from 'pinia'
import { markRaw } from 'vue'
import { Document, ChatDotRound, Setting, User, Moon, Sunny, SwitchButton, Delete } from '@element-plus/icons-vue'
import { currentTheme, toggleTheme } from '../utils/theme'

const router = useRouter()
const userStore = useUserStore()
const chatStore = useChatStore()

// 从store中获取引用
const { userInfo } = storeToRefs(userStore)
const { chatSessions, currentSessionId, loading } = storeToRefs(chatStore)

// 消息输入
const messageInput = ref('')
const inputRef = ref<HTMLTextAreaElement>()

// 消息容器滚动
const messagesContainer = ref<HTMLElement>()

// 侧边栏状态
const isCollapsed = ref(false)
const isMobile = ref(window.innerWidth <= 768)

// 清空消息确认
const confirmClearMessages = () => {
  ElMessageBox.confirm(
    '确定要清空所有消息吗？此操作不可恢复。',
    '清空确认',
    {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    }
  ).then(() => {
    if (currentSessionId.value) {
      chatStore.clearSession(currentSessionId.value)
    }
  }).catch(() => {
    // 取消操作
  })
}

// 删除会话确认
const confirmDeleteSession = (sessionId: string) => {
  ElMessageBox.confirm(
    '确定要删除此会话吗？此操作不可恢复。',
    '删除确认',
    {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    }
  ).then(() => {
    chatStore.deleteSession(sessionId)
  }).catch(() => {
    // 取消操作
  })
}

// 当前会话
const currentSession = computed(() => {
  return chatStore.getCurrentSession()
})

// 发送消息
const sendMessage = async () => {
  if (!messageInput.value.trim()) return
  
  const message = messageInput.value
  messageInput.value = ''
  
  await chatStore.sendMessage(message)
  
  // 自动聚焦输入框
  if (inputRef.value) {
    inputRef.value.focus()
  }
}

// 监听Enter键发送消息
const handleKeyDown = (e: KeyboardEvent) => {
  if (e.key === 'Enter' && !e.shiftKey) {
    e.preventDefault()
    sendMessage()
  }
}

// 滚动到底部
const scrollToBottom = () => {
  if (messagesContainer.value) {
    setTimeout(() => {
      messagesContainer.value!.scrollTop = messagesContainer.value!.scrollHeight
    }, 10)
  }
}

// 响应式调整
const handleResize = () => {
  isMobile.value = window.innerWidth <= 768
  if (isMobile.value) {
    isCollapsed.value = true
  }
}

// 导航图标
const navItems = [
  { icon: markRaw(ChatDotRound), label: '聊天', path: '/chat' },
  { icon: markRaw(User), label: '个人中心', path: '/profile' },
  { icon: markRaw(Setting), label: '设置', path: '/settings' }
]

// 生命周期钩子
onMounted(() => {
  // 初始化聊天会话
  chatStore.initialize()
  
  // 监听窗口大小变化
  window.addEventListener('resize', handleResize)
  handleResize()
  
  // 自动聚焦输入框
  if (inputRef.value) {
    inputRef.value.focus()
  }
})

onUnmounted(() => {
  // 保存聊天会话
  chatStore.saveToLocalStorage()
  
  // 移除事件监听
  window.removeEventListener('resize', handleResize)
})

// 当前会话或消息变化时滚动到底部
watch(
  () => currentSession.value?.messages.length,
  () => scrollToBottom(),
  { immediate: true }
)

// 当前会话ID变化时滚动到底部
watch(
  () => currentSessionId.value,
  () => scrollToBottom(),
  { immediate: true }
)

// 处理注销
const handleLogout = () => {
  ElMessageBox.confirm(
    '确定要退出登录吗？',
    '退出登录',
    {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    }
  ).then(() => {
    // 保存聊天会话
    chatStore.saveToLocalStorage()
    
    // 注销
    userStore.logout()
  }).catch(() => {
    // 取消操作
  })
}

// 切换主题
const isDarkMode = computed(() => currentTheme.value === 'dark')
const handleToggleTheme = () => {
  toggleTheme()
}
</script>

<template>
  <div class="chat-layout">
    <!-- 侧边导航 -->
    <div class="chat-sidebar" :class="{ 'is-collapsed': isCollapsed }">
      <div class="sidebar-header">
        <div class="app-logo">DeepSeek</div>
        <el-button
          v-if="isMobile"
          circle
          text
          class="collapse-button"
          @click="isCollapsed = !isCollapsed"
        >
          <el-icon><i class="el-icon-arrow-left" /></el-icon>
        </el-button>
      </div>
      
      <!-- 会话列表 -->
      <div class="session-list">
        <div class="session-actions">
          <el-button type="primary" round @click="chatStore.createNewSession">
            <el-icon class="mr-5"><ChatDotRound /></el-icon>
            新对话
          </el-button>
        </div>
        
        <el-scrollbar>
          <div
            v-for="session in chatSessions"
            :key="session.id"
            class="session-item"
            :class="{ 'is-active': session.id === currentSessionId }"
            @click="chatStore.switchSession(session.id)"
          >
            <div class="session-title">{{ session.title }}</div>
            <div class="session-actions">
              <el-button
                circle
                text
                size="small"
                @click.stop="confirmDeleteSession(session.id)"
                class="delete-button"
              >
                <el-icon class="delete-icon"><Delete /></el-icon>
              </el-button>
            </div>
          </div>
        </el-scrollbar>
      </div>
      
      <!-- 用户信息 -->
      <div class="user-info">
        <div class="user-details">
          <div class="username">{{ userInfo?.username }}</div>
          <div class="user-email">{{ userInfo?.email || '' }}</div>
        </div>
        <el-dropdown>
          <span class="el-dropdown-link">
            <el-avatar :size="32" class="user-avatar">
              {{ userInfo?.username?.charAt(0).toUpperCase() }}
            </el-avatar>
          </span>
          <template #dropdown>
            <el-dropdown-menu>
              <el-dropdown-item @click="router.push('/profile')">
                <el-icon><User /></el-icon> 个人中心
              </el-dropdown-item>
              <el-dropdown-item @click="handleToggleTheme">
                <el-icon><component :is="isDarkMode ? Sunny : Moon" /></el-icon>
                {{ isDarkMode ? '切换亮色模式' : '切换暗色模式' }}
              </el-dropdown-item>
              <el-dropdown-item divided @click="handleLogout">
                <el-icon><SwitchButton /></el-icon> 退出登录
              </el-dropdown-item>
            </el-dropdown-menu>
          </template>
        </el-dropdown>
      </div>
    </div>
    
    <!-- 主内容区 -->
    <div class="chat-main">
      <div class="chat-header">
        <div class="title-section">
          <el-button
            v-if="isMobile"
            circle
            text
            class="menu-button"
            @click="isCollapsed = !isCollapsed"
          >
            <el-icon><i class="el-icon-menu" /></el-icon>
          </el-button>
          <h2>{{ currentSession?.title || '新对话' }}</h2>
        </div>
        
        <div class="nav-section">
          <div class="nav-items">
            <div
              v-for="(item, index) in navItems"
              :key="index"
              class="nav-item"
              @click="router.push(item.path)"
            >
              <el-icon><component :is="item.icon" /></el-icon>
              <span>{{ item.label }}</span>
            </div>
            
            <!-- 添加主题切换按钮 -->
            <div class="nav-item" @click="handleToggleTheme">
              <el-icon><component :is="isDarkMode ? Sunny : Moon" /></el-icon>
              <span>{{ isDarkMode ? '亮色' : '暗色' }}</span>
            </div>
          </div>
          
          <el-dropdown v-if="currentSession">
            <el-button text>
              <el-icon><i class="el-icon-more" /></el-icon>
            </el-button>
            <template #dropdown>
              <el-dropdown-menu>
                <el-dropdown-item @click="confirmClearMessages">
                  清空消息
                </el-dropdown-item>
                <el-dropdown-item divided @click="handleLogout">
                  退出登录
                </el-dropdown-item>
              </el-dropdown-menu>
            </template>
          </el-dropdown>
        </div>
      </div>
      
      <!-- 聊天内容区 -->
      <div class="chat-content">
        <div
          v-if="!currentSession || currentSession.messages.length === 0"
          class="empty-chat"
        >
          <div class="empty-chat-content">
            <h3>开始新的对话</h3>
            <p>AI助手可以回答您的问题，请在下方输入您的问题开始对话。</p>
          </div>
        </div>
        
        <el-scrollbar v-else ref="messagesContainer" class="messages-container">
          <div
            v-for="message in currentSession.messages"
            :key="message.id"
            class="message-item"
            :class="{ 'user-message': message.role === 'user' }"
          >
            <div class="message-avatar">
              <el-avatar 
                :size="36" 
                :icon="message.role === 'user' ? User : Document"
              />
            </div>
            <div class="message-content">
              <div class="message-text">{{ message.content }}</div>
              <div class="message-time">
                {{ new Date(message.timestamp).toLocaleTimeString() }}
              </div>
            </div>
          </div>
        </el-scrollbar>
      </div>
      
      <!-- 消息输入区 -->
      <div class="chat-input-area">
        <div class="input-container">
          <el-input
            v-model="messageInput"
            ref="inputRef"
            type="textarea"
            :rows="3"
            placeholder="请输入您的问题..."
            resize="none"
            @keydown="handleKeyDown"
          />
          
          <div class="input-actions">
            <el-button
              type="primary"
              :loading="loading"
              @click="sendMessage"
              :disabled="!messageInput.trim()"
            >
              发送
            </el-button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.chat-layout {
  display: flex;
  height: 100vh;
  overflow: hidden;
}

.chat-sidebar {
  display: flex;
  flex-direction: column;
  width: 260px;
  background-color: var(--card-bg);
  border-right: 1px solid var(--border-color);
  transition: all 0.3s;
  z-index: 10;
}

.chat-sidebar.is-collapsed {
  width: 0;
  overflow: hidden;
}

.sidebar-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16px;
  border-bottom: 1px solid var(--border-color);
}

.app-logo {
  font-size: 18px;
  font-weight: bold;
  color: var(--primary-color);
}

.session-list {
  flex: 1;
  overflow: hidden;
  display: flex;
  flex-direction: column;
}

.session-actions {
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 16px;
  opacity: 1;
}

.session-actions .el-button {
  width: 100%;
  margin-bottom: 10px;
  transition: all 0.3s ease;
}

.session-actions .el-button:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
}

.session-item {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 12px 16px;
  cursor: pointer;
  border-bottom: 1px solid var(--border-color);
  transition: all 0.3s;
  color: var(--text-color);
}

.session-item:hover {
  background-color: var(--hover-color);
}

.session-item.is-active {
  background-color: var(--active-color);
  border-left: 3px solid var(--primary-color);
}

.session-item .session-actions {
  padding: 0;
  opacity: 0;
  transition: opacity 0.3s;
}

.session-item:hover .session-actions {
  opacity: 1;
}

.session-title {
  flex: 1;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.user-info {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 16px;
  border-top: 1px solid var(--border-color);
  background-color: var(--card-bg);
}

.user-details {
  display: flex;
  flex-direction: column;
}

.username {
  font-weight: bold;
  color: var(--text-color);
}

.user-email {
  font-size: 12px;
  color: var(--text-muted);
  margin-top: 4px;
}

.chat-main {
  flex: 1;
  display: flex;
  flex-direction: column;
  overflow: hidden;
  background-color: var(--page-bg);
}

.chat-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16px;
  border-bottom: 1px solid var(--border-color);
  background-color: var(--header-bg);
  color: var(--text-color);
}

.title-section {
  display: flex;
  align-items: center;
}

.title-section h2 {
  margin: 0 0 0 10px;
  font-size: 18px;
  color: var(--heading-color);
}

.nav-section {
  display: flex;
  align-items: center;
}

.nav-items {
  display: flex;
  gap: 16px;
}

.nav-item {
  display: flex;
  flex-direction: column;
  align-items: center;
  cursor: pointer;
  padding: 4px 8px;
  border-radius: 4px;
  transition: background-color 0.3s;
  color: var(--text-color);
  position: relative;
}

.nav-item .el-icon {
  font-size: 20px;
  margin-bottom: 4px;
  color: var(--text-color);
}

.nav-item:hover {
  background-color: var(--hover-color);
}

.nav-item span {
  font-size: 12px;
  margin-top: 4px;
}

.chat-content {
  flex: 1;
  overflow: hidden;
  position: relative;
  background-color: var(--page-bg);
}

.empty-chat {
  display: flex;
  justify-content: center;
  align-items: center;
  height: 100%;
  background-color: var(--page-bg);
}

.empty-chat-content {
  text-align: center;
  max-width: 400px;
  padding: 20px;
}

.empty-chat-content h3 {
  margin-bottom: 12px;
  color: var(--heading-color);
}

.empty-chat-content p {
  color: var(--text-muted);
  line-height: 1.6;
}

.messages-container {
  height: 100%;
  padding: 16px;
}

.message-item {
  display: flex;
  margin-bottom: 16px;
}

.message-avatar {
  margin-right: 12px;
}

.message-content {
  flex: 1;
  background-color: var(--assistant-bg);
  padding: 12px;
  border-radius: 8px;
  max-width: 80%;
  box-shadow: var(--message-shadow);
}

.user-message .message-content {
  background-color: var(--user-bg);
  color: var(--text-color);
}

.message-text {
  word-break: break-word;
  white-space: pre-wrap;
}

.message-time {
  font-size: 12px;
  color: var(--text-muted);
  margin-top: 8px;
  text-align: right;
}

.chat-input-area {
  padding: 16px;
  border-top: 1px solid var(--border-color);
  background-color: var(--card-bg);
}

.input-container {
  display: flex;
  flex-direction: column;
}

.input-actions {
  display: flex;
  justify-content: flex-end;
  margin-top: 12px;
}

/* 响应式样式 */
@media (max-width: 768px) {
  .chat-sidebar {
    position: absolute;
    height: 100%;
  }
  
  .nav-items {
    display: none;
  }
}

.delete-button {
  color: var(--text-color) !important;
  transition: all 0.3s;
}

.delete-button:hover {
  color: #f56c6c !important;
  background-color: rgba(245, 108, 108, 0.1);
}

.delete-icon {
  font-size: 16px;
}

.mr-5 {
  margin-right: 5px;
}
</style> 