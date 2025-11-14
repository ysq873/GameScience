<template>
  <div class="container">
    <h2>已购模型库</h2>
    <el-table :data="list" v-loading="loading">
      <el-table-column prop="ModelId" label="模型ID" />
      <el-table-column label="操作" width="220">
        <template #default="scope">
          <el-button type="primary" @click="download(scope.row.ModelId)">下载</el-button>
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
    }
  }
}
</script>

<style>
.container { max-width: 960px; margin: 24px auto; }
</style>
