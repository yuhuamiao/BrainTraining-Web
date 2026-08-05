<template>
  <section class="game-page">
    <div class="game-topline"><router-link class="back-link" to="/">返回训练</router-link><span class="eyebrow">工作记忆 · 3–7 站</span></div>
    <div class="game-surface">
      <header class="game-header">
        <h1>公交人数</h1>
        <div class="game-meta">
          <div><span>站点</span><strong>{{ phase === 'playing' ? `${eventIndex + 1}/${session.events.length}` : '—' }}</strong></div>
          <div><span>难度</span><strong>{{ difficultyLabel }}</strong></div>
        </div>
      </header>

      <div v-if="phase === 'setup'" class="game-body game-setup">
        <div class="setup-inner">
          <img class="bus-setup-image" :src="busImage" alt="公交车" />
          <h2>记住车上人数</h2><p>从起始人数开始，连续计算每站上下车后的最终人数。</p>
          <div class="segmented difficulty-control"><button v-for="item in difficulties" :key="item.value" :class="{ active: difficulty === item.value }" type="button" @click="difficulty = item.value">{{ item.label }}</button></div>
          <button class="button start-button" type="button" @click="startGame">开始发车</button>
        </div>
      </div>

      <div v-else-if="phase === 'playing'" class="game-body bus-play">
        <div class="bus-scene">
          <div class="route-line"><i v-for="(_, index) in session.events" :key="index" :class="{ passed: index < eventIndex, current: index === eventIndex }"></i></div>
          <div v-if="showStart" class="event-display start-count"><span>起始人数</span><strong>{{ session.startCount }}</strong><small>人</small></div>
          <div v-else class="event-display" :class="currentEvent.entering ? 'entering' : 'exiting'">
            <span>{{ currentEvent.entering ? '上车' : '下车' }}</span><strong>{{ currentEvent.entering ? '+' : '−' }}{{ currentEvent.amount }}</strong><small>人</small>
          </div>
          <div class="bus-visual">
            <img class="bus-image" :src="busImage" alt="公交车" />
            <div v-if="!showStart" :class="['passenger-strip', currentEvent.entering ? 'move-in' : 'move-out']">
              <img v-for="index in currentEvent.amount" :key="index" :src="currentEvent.entering ? passengerEnter : passengerExit" alt="" />
            </div>
          </div>
        </div>
      </div>

      <div v-else-if="phase === 'answer'" class="game-body bus-answer">
        <div>
          <p class="eyebrow">终点站</p><h2>现在车上有多少人？</h2>
          <form @submit.prevent="submitAnswer"><input v-model.number="answer" class="input answer-input" type="number" min="0" inputmode="numeric" autofocus required /><button class="button" type="submit">提交答案</button></form>
        </div>
      </div>

      <div v-else class="game-body game-result">
        <div class="result-inner"><p class="eyebrow">{{ isCorrect ? 'Correct' : 'Keep going' }}</p><h2>{{ isCorrect ? '计算正确' : '这次差一点' }}</h2><div class="result-number">{{ session.finalCount }} 人</div><p>你的答案：{{ answer }} · {{ saveMessage }}</p><div class="result-actions"><button class="button" type="button" @click="reset">再练一次</button><router-link class="button button-secondary" to="/">选择其他训练</router-link></div></div>
      </div>
    </div>
  </section>
</template>

<script setup>
import { computed, onBeforeUnmount, ref } from 'vue'
import busImage from '@/assets/images/bus.png'
import passengerEnter from '@/assets/images/passenger_enter.png'
import passengerExit from '@/assets/images/passenger_exit.png'
import { apiMessage, submitBusScore } from '@/services/api'
import { generateBusSession } from '@/utils/games/bus'

const phase = ref('setup')
const difficulty = ref('easy')
const session = ref({ startCount: 0, events: [], finalCount: 0, delay: 1000 })
const eventIndex = ref(0)
const showStart = ref(true)
const answer = ref(null)
const isCorrect = ref(false)
const saveMessage = ref('正在保存成绩…')
const difficulties = [{ value: 'easy', label: '简单' }, { value: 'medium', label: '中等' }, { value: 'hard', label: '困难' }]
const difficultyLabel = computed(() => difficulties.find((item) => item.value === difficulty.value)?.label)
const currentEvent = computed(() => session.value.events[eventIndex.value] || { entering: true, amount: 0 })
let eventTimer = null

function startGame() {
  window.clearTimeout(eventTimer)
  session.value = generateBusSession(difficulty.value)
  eventIndex.value = 0
  showStart.value = true
  answer.value = null
  saveMessage.value = '正在保存成绩…'
  phase.value = 'playing'
  eventTimer = window.setTimeout(showEvent, 1500)
}
function showEvent() {
  showStart.value = false
  eventTimer = window.setTimeout(() => {
    if (eventIndex.value >= session.value.events.length - 1) { phase.value = 'answer'; return }
    eventIndex.value += 1
    showEvent()
  }, session.value.delay)
}
async function submitAnswer() {
  if (answer.value === null || answer.value < 0) return
  isCorrect.value = Number(answer.value) === session.value.finalCount
  phase.value = 'result'
  try {
    await submitBusScore({ successNum: isCorrect.value ? 1 : 0, trainingNum: 1, accuracy: isCorrect.value ? 1 : 0, level: difficulty.value })
    saveMessage.value = '成绩已保存'
  } catch (requestError) { saveMessage.value = apiMessage(requestError, '成绩暂未保存') }
}
function reset() { window.clearTimeout(eventTimer); phase.value = 'setup' }
onBeforeUnmount(() => window.clearTimeout(eventTimer))
</script>

<style scoped>
.difficulty-control { margin-bottom: 18px; }.start-button { display: flex; margin: 0 auto; }
.bus-setup-image { width: 230px; max-width: 80%; margin-bottom: 12px; }
.bus-play, .bus-answer { min-height: 430px; display: grid; place-items: center; text-align: center; }
.bus-scene { width: min(100%, 700px); }
.route-line { display: flex; justify-content: center; gap: 20px; margin-bottom: 28px; }
.route-line i { width: 13px; aspect-ratio: 1; border: 2px solid #bbc5be; border-radius: 50%; background: white; }
.route-line i.passed { border-color: var(--green); background: var(--green); }.route-line i.current { border-color: var(--yellow); box-shadow: 0 0 0 4px rgba(227,179,62,.2); }
.event-display { min-height: 92px; display: flex; align-items: baseline; justify-content: center; gap: 8px; }
.event-display span { margin-right: 8px; color: var(--muted); font-size: 14px; font-weight: 750; }.event-display strong { color: var(--green); font-size: 52px; }.event-display.exiting strong { color: var(--coral); }.event-display small { color: var(--muted); }
.bus-visual { min-height: 220px; position: relative; display: grid; place-items: center; }
.bus-image { width: min(82%, 420px); height: auto; }.passenger-strip { position: absolute; right: 1%; bottom: 12px; max-width: 45%; display: flex; justify-content: center; flex-wrap: wrap; }
.passenger-strip img { width: 34px; height: 50px; object-fit: contain; animation: passengerMove 650ms ease both; }.move-out img { animation-direction: reverse; }
@keyframes passengerMove { from { transform: translateX(60px); opacity: 0; } to { transform: translateX(0); opacity: 1; } }
.bus-answer > div { width: min(100%, 460px); }.bus-answer h2 { margin: 0 0 24px; font-size: 28px; }.bus-answer form { display: grid; grid-template-columns: 1fr auto; gap: 10px; }.answer-input { min-height: 54px; text-align: center; font-size: 24px; font-weight: 800; }
@media (max-width: 600px) { .bus-visual { min-height: 180px; } .event-display strong { font-size: 42px; } .bus-answer form { grid-template-columns: 1fr; } .passenger-strip img { width: 26px; height: 42px; } }
</style>
