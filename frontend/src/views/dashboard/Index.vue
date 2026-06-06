<template>
  <div class="dashboard">
    <!-- Skeleton loading -->
    <template v-if="loading">
      <div class="stats-grid">
        <div v-for="i in 4" :key="i" class="stat-card skeleton-stat">
          <div class="stat-icon-wrapper skeleton-shimmer" style="width:44px;height:44px;border-radius:var(--radius-md);"></div>
          <div class="stat-content">
            <div class="skeleton-shimmer" style="width:80px;height:12px;border-radius:4px;margin-bottom:6px;"></div>
            <div class="skeleton-shimmer" style="width:120px;height:24px;border-radius:4px;"></div>
          </div>
        </div>
      </div>
      <div class="section-grid section-gap">
        <div v-for="i in 2" :key="i" class="section-card">
          <div class="skeleton-shimmer" style="height:300px;"></div>
        </div>
      </div>
    </template>

    <template v-else>
      <!-- Stats row -->
      <div class="stats-grid">
        <div class="stat-card stat-card--success stagger-1">
          <div class="stat-card-bg"></div>
          <div class="stat-icon">
            <TrendingUp :size="22" :stroke-width="1.5" />
          </div>
          <div class="stat-content">
            <span class="stat-label">{{ t('dashboard.monthlyExpense') }}</span>
            <span class="stat-value">{{ formatAmount(stats?.monthly_expense || 0, 'CNY') }}</span>
          </div>
          <span class="stat-badge">{{ t('dashboard.monthlyExpense').split(' ')[0] }}</span>
        </div>
        <div class="stat-card stat-card--info stagger-2">
          <div class="stat-card-bg"></div>
          <div class="stat-icon">
            <BarChart3 :size="22" :stroke-width="1.5" />
          </div>
          <div class="stat-content">
            <span class="stat-label">{{ t('dashboard.yearlyExpense') }}</span>
            <span class="stat-value">{{ formatAmount(stats?.yearly_expense || 0, 'CNY') }}</span>
          </div>
          <span class="stat-badge">{{ t('dashboard.yearlyExpense').split(' ')[0] }}</span>
        </div>
        <div class="stat-card stat-card--accent stagger-3">
          <div class="stat-card-bg"></div>
          <div class="stat-icon">
            <CreditCard :size="22" :stroke-width="1.5" />
          </div>
          <div class="stat-content">
            <span class="stat-label">{{ t('dashboard.subscriptionCount') }}</span>
            <span class="stat-value">{{ stats?.subscription_count || 0 }}</span>
          </div>
        </div>
        <div class="stat-card stat-card--warning stagger-4">
          <div class="stat-card-bg"></div>
          <div class="stat-icon">
            <Package :size="22" :stroke-width="1.5" />
          </div>
          <div class="stat-content">
            <span class="stat-label">{{ t('dashboard.assetCount') }}</span>
            <span class="stat-value">{{ stats?.asset_count || 0 }}</span>
          </div>
        </div>
      </div>

      <!-- Middle row: Chart + Renewals -->
      <div class="section-grid section-gap">
        <div class="section-card stagger-3">
          <div class="section-card-header">
            <h3 class="section-card-title">{{ t('dashboard.categoryExpense') }}</h3>
          </div>
          <div class="chart-wrapper">
            <div ref="chartRef" class="chart-container"></div>
            <n-empty v-if="!stats?.category_expense?.length" :description="t('common.noData')" class="chart-empty" />
          </div>
        </div>
        <div class="section-card stagger-4">
          <div class="section-card-header">
            <h3 class="section-card-title">{{ t('dashboard.upcomingRenewals') }}</h3>
            <n-tag v-if="stats?.upcoming_renewal_list?.length" size="small" round class="renewal-count-tag">
              {{ stats.upcoming_renewal_list.length }}
            </n-tag>
          </div>
          <div class="renewal-list">
            <n-empty v-if="!stats?.upcoming_renewal_list?.length" :description="t('common.noData')" />
            <div v-else class="renewal-items">
              <div v-for="item in stats?.upcoming_renewal_list || []" :key="item.id" class="renewal-item">
                <div class="renewal-icon">
                  <CreditCard :size="16" :stroke-width="1.5" />
                </div>
                <div class="renewal-info">
                  <span class="renewal-name">{{ item.name }}</span>
                  <span class="renewal-date">{{ item.next_payment_date }}</span>
                </div>
                <span class="renewal-amount">{{ formatAmount(item.amount, item.currency) }}</span>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- Bottom row: Warning + Expiring -->
      <div class="section-grid section-gap">
        <div class="section-card stagger-5">
          <div class="section-card-header">
            <h3 class="section-card-title">{{ t('dashboard.warningAssets') }}</h3>
            <n-tag v-if="stats?.warning_assets" type="warning" size="small" round>{{ stats.warning_assets }}</n-tag>
          </div>
          <div class="renewal-list">
            <n-empty v-if="!warningAssetList.length" :description="t('common.noData')" />
            <div v-else class="renewal-items">
              <div v-for="item in warningAssetList" :key="item.id" class="renewal-item renewal-item--warning">
                <div class="renewal-icon renewal-icon--warning">
                  <AlertTriangle :size="16" :stroke-width="1.5" />
                </div>
                <div class="renewal-info">
                  <span class="renewal-name">{{ item.name }}</span>
                  <span class="renewal-date">{{ t('asset.expireDate') }}: {{ item.expire_date }}</span>
                </div>
                <n-tag type="warning" size="small">{{ t('asset.warning') }}</n-tag>
              </div>
            </div>
          </div>
        </div>
        <div class="section-card stagger-6">
          <div class="section-card-header">
            <h3 class="section-card-title">{{ t('dashboard.expiringAssets') }}</h3>
          </div>
          <div class="renewal-list">
            <n-empty v-if="!expiringSoonAssetList.length" :description="t('common.noData')" />
            <div v-else class="renewal-items">
              <div v-for="item in expiringSoonAssetList" :key="item.id" class="renewal-item">
                <div class="renewal-icon renewal-icon--info">
                  <Clock :size="16" :stroke-width="1.5" />
                </div>
                <div class="renewal-info">
                  <span class="renewal-name">{{ item.name }}</span>
                  <span class="renewal-date">{{ typeLabel(item.asset_type) }} · {{ t('asset.expireDate') }}: {{ item.expire_date }}</span>
                </div>
                <n-tag type="info" size="small">{{ typeLabel(item.asset_type) }}</n-tag>
              </div>
            </div>
          </div>
        </div>
      </div>
    </template>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted, nextTick, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { NTag, NEmpty } from 'naive-ui'
import { TrendingUp, BarChart3, CreditCard, Package, AlertTriangle, Clock } from '@lucide/vue'
import * as echarts from 'echarts/core'
import { PieChart } from 'echarts/charts'
import { TooltipComponent, LegendComponent } from 'echarts/components'
import { CanvasRenderer } from 'echarts/renderers'
import { type DashboardData } from '@/api/dashboard'

echarts.use([PieChart, TooltipComponent, LegendComponent, CanvasRenderer])
import { useDashboardStore } from '@/stores/dashboard'
import { formatAmount } from '@/utils/currency'
import { useAssetLabels } from '@/utils/mappings'
import { useTheme } from '@/composables/useTheme'

const { t } = useI18n()
const { isDark } = useTheme()
const dashboardStore = useDashboardStore()
const { typeLabel } = useAssetLabels()
const loading = computed(() => !dashboardStore.loaded)
const stats = computed<DashboardData | null>(() => dashboardStore.data)
const chartRef = ref<HTMLElement | null>(null)
let chart: echarts.ECharts | null = null

const warningAssetList = computed(() =>
  (stats.value?.expiring_asset_list || []).filter(item => item.status === 'warning')
)

const expiringSoonAssetList = computed(() =>
  (stats.value?.expiring_asset_list || []).filter(item =>
    item.status !== 'warning' && item.expire_date &&
    new Date(item.expire_date) <= new Date(Date.now() + 30 * 24 * 60 * 60 * 1000)
  )
)

function handleResize() { chart?.resize() }

watch(isDark, () => {
  if (chart) renderChart()
})

onMounted(async () => {
  await dashboardStore.ensureLoaded()
  await nextTick()
  renderChart()
})

onUnmounted(() => {
  window.removeEventListener('resize', handleResize)
  chart?.dispose()
})

function renderChart() {
  if (!chartRef.value || !stats.value?.category_expense?.length) return

  if (!chart) {
    chart = echarts.init(chartRef.value)
    window.addEventListener('resize', handleResize)
  }

  const data = stats.value.category_expense.map(item => ({
    name: item.category_name,
    value: Number(item.amount.toFixed(2)),
    itemStyle: { color: item.color },
  }))

  const dark = isDark.value

  chart.setOption({
    tooltip: {
      trigger: 'item',
      formatter: '{b}: {c} ({d}%)',
      backgroundColor: dark ? '#1B2336' : '#FFFFFF',
      borderColor: dark ? '#334155' : '#E2E8F0',
      textStyle: {
        color: dark ? '#F1F5F9' : '#1E293B',
        fontFamily: "'DM Sans', system-ui, sans-serif",
        fontSize: 13,
      },
      borderWidth: 1,
      padding: [8, 12],
      extraCssText: 'border-radius: 8px; box-shadow: 0 4px 12px rgba(0,0,0,0.1);',
    },
    series: [{
      type: 'pie',
      radius: ['50%', '78%'],
      center: ['50%', '50%'],
      avoidLabelOverlap: false,
      itemStyle: {
        borderRadius: 6,
        borderColor: dark ? '#141C2E' : '#FFFFFF',
        borderWidth: 3,
      },
      label: {
        show: true,
        color: dark ? '#94A3B8' : '#475569',
        fontFamily: "'DM Sans', system-ui, sans-serif",
        fontSize: 12,
        formatter: '{b}',
      },
      emphasis: {
        scaleSize: 8,
        label: {
          show: true,
          fontSize: 14,
          fontWeight: 600,
          fontFamily: "'Space Grotesk', system-ui, sans-serif",
        },
        itemStyle: {
          shadowBlur: 20,
          shadowColor: 'rgba(0, 0, 0, 0.15)',
        },
      },
      animationType: 'scale',
      animationEasing: 'elasticOut',
      animationDelay: (idx: number) => idx * 80,
      data,
    }],
  }, true)
}
</script>

<style scoped>
.dashboard {
  width: 100%;
  max-width: 1440px;
  margin: 0 auto;
}

/* ---- Stats Grid ---- */
.stats-grid {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: var(--spacing-md);
}

.stat-card {
  position: relative;
  padding: var(--spacing-lg);
  border-radius: var(--radius-xl);
  background: var(--color-bg-card);
  border: 1px solid var(--color-border);
  display: flex;
  align-items: flex-start;
  gap: var(--spacing-md);
  transition: box-shadow var(--transition-smooth), transform var(--transition-smooth), border-color var(--transition-smooth);
  overflow: hidden;
  cursor: default;
}

.stat-card:hover {
  box-shadow: var(--shadow-card-hover);
  transform: translateY(-2px);
  border-color: transparent;
}

.stat-card-bg {
  position: absolute;
  inset: 0;
  opacity: 0;
  transition: opacity var(--transition-smooth);
  pointer-events: none;
}

.stat-card:hover .stat-card-bg {
  opacity: 1;
}

.stat-card--success .stat-card-bg { background: linear-gradient(135deg, rgba(34, 197, 94, 0.04) 0%, transparent 60%); }
.stat-card--info .stat-card-bg { background: linear-gradient(135deg, rgba(59, 130, 246, 0.04) 0%, transparent 60%); }
.stat-card--accent .stat-card-bg { background: linear-gradient(135deg, rgba(34, 197, 94, 0.04) 0%, transparent 60%); }
.stat-card--warning .stat-card-bg { background: linear-gradient(135deg, rgba(245, 158, 11, 0.04) 0%, transparent 60%); }

.stat-icon {
  width: 44px;
  height: 44px;
  border-radius: var(--radius-md);
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
  position: relative;
}

.stat-card--success .stat-icon {
  background: var(--color-success-light);
  color: var(--color-success);
}

.stat-card--info .stat-icon {
  background: var(--color-info-light);
  color: var(--color-info);
}

.stat-card--accent .stat-icon {
  background: var(--color-accent-light);
  color: var(--color-accent);
}

.stat-card--warning .stat-icon {
  background: var(--color-warning-light);
  color: var(--color-warning);
}

.stat-content {
  display: flex;
  flex-direction: column;
  gap: 4px;
  min-width: 0;
}

.stat-label {
  font-size: 13px;
  color: var(--color-text-muted);
  font-weight: 500;
}

.stat-value {
  font-family: var(--font-heading);
  font-size: 24px;
  font-weight: 700;
  color: var(--color-text-primary);
  letter-spacing: -0.02em;
  line-height: 1.2;
}

.stat-badge {
  position: absolute;
  top: var(--spacing-sm);
  right: var(--spacing-sm);
  font-size: 10px;
  text-transform: uppercase;
  letter-spacing: 0.05em;
  color: var(--color-text-muted);
  opacity: 0.6;
}

/* ---- Section Grid ---- */
.section-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: var(--spacing-md);
}

.section-card {
  background: var(--color-bg-card);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-xl);
  overflow: hidden;
  transition: border-color var(--transition-smooth);
}

.section-card:hover {
  border-color: var(--color-border-strong);
}

.section-card-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: var(--spacing-md) var(--spacing-lg);
  border-bottom: 1px solid var(--color-border);
}

.section-card-title {
  font-family: var(--font-heading);
  font-size: 15px;
  font-weight: 600;
  color: var(--color-text-primary);
  margin: 0;
  letter-spacing: -0.01em;
}

.renewal-count-tag {
  background: var(--color-accent-light) !important;
  color: var(--color-accent) !important;
  border: none !important;
}

/* ---- Chart ---- */
.chart-wrapper {
  position: relative;
  padding: var(--spacing-md);
}

.chart-container {
  height: 300px;
}

.chart-empty {
  position: absolute;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
}

/* ---- Renewal List ---- */
.renewal-list {
  padding: var(--spacing-sm) var(--spacing-lg) var(--spacing-lg);
}

.renewal-items {
  display: flex;
  flex-direction: column;
}

.renewal-item {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: var(--spacing-sm) 0;
  border-bottom: 1px solid var(--color-border);
  gap: var(--spacing-sm);
  transition: background var(--transition-fast);
  border-radius: var(--radius-sm);
}

.renewal-item:last-child {
  border-bottom: none;
}

.renewal-icon {
  width: 32px;
  height: 32px;
  border-radius: var(--radius-sm);
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
  background: var(--color-accent-light);
  color: var(--color-accent);
}

.renewal-icon--warning {
  background: var(--color-warning-light);
  color: var(--color-warning);
}

.renewal-icon--info {
  background: var(--color-info-light);
  color: var(--color-info);
}

.renewal-item--warning {
  background: var(--color-warning-light);
  margin: 0 calc(-1 * var(--spacing-lg));
  padding: var(--spacing-sm) var(--spacing-lg);
  border-radius: var(--radius-sm);
  border-bottom: none;
}

.renewal-info {
  display: flex;
  flex-direction: column;
  gap: 2px;
  min-width: 0;
  flex: 1;
}

.renewal-name {
  font-weight: 500;
  font-size: 14px;
  color: var(--color-text-primary);
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.renewal-date {
  font-size: 12px;
  color: var(--color-text-muted);
  font-family: var(--font-mono);
  font-size: 11px;
  letter-spacing: -0.01em;
}

.renewal-amount {
  font-family: var(--font-heading);
  font-weight: 600;
  font-size: 14px;
  color: var(--color-text-primary);
  white-space: nowrap;
  flex-shrink: 0;
}

/* ---- Skeleton ---- */
.skeleton-stat {
  display: flex;
  gap: var(--spacing-md);
  align-items: flex-start;
  padding: var(--spacing-lg);
  border-radius: var(--radius-xl);
  background: var(--color-bg-card);
  border: 1px solid var(--color-border);
}

/* ---- Responsive ---- */
@media (max-width: 1024px) {
  .stats-grid {
    grid-template-columns: repeat(2, 1fr);
  }
}

@media (max-width: 768px) {
  .stats-grid {
    grid-template-columns: repeat(2, 1fr);
    gap: var(--spacing-sm);
  }
  .stat-card {
    padding: var(--spacing-md);
    flex-direction: column;
    gap: var(--spacing-sm);
  }
  .stat-icon {
    width: 36px;
    height: 36px;
  }
  .stat-value {
    font-size: 20px;
  }
  .section-grid {
    grid-template-columns: 1fr;
  }
  .chart-container {
    height: 240px;
  }
  .renewal-item {
    padding: var(--spacing-xs) 0;
  }
}
</style>
