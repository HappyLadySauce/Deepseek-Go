import { defineStore } from 'pinia'
import { ref } from 'vue'
import axios from 'axios'
import { ElMessage } from 'element-plus'

// 消息类型
export interface Message {
  id: string
  content: string
  role: 'user' | 'assistant'
  timestamp: Date
}

// 聊天会话类型
export interface ChatSession {
  id: string
  title: string
  messages: Message[]
  createdAt: Date
  updatedAt: Date
}

export const useChatStore = defineStore('chat', () => {
  // 状态
  const chatSessions = ref<ChatSession[]>([])
  const currentSessionId = ref<string | null>(null)
  const loading = ref(false)

  // 获取当前会话
  const getCurrentSession = (): ChatSession | undefined => {
    if (!currentSessionId.value) return undefined
    return chatSessions.value.find(session => session.id === currentSessionId.value)
  }

  // 创建新会话
  const createNewSession = (): string => {
    const newSession: ChatSession = {
      id: Date.now().toString(),
      title: '新对话',
      messages: [],
      createdAt: new Date(),
      updatedAt: new Date()
    }
    
    chatSessions.value.unshift(newSession)
    currentSessionId.value = newSession.id
    return newSession.id
  }

  // 切换会话
  const switchSession = (sessionId: string) => {
    const session = chatSessions.value.find(s => s.id === sessionId)
    if (session) {
      currentSessionId.value = sessionId
    }
  }

  // 发送消息
  const sendMessage = async (content: string): Promise<boolean> => {
    if (!content.trim()) {
      ElMessage.warning('请输入消息内容')
      return false
    }
    
    // 确保有当前会话
    if (!currentSessionId.value) {
      createNewSession()
    }
    
    const session = getCurrentSession()
    if (!session) return false
    
    // 创建用户消息
    const userMessage: Message = {
      id: Date.now().toString(),
      content,
      role: 'user',
      timestamp: new Date()
    }
    
    // 添加到会话
    session.messages.push(userMessage)
    session.updatedAt = new Date()
    
    // 如果是第一条消息，更新会话标题
    if (session.messages.length === 1) {
      session.title = content.length > 20 ? content.substring(0, 20) + '...' : content
    }
    
    loading.value = true
    
    try {
      // 发送请求到AI聊天API
      const response = await axios.post('/api/v1/chat', {
        message: content,
        session_id: session.id
      })
      
      // 添加AI回复
      const aiMessage: Message = {
        id: (Date.now() + 1).toString(),
        content: response.data.message,
        role: 'assistant',
        timestamp: new Date()
      }
      
      session.messages.push(aiMessage)
      session.updatedAt = new Date()
      
      return true
    } catch (error: any) {
      ElMessage.error(error.response?.data?.error || 'AI响应失败，请稍后再试')
      return false
    } finally {
      loading.value = false
    }
  }

  // 删除会话
  const deleteSession = (sessionId: string) => {
    const index = chatSessions.value.findIndex(s => s.id === sessionId)
    if (index !== -1) {
      chatSessions.value.splice(index, 1)
      
      // 如果删除的是当前会话，切换到其他会话或创建新会话
      if (currentSessionId.value === sessionId) {
        if (chatSessions.value.length > 0) {
          currentSessionId.value = chatSessions.value[0].id
        } else {
          currentSessionId.value = null
        }
      }
    }
  }

  // 清空会话消息
  const clearSession = (sessionId: string) => {
    const session = chatSessions.value.find(s => s.id === sessionId)
    if (session) {
      session.messages = []
      session.updatedAt = new Date()
    }
  }

  // 初始化
  const initialize = () => {
    // 从本地存储加载会话
    const savedSessions = localStorage.getItem('chatSessions')
    if (savedSessions) {
      try {
        const parsed = JSON.parse(savedSessions)
        chatSessions.value = parsed.map((session: any) => ({
          ...session,
          createdAt: new Date(session.createdAt),
          updatedAt: new Date(session.updatedAt),
          messages: session.messages.map((msg: any) => ({
            ...msg,
            timestamp: new Date(msg.timestamp)
          }))
        }))
      } catch (e) {
        console.error('Failed to parse saved chat sessions', e)
      }
    }
    
    // 从本地存储加载当前会话ID
    const savedCurrentId = localStorage.getItem('currentChatSessionId')
    if (savedCurrentId) {
      currentSessionId.value = savedCurrentId
    }
    
    // 如果没有会话，创建一个新会话
    if (chatSessions.value.length === 0) {
      createNewSession()
    }
  }

  // 保存会话到本地存储
  const saveToLocalStorage = () => {
    localStorage.setItem('chatSessions', JSON.stringify(chatSessions.value))
    if (currentSessionId.value) {
      localStorage.setItem('currentChatSessionId', currentSessionId.value)
    }
  }

  return {
    chatSessions,
    currentSessionId,
    loading,
    getCurrentSession,
    createNewSession,
    switchSession,
    sendMessage,
    deleteSession,
    clearSession,
    initialize,
    saveToLocalStorage
  }
}) 