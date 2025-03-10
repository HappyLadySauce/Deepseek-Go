<script setup lang="ts">
import { ref } from 'vue'
import { Plus, Message, Delete, Edit } from '@element-plus/icons-vue'

const activeChat = ref('1')
const chatList = ref([
  { id: '1', title: '关于人工智能的讨论', time: '10:30', preview: '人工智能正在改变我们的生活方式...' },
  { id: '2', title: '项目开发计划', time: '昨天', preview: '让我们讨论一下新项目的具体实施计划' },
  { id: '3', title: '技术架构设计', time: '2天前', preview: '针对新系统的技术架构，我们需要考虑以下几点' },
])
</script>

<template>
  <div class="aside-container">
    <div class="new-chat">
      <el-button type="primary" class="new-chat-btn" :icon="Plus">
        新建对话
      </el-button>
    </div>

    <div class="chat-list">
      <el-scrollbar height="calc(100vh - 180px)">
        <div
          v-for="chat in chatList"
          :key="chat.id"
          class="chat-item"
          :class="{ 'active': activeChat === chat.id }"
          @click="activeChat = chat.id"
        >
          <el-icon class="chat-icon"><Message /></el-icon>
          <div class="chat-content">
            <div class="chat-header">
              <span class="chat-title">{{ chat.title }}</span>
              <span class="chat-time">{{ chat.time }}</span>
            </div>
            <div class="chat-preview">{{ chat.preview }}</div>
          </div>
          <div class="chat-actions">
            <el-icon class="action-icon"><Edit /></el-icon>
            <el-icon class="action-icon"><Delete /></el-icon>
          </div>
        </div>
      </el-scrollbar>
    </div>
  </div>
</template>

<style scoped>
.aside-container {
  height: 100%;
  background: #f5f7fa;
  border-right: 1px solid #e6e6e6;
}

.new-chat {
  padding: 16px;
  border-bottom: 1px solid #e6e6e6;
}

.new-chat-btn {
  width: 100%;
}

.chat-list {
  padding: 8px;
}

.chat-item {
  display: flex;
  align-items: flex-start;
  padding: 12px;
  margin-bottom: 4px;
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.3s;
  position: relative;
}

.chat-item:hover {
  background: #ecf5ff;
}

.chat-item.active {
  background: #ecf5ff;
}

.chat-icon {
  font-size: 20px;
  color: #409EFF;
  margin-right: 12px;
  margin-top: 2px;
}

.chat-content {
  flex: 1;
  min-width: 0;
}

.chat-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 4px;
}

.chat-title {
  font-weight: 500;
  color: #333;
  font-size: 14px;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.chat-time {
  font-size: 12px;
  color: #999;
}

.chat-preview {
  font-size: 12px;
  color: #666;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.chat-actions {
  display: none;
  position: absolute;
  right: 12px;
  top: 50%;
  transform: translateY(-50%);
  background: #fff;
  border-radius: 4px;
  padding: 4px;
  box-shadow: 0 2px 12px rgba(0,0,0,0.1);
}

.chat-item:hover .chat-actions {
  display: flex;
  gap: 8px;
}

.action-icon {
  font-size: 16px;
  color: #666;
  cursor: pointer;
  padding: 4px;
  border-radius: 4px;
  transition: all 0.3s;
}

.action-icon:hover {
  background: #f5f7fa;
  color: #409EFF;
}
</style>
