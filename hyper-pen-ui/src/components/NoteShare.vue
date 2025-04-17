<template>
  <div class="note-share">
    <el-dialog
      v-model="dialogVisible"
      title="分享笔记"
      width="40%"
    >
      <div class="share-content">
        <div class="share-links" v-if="shareLinks.length > 0">
          <div v-for="link in shareLinks" :key="link.id" class="share-link-item">
            <div class="link-info">
              <span class="link-url">{{ getShareUrl(link.token) }}</span>
              <span class="link-expiry">有效期至：{{ formatDate(link.expiresAt) }}</span>
            </div>
            <div class="link-actions">
              <el-button type="primary" link @click="copyLink(link.token)">复制链接</el-button>
              <el-button type="danger" link @click="deleteLink(link.id)">删除</el-button>
            </div>
          </div>
        </div>
        <div v-else class="no-links">
          <p>暂无分享链接</p>
        </div>
        
        <div class="create-link">
          <el-button type="primary" @click="createShareLink">创建分享链接</el-button>
        </div>
      </div>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import dayjs from 'dayjs'

const props = defineProps({
  noteId: {
    type: String,
    required: true
  }
})

const dialogVisible = ref(false)
const shareLinks = ref([])

const fetchShareLinks = async () => {
  try {
    const token = localStorage.getItem('token')
    const response = await fetch(`/api/notes/${props.noteId}/share-links`, {
      headers: {
        'Authorization': `Bearer ${token}`
      }
    })
    if (!response.ok) throw new Error('获取分享链接失败')
    shareLinks.value = await response.json()
  } catch (error) {
    ElMessage.error(error.message)
  }
}

const createShareLink = async () => {
  try {
    const token = localStorage.getItem('token')
    const response = await fetch(`/api/notes/${props.noteId}/share-links`, {
      method: 'POST',
      headers: {
        'Authorization': `Bearer ${token}`
      }
    })
    if (!response.ok) throw new Error('创建分享链接失败')
    await fetchShareLinks()
    ElMessage.success('创建成功')
  } catch (error) {
    ElMessage.error(error.message)
  }
}

const deleteLink = async (linkId) => {
  try {
    await ElMessageBox.confirm('确定要删除这个分享链接吗？', '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })
    
    const token = localStorage.getItem('token')
    const response = await fetch(`/api/share-links/${linkId}`, {
      method: 'DELETE',
      headers: {
        'Authorization': `Bearer ${token}`
      }
    })
    
    if (!response.ok) throw new Error('删除分享链接失败')
    
    await fetchShareLinks()
    ElMessage.success('删除成功')
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error(error.message)
    }
  }
}

const copyLink = (token) => {
  const url = getShareUrl(token)
  navigator.clipboard.writeText(url)
    .then(() => {
      ElMessage.success('链接已复制到剪贴板')
    })
    .catch(() => {
      ElMessage.error('复制失败')
    })
}

const getShareUrl = (token) => {
  return `${window.location.origin}/shared/${token}`
}

const formatDate = (date) => {
  return dayjs(date).format('YYYY-MM-DD HH:mm')
}

const show = () => {
  dialogVisible.value = true
  fetchShareLinks()
}

defineExpose({
  show
})
</script>

<style scoped>
.share-content {
  padding: 20px;
}

.share-links {
  margin-bottom: 20px;
}

.share-link-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 10px;
  border: 1px solid #dcdfe6;
  border-radius: 4px;
  margin-bottom: 10px;
}

.link-info {
  display: flex;
  flex-direction: column;
}

.link-url {
  font-size: 14px;
  color: #409EFF;
  margin-bottom: 5px;
}

.link-expiry {
  font-size: 12px;
  color: #909399;
}

.link-actions {
  display: flex;
  gap: 10px;
}

.no-links {
  text-align: center;
  padding: 20px;
  color: #909399;
}

.create-link {
  text-align: center;
  margin-top: 20px;
}
</style> 