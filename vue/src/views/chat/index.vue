<script setup lang="ts">
import { ref, onMounted, nextTick, onUnmounted, watch } from 'vue';
import { useRoute } from 'vue-router';
import { useChatStore } from '@/stores/chat';
import { uploadKnowledgeFile, type KnowledgeFile } from '@/utils/chatApi';
import { ElMessage, ElMessageBox } from 'element-plus';
import { 
  ChatLineSquare as Send, 
  Plus, 
  Setting, 
  EditPen as Edit, 
  Delete, 
  Document, 
  Upload, 
  CloseBold as Close 
} from '@element-plus/icons-vue';

// 使用Pinia状态管理
const chatStore = useChatStore();
const route = useRoute();

// 页面状态
const chatContainerRef = ref<HTMLElement | null>(null);
const inputRef = ref<HTMLTextAreaElement | null>(null);
const showAISettings = ref(false);
const fileUploading = ref(false);
const fileList = ref<any[]>([]);
const sessionInput = ref('');
const modelOptions = ref<{[key: string]: Array<{name: string, provider: string, description: string}>}>({
  'deepseek': [
    { name: "deepseek-chat", provider: "deepseek", description: "基础模型"},
    { name: "deepseek-reasoner", provider: "deepseek", description: "深度思考模型"}
  ],
  'kimi': [
    { name: 'moonshot-v1-8k', provider: 'kimi', description: '基础模型，支持8K上下文' },
    { name: 'moonshot-v1-32k', provider: 'kimi', description: '基础模型，支持32K上下文' },
    { name: 'moonshot-v1-128k', provider: 'kimi', description: '基础模型，支持128K上下文' },
    { name: 'moonshot-v1-auto', provider: 'kimi', description: '自动选择模型，根据上下文长度' }
  ]
});

// 自定义配置表单
const configForm = ref({
  model_name: '',
  provider: '',
  temperature: 0.7,
  max_tokens: 2048,
  is_default: false
});

// 监听当前配置变化，同步到配置表单
watch(() => chatStore.currentConfig, (newConfig) => {
  if (newConfig) {
    configForm.value.model_name = newConfig.model_name;
    configForm.value.provider = newConfig.provider;
    configForm.value.temperature = newConfig.temperature;
    configForm.value.max_tokens = newConfig.max_tokens;
    configForm.value.is_default = newConfig.is_default;
  }
}, { immediate: true, deep: true });

// 保存AI配置
async function saveAIConfig() {
  try {
    // 这里可以调用API保存配置
    // 暂时只设置当前配置
    if (chatStore.currentConfigId && chatStore.currentConfigId > 0) {
      // 更新现有配置
      console.log('更新AI配置:', configForm.value);
    } else {
      // 创建新配置
      console.log('创建新AI配置:', configForm.value);
    }
    
    // 更新store中的AI配置
    if (chatStore.currentConfig) {
      chatStore.currentConfig.model_name = configForm.value.model_name;
      chatStore.currentConfig.provider = configForm.value.provider;
      chatStore.currentConfig.temperature = configForm.value.temperature;
      chatStore.currentConfig.max_tokens = configForm.value.max_tokens;
      chatStore.currentConfig.is_default = configForm.value.is_default;
    }
    
    ElMessage.success('AI配置已保存');
    showAISettings.value = false;
  } catch (error) {
    console.error('保存AI配置失败:', error);
    ElMessage.error('保存AI配置失败');
  }
}

// 重置配置表单
function resetConfigForm() {
  if (chatStore.currentConfig) {
    configForm.value.model_name = chatStore.currentConfig.model_name;
    configForm.value.provider = chatStore.currentConfig.provider;
    configForm.value.temperature = chatStore.currentConfig.temperature;
    configForm.value.max_tokens = chatStore.currentConfig.max_tokens;
    configForm.value.is_default = chatStore.currentConfig.is_default;
  } else {
    configForm.value = {
      model_name: '',
      provider: '',
      temperature: 0.7,
      max_tokens: 2048,
      is_default: false
    };
  }
}

// 监听消息变化，自动滚动到底部
watch(() => chatStore.messages, () => {
  scrollToBottom();
}, { deep: true });

// 初始化页面
onMounted(async () => {
  try {
    console.log('正在加载聊天页面数据...');
    // 加载数据
    const results = await Promise.allSettled([
      chatStore.loadAIConfigs().catch(e => {
        console.error('加载AI配置失败:', e);
        return null;
      }),
      chatStore.loadSessions().catch(e => {
        console.error('加载会话列表失败:', e);
        return null;
      }),
      chatStore.loadKnowledgeFiles().catch(e => {
        console.error('加载知识库文件失败:', e);
        return null;
      })
    ]);
    
    console.log('数据加载结果:', results);
    
    // 检查会话加载状态
    const sessionsResult = results[1];
    if (sessionsResult.status === 'fulfilled') {
      console.log('会话加载成功');
    } else {
      console.error('会话加载失败，使用模拟数据:', sessionsResult.reason);
      // 使用模拟数据
      if (!chatStore.sessions || chatStore.sessions.length === 0) {
        console.log('使用模拟会话数据');
        chatStore.sessions = [
          {
            id: 1,
            title: "关于AI的讨论",
            last_message: "人工智能正在改变我们的生活方式...",
            created_at: new Date().toISOString(),
            updated_at: new Date().toISOString()
          },
          {
            id: 2,
            title: "技术问题解答",
            last_message: "你可以使用Vue3的Composition API来解决这个问题...",
            created_at: new Date(Date.now() - 24 * 60 * 60 * 1000).toISOString(),
            updated_at: new Date(Date.now() - 24 * 60 * 60 * 1000).toISOString()
          }
        ];
      }
    }
    
    // 检查知识库文件加载状态
    const filesResult = results[2];
    if (filesResult.status !== 'fulfilled' && (!chatStore.knowledgeFiles || chatStore.knowledgeFiles.length === 0)) {
      console.log('使用模拟知识库文件数据');
      chatStore.knowledgeFiles = [
        {
          id: 1,
          file_name: "Vue3开发手册.pdf",
          file_size: 1254789,
          file_type: ".pdf",
          status: "completed",
          created_at: new Date().toISOString()
        },
        {
          id: 2,
          file_name: "深度学习入门.docx",
          file_size: 852147,
          file_type: ".docx",
          status: "completed",
          created_at: new Date(Date.now() - 2 * 24 * 60 * 60 * 1000).toISOString()
        }
      ];
    }
    
    console.log('数据加载完成');
  } catch (error) {
    console.error("加载数据失败:", error);
  }
  
  // 如果URL有会话ID参数，加载对应会话
  const sessionId = route.query.sessionId ? Number(route.query.sessionId) : null;
  if (sessionId) {
    console.log('从URL加载会话ID:', sessionId);
    chatStore.selectSession(sessionId);
  } else {
    console.log('没有会话ID参数，不加载特定会话');
  }
  
  // 设置消息观察器
  setupMessageObserver();
  
  // 聚焦输入框
  nextTick(() => {
    if (inputRef.value) {
      inputRef.value.focus();
    }
  });
});

// 滚动到底部
function scrollToBottom() {
  nextTick(() => {
    if (chatContainerRef.value) {
      chatContainerRef.value.scrollTop = chatContainerRef.value.scrollHeight;
    }
  });
}

// 使用MutationObserver监听消息列表变化，自动滚动
let messageObserver: MutationObserver | null = null;

// 设置消息容器的观察器
function setupMessageObserver() {
  nextTick(() => {
    if (chatContainerRef.value && !messageObserver) {
      messageObserver = new MutationObserver(scrollToBottom);
      messageObserver.observe(chatContainerRef.value, {
        childList: true,
        subtree: true
      });
    }
  });
}

// 在组件卸载时清理观察器
onUnmounted(() => {
  if (messageObserver) {
    messageObserver.disconnect();
    messageObserver = null;
  }
});

// 监听消息输入框高度
function adjustInputHeight(event: Event) {
  const textarea = event.target as HTMLTextAreaElement;
  if (!textarea) return;
  
  textarea.style.height = 'auto';
  textarea.style.height = `${textarea.scrollHeight}px`;
}

// 处理文件上传
async function handleFileUpload(file: File) {
  // 检查文件类型和大小
  const allowedTypes = ['.pdf', '.docx', '.doc', '.txt', '.md'];
  const fileExt = file.name.substring(file.name.lastIndexOf('.')).toLowerCase();
  
  if (!allowedTypes.includes(fileExt)) {
    ElMessage.error('仅支持PDF、Word、TXT和Markdown文件');
    return false;
  }
  
  if (file.size > 10 * 1024 * 1024) { // 10MB限制
    ElMessage.error('文件大小不能超过10MB');
    return false;
  }
  
  try {
    fileUploading.value = true;
    const result = await uploadKnowledgeFile(file);
    ElMessage.success('文件上传成功，正在处理中');
    
    // 刷新知识库列表
    await chatStore.loadKnowledgeFiles();
    
    // 清空上传列表
    fileList.value = [];
    
    return false; // 阻止默认上传行为
  } catch (error) {
    console.error('文件上传失败:', error);
    ElMessage.error('文件上传失败');
    return false;
  } finally {
    fileUploading.value = false;
  }
}

// 删除知识库文件
function handleDeleteFile(id: number) {
  ElMessageBox.confirm(
    '确定要删除这个文件吗？这个操作不可恢复。',
    '删除确认',
    {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    }
  ).then(() => {
    chatStore.removeKnowledgeFile(id);
  }).catch(() => {
    // 用户取消删除
  });
}

// 重命名会话
function handleRenameSession(id: number) {
  const session = chatStore.sessions.find(s => s.id === id);
  if (session) {
    sessionInput.value = session.title;
    ElMessageBox.prompt(
      '请输入新的会话名称',
      '重命名会话',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        inputValue: session.title,
        inputValidator: (value) => {
          return value.trim().length > 0 ? true : '会话名称不能为空';
        }
      }
    ).then(({ value }) => {
      chatStore.renameSession(id, value);
    }).catch(() => {
      // 用户取消重命名
    });
  }
}

// 删除会话
function handleDeleteSession(id: number) {
  ElMessageBox.confirm(
    '确定要删除这个会话吗？所有相关消息都将被删除，这个操作不可恢复。',
    '删除确认',
    {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    }
  ).then(() => {
    chatStore.removeSession(id);
  }).catch(() => {
    // 用户取消删除
  });
}

// 发送消息
async function sendMessage() {
  console.log('发送消息按钮被点击');
  
  // 如果正在流式生成，则先停止
  if (chatStore.streaming) {
    console.log('停止生成');
    chatStore.cancelStream();
    return;
  }
  
  if (!chatStore.newMessage.trim()) {
    console.log('消息为空，不发送');
    return;
  }
  
  // 检查是否有会话ID
  if (!chatStore.currentSessionId) {
    console.log('没有会话ID，先创建新会话');
    try {
      await chatStore.createNewSession();
    } catch (error) {
      console.error('创建新会话失败:', error);
      ElMessage.error('创建新会话失败，无法发送消息');
      return;
    }
  }
  
  // 如果配置表单已修改，先应用配置
  if (chatStore.currentConfig && (
    chatStore.currentConfig.model_name !== configForm.value.model_name ||
    chatStore.currentConfig.provider !== configForm.value.provider ||
    chatStore.currentConfig.temperature !== configForm.value.temperature ||
    chatStore.currentConfig.max_tokens !== configForm.value.max_tokens
  )) {
    console.log('应用临时AI配置');
    // 应用临时配置
    chatStore.currentConfig.model_name = configForm.value.model_name;
    chatStore.currentConfig.provider = configForm.value.provider;
    chatStore.currentConfig.temperature = configForm.value.temperature;
    chatStore.currentConfig.max_tokens = configForm.value.max_tokens;
  }
  
  // 保存消息内容
  const messageText = chatStore.newMessage;
  console.log('开始发送消息:', messageText, '会话ID:', chatStore.currentSessionId);
  
  // 立即清空输入框
  chatStore.newMessage = '';

  try {
    // 调用API发送消息 - 让store处理所有消息添加逻辑
    await chatStore.submitMessage(messageText);
    console.log('消息发送成功');
    // 滚动到底部
    scrollToBottom();
  } catch (error) {
    console.error('消息发送失败:', error);
    ElMessage.error('消息发送失败，请稍后重试');
  }
}

// 发送消息（支持Enter键发送，Shift+Enter换行）
function handleKeyDown(event: KeyboardEvent) {
  if (event.key === 'Enter' && !event.shiftKey) {
    event.preventDefault();
    sendMessage();
  }
}

// 添加新的函数用于创建新会话
async function createNewSession() {
  console.log('创建新会话');
  
  try {
    // 先清空当前会话ID和消息
    chatStore.selectSession(null);
    
    // 获取默认AI配置（确保有配置ID）
    if (!chatStore.currentConfig) {
      await chatStore.loadAIConfigs();
    }
    
    // 发送创建新会话的API请求
    await chatStore.createNewSession();
    
    // 滚动到底部
    nextTick(() => {
      scrollToBottom();
    });
    
    // 聚焦输入框
    nextTick(() => {
      if (inputRef.value) {
        inputRef.value.focus();
      }
    });
  } catch (error) {
    console.error('创建新会话失败:', error);
    ElMessage.error('创建新会话失败');
  }
}
</script>

<template>
  <div class="chat-page">
    <!-- 左侧面板：历史记录 + AI配置 -->
    <div class="left-panel">
      <!-- AI参数设置区域 -->
      <div class="ai-panel">
        <div class="panel-header">
          <h2>AI参数设置</h2>
          <el-button type="primary" text :icon="Setting" @click="showAISettings = !showAISettings">
            {{ showAISettings ? '收起' : '展开' }}
          </el-button>
        </div>
        
        <div v-if="showAISettings" class="ai-settings">
          <el-form :model="configForm" label-width="80px" size="small">
            <el-form-item label="模型提供商">
              <el-select 
                v-model="configForm.provider" 
                placeholder="选择提供商"
                @change="configForm.model_name = ''"
              >
                <el-option label="DeepSeek" value="deepseek" />
                <el-option label="Kimi" value="kimi" />
              </el-select>
            </el-form-item>
            
            <el-form-item label="模型">
              <el-select v-model="configForm.model_name" placeholder="选择模型">
                <el-option-group 
                  v-for="(models, provider) in modelOptions" 
                  :key="provider" 
                  :label="provider === 'deepseek' ? 'DeepSeek' : 'Kimi'"
                  v-show="configForm.provider === '' || configForm.provider === provider"
                >
                  <el-option 
                    v-for="model in models" 
                    :key="model.name" 
                    :label="model.name" 
                    :value="model.name"
                  >
                    <div>
                      <div>{{ model.name }}</div>
                      <small class="text-muted">{{ model.description }}</small>
                    </div>
                  </el-option>
                </el-option-group>
              </el-select>
            </el-form-item>
            
            <el-form-item label="温度">
              <el-slider v-model="configForm.temperature" :min="0" :max="1" :step="0.1" show-input />
            </el-form-item>
            
            <el-form-item label="最大长度">
              <el-input-number v-model="configForm.max_tokens" :min="128" :max="8192" :step="128" />
            </el-form-item>
            
            <el-form-item>
              <el-button type="primary" size="small" @click="saveAIConfig">保存配置</el-button>
              <el-button size="small" @click="resetConfigForm">重置</el-button>
            </el-form-item>
          </el-form>
        </div>
      </div>
      
      <!-- 历史记录区域 -->
      <div class="history-panel">
        <div class="panel-header">
          <h2>历史会话</h2>
          <el-button type="primary" text :icon="Plus" @click="createNewSession">
            新会话
          </el-button>
        </div>
        
        <div class="session-list">
          <el-empty v-if="!chatStore.sessions || chatStore.sessions.length === 0" description="暂无历史会话" />
          
          <div 
            v-for="session in chatStore.sessions || []" 
            :key="session.id"
            class="session-item"
            :class="{ 'active': chatStore.currentSessionId === session.id }"
            @click="chatStore.selectSession(session.id)"
          >
            <div class="session-content">
              <div class="session-title">{{ session.title }}</div>
              <div class="session-preview">{{ session.last_message }}</div>
              <div class="session-time">{{ new Date(session.updated_at).toLocaleString() }}</div>
            </div>
            
            <div class="session-actions">
              <el-button 
                type="primary" 
                text 
                :icon="Edit" 
                @click.stop="handleRenameSession(session.id)" 
                title="重命名"
              />
              <el-button 
                type="danger" 
                text 
                :icon="Delete" 
                @click.stop="handleDeleteSession(session.id)" 
                title="删除"
              />
            </div>
          </div>
        </div>
      </div>
    </div>
    
    <!-- 中间面板：聊天窗口 -->
    <div class="chat-container">
      <!-- 消息列表 -->
      <div class="messages-container" ref="chatContainerRef">
        <el-empty 
          v-if="!chatStore.messages || chatStore.messages.length === 0" 
          description="开始新的对话" 
          :image-size="120"
        >
          <template #description>
            <p>有问题需要解答？向AI提问获取帮助。</p>
          </template>
        </el-empty>
        
        <div v-else class="message-list">
          <div 
            v-for="message in chatStore.messages || []" 
            :key="message.id"
            class="message-item"
            :class="message.role"
          >
            <div class="message-avatar">
              <el-avatar 
                :size="36" 
                :src="message.role === 'user' ? '' : '/src/assets/logo.png'" 
                :icon="message.role === 'user' ? 'UserFilled' : ''"
              />
            </div>
            
            <div class="message-content">
              <div class="message-text">{{ message.content }}</div>
              <div class="message-time">{{ new Date(message.created_at).toLocaleString() }}</div>
            </div>
          </div>
        </div>
      </div>
      
      <!-- 输入区域 -->
      <div class="input-container">
        <div class="selected-knowledge" v-if="chatStore.selectedKnowledgeIds && chatStore.selectedKnowledgeIds.length > 0">
          <div class="knowledge-tag" v-for="id in chatStore.selectedKnowledgeIds" :key="id">
            {{ chatStore.knowledgeFiles && chatStore.knowledgeFiles.find(f => f.id === id)?.file_name }}
            <el-button 
              type="danger" 
              text 
              :icon="Close" 
              size="small"
              @click="chatStore.toggleKnowledgeFile(id)"
            />
          </div>
          <el-button 
            type="info" 
            text 
            size="small"
            @click="chatStore.clearSelectedKnowledge"
          >
            清空
          </el-button>
        </div>
        
        <div class="input-area">
          <el-input
            v-model="chatStore.newMessage"
            type="textarea"
            :rows="1"
            ref="inputRef"
            placeholder="输入消息，按Enter发送，Shift+Enter换行"
            resize="none"
            @input="adjustInputHeight"
            @keydown="handleKeyDown"
          />
          
          <el-button
            :type="chatStore.streaming ? 'danger' : 'primary'"
            class="send-button"
            :icon="chatStore.streaming ? Close : Send"
            :loading="chatStore.loading"
            @click="sendMessage"
          >
            {{ chatStore.streaming ? '停止生成' : '发送' }}
          </el-button>
        </div>
      </div>
    </div>
    
    <!-- 右侧面板：知识库 -->
    <div class="knowledge-panel">
      <div class="panel-header">
        <h2>知识库</h2>
        <el-upload
          class="upload-button"
          action="#"
          :http-request="({ file }: any) => handleFileUpload(file)"
          :show-file-list="false"
          :before-upload="() => true"
          accept=".pdf,.docx,.doc,.txt,.md"
          :disabled="fileUploading"
        >
          <el-button 
            type="primary" 
            text 
            :icon="Upload" 
            :loading="fileUploading"
          >
            上传文件
          </el-button>
        </el-upload>
      </div>
      
      <div class="knowledge-list">
        <el-empty v-if="!chatStore.knowledgeFiles || chatStore.knowledgeFiles.length === 0" description="暂无知识库文件" />
        
        <div 
          v-for="file in chatStore.knowledgeFiles || []" 
          :key="file.id"
          class="knowledge-item"
          :class="{ 'selected': chatStore.selectedKnowledgeIds && chatStore.selectedKnowledgeIds.includes(file.id) }"
          @click="chatStore.toggleKnowledgeFile(file.id)"
        >
          <div class="file-info">
            <el-icon><Document /></el-icon>
            <div class="file-details">
              <div class="file-name">{{ file.file_name }}</div>
              <div class="file-meta">
                <span>{{ (file.file_size / 1024).toFixed(1) }}KB</span>
                <span>{{ file.status }}</span>
              </div>
            </div>
          </div>
          
          <div class="file-actions">
            <el-button 
              type="danger" 
              text 
              :icon="Delete" 
              @click.stop="handleDeleteFile(file.id)" 
              title="删除"
            />
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.chat-page {
  display: flex;
  height: calc(100vh - 130px);
  padding: 16px;
  gap: 16px;
  background-color: var(--page-bg);
}

.left-panel, .knowledge-panel {
  width: 260px;
  display: flex;
  flex-direction: column;
  gap: 16px;
  overflow: hidden;
}

.left-panel {
  display: flex;
  flex-direction: column;
}

.ai-panel {
  background-color: var(--card-bg);
  border-radius: 8px;
  overflow: hidden;
  box-shadow: var(--card-shadow);
}

.history-panel {
  flex: 1;
  background-color: var(--card-bg);
  border-radius: 8px;
  display: flex;
  flex-direction: column;
  overflow: hidden;
  box-shadow: var(--card-shadow);
}

.knowledge-panel {
  background-color: var(--card-bg);
  border-radius: 8px;
  overflow: hidden;
  box-shadow: var(--card-shadow);
}

.panel-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 12px 16px;
  border-bottom: 1px solid var(--border-color);
}

.panel-header h2 {
  font-size: 16px;
  margin: 0;
  color: var(--heading-color);
}

.ai-settings {
  padding: 16px;
}

.session-list, .knowledge-list {
  flex: 1;
  overflow-y: auto;
  padding: 8px;
}

.session-item, .knowledge-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 12px;
  margin-bottom: 8px;
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.3s;
  background-color: var(--item-bg);
}

.session-item:hover, .knowledge-item:hover {
  background-color: var(--hover-color);
}

.session-item.active {
  background-color: var(--active-color);
}

.session-content {
  flex: 1;
  overflow: hidden;
}

.session-title {
  font-weight: 500;
  margin-bottom: 4px;
  color: var(--text-color);
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.session-preview {
  font-size: 12px;
  color: var(--text-light);
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.session-time, .message-time {
  font-size: 11px;
  color: var(--text-muted);
  margin-top: 4px;
}

.session-actions, .file-actions {
  display: flex;
  gap: 4px;
  opacity: 0;
  transition: opacity 0.3s;
}

.session-item:hover .session-actions,
.knowledge-item:hover .file-actions {
  opacity: 1;
}

.chat-container {
  flex: 1;
  display: flex;
  flex-direction: column;
  background-color: var(--card-bg);
  border-radius: 8px;
  overflow: hidden;
  box-shadow: var(--card-shadow);
}

.messages-container {
  flex: 1;
  overflow-y: auto;
  padding: 16px;
  display: flex;
  flex-direction: column;
}

.message-list {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.message-item {
  display: flex;
  gap: 12px;
  max-width: 80%;
}

.message-item.user {
  align-self: flex-end;
  flex-direction: row-reverse;
}

.message-content {
  background-color: var(--message-bg);
  padding: 12px 16px;
  border-radius: 12px;
  box-shadow: var(--message-shadow);
}

.message-item.assistant .message-content {
  background-color: var(--assistant-bg);
  border-top-left-radius: 4px;
}

.message-item.user .message-content {
  background-color: var(--user-bg);
  border-top-right-radius: 4px;
}

.message-text {
  white-space: pre-wrap;
  word-break: break-word;
  color: var(--text-color);
}

.input-container {
  padding: 16px;
  border-top: 1px solid var(--border-color);
}

.input-area {
  display: flex;
  gap: 12px;
  align-items: flex-end;
}

.send-button {
  flex-shrink: 0;
}

.text-muted {
  color: var(--text-muted);
  font-size: 12px;
}

.file-info {
  display: flex;
  align-items: center;
  gap: 8px;
  flex: 1;
  overflow: hidden;
}

.file-details {
  flex: 1;
  overflow: hidden;
}

.file-name {
  font-weight: 500;
  margin-bottom: 4px;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.file-meta {
  display: flex;
  gap: 8px;
  font-size: 11px;
  color: var(--text-muted);
}

.knowledge-item.selected {
  background-color: var(--active-color);
  border: 1px solid var(--primary-color);
}

.selected-knowledge {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
  margin-bottom: 8px;
}

.knowledge-tag {
  display: flex;
  align-items: center;
  gap: 4px;
  padding: 4px 8px;
  background-color: var(--tag-bg);
  border-radius: 4px;
  font-size: 12px;
}

/* 自定义滚动条 */
::-webkit-scrollbar {
  width: 6px;
}

::-webkit-scrollbar-thumb {
  background-color: var(--scrollbar-color);
  border-radius: 10px;
}

::-webkit-scrollbar-track {
  background-color: var(--scrollbar-track);
}
</style> 