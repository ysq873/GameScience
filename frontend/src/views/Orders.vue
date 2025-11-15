<template>
  <div class="container">
    <h2>我的订单</h2>
    <el-table :data="list" v-loading="loading">
      <el-table-column prop="id" label="订单ID" width="100" />
      <el-table-column label="总金额（元）">
        <template #default="scope">{{ formatYuan(scope.row.total_cents) }}</template>
      </el-table-column>
      <el-table-column label="状态">
        <template #default="scope">{{ statusText(scope.row) }}</template>
      </el-table-column>
      <el-table-column label="操作" width="260">
        <template #default="scope">
          <el-button v-if="isPendingOrExpired(scope.row)" :disabled="isExpired(scope.row)" type="primary" @click="goPay(scope.row.id)">支付</el-button>
          <el-button v-if="isPaid(scope.row)" type="danger" @click="refund(scope.row.id)">退款</el-button>
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
    async refund(id) { await refundOrder(id); this.fetch() },
    statusCode(row) { return Number.isFinite(row.status_code) ? row.status_code : (typeof row.status === 'number' ? row.status : undefined) },
    statusText(row) {
      const code = this.statusCode(row)
      const codeMap = {0:'待支付',1:'已支付',2:'已过期',3:'已退款'}
      if (Number.isFinite(code)) return codeMap[code] || String(code)
      const s = String(row.status || '').toLowerCase()
      const strMap = {pending:'待支付',paid:'已支付',expired:'已过期',refunded:'已退款'}
      return strMap[s] || (row.status === undefined ? '' : String(row.status))
    },
    isExpired(row) { const c = this.statusCode(row); return c === 2 || row.status === 'expired' },
    isPaid(row) { const c = this.statusCode(row); return c === 1 || row.status === 'paid' },
    isPendingOrExpired(row) { const c = this.statusCode(row); if (Number.isFinite(c)) return c === 0 || c === 2; return row.status === 'pending' || row.status === 'expired' },
    formatYuan(c) { return (Number(c) / 100).toFixed(2) }
  }
}
</script>

<style>
.container { max-width: 960px; margin: 24px auto; }
</style>
