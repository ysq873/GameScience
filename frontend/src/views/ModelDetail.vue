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
            <el-button v-if="!purchasedByMe" type="primary" :disabled="model.status !== 1" @click="buy">购买</el-button>
            <el-button v-else type="success" @click="$router.push('/library')">去已购库</el-button>
            <el-button :type="isFav ? 'warning' : 'default'" :loading="favLoading" :disabled="favLoading" @click="toggleFav" style="margin-left:8px">{{ isFav ? '取消收藏' : '收藏' }}</el-button>
          </div>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script>
import { getModel } from '@/api/models'
import { createOrder } from '@/api/orders'
import { addFavorite, removeFavorite, listFavorites } from '@/api/favorites'

export default {
  name: 'ModelDetail',
  data() {
    return { model: {}, loading: false, purchasedByMe: false, isFav: false, favLoading: false }
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
      try {
        const lp = await import('@/api/downloads').then(m => m.listPurchases())
        const list = lp.data.list || []
        this.purchasedByMe = list.some(it => Number(it.model_id) === Number(id))
      } catch {}
      try {
        const favs = await listFavorites()
        const lst = favs.data.list || []
        this.isFav = lst.some(it => Number(it.id) === Number(id))
      } catch {}
    } finally { this.loading = false }
  },
  methods: {
    async buy() {
      if (this.purchasedByMe) { this.$message.warning('已购买过该模型'); this.$router.push('/library'); return }
      try {
        const res = await createOrder([Number(this.$route.params.id)])
        const orderId = res.data.order_id
        this.$router.push(`/pay/${orderId}`)
      } catch (e) {
        const code = e.response?.data?.code
        if (code === 'already_purchased') { this.$message.warning('已购买过该模型'); this.$router.push('/library'); return }
        this.$message.error(e.response?.data?.message || '下单失败')
      }
    },
    formatYuan(c) { return (Number(c) / 100).toFixed(2) },
    isAbsolute(u) { return /^https?:\/\//i.test(u) },
    coverSrc(c) {
      if (!c) return ''
      const norm = String(c).replace(/\\/g, '/')
      return `/api/static?file=${encodeURIComponent(norm)}`
    },
    statusText(s) { if (s === 1) return (this.purchasedByMe ? '已购买' : '上架'); if (s === 2) return '下架'; return '待定' },
    statusType(s) { if (s === 1) return 'success'; if (s === 2) return 'info'; return 'warning' },
    async toggleFav() {
      const id = Number(this.$route.params.id)
      this.favLoading = true
      try {
        if (this.isFav) {
          await removeFavorite(id)
          this.isFav = false
          this.$message.success('已取消收藏')
        } else {
          await addFavorite(id)
          this.isFav = true
          this.$message.success('已收藏')
        }
      } catch (e) {
        this.$message.error(e.response?.data?.message || '操作失败')
      } finally {
        this.favLoading = false
      }
    }
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
