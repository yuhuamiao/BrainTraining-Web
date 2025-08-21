<template>
  <div class="difficulty-selector">
    <h2>选择难度</h2>
    <div class="difficulty-options">
      <button 
        v-for="level in difficultyLevels" 
        :key="level.id"
        @click="selectDifficulty(level.id)"
        :class="['difficulty-btn', level.id]"
      >
        <div class="difficulty-icon">
          <img :src="level.icon" :alt="level.label" />
        </div>
        <div class="difficulty-label">{{ level.label }}</div>
        <div class="difficulty-desc">{{ level.description }}</div>
      </button>
    </div>
  </div>
</template>

<script>
export default {
  data() {
    return {
      difficultyLevels: [
        {
          id: 'easy',
          label: '简单',
          description: '3轮变化，人数变化小',
          icon: require('@/assets/images/difficulty_easy.png')
        },
        {
          id: 'medium',
          label: '中等',
          description: '5轮变化，人数变化适中',
          icon: require('@/assets/images/difficulty_medium.png')
        },
        {
          id: 'hard',
          label: '困难',
          description: '7轮变化，人数变化大且速度快',
          icon: require('@/assets/images/difficulty_hard.png')
        }
      ]
    };
  },
  methods: {
    selectDifficulty(level) {
      this.$emit('difficulty-selected', level);
    }
  }
};
</script>

<style scoped>
.difficulty-selector {
  text-align: center;
}

.difficulty-options {
  display: flex;
  justify-content: center;
  gap: 20px;
  margin-top: 30px;
}

.difficulty-btn {
  width: 150px;
  padding: 20px 10px;
  border: none;
  border-radius: 10px;
  background: #f5f5f5;
  cursor: pointer;
  transition: all 0.3s;
}

.difficulty-btn:hover {
  transform: translateY(-5px);
  box-shadow: 0 5px 15px rgba(0,0,0,0.1);
}

.difficulty-btn.easy:hover { background: #e8f5e9; }
.difficulty-btn.medium:hover { background: #fff8e1; }
.difficulty-btn.hard:hover { background: #ffebee; }

.difficulty-icon img {
  width: 60px;
  height: 60px;
}

.difficulty-label {
  font-size: 1.2rem;
  font-weight: bold;
  margin: 10px 0;
}

.difficulty-desc {
  font-size: 0.9rem;
  color: #666;
}
</style>