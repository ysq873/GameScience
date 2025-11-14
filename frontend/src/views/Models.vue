<template>
  <div class="container">
    <h2>全部作品</h2>
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
  </div>
  </template>

<script>
import { listModels } from '@/api/models'

export default {
  name: 'Models',
  data() {
    return { list: [], loading: false, page: 1, size: 10, total: 0 }
  },
  mounted() { this.fetch() },
  methods: {
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
    goDetail(id) { this.$router.push(`/models/${id}`) }
  }
}
</script>

<style>
.container { max-width: 960px; margin: 24px auto; }
.pager { margin-top: 16px; display: flex; justify-content: center; }
</style>
