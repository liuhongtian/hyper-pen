<template>
  <div class="home-container">
    <el-card class="intro-card">
      <template #header>
        <div class="card-header">
          <h1>欢迎使用 Hyper Pen</h1>
        </div>
      </template>
      
      <div class="intro-content">
        <h2>系统介绍</h2>
        <p>Hyper Pen 是一个功能强大的在线笔记应用，支持以下特性：</p>
        <ul>
          <li>使用 Markdown 格式编写笔记</li>
          <li>实时预览功能</li>
          <li>笔记分类和标签管理</li>
          <li>笔记分享功能</li>
          <li>支持 GitHub 和微信登录</li>
        </ul>
        
        <div class="action-buttons">
          <el-button type="primary" @click="$router.push('/login')">立即登录</el-button>
          <el-button @click="$router.push('/register')">注册账号</el-button>
        </div>
      </div>
    </el-card>
  </div>
</template>

<script setup>
import { useRouter } from 'vue-router'

const router = useRouter()
</script>

<style scoped>
.home-container {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 100vh;
  padding: 20px;
}

.intro-card {
  max-width: 800px;
  width: 100%;
}

.card-header {
  text-align: center;
}

.intro-content {
  padding: 20px;
  
  h2 {
    margin-bottom: 20px;
  }
  
  p {
    margin-bottom: 15px;
    font-size: 16px;
  }
  
  ul {
    margin-bottom: 30px;
    padding-left: 20px;
    
    li {
      margin-bottom: 10px;
      font-size: 16px;
    }
  }
}

.action-buttons {
  display: flex;
  justify-content: center;
  gap: 20px;
  margin-top: 30px;
}
</style> 