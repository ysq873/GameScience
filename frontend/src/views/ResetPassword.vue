<template>
  <div class="reset-container">
    <div class="reset-form">
      <h2>重置密码</h2>

      <!-- Step 1: 输入邮箱 -->
      <el-form v-if="step === 1" :model="form" :rules="rules" ref="resetForm">
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
            @click="handleSend"
            style="width: 100%"
          >
            {{ loading ? '发送中...' : '发送重置邮件' }}
          </el-button>
        </el-form-item>
      </el-form>

      <!-- Step 2: 输入验证码 -->
      <el-form v-else-if="step === 2" :model="verifyForm" ref="verifyForm">
        <el-form-item>
          <el-input
            v-model="verifyForm.code"
            placeholder="请输入邮箱验证码"
            size="large"
            prefix-icon="Lock"
          />
        </el-form-item>
        <el-form-item>
          <el-button
            type="primary"
            size="large"
            :loading="loading"
            @click="handleVerify"
            style="width: 100%"
          >
            {{ loading ? '验证中...' : '验证验证码' }}
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
import { initRecoveryFlow, recovery } from '@/api/auth'
import axios from 'axios'

export default {
  name: 'ResetPassword',
  data() {
    return {
      step: 1, // 1: 邮箱输入阶段, 2: 验证码输入阶段
      flowId: '',
      actionUrl: '',
      csrf_token: '',
      loading: false,
      form: {
        email: ''
      },
      verifyForm: {
        code: ''
      },
      rules: {
        email: [
          { required: true, message: '请输入邮箱地址', trigger: 'blur' },
          { type: 'email', message: '请输入正确的邮箱地址', trigger: 'blur' }
        ]
      }
    }
  },

  async created() {
    await this.initRecoveryFlow()
  },

  methods: {
    // 初始化恢复流程
    async initRecoveryFlow() {
      try {
        const res = await initRecoveryFlow()
        this.flowId = res.data.id || ''
        this.actionUrl = res.data.ui.action || ''
        this.csrf_token = res.data.ui.nodes.find(
          (n) => n?.attributes.name === 'csrf_token'
        )?.attributes?.value || ''
        if (this.flowId && this.actionUrl) {
          console.log('Recovery flow initialized:', this.flowId)
        } else {
          this.$message.error('初始化重置流程失败')
        }
      } catch (error) {
        this.$message.error('初始化重置流程失败')
      }
    },

    // Step 1: 发送验证码邮件
    async handleSend() {
      this.$refs.resetForm.validate(async (valid) => {
        if (!valid) return
        this.loading = true
      try {
        this.form.csrf_token = this.csrf_token
        this.form.method = 'code'

          await recovery(this.flowId, this.form)
        this.$message.success('验证码已发送到邮箱，请输入验证码')
        this.step = 2 // 进入验证码输入阶段
        } catch (error) {
          this.$message.error('发送失败：' + (error.response?.data?.message || error.message))
        } finally {
          this.loading = false
        }
      })
    },

    // Step 2: 验证验证码
    async handleVerify() {
      if (!this.verifyForm.code) {
        this.$message.error('请输入验证码')
        return
      }
      this.loading = true
      try {
        const response = await recovery(this.flowId, {
          method: 'code',
          code: this.verifyForm.code,
          csrf_token: this.csrf_token
        })
        this.$message.success('验证成功，请设置新密码')
      } catch (error) {
        const redirectUrl = error.response?.data?.redirect_browser_to
        if (redirectUrl) {
          console.log('Kratos 要求跳转到:', redirectUrl)
          window.location.href = redirectUrl
        } else {
          console.error(error)
          this.$message.error(error.response?.data?.error?.message || '验证码错误')
        }
      } finally {
        this.loading = false
      }
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
