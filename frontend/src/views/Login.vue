<template>
  <div class="login-container">
    <div class="login-form">
      <h2>用户登录</h2>
      <el-form :model="form" :rules="rules" ref="loginForm">
        <el-form-item prop="identifier">
          <el-input
            v-model="form.identifier"
            placeholder="请输入邮箱"
            size="large"
            prefix-icon="Message"
          />
        </el-form-item>
        <el-form-item prop="password">
          <el-input
            v-model="form.password"
            type="password"
            placeholder="请输入密码"
            size="large"
            prefix-icon="Lock"
            show-password
          />
        </el-form-item>
        <el-form-item>
          <el-button
            type="primary"
            size="large"
            :loading="loading"
            @click="handleLogin"
            style="width: 100%"
          >
            {{ loading ? '登录中...' : '登录' }}
          </el-button>
        </el-form-item>
      </el-form>
      <div class="form-links">
        <router-link to="/register">没有账号？去注册&nbsp;&nbsp;&nbsp;&nbsp;</router-link>
        <router-link to="/reset-password">忘记密码？</router-link>
      </div>
    </div>
  </div>
</template>

<script>
import { initLoginFlow, submitLogin, getSession } from '@/api/auth'

export default {
  name: 'Login',
  data() {
    return {
      flowId: '',
      form: {
        identifier: '',
        password: '',
        method: 'password'
      },
      rules: {
        identifier: [
          { required: true, message: '请输入邮箱地址', trigger: 'blur' },
          { type: 'email', message: '请输入正确的邮箱地址', trigger: 'blur' }
        ],
        password: [
          { required: true, message: '请输入密码', trigger: 'blur' }
        ]
      },
      loading: false
    }
  },
  async created() {
    await this.initLoginFlow()
  },
  methods: {
    async initLoginFlow() {
      try {
        const response = await initLoginFlow()
        this.flowId = response.data.id
        this.csrf_token = response.data.ui.nodes.find(n => n.attributes.name === 'csrf_token')?.attributes?.value
        console.log('Login flow initialized:', this.flowId)
        console.log('Login flow csrf_token:', this.csrf_token)
      } catch (error) {
        this.$message.error('初始化登录流程失败')
      }
    },
    async handleLogin() {
      this.$refs.loginForm.validate(async (valid) => {
        if (!valid || !this.flowId) return
        
        this.loading = true
        try {
          this.form.csrf_token = this.csrf_token
          // 提交登录到 Kratos
          const response = await submitLogin(this.flowId, this.form)
          console.log('Login response:', response)
          
          // 登录成功，Kratos 会设置 session cookie
          if (response.data?.session) {
            this.$message.success('登录成功')
            
            // 存储用户信息到 store
            this.$store.commit('SET_USER', {
              id: response.data.session.identity.id,
              email: response.data.session.identity.traits.email,
              name: response.data.session.identity.traits.name,
              session: response.data.session
            })
            
            // 跳转到个人中心
            this.$router.push('/profile')
          } else {
            this.$message.error('登录失败，请重试')
            await this.initLoginFlow() // 重新初始化流程
          }
        } catch (error) {
          console.error('Login error:', error)
          this.$message.error('登录失败：' + (error.response?.data?.ui?.messages?.[0]?.text || '请检查输入'))
          await this.initLoginFlow() // 重新初始化流程
        } finally {
          this.loading = false
        }
      })
    }
  }
}
</script>

<style scoped>
.login-container {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 100vh;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
}

.login-form {
  background: white;
  padding: 2.5rem;
  border-radius: 12px;
  box-shadow: 0 10px 30px rgba(0, 0, 0, 0.2);
  width: 450px;
}

.login-form h2 {
  text-align: center;
  margin-bottom: 2rem;
  color: #333;
  font-size: 1.8rem;
}

.name-group {
  display: flex;
  gap: 1rem;
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