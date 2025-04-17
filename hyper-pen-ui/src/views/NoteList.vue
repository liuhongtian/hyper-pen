<template>
  <div class="note-list-container">
    <el-row :gutter="20">
      <el-col :span="isSidebarVisible ? 6 : 1">
        <el-card class="note-list-card">
          <template #header>
            <div class="card-header">
              <el-button @click="toggleSidebar" :icon="isSidebarVisible ? ArrowLeft : ArrowRight" circle />
              <div class="action-buttons">
                <el-button type="success" circle plain @click="createNewNote"><el-icon><DocumentAdd /></el-icon></el-button>
                <el-button type="warning" circle plain @click="showImportExport"><el-icon><Upload /></el-icon></el-button>
              </div>
            </div>
          </template>

          <el-tabs v-model="activeTab">
            <el-tab-pane label="笔记列表" name="list">
              <el-input v-model="searchQuery" placeholder="搜索笔记" class="search-input" clearable />

              <el-menu :default-active="activeNoteId" class="note-menu" @select="handleNoteSelect">
                <el-menu-item v-for="note in filteredNotes" :key="note.id" :index="note.id">
                  <template #title>
                    <div class="note-item">
                      <span class="note-title">{{ note.title }}</span>
                      <span class="note-date">{{ formatDate(note.updatedAt) }}</span>
                    </div>
                  </template>
                </el-menu-item>
              </el-menu>
            </el-tab-pane>

            <el-tab-pane label="分类管理" name="categories">
              <CategoryManager />
            </el-tab-pane>

            <el-tab-pane label="标签管理" name="tags">
              <TagManager />
            </el-tab-pane>

            <el-tab-pane label="搜索笔记" name="search">
              <NoteSearch />
            </el-tab-pane>
          </el-tabs>
        </el-card>
      </el-col>

      <el-col :span="isSidebarVisible ? 18 : 23">
        <el-card class="editor-card">
          <template #header>
            <div class="card-header">
              <el-button @click="toggleSidebar" :icon="isSidebarVisible ? ArrowLeft : ArrowRight" circle />
              <el-input v-model="currentNote.title" placeholder="请输入笔记标题" class="note-title-input" />
              <div class="action-buttons">
                <el-button type="primary" @click="saveNote">保存</el-button>
                <el-button @click="showShareDialog" v-if="currentNote.id">分享</el-button>
                <el-button @click="deleteNote" type="danger">删除</el-button>
              </div>
            </div>
          </template>

          <el-form :model="currentNote">
            <el-row :gutter="20">
              <el-col :span="6">
                <el-form-item>
                  <el-select v-model="currentNote.category_id" placeholder="选择分类" clearable>
                    <el-option v-for="category in categories" :key="category.id" :label="category.name"
                      :value="category.id" />
                  </el-select>
                </el-form-item>
              </el-col>
              <el-col :span="18">
                <el-form-item>
                  <el-select v-model="currentNote.tag_ids" multiple placeholder="选择标签" clearable>
                    <el-option v-for="tag in tags" :key="tag.id" :label="tag.name" :value="tag.id">
                      <el-tag :color="tag.color">{{ tag.name }}</el-tag>
                    </el-option>
                  </el-select>
                </el-form-item>
              </el-col>
            </el-row>
          </el-form>

          <el-row :gutter="20">
            <el-col :span="12">
              <div class="editor-section">
                <el-input type="textarea" v-model="currentNote.content" :rows="20"
                  placeholder="请输入笔记内容（支持Markdown格式）" />
              </div>
            </el-col>
            <el-col :span="12">
              <div class="preview-section">
                <div class="markdown-body" v-html="renderedContent"></div>
              </div>
            </el-col>
          </el-row>
        </el-card>
      </el-col>
    </el-row>

    <NoteShare ref="shareDialog" :note-id="currentNote.id" />
    <NoteImportExport ref="importExportDialog" :note-id="currentNote.id" />
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { marked } from 'marked'
import { ElMessage, ElMessageBox } from 'element-plus'
import dayjs from 'dayjs'
import { ArrowLeft, ArrowRight, DocumentAdd, Upload, Download } from '@element-plus/icons-vue'
import 'highlight.js/styles/github.css'
import CategoryManager from '@/components/CategoryManager.vue'
import TagManager from '@/components/TagManager.vue'
import NoteSearch from '@/components/NoteSearch.vue'
import NoteShare from '@/components/NoteShare.vue'
import NoteImportExport from '@/components/NoteImportExport.vue'

const searchQuery = ref('')
const notes = ref([])
const activeNoteId = ref('')
const activeTab = ref('list')
const categories = ref([])
const tags = ref([])
const shareDialog = ref(null)
const importExportDialog = ref(null)
const isSidebarVisible = ref(true)

const currentNote = ref({
  id: '',
  title: '',
  content: '',
  category_id: '',
  tag_ids: []
})

const renderedContent = computed(() => {
  return marked(currentNote.value.content)
})

const filteredNotes = computed(() => {
  if (!searchQuery.value) return notes.value
  return notes.value.filter(note =>
    note.title.toLowerCase().includes(searchQuery.value.toLowerCase())
  )
})

const formatDate = (date) => {
  return dayjs(date).format('YYYY-MM-DD HH:mm')
}

const toggleSidebar = () => {
  isSidebarVisible.value = !isSidebarVisible.value
}

const fetchNotes = async () => {
  try {
    const token = localStorage.getItem('token')
    const response = await fetch('/api/notes', {
      headers: {
        'Authorization': `Bearer ${token}`
      }
    })
    if (!response.ok) throw new Error('获取笔记列表失败')
    notes.value = await response.json()
    if (notes.value.length > 0) {
      handleNoteSelect(notes.value[0].id)
    }
  } catch (error) {
    ElMessage.error(error.message)
  }
}

const fetchCategories = async () => {
  try {
    const token = localStorage.getItem('token')
    const response = await fetch('/api/categories', {
      headers: {
        'Authorization': `Bearer ${token}`
      }
    })
    if (!response.ok) throw new Error('获取分类列表失败')
    categories.value = await response.json()
  } catch (error) {
    ElMessage.error(error.message)
  }
}

const fetchTags = async () => {
  try {
    const token = localStorage.getItem('token')
    const response = await fetch('/api/tags', {
      headers: {
        'Authorization': `Bearer ${token}`
      }
    })
    if (!response.ok) throw new Error('获取标签列表失败')
    tags.value = await response.json()
  } catch (error) {
    ElMessage.error(error.message)
  }
}

const handleNoteSelect = (noteId) => {
  const note = notes.value.find(n => n.id === noteId)
  if (note) {
    activeNoteId.value = noteId
    currentNote.value = { ...note }
  }
}

const createNewNote = () => {
  currentNote.value = {
    id: '',
    title: '新笔记',
    content: '',
    category_id: '',
    tag_ids: []
  }
  activeNoteId.value = ''
}

const saveNote = async () => {
  try {
    const token = localStorage.getItem('token')
    const url = currentNote.value.id ? `/api/notes/${currentNote.value.id}` : '/api/notes'
    const method = currentNote.value.id ? 'PUT' : 'POST'
    //alert(JSON.stringify(currentNote.value))

    const response = await fetch(url, {
      method,
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${token}`
      },
      body: JSON.stringify(currentNote.value)
    })

    if (!response.ok) throw new Error('保存笔记失败')

    ElMessage.success('保存成功')
    await fetchNotes()
  } catch (error) {
    ElMessage.error(error.message)
  }
}

const deleteNote = async () => {
  if (!currentNote.value.id) return

  try {
    await ElMessageBox.confirm('确定要删除这篇笔记吗？', '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })

    const token = localStorage.getItem('token')
    const response = await fetch(`/api/notes/${currentNote.value.id}`, {
      method: 'DELETE',
      headers: {
        'Authorization': `Bearer ${token}`
      }
    })

    if (!response.ok) throw new Error('删除笔记失败')

    ElMessage.success('删除成功')
    createNewNote()
    await fetchNotes()
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error(error.message)
    }
  }
}

const showShareDialog = () => {
  if (shareDialog.value) {
    shareDialog.value.show()
  }
}

const showImportExport = () => {
  if (importExportDialog.value) {
    importExportDialog.value.show()
  }
}

onMounted(() => {
  fetchNotes()
  fetchCategories()
  fetchTags()
})
</script>

<style scoped>
.note-list-container {
  padding: 20px;
  height: calc(100vh - 60px);
}

.note-list-card {
  height: 100%;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.action-buttons {
  display: flex;
  gap: 10px;
  align-items: center;
}

.search-input {
  margin-bottom: 10px;
}

.note-menu {
  height: calc(100% - 100px);
  overflow-y: auto;
}

.note-item {
  display: flex;
  flex-direction: column;
}

.note-title {
  font-size: 14px;
  margin-bottom: 4px;
}

.note-date {
  font-size: 12px;
  color: #909399;
}

.editor-card {
  height: 100%;
}

.note-title-input {
  width: 300px;
  margin: 0 10px;
}

.editor-section,
.preview-section {
  height: calc(100vh - 300px);
  display: flex;
  flex-direction: column;
}

.editor-section h3,
.preview-section h3 {
  margin: 0 0 10px 0;
  font-size: 16px;
  color: #606266;
}

.editor-section :deep(.el-textarea__inner) {
  height: 100%;
  resize: none;
}

.preview-section {
  border: 1px solid #dcdfe6;
  border-radius: 4px;
  padding: 10px;
  overflow-y: auto;
}

.markdown-body {
  height: 100%;
  overflow-y: auto;
}

:deep(.el-card__header) {
  padding: 10px 20px;
}

:deep(.el-card__body) {
  padding: 10px;
}
</style>