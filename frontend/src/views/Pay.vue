<template>
  <div class="container" v-loading="loading">
    <h2>支付订单 {{ orderId }}</h2>
    <div v-if="payInfo">
      <p>支付单号：{{ payInfo.payment_id }}</p>
      <p>幂等键：{{ payInfo.idempotency_key }}</p>
      <div class="actions">
        <el-button type="primary" @click="callback('succeeded')">模拟支付成功</el-button>
        <el-button type="warning" @click="callback('failed')">模拟支付失败</el-button>
      </div>
    </div>
    <div v-else>
      <el-button type="primary" @click="pay">发起支付</el-button>
    </div>
  </div>
</template>

<script>
import { payOrder, mockCallback } from '@/api/orders'

export default {
  name: 'Pay',
  data() {
    return { orderId: this.$route.params.orderId, loading: false, payInfo: null }
  },
  methods: {
    async pay() {
      this.loading = true
      try {
        const res = await payOrder(this.orderId)
        this.payInfo = res.data
      } finally { this.loading = false }
    },
    async callback(status) {
      await mockCallback(Number(this.orderId), status, this.payInfo.idempotency_key)
      this.$router.push('/orders')
    }
  }
}
</script>

<style>
.container { max-width: 720px; margin: 24px auto; }
.actions { margin-top: 12px; display: flex; gap: 12px; }
</style>
