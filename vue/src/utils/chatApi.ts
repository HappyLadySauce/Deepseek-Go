import axios from '@/axios.ts';

// 会话类型
export interface Session {
  id: number;
  title: string;
  last_message: string;
  created_at: string;
  updated_at: string;
}

// 消息类型
export interface Message {
  id: number;
  role: 'user' | 'assistant';
  content: string;
  created_at: string;
}

// AI配置类型
export interface AIConfig {
  id: number;
  model_name: string;
  temperature: number;
  max_tokens: number;
  provider: string;
  is_default: boolean;
  created_at: string;
  updated_at?: string;
}

// 模型选项类型
export interface ModelOption {
  name: string;
  provider: string;
  description: string;
}

// 知识库文件类型
export interface KnowledgeFile {
  id: number;
  file_name: string;
  file_size: number;
  file_type: string;
  status: 'processing' | 'completed' | 'failed';
  created_at: string;
  processed_at?: string;
  vector_count?: number;
}

// 分页数据结构
export interface PaginatedData<T> {
  data: T[];
  total: number;
  page: number;
  page_size: number;
}

// 获取会话列表
export async function getSessions(page = 1, pageSize = 10): Promise<PaginatedData<Session>> {
  const response = await axios.get(`/chat/sessions?page=${page}&page_size=${pageSize}`);
  return response.data.data;
}

// 获取会话消息
export async function getSessionMessages(sessionId: number, page = 1, pageSize = 20): Promise<PaginatedData<Message>> {
  const response = await axios.get(`/chat/sessions/${sessionId}?page=${page}&page_size=${pageSize}`);
  return response.data.data;
}

// 更新会话标题
export async function updateSession(sessionId: number, title: string): Promise<Session> {
  const response = await axios.put(`/chat/sessions/${sessionId}`, { title });
  return response.data.data;
}

// 删除会话
export async function deleteSession(sessionId: number): Promise<void> {
  await axios.delete(`/chat/sessions/${sessionId}`);
}

// 发送消息（非流式）
export async function sendMessage(
  message: string, 
  sessionId = 0, 
  aiConfigId = 0, 
  knowledgeIds: number[] = []
): Promise<{ message: Message, session: Session }> {
  const response = await axios.post('/chat/completions', {
    session_id: sessionId,
    message,
    ai_config_id: aiConfigId,
    knowledge_ids: knowledgeIds
  });
  return {
    message: response.data.data,
    session: response.data.session
  };
}

// 获取AI配置列表
export async function getAIConfigs(): Promise<AIConfig[]> {
  const response = await axios.get('/ai-config/');
  return response.data.data;
}

// 获取默认AI配置
export async function getDefaultAIConfig(): Promise<AIConfig> {
  const response = await axios.get('/ai-config/default');
  return response.data.data;
}

// 创建AI配置
export async function createAIConfig(config: Omit<AIConfig, 'id' | 'created_at'>): Promise<AIConfig> {
  const response = await axios.post('/ai-config/', config);
  return response.data.data;
}

// 更新AI配置
export async function updateAIConfig(id: number, config: Partial<AIConfig>): Promise<AIConfig> {
  const response = await axios.put(`/ai-config/${id}`, config);
  return response.data.data;
}

// 删除AI配置
export async function deleteAIConfig(id: number): Promise<void> {
  await axios.delete(`/ai-config/${id}`);
}

// 获取可用的AI模型列表
export async function getAvailableModels(): Promise<Record<string, ModelOption[]>> {
  const response = await axios.get('/ai-config/models');
  return response.data.data;
}

// 获取知识库文件列表
export async function getKnowledgeFiles(page = 1, pageSize = 10): Promise<PaginatedData<KnowledgeFile>> {
  const response = await axios.get(`/knowledge/files?page=${page}&page_size=${pageSize}`);
  return response.data.data;
}

// 获取知识库文件详情
export async function getKnowledgeFileDetail(id: number): Promise<KnowledgeFile> {
  const response = await axios.get(`/knowledge/files/${id}`);
  return response.data.data;
}

// 删除知识库文件
export async function deleteKnowledgeFile(id: number): Promise<void> {
  await axios.delete(`/knowledge/files/${id}`);
}

// 上传知识库文件
export async function uploadKnowledgeFile(file: File): Promise<KnowledgeFile> {
  const formData = new FormData();
  formData.append('file', file);
  
  const response = await axios.post('/knowledge/upload', formData, {
    headers: {
      'Content-Type': 'multipart/form-data'
    }
  });
  
  return response.data.data;
}

// 创建EventSource用于流式聊天
export function createChatStream(
  message: string, 
  sessionId = 0, 
  aiConfigId = 0, 
  knowledgeIds: number[] = []
): { eventSource: EventSource, controller: AbortController } {
  const controller = new AbortController();
  const params = new URLSearchParams({
    session_id: sessionId.toString(),
    message,
    ai_config_id: aiConfigId.toString(),
    knowledge_ids: JSON.stringify(knowledgeIds)
  });
  
  // 获取token并添加到URL
  const token = localStorage.getItem('token');
  
  console.log('创建流式聊天连接，参数:', {
    sessionId,
    message: message.substring(0, 20) + (message.length > 20 ? '...' : ''),
    aiConfigId,
    knowledgeIds
  });
  
  // 创建EventSource
  const baseUrl = "http://localhost:14020";
  const eventSource = new EventSource(`${baseUrl}/api/v1/chat/stream?${params.toString()}&token=${token}`);
  
  // 添加事件监听器记录连接状态
  eventSource.onopen = () => {
    console.log('流式连接已打开');
  };
  
  eventSource.onerror = (error) => {
    console.error('流式连接错误:', error);
  };
  
  // 设置中止控制器，用于关闭连接
  controller.signal.addEventListener('abort', () => {
    console.log('手动关闭流式连接');
    eventSource.close();
  });
  
  return { eventSource, controller };
} 