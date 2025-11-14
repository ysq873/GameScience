<template>
  <div class="container">
    <h2>已购模型库</h2>
    <el-table :data="list" v-loading="loading">
      <el-table-column label="封面" width="160">
        <template #default="scope">
          <template v-if="scope.row.cover_url">
            <img :src="isAbsolute(scope.row.cover_url) ? scope.row.cover_url : coverSrc(scope.row.cover_url)" style="width:120px;height:72px;object-fit:cover;border-radius:4px" />
          </template>
          <span v-else>无封面</span>
        </template>
      </el-table-column>
      <el-table-column prop="title" label="标题" />
      <el-table-column prop="model_id" label="模型ID" width="120" />
      <el-table-column label="操作" width="220">
        <template #default="scope">
          <el-button type="primary" @click="download(scope.row.model_id)">下载</el-button>
        </template>
      </el-table-column>
    </el-table>
  </div>
</template>

<script>
import { listPurchases, generateDownloadToken } from '@/api/downloads'

export default {
  name: 'Library',
  data() { return { list: [], loading: false } },
  mounted() { this.fetch() },
  methods: {
    async fetch() {
      this.loading = true
      try {
        const res = await listPurchases()
        this.list = res.data.list || []
      } finally { this.loading = false }
    },
    async download(modelId) {
      const res = await generateDownloadToken(modelId)
      const token = res.data.token
      window.location.href = `/api/download?token=${token}`
    },
    isAbsolute(u) { return /^https?:\/\//i.test(u) },
    coverSrc(c) { if (!c) return ''; const norm = String(c).replace(/\\/g, '/'); return `/api/static?file=${encodeURIComponent(norm)}` }
  }
}
</script>

<style>
.container { max-width: 960px; margin: 24px auto; }
</style>
