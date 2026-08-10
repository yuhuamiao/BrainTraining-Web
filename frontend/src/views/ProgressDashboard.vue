<template>
  <section>
    <div class="page-heading">
      <div>
        <p class="eyebrow">Progress</p>
        <h1>训练记录</h1>
        <p>按项目回看完成次数、准确率和用时。</p>
      </div>
    </div>

    <div class="progress-layout">
      <aside class="profile-panel">
        <span class="profile-avatar">{{ initial }}</span>
        <div class="profile-name-row">
          <div>
            <span>训练者</span>
            <strong>{{ authStore.user?.username }}</strong>
          </div>
          <button class="button button-ghost button-small" type="button" @click="editing = !editing">{{ editing ? '取消' : '编辑' }}</button>
        </div>
        <form v-if="editing" class="name-form" @submit.prevent="saveName">
          <input v-model.trim="username" class="input" minlength="3" maxlength="50" required />
          <button class="button button-small" :disabled="saving">保存</button>
        </form>
        <p v-if="profileMessage" class="profile-message">{{ profileMessage }}</p>

        <dl class="profile-stats">
          <div><dt>总训练</dt><dd>{{ records.length }}</dd></div>
          <div><dt>覆盖项目</dt><dd>{{ completedTypes }} / 5</dd></div>
          <div><dt>平均准确率</dt><dd>{{ averageAccuracy }}%</dd></div>
        </dl>
      </aside>

      <div class="records-panel">
        <div class="records-toolbar">
          <div class="segmented score-filters">
            <button v-for="filter in filters" :key="filter.key" :class="{ active: activeFilter === filter.key }" type="button" @click="activeFilter = filter.key">{{ filter.label }}</button>
          </div>
          <span>{{ filteredRecords.length }} 条</span>
        </div>

        <div v-if="loading" class="records-empty">正在读取记录…</div>
        <div v-else-if="error" class="records-empty error-text">{{ error }}</div>
        <div v-else-if="!filteredRecords.length" class="records-empty">
          <strong>还没有训练记录</strong>
          <router-link class="button" to="/">开始第一次训练</router-link>
        </div>
        <div v-else class="records-list">
          <article v-for="record in filteredRecords" :key="`${record.type}-${record.id}`" class="record-row">
            <span class="record-code" :style="{ background: record.color }">{{ record.code }}</span>
            <div class="record-main">
              <strong>{{ record.label }}</strong>
              <span>{{ formatDate(record.createdAt) }}</span>
            </div>
            <div><span>结果</span><strong>{{ resultLabel(record) }}</strong></div>
            <div><span>{{ record.timeElapsed ? '用时' : '难度' }}</span><strong>{{ record.timeElapsed ? `${record.timeElapsed.toFixed(1)} 秒` : difficultyLabel(record.level) }}</strong></div>
          </article>
        </div>
      </div>
    </div>
  </section>
</template>

<script setup>
import { computed, onMounted, ref } from 'vue'
import { apiMessage, getUserScores, updateMe } from '@/services/api'
import { useAuthStore } from '@/stores/authStore'

const authStore = useAuthStore()
const loading = ref(true)
const error = ref('')
const rawScores = ref({})
const activeFilter = ref('all')
const editing = ref(false)
const saving = ref(false)
const username = ref(authStore.user?.username || '')
const profileMessage = ref('')

const metadata = {
  schulte: { label: '舒尔特方格', code: 'S', color: '#2e6b50' },
  colorWord: { label: '多色文字', code: 'C', color: '#dd6b52' },
  memory: { label: '瞬时记忆', code: 'M', color: '#3d6f9f' },
  bus: { label: '公交人数', code: 'B', color: '#d29a1d' },
  sudoku: { label: '数独挑战', code: '9', color: '#725c8e' },
}
const filters = [{ key: 'all', label: '全部' }, ...Object.entries(metadata).map(([key, value]) => ({ key, label: value.label }))]
const records = computed(() => Object.entries(metadata).flatMap(([type, meta]) => (rawScores.value[type] || []).map((record) => ({ ...record, ...meta, type }))).sort((a, b) => new Date(b.createdAt) - new Date(a.createdAt)))
const filteredRecords = computed(() => activeFilter.value === 'all' ? records.value : records.value.filter((record) => record.type === activeFilter.value))
const completedTypes = computed(() => Object.keys(metadata).filter((key) => rawScores.value[key]?.length).length)
const accuracyRecords = computed(() => records.value.filter((record) => typeof record.accuracy === 'number'))
const averageAccuracy = computed(() => accuracyRecords.value.length ? Math.round(accuracyRecords.value.reduce((sum, record) => sum + record.accuracy, 0) / accuracyRecords.value.length * 100) : 0)
const initial = computed(() => (authStore.user?.username || '训').slice(0, 1).toUpperCase())

function resultLabel(record) {
  if (record.type === 'schulte') return `${record.successNum || 0} / 25`
  if (record.type === 'sudoku') return record.isPassed ? `完成 · ${record.mistakes || 0} 错` : '未完成'
  return `${Math.round((record.accuracy || 0) * 100)}%`
}
function difficultyLabel(level) { return ({ easy: '简单', medium: '中等', hard: '困难' })[level] || '—' }
function formatDate(value) { return new Intl.DateTimeFormat('zh-CN', { month: '2-digit', day: '2-digit', hour: '2-digit', minute: '2-digit' }).format(new Date(value)) }

async function saveName() {
  saving.value = true
  profileMessage.value = ''
  try {
    await updateMe({ username: username.value })
    authStore.setUser({ ...authStore.user, username: username.value })
    editing.value = false
    profileMessage.value = '用户名已更新'
  } catch (requestError) {
    profileMessage.value = apiMessage(requestError)
  } finally {
    saving.value = false
  }
}

onMounted(async () => {
  try { rawScores.value = await getUserScores() }
  catch (requestError) { error.value = apiMessage(requestError) }
  finally { loading.value = false }
})
</script>

<style scoped>
.progress-layout { display: grid; grid-template-columns: 270px minmax(0, 1fr); gap: 28px; align-items: start; }
.profile-panel { position: sticky; top: 100px; padding-right: 28px; border-right: 1px solid var(--line); }
.profile-avatar { width: 64px; aspect-ratio: 1; display: grid; place-items: center; border-radius: 50%; background: var(--ink); color: white; font-size: 24px; font-weight: 800; }
.profile-name-row { display: flex; align-items: end; justify-content: space-between; margin: 18px 0; }
.profile-name-row div { min-width: 0; display: grid; gap: 3px; }
.profile-name-row span { color: var(--muted); font-size: 12px; }
.profile-name-row strong { overflow-wrap: anywhere; font-size: 20px; }
.name-form { display: grid; grid-template-columns: 1fr auto; gap: 8px; margin-bottom: 12px; }
.profile-message { color: var(--muted); font-size: 12px; }
.profile-stats { margin: 28px 0 0; border-top: 1px solid var(--line); }
.profile-stats div { display: flex; justify-content: space-between; padding: 15px 0; border-bottom: 1px solid var(--line); }
.profile-stats dt { color: var(--muted); font-size: 13px; }
.profile-stats dd { margin: 0; font-weight: 800; }
.records-panel { min-width: 0; }
.records-toolbar { min-height: 48px; display: flex; align-items: center; justify-content: space-between; gap: 16px; margin-bottom: 14px; }
.records-toolbar > span { color: var(--muted); font-size: 12px; white-space: nowrap; }
.score-filters { min-width: 0; max-width: 100%; flex: 1; overflow-x: auto; }
.score-filters button { white-space: nowrap; }
.records-list { border-top: 1px solid var(--line); }
.record-row { min-height: 82px; display: grid; grid-template-columns: 38px minmax(150px, 1fr) 120px 120px; align-items: center; gap: 16px; border-bottom: 1px solid var(--line); }
.record-code { width: 34px; aspect-ratio: 1; display: grid; place-items: center; border-radius: 5px; color: white; font-weight: 850; }
.record-main, .record-row > div:not(.record-main) { display: grid; gap: 4px; }
.record-main span, .record-row > div:not(.record-main) span { color: var(--muted); font-size: 11px; }
.record-row > div:not(.record-main) strong { font-size: 14px; }
.records-empty { min-height: 330px; display: flex; flex-direction: column; align-items: center; justify-content: center; gap: 20px; border-block: 1px solid var(--line); color: var(--muted); text-align: center; }
.error-text { color: var(--red); }
@media (max-width: 760px) {
  .progress-layout { grid-template-columns: 1fr; }
  .profile-panel { position: static; padding: 0 0 24px; border-right: 0; border-bottom: 1px solid var(--line); }
  .record-row { grid-template-columns: 34px minmax(120px, 1fr) 80px; gap: 10px; }
  .record-row > div:last-child { display: none; }
  .records-toolbar { align-items: start; }
}
</style>
