<template>
  <div id="app">
    <nav class="navbar" v-if="$route.meta.showNav">
      <div class="nav-container">
        <router-link to="/" class="nav-brand">用户认证系统</router-link>
        <div class="nav-links">
          <router-link to="/profile" v-if="user.id">个人中心</router-link>
          <router-link to="/login" v-if="!user.id">登录</router-link>
          <router-link to="/register" v-if="!user.id">注册</router-link>
          <el-button @click="handleLogout" v-if="user.id">退出</el-button>
        </div>
      </div>
    </nav>
    <main>
      <router-view />
    </main>
  </div>
</template>

<script>
import { mapState } from 'vuex'
import { logout } from '@/api/auth'


export default {
  name: 'App',
  computed: {
    ...mapState(['user'])
  },
  methods: {
    async handleLogout() {
      this.$confirm('确定要退出登录吗？', '提示', {
        type: 'warning'
      }).then(async () => {
        try {
          const res = await logout()

        // 拿到 Kratos 返回的 logout_url
        const logoutUrl = res.data.logout_url
        console.log(logoutUrl)
        if (logoutUrl) {
          // 清除本地用户状态
          this.$store.commit('CLEAR_USER')

          window.location.replace(logoutUrl)
          this.$store.commit('CLEAR_USER')
        } 
      }
        catch (error) {
          this.$message.error('退出登录失败')
          this.$store.commit('CLEAR_USER')
        }
      })
    }
  }
}
</script>

<style>
* {
  margin: 0;
  padding: 0;
  box-sizing: border-box;
}

body {
  font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
  background-color: #f5f5f5;
  color: #333;
}

.navbar {
  background: #fff;
  box-shadow: 0 2px 4px rgba(0,0,0,0.1);
  padding: 1rem 0;
}

.nav-container {
  max-width: 1200px;
  margin: 0 auto;
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 0 2rem;
}

.nav-brand {
  font-size: 1.5rem;
  font-weight: bold;
  color: #007bff;
  text-decoration: none;
}

.nav-links {
  display: flex;
  gap: 2rem;
}

.nav-links a {
  color: #333;
  text-decoration: none;
  padding: 0.5rem 1rem;
  border-radius: 4px;
  transition: background-color 0.3s;
}

.nav-links a:hover {
  background-color: #f8f9fa;
}

.nav-links a.router-link-active {
  color: #007bff;
  background-color: #e7f3ff;
}

main {
  min-height: calc(100vh - 80px);
}
</style>