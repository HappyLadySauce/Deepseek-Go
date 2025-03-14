<script setup lang="ts">
import { ref, reactive } from 'vue'
import { ElMessage } from 'element-plus'
import { useRouter } from 'vue-router'
import { ArrowLeft } from '@element-plus/icons-vue'
import { currentTheme } from '../utils/theme'

const router = useRouter()

// 主题设置
const theme = ref('light')
const themeOptions = [
  { label: '浅色主题', value: 'light' },
  { label: '深色主题', value: 'dark' },
  { label: '跟随系统', value: 'auto' }
]

// 语言设置
const language = ref('zh-CN')
const languageOptions = [
  { label: '简体中文', value: 'zh-CN' },
  { label: '繁體中文', value: 'zh-TW' },
  { label: 'English', value: 'en-US' }
]

// 聊天设置
const chatSettings = reactive({
  enableEnterToSend: true,
  enableAutoSave: true,
  messageDisplayLimit: 50
})

// 保存设置
const saveSettings = () => {
  // 实际应用中应该保存到后端或本地存储
  localStorage.setItem('theme', theme.value)
  localStorage.setItem('language', language.value)
  localStorage.setItem('chatSettings', JSON.stringify(chatSettings))
  
  ElMessage.success('设置已保存')
}

// 重置设置
const resetSettings = () => {
  theme.value = 'light'
  language.value = 'zh-CN'
  chatSettings.enableEnterToSend = true
  chatSettings.enableAutoSave = true
  chatSettings.messageDisplayLimit = 50
  
  ElMessage.success('设置已重置')
}
</script>

<template>
  <div class="settings-container">
    <div class="settings-header">
      <div class="header-content">
        <el-button 
          class="back-button"
          type="primary"
          plain
          round
          @click="router.push('/chat')" 
          :icon="ArrowLeft"
        >
          返回聊天
        </el-button>
        <h1>设置</h1>
      </div>
    </div>
    
    <div class="settings-content">
      <el-card class="settings-card">
        <template #header>
          <div class="card-header">
            <h3>外观设置</h3>
          </div>
        </template>
        
        <div class="settings-item">
          <span class="settings-label">主题</span>
          <el-radio-group v-model="theme">
            <el-radio 
              v-for="option in themeOptions"
              :key="option.value"
              :label="option.value"
            >
              {{ option.label }}
            </el-radio>
          </el-radio-group>
        </div>
        
        <div class="settings-item">
          <span class="settings-label">语言</span>
          <el-select v-model="language" class="settings-select">
            <el-option
              v-for="option in languageOptions"
              :key="option.value"
              :label="option.label"
              :value="option.value"
            />
          </el-select>
        </div>
      </el-card>
      
      <el-card class="settings-card">
        <template #header>
          <div class="card-header">
            <h3>聊天设置</h3>
          </div>
        </template>
        
        <div class="settings-item">
          <span class="settings-label">按Enter键发送消息</span>
          <el-switch v-model="chatSettings.enableEnterToSend" />
        </div>
        
        <div class="settings-item">
          <span class="settings-label">自动保存聊天记录</span>
          <el-switch v-model="chatSettings.enableAutoSave" />
        </div>
        
        <div class="settings-item">
          <span class="settings-label">消息显示数量限制</span>
          <el-slider
            v-model="chatSettings.messageDisplayLimit"
            :min="10"
            :max="200"
            :step="10"
            :marks="{
              10: '10',
              50: '50',
              100: '100',
              200: '200'
            }"
            class="settings-slider"
          />
        </div>
      </el-card>
      
      <div class="settings-actions">
        <el-button type="primary" @click="saveSettings">保存设置</el-button>
        <el-button @click="resetSettings">重置设置</el-button>
      </div>
    </div>
  </div>
</template>

<style scoped>
.settings-container {
  min-height: 100vh;
  background-color: var(--page-bg);
  padding-bottom: 40px;
}

.settings-header {
  background-color: #4facfe;
  background-image: linear-gradient(135deg, #4facfe 0%, #00f2fe 100%);
  color: white;
  padding: 40px 0;
  margin-bottom: 30px;
}

.header-content {
  max-width: 800px;
  margin: 0 auto;
  padding: 0 24px;
  display: flex;
  align-items: center;
}

.settings-header h1 {
  font-size: 28px;
  color: white;
  margin: 0 0 0 20px;
}

.back-button {
  transition: all 0.3s ease;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.1);
}

.back-button:hover {
  transform: translateX(-5px);
  box-shadow: 0 4px 15px rgba(0, 0, 0, 0.2);
}

.settings-content {
  max-width: 800px;
  margin: 0 auto;
  padding: 0 24px;
}

.settings-card {
  margin-bottom: 24px;
  border-radius: 8px;
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.1);
  transition: all 0.3s ease;
  background-color: var(--card-bg);
}

.settings-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 6px 20px rgba(0, 0, 0, 0.15);
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  border-bottom: 1px solid var(--border-color);
  padding-bottom: 12px;
}

.card-header h3 {
  margin: 0;
  font-size: 18px;
  color: var(--heading-color);
}

.settings-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16px 0;
  border-bottom: 1px solid var(--border-color);
}

.settings-item:last-child {
  border-bottom: none;
}

.settings-label {
  font-size: 16px;
  color: var(--text-color);
}

.settings-select {
  width: 200px;
}

.settings-slider {
  width: 60%;
  margin-left: 20px;
}

.settings-actions {
  display: flex;
  justify-content: center;
  gap: 16px;
  margin-top: 24px;
}

.settings-actions .el-button {
  min-width: 120px;
  padding: 12px 20px;
  font-weight: 500;
}

@media (max-width: 768px) {
  .settings-item {
    flex-direction: column;
    align-items: flex-start;
  }
  
  .settings-label {
    margin-bottom: 12px;
  }
  
  .settings-select,
  .settings-slider {
    width: 100%;
    margin-left: 0;
  }
  
  .header-content {
    flex-direction: column;
    align-items: flex-start;
  }
  
  .settings-header h1 {
    margin-top: 15px;
    margin-left: 0;
  }
}
</style> 