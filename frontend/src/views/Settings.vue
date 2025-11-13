<template>
  <div class="NewPassword-container">
    <div class="NewPassword-form">
      <h2>设置新密码</h2>
      <el-form :model="form" :rules="rules" ref="passwordForm">
        <el-form-item prop="password">
          <el-input
            v-model="form.password"
            type="password"
            placeholder="请输入新密码"
            size="large"
            show-password
          />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" size="large" @click="handleReset" style="width: 100%">
            提交新密码
          </el-button>
        </el-form-item>
      </el-form>
    </div>
  </div>
</template>

<script>

import { settings, initSettings, getSession } from '@/api/auth'

export default {
  name: 'Settings',
  data() {
    return {
      form: { password: '' },
      rules: {
        password: [
          { required: true, message: '请输入新密码', trigger: 'blur' },
          { min: 6, message: '密码至少6位', trigger: 'blur' }
        ],
      },
      flowId: '',
      actionUrl: '',
      csrf_token: ''
    }
  },
  async created() {
    this.flowId = new URLSearchParams(window.location.search).get('flow')
    const res = await initSettings(this.flowId)
    this.actionUrl = res.data.ui.action
    this.csrf_token = res.data.ui.nodes.find(n => n.attributes?.name === 'csrf_token')?.attributes?.value || ''
  },
  methods: {
    async handleReset() {
      try {
        if (!this.flowId) return this.$message.error('缺少恢复令牌')
        const resp = await settings(this.flowId, {
          method: 'password',
          password: this.form.password,
          csrf_token: this.csrf_token
        })
        const redirectUrl = (resp && resp.data && (resp.data.redirect_browser_to || (resp.data.continue_with && resp.data.continue_with.find && resp.data.continue_with.find(i => i.action === 'redirect_browser_to')?.redirect_browser_to)))
        if (redirectUrl && !/\/settings(\?|$)/.test(redirectUrl)) {
          window.location.href = redirectUrl
          return
        }
        const me = await getSession().catch(() => null)
        if (me && me.data) {
          this.$message.success('密码重置成功')
          this.$router.push('/profile')
        } else {
          this.$message.success('密码重置成功，请重新登录')
          this.$router.push('/login')
        }
      } catch (error) {
        console.error(error)
        this.$message.error('密码重置失败：' + (error.response?.data?.error?.message || error.message))
      }
    }
  }
}
</script>


<style scoped>
.NewPassword-container {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 100vh;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
}

.NewPassword-form {
  background: white;
  padding: 2.5rem;
  border-radius: 12px;
  box-shadow: 0 10px 30px rgba(0, 0, 0, 0.2);
  width: 400px;
}

.NewPassword-form h2 {
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
