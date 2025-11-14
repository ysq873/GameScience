<template>
  <div class="login-container">
    <div class="login-form">
      <h2>用户登录</h2>
      <el-form :model="form" :rules="rules" ref="loginForm" @submit.prevent>
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
        <el-form-item>
          <el-button
            type="default"
            size="large"
            :loading="loading"
            @click="handleLoginGithub"
            style="width: 100%"
          >
            使用 GitHub 登录
          </el-button>
        </el-form-item>
      </el-form>
      <div class="form-links">
        <router-link to="/register">没有账号？去注册&nbsp;&nbsp;&nbsp;&nbsp;</router-link>
        <router-link to="/reset-password">忘记密码？</router-link>
      </div>
      <div v-if="alreadyLogged" class="logged-tip">
        <el-alert title="当前处于恢复账号流程，请先退出或去重置密码" type="info" show-icon />
        <div style="margin-top: 10px; display: flex; gap: 12px; justify-content: center;">
          <el-button @click="$router.push({ path: '/reset-password', query: { from: 'recovery' } })" type="primary">去设置新密码</el-button>
          <el-button @click="logoutAndRedirect" type="warning">退出登录</el-button>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { initLoginFlow, getLoginFlow, submitLogin, getSession, startOidc } from '@/api/auth'
import { Message, Lock } from '@element-plus/icons-vue'

export default {
  name: 'Login',
  components: { Message, Lock },
  data() {
    return {
      flowId: '',
      actionUrl: '',
      csrf_token: '',
      form: {
        identifier: '',
        password: '',
        method: 'password',
        // csrf_token 会在提交前动态塞入
      },
      rules: {
        identifier: [
          { required: true, message: '请输入邮箱地址', trigger: 'blur' },
          { type: 'email', message: '请输入正确的邮箱地址', trigger: 'blur' }
        ],
        password: [{ required: true, message: '请输入密码', trigger: 'blur' }]
      },
      loading: false
      , alreadyLogged: false
    }
  },
  async created() {
    const logged = await this.checkSession()
    if (!logged) {
      await this.prepareLoginFlow()
    }
  },
  methods: {
    isRecoveryCtx() {
      return this.$route.query.from === 'recovery'
    },
    async checkSession() {
      try {
        const me = await getSession()
        if (me?.data?.identity?.id) {
          const methods = Array.isArray(me.data.authentication_methods) ? me.data.authentication_methods : []
          const isRecoverySession = methods.some(m => m.method === 'code_recovery')
          if (isRecoverySession) {
            this.alreadyLogged = true
          } else {
            this.$store.commit('SET_USER', {
              id: me.data.identity.id,
              email: me.data.identity.traits?.email,
              name: me.data.identity.traits?.name,
              session: me.data
            })
            this.$router.replace('/models')
          }
          return true
        }
      } catch (e) {}
      return false
    },
    async prepareLoginFlow () {
    try {
      const urlFlow = this.$route.query.flow // 这里读 URL 上是否已有 flow
      if (urlFlow) {
        // 已有 flow，直接读取
        const { data } = await getLoginFlow(urlFlow)
        this.flowId   = data.id
        this.actionUrl = data.ui.action
        this.csrf_token = data.ui.nodes.find(n => n.attributes?.name === 'csrf_token')?.attributes?.value || ''
        return
      }

      // 没有 flow：创建一个新的（Browser 端点 + Accept: application/json + withCredentials）
      const { data } = await initLoginFlow()
      this.flowId    = data.id
      this.actionUrl = data.ui.action
      this.csrf_token = data.ui.nodes.find(n => n.attributes?.name === 'csrf_token')?.attributes?.value || ''

      // 关键：把 flow 写回 URL，之后刷新/回退都能从 URL 复用
      this.$router.replace({
        path: this.$route.path,
        query: { ...this.$route.query, flow: this.flowId }
      })
    } catch (e) {
      // flow 过期或无效时，重新创建
      await this.recreateLoginFlow()
    }
  },
  async recreateLoginFlow () {
    const { data } = await initLoginFlow()
    this.flowId     = data.id
    this.actionUrl  = data.ui.action
    this.csrf_token = data.ui.nodes.find(n => n.attributes?.name === 'csrf_token')?.attributes?.value || ''
    this.$router.replace({ path: this.$route.path, query: { ...this.$route.query, flow: this.flowId } })
  },

    async handleLogin() {
      this.$refs.loginForm.validate(async (valid) => {
        if (!valid || !this.flowId || !this.actionUrl) return
        this.loading = true
        try {
          const payload = { ...this.form }
          if (this.csrf_token) payload.csrf_token = this.csrf_token

          // **关键：提交到 ui.action，而不是手拼 URL**
          const resp = await submitLogin(this.flowId, payload)

          // Browser Flow 不一定回 session，这里统一以 whoami 为准
          const me = await getSession()

          if (me?.data) {
            this.$message.success('登录成功')
            this.$store.commit('SET_USER', {
              id: me.data.identity.id,
              email: me.data.identity.traits?.email,
              name: me.data.identity.traits?.name,
              session: me.data
            })
            this.$router.push('/models')
          } else {
            this.$message.error('登录失败，请重试')
            await this.prepareLoginFlow()
          }
        } catch (error) {
          console.error('Login error:', error)
          const uiMsg = error.response?.data?.ui?.messages?.[0]?.text
          const fieldMsg = error.response?.data?.ui?.nodes?.flatMap(n => n.messages || [])?.[0]?.text
          this.$message.error('登录失败：' + (uiMsg || fieldMsg || '请检查输入'))
          await this.prepareLoginFlow() // 处理 flow 过期等情况
        } finally {
          this.loading = false
        }
      })
    }
    ,
    async logoutAndRedirect() {
      try {
        const res = await logout()
        const logoutUrl = res.data.logout_url
        if (logoutUrl) {
          this.$store.commit('CLEAR_USER')
          window.location.replace(logoutUrl)
        } else {
          this.$message.error('退出登录失败')
        }
      } catch(e) {
        this.$message.error('退出登录失败')
      }
    },
    async handleLoginGithub() {
      if (!this.flowId || !this.csrf_token) {
        await this.prepareLoginFlow()
      }
      this.loading = true
      try {
        const resp = await startOidc(this.flowId, 'github', this.csrf_token)
        const redirectUrl = resp?.data?.redirect_browser_to || resp?.data?.continue_with?.find?.(i => i.action === 'redirect_browser_to')?.redirect_browser_to
        if (redirectUrl) {
          window.location.href = redirectUrl
          return
        }
        this.$message.error('未获取到跳转地址')
      } catch (error) {
        if (error?.code === 'ERR_CANCELED' || /browser location has changed/i.test(error?.message || '')) {
          return
        }
        const redirectUrl = error.response?.data?.redirect_browser_to || error.response?.data?.continue_with?.find?.(i => i.action === 'redirect_browser_to')?.redirect_browser_to
        if (redirectUrl) {
          window.location.href = redirectUrl
        } else {
          const uiMsg = error.response?.data?.ui?.messages?.[0]?.text
          const fieldMsg = error.response?.data?.ui?.nodes?.flatMap(n => n.messages || [])?.[0]?.text
          this.$message.error('GitHub 登录失败：' + (uiMsg || fieldMsg || error.message))
          await this.prepareLoginFlow()
        }
      } finally {
        this.loading = false
      }
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

.logged-tip {
  margin-top: 12px;
  text-align: center;
}
</style>
