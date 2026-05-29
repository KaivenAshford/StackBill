<template>
  <div class="dashboard">
    <n-grid :cols="4" :x-gap="16" :y-gap="16" responsive="screen" item-responsive>
      <n-gi span="4 m:2 l:1">
        <n-card>
          <n-statistic :label="t('dashboard.monthlyExpense')">
            <template #prefix>$</template>
            {{ stats?.monthly_expense?.toFixed(2) || '0.00' }}
          </n-statistic>
        </n-card>
      </n-gi>
      <n-gi span="4 m:2 l:1">
        <n-card>
          <n-statistic :label="t('dashboard.yearlyExpense')">
            <template #prefix>$</template>
            {{ stats?.yearly_expense?.toFixed(2) || '0.00' }}
          </n-statistic>
        </n-card>
      </n-gi>
      <n-gi span="4 m:2 l:1">
        <n-card>
          <n-statistic :label="t('dashboard.subscriptionCount')" :value="stats?.subscription_count || 0" />
        </n-card>
      </n-gi>
      <n-gi span="4 m:2 l:1">
        <n-card>
          <n-statistic :label="t('dashboard.assetCount')" :value="stats?.asset_count || 0" />
        </n-card>
      </n-gi>
    </n-grid>

    <n-grid :cols="2" :x-gap="16" :y-gap="16" style="margin-top:16px;" responsive="screen" item-responsive>
      <n-gi span="2 l:1">
        <n-card title="Category Expense">
          <div ref="chartRef" style="height:300px;"></div>
        </n-card>
      </n-gi>
      <n-gi span="2 l:1">
        <n-card :title="t('dashboard.upcomingRenewals')">
          <n-empty v-if="!stats?.upcoming_renewal_list?.length" :description="t('common.noData')" />
          <n-list v-else>
            <n-list-item v-for="item in stats?.upcoming_renewal_list || []" :key="item.id">
              <n-thing :title="item.name">
                <template #description>
                  <span style="color:#999;">{{ item.next_payment_date }} · {{ item.amount }} {{ item.currency }}</span>
                </template>
              </n-thing>
            </n-list-item>
          </n-list>
        </n-card>
      </n-gi>
    </n-grid>

    <n-grid :cols="2" :x-gap="16" :y-gap="16" style="margin-top:16px;" responsive="screen" item-responsive>
      <n-gi span="2 l:1">
        <n-card title="Warning Assets">
          <n-tag type="warning" size="large">{{ stats?.warning_assets || 0 }} warnings</n-tag>
          <n-empty v-if="!stats?.expiring_asset_list?.length" style="margin-top:16px;" :description="t('common.noData')" />
          <n-list v-else style="margin-top:16px;">
            <n-list-item v-for="item in stats?.expiring_asset_list || []" :key="item.id">
              <n-thing :title="item.name">
                <template #description>
                  <span style="color:#999;">{{ t('asset.expireDate') }}: {{ item.expire_date }} · {{ item.status }}</span>
                </template>
              </n-thing>
            </n-list-item>
          </n-list>
        </n-card>
      </n-gi>
      <n-gi span="2 l:1">
        <n-card :title="t('dashboard.expiringAssets')">
          <n-empty v-if="!stats?.expiring_asset_list?.length" :description="t('common.noData')" />
          <n-list v-else>
            <n-list-item v-for="item in stats?.expiring_asset_list || []" :key="item.id">
              <n-thing :title="item.name">
                <template #description>
                  <span style="color:#999;">{{ item.asset_type }} · {{ t('asset.expireDate') }}: {{ item.expire_date }}</span>
                </template>
              </n-thing>
            </n-list-item>
          </n-list>
        </n-card>
      </n-gi>
    </n-grid>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, nextTick } from 'vue'
import { useI18n } from 'vue-i18n'
import { NGrid, NGi, NCard, NStatistic, NList, NListItem, NThing, NEmpty, NTag } from 'naive-ui'
import * as echarts from 'echarts'
import { getDashboard, type DashboardData } from '@/api/dashboard'

const { t } = useI18n()
const stats = ref<DashboardData | null>(null)
const chartRef = ref<HTMLElement | null>(null)
let chart: echarts.ECharts | null = null

onMounted(async () => {
  const res = await getDashboard()
  stats.value = res.data
  await nextTick()
  renderChart()
})

function renderChart() {
  if (!chartRef.value || !stats.value?.category_expense?.length) return

  chart = echarts.init(chartRef.value)

  const data = stats.value.category_expense.map(item => ({
    name: item.category_name,
    value: Number(item.amount.toFixed(2)),
    itemStyle: { color: item.color },
  }))

  chart.setOption({
    tooltip: { trigger: 'item', formatter: '{b}: ${c} ({d}%)' },
    series: [{
      type: 'pie',
      radius: ['40%', '70%'],
      avoidLabelOverlap: false,
      itemStyle: { borderRadius: 6 },
      label: { show: true, formatter: '{b}\n${c}' },
      data,
    }],
  })

  window.addEventListener('resize', () => chart?.resize())
}
</script>
