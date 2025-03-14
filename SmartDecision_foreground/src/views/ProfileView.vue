<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { useUserStore } from '../stores/user'
import { storeToRefs } from 'pinia'
import { ElMessage, ElMessageBox } from 'element-plus'
import type { FormInstance, FormRules } from 'element-plus'
import { User, Lock, Edit, Camera, ArrowLeft } from '@element-plus/icons-vue'
import { useRouter } from 'vue-router'

const userStore = useUserStore()
const { userInfo, loading } = storeToRefs(userStore)
const router = useRouter()

// 表单引用
const profileFormRef = ref<FormInstance>()
const passwordFormRef = ref<FormInstance>()

// 头像上传
const avatarUrl = ref('')
const isUploadingAvatar = ref(false)

// 激活的标签页
const activeTab = ref('profile')

// 表单数据
const profileForm = reactive({
  username: '',
  email: '',
  nickname: '',
  bio: ''
})

const passwordForm = reactive({
  oldPassword: '',
  newPassword: '',
  confirmPassword: ''
})

// 是否编辑模式
const isEditingProfile = ref(false)

// 加载用户数据
onMounted(() => {
  if (userInfo.value) {
    profileForm.username = userInfo.value.username
    profileForm.email = userInfo.value.email
    profileForm.nickname = userInfo.value.username // 假设昵称初始值为用户名
    profileForm.bio = '我是一个DeepSeek AI的用户' // 默认简介
    
    // 设置默认头像
    avatarUrl.value = generateAvatarUrl(userInfo.value.username)
  }
})

// 生成基于用户名的头像URL
const generateAvatarUrl = (username: string) => {
  const colors = ['4facfe', '00f2fe', '0acffe', '495aff', '6a75ff']
  const colorIndex = username.charCodeAt(0) % colors.length
  const bgColor = colors[colorIndex]
  const initial = username.charAt(0).toUpperCase()
  return `https://ui-avatars.com/api/?name=${initial}&background=${bgColor}&color=fff&size=256`
}

// 表单验证规则
const profileRules = reactive<FormRules>({
  username: [
    { required: true, message: '请输入用户名', trigger: 'blur' },
    { min: 3, message: '用户名长度至少为3个字符', trigger: 'blur' }
  ],
  nickname: [
    { required: true, message: '请输入昵称', trigger: 'blur' }
  ],
  bio: [
    { max: 200, message: '简介最多200字符', trigger: 'blur' }
  ]
})

// 密码验证规则
const validatePass = (rule: any, value: string, callback: any) => {
  if (value === '') {
    callback(new Error('请输入新密码'))
  } else {
    if (passwordForm.confirmPassword !== '') {
      if (passwordFormRef.value) {
        passwordFormRef.value.validateField('confirmPassword')
      }
    }
    callback()
  }
}

const validateConfirmPass = (rule: any, value: string, callback: any) => {
  if (value === '') {
    callback(new Error('请再次输入新密码'))
  } else if (value !== passwordForm.newPassword) {
    callback(new Error('两次输入密码不一致'))
  } else {
    callback()
  }
}

const passwordRules = reactive<FormRules>({
  oldPassword: [
    { required: true, message: '请输入当前密码', trigger: 'blur' },
    { min: 6, message: '密码长度至少为6个字符', trigger: 'blur' }
  ],
  newPassword: [
    { required: true, message: '请输入新密码', trigger: 'blur' },
    { min: 6, message: '密码长度至少为6个字符', trigger: 'blur' },
    { validator: validatePass, trigger: 'blur' }
  ],
  confirmPassword: [
    { required: true, message: '请再次输入新密码', trigger: 'blur' },
    { validator: validateConfirmPass, trigger: 'blur' }
  ]
})

// 开始编辑资料
const startEditing = () => {
  isEditingProfile.value = true
}

// 取消编辑
const cancelEditing = () => {
  isEditingProfile.value = false
  
  // 恢复原始数据
  if (userInfo.value) {
    profileForm.username = userInfo.value.username
    profileForm.nickname = userInfo.value.username
  }
}

// 更新个人资料
const updateProfile = async (formEl: FormInstance | undefined) => {
  if (!formEl) return
  
  await formEl.validate(async (valid) => {
    if (valid) {
      // 这里应该调用API更新个人资料
      setTimeout(() => {
        ElMessage.success('个人资料更新成功')
        isEditingProfile.value = false
        
        // 更新本地存储的用户信息
        if (userInfo.value) {
          userInfo.value.username = profileForm.username
        }
      }, 800)
    } else {
      ElMessage.error('请正确填写表单')
      return false
    }
  })
}

// 更新密码
const updatePassword = async (formEl: FormInstance | undefined) => {
  if (!formEl) return
  
  await formEl.validate(async (valid) => {
    if (valid) {
      // 这里应该调用API更新密码
      setTimeout(() => {
        ElMessage.success('密码修改成功')
        
        // 清空密码字段
        passwordForm.oldPassword = ''
        passwordForm.newPassword = ''
        passwordForm.confirmPassword = ''
      }, 800)
    } else {
      ElMessage.error('请正确填写表单')
      return false
    }
  })
}

// 处理头像上传
const handleAvatarUpload = () => {
  ElMessageBox.alert('头像上传功能正在开发中', '提示', {
    confirmButtonText: '确定',
    type: 'info'
  })
}

// 获取加入天数
const getJoinDays = () => {
  // 假设用户是30天前加入的
  return 30
}
</script>

<template>
  <div class="profile-container">
    <div class="profile-header">
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
        <h1>个人中心</h1>
      </div>
    </div>
    
    <div class="profile-content">
      <el-tabs v-model="activeTab" class="profile-tabs">
        <el-tab-pane label="个人资料" name="profile">
          <div class="profile-card">
            <div class="profile-avatar-section">
              <div class="avatar-container">
                <el-avatar :size="120" :src="avatarUrl">
                  {{ profileForm.username ? profileForm.username.charAt(0).toUpperCase() : 'U' }}
                </el-avatar>
                <div class="avatar-edit-overlay" @click="handleAvatarUpload">
                  <el-icon><Camera /></el-icon>
                </div>
              </div>
              <div class="profile-name">
                <h2>{{ profileForm.nickname || profileForm.username }}</h2>
                <p class="user-status">
                  <el-tag size="small" type="success">已认证用户</el-tag>
                </p>
                <p class="user-joined">
                  已加入 {{ getJoinDays() }} 天
                </p>
              </div>
            </div>
            
            <div class="profile-edit-actions" v-if="!isEditingProfile">
              <el-button type="primary" @click="startEditing" :icon="Edit">
                编辑资料
              </el-button>
            </div>
            
            <el-form
              ref="profileFormRef"
              :model="profileForm"
              :rules="profileRules"
              label-position="top"
              class="profile-form"
            >
              <el-form-item label="用户名" prop="username">
                <el-input v-model="profileForm.username" :disabled="!isEditingProfile" />
              </el-form-item>
              
              <el-form-item label="邮箱" prop="email">
                <el-input v-model="profileForm.email" disabled />
                <div class="form-item-tip">邮箱地址不可修改</div>
              </el-form-item>
              
              <el-form-item label="昵称" prop="nickname">
                <el-input v-model="profileForm.nickname" :disabled="!isEditingProfile" />
              </el-form-item>
              
              <el-form-item label="个人简介" prop="bio">
                <el-input
                  v-model="profileForm.bio"
                  :disabled="!isEditingProfile"
                  type="textarea"
                  :rows="4"
                  placeholder="介绍一下自己吧"
                />
              </el-form-item>
              
              <div class="form-actions" v-if="isEditingProfile">
                <el-button @click="cancelEditing">取消</el-button>
                <el-button
                  type="primary"
                  @click="updateProfile(profileFormRef)"
                  :loading="loading"
                >
                  保存修改
                </el-button>
              </div>
            </el-form>
          </div>
        </el-tab-pane>
        
        <el-tab-pane label="修改密码" name="password">
          <div class="profile-card">
            <h3>修改密码</h3>
            <p class="section-desc">为了保障您的账号安全，建议定期修改密码</p>
            
            <el-form
              ref="passwordFormRef"
              :model="passwordForm"
              :rules="passwordRules"
              label-position="top"
              class="password-form"
            >
              <el-form-item label="当前密码" prop="oldPassword">
                <el-input
                  v-model="passwordForm.oldPassword"
                  type="password"
                  show-password
                  placeholder="输入当前密码"
                />
              </el-form-item>
              
              <el-form-item label="新密码" prop="newPassword">
                <el-input
                  v-model="passwordForm.newPassword"
                  type="password"
                  show-password
                  placeholder="输入新密码"
                />
              </el-form-item>
              
              <el-form-item label="确认新密码" prop="confirmPassword">
                <el-input
                  v-model="passwordForm.confirmPassword"
                  type="password"
                  show-password
                  placeholder="再次输入新密码"
                />
              </el-form-item>
              
              <div class="form-actions">
                <el-button
                  type="primary"
                  @click="updatePassword(passwordFormRef)"
                  :loading="loading"
                >
                  修改密码
                </el-button>
              </div>
            </el-form>
          </div>
        </el-tab-pane>
        
        <el-tab-pane label="账号安全" name="security">
          <div class="profile-card">
            <h3>安全设置</h3>
            <p class="section-desc">管理您的账号安全选项</p>
            
            <div class="security-items">
              <div class="security-item">
                <div class="security-info">
                  <el-icon><Lock /></el-icon>
                  <div>
                    <h4>双因素认证</h4>
                    <p>增强账号安全性，使用验证器生成一次性验证码</p>
                  </div>
                </div>
                <el-button>设置</el-button>
              </div>
              
              <div class="security-item">
                <div class="security-info">
                  <el-icon><User /></el-icon>
                  <div>
                    <h4>账号登录历史</h4>
                    <p>查看您的账号登录记录和设备信息</p>
                  </div>
                </div>
                <el-button>查看</el-button>
              </div>
            </div>
          </div>
        </el-tab-pane>
      </el-tabs>
    </div>
  </div>
</template>

<style scoped>
.profile-container {
  min-height: 100vh;
  background-color: var(--page-bg);
}

.profile-header {
  background-color: #4facfe;
  background-image: linear-gradient(135deg, #4facfe 0%, #00f2fe 100%);
  color: white;
  padding: 40px 0 60px;
}

.header-content {
  max-width: 1200px;
  margin: 0 auto;
  padding: 0 24px;
  display: flex;
  align-items: center;
  position: relative;
}

.profile-header h1 {
  font-size: 28px;
  margin-bottom: 8px;
  margin-left: 20px;
}

.back-button {
  transition: all 0.3s ease;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.1);
}

.back-button:hover {
  transform: translateX(-5px);
  box-shadow: 0 4px 15px rgba(0, 0, 0, 0.2);
}

.profile-content {
  max-width: 1200px;
  margin: -40px auto 0;
  padding: 0 24px 40px;
  position: relative;
}

.profile-tabs {
  background-color: white;
  border-radius: 8px;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.1);
  overflow: hidden;
}

.profile-card {
  padding: 30px;
}

.profile-avatar-section {
  display: flex;
  align-items: center;
  margin-bottom: 30px;
}

.avatar-container {
  position: relative;
  margin-right: 24px;
}

.avatar-edit-overlay {
  position: absolute;
  bottom: 0;
  right: 0;
  background-color: rgba(0, 0, 0, 0.6);
  width: 36px;
  height: 36px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
  cursor: pointer;
  transition: all 0.3s;
}

.avatar-edit-overlay:hover {
  background-color: rgba(0, 0, 0, 0.8);
}

.profile-name h2 {
  font-size: 24px;
  margin: 0 0 8px;
}

.user-status {
  margin-bottom: 6px;
}

.user-joined {
  color: #909399;
  font-size: 14px;
}

.profile-edit-actions {
  display: flex;
  justify-content: flex-end;
  margin-bottom: 20px;
}

.profile-form, .password-form {
  max-width: 600px;
}

.form-item-tip {
  font-size: 12px;
  color: #909399;
  margin-top: 4px;
}

.form-actions {
  margin-top: 30px;
  display: flex;
  justify-content: flex-end;
  gap: 12px;
}

h3 {
  font-size: 20px;
  color: #303133;
  margin: 0 0 10px;
}

.section-desc {
  color: #909399;
  margin-bottom: 30px;
}

.security-items {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.security-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 20px;
  background-color: #f5f7fa;
  border-radius: 8px;
}

.security-info {
  display: flex;
  align-items: center;
  gap: 16px;
}

.security-info .el-icon {
  font-size: 24px;
  color: #409eff;
}

.security-info h4 {
  font-size: 16px;
  margin: 0 0 4px;
}

.security-info p {
  color: #909399;
  margin: 0;
}

@media (max-width: 768px) {
  .profile-avatar-section {
    flex-direction: column;
    text-align: center;
  }
  
  .avatar-container {
    margin-right: 0;
    margin-bottom: 20px;
  }
  
  .security-item {
    flex-direction: column;
    gap: 20px;
    align-items: flex-start;
  }
  
  .header-content {
    flex-direction: column;
    align-items: flex-start;
  }
  
  .profile-header h1 {
    margin-top: 15px;
    margin-left: 0;
  }
}
</style> 