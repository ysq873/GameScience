<template>
  <div class="detail-container" v-loading="loading">
    <div class="detail-header">
      <el-button link @click="$router.push('/models')">返回作品</el-button>
    </div>
    <el-row :gutter="20">
      <el-col :span="12">
        <el-card>
          <el-skeleton :loading="loading" animated>
            <template #template>
              <el-skeleton-item variant="image" style="width:100%;height:320px" />
            </template>
            <template #default>
              <el-image v-if="coverUrl" :src="coverUrl" fit="cover" style="width:100%;height:320px;border-radius:6px" />
              <div v-else class="no-cover">无封面</div>
            </template>
          </el-skeleton>
        </el-card>
      </el-col>
      <el-col :span="12">
        <el-card>
          <div class="title">{{ model.title }}</div>
          <div class="desc">{{ model.description }}</div>
          <div class="meta">
            <span class="price">￥{{ formatYuan(model.price_cents) }}</span>
            <el-tag size="small" :type="statusType(model.status)">{{ statusText(model.status) }}</el-tag>
          </div>
          <div class="actions">
            <el-button type="primary" :disabled="model.status !== 1" @click="buy">购买</el-button>
          </div>
        </el-card>
      </el-col>
    </el-row>
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
  computed: {
    coverUrl() {
      const c = this.model.cover_url
      if (!c) return ''
      if (/^https?:\/\//i.test(c)) return c
      const norm = String(c).replace(/\\/g, '/')
      return `/api/static?file=${encodeURIComponent(norm)}`
    }
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
    },
    formatYuan(c) { return (Number(c) / 100).toFixed(2) },
    isAbsolute(u) { return /^https?:\/\//i.test(u) },
    coverSrc(c) {
      if (!c) return ''
      const norm = String(c).replace(/\\/g, '/')
      return `/api/static?file=${encodeURIComponent(norm)}`
    },
    statusText(s) { if (s === 1) return '上架'; if (s === 2) return '下架'; return '待定' },
    statusType(s) { if (s === 1) return 'success'; if (s === 2) return 'info'; return 'warning' }
  }
}
</script>

<style>
.detail-container { max-width: 960px; margin: 24px auto; }
.detail-header { display:flex; justify-content:flex-start; margin-bottom: 12px; }
.title { font-size: 22px; font-weight: 600; margin-bottom: 8px; }
.desc { color: #666; margin-bottom: 12px; }
.meta { display:flex; align-items:center; gap: 12px; margin-bottom: 16px; }
.price { font-size: 18px; color: #333; }
.actions { margin-top: 8px; }
.no-cover { width:100%; height:320px; display:flex; align-items:center; justify-content:center; color:#999; background:#f5f5f5; border-radius:6px }
</style>
