<template>
  <div class="wechat-login">
    <el-dialog
      v-model="dialogVisible"
      title="微信扫码登录"
      width="400px"
      :close-on-click-modal="false"
      :close-on-press-escape="false"
      :show-close="false"
    >
      <div class="qrcode-container">
        <img :src="qrcodeUrl" alt="微信登录二维码" />
      </div>
      <div class="status-text">{{ statusText }}</div>
      <template #footer>
        <el-button @click="handleCancel">取消</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted, watch } from 'vue'
import { ElMessage } from 'element-plus'
import axios from 'axios'

const props = defineProps({
  visible: {
    type: Boolean,
    default: false
  }
})

const emit = defineEmits(['update:visible', 'success'])

const dialogVisible = ref(false)
const qrcodeUrl = ref('')
const statusText = ref('请使用微信扫码登录')
const timer = ref(null)

// 监听visible属性变化
watch(() => props.visible, (val) => {
  dialogVisible.value = val
  if (val) {
    startWechatLogin()
  } else {
    stopPolling()
  }
})

// 监听dialogVisible变化
watch(dialogVisible, (val) => {
  if (!val) {
    emit('update:visible', false)
  }
})

// 开始微信登录流程
const startWechatLogin = async () => {
  try {
    const response = await axios.get('/api/auth/wechat/login')
    qrcodeUrl.value = response.data.qrcode_url
    startPolling()
  } catch (error) {
    ElMessage.error('获取二维码失败')
    handleCancel()
  }
}

// 开始轮询检查登录状态
const startPolling = () => {
  timer.value = setInterval(async () => {
    try {
      const response = await axios.get('/api/auth/wechat/check')
      if (response.data.status === 'success') {
        stopPolling()
        emit('success', response.data)
        handleCancel()
      }
    } catch (error) {
      console.error('检查登录状态失败:', error)
    }
  }, 2000)
}

// 停止轮询
const stopPolling = () => {
  if (timer.value) {
    clearInterval(timer.value)
    timer.value = null
  }
}

// 取消登录
const handleCancel = () => {
  stopPolling()
  dialogVisible.value = false
}

// 组件卸载时清理定时器
onUnmounted(() => {
  stopPolling()
})
</script>

<style scoped>
.wechat-login {
  .qrcode-container {
    display: flex;
    justify-content: center;
    margin: 20px 0;
    
    img {
      width: 200px;
      height: 200px;
    }
  }
  
  .status-text {
    text-align: center;
    color: #666;
    margin-bottom: 20px;
  }
}
</style> 