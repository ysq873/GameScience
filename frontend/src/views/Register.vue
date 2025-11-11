<template>
  <div class="register-container">
    <div class="register-form">
      <h2>用户注册</h2>
      <el-form :model="form" :rules="rules" ref="registerForm">
        <el-form-item prop="email">
          <el-input
            v-model="form.email"
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
        <el-form-item prop="confirmPassword">
          <el-input
            v-model="form.confirmPassword"
            type="password"
            placeholder="请确认密码"
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
            @click="handleRegister"
            style="width: 100%"
          >
            {{ loading ? '注册中...' : '注册' }}
          </el-button>
        </el-form-item>
      </el-form>
      <div class="form-links">
        <router-link to="/login">已有账号？去登录</router-link>
      </div>
    </div>
  </div>
</template>

<script>
import { register, initRegistrationFlow } from '@/api/auth'

export default {
  name: 'Register',
  data() {
    const validatePass2 = (rule, value, callback) => {
      if (value === '') {
        callback(new Error('请再次输入密码'))
      } else if (value !== this.form.password) {
        callback(new Error('两次输入密码不一致!'))
      } else {
        callback()
      }
    }

    return {
      form: {
        email: '',
        password: '',
        confirmPassword: '',
      },
      rules: {
        email: [
          { required: true, message: '请输入邮箱地址', trigger: 'blur' },
          { type: 'email', message: '请输入正确的邮箱地址', trigger: 'blur' }
        ],
        password: [
          { required: true, message: '请输入密码', trigger: 'blur' },
          { min: 6, message: '密码长度不能少于6位', trigger: 'blur' }
        ],
        confirmPassword: [
          { required: true, message: '请确认密码', trigger: 'blur' },
          { validator: validatePass2, trigger: 'blur' }
        ],
      },
      loading: false,
      flowId: null,
      csrf_token: null
    }
  },

  async created() {
    await this.initRegistrationFlow()
  },

  methods: {
    async initRegistrationFlow() {
      try {
        const response = await initRegistrationFlow()
        this.flowId = response.data.id
        this.csrf_token = response.data.ui.nodes.find(
          (n) => n.attributes.name === 'csrf_token'
        )?.attributes?.value
        console.log('✅ Registration flow initialized:', this.flowId)
      } catch (error) {
        console.error(error)
        this.$message.error('初始化注册流程失败')
      }
    },

    async handleRegister() {
      this.$refs.registerForm.validate(async (valid) => {
        if (!valid) return

        this.loading = true
        try {
          // ✅ 按 Kratos 规范构造请求体
          const payload = {
            method: 'password',
            traits: {
              email: this.form.email
            },
            password: this.form.password,
            csrf_token: this.csrf_token
          }

          const response = await register(this.flowId, payload)

          this.$message.success('注册成功')
          // 存储用户信息到 store
          this.$store.commit('SET_USER', {
              id: response.data.session.identity.id,
              email: response.data.session.identity.traits.email,
              name: response.data.session.identity.traits.name,
              session: response.data.session
            })
          this.$router.push('/profile')

        } catch (error) {
          console.error(error)
          this.$message.error('注册失败：' + (error.response?.data?.error?.reason || error.message))
        } finally {
          this.loading = false
        }
      })
    }
  }
}
</script>

<style scoped>
.register-container {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 100vh;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
}

.register-form {
  background: white;
  padding: 2.5rem;
  border-radius: 12px;
  box-shadow: 0 10px 30px rgba(0, 0, 0, 0.2);
  width: 450px;
}

.register-form h2 {
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
