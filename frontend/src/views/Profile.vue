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
            <label>余额:</label>
            <span>￥{{ (Number(user.balance_cents||0)/100).toFixed(2) }}</span>
          </div>
          
        </div>
      </el-card>

      <el-card class="favorites-card">
        <template #header>
          <div class="card-header">
            <span>我的收藏</span>
          </div>
        </template>
        <el-table :data="favList" style="width: 100%; margin-top: 1rem">
          <el-table-column label="封面" width="160">
            <template #default="scope">
              <template v-if="scope.row.cover_url">
                <img :src="isAbsolute(scope.row.cover_url) ? scope.row.cover_url : coverSrc(scope.row.cover_url)" style="width:120px;height:72px;object-fit:cover;border-radius:4px" />
              </template>
              <span v-else>无封面</span>
            </template>
          </el-table-column>
          <el-table-column prop="title" label="标题" />
          <el-table-column label="价格">
            <template #default="scope">
              {{ (Number(scope.row.price_cents) / 100).toFixed(2) }}
            </template>
          </el-table-column>
          <el-table-column label="操作" width="160">
            <template #default="scope">
              <el-button type="primary" link @click="$router.push('/models/'+scope.row.id)">查看</el-button>
              <el-button type="danger" link @click="removeFavoriteModel(scope.row.id)">取消收藏</el-button>
            </template>
          </el-table-column>
        </el-table>
      </el-card>

      <el-card class="my-models-card" style="margin-top: 1.5rem">
        <template #header>
          <div class="card-header">
            <span>我的作品</span>
            <el-button type="primary" style="float:right" @click="uploadDialogVisible = true">上传我的作品</el-button>
          </div>
        </template>
        <el-table :data="myList" v-loading="myLoading" style="width:100%">
          <el-table-column prop="title" label="标题" />
          <el-table-column label="价格">
            <template #default="scope">
              {{ (Number(scope.row.price_cents) / 100).toFixed(2) }}
            </template>
          </el-table-column>
          <el-table-column label="封面" width="160">
            <template #default="scope">
              <template v-if="scope.row.cover_url">
                <img :src="isAbsolute(scope.row.cover_url) ? scope.row.cover_url : coverSrc(scope.row.cover_url)" style="width:120px;height:72px;object-fit:cover;border-radius:4px" />
              </template>
              <span v-else>无封面</span>
            </template>
          </el-table-column>
          <el-table-column label="状态" width="120">
            <template #default="scope">
              {{ statusText(scope.row.status) }}
            </template>
          </el-table-column>
          <el-table-column label="操作" width="220">
            <template #default="scope">
              <el-button size="small" type="success" @click="setStatus(scope.row.id, 'listed')">上架</el-button>
              <el-button size="small" @click="setStatus(scope.row.id, 'delisted')">下架</el-button>
            </template>
          </el-table-column>
        </el-table>
        <div style="margin-top: 12px; display:flex; justify-content:center">
          <el-pagination layout="prev, pager, next" :page-size="mySize" :current-page="myPage" @current-change="changeMyPage" />
        </div>
      </el-card>

      <el-dialog v-model="uploadDialogVisible" title="上传我的作品" width="520px">
        <el-form :model="upload" label-width="100px">
          <el-form-item label="标题">
            <el-input v-model="upload.title" />
          </el-form-item>
          <el-form-item label="描述">
            <el-input type="textarea" v-model="upload.description" />
          </el-form-item>
          <el-form-item label="价格(分)">
            <el-input v-model.number="upload.price_cents" />
          </el-form-item>
          <el-form-item label="模型文件">
            <input type="file" @change="onFileChange" />
          </el-form-item>
          <el-form-item label="封面图片">
            <input type="file" @change="onCoverChange" />
          </el-form-item>
        </el-form>
        <template #footer>
          <el-button @click="uploadDialogVisible = false">取消</el-button>
          <el-button type="primary" :loading="uploading" @click="submitUpload">提交</el-button>
        </template>
      </el-dialog>

      <div class="actions">
        <el-button type="primary" @click="$router.push('/')">返回首页</el-button>
        <el-button @click="handleLogout">退出登录</el-button>
      </div>
    </div>
  </div>
</template>

<script>
import { getProfile } from '@/api/user'
import { logout } from '@/api/auth'
import { getMyModels, updateStatus, uploadModel } from '@/api/models'
import { listFavorites, removeFavorite } from '@/api/favorites'


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
      favList: [],
      myList: [],
      myLoading: false,
      myPage: 1,
      mySize: 10,
      uploadDialogVisible: false,
      uploading: false,
      upload: { title: '', description: '', price_cents: null, file: null, cover: null }
    }
  },
  async created() {
    await this.loadProfile()
    await this.loadFavorites()
    await this.fetchMyModels()
  },
  methods: {
    async loadProfile() {
      try {
        const response = await getProfile()
        this.user = response.data
      } catch (error) {
        this.$message.error('请先登录后查看个人中心')
        this.$router.push('/login')
      }
    },
    async loadFavorites() {
      const res = await listFavorites()
      this.favList = res.data.list || []
    },
    async fetchMyModels() {
      this.myLoading = true
      try {
        const res = await getMyModels({ page: this.myPage, size: this.mySize })
        this.myList = res.data.list || []
      } finally { this.myLoading = false }
    },
    changeMyPage(p) { this.myPage = p; this.fetchMyModels() },
    async setStatus(id, status) {
      await updateStatus(id, status)
      this.$message.success(status === 1 ? '已上架' : (status === 2 ? '已下架' : '状态已更新'))
      this.fetchMyModels()
    },
    statusText(s) {
      if (s === 1) return '上架'
      if (s === 2) return '下架'
      return '待定'
    },
    onFileChange(e) { this.upload.file = e.target.files?.[0] || null },
    onCoverChange(e) { this.upload.cover = e.target.files?.[0] || null },
    async submitUpload() {
      const { title, price_cents, file } = this.upload
      if (!title || !price_cents || !file) { this.$message.error('标题、价格、模型文件为必填'); return }
      const fd = new FormData()
      fd.append('title', title)
      fd.append('description', this.upload.description || '')
      fd.append('price_cents', String(price_cents))
      fd.append('file', file)
      if (this.upload.cover) fd.append('cover', this.upload.cover)
      this.uploading = true
      try {
        await uploadModel(fd)
        this.$message.success('上传成功，作品为草稿状态')
        this.uploadDialogVisible = false
        this.upload = { title: '', description: '', price_cents: null, file: null, cover: null }
        this.fetchMyModels()
      } catch (e) {
        this.$message.error('上传失败')
      } finally { this.uploading = false }
    },
    isAbsolute(u) { return /^https?:\/\//i.test(u) },
    coverLink(u) { return this.isAbsolute(u) ? u : u },
    coverSrc(c) {
      if (!c) return ''
      const norm = String(c).replace(/\\/g, '/')
      return `/api/static?file=${encodeURIComponent(norm)}`
    },
    async removeFavoriteModel(id) {
      await removeFavorite(id)
      this.$message.success('已取消收藏')
      await this.loadFavorites()
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
    
