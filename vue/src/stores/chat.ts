/**
 * 聊天功能状态管理模块
 * 
 * 该模块负责管理聊天相关的所有状态和操作，包括：
 * - 会话列表管理（加载、创建、重命名、删除）
 * - 消息管理（加载、发送、接收）
 * - AI配置管理（加载、选择）
 * - 知识库文件管理（加载、选择、删除）
 * - 流式聊天功能
 * 
 * 适配了后端返回的数据格式，将后端使用的大驼峰命名转换为前端使用的小写命名。
 */
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
  createChatStream,
  createSession
} from '@/utils/chatApi';
import axios from '@/axios';
import { ElMessage } from 'element-plus';

export const useChatStore = defineStore('chat', () => {
  // ==================== 状态定义 ====================
  
  /**
   * 所有聊天会话列表
   * @type {Session[]} 会话对象数组，每个会话包含id、标题、最后消息等信息
   */
  const sessions = ref<Session[]>([]);
  
  /**
   * 当前选中的会话ID
   * @type {number|null} 当前会话ID，如果为null表示没有选中会话
   */
  const currentSessionId = ref<number | null>(null);
  
  /**
   * 当前会话的消息列表
   * @type {Message[]} 消息对象数组，每条消息包含id、角色、内容、创建时间等
   */
  const messages = ref<Message[]>([]);
  
  /**
   * 是否正在加载数据
   * @type {boolean} true表示正在加载，false表示加载完成
   */
  const loading = ref(false);
  
  /**
   * 是否正在进行流式传输
   * @type {boolean} true表示正在流式传输中，false表示已结束
   */
  const streaming = ref(false);
  
  /**
   * 流式传输控制器，用于中断流式传输
   * @type {AbortController|null} 控制器对象或null
   */
  const streamController = ref<AbortController | null>(null);
  
  /**
   * 所有AI配置列表
   * @type {AIConfig[]} AI配置对象数组，包含模型名称、温度、最大token等
   */
  const aiConfigs = ref<AIConfig[]>([]);
  
  /**
   * 当前选中的AI配置ID
   * @type {number} 配置ID，0表示使用默认配置
   */
  const currentConfigId = ref<number>(0);
  
  /**
   * 所有知识库文件列表
   * @type {KnowledgeFile[]} 知识库文件对象数组，包含文件名、大小、状态等
   */
  const knowledgeFiles = ref<KnowledgeFile[]>([]);
  
  /**
   * 当前选中的知识库文件ID列表
   * @type {number[]} 知识库文件ID数组
   */
  const selectedKnowledgeIds = ref<number[]>([]);
  
  /**
   * 新消息输入内容
   * @type {string} 用户正在输入的消息文本
   */
  const newMessage = ref('');
  
  /**
   * 会话总数
   * @type {number} 用户拥有的全部会话数量
   */
  const totalSessions = ref(0);
  
  /**
   * 当前页码
   * @type {number} 会话列表分页当前页码
   */
  const currentPage = ref(1);
  
  /**
   * 每页大小
   * @type {number} 会话列表分页每页显示数量
   */
  const pageSize = ref(10);

  // ==================== 计算属性 ====================
  
  /**
   * 当前选中的会话对象
   * 根据currentSessionId查找对应的会话，找不到则返回null
   * @returns {Session|null} 当前会话对象或null
   */
  const currentSession = computed(() => 
    sessions.value.find(s => s.id === currentSessionId.value) || null
  );

  /**
   * 当前选中的AI配置对象
   * 首先查找指定ID的配置，找不到则查找默认配置，仍找不到则返回null
   * @returns {AIConfig|null} 当前AI配置对象或null
   */
  const currentConfig = computed(() => 
    aiConfigs.value.find(c => c.id === currentConfigId.value) || 
    aiConfigs.value.find(c => c.is_default) || 
    null
  );

  // ==================== 会话管理函数 ====================
  
  /**
   * 加载会话列表
   * 从服务器获取会话列表数据，并处理后端不同的数据格式
   * 
   * @param {number} page - 页码，默认为1
   * @param {number} size - 每页大小，默认为10
   * @returns {Promise<void>}
   */
  async function loadSessions(page = 1, size = 10) {
    try {
      console.log('开始加载会话列表, 页码:', page, '每页大小:', size);
      // 直接使用axios获取数据
      const response = await axios.get(`/chat/sessions?page=${page}&page_size=${size}`);
      const result = response.data;
      console.log('会话列表加载成功:', result);
      
      // 处理后端返回的会话数据，适配字段名称
      if (result && result.sessions && Array.isArray(result.sessions)) {
        // 将后端大驼峰格式转换为前端所需格式
        const formattedSessions = result.sessions.map((session: any) => ({
          id: session.ID,                  // 会话ID
          title: session.title,            // 会话标题
          last_message: session.last_message, // 最后一条消息内容
          created_at: session.CreatedAt,   // 创建时间
          updated_at: session.UpdatedAt    // 更新时间
        }));
        
        sessions.value = formattedSessions;
        console.log('设置会话列表:', sessions.value);
        totalSessions.value = result.total || 0;
        currentPage.value = page;
        pageSize.value = size;
      } else if (result && result.data && Array.isArray(result.data)) {
        // 兼容处理老格式数据
        sessions.value = result.data;
        console.log('设置会话列表(使用data字段):', sessions.value);
        totalSessions.value = result.total || 0;
        currentPage.value = page;
        pageSize.value = size;
      } else {
        console.error('API返回的会话列表数据格式不正确:', result);
        ElMessage.warning('会话列表数据格式不正确');
      }
    } catch (error) {
      console.error('加载会话列表失败:', error);
      ElMessage.error('加载会话列表失败');
    }
  }

  /**
   * 加载会话消息
   * 根据会话ID从服务器获取特定会话的消息列表
   * 
   * @param {number} sessionId - 会话ID
   * @returns {Promise<void>}
   */
  async function loadMessages(sessionId: number) {
    if (!sessionId) return;
    
    try {
      console.log('开始加载会话消息, 会话ID:', sessionId);
      loading.value = true;
      currentSessionId.value = sessionId;
      
      // 直接用axios获取消息
      const response = await axios.get(`/chat/sessions/${sessionId}`);
      const result = response.data;
      console.log('会话消息加载成功:', result);
      
      // 检查消息格式，适配不同返回结构
      if (result && result.messages && Array.isArray(result.messages)) {
        // 将后端格式转换为前端需要的格式
        const formattedMessages = result.messages.map((msg: any) => ({
          id: msg.ID,              // 消息ID
          role: msg.role,          // 角色(user/assistant)
          content: msg.content,    // 消息内容
          created_at: msg.CreatedAt // 创建时间
        }));
        messages.value = formattedMessages;
      } else if (result && result.data && Array.isArray(result.data)) {
        // 处理老格式数据
        messages.value = result.data;
      } else {
        console.error('API返回的消息数据格式不正确:', result);
        messages.value = [];
        ElMessage.warning('消息数据格式不正确');
      }
    } catch (error) {
      console.error('加载会话消息失败:', error);
      ElMessage.error('加载会话消息失败');
      messages.value = [];
    } finally {
      loading.value = false;
    }
  }

  /**
   * 创建新会话或切换会话
   * 如果sessionId为null，则准备创建新会话；否则切换到指定会话
   * 
   * @param {number|null} sessionId - 会话ID，null表示准备创建新会话
   * @returns {Promise<void>}
   */
  async function selectSession(sessionId: number | null) {
    console.log('选择会话, ID:', sessionId);
    
    // 如果正在流式传输，先取消
    if (streaming.value && streamController.value) {
      console.log('取消当前流式传输');
      streamController.value.abort();
      streaming.value = false;
    }
    
    // 清空消息
    if (!messages.value) {
      messages.value = [];
    } else {
      messages.value = [];
    }
    
    // 如果是新会话，则不设置currentSessionId
    if (sessionId === 0 || sessionId === null) {
      console.log('准备创建新会话，当前设置为null');
      currentSessionId.value = null;
      return;
    } 
    
    // 如果是临时会话（前端创建的）
    if (sessionId && sessionId > 9999999) {
      console.log('选择临时会话:', sessionId);
      currentSessionId.value = sessionId;
      
      // 由于是临时会话，消息可能为空，不需要从服务器加载
      return;
    }
    
    // 处理正常会话ID
    console.log('切换到会话:', sessionId);
    currentSessionId.value = sessionId;
    
    // 加载选定会话的消息
    await loadMessages(sessionId);
  }

  /**
   * 更新会话标题
   * 向服务器发送请求修改指定会话的标题
   * 
   * @param {number} sessionId - 会话ID
   * @param {string} title - 新标题
   * @returns {Promise<void>}
   */
  async function renameSession(sessionId: number, title: string) {
    try {
      const response = await axios.put(`/chat/sessions/${sessionId}`, { title });
      const updatedSession = response.data.data;
      
      // 格式转换
      const formattedSession: Session = {
        id: updatedSession.ID,
        title: updatedSession.title,
        last_message: updatedSession.last_message,
        created_at: updatedSession.CreatedAt,
        updated_at: updatedSession.UpdatedAt
      };
      
      const index = sessions.value.findIndex(s => s.id === sessionId);
      if (index !== -1) {
        sessions.value[index] = formattedSession;
      }
      ElMessage.success('会话已重命名');
    } catch (error) {
      console.error('重命名会话失败:', error);
      ElMessage.error('重命名会话失败');
    }
  }

  /**
   * 删除会话
   * 向服务器发送请求删除指定会话及其所有消息
   * 
   * @param {number} sessionId - 会话ID
   * @returns {Promise<void>}
   */
  async function removeSession(sessionId: number) {
    try {
      await axios.delete(`/chat/sessions/${sessionId}`);
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

  // ==================== 消息发送函数 ====================
  
  /**
   * 发送消息
   * 向当前会话发送消息并接收AI回复
   * 如果当前没有会话，会先创建一个新会话
   * 
   * @param {string} message - 消息内容，默认使用newMessage的值
   * @returns {Promise<any>} 服务器响应
   */
  async function submitMessage(message = newMessage.value) {
    if (!message.trim()) return;
    
    try {
      // 检查是否有会话ID
      if (!currentSessionId.value) {
        console.log('没有会话ID，先创建新会话');
        await createNewSession();
      }
      
      // 确保会话ID存在
      if (!currentSessionId.value) {
        throw new Error('创建会话失败，无法获取会话ID');
      }
      
      console.log('发送消息:', message, '会话ID:', currentSessionId.value, '配置ID:', currentConfigId.value);
      
      // 添加用户消息到列表
      const userMessage: Message = {
        id: Date.now(), // 临时ID
        role: 'user',
        content: message,
        created_at: new Date().toISOString()
      };
      
      if (!messages.value) {
        messages.value = [];
      }
      
      messages.value.push(userMessage);
      newMessage.value = '';
      loading.value = true;
      
      // 使用axios直接发送消息到服务器
      const response = await axios.post('/chat/completions', {
        session_id: currentSessionId.value,
        message: message,
        ai_config_id: currentConfigId.value,
        knowledge_ids: selectedKnowledgeIds.value
      });
      
      const result = response.data;
      console.log('收到服务器响应:', result);
      
      // 处理后端返回的消息
      if (result.data && typeof result.data === 'object') {
        // 统一处理返回的消息格式
        const assistantMessage: Message = {
          id: result.data.ID || Date.now() + 1,
          role: 'assistant',
          content: result.data.content || result.data.message || '',
          created_at: result.data.CreatedAt || new Date().toISOString()
        };
        
        messages.value.push(assistantMessage);
      }
      
      // 如果响应中包含会话信息，更新会话列表
      if (result.session) {
        // 将后端会话格式转换为前端格式
        const updatedSession: Session = {
          id: result.session.ID,
          title: result.session.title,
          last_message: result.session.last_message,
          created_at: result.session.CreatedAt,
          updated_at: result.session.UpdatedAt
        };
        
        // 查找并更新会话列表中的会话
        const index = sessions.value.findIndex(s => s.id === currentSessionId.value);
        if (index !== -1) {
          console.log('更新会话信息:', updatedSession);
          sessions.value[index] = updatedSession;
        } else {
          // 找不到会话，可能是API不一致，尝试刷新会话列表
          console.log('找不到会话，刷新会话列表');
          await loadSessions();
        }
      }
      
      return response.data;
    } catch (error: any) {
      console.error('发送消息失败:', error);
      if (error.response && error.response.status === 404) {
        ElMessage.error('会话不存在，请创建新会话');
      } else if (error.response && error.response.status === 401) {
        ElMessage.error('未授权操作，请重新登录');
      } else {
        ElMessage.error('发送消息失败: ' + (error.message || '未知错误'));
      }
      throw error;
    } finally {
      loading.value = false;
    }
  }

  /**
   * 流式发送消息
   * 使用SSE(Server-Sent Events)流式接收AI回复
   * 如果当前没有会话，会先创建一个新会话
   * 
   * @param {string} message - 消息内容，默认使用newMessage的值
   * @returns {Promise<void>}
   */
  async function streamMessage(message = newMessage.value) {
    if (!message.trim() || streaming.value) return;
    
    try {
      // 检查是否有会话ID
      if (!currentSessionId.value) {
        console.log('没有会话ID，先创建新会话');
        await createNewSession();
      }
      
      // 确保会话ID存在
      if (!currentSessionId.value) {
        throw new Error('创建会话失败，无法获取会话ID');
      }
      
      console.log('开始流式发送消息:', message, '会话ID:', currentSessionId.value);
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
        currentSessionId.value,
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

  /**
   * 取消流式传输
   * 中断当前的流式传输连接
   */
  function cancelStream() {
    if (streamController.value) {
      streamController.value.abort();
      streaming.value = false;
      streamController.value = null;
    }
  }

  // ==================== AI配置函数 ====================
  
  /**
   * 加载AI配置
   * 从服务器获取所有AI配置项，并设置默认配置
   * 
   * @returns {Promise<void>}
   */
  async function loadAIConfigs() {
    try {
      console.log('开始加载AI配置');
      const response = await axios.get('/ai-config/');
      const result = response.data;
      console.log('AI配置加载成功:', result);
      
      // 处理后端返回的AI配置数据，适配字段名称
      if (result.data && Array.isArray(result.data)) {
        // 将后端格式转换为前端格式
        const formattedConfigs = result.data.map((config: any) => ({
          id: config.ID || config.id,
          model_name: config.model_name,
          temperature: config.temperature,
          max_tokens: config.max_tokens,
          provider: config.provider,
          is_default: config.is_default,
          created_at: config.CreatedAt || config.created_at,
          updated_at: config.UpdatedAt || config.updated_at
        }));
        
        aiConfigs.value = formattedConfigs;
        
        // 如果没有设置当前配置，则使用默认配置
        if (!currentConfigId.value) {
          const defaultConfig = aiConfigs.value.find(c => c.is_default);
          if (defaultConfig) {
            currentConfigId.value = defaultConfig.id;
          }
        }
      } else {
        console.error('API返回的AI配置数据格式不正确:', result);
        ElMessage.warning('AI配置数据格式不正确');
      }
    } catch (error) {
      console.error('加载AI配置失败:', error);
      ElMessage.error('加载AI配置失败');
    }
  }

  /**
   * 设置当前AI配置
   * 更新当前使用的AI配置ID
   * 
   * @param {number} configId - 配置ID
   */
  function setCurrentConfig(configId: number) {
    currentConfigId.value = configId;
  }

  // ==================== 知识库文件函数 ====================
  
  /**
   * 加载知识库文件
   * 从服务器获取知识库文件列表
   * 
   * @param {number} page - 页码，默认为1
   * @param {number} size - 每页大小，默认为10
   * @returns {Promise<void>}
   */
  async function loadKnowledgeFiles(page = 1, size = 10) {
    try {
      console.log('开始加载知识库文件列表');
      const response = await axios.get(`/knowledge/files?page=${page}&page_size=${size}`);
      const result = response.data;
      
      if (result.files && Array.isArray(result.files)) {
        // 将后端格式转换为前端格式
        const formattedFiles = result.files.map((file: any) => ({
          id: file.ID,
          file_name: file.file_name,
          file_size: file.file_size,
          file_type: file.file_type,
          status: file.status,
          created_at: file.CreatedAt,
          processed_at: file.processed_at,
          vector_count: file.vector_count
        }));
        
        knowledgeFiles.value = formattedFiles;
        console.log('知识库文件加载成功:', knowledgeFiles.value);
      } else if (result.data && Array.isArray(result.data)) {
        // 兼容老格式
        knowledgeFiles.value = result.data;
        console.log('知识库文件加载成功(使用data字段):', knowledgeFiles.value);
      } else {
        console.error('API返回的知识库文件数据格式不正确:', result);
        ElMessage.warning('知识库文件数据格式不正确');
      }
    } catch (error) {
      console.error('加载知识库文件失败:', error);
      ElMessage.error('加载知识库文件失败');
    }
  }

  /**
   * 切换知识库文件选择状态
   * 在已选中知识库文件列表中添加或移除指定ID
   * 
   * @param {number} id - 知识库文件ID
   */
  function toggleKnowledgeFile(id: number) {
    const index = selectedKnowledgeIds.value.indexOf(id);
    if (index === -1) {
      selectedKnowledgeIds.value.push(id);
    } else {
      selectedKnowledgeIds.value.splice(index, 1);
    }
  }

  /**
   * 删除知识库文件
   * 向服务器发送请求删除指定知识库文件
   * 
   * @param {number} id - 知识库文件ID
   * @returns {Promise<void>}
   */
  async function removeKnowledgeFile(id: number) {
    try {
      await axios.delete(`/knowledge/files/${id}`);
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

  /**
   * 清空选中的知识库文件
   * 清空当前选中的知识库文件列表
   */
  function clearSelectedKnowledge() {
    selectedKnowledgeIds.value = [];
  }
  
  // ==================== 创建新会话函数 ====================
  
  /**
   * 创建新会话
   * 向服务器发送请求创建新会话
   * 成功后会自动添加一条欢迎消息，并更新URL
   * 
   * @returns {Promise<Session>} 创建成功的会话对象
   */
  async function createNewSession() {
    try {
      console.log('开始创建新会话');
      
      // 直接调用axios创建新会话
      const response = await axios.post('/chat/sessions', { title: "新会话" });
      const backendSession = response.data.data;
      console.log('新会话创建成功:', backendSession);
      
      // 添加到会话列表
      if (!sessions.value) {
        sessions.value = [];
      }
      
      // 将API返回格式转换为前端格式
      const newSession = {
        id: backendSession.ID,
        title: backendSession.title,
        last_message: backendSession.last_message || "",
        created_at: backendSession.CreatedAt,
        updated_at: backendSession.UpdatedAt
      };
      
      // 将新会话添加到会话列表的最前面
      sessions.value = [newSession, ...sessions.value];
      
      // 设置当前会话ID
      currentSessionId.value = newSession.id;
      
      // 清空消息列表
      messages.value = [];
      
      // 如果有欢迎消息，可以添加
      messages.value.push({
        id: Date.now(),
        role: 'assistant',
        content: '你好！我是DeepSeek AI助手。请问有什么我可以帮助你的？',
        created_at: new Date().toISOString()
      });
      
      // 更新URL
      if (window.history && window.location) {
        const url = new URL(window.location.href);
        url.searchParams.set('sessionId', newSession.id.toString());
        window.history.pushState({}, '', url.toString());
      }
      
      return newSession;
    } catch (error) {
      console.error('创建新会话失败:', error);
      ElMessage.error('创建新会话失败');
      throw error;
    }
  }

  return {
    // 状态
    sessions,
    messages,
    currentSessionId, 
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
    clearSelectedKnowledge,
    createNewSession
  };
}); 