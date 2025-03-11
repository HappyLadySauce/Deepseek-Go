import { defineStore } from 'pinia';
import { ref, computed } from 'vue';
import type { 
  Session, 
  Message, 
  AIConfig, 
  KnowledgeFile
} from '@/utils/chatApi';
import { 
  getSessions, 
  getSessionMessages, 
  updateSession, 
  deleteSession, 
  sendMessage,
  getAIConfigs,
  getDefaultAIConfig,
  getKnowledgeFiles,
  deleteKnowledgeFile,
  createChatStream
} from '@/utils/chatApi';
import { ElMessage } from 'element-plus';

export const useChatStore = defineStore('chat', () => {
  // 状态
  const sessions = ref<Session[]>([]);
  const currentSessionId = ref<number | null>(null);
  const messages = ref<Message[]>([]);
  const loading = ref(false);
  const streaming = ref(false);
  const streamController = ref<AbortController | null>(null);
  const aiConfigs = ref<AIConfig[]>([]);
  const currentConfigId = ref<number>(0); // 0表示使用默认配置
  const knowledgeFiles = ref<KnowledgeFile[]>([]);
  const selectedKnowledgeIds = ref<number[]>([]);
  const newMessage = ref('');
  const totalSessions = ref(0);
  const currentPage = ref(1);
  const pageSize = ref(10);

  // 计算属性
  const currentSession = computed(() => 
    sessions.value.find(s => s.id === currentSessionId.value) || null
  );

  const currentConfig = computed(() => 
    aiConfigs.value.find(c => c.id === currentConfigId.value) || 
    aiConfigs.value.find(c => c.is_default) || 
    null
  );

  // 加载会话列表
  async function loadSessions(page = 1, size = 10) {
    try {
      console.log('开始加载会话列表, 页码:', page, '每页大小:', size);
      const data = await getSessions(page, size);
      console.log('会话列表加载成功:', data);
      sessions.value = data.data;
      totalSessions.value = data.total;
      currentPage.value = page;
      pageSize.value = size;
    } catch (error) {
      console.error('加载会话列表失败:', error);
      ElMessage.error('加载会话列表失败');
    }
  }

  // 加载会话消息
  async function loadMessages(sessionId: number) {
    if (!sessionId) return;
    
    try {
      console.log('开始加载会话消息, 会话ID:', sessionId);
      loading.value = true;
      currentSessionId.value = sessionId;
      const data = await getSessionMessages(sessionId);
      console.log('会话消息加载成功:', data);
      messages.value = data.data;
    } catch (error) {
      console.error('加载会话消息失败:', error);
      ElMessage.error('加载会话消息失败');
    } finally {
      loading.value = false;
    }
  }

  // 创建新会话或切换会话
  async function selectSession(sessionId: number | null) {
    console.log('选择会话, ID:', sessionId);
    // 如果正在流式传输，先取消
    if (streaming.value && streamController.value) {
      console.log('取消当前流式传输');
      streamController.value.abort();
      streaming.value = false;
    }
    
    // 清空消息
    messages.value = [];
    
    // 如果是新会话，则不设置currentSessionId
    if (sessionId === 0 || sessionId === null) {
      console.log('创建新会话');
      currentSessionId.value = null;
      return;
    }
    
    // 加载选定会话的消息
    await loadMessages(sessionId);
  }

  // 更新会话标题
  async function renameSession(sessionId: number, title: string) {
    try {
      const updated = await updateSession(sessionId, title);
      const index = sessions.value.findIndex(s => s.id === sessionId);
      if (index !== -1) {
        sessions.value[index] = updated;
      }
      ElMessage.success('会话已重命名');
    } catch (error) {
      console.error('重命名会话失败:', error);
      ElMessage.error('重命名会话失败');
    }
  }

  // 删除会话
  async function removeSession(sessionId: number) {
    try {
      await deleteSession(sessionId);
      sessions.value = sessions.value.filter(s => s.id !== sessionId);
      if (currentSessionId.value === sessionId) {
        currentSessionId.value = null;
        messages.value = [];
      }
      ElMessage.success('会话已删除');
    } catch (error) {
      console.error('删除会话失败:', error);
      ElMessage.error('删除会话失败');
    }
  }

  // 发送消息
  async function submitMessage(message = newMessage.value) {
    if (!message.trim()) return;
    
    try {
      console.log('发送消息:', message, '会话ID:', currentSessionId.value || 0);
      // 添加用户消息到列表
      const userMessage: Message = {
        id: Date.now(), // 临时ID
        role: 'user',
        content: message,
        created_at: new Date().toISOString()
      };
      
      messages.value.push(userMessage);
      newMessage.value = '';
      loading.value = true;
      
      // 发送消息到服务器
      const response = await sendMessage(
        message, 
        currentSessionId.value || 0,
        currentConfigId.value,
        selectedKnowledgeIds.value
      );
      
      console.log('收到服务器响应:', response);
      
      // 更新消息列表
      messages.value.push(response.message);
      
      // 如果是新会话，更新会话列表
      if (!currentSessionId.value) {
        currentSessionId.value = response.session.id;
        sessions.value = [response.session, ...sessions.value];
      } else {
        // 更新现有会话的最后消息
        const index = sessions.value.findIndex(s => s.id === currentSessionId.value);
        if (index !== -1) {
          sessions.value[index] = response.session;
        }
      }
    } catch (error) {
      console.error('发送消息失败:', error);
      ElMessage.error('发送消息失败');
    } finally {
      loading.value = false;
    }
  }

  // 流式发送消息
  function streamMessage(message = newMessage.value) {
    if (!message.trim() || streaming.value) return;
    
    try {
      console.log('开始流式发送消息:', message, '会话ID:', currentSessionId.value || 0);
      // 添加用户消息到列表
      const userMessage: Message = {
        id: Date.now(), // 临时ID
        role: 'user',
        content: message,
        created_at: new Date().toISOString()
      };
      
      messages.value.push(userMessage);
      newMessage.value = '';
      
      // 准备AI回复消息
      const aiMessage: Message = {
        id: Date.now() + 1, // 临时ID
        role: 'assistant',
        content: '',
        created_at: new Date().toISOString()
      };
      
      messages.value.push(aiMessage);
      streaming.value = true;
      
      // 创建流式连接
      const { eventSource, controller } = createChatStream(
        message,
        currentSessionId.value || 0,
        currentConfigId.value,
        selectedKnowledgeIds.value
      );
      
      streamController.value = controller;
      
      // 处理流式数据
      eventSource.onmessage = (event) => {
        console.log('流式数据:', event.data);
        if (event.data === '[DONE]') {
          // 流结束
          console.log('流式传输完成');
          streaming.value = false;
          eventSource.close();
          streamController.value = null;
          
          // 更新会话列表（这里可能需要调用额外的API来获取更新后的会话信息）
          loadSessions();
          return;
        }
        
        try {
          const data = JSON.parse(event.data);
          if (data.delta && data.delta.content) {
            // 更新最后一条消息的内容
            aiMessage.content += data.delta.content;
            // 强制更新最后一条消息以触发UI更新
            messages.value = [...messages.value];
          }
        } catch (e) {
          console.error('解析流数据失败:', e);
        }
      };
      
      eventSource.onerror = (error) => {
        console.error('流式连接错误:', error);
        ElMessage.error('聊天连接中断');
        streaming.value = false;
        eventSource.close();
        streamController.value = null;
      };
    } catch (error) {
      console.error('初始化流式聊天失败:', error);
      ElMessage.error('初始化聊天失败');
      streaming.value = false;
    }
  }

  // 取消流式传输
  function cancelStream() {
    if (streamController.value) {
      streamController.value.abort();
      streaming.value = false;
      streamController.value = null;
    }
  }

  // 加载AI配置
  async function loadAIConfigs() {
    try {
      aiConfigs.value = await getAIConfigs();
      // 如果没有设置当前配置，则使用默认配置
      if (!currentConfigId.value) {
        const defaultConfig = aiConfigs.value.find(c => c.is_default);
        if (defaultConfig) {
          currentConfigId.value = defaultConfig.id;
        }
      }
    } catch (error) {
      console.error('加载AI配置失败:', error);
    }
  }

  // 设置当前AI配置
  function setCurrentConfig(configId: number) {
    currentConfigId.value = configId;
  }

  // 加载知识库文件
  async function loadKnowledgeFiles() {
    try {
      const data = await getKnowledgeFiles();
      knowledgeFiles.value = data.data;
    } catch (error) {
      console.error('加载知识库文件失败:', error);
    }
  }

  // 选择/取消选择知识库文件
  function toggleKnowledgeFile(id: number) {
    const index = selectedKnowledgeIds.value.indexOf(id);
    if (index === -1) {
      selectedKnowledgeIds.value.push(id);
    } else {
      selectedKnowledgeIds.value.splice(index, 1);
    }
  }

  // 删除知识库文件
  async function removeKnowledgeFile(id: number) {
    try {
      await deleteKnowledgeFile(id);
      knowledgeFiles.value = knowledgeFiles.value.filter(f => f.id !== id);
      // 如果被删除的文件在已选择列表中，也从那里移除
      const index = selectedKnowledgeIds.value.indexOf(id);
      if (index !== -1) {
        selectedKnowledgeIds.value.splice(index, 1);
      }
      ElMessage.success('文件已删除');
    } catch (error) {
      console.error('删除文件失败:', error);
      ElMessage.error('删除文件失败');
    }
  }

  // 清空选中的知识库
  function clearSelectedKnowledge() {
    selectedKnowledgeIds.value = [];
  }

  return {
    // 状态
    sessions,
    currentSessionId,
    messages,
    loading,
    streaming,
    aiConfigs,
    currentConfigId,
    knowledgeFiles,
    selectedKnowledgeIds,
    newMessage,
    totalSessions,
    currentPage,
    pageSize,
    
    // 计算属性
    currentSession,
    currentConfig,
    
    // 方法
    loadSessions,
    loadMessages,
    selectSession,
    renameSession,
    removeSession,
    submitMessage,
    streamMessage,
    cancelStream,
    loadAIConfigs,
    setCurrentConfig,
    loadKnowledgeFiles,
    toggleKnowledgeFile,
    removeKnowledgeFile,
    clearSelectedKnowledge
  };
}); 