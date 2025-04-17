<template>
  <div class="category-manager">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>分类管理</span>
          <el-button type="primary" @click="showAddDialog">添加分类</el-button>
        </div>
      </template>
      
      <el-table :data="categories" style="width: 100%">
        <el-table-column prop="name" label="分类名称" />
        <el-table-column prop="noteCount" label="笔记数量" width="100" />
        <el-table-column label="操作" width="200">
          <template #default="{ row }">
            <el-button type="primary" link @click="editCategory(row)">编辑</el-button>
            <el-button type="danger" link @click="deleteCategory(row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>
    
    <el-dialog
      v-model="dialogVisible"
      :title="isEdit ? '编辑分类' : '添加分类'"
      width="30%"
    >
      <el-form :model="categoryForm" label-width="80px">
        <el-form-item label="分类名称">
          <el-input v-model="categoryForm.name" />
        </el-form-item>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="dialogVisible = false">取消</el-button>
          <el-button type="primary" @click="saveCategory">确定</el-button>
        </span>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'

const categories = ref([])
const dialogVisible = ref(false)
const isEdit = ref(false)
const categoryForm = ref({
  id: '',
  name: ''
})

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

const showAddDialog = () => {
  isEdit.value = false
  categoryForm.value = {
    id: '',
    name: ''
  }
  dialogVisible.value = true
}

const editCategory = (category) => {
  isEdit.value = true
  categoryForm.value = { ...category }
  dialogVisible.value = true
}

const saveCategory = async () => {
  try {
    const token = localStorage.getItem('token')
    const url = isEdit.value ? `/api/categories/${categoryForm.value.id}` : '/api/categories'
    const method = isEdit.value ? 'PUT' : 'POST'
    
    const response = await fetch(url, {
      method,
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${token}`
      },
      body: JSON.stringify(categoryForm.value)
    })
    
    if (!response.ok) throw new Error('保存分类失败')
    
    ElMessage.success('保存成功')
    dialogVisible.value = false
    await fetchCategories()
  } catch (error) {
    ElMessage.error(error.message)
  }
}

const deleteCategory = async (category) => {
  try {
    await ElMessageBox.confirm('确定要删除这个分类吗？', '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })
    
    const token = localStorage.getItem('token')
    const response = await fetch(`/api/categories/${category.id}`, {
      method: 'DELETE',
      headers: {
        'Authorization': `Bearer ${token}`
      }
    })
    
    if (!response.ok) throw new Error('删除分类失败')
    
    ElMessage.success('删除成功')
    await fetchCategories()
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error(error.message)
    }
  }
}

onMounted(() => {
  fetchCategories()
})
</script>

<style scoped>
.category-manager {
  padding: 20px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}
</style> 