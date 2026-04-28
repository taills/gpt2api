<script setup lang="ts">
import { onMounted, ref, reactive } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { formatDateTime, nullVal } from '@/utils/format'
import {
  listMyKeys,
  createMyKey,
  updateMyKey,
  deleteMyKey,
  type ApiKey,
  type CreateKeyInput,
  type UpdateKeyInput,
} from '@/api/me'

// ========== 列表 ==========
const loading = ref(false)
const rows = ref<ApiKey[]>([])
const total = ref(0)
const pager = reactive({ page: 1, page_size: 10 })

async function fetchList() {
  loading.value = true
  try {
    const data = await listMyKeys({ page: pager.page, page_size: pager.page_size })
    rows.value = data.list || []
    total.value = data.total || 0
  } catch (e: any) {
    ElMessage.error(e?.message || '加载失败')
  } finally {
    loading.value = false
  }
}

function onPageChange(p: number) {
  pager.page = p
  fetchList()
}

// ========== 创建 ==========
const createVisible = ref(false)
const createLoading = ref(false)
const createForm = reactive<CreateKeyInput>({ name: '', quota_limit: 0, rpm: 0, tpm: 0 })

// 新建成功后展示完整密钥(仅此一次)
const newKeyVisible = ref(false)
const newKeyValue = ref('')

function openCreate() {
  createForm.name = ''
  createForm.quota_limit = 0
  createForm.rpm = 0
  createForm.tpm = 0
  createVisible.value = true
}

async function submitCreate() {
  if (!createForm.name.trim()) {
    ElMessage.warning('请输入 Key 名称')
    return
  }
  createLoading.value = true
  try {
    const res = await createMyKey({
      name: createForm.name.trim(),
      quota_limit: createForm.quota_limit || 0,
      rpm: createForm.rpm || 0,
      tpm: createForm.tpm || 0,
    })
    createVisible.value = false
    newKeyValue.value = res.key
    newKeyVisible.value = true
    await fetchList()
  } catch (e: any) {
    ElMessage.error(e?.message || '创建失败')
  } finally {
    createLoading.value = false
  }
}

async function copyKey() {
  try {
    await navigator.clipboard.writeText(newKeyValue.value)
    ElMessage.success('已复制到剪贴板')
  } catch {
    ElMessage.warning('请手动复制')
  }
}

// ========== 编辑 ==========
const editVisible = ref(false)
const editLoading = ref(false)
const editTarget = ref<ApiKey | null>(null)
const editForm = reactive<UpdateKeyInput>({ name: '', quota_limit: 0, rpm: 0, tpm: 0 })

function openEdit(row: ApiKey) {
  editTarget.value = row
  editForm.name = row.name
  editForm.quota_limit = row.quota_limit
  editForm.rpm = row.rpm
  editForm.tpm = row.tpm
  editVisible.value = true
}

async function submitEdit() {
  if (!editTarget.value) return
  if (!editForm.name?.trim()) {
    ElMessage.warning('请输入 Key 名称')
    return
  }
  editLoading.value = true
  try {
    await updateMyKey(editTarget.value.id, {
      name: editForm.name!.trim(),
      quota_limit: editForm.quota_limit,
      rpm: editForm.rpm,
      tpm: editForm.tpm,
    })
    editVisible.value = false
    ElMessage.success('已更新')
    await fetchList()
  } catch (e: any) {
    ElMessage.error(e?.message || '更新失败')
  } finally {
    editLoading.value = false
  }
}

// ========== 启用/禁用 ==========
async function toggleEnabled(row: ApiKey, val: boolean) {
  try {
    await updateMyKey(row.id, { enabled: val })
    row.enabled = val
  } catch (e: any) {
    ElMessage.error(e?.message || '操作失败')
    row.enabled = !val // 回滚
  }
}

// ========== 删除 ==========
async function handleDelete(row: ApiKey) {
  await ElMessageBox.confirm(
    `确认删除 Key「${row.name}」？此操作不可撤销，使用该 Key 的请求将立即失效。`,
    '删除确认',
    { type: 'warning', confirmButtonText: '确认删除', cancelButtonText: '取消' },
  )
  try {
    await deleteMyKey(row.id)
    ElMessage.success('已删除')
    await fetchList()
  } catch (e: any) {
    ElMessage.error(e?.message || '删除失败')
  }
}

// ========== 工具 ==========
function formatQuota(v: number) {
  if (!v || v === 0) return '不限'
  if (v >= 10_000) return (v / 10_000).toFixed(2) + ' 积分'
  return v + ' 厘'
}

function formatLastUsed(row: ApiKey) {
  return nullVal(row.last_used_at) ? formatDateTime(row.last_used_at) : '从未使用'
}

onMounted(fetchList)
</script>

<template>
  <div class="keys-page">
    <!-- 操作栏 -->
    <div class="toolbar">
      <span class="toolbar-title">我的 API Key</span>
      <el-button type="primary" @click="openCreate">
        <el-icon><Plus /></el-icon> 新建 Key
      </el-button>
    </div>

    <!-- 列表 -->
    <el-table :data="rows" v-loading="loading" border stripe class="keys-table">
      <el-table-column label="名称" prop="name" min-width="140" />
      <el-table-column label="Key 前缀" prop="key_prefix" min-width="160">
        <template #default="{ row }">
          <el-tag type="info" class="mono">{{ row.key_prefix }}…</el-tag>
        </template>
      </el-table-column>
      <el-table-column label="配额上限" min-width="110">
        <template #default="{ row }">{{ formatQuota(row.quota_limit) }}</template>
      </el-table-column>
      <el-table-column label="已用" min-width="110">
        <template #default="{ row }">{{ formatQuota(row.quota_used) }}</template>
      </el-table-column>
      <el-table-column label="RPM" prop="rpm" min-width="80">
        <template #default="{ row }">{{ row.rpm || '不限' }}</template>
      </el-table-column>
      <el-table-column label="最后使用" min-width="150">
        <template #default="{ row }">{{ formatLastUsed(row) }}</template>
      </el-table-column>
      <el-table-column label="创建时间" min-width="150">
        <template #default="{ row }">{{ formatDateTime(row.created_at) }}</template>
      </el-table-column>
      <el-table-column label="启用" width="80" align="center">
        <template #default="{ row }">
          <el-switch
            :model-value="row.enabled"
            @change="(v: string | number | boolean) => toggleEnabled(row, v as boolean)"
          />
        </template>
      </el-table-column>
      <el-table-column label="操作" width="120" align="center" fixed="right">
        <template #default="{ row }">
          <el-button link type="primary" size="small" @click="openEdit(row)">编辑</el-button>
          <el-button link type="danger" size="small" @click="handleDelete(row)">删除</el-button>
        </template>
      </el-table-column>
    </el-table>

    <el-pagination
      v-if="total > pager.page_size"
      class="pagination"
      background
      layout="total, prev, pager, next"
      :total="total"
      :page-size="pager.page_size"
      :current-page="pager.page"
      @current-change="onPageChange"
    />

    <!-- 新建对话框 -->
    <el-dialog v-model="createVisible" title="新建 API Key" width="480px" :close-on-click-modal="false">
      <el-form label-width="90px">
        <el-form-item label="名称" required>
          <el-input v-model="createForm.name" placeholder="用于识别此 Key 的备注名" maxlength="64" show-word-limit />
        </el-form-item>
        <el-form-item label="配额上限">
          <el-input-number v-model="createForm.quota_limit" :min="0" :step="10000" style="width:100%" />
          <div class="form-hint">单位：厘（10000 厘 = 1 积分），0 表示不限</div>
        </el-form-item>
        <el-form-item label="RPM 限速">
          <el-input-number v-model="createForm.rpm" :min="0" style="width:100%" />
          <div class="form-hint">每分钟最大请求数，0 表示不限</div>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="createVisible = false">取消</el-button>
        <el-button type="primary" :loading="createLoading" @click="submitCreate">创建</el-button>
      </template>
    </el-dialog>

    <!-- 新 Key 展示（仅此一次） -->
    <el-dialog v-model="newKeyVisible" title="Key 已创建 — 请立即保存" width="520px" :close-on-click-modal="false">
      <el-alert type="warning" :closable="false" style="margin-bottom:16px">
        此密钥只显示一次，关闭后将无法再查看完整内容，请立即复制保存。
      </el-alert>
      <el-input :value="newKeyValue" readonly class="mono-input">
        <template #append>
          <el-button @click="copyKey">复制</el-button>
        </template>
      </el-input>
      <template #footer>
        <el-button type="primary" @click="newKeyVisible = false">我已保存，关闭</el-button>
      </template>
    </el-dialog>

    <!-- 编辑对话框 -->
    <el-dialog v-model="editVisible" title="编辑 API Key" width="480px" :close-on-click-modal="false">
      <el-form label-width="90px">
        <el-form-item label="名称" required>
          <el-input v-model="editForm.name" placeholder="备注名" maxlength="64" show-word-limit />
        </el-form-item>
        <el-form-item label="配额上限">
          <el-input-number v-model="editForm.quota_limit" :min="0" :step="10000" style="width:100%" />
          <div class="form-hint">单位：厘（10000 厘 = 1 积分），0 表示不限</div>
        </el-form-item>
        <el-form-item label="RPM 限速">
          <el-input-number v-model="editForm.rpm" :min="0" style="width:100%" />
          <div class="form-hint">每分钟最大请求数，0 表示不限</div>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="editVisible = false">取消</el-button>
        <el-button type="primary" :loading="editLoading" @click="submitEdit">保存</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<style scoped>
.keys-page {
  padding: 20px;
}
.toolbar {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 16px;
}
.toolbar-title {
  font-size: 16px;
  font-weight: 600;
}
.keys-table {
  width: 100%;
}
.mono {
  font-family: 'SFMono-Regular', Consolas, monospace;
  font-size: 12px;
}
.mono-input :deep(input) {
  font-family: 'SFMono-Regular', Consolas, monospace;
  font-size: 13px;
}
.pagination {
  margin-top: 16px;
  justify-content: flex-end;
}
.form-hint {
  font-size: 12px;
  color: var(--el-text-color-secondary);
  margin-top: 4px;
  line-height: 1.4;
}
</style>
