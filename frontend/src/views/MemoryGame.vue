<template>
  <section class="game-page">
    <div class="game-topline"><router-link class="back-link" to="/">返回训练</router-link><span class="eyebrow">空间记忆 · 10 轮</span></div>
    <div class="game-surface">
      <header class="game-header">
        <h1>瞬时记忆</h1>
        <div class="game-meta">
          <div><span>轮次</span><strong>{{ phase === 'setup' ? '—' : `${Math.min(round, 10)}/10` }}</strong></div>
          <div><span>{{ phase === 'memorize' ? '记忆' : phase === 'recall' ? '作答' : '阶段' }}</span><strong>{{ stageRemaining.toFixed(1) }}s</strong></div>
        </div>
      </header>

      <div v-if="phase === 'setup'" class="game-body game-setup">
        <div class="setup-inner">
          <div class="memory-sample" aria-hidden="true"><i></i><i class="on"></i><i></i><i class="on"></i><i></i><i></i><i class="on"></i><i></i><i></i><i class="on"></i><i></i><i></i></div>
          <h2>记住亮起的位置</h2><p>方格隐藏后，选回刚才亮起的全部位置。</p>
          <div class="segmented difficulty-control"><button v-for="item in difficulties" :key="item.value" :class="{ active: difficulty === item.value }" type="button" @click="difficulty = item.value">{{ item.label }}</button></div>
          <button class="button start-button" type="button" @click="startSession">开始 10 轮训练</button>
          <p v-if="error" class="status-line form-error">{{ error }}</p>
        </div>
      </div>

      <div v-else-if="phase !== 'result'" class="game-body memory-play">
        <div v-if="phase === 'transition' || phase === 'loading'" class="memory-transition"><i></i><span>{{ phase === 'loading' ? '生成下一组位置' : '重组方格' }}</span></div>
        <div v-else class="memory-grid" :style="gridStyle">
          <button v-for="index in totalCells" :key="index" type="button" :class="{ target: phase === 'memorize' && positions.includes(index - 1), selected: phase === 'recall' && selected.has(index - 1) }" :disabled="phase !== 'recall'" :aria-label="`方格 ${index}`" @click="toggleCell(index - 1)"></button>
        </div>
        <div v-if="phase === 'recall'" class="memory-actions"><span>已选 {{ selected.size }} / {{ positions.length }}</span><button class="button" type="button" @click="submitRound">确认位置</button></div>
        <p class="status-line">{{ phaseLabel }}</p>
      </div>

      <div v-else class="game-body game-result">
        <div class="result-inner"><p class="eyebrow">Session complete</p><h2>10 轮记忆完成</h2><div class="result-number">{{ successes }}/10</div><p>完整还原 {{ successes }} 轮 · {{ saveMessage }}</p><div class="result-actions"><button class="button" type="button" @click="reset">再练一次</button><router-link class="button button-secondary" to="/">选择其他训练</router-link></div></div>
      </div>
    </div>
  </section>
</template>

<script setup>
import { computed, onBeforeUnmount, ref } from 'vue'
import { apiMessage, getMemoryMatrix, submitMemoryScore } from '@/services/api'

const phase = ref('setup')
const difficulty = ref('easy')
const round = ref(0)
const successes = ref(0)
const positions = ref([])
const totalCells = ref(30)
const displayTime = ref(3)
const answerTime = ref(8)
const stageRemaining = ref(0)
const selected = ref(new Set())
const error = ref('')
const saveMessage = ref('正在保存成绩…')
const difficulties = [{ value: 'easy', label: '简单' }, { value: 'medium', label: '中等' }, { value: 'hard', label: '困难' }]
const columns = computed(() => ({ 30: 6, 42: 7, 56: 8 })[totalCells.value] || 6)
const rows = computed(() => totalCells.value / columns.value)
const gridStyle = computed(() => ({ gridTemplateColumns: `repeat(${columns.value}, 1fr)`, gridTemplateRows: `repeat(${rows.value}, 1fr)` }))
const phaseLabel = computed(() => phase.value === 'memorize' ? `记住 ${positions.value.length} 个蓝色方格` : phase.value === 'recall' ? '选回刚才亮起的位置' : '')
let stageTimer = null
let transitionTimer = null

function startSession() { round.value = 1; successes.value = 0; error.value = ''; saveMessage.value = '正在保存成绩…'; loadRound() }
async function loadRound() {
  clearTimers()
  phase.value = 'loading'
  selected.value = new Set()
  try {
    const response = await getMemoryMatrix(difficulty.value)
    positions.value = response.positions
    totalCells.value = response.totalCells
    displayTime.value = response.displayTime
    answerTime.value = response.answerTime
    phase.value = 'memorize'
    runStage(displayTime.value, beginRecall)
  } catch (requestError) {
    error.value = apiMessage(requestError, '无法生成记忆方格')
    phase.value = 'setup'
  }
}
function beginRecall() {
  phase.value = 'transition'; stageRemaining.value = 0
  transitionTimer = window.setTimeout(() => { phase.value = 'recall'; runStage(answerTime.value, submitRound) }, 650)
}
function runStage(seconds, done) {
  const start = performance.now(); stageRemaining.value = seconds
  stageTimer = window.setInterval(() => {
    stageRemaining.value = Math.max(0, seconds - (performance.now() - start) / 1000)
    if (stageRemaining.value <= 0) { window.clearInterval(stageTimer); done() }
  }, 50)
}
function toggleCell(index) {
  const next = new Set(selected.value)
  if (next.has(index)) next.delete(index); else next.add(index)
  selected.value = next
}
function submitRound() {
  if (phase.value !== 'recall') return
  window.clearInterval(stageTimer)
  const exact = selected.value.size === positions.value.length && positions.value.every((position) => selected.value.has(position))
  if (exact) successes.value += 1
  phase.value = 'transition'; stageRemaining.value = 0
  if (round.value >= 10) { transitionTimer = window.setTimeout(finish, 450); return }
  round.value += 1
  transitionTimer = window.setTimeout(loadRound, 450)
}
async function finish() {
  phase.value = 'result'
  try {
    await submitMemoryScore({ successNum: successes.value, trainingNum: 10, accuracy: successes.value / 10, level: difficulty.value })
    saveMessage.value = '成绩已保存'
  } catch (requestError) { saveMessage.value = apiMessage(requestError, '成绩暂未保存') }
}
function reset() { clearTimers(); phase.value = 'setup'; stageRemaining.value = 0 }
function clearTimers() { window.clearInterval(stageTimer); window.clearTimeout(transitionTimer) }
onBeforeUnmount(clearTimers)
</script>

<style scoped>
.difficulty-control { margin-bottom: 18px; }.start-button { display: flex; margin: 0 auto; }
.memory-sample { width: 180px; display: grid; grid-template-columns: repeat(4, 1fr); gap: 5px; margin: 0 auto 24px; }
.memory-sample i { aspect-ratio: 1; border: 1px solid var(--line); background: #f2f4f2; }.memory-sample i.on { background: var(--blue); border-color: var(--blue); }
.memory-play { display: grid; justify-items: center; }
.memory-grid { width: min(100%, 630px); aspect-ratio: 7 / 6; display: grid; gap: 5px; }
.memory-grid button { min-width: 0; border: 1px solid #cfd6d1; border-radius: 3px; background: #f6f8f6; transition: background 140ms ease, transform 140ms ease; }
.memory-grid button.target { border-color: var(--blue); background: var(--blue); transform: scale(.92); }
.memory-grid button.selected { border-color: var(--green); background: #cfe2d5; box-shadow: inset 0 0 0 2px var(--green); }
.memory-actions { width: min(100%, 630px); display: flex; align-items: center; justify-content: space-between; gap: 12px; margin-top: 20px; }.memory-actions span { color: var(--muted); font-size: 13px; }
.memory-transition { min-height: 360px; display: grid; place-items: center; align-content: center; gap: 20px; color: var(--muted); }
.memory-transition i { width: 54px; aspect-ratio: 1; border: 3px solid #dce2dd; border-top-color: var(--blue); border-radius: 50%; animation: spin 650ms linear infinite; }
@keyframes spin { to { transform: rotate(360deg); } }
@media (max-width: 600px) { .memory-grid { gap: 3px; } .memory-actions { align-items: stretch; flex-direction: column; } }
</style>
