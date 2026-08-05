<template>
  <section>
    <div class="page-heading">
      <div>
        <p class="eyebrow">Daily practice</p>
        <h1>今天，练哪一项？</h1>
        <p>五种练习，覆盖速度、抑制控制、空间记忆、心算与逻辑。</p>
      </div>
      <router-link class="button button-secondary" to="/progress">查看训练记录</router-link>
    </div>

    <div class="snapshot-band">
      <div><span>累计训练</span><strong>{{ totalSessions }}</strong><small>次</small></div>
      <div><span>已完成项目</span><strong>{{ completedTypes }}</strong><small>/ 5</small></div>
      <div><span>本次建议</span><strong class="recommendation">{{ recommendedTitle }}</strong></div>
    </div>

    <div class="training-grid">
      <router-link
        v-for="training in trainings"
        :key="training.route"
        :to="training.route"
        class="training-card"
        :style="{ '--accent': training.color }"
      >
        <div class="training-copy">
          <span class="training-index">{{ training.index }}</span>
          <p>{{ training.focus }}</p>
          <h2>{{ training.title }}</h2>
          <dl>
            <div><dt>时长</dt><dd>{{ training.duration }}</dd></div>
            <div><dt>节奏</dt><dd>{{ training.pace }}</dd></div>
          </dl>
        </div>

        <div class="training-preview" :class="training.preview" aria-hidden="true">
          <template v-if="training.preview === 'preview-schulte'">
            <i v-for="number in [17, 3, 21, 8, 12, 1, 25, 6, 14]" :key="number">{{ number }}</i>
          </template>
          <template v-else-if="training.preview === 'preview-color'">
            <b>红</b><b>蓝</b><b>绿</b>
          </template>
          <template v-else-if="training.preview === 'preview-memory'">
            <i v-for="cell in 20" :key="cell" :class="{ on: [2, 7, 9, 14, 18].includes(cell) }"></i>
          </template>
          <img v-else-if="training.preview === 'preview-bus'" :src="busImage" alt="" />
          <template v-else>
            <i v-for="(number, index) in [5, 3, 0, 0, 7, 0, 6, 0, 2]" :key="index">{{ number || '' }}</i>
          </template>
        </div>
        <span class="enter-label">开始</span>
      </router-link>
    </div>
  </section>
</template>

<script setup>
import { computed, onMounted, ref } from 'vue'
import busImage from '@/assets/images/bus.png'
import { getUserScores } from '@/services/api'

const scores = ref({})
const trainings = [
  { index: '01', title: '舒尔特方格', focus: '视觉搜索', duration: '30 秒', pace: '快速', route: '/schulte', color: '#2e6b50', preview: 'preview-schulte' },
  { index: '02', title: '多色文字', focus: '抑制控制', duration: '30 秒', pace: '连续', route: '/color-words', color: '#dd6b52', preview: 'preview-color' },
  { index: '03', title: '瞬时记忆', focus: '空间记忆', duration: '10 轮', pace: '专注', route: '/memory', color: '#3d6f9f', preview: 'preview-memory' },
  { index: '04', title: '公交人数', focus: '工作记忆', duration: '3–7 站', pace: '变化', route: '/bus-game', color: '#e3b33e', preview: 'preview-bus' },
  { index: '05', title: '数独挑战', focus: '逻辑推理', duration: '不限时', pace: '沉浸', route: '/sudoku', color: '#725c8e', preview: 'preview-sudoku' },
]

const scoreKeys = ['schulte', 'colorWord', 'memory', 'bus', 'sudoku']
const totalSessions = computed(() => scoreKeys.reduce((sum, key) => sum + (scores.value[key]?.length || 0), 0))
const completedTypes = computed(() => scoreKeys.filter((key) => scores.value[key]?.length).length)
const recommendedTitle = computed(() => trainings.find((_, index) => !scores.value[scoreKeys[index]]?.length)?.title || trainings[totalSessions.value % trainings.length].title)

onMounted(async () => {
  try {
    scores.value = await getUserScores()
  } catch {
    scores.value = {}
  }
})
</script>

<style scoped>
.snapshot-band { display: grid; grid-template-columns: 1fr 1fr 1.4fr; margin-bottom: 30px; border-block: 1px solid var(--line); }
.snapshot-band > div { min-height: 94px; display: flex; align-items: baseline; gap: 6px; padding: 24px 26px; border-right: 1px solid var(--line); }
.snapshot-band > div:last-child { border-right: 0; }
.snapshot-band span { margin-right: auto; color: var(--muted); font-size: 13px; }
.snapshot-band strong { font-size: 29px; }
.snapshot-band small { color: var(--muted); }
.snapshot-band .recommendation { font-size: 18px; color: var(--green); }
.training-grid { display: grid; grid-template-columns: repeat(2, minmax(0, 1fr)); gap: 16px; }
.training-card { min-height: 242px; position: relative; display: grid; grid-template-columns: 1.1fr 0.9fr; overflow: hidden; border: 1px solid var(--line); border-top: 4px solid var(--accent); border-radius: 7px; background: white; color: var(--ink); text-decoration: none; transition: transform 180ms ease, box-shadow 180ms ease; }
.training-card:hover { transform: translateY(-3px); box-shadow: var(--shadow); }
.training-card:last-child { grid-column: 1 / -1; grid-template-columns: 1fr 1fr; }
.training-copy { padding: 24px; }
.training-index { color: var(--accent); font-size: 12px; font-weight: 850; }
.training-copy p { margin: 28px 0 5px; color: var(--muted); font-size: 12px; font-weight: 700; }
.training-copy h2 { margin: 0; font-size: 22px; }
.training-copy dl { display: flex; gap: 24px; margin: 24px 0 0; }
.training-copy dl div { display: grid; gap: 3px; }
.training-copy dt { color: var(--muted); font-size: 11px; }
.training-copy dd { margin: 0; font-size: 13px; font-weight: 750; }
.enter-label { position: absolute; right: 16px; bottom: 14px; color: var(--accent); font-size: 13px; font-weight: 800; }
.training-preview { min-width: 0; display: grid; place-content: center; padding: 22px; background: #f0f3f0; }
.preview-schulte, .preview-sudoku { grid-template-columns: repeat(3, 38px); grid-template-rows: repeat(3, 38px); gap: 3px; }
.preview-schulte i, .preview-sudoku i { display: grid; place-items: center; border: 1px solid #c8d0ca; background: white; font-size: 12px; font-style: normal; font-weight: 750; }
.preview-color { grid-template-columns: 1fr; gap: 4px; font-size: 30px; line-height: 1; }
.preview-color b:nth-child(1) { color: #3d6f9f; }
.preview-color b:nth-child(2) { color: #dd6b52; }
.preview-color b:nth-child(3) { color: #2e6b50; }
.preview-memory { grid-template-columns: repeat(5, 19px); gap: 5px; }
.preview-memory i { width: 19px; aspect-ratio: 1; border: 1px solid #c8d0ca; background: white; }
.preview-memory i.on { border-color: #3d6f9f; background: #3d6f9f; }
.preview-bus img { width: min(100%, 210px); height: auto; }
@media (max-width: 760px) {
  .snapshot-band { grid-template-columns: 1fr 1fr; }
  .snapshot-band > div { min-height: 76px; padding: 16px 12px; }
  .snapshot-band > div:nth-child(2) { border-right: 0; }
  .snapshot-band > div:last-child { grid-column: 1 / -1; border-top: 1px solid var(--line); }
  .training-grid { grid-template-columns: 1fr; }
  .training-card, .training-card:last-child { grid-column: auto; grid-template-columns: 1.1fr 0.9fr; min-height: 220px; }
  .training-copy { padding: 19px 16px; }
  .training-copy dl { gap: 14px; }
  .training-preview { padding: 10px; }
  .preview-schulte, .preview-sudoku { grid-template-columns: repeat(3, 31px); grid-template-rows: repeat(3, 31px); }
}
</style>
