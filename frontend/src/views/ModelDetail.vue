<template>
  <div class="container" v-loading="loading">
    <h2>{{ model.title }}</h2>
    <p>{{ model.description }}</p>
    <p>价格：{{ model.price_cents }}</p>
    <div class="actions">
      <el-button type="primary" @click="buy">购买</el-button>
    </div>
  </div>
</template>

<script>
import { getModel } from '@/api/models'
import { createOrder } from '@/api/orders'

export default {
  name: 'ModelDetail',
  data() {
    return { model: {}, loading: false }
  },
  async mounted() {
    this.loading = true
    try {
      const id = this.$route.params.id
      const res = await getModel(id)
      this.model = res.data
    } finally { this.loading = false }
  },
  methods: {
    async buy() {
      const res = await createOrder([Number(this.$route.params.id)])
      const orderId = res.data.order_id
      this.$router.push(`/pay/${orderId}`)
    }
  }
}
</script>

<style>
.container { max-width: 720px; margin: 24px auto; }
.actions { margin-top: 12px; }
</style>
