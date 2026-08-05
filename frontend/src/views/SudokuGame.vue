<template>
  <section class="game-page sudoku-page">
    <div class="game-topline"><router-link class="back-link" to="/">返回训练</router-link><span class="eyebrow">逻辑推理 · 9×9</span></div>
    <div class="game-surface">
      <header class="game-header">
        <h1>数独挑战</h1>
        <div class="game-meta"><div><span>错误</span><strong>{{ mistakes }}</strong></div><div><span>用时</span><strong>{{ formattedTime }}</strong></div></div>
      </header>

      <div v-if="phase === 'setup'" class="game-body game-setup">
        <div class="setup-inner"><div class="sudoku-sample" aria-hidden="true"><i>5</i><i></i><i>4</i><i></i><i>7</i><i></i><i>6</i><i></i><i>2</i></div><h2>填满 9 × 9 方格</h2><p>每行、每列和每个九宫格使用 1 到 9，且不重复。</p><div class="segmented difficulty-control"><button v-for="item in difficulties" :key="item.value" :class="{ active: difficulty === item.value }" type="button" @click="difficulty = item.value">{{ item.label }}</button></div><button class="button start-button" type="button" @click="startGame">生成新题</button></div>
      </div>

      <div v-else-if="phase === 'playing'" class="game-body sudoku-play">
        <div class="sudoku-workspace">
          <div class="sudoku-board" role="grid" aria-label="数独棋盘">
            <button v-for="(value, index) in values" :key="index" type="button" role="gridcell" :class="cellClass(index, value)" :aria-label="`第 ${Math.floor(index / 9) + 1} 行第 ${index % 9 + 1} 列${value ? `，数字 ${value}` : ''}`" @click="selectCell(index)">{{ value || '' }}</button>
          </div>
          <div class="sudoku-controls">
            <p>{{ selectedIndex === null ? '选择一个空格' : `第 ${Math.floor(selectedIndex / 9) + 1} 行 · 第 ${selectedIndex % 9 + 1} 列` }}</p>
            <div class="number-pad"><button v-for="number in 9" :key="number" type="button" :disabled="selectedIndex === null || givens[selectedIndex]" @click="enterNumber(number)">{{ number }}</button></div>
            <div class="sudoku-actions"><button class="button button-secondary" type="button" title="清除当前格" :disabled="selectedIndex === null || givens[selectedIndex]" @click="enterNumber(0)">清除</button><button class="button button-secondary" type="button" @click="resetPuzzle">重置</button><button class="button" type="button" @click="startGame">新题</button></div>
          </div>
        </div>
      </div>

      <div v-else class="game-body game-result"><div class="result-inner"><p class="eyebrow">Puzzle complete</p><h2>数独完成</h2><div class="result-number">{{ formattedTime }}</div><p>{{ mistakes }} 次错误 · {{ saveMessage }}</p><div class="result-actions"><button class="button" type="button" @click="startGame">下一题</button><router-link class="button button-secondary" to="/">选择其他训练</router-link></div></div></div>
    </div>
  </section>
</template>

<script setup>
import { computed, onBeforeUnmount, ref } from 'vue'
import { getSudoku } from 'sudoku-gen'
import { apiMessage, submitSudokuScore } from '@/services/api'

const phase = ref('setup')
const difficulty = ref('easy')
const puzzle = ref('')
const solution = ref('')
const values = ref([])
const givens = ref([])
const wrongCells = ref(new Set())
const selectedIndex = ref(null)
const mistakes = ref(0)
const elapsed = ref(0)
const saveMessage = ref('正在保存成绩…')
const difficulties = [{ value: 'easy', label: '简单' }, { value: 'medium', label: '中等' }, { value: 'hard', label: '困难' }]
const formattedTime = computed(() => `${String(Math.floor(elapsed.value / 60)).padStart(2, '0')}:${String(elapsed.value % 60).padStart(2, '0')}`)
let timerId = null

function startGame() {
  window.clearInterval(timerId)
  const generated = getSudoku(difficulty.value)
  puzzle.value = generated.puzzle
  solution.value = generated.solution
  givens.value = [...puzzle.value].map((character) => character !== '-')
  values.value = [...puzzle.value].map((character) => character === '-' ? 0 : Number(character))
  wrongCells.value = new Set()
  selectedIndex.value = null
  mistakes.value = 0
  elapsed.value = 0
  saveMessage.value = '正在保存成绩…'
  phase.value = 'playing'
  timerId = window.setInterval(() => { elapsed.value += 1 }, 1000)
}
function selectCell(index) { selectedIndex.value = index }
function enterNumber(number) {
  const index = selectedIndex.value
  if (index === null || givens.value[index]) return
  const nextValues = [...values.value]
  nextValues[index] = number
  values.value = nextValues
  const nextWrong = new Set(wrongCells.value)
  if (number && String(number) !== solution.value[index]) {
    if (!nextWrong.has(index)) mistakes.value += 1
    nextWrong.add(index)
  } else nextWrong.delete(index)
  wrongCells.value = nextWrong
  if (values.value.every((value, cellIndex) => String(value) === solution.value[cellIndex])) finish()
}
function resetPuzzle() {
  values.value = [...puzzle.value].map((character) => character === '-' ? 0 : Number(character))
  wrongCells.value = new Set(); selectedIndex.value = null; mistakes.value = 0; elapsed.value = 0
}
function cellClass(index, value) {
  const selected = selectedIndex.value
  const sameGroup = selected !== null && (Math.floor(index / 9) === Math.floor(selected / 9) || index % 9 === selected % 9 || (Math.floor(index / 27) === Math.floor(selected / 27) && Math.floor((index % 9) / 3) === Math.floor((selected % 9) / 3)))
  return { given: givens.value[index], selected: selected === index, related: sameGroup, same: selected !== null && value && value === values.value[selected], wrong: wrongCells.value.has(index) }
}
async function finish() {
  window.clearInterval(timerId); phase.value = 'result'
  try {
    await submitSudokuScore({ level: difficulty.value, isPassed: true, mistakes: mistakes.value, timeElapsed: elapsed.value })
    saveMessage.value = '成绩已保存'
  } catch (requestError) { saveMessage.value = apiMessage(requestError, '成绩暂未保存') }
}
onBeforeUnmount(() => window.clearInterval(timerId))
</script>

<style scoped>
.difficulty-control { margin-bottom: 18px; }.start-button { display: flex; margin: 0 auto; }
.sudoku-sample { width: 156px; display: grid; grid-template-columns: repeat(3, 1fr); margin: 0 auto 24px; border: 2px solid var(--ink); }.sudoku-sample i { aspect-ratio: 1; display: grid; place-items: center; border: 1px solid var(--line); font-style: normal; font-weight: 800; }
.sudoku-play { display: grid; place-items: center; }.sudoku-workspace { width: min(100%, 780px); display: grid; grid-template-columns: minmax(360px, 520px) 200px; gap: 28px; align-items: center; }
.sudoku-board { width: 100%; aspect-ratio: 1; display: grid; grid-template-columns: repeat(9, 1fr); border: 2px solid var(--ink); background: var(--ink); gap: 1px; }
.sudoku-board button { min-width: 0; border: 0; background: white; color: var(--blue); font-size: 20px; font-weight: 750; }.sudoku-board button:nth-child(9n + 3), .sudoku-board button:nth-child(9n + 6) { border-right: 2px solid var(--ink); }.sudoku-board button:nth-child(n + 19):nth-child(-n + 27), .sudoku-board button:nth-child(n + 46):nth-child(-n + 54) { border-bottom: 2px solid var(--ink); }
.sudoku-board button.given { color: var(--ink); font-weight: 850; }.sudoku-board button.related { background: #edf2ee; }.sudoku-board button.same { background: #dce9e0; }.sudoku-board button.selected { background: #c6dfce; box-shadow: inset 0 0 0 2px var(--green); }.sudoku-board button.wrong { background: #f7ded9; color: var(--red); }
.sudoku-controls > p { min-height: 20px; margin: 0 0 16px; color: var(--muted); font-size: 13px; }.number-pad { display: grid; grid-template-columns: repeat(3, 1fr); gap: 5px; }.number-pad button { aspect-ratio: 1; border: 1px solid var(--line); border-radius: 4px; background: #f2f5f2; color: var(--ink); font-size: 20px; font-weight: 800; }.number-pad button:hover:not(:disabled) { border-color: var(--green); background: white; }.number-pad button:disabled { opacity: .45; }
.sudoku-actions { display: grid; gap: 7px; margin-top: 14px; }
@media (max-width: 760px) { .sudoku-workspace { grid-template-columns: 1fr; gap: 18px; }.sudoku-board button { font-size: 17px; }.sudoku-controls { display: grid; grid-template-columns: 1fr; }.number-pad { grid-template-columns: repeat(9, 1fr); gap: 2px; }.number-pad button { font-size: 15px; }.sudoku-actions { grid-template-columns: repeat(3, 1fr); }.sudoku-actions .button { padding: 0 8px; } }
</style>
