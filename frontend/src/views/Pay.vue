<template>
  <div class="pay-container" v-loading="loading">
    <el-card class="pay-card">
      <div class="pay-header">
        <div class="order-id"><span>订单号</span><strong>#{{ orderId }}</strong></div>
        <div class="pay-tools">
          <span>支付剩余时间：</span>
          <strong :class="{expired}">{{ countdown }}</strong>
        </div>
      </div>
      <el-table :data="items" size="small" style="width:100%" class="items-table">
        <el-table-column label="商品" min-width="360">
          <template #default="scope">
            <div class="item-row">
              <img v-if="scope.row.cover_url" :src="isAbsolute(scope.row.cover_url) ? scope.row.cover_url : coverSrc(scope.row.cover_url)" class="thumb" />
              <div class="title">{{ scope.row.title }}</div>
            </div>
          </template>
        </el-table-column>
        <el-table-column label="单价（元）" width="120">
          <template #default="scope">{{ formatYuan(scope.row.price_cents) }}</template>
        </el-table-column>
        <el-table-column label="数量" width="90">
          <template #default="scope">{{ scope.row.quantity || 1 }}</template>
        </el-table-column>
        <el-table-column label="小计（元）" width="120">
          <template #default="scope">{{ formatYuan(scope.row.subtotal_cents || scope.row.price_cents) }}</template>
        </el-table-column>
      </el-table>
      <div class="total">合计：<strong>￥{{ formatYuan(detail.total_cents) }}</strong></div>
      <div class="actions actions-bottom">
        <el-button type="primary" :disabled="expired || !!payInfo" @click="pay">发起支付</el-button>
        <el-button type="success" :disabled="!payInfo" @click="callback('succeeded')">模拟支付成功</el-button>
        <el-button type="warning" :disabled="!payInfo" @click="callback('failed')">模拟支付失败</el-button>
      </div>
      <div v-if="expired" class="expired-tip">订单已过期，请返回重新下单</div>
    </el-card>
  </div>
</template>

<script>
import { payOrder, mockCallback, getOrderDetail } from '@/api/orders'

export default {
  name: 'Pay',
  data() {
    return { orderId: this.$route.params.orderId, loading: false, payInfo: null, detail: {}, items: [], countdown: '', expired: false, _timer: null }
  },
  async mounted() {
    this.loading = true
    try {
      const { data } = await getOrderDetail(this.orderId)
      this.detail = data || {}
      this.items = data.items || []
      let left = Number.isFinite(Number(data.expires_seconds_left))
        ? Number(data.expires_seconds_left)
        : Math.max(0, Math.floor((new Date(data.expires_at || data.created_at).getTime() - new Date(data.server_now || Date.now()).getTime()) / 1000))
      {
        const mm = String(Math.floor(left / 60)).padStart(2, '0')
        const ss = String(left % 60).padStart(2, '0')
        this.countdown = `${mm}:${ss}`
        this.expired = left <= 0
      }
      this._timer = setInterval(() => {
        left = Math.max(0, left - 1)
        const mm = String(Math.floor(left / 60)).padStart(2, '0')
        const ss = String(left % 60).padStart(2, '0')
        this.countdown = `${mm}:${ss}`
        this.expired = left <= 0
        if (this.expired && this._timer) { clearInterval(this._timer); this._timer = null }
      }, 1000)
    } finally { this.loading = false }
  },
  beforeUnmount() { if (this._timer) clearInterval(this._timer) },
  methods: {
    async pay() {
      this.loading = true
      try {
        const res = await payOrder(this.orderId)
        this.payInfo = res.data
      } catch (e) {
        this.$message.error(e.response?.data || '发起支付失败')
      } finally { this.loading = false }
    },
    async callback(status) {
      await mockCallback(Number(this.orderId), status, this.payInfo.idempotency_key)
      this.$router.push('/orders')
    },
    formatYuan(c) { return (Number(c) / 100).toFixed(2) },
    isAbsolute(u) { return /^https?:\/\//i.test(u) },
    coverSrc(c) { if (!c) return ''; const norm = String(c).replace(/\\/g, '/'); return `/api/static?file=${encodeURIComponent(norm)}` }
  }
}
</script>

<style>
.pay-container { max-width: 960px; margin: 24px auto; }
.pay-card { padding-bottom: 8px }
.pay-header { display:flex; justify-content: space-between; align-items:center; margin-bottom: 12px }
.order-id span { margin-right: 6px }
.pay-tools { display:flex; align-items:center; gap: 12px }
.expired { color: #f56c6c }
.actions { display:flex; gap: 12px }
.actions-bottom { justify-content: flex-end; margin-top: 8px }
.items-table .item-row { display:flex; align-items:center; gap: 12px }
.items-table .thumb { width: 64px; height: 64px; object-fit: cover; border-radius: 4px }
.items-table .title { font-size: 14px }
.total { text-align: right; margin-top: 8px }
.expired-tip { margin-top: 8px; color: #f56c6c }
</style>
