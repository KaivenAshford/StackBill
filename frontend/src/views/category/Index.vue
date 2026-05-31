<template>
  <div>
    <div class="page-toolbar">
      <h2 class="page-heading">{{ t('nav.categories') }}</h2>
      <n-button type="primary" @click="openCreate">
        <template #icon>
          <Plus :size="16" :stroke-width="1.5" />
        </template>
        {{ t('common.create') }}
      </n-button>
    </div>

    <!-- Single stable grid — skeleton or real cards share the same container -->
    <div class="category-grid" v-if="displayItems.length > 0">
      <div v-for="(cat, idx) in displayItems" :key="cat._key" class="category-card" :class="{ 'category-card--skeleton': cat._skeleton }">
        <div class="category-color" :style="{ background: cat.color }"></div>
        <div class="category-info">
          <span class="category-name">{{ cat.name }}</span>
          <n-tag v-if="!cat._skeleton" size="small" round>{{ cat.type === 'subscription' ? t('category.subscription') : t('category.asset') }}</n-tag>
        </div>
        <div v-if="!cat._skeleton" class="category-actions">
          <n-button quaternary circle size="small" :aria-label="t('common.edit')" @click="startEdit(cat as Category)">
            <template #icon>
              <Pencil :size="14" :stroke-width="1.5" />
            </template>
          </n-button>
          <n-button quaternary circle size="small" type="error" :aria-label="t('common.delete')" @click="confirmDeleteCategory((cat as Category).id)">
            <template #icon>
              <Trash2 :size="14" :stroke-width="1.5" />
            </template>
          </n-button>
        </div>
      </div>
    </div>
    <n-empty v-if="categoryStore.loaded && categoryStore.categories.length === 0" :description="t('common.noData')" class="empty-state-gap" />

    <n-modal v-model:show="showCreate" :title="editing ? t('common.edit') : t('common.create')" preset="card" style="max-width:480px; width: calc(100vw - 32px);">
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
        <div class="modal-actions">
          <n-button @click="showCreate = false">{{ t('common.cancel') }}</n-button>
          <n-button type="primary" :loading="saving" @click="handleSave">{{ t('common.save') }}</n-button>
        </div>
      </template>
    </n-modal>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { NButton, NModal, NForm, NFormItem, NInput, NSelect, NColorPicker, NTag, NEmpty, useMessage, useDialog } from 'naive-ui'
import { Plus, Pencil, Trash2 } from '@lucide/vue'
import { createCategory, updateCategory, deleteCategory } from '@/api/category'
import type { Category } from '@/types'
import { useCategoryStore } from '@/stores/category'

interface SkeletonItem { _key: string; _skeleton: true; name: string; color: string }
interface RealItem extends Category { _key: string; _skeleton: false }
type DisplayItem = SkeletonItem | RealItem

const SKELETON_COUNT = 6
const SKELETON_COLORS = ['#3B82F6', '#8B5CF6', '#22C55E', '#F59E0B', '#EF4444', '#14B8A6']

function makeSkeletons(): SkeletonItem[] {
  return Array.from({ length: SKELETON_COUNT }, (_, i) => ({
    _key: `skel-${i}`,
    _skeleton: true as const,
    name: '',
    color: SKELETON_COLORS[i % SKELETON_COLORS.length],
  }))
}

const { t } = useI18n()
const message = useMessage()
const dialog = useDialog()
const categoryStore = useCategoryStore()

const showCreate = ref(false)
const editing = ref<Category | null>(null)
const saving = ref(false)
const form = reactive({ name: '', type: 'subscription', color: '#1890ff', icon: '' })

const displayItems = computed<DisplayItem[]>(() => {
  if (!categoryStore.loaded) return makeSkeletons()
  return categoryStore.categories.map(c => ({ ...c, _key: `cat-${c.id}`, _skeleton: false as const }))
})

const typeOptions = [
  { label: t('category.subscription'), value: 'subscription' },
  { label: t('category.asset'), value: 'asset' },
]

onMounted(() => categoryStore.ensureLoaded())

function openCreate() {
  editing.value = null
  form.name = ''
  form.type = 'subscription'
  form.color = '#1890ff'
  form.icon = ''
  showCreate.value = true
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
    await categoryStore.refresh()
  } catch (e: unknown) {
    message.error((e as Error).message || t('common.failed'))
  } finally {
    saving.value = false
  }
}

function confirmDeleteCategory(id: number) {
  dialog.warning({
    title: t('common.confirm'),
    content: t('common.confirmDelete'),
    positiveText: t('common.confirm'),
    negativeText: t('common.cancel'),
    onPositiveClick: () => handleDelete(id),
  })
}

async function handleDelete(id: number) {
  try {
    await deleteCategory(id)
    message.success(t('common.success'))
    await categoryStore.refresh()
  } catch (e: unknown) {
    message.error((e as Error).message || t('common.failed'))
  }
}
</script>

<style scoped>
.category-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(280px, 1fr));
  gap: var(--spacing-md);
}

.category-card {
  display: flex;
  align-items: center;
  gap: var(--spacing-md);
  padding: var(--spacing-md) var(--spacing-lg);
  background: var(--color-bg-card);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-xl);
  transition: box-shadow var(--transition-smooth), transform var(--transition-smooth), border-color var(--transition-smooth);
  position: relative;
  overflow: hidden;
  min-height: 72px;
}

.category-card:hover:not(.category-card--skeleton) {
  box-shadow: var(--shadow-card-hover);
  transform: translateY(-2px);
  border-color: transparent;
}

.category-card--skeleton {
  pointer-events: none;
  opacity: 0.7;
}

.category-card--skeleton .category-color {
  opacity: 0.3;
}

.category-card--skeleton .category-name {
  color: transparent;
  background: var(--color-bg-muted);
  border-radius: 4px;
  width: 80px;
  height: 14px;
}

.category-card--skeleton .category-info::after {
  content: '';
  display: block;
  width: 48px;
  height: 20px;
  border-radius: 10px;
  background: var(--color-bg-muted);
}

.category-card:not(.category-card--skeleton)::after {
  content: '';
  position: absolute;
  inset: 0;
  opacity: 0;
  transition: opacity var(--transition-smooth);
  pointer-events: none;
}

.category-card:not(.category-card--skeleton):hover::after {
  opacity: 1;
  background: linear-gradient(135deg, rgba(34, 197, 94, 0.02) 0%, transparent 60%);
}

.category-color {
  width: 6px;
  height: 40px;
  border-radius: var(--radius-full);
  flex-shrink: 0;
  transition: height var(--transition-smooth), opacity var(--transition-base);
}

.category-card:not(.category-card--skeleton):hover .category-color {
  height: 44px;
}

.category-info {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 4px;
  min-width: 0;
  position: relative;
  z-index: 1;
}

.category-name {
  font-weight: 500;
  font-size: 14px;
  color: var(--color-text-primary);
}

.category-actions {
  display: flex;
  gap: var(--spacing-xs);
  flex-shrink: 0;
  position: relative;
  z-index: 1;
}

.modal-actions {
  display: flex;
  justify-content: flex-end;
  gap: var(--spacing-sm);
}
</style>
