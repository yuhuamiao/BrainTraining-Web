<template>
  <section class="game-page">
    <div class="game-topline"><router-link class="back-link" to="/">返回训练</router-link><span class="eyebrow">抑制控制 · 30 秒</span></div>
    <div class="game-surface">
      <header class="game-header">
        <h1>多色文字</h1>
        <div class="game-meta">
          <div><span>正确</span><strong>{{ correct }}</strong></div>
          <div><span>剩余</span><strong>{{ remaining }}s</strong></div>
        </div>
      </header>

      <div v-if="phase === 'setup'" class="game-body game-setup">
        <div class="setup-inner">
          <div class="stroop-sample" aria-hidden="true"><i>红</i><i>蓝</i><i>绿</i></div>
          <h2>判断文字的颜色</h2>
          <p>忽略文字含义，选择当前高亮文字的实际颜色。</p>
          <div class="segmented difficulty-control">
            <button v-for="item in difficulties" :key="item.value" :class="{ active: difficulty === item.value }" type="button" @click="difficulty = item.value">{{ item.label }}</button>
          </div>
          <button class="button start-button" type="button" :disabled="loading" @click="startGame">{{ loading ? '准备中…' : '开始训练' }}</button>
          <p v-if="error" class="status-line form-error">{{ error }}</p>
        </div>
      </div>

      <div v-else-if="phase === 'playing'" class="game-body color-play">
        <div class="color-grid" :style="{ gridTemplateColumns: `repeat(${matrix.length}, 1fr)` }">
          <div v-for="(cell, index) in flatMatrix" :key="index" :class="['color-cell', { active: index === activeCell }]" :style="index === activeCell ? { color: colorMap[currentColor] } : undefined">{{ cell.word }}</div>
        </div>
        <div class="answer-colors">
          <button v-for="color in colors" :key="color.value" type="button" :disabled="!answering" @click="answer(color.value)">
            <i :style="{ background: colorMap[color.value] }"></i>{{ color.label }}
          </button>
        </div>
        <p class="status-line">{{ feedback }}</p>
      </div>

      <div v-else class="game-body game-result">
        <div class="result-inner">
          <p class="eyebrow">Session complete</p><h2>抑制控制训练完成</h2>
          <div class="result-number">{{ accuracy }}%</div>
          <p>正确 {{ correct }} / {{ questions }} · {{ saveMessage }}</p>
          <div class="result-actions"><button class="button" type="button" @click="reset">再练一次</button><router-link class="button button-secondary" to="/">选择其他训练</router-link></div>
        </div>
      </div>
    </div>
  </section>
</template>

<script setup>
import { computed, onBeforeUnmount, ref } from 'vue'
import { apiMessage, getColorMatrix, submitColorWordScore } from '@/services/api'

const phase = ref('setup')
const difficulty = ref('easy')
const matrix = ref([])
const loading = ref(false)
const error = ref('')
const remaining = ref(30)
const currentColor = ref('red')
const activeCell = ref(-1)
const answering = ref(false)
const correct = ref(0)
const questions = ref(0)
const feedback = ref('等待第一个颜色')
const saveMessage = ref('正在保存成绩…')
const flatMatrix = computed(() => matrix.value.flat())
const accuracy = computed(() => questions.value ? Math.round(correct.value / questions.value * 100) : 0)
const difficulties = [{ value: 'easy', label: '简单' }, { value: 'medium', label: '中等' }, { value: 'hard', label: '困难' }]
const colors = [{ value: 'red', label: '红' }, { value: 'blue', label: '蓝' }, { value: 'green', label: '绿' }, { value: 'yellow', label: '黄' }, { value: 'black', label: '黑' }]
const colorMap = { red: '#d6534d', blue: '#356da0', green: '#2e7a52', yellow: '#c58e16', black: '#202724' }
let socket = null
let timerId = null
let fallbackId = null
let fallbackGuard = null

async function startGame() {
  cleanup()
  loading.value = true
  error.value = ''
  try {
    matrix.value = (await getColorMatrix(difficulty.value)).matrix
  } catch (requestError) {
    error.value = apiMessage(requestError, '无法获取文字矩阵')
    loading.value = false
    return
  }
  loading.value = false
  phase.value = 'playing'
  remaining.value = 30
  correct.value = 0
  questions.value = 0
  feedback.value = '等待第一个颜色'
  saveMessage.value = '正在保存成绩…'
  connectStream()
  timerId = window.setInterval(() => {
    remaining.value -= 1
    if (remaining.value <= 0) finish()
  }, 1000)
}

function connectStream() {
  const configured = import.meta.env.VITE_API_BASE_URL
  let url
  if (configured?.startsWith('http')) {
    const apiUrl = new URL(configured)
    apiUrl.protocol = apiUrl.protocol === 'https:' ? 'wss:' : 'ws:'
    apiUrl.pathname = `${apiUrl.pathname.replace(/\/$/, '')}/color_words/color`
    url = apiUrl.toString()
  } else {
    const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:'
    url = `${protocol}//${window.location.host}/api/v1/color_words/color`
  }
  socket = new WebSocket(url)
  socket.onmessage = (event) => nextPrompt(JSON.parse(event.data).color)
  socket.onerror = startFallback
  socket.onclose = () => { if (phase.value === 'playing' && questions.value === 0) startFallback() }
  fallbackGuard = window.setTimeout(() => { if (questions.value === 0) startFallback() }, 1800)
}

function startFallback() {
  if (fallbackId || phase.value !== 'playing') return
  socket?.close()
  nextPrompt(randomColor())
  fallbackId = window.setInterval(() => nextPrompt(randomColor()), 1000)
}

function randomColor() { return colors[Math.floor(Math.random() * colors.length)].value }
function nextPrompt(color) {
  if (phase.value !== 'playing' || !flatMatrix.value.length) return
  currentColor.value = colorMap[color] ? color : randomColor()
  activeCell.value = Math.floor(Math.random() * flatMatrix.value.length)
  answering.value = true
  questions.value += 1
  feedback.value = '选择文字显示的颜色'
}
function answer(value) {
  if (!answering.value) return
  answering.value = false
  if (value === currentColor.value) { correct.value += 1; feedback.value = '正确' }
  else feedback.value = '继续专注颜色'
}

async function finish() {
  if (phase.value !== 'playing') return
  cleanup()
  phase.value = 'result'
  try {
    await submitColorWordScore({ successNum: correct.value, trainingNum: Math.max(questions.value, 1), accuracy: questions.value ? correct.value / questions.value : 0, level: difficulty.value })
    saveMessage.value = '成绩已保存'
  } catch (requestError) { saveMessage.value = apiMessage(requestError, '成绩暂未保存') }
}
function reset() { phase.value = 'setup'; remaining.value = 30; error.value = '' }
function cleanup() { window.clearInterval(timerId); window.clearInterval(fallbackId); window.clearTimeout(fallbackGuard); fallbackId = null; socket?.close(); socket = null }
onBeforeUnmount(cleanup)
</script>

<style scoped>
.stroop-sample { display: flex; justify-content: center; gap: 10px; margin-bottom: 24px; font-size: 34px; font-weight: 850; }
.stroop-sample i { font-style: normal; }.stroop-sample i:nth-child(1) { color: var(--blue); }.stroop-sample i:nth-child(2) { color: var(--coral); }.stroop-sample i:nth-child(3) { color: var(--green); }
.difficulty-control { margin-bottom: 18px; }.start-button { display: flex; margin: 0 auto; }
.color-play { display: grid; justify-items: center; }
.color-grid { width: min(100%, 580px); aspect-ratio: 1; display: grid; gap: 4px; }
.color-cell { min-width: 0; display: grid; place-items: center; border: 1px solid #d5dbd6; border-radius: 3px; background: #f7f8f6; color: #adb5af; font-size: clamp(12px, 2.3vw, 21px); font-weight: 850; }
.color-cell.active { z-index: 1; border: 3px solid currentColor; background: white; box-shadow: 0 5px 14px rgba(23,32,29,.12); transform: scale(1.08); }
.answer-colors { display: flex; flex-wrap: wrap; justify-content: center; gap: 8px; margin-top: 24px; }
.answer-colors button { min-width: 86px; min-height: 42px; display: flex; align-items: center; justify-content: center; gap: 8px; border: 1px solid var(--line); border-radius: 5px; background: #eef1ee; color: var(--ink); font-weight: 750; }
.answer-colors button:hover:not(:disabled) { border-color: #aab4ad; background: white; }.answer-colors button:disabled { opacity: .46; }
.answer-colors i { width: 12px; aspect-ratio: 1; border-radius: 50%; }
@media (max-width: 600px) { .color-grid { gap: 2px; } .answer-colors button { min-width: 61px; padding: 0 9px; } }
</style>
