import { createRouter, createWebHistory } from 'vue-router'
import Home from '@/views/Home.vue'
import Login from '@/views/Login.vue'
import Register from '@/views/Register.vue'
import NoteList from '@/views/NoteList.vue'
import SharedNote from '../views/SharedNote.vue'

const routes = [
  {
    path: '/',
    name: 'home',
    component: Home
  },
  {
    path: '/login',
    name: 'login',
    component: Login,
    meta: { requiresGuest: true }
  },
  {
    path: '/register',
    name: 'register',
    component: Register,
    meta: { requiresGuest: true }
  },
  {
    path: '/notes',
    name: 'notes',
    component: NoteList,
    meta: { requiresAuth: true }
  },
  {
    path: '/shared/:token',
    name: 'SharedNote',
    component: SharedNote,
    meta: {
      title: '共享笔记',
      requiresAuth: false
    }
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

// 路由守卫
router.beforeEach((to, from, next) => {
  // 设置页面标题
  document.title = to.meta.title ? `${to.meta.title} - Hyper Pen` : 'Hyper Pen'
  
  const isAuthenticated = !!localStorage.getItem('token')
  
  if (to.meta.requiresAuth && !isAuthenticated) {
    // 需要登录但未登录，重定向到登录页
    next({ name: 'login' })
  } else if (to.meta.requiresGuest && isAuthenticated) {
    // 已登录用户访问登录/注册页，重定向到笔记页
    next({ name: 'notes' })
  } else if (to.name === 'home' && isAuthenticated) {
    // 已登录用户访问首页，重定向到笔记页
    next({ name: 'notes' })
  } else {
    next()
  }
})

export default router 