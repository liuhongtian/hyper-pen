<template>
  <div class="note-import-export">
    <el-dialog
      v-model="dialogVisible"
      title="导入/导出笔记"
      width="40%"
    >
      <div class="import-export-content">
        <el-tabs v-model="activeTab">
          <el-tab-pane label="导入笔记" name="import">
            <div class="import-section">
              <el-upload
                class="upload-demo"
                drag
                action="/api/notes/import"
                :headers="uploadHeaders"
                :on-success="handleImportSuccess"
                :on-error="handleImportError"
                :before-upload="beforeImportUpload"
              >
                <el-icon class="el-icon--upload"><upload-filled /></el-icon>
                <div class="el-upload__text">
                  将文件拖到此处，或<em>点击上传</em>
                </div>
                <template #tip>
                  <div class="el-upload__tip">
                    支持导入 Markdown 文件 (.md) 或 JSON 文件
                  </div>
                </template>
              </el-upload>
            </div>
          </el-tab-pane>
          
          <el-tab-pane label="导出笔记" name="export">
            <div class="export-section">
              <el-form :model="exportForm" label-width="80px">
                <el-form-item label="导出格式">
                  <el-radio-group v-model="exportForm.format">
                    <el-radio label="markdown">Markdown</el-radio>
                    <el-radio label="json">JSON</el-radio>
                  </el-radio-group>
                </el-form-item>
                
                <el-form-item label="导出范围">
                  <el-select v-model="exportForm.range" placeholder="请选择导出范围">
                    <el-option label="所有笔记" value="all" />
                    <el-option label="当前笔记" value="current" />
                  </el-select>
                </el-form-item>
                
                <el-form-item>
                  <el-button type="primary" @click="handleExport">导出</el-button>
                </el-form-item>
              </el-form>
            </div>
          </el-tab-pane>
        </el-tabs>
      </div>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { ElMessage } from 'element-plus'
import { UploadFilled } from '@element-plus/icons-vue'

const props = defineProps({
  noteId: {
    type: String,
    required: true
  }
})

const dialogVisible = ref(false)
const activeTab = ref('import')
const exportForm = ref({
  format: 'markdown',
  range: 'current'
})

const uploadHeaders = {
  'Authorization': `Bearer ${localStorage.getItem('token')}`
}

const beforeImportUpload = (file) => {
  const isMarkdown = file.type === 'text/markdown' || file.name.endsWith('.md')
  const isJson = file.type === 'application/json' || file.name.endsWith('.json')
  
  if (!isMarkdown && !isJson) {
    ElMessage.error('只能上传 Markdown 或 JSON 文件！')
    return false
  }
  return true
}

const handleImportSuccess = (response) => {
  ElMessage.success('导入成功')
  dialogVisible.value = false
}

const handleImportError = () => {
  ElMessage.error('导入失败')
}

const handleExport = async () => {
  try {
    const token = localStorage.getItem('token')
    const url = exportForm.value.range === 'all' 
      ? `/api/notes/export?format=${exportForm.value.format}`
      : `/api/notes/${props.noteId}/export?format=${exportForm.value.format}`
    
    const response = await fetch(url, {
      headers: {
        'Authorization': `Bearer ${token}`
      }
    })
    
    if (!response.ok) throw new Error('导出失败')
    
    const blob = await response.blob()
    const downloadUrl = window.URL.createObjectURL(blob)
    const link = document.createElement('a')
    link.href = downloadUrl
    link.download = `notes.${exportForm.value.format}`
    document.body.appendChild(link)
    link.click()
    document.body.removeChild(link)
    window.URL.revokeObjectURL(downloadUrl)
    
    ElMessage.success('导出成功')
    dialogVisible.value = false
  } catch (error) {
    ElMessage.error(error.message)
  }
}

const show = () => {
  dialogVisible.value = true
}

defineExpose({
  show
})
</script>

<style scoped>
.import-export-content {
  padding: 20px;
}

.import-section,
.export-section {
  padding: 20px;
}

.upload-demo {
  width: 100%;
}

.el-upload__tip {
  margin-top: 10px;
  color: #909399;
}
</style> 