<template>
  <div class="tag-manager">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>标签管理</span>
          <el-button type="primary" @click="showAddDialog">添加标签</el-button>
        </div>
      </template>
      
      <el-table :data="tags" style="width: 100%">
        <el-table-column prop="name" label="标签名称" />
        <el-table-column prop="color" label="颜色">
          <template #default="{ row }">
            <el-tag :color="row.color">{{ row.name }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="noteCount" label="笔记数量" width="100" />
        <el-table-column label="操作" width="200">
          <template #default="{ row }">
            <el-button type="primary" link @click="editTag(row)">编辑</el-button>
            <el-button type="danger" link @click="deleteTag(row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>
    
    <el-dialog
      v-model="dialogVisible"
      :title="isEdit ? '编辑标签' : '添加标签'"
      width="30%"
    >
      <el-form :model="tagForm" label-width="80px">
        <el-form-item label="标签名称">
          <el-input v-model="tagForm.name" />
        </el-form-item>
        <el-form-item label="标签颜色">
          <el-color-picker v-model="tagForm.color" />
        </el-form-item>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="dialogVisible = false">取消</el-button>
          <el-button type="primary" @click="saveTag">确定</el-button>
        </span>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'

const tags = ref([])
const dialogVisible = ref(false)
const isEdit = ref(false)
const tagForm = ref({
  id: '',
  name: '',
  color: '#409EFF'
})

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

const showAddDialog = () => {
  isEdit.value = false
  tagForm.value = {
    id: '',
    name: '',
    color: '#409EFF'
  }
  dialogVisible.value = true
}

const editTag = (tag) => {
  isEdit.value = true
  tagForm.value = { ...tag }
  dialogVisible.value = true
}

const saveTag = async () => {
  try {
    const token = localStorage.getItem('token')
    const url = isEdit.value ? `/api/tags/${tagForm.value.id}` : '/api/tags'
    const method = isEdit.value ? 'PUT' : 'POST'
    
    const response = await fetch(url, {
      method,
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${token}`
      },
      body: JSON.stringify(tagForm.value)
    })
    
    if (!response.ok) throw new Error('保存标签失败')
    
    ElMessage.success('保存成功')
    dialogVisible.value = false
    await fetchTags()
  } catch (error) {
    ElMessage.error(error.message)
  }
}

const deleteTag = async (tag) => {
  try {
    await ElMessageBox.confirm('确定要删除这个标签吗？', '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })
    
    const token = localStorage.getItem('token')
    const response = await fetch(`/api/tags/${tag.id}`, {
      method: 'DELETE',
      headers: {
        'Authorization': `Bearer ${token}`
      }
    })
    
    if (!response.ok) throw new Error('删除标签失败')
    
    ElMessage.success('删除成功')
    await fetchTags()
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error(error.message)
    }
  }
}

onMounted(() => {
  fetchTags()
})
</script>

<style scoped>
.tag-manager {
  padding: 20px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}
</style> 