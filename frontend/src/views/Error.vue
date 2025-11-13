<template>
  <div class="error-container">
    <div class="error-card">
      <h2>发生错误</h2>
      <div v-if="loading" class="loading">加载中...</div>
      <div v-else>
        <p v-if="message">{{ message }}</p>
        <el-button type="primary" @click="goHome">返回主页</el-button>
      </div>
    </div>
  </div>
</template>

<script>
import { getError } from '@/api/auth'

export default {
  name: 'Error',
  data() {
    return {
      loading: true,
      message: ''
    }
  },
  async created() {
    const id = this.$route.query.id
    if (!id) {
      this.message = '未知错误'
      this.loading = false
      return
    }
    try {
      const { data } = await getError(id)
      this.message = data?.error?.message || '发生未知错误'
    } catch (e) {
      this.message = e.response?.data?.error?.message || e.message || '加载错误信息失败'
    } finally {
      this.loading = false
    }
  },
  methods: {
    goHome() {
      this.$router.push('/')
    }
  }
}
</script>

<style scoped>
.error-container {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 100vh;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
}

.error-card {
  background: white;
  padding: 2.5rem;
  border-radius: 12px;
  box-shadow: 0 10px 30px rgba(0, 0, 0, 0.2);
  width: 450px;
}

.error-card h2 {
  text-align: center;
  margin-bottom: 2rem;
  color: #333;
  font-size: 1.8rem;
}

.loading {
  text-align: center;
  color: #666;
}
</style>
