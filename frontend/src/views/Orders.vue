<template>
  <div class="container">
    <h2>我的订单</h2>
    <el-table :data="list" v-loading="loading">
      <el-table-column prop="id" label="订单ID" width="100" />
      <el-table-column prop="total_cents" label="总金额" />
      <el-table-column prop="status" label="状态" />
      <el-table-column label="操作" width="260">
        <template #default="scope">
          <el-button v-if="scope.row.status==='pending'" type="primary" @click="goPay(scope.row.id)">支付</el-button>
          <el-button v-if="scope.row.status==='paid'" type="danger" @click="refund(scope.row.id)">退款</el-button>
        </template>
      </el-table-column>
    </el-table>
  </div>
</template>

<script>
import { listOrders, refundOrder } from '@/api/orders'

export default {
  name: 'Orders',
  data() {
    return { list: [], loading: false }
  },
  mounted() { this.fetch() },
  methods: {
    async fetch() {
      this.loading = true
      try {
        const res = await listOrders()
        this.list = res.data.list || []
      } finally { this.loading = false }
    },
    goPay(id) { this.$router.push(`/pay/${id}`) },
    async refund(id) { await refundOrder(id); this.fetch() }
  }
}
</script>

<style>
.container { max-width: 960px; margin: 24px auto; }
</style>
