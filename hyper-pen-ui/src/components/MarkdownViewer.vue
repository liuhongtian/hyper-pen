<template>
  <div class="markdown-viewer" v-html="renderedContent"></div>
</template>

<script setup>
import { computed } from 'vue'
import { marked } from 'marked'
import DOMPurify from 'dompurify'

const props = defineProps({
  content: {
    type: String,
    required: true
  }
})

// 配置 marked
marked.setOptions({
  breaks: true,
  gfm: true
})

// 渲染 Markdown 内容
const renderedContent = computed(() => {
  const html = marked(props.content)
  return DOMPurify.sanitize(html)
})
</script>

<style scoped>
.markdown-viewer {
  line-height: 1.6;
  
  :deep(h1) {
    font-size: 2em;
    margin: 0.67em 0;
  }
  
  :deep(h2) {
    font-size: 1.5em;
    margin: 0.83em 0;
  }
  
  :deep(h3) {
    font-size: 1.17em;
    margin: 1em 0;
  }
  
  :deep(p) {
    margin: 1em 0;
  }
  
  :deep(ul), :deep(ol) {
    margin: 1em 0;
    padding-left: 2em;
  }
  
  :deep(blockquote) {
    margin: 1em 0;
    padding: 0 1em;
    border-left: 4px solid #ddd;
    color: #666;
  }
  
  :deep(code) {
    padding: 0.2em 0.4em;
    background-color: #f6f8fa;
    border-radius: 3px;
    font-family: monospace;
  }
  
  :deep(pre) {
    padding: 1em;
    background-color: #f6f8fa;
    border-radius: 3px;
    overflow-x: auto;
    
    code {
      padding: 0;
      background-color: transparent;
    }
  }
  
  :deep(table) {
    border-collapse: collapse;
    margin: 1em 0;
    width: 100%;
    
    th, td {
      padding: 0.5em;
      border: 1px solid #ddd;
    }
    
    th {
      background-color: #f6f8fa;
    }
  }
}
</style> 