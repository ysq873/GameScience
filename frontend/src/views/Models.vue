<template>
  <div class="container">
    <div class="header">
      <h2>全部作品</h2>
      <el-button v-if="$store.state.user.id" type="primary" @click="uploadDialogVisible = true">上传作品</el-button>
    </div>
    <el-table :data="list" style="width: 100%" v-loading="loading">
      <el-table-column prop="id" label="ID" width="80" />
      <el-table-column prop="title" label="标题" />
      <el-table-column prop="price_cents" label="价格" />
      <el-table-column label="操作" width="200">
        <template #default="scope">
          <el-button type="primary" @click="goDetail(scope.row.id)">详情</el-button>
        </template>
      </el-table-column>
    </el-table>
    <div class="pager">
      <el-pagination
        layout="prev, pager, next"
        :total="total"
        :page-size="size"
        :current-page="page"
        @current-change="changePage" />
    </div>
    <el-dialog v-model="uploadDialogVisible" title="上传作品" width="520px">
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
  </div>
  </template>

<script>
import { listModels, uploadModel } from '@/api/models'
import { getSession } from '@/api/auth'

export default {
  name: 'Models',
  data() {
    return { list: [], loading: false, page: 1, size: 10, total: 0,
      uploadDialogVisible: false,
      uploading: false,
      upload: { title: '', description: '', price_cents: null, file: null, cover: null }
    }
  },
  mounted() { this.ensureUser().then(() => this.fetch()) },
  methods: {
    async ensureUser() {
      try {
        const me = await getSession()
        if (me?.data?.identity?.id) {
          this.$store.commit('SET_USER', {
            id: me.data.identity.id,
            email: me.data.identity.traits?.email,
            name: me.data.identity.traits?.name,
            session: me.data
          })
          return
        }
      } catch (e) {}
      this.$router.replace('/login')
    },
    async fetch() {
      this.loading = true
      try {
        const res = await listModels({ page: this.page, size: this.size })
        const data = res.data
        this.list = data.list || []
        this.total = data.total || this.list.length
      } finally { this.loading = false }
    },
    changePage(p) { this.page = p; this.fetch() },
    goDetail(id) { this.$router.push(`/models/${id}`) },
    onFileChange(e) { this.upload.file = e.target.files?.[0] || null },
    onCoverChange(e) { this.upload.cover = e.target.files?.[0] || null },
    async submitUpload() {
      if (!this.$store.state.user.id) { this.$message.error('请先登录'); return }
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
        this.fetch()
      } catch (e) {
        this.$message.error('上传失败')
      } finally { this.uploading = false }
    }
  }
}
</script>

<style>
.container { max-width: 960px; margin: 24px auto; }
.header { display:flex; justify-content: space-between; align-items:center; margin-bottom:12px; }
.pager { margin-top: 16px; display: flex; justify-content: center; }
</style>
