<template>
  <div class="profile-container">
    <div class="profile-card">
      <h2>个人资料</h2>

      <el-card class="info-card">
        <template #header>
          <div class="card-header">
            <span>基本信息</span>
          </div>
        </template>
        <div class="profile-info">
          <div class="info-item">
            <label>邮箱:</label>
            <span>{{ user.email }}</span>
          </div>
          <div class="info-item">
            <label>姓名:</label>
            <span>{{ user.name.first }} {{ user.name.last }}</span>
          </div>
        </div>
      </el-card>

      <el-card class="favorites-card">
        <template #header>
          <div class="card-header">
            <span>我的收藏</span>
          </div>
        </template>
        <div class="add-favorite">
          <el-input
            v-model="newFavorite"
            placeholder="输入收藏项目"
            style="width: 300px; margin-right: 1rem"
            @keyup.enter="addFavorite"
          />
          <el-button type="primary" @click="addFavorite" :loading="addingFavorite">
            添加收藏
          </el-button>
        </div>
        <el-table :data="user.favorites" style="width: 100%; margin-top: 1rem">
          <el-table-column prop="item" label="收藏项目">
            <template #default="scope">
              {{ scope.row }}
            </template>
          </el-table-column>
          <el-table-column label="操作" width="100">
            <template #default="scope">
              <el-button
                type="danger"
                link
                @click="removeFavorite(scope.$index)"
              >
                删除
              </el-button>
            </template>
          </el-table-column>
        </el-table>
      </el-card>

      <div class="actions">
        <el-button type="primary" @click="$router.push('/')">返回首页</el-button>
        <el-button @click="handleLogout">退出登录</el-button>
      </div>
    </div>
  </div>
</template>

<script>
import { addFavorite, getProfile } from '@/api/user'
import { logout } from '@/api/auth'

export default {
  name: 'Profile',
  data() {
    return {
      user: {
        id: '',
        email: '',
        name: {
          first: '',
          last: ''
        },
        favorites: []
      },
      newFavorite: '',
      addingFavorite: false
    }
  },
  async created() {
    await this.loadProfile()
  },
  methods: {
    async loadProfile() {
      try {
        const response = await getProfile()
        this.user = response.data
      } catch (error) {
        this.$message.error('获取用户资料失败')
      }
    },
    async addFavorite() {
      if (!this.newFavorite.trim()) {
        this.$message.warning('请输入收藏内容')
        return
      }

      this.addingFavorite = true
      try {
        const response = await addFavorite(this.newFavorite.trim())
        this.user = response.data
        this.newFavorite = ''
        this.$message.success('添加收藏成功')
      } catch (error) {
        this.$message.error('添加收藏失败')
      } finally {
        this.addingFavorite = false
      }
    },
    removeFavorite(index) {
      this.$confirm('确定要删除这个收藏吗？', '提示', {
        type: 'warning'
      }).then(() => {
        this.user.favorites.splice(index, 1)
        this.$message.success('删除成功')
      })
    },
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

<style scoped>
.profile-container {
  max-width: 800px;
  margin: 2rem auto;
  padding: 0 1rem;
}

.profile-card h2 {
  text-align: center;
  margin-bottom: 2rem;
  color: #333;
}

.info-card, .favorites-card {
  margin-bottom: 2rem;
}

.card-header {
  font-weight: bold;
  font-size: 1.1rem;
}

.profile-info {
  display: flex;
  flex-direction: column;
  gap: 1rem;
}

.info-item {
  display: flex;
  align-items: center;
}

.info-item label {
  font-weight: bold;
  width: 80px;
  margin-right: 1rem;
}

.add-favorite {
  display: flex;
  align-items: center;
  margin-bottom: 1rem;
}

.actions {
  text-align: center;
  margin-top: 2rem;
}
</style>