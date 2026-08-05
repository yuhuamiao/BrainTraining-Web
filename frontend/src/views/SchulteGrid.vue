<template>
  <section class="game-page">
    <div class="game-topline">
      <router-link class="back-link" to="/">返回训练</router-link>
      <span class="eyebrow">视觉搜索 · 30 秒</span>
    </div>

    <div class="game-surface">
      <header class="game-header">
        <h1>舒尔特方格</h1>
        <div class="game-meta">
          <div><span>下一个</span><strong>{{ phase === 'playing' ? Math.min(nextNumber, 25) : '—' }}</strong></div>
          <div><span>剩余</span><strong>{{ remaining.toFixed(1) }}s</strong></div>
        </div>
      </header>

      <div v-if="phase === 'setup'" class="game-body game-setup">
        <div class="setup-inner">
          <div class="setup-board" aria-hidden="true"><i>1</i><i>9</i><i>4</i><i>7</i></div>
          <h2>从 1 找到 25</h2>
          <p>保持视线在表格中央，按数字顺序点击。</p>
          <button class="button" type="button" :disabled="loading" @click="startGame">{{ loading ? '准备中…' : '开始 30 秒训练' }}</button>
          <p v-if="error" class="status-line form-error">{{ error }}</p>
        </div>
      </div>

      <div v-else-if="phase === 'playing'" class="game-body schulte-play">
        <div class="progress-track"><i :style="{ width: `${(nextNumber - 1) / 25 * 100}%` }"></i></div>
        <div class="schulte-grid">
          <button
            v-for="number in numbers"
            :key="number"
            type="button"
            :class="['schulte-cell', `tone-${number % 5}`, { done: number < nextNumber, wrong: wrongNumber === number }]"
            :disabled="number < nextNumber"
            :aria-label="`数字 ${number}`"
            @click="selectNumber(number)"
          >{{ number }}</button>
        </div>
        <p class="status-line">{{ feedback }}</p>
      </div>

      <div v-else class="game-body game-result">
        <div class="result-inner">
          <p class="eyebrow">{{ completed ? 'Completed' : 'Time up' }}</p>
          <h2>{{ completed ? '顺序完成' : `找到 ${nextNumber - 1} 个数字` }}</h2>
          <div class="result-number">{{ elapsed.toFixed(1) }}s</div>
          <p>{{ saveMessage }}</p>
          <div class="result-actions">
            <button class="button" type="button" @click="startGame">再练一次</button>
            <router-link class="button button-secondary" to="/">选择其他训练</router-link>
          </div>
        </div>
      </div>
    </div>
  </section>
</template>

<script setup>
import { computed, onBeforeUnmount, onMounted, ref } from 'vue'
import { apiMessage, getSchulteMatrix, submitSchulteScore } from '@/services/api'

const phase = ref('setup')
const matrix = ref([])
const loading = ref(false)
const error = ref('')
const nextNumber = ref(1)
const remaining = ref(30)
const elapsed = ref(0)
const completed = ref(false)
const wrongNumber = ref(null)
const feedback = ref('点击数字 1 开始')
const saveMessage = ref('正在保存成绩…')
const numbers = computed(() => matrix.value.flat())
let timerId = null
let wrongTimer = null
let startedAt = 0

async function loadMatrix() {
  loading.value = true
  error.value = ''
  try {
    matrix.value = (await getSchulteMatrix()).matrix
  } catch (requestError) {
    error.value = apiMessage(requestError, '无法获取方格')
  } finally {
    loading.value = false
  }
}

async function startGame() {
  clearTimers()
  await loadMatrix()
  if (!matrix.value.length) return
  phase.value = 'playing'
  nextNumber.value = 1
  remaining.value = 30
  elapsed.value = 0
  completed.value = false
  feedback.value = '点击数字 1 开始'
  saveMessage.value = '正在保存成绩…'
  startedAt = performance.now()
  timerId = window.setInterval(tick, 50)
}

function tick() {
  elapsed.value = (performance.now() - startedAt) / 1000
  remaining.value = Math.max(0, 30 - elapsed.value)
  if (remaining.value <= 0) finish(false)
}

function selectNumber(number) {
  if (number !== nextNumber.value) {
    wrongNumber.value = number
    feedback.value = `下一个是 ${nextNumber.value}`
    window.clearTimeout(wrongTimer)
    wrongTimer = window.setTimeout(() => { wrongNumber.value = null }, 260)
    return
  }
  nextNumber.value += 1
  feedback.value = nextNumber.value <= 25 ? `继续找 ${nextNumber.value}` : '完成'
  if (nextNumber.value === 26) finish(true)
}

async function finish(isComplete) {
  if (phase.value !== 'playing') return
  window.clearInterval(timerId)
  elapsed.value = Math.min(30, (performance.now() - startedAt) / 1000)
  remaining.value = Math.max(0, 30 - elapsed.value)
  completed.value = isComplete
  phase.value = 'result'
  try {
    await submitSchulteScore({ isPassed: isComplete, successNum: nextNumber.value - 1, trainingNum: 25, timeElapsed: Number(elapsed.value.toFixed(2)) })
    saveMessage.value = '成绩已保存'
  } catch (requestError) {
    saveMessage.value = apiMessage(requestError, '成绩暂未保存')
  }
}

function clearTimers() {
  window.clearInterval(timerId)
  window.clearTimeout(wrongTimer)
}

onMounted(loadMatrix)
onBeforeUnmount(clearTimers)
</script>

<style scoped>
.setup-board { width: 126px; display: grid; grid-template-columns: 1fr 1fr; gap: 5px; margin: 0 auto 24px; transform: rotate(-4deg); }
.setup-board i { aspect-ratio: 1; display: grid; place-items: center; border: 1px solid var(--line); background: #f1f4f1; color: var(--green); font-size: 24px; font-style: normal; font-weight: 800; }
.schulte-play { display: grid; justify-items: center; }
.progress-track { width: min(100%, 560px); height: 4px; margin-bottom: 20px; overflow: hidden; border-radius: 2px; background: #e4e9e5; }
.progress-track i { display: block; height: 100%; background: var(--green); transition: width 120ms linear; }
.schulte-grid { width: min(100%, 560px); aspect-ratio: 1; display: grid; grid-template-columns: repeat(5, 1fr); gap: 7px; }
.schulte-cell { min-width: 0; border: 1px solid #cdd5cf; border-radius: 5px; background: white; color: var(--ink); font-size: 24px; font-weight: 800; transition: transform 100ms ease, background 140ms ease, opacity 140ms ease; }
.schulte-cell:hover:not(:disabled) { transform: translateY(-2px); border-color: var(--green); }
.schulte-cell.done { border-color: transparent; background: #e6ede8; color: #8b9690; opacity: 0.64; }
.schulte-cell.wrong { background: #f9e2de; transform: scale(0.95); }
.tone-0 { color: #2e6b50; }.tone-1 { color: #3d6f9f; }.tone-2 { color: #9c5e40; }.tone-3 { color: #725c8e; }.tone-4 { color: #5a625e; }
@media (max-width: 600px) { .schulte-grid { gap: 4px; } .schulte-cell { font-size: 18px; } }
</style>
