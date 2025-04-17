<template>
  <div class="note-search">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>搜索笔记</span>
        </div>
      </template>
      
      <el-form :model="searchForm" label-width="80px">
        <el-form-item label="关键词">
          <el-input
            v-model="searchForm.keyword"
            placeholder="输入关键词搜索"
            clearable
            @keyup.enter="handleSearch"
          />
        </el-form-item>
        
        <el-form-item label="分类">
          <el-select
            v-model="searchForm.categoryId"
            placeholder="选择分类"
            clearable
          >
            <el-option
              v-for="category in categories"
              :key="category.id"
              :label="category.name"
              :value="category.id"
            />
          </el-select>
        </el-form-item>
        
        <el-form-item label="标签">
          <el-select
            v-model="searchForm.tagIds"
            multiple
            placeholder="选择标签"
            clearable
          >
            <el-option
              v-for="tag in tags"
              :key="tag.id"
              :label="tag.name"
              :value="tag.id"
            >
              <el-tag :color="tag.color">{{ tag.name }}</el-tag>
            </el-option>
          </el-select>
        </el-form-item>
        
        <el-form-item label="时间范围">
          <el-date-picker
            v-model="searchForm.dateRange"
            type="daterange"
            range-separator="至"
            start-placeholder="开始日期"
            end-placeholder="结束日期"
            value-format="YYYY-MM-DD"
          />
        </el-form-item>
        
        <el-form-item>
          <el-button type="primary" @click="handleSearch">搜索</el-button>
          <el-button @click="resetSearch">重置</el-button>
        </el-form-item>
      </el-form>
      
      <el-table
        v-loading="loading"
        :data="searchResults"
        style="width: 100%"
        @row-click="handleNoteClick"
      >
        <el-table-column prop="title" label="标题" />
        <el-table-column prop="category.name" label="分类" width="120" />
        <el-table-column label="标签" width="200">
          <template #default="{ row }">
            <el-tag
              v-for="tag in row.tags"
              :key="tag.id"
              :color="tag.color"
              class="tag-item"
            >
              {{ tag.name }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="updatedAt" label="更新时间" width="180">
          <template #default="{ row }">
            {{ formatDate(row.updatedAt) }}
          </template>
        </el-table-column>
      </el-table>
      
      <div class="pagination">
        <el-pagination
          v-model:current-page="currentPage"
          v-model:page-size="pageSize"
          :total="total"
          :page-sizes="[10, 20, 50, 100]"
          layout="total, sizes, prev, pager, next"
          @size-change="handleSizeChange"
          @current-change="handleCurrentChange"
        />
      </div>
    </el-card>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import dayjs from 'dayjs'

const router = useRouter()

const searchForm = ref({
  keyword: '',
  categoryId: '',
  tagIds: [],
  dateRange: []
})

const categories = ref([])
const tags = ref([])
const searchResults = ref([])
const loading = ref(false)
const currentPage = ref(1)
const pageSize = ref(10)
const total = ref(0)

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

const handleSearch = async () => {
  try {
    loading.value = true
    const token = localStorage.getItem('token')
    const params = new URLSearchParams({
      page: currentPage.value,
      pageSize: pageSize.value,
      keyword: searchForm.value.keyword,
      categoryId: searchForm.value.categoryId,
      ...(searchForm.value.tagIds.length > 0 && { tagIds: searchForm.value.tagIds.join(',') }),
      ...(searchForm.value.dateRange?.length === 2 && {
        startDate: searchForm.value.dateRange[0],
        endDate: searchForm.value.dateRange[1]
      })
    })
    
    const response = await fetch(`/api/notes/search?${params}`, {
      headers: {
        'Authorization': `Bearer ${token}`
      }
    })
    
    if (!response.ok) throw new Error('搜索失败')
    
    const data = await response.json()
    searchResults.value = data.items
    total.value = data.total
  } catch (error) {
    ElMessage.error(error.message)
  } finally {
    loading.value = false
  }
}

const resetSearch = () => {
  searchForm.value = {
    keyword: '',
    categoryId: '',
    tagIds: [],
    dateRange: []
  }
  currentPage.value = 1
  handleSearch()
}

const handleSizeChange = (val) => {
  pageSize.value = val
  handleSearch()
}

const handleCurrentChange = (val) => {
  currentPage.value = val
  handleSearch()
}

const handleNoteClick = (row) => {
  router.push(`/notes/${row.id}`)
}

const formatDate = (date) => {
  return dayjs(date).format('YYYY-MM-DD HH:mm')
}

onMounted(() => {
  fetchCategories()
  fetchTags()
  handleSearch()
})
</script>

<style scoped>
.note-search {
  padding: 20px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.tag-item {
  margin-right: 5px;
  margin-bottom: 5px;
}

.pagination {
  margin-top: 20px;
  display: flex;
  justify-content: center;
}
</style> 