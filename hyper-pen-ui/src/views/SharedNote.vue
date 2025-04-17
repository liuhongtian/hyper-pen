<template>
  <div class="shared-note">
    <el-card v-if="note" class="note-card">
      <template #header>
        <div class="note-header">
          <h2>{{ note.title }}</h2>
          <div class="note-meta">
            <el-tag v-if="note.category" type="info">
              {{ note.category.name }}
            </el-tag>
            <el-tag
              v-for="tag in note.tags"
              :key="tag.id"
              type="success"
              class="ml-2"
            >
              {{ tag.name }}
            </el-tag>
          </div>
        </div>
      </template>
      
      <markdown-viewer :content="note.content" />
      
      <div class="note-footer">
        <div class="note-info">
          <span>最后更新：{{ formatDate(note.updated_at) }}</span>
        </div>
      </div>
    </el-card>
    
    <el-empty v-else-if="!loading" description="笔记不存在或已过期" />
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import axios from 'axios'
import dayjs from 'dayjs'
import MarkdownViewer from '../components/MarkdownViewer.vue'

const route = useRoute()
const note = ref(null)
const loading = ref(true)

// 获取共享笔记
const fetchSharedNote = async () => {
  try {
    const token = route.params.token
    const response = await axios.get(`/api/shared/${token}`)
    note.value = response.data
  } catch (error) {
    console.error('获取共享笔记失败:', error)
  } finally {
    loading.value = false
  }
}

// 格式化日期
const formatDate = (date) => {
  return dayjs(date).format('YYYY-MM-DD HH:mm')
}

onMounted(() => {
  fetchSharedNote()
})
</script>

<style scoped>
.shared-note {
  max-width: 800px;
  margin: 0 auto;
  padding: 20px;
  
  .note-card {
    .note-header {
      h2 {
        margin: 0 0 10px 0;
      }
      
      .note-meta {
        display: flex;
        align-items: center;
        flex-wrap: wrap;
        gap: 8px;
      }
    }
    
    .note-content {
      min-height: 200px;
    }
    
    .note-footer {
      margin-top: 20px;
      padding-top: 20px;
      border-top: 1px solid #eee;
      
      .note-info {
        color: #666;
        font-size: 14px;
      }
    }
  }
}
</style> 