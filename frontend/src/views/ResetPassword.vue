<template>
  <div class="reset-container">
    <div class="reset-form">
      <h2>重置密码</h2>
      <el-form :model="form" :rules="rules" ref="resetForm">
        <el-form-item prop="email">
          <el-input
            v-model="form.email"
            placeholder="请输入注册邮箱"
            size="large"
            prefix-icon="Message"
          />
        </el-form-item>
        <el-form-item>
          <el-button
            type="primary"
            size="large"
            :loading="loading"
            @click="handleReset"
            style="width: 100%"
          >
            {{ loading ? '发送中...' : '发送重置邮件' }}
          </el-button>
        </el-form-item>
      </el-form>
      <div class="form-links">
        <router-link to="/login">返回登录</router-link>
      </div>
    </div>
  </div>
</template>

<script>
import { resetPassword } from '@/api/auth'

export default {
  name: 'ResetPassword',
  data() {
    return {
      form: {
        email: ''
      },
      rules: {
        email: [
          { required: true, message: '请输入邮箱地址', trigger: 'blur' },
          { type: 'email', message: '请输入正确的邮箱地址', trigger: 'blur' }
        ]
      },
      loading: false
    }
  },
  methods: {
    async handleReset() {
      this.$refs.resetForm.validate(async (valid) => {
        if (!valid) return
        
        this.loading = true
        try {
          await resetPassword(this.form)
          this.$message.success('重置邮件已发送，请查收邮箱')
          this.$router.push('/login')
        } catch (error) {
          this.$message.error('发送失败：' + (error.response?.data?.message || error.message))
        } finally {
          this.loading = false
        }
      })
    }
  }
}
</script>

<style scoped>
.reset-container {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 100vh;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
}

.reset-form {
  background: white;
  padding: 2.5rem;
  border-radius: 12px;
  box-shadow: 0 10px 30px rgba(0, 0, 0, 0.2);
  width: 400px;
}

.reset-form h2 {
  text-align: center;
  margin-bottom: 2rem;
  color: #333;
  font-size: 1.8rem;
}

.form-links {
  text-align: center;
  margin-top: 1rem;
}

.form-links a {
  color: #666;
  text-decoration: none;
  font-size: 0.9rem;
}

.form-links a:hover {
  color: #007bff;
}
</style>