<template>
  <div class="login-container">
    <el-card class="login-card">
      <template #header>
        <div class="card-header">
          <h2>登录</h2>
        </div>
      </template>
      
      <el-form :model="loginForm" :rules="rules" ref="loginFormRef" label-width="80px">
        <el-form-item label="用户名" prop="username">
          <el-input v-model="loginForm.username" />
        </el-form-item>
        <el-form-item label="密码" prop="password">
          <el-input v-model="loginForm.password" type="password" />
        </el-form-item>
        
        <el-form-item>
          <el-button type="primary" @click="handleLogin">登录</el-button>
          <el-button @click="$router.push('/register')">注册</el-button>
        </el-form-item>
        
        <el-divider>第三方登录</el-divider>
        
        <div class="social-login">
          <el-button type="primary" plain @click="handleGithubLogin">
            <el-icon><ChromeFilled /></el-icon>
            GitHub登录
          </el-button>
          <el-button type="success" plain @click="handleWechatLogin">
            <el-icon><Eleme /></el-icon>
            微信登录
          </el-button>
        </div>
      </el-form>
    </el-card>
    
    <WechatLogin
      v-model:visible="showWechatLogin"
      @success="handleWechatLoginSuccess"
    />
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { ChromeFilled, Eleme } from '@element-plus/icons-vue'
import { login } from '../utils/auth'
import WechatLogin from '@/components/WechatLogin.vue'

const router = useRouter()
const loginFormRef = ref(null)

const loginForm = ref({
  username: '',
  password: ''
})

const rules = {
  username: [
    { required: true, message: '请输入用户名', trigger: 'blur' }
  ],
  password: [
    { required: true, message: '请输入密码', trigger: 'blur' }
  ]
}

const showWechatLogin = ref(false)

const handleLogin = async () => {
  try {
    const response = await fetch('/api/auth/login', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json'
      },
      body: JSON.stringify({
        username: loginForm.value.username,
        password: loginForm.value.password
      })
    })

    if (!response.ok) {
      const error = await response.json()
      throw new Error(error.error || '登录失败')
    }

    const data = await response.json()
    localStorage.setItem('token', data.token)
    localStorage.setItem('user', JSON.stringify(data.user))
    ElMessage.success('登录成功')
    router.push('/notes')
  } catch (error) {
    ElMessage.error(error.message)
  }
}

const handleGithubLogin = () => {
  // TODO: 实现GitHub登录
}

const handleWechatLogin = () => {
  showWechatLogin.value = true
}

const handleWechatLoginSuccess = (data) => {
  // 保存token和用户信息
  localStorage.setItem('token', data.token)
  localStorage.setItem('user', JSON.stringify(data.user))
  
  ElMessage.success('登录成功')
  router.push('/notes')
}
</script>

<style scoped>
.login-container {
  display: flex;
  justify-content: center;
  align-items: center;
  height: 100%;
}

.login-card {
  width: 400px;
}

.card-header {
  text-align: center;
}

.social-login {
  display: flex;
  justify-content: space-around;
  margin-top: 20px;
}
</style> 