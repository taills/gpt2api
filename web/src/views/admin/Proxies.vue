<template>
  <div class="proxies-page">
    <!-- 工具栏 -->
    <el-row class="toolbar" :gutter="12" align="middle">
      <el-col :span="12">
        <el-button type="primary" @click="openCreate">新增代理</el-button>
        <el-button @click="openImport">批量导入</el-button>
        <el-button :loading="probingAll" @click="probeAll">全量探测</el-button>
      </el-col>
      <el-col :span="12" style="text-align:right">
        <el-button :icon="RefreshRight" circle @click="load" />
      </el-col>
    </el-row>

    <!-- 表格 -->
    <el-table :data="rows" stripe border v-loading="loading" style="width:100%">
      <el-table-column label="ID" prop="id" width="70" />
      <el-table-column label="方案" prop="scheme" width="90" />
      <el-table-column label="地址" min-width="180">
        <template #default="{ row }">{{ row.host }}:{{ row.port }}</template>
      </el-table-column>
      <el-table-column label="用户名" prop="username" min-width="120" show-overflow-tooltip />
      <el-table-column label="国家/运营商" min-width="120">
        <template #default="{ row }">
          <span v-if="row.country || row.isp">{{ [row.country, row.isp].filter(Boolean).join(' / ') }}</span>
          <span v-else class="dim">-</span>
        </template>
      </el-table-column>
      <el-table-column label="健康分" width="90" align="center">
        <template #default="{ row }">
          <el-tag :type="scoreType(row.health_score)" size="small">{{ row.health_score }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column label="状态" width="80" align="center">
        <template #default="{ row }">
          <el-switch
            :model-value="row.enabled"
            @change="(v: boolean) => toggleEnabled(row, v)"
          />
        </template>
      </el-table-column>
      <el-table-column label="备注" prop="remark" min-width="120" show-overflow-tooltip />
      <el-table-column label="最后探测" min-width="160">
        <template #default="{ row }">
          <span v-if="lastProbeAt(row)" class="time">{{ lastProbeAt(row) }}</span>
          <span v-else class="dim">-</span>
        </template>
      </el-table-column>
      <el-table-column label="最后错误" min-width="160" show-overflow-tooltip>
        <template #default="{ row }">
          <span v-if="row.last_error" class="err">{{ row.last_error }}</span>
          <span v-else class="dim">-</span>
        </template>
      </el-table-column>
      <el-table-column label="创建时间" width="160">
        <template #default="{ row }">{{ fmtTime(row.created_at) }}</template>
      </el-table-column>
      <el-table-column label="操作" width="180" fixed="right">
        <template #default="{ row }">
          <el-button size="small" @click="openEdit(row)">编辑</el-button>
          <el-button
            size="small"
            :loading="probingId === row.id"
            @click="probeOne(row)"
          >探测</el-button>
          <el-button size="small" type="danger" @click="remove(row)">删除</el-button>
        </template>
      </el-table-column>
    </el-table>

    <!-- 分页 -->
    <el-pagination
      class="pager"
      background
      layout="total, sizes, prev, pager, next"
      :total="total"
      v-model:current-page="pager.page"
      v-model:page-size="pager.page_size"
      :page-sizes="[20, 50, 100]"
      @change="load"
    />

    <!-- 新增/编辑对话框 -->
    <el-dialog v-model="dialogVisible" :title="editing ? '编辑代理' : '新增代理'" width="520px">
      <el-form :model="form" label-width="90px" ref="formRef" :rules="rules">
        <el-form-item label="方案" prop="scheme">
          <el-select v-model="form.scheme" style="width:100%">
            <el-option label="http" value="http" />
            <el-option label="https" value="https" />
            <el-option label="socks5" value="socks5" />
          </el-select>
        </el-form-item>
        <el-form-item label="主机" prop="host">
          <el-input v-model="form.host" placeholder="如 127.0.0.1 或 proxy.example.com" />
        </el-form-item>
        <el-form-item label="端口" prop="port">
          <el-input-number v-model="form.port" :min="1" :max="65535" style="width:100%" />
        </el-form-item>
        <el-form-item label="用户名">
          <el-input v-model="form.username" placeholder="无需认证可留空" />
        </el-form-item>
        <el-form-item label="密码">
          <el-input
            v-model="form.password"
            type="password"
            show-password
            :placeholder="editing ? '留空不修改密码' : '无需认证可留空'"
          />
        </el-form-item>
        <el-form-item label="国家">
          <el-input v-model="form.country" placeholder="如 CN" />
        </el-form-item>
        <el-form-item label="运营商">
          <el-input v-model="form.isp" placeholder="如 Aliyun" />
        </el-form-item>
        <el-form-item label="启用">
          <el-switch v-model="form.enabled" />
        </el-form-item>
        <el-form-item label="备注">
          <el-input v-model="form.remark" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="saving" @click="save">保存</el-button>
      </template>
    </el-dialog>

    <!-- 批量导入对话框 -->
    <el-dialog v-model="importVisible" title="批量导入代理" width="620px">
      <el-form :model="importForm" label-width="90px">
        <el-form-item label="代理列表">
          <el-input
            v-model="importForm.text"
            type="textarea"
            :rows="8"
            placeholder="每行一条，格式：scheme://[user:pass@]host:port 或 host:port:user:pass"
          />
        </el-form-item>
        <el-form-item label="默认启用">
          <el-switch v-model="importForm.enabled" />
        </el-form-item>
        <el-form-item label="国家">
          <el-input v-model="importForm.country" />
        </el-form-item>
        <el-form-item label="运营商">
          <el-input v-model="importForm.isp" />
        </el-form-item>
        <el-form-item label="备注">
          <el-input v-model="importForm.remark" />
        </el-form-item>
        <el-form-item label="覆盖已有">
          <el-switch v-model="importForm.overwrite" />
          <span class="hint">开启后已存在的代理会被更新</span>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="importVisible = false">取消</el-button>
        <el-button type="primary" :loading="importing" @click="doImport">导入</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { RefreshRight } from '@element-plus/icons-vue'
import { ElMessage, ElMessageBox, ElNotification } from 'element-plus'
import type { FormInstance, FormRules } from 'element-plus'
import * as proxyApi from '@/api/proxies'
import type { Proxy } from '@/api/proxies'

// ——— 列表 ———
const rows = ref<Proxy[]>([])
const total = ref(0)
const loading = ref(false)
const pager = reactive({ page: 1, page_size: 20 })

async function load() {
  loading.value = true
  try {
    const res = await proxyApi.listProxies({ page: pager.page, page_size: pager.page_size })
    rows.value = res.list ?? []
    total.value = res.total ?? 0
  } catch (e: any) {
    ElMessage.error(e?.message ?? '加载失败')
  } finally {
    loading.value = false
  }
}

onMounted(load)

// ——— 健康分颜色 ———
function scoreType(score: number): 'success' | 'warning' | 'danger' {
  if (score >= 71) return 'success'
  if (score >= 41) return 'warning'
  return 'danger'
}

// ——— 时间格式化 ———
function fmtTime(s?: string | null): string {
  if (!s) return '-'
  return new Date(s).toLocaleString()
}

function lastProbeAt(row: Proxy): string {
  const v = row.last_probe_at
  if (!v) return ''
  if (typeof v === 'string') return fmtTime(v)
  if (typeof v === 'object' && v.Valid) return fmtTime(v.Time)
  return ''
}

// ——— 启用/禁用 ———
async function toggleEnabled(row: Proxy, enabled: boolean) {
  try {
    await proxyApi.updateProxy(row.id, { scheme: row.scheme, host: row.host, port: row.port, enabled })
    row.enabled = enabled
  } catch (e: any) {
    ElMessage.error(e?.message ?? '操作失败')
  }
}

// ——— 新增/编辑 ———
const dialogVisible = ref(false)
const editing = ref<Proxy | null>(null)
const saving = ref(false)
const formRef = ref<FormInstance>()

const form = reactive<proxyApi.ProxyCreate & { password: string }>({
  scheme: 'http',
  host: '',
  port: 8080,
  username: '',
  password: '',
  country: '',
  isp: '',
  enabled: true,
  remark: '',
})

const rules: FormRules = {
  scheme: [{ required: true, message: '请选择方案' }],
  host: [{ required: true, message: '请填写主机' }],
  port: [{ required: true, message: '请填写端口' }],
}

function openCreate() {
  editing.value = null
  Object.assign(form, { scheme: 'http', host: '', port: 8080, username: '', password: '', country: '', isp: '', enabled: true, remark: '' })
  dialogVisible.value = true
}

function openEdit(row: Proxy) {
  editing.value = row
  Object.assign(form, { scheme: row.scheme, host: row.host, port: row.port, username: row.username, password: '', country: row.country, isp: row.isp, enabled: row.enabled, remark: row.remark })
  dialogVisible.value = true
}

async function save() {
  await formRef.value?.validate()
  saving.value = true
  const body: proxyApi.ProxyCreate = {
    scheme: form.scheme, host: form.host, port: form.port,
    username: form.username || undefined,
    country: form.country || undefined,
    isp: form.isp || undefined,
    enabled: form.enabled,
    remark: form.remark || undefined,
  }
  if (form.password) body.password = form.password
  try {
    if (editing.value) {
      await proxyApi.updateProxy(editing.value.id, body)
      ElMessage.success('更新成功')
    } else {
      await proxyApi.createProxy(body)
      ElMessage.success('创建成功')
    }
    dialogVisible.value = false
    load()
  } catch (e: any) {
    ElMessage.error(e?.message ?? '保存失败')
  } finally {
    saving.value = false
  }
}

// ——— 删除 ———
async function remove(row: Proxy) {
  await ElMessageBox.confirm(`确定删除代理 ${row.host}:${row.port} 吗？`, '删除确认', { type: 'warning' })
  try {
    await proxyApi.deleteProxy(row.id)
    ElMessage.success('已删除')
    load()
  } catch (e: any) {
    ElMessage.error(e?.message ?? '删除失败')
  }
}

// ——— 单条探测 ———
const probingId = ref<number | null>(null)

async function probeOne(row: Proxy) {
  probingId.value = row.id
  try {
    const res = await proxyApi.probeProxy(row.id)
    if (res.ok) {
      ElMessage.success(`探测成功，延迟 ${res.latency_ms} ms`)
    } else {
      ElMessage.warning(`探测失败：${res.error ?? '未知'}`)
    }
    row.health_score = res.health_score
    load()
  } catch (e: any) {
    ElMessage.error(e?.message ?? '探测失败')
  } finally {
    probingId.value = null
  }
}

// ——— 全量探测 ———
const probingAll = ref(false)

async function probeAll() {
  probingAll.value = true
  try {
    const res = await proxyApi.probeAllProxies()
    ElNotification({
      title: '全量探测完成',
      message: `共 ${res.total} 条，成功 ${res.ok}，失败 ${res.bad}`,
      type: res.bad === 0 ? 'success' : 'warning',
      duration: 5000,
    })
    load()
  } catch (e: any) {
    ElMessage.error(e?.message ?? '探测失败')
  } finally {
    probingAll.value = false
  }
}

// ——— 批量导入 ———
const importVisible = ref(false)
const importing = ref(false)
const importForm = reactive({
  text: '',
  enabled: true,
  country: '',
  isp: '',
  remark: '',
  overwrite: false,
})

function openImport() {
  Object.assign(importForm, { text: '', enabled: true, country: '', isp: '', remark: '', overwrite: false })
  importVisible.value = true
}

async function doImport() {
  if (!importForm.text.trim()) {
    ElMessage.warning('请填写代理列表')
    return
  }
  importing.value = true
  try {
    const res = await proxyApi.importProxies({
      text: importForm.text,
      enabled: importForm.enabled,
      country: importForm.country || undefined,
      isp: importForm.isp || undefined,
      remark: importForm.remark || undefined,
      overwrite: importForm.overwrite,
    })
    ElMessage.success(`导入完成：新增 ${res.created}，更新 ${res.updated}，跳过 ${res.skipped}，无效 ${res.invalid}`)
    importVisible.value = false
    load()
  } catch (e: any) {
    ElMessage.error(e?.message ?? '导入失败')
  } finally {
    importing.value = false
  }
}
</script>

<style scoped>
.proxies-page { padding: 16px; }
.toolbar { margin-bottom: 12px; }
.pager { margin-top: 16px; justify-content: flex-end; display: flex; }
.dim { color: #999; }
.err { color: #f56c6c; }
.time { color: #606266; }
.hint { margin-left: 8px; color: #909399; font-size: 12px; }
</style>
