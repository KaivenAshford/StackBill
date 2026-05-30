<template>
  <div>
    <div style="display:flex;justify-content:space-between;margin-bottom:16px;">
      <h2 style="margin:0">{{ t('nav.categories') }}</h2>
      <n-button type="primary" @click="showCreate = true">{{ t('common.create') }}</n-button>
    </div>
    <n-spin :show="loading">
      <n-data-table :columns="columns" :data="categories" :bordered="false" />
    </n-spin>
    <n-modal v-model:show="showCreate" :title="editing ? t('common.edit') : t('common.create')" preset="card" style="width:480px;">
      <n-form :model="form" label-placement="left" label-width="80">
        <n-form-item :label="t('category.name')">
          <n-input v-model:value="form.name" />
        </n-form-item>
        <n-form-item :label="t('category.type')">
          <n-select v-model:value="form.type" :options="typeOptions" />
        </n-form-item>
        <n-form-item :label="t('category.color')">
          <n-color-picker v-model:value="form.color" :modes="['hex']" :show-alpha="false" />
        </n-form-item>
        <n-form-item :label="t('category.icon')">
          <n-input v-model:value="form.icon" />
        </n-form-item>
      </n-form>
      <template #action>
        <n-button @click="showCreate = false">{{ t('common.cancel') }}</n-button>
        <n-button type="primary" :loading="saving" @click="handleSave" style="margin-left:8px;">{{ t('common.save') }}</n-button>
      </template>
    </n-modal>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, h } from 'vue'
import { useI18n } from 'vue-i18n'
import { NDataTable, NButton, NModal, NForm, NFormItem, NInput, NSelect, NColorPicker, NTag, NSpace, NSpin, useMessage } from 'naive-ui'
import { listCategories, createCategory, updateCategory, deleteCategory } from '@/api/category'
import type { Category } from '@/types'

const { t } = useI18n()
const message = useMessage()

const categories = ref<Category[]>([])
const loading = ref(true)
const showCreate = ref(false)
const editing = ref<Category | null>(null)
const saving = ref(false)
const form = reactive({ name: '', type: 'subscription', color: '#1890ff', icon: '' })

const typeOptions = [
  { label: t('category.subscription'), value: 'subscription' },
  { label: t('category.asset'), value: 'asset' },
]

const columns = [
  { title: t('category.name'), key: 'name' },
  { title: t('category.type'), key: 'type', render: (row: Category) => h(NTag, { size: 'small' }, { default: () => row.type === 'subscription' ? t('category.subscription') : t('category.asset') }) },
  { title: t('category.color'), key: 'color', render: (row: Category) => h('div', { style: { width: '20px', height: '20px', borderRadius: '4px', background: row.color } }) },
  { title: t('common.edit'), key: 'actions', render: (row: Category) => h(NSpace, null, {
    default: () => [
      h(NButton, { size: 'small', onClick: () => startEdit(row) }, { default: () => t('common.edit') }),
      h(NButton, { size: 'small', type: 'error', onClick: () => handleDelete(row.id) }, { default: () => t('common.delete') }),
    ],
  })},
]

onMounted(() => fetchData())

async function fetchData() {
  loading.value = true
  try {
    const res = await listCategories()
    categories.value = res.data
  } finally {
    loading.value = false
  }
}

function startEdit(cat: Category) {
  editing.value = cat
  form.name = cat.name
  form.type = cat.type
  form.color = cat.color
  form.icon = cat.icon
  showCreate.value = true
}

async function handleSave() {
  saving.value = true
  try {
    if (editing.value) {
      await updateCategory(editing.value.id, { name: form.name, type: form.type, color: form.color, icon: form.icon })
    } else {
      await createCategory({ name: form.name, type: form.type, color: form.color, icon: form.icon } as Partial<Category>)
    }
    message.success(t('common.success'))
    showCreate.value = false
    editing.value = null
    form.name = ''
    form.type = 'subscription'
    form.color = '#1890ff'
    form.icon = ''
    await fetchData()
  } catch (e: unknown) {
    message.error((e as Error).message || t('common.failed'))
  } finally {
    saving.value = false
  }
}

async function handleDelete(id: number) {
  try {
    await deleteCategory(id)
    message.success(t('common.success'))
    await fetchData()
  } catch (e: unknown) {
    message.error((e as Error).message || t('common.failed'))
  }
}
</script>
