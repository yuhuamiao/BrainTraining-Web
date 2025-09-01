<template>
  <div class="bus-game-container">
    <!-- 难度选择器组件 -->
    <DifficultySelector 
      v-if="gameState === 'difficultySelect'"
      @difficulty-selected="startGame"
    />
    
    <!-- 游戏主界面 -->
    <div v-if="gameState === 'playing'">
      <!-- 公交车显示区域 -->
      <div class="bus-display">
        <!-- 公交车图标 -->
        <div class="bus-icon">
          <img src="@/assets/images/bus.png" alt="公交车" />
          <!-- 当前人数显示 -->
          <div class="passenger-count">{{ currentPassengers }}</div>
        </div>
        
        <!-- 乘客动画区域 -->
        <div class="passenger-animation-area">
          <!-- 使用TransitionGroup实现进出动画 -->
          <TransitionGroup name="passenger" tag="div">
            <div 
              v-for="passenger in visiblePassengers"
              :key="`${passenger.id}-${passenger.type}`"
              class="passenger"
              :class="[passenger.type, `speed-${currentDifficulty}`]"
            >
              <img :src="passengerImage(passenger.type)" alt="乘客" />
            </div>
          </TransitionGroup>
        </div>
      </div>
      
      <!-- 轮次信息 -->
      <div class="round-info">
        第 {{ currentRound }} / {{ totalRounds }} 轮
      </div>

      <!-- 当前轮次答案输入区域 -->
      <div v-if="gameState === 'playing' && showRoundInput" class="round-input">
        <h3>当前有多少人？</h3>
        <input 
          type="number" 
          v-model.number="roundAnswer" 
          min="0"
          placeholder="输入当前人数"
          @keyup.enter="submitRoundAnswer"
        />
        <button @click="submitRoundAnswer">提交</button>
      </div>
    </div>
    
    <!-- 最终答案输入区域 -->
    <div v-if="gameState === 'inputFinalAnswer'" class="answer-input">
      <h3>当前有 {{ currentPassengers }} 人，最终有多少人？</h3>
      <input 
        type="number" 
        v-model.number="userAnswer" 
        min="0"
        placeholder="输入最终人数"
        @keyup.enter="submitAnswer"
      />
      <button @click="submitAnswer">提交</button>
    </div>
    
    <!-- 结果反馈 -->
    <ResultDisplay 
      v-if="gameState === 'showResult'"
      :isCorrect="isAnswerCorrect"
      :correctAnswer="correctAnswer"
      :successNum="successNum"
      :trainingNum="trainingNum"
      @restart="resetGame"
    />
  </div>
</template>

<script>
import { ref, computed, onUnmounted } from 'vue';
import { useGameStore } from '@/stores/gameStore';
import DifficultySelector from '@/components/BusNumberGame/DifficultySelector.vue';
import ResultDisplay from '@/components/BusNumberGame/ResultDisplay.vue';

export default {
  components: {
    DifficultySelector,
    ResultDisplay
  },
  setup() {
    const gameStore = useGameStore();
    
    // 游戏状态：difficultySelect, playing, inputFinalAnswer, showResult
    const gameState = ref('difficultySelect');
    
    // 当前难度（从store获取）
    const currentDifficulty = computed(() => gameStore.difficulty);
    
    // 当前轮次（从store获取）
    const currentRound = computed(() => gameStore.currentRound);
    
    // 总轮次（从store获取）
    const totalRounds = computed(() => gameStore.totalRounds);
    
    // 当前乘客数量（从store获取）
    const currentPassengers = computed(() => gameStore.currentPassengers);
    
    // 用户输入的最终答案
    const userAnswer = ref(null);
    
    // 用户输入的每轮答案
    const roundAnswer = ref(null);
    
    // 是否显示每轮输入框
    const showRoundInput = ref(false);
    
    // 正确答案（计算属性）
    const correctAnswer = computed(() => gameStore.correctAnswer);
    
    // 每轮正确答案（计算属性）
    const correctRoundAnswer = computed(() => gameStore.correctRoundAnswer);
    
    // 答案是否正确
    const isAnswerCorrect = ref(false);
    
    // 成功次数和训练次数
    const successNum = ref(0);
    const trainingNum = ref(0);
    
    // 可见的乘客列表（用于动画）
    const visiblePassengers = ref([]);
    
    // 乘客ID计数器（确保唯一key）
    let passengerIdCounter = 0;
    
    // 存储所有定时器的引用，用于组件卸载时清理
    const timers = ref([]);
    
    // 根据乘客类型获取图片
    const passengerImage = (type) => {
      return type === 'entering' 
        ? require('@/assets/images/passenger_enter.png')
        : require('@/assets/images/passenger_exit.png');
    };
    
    // 设置定时器并跟踪它
    const setTrackedTimeout = (callback, delay) => {
      const timerId = setTimeout(() => {
        // 执行完成后从列表中移除
        const index = timers.value.indexOf(timerId);
        if (index > -1) {
          timers.value.splice(index, 1);
        }
        callback();
      }, delay);
      
      timers.value.push(timerId);
      return timerId;
    };
    
    // 开始游戏（选择难度后调用）
    const startGame = (difficulty) => {
      gameStore.startNewGame(difficulty);
      gameState.value = 'playing';
      successNum.value = 0;
      trainingNum.value = 0;
      nextRound();
    };
    
    // 进行下一轮
    const nextRound = () => {
      if (currentRound.value >= totalRounds.value) {
        // 所有轮次结束，进入最终答题环节
        gameState.value = 'inputFinalAnswer';
        return;
      }
      
      // 生成本轮变化数据
      const roundData = gameStore.generateRound();
      
      // 清空可见乘客
      visiblePassengers.value = [];
      
      // 根据变化数量生成乘客
      for (let i = 0; i < roundData.change; i++) {
        // 为每个乘客创建唯一ID
        passengerIdCounter++;
        
        // 添加乘客到可见列表（带延迟实现逐个出现效果）
        setTrackedTimeout(() => {
          visiblePassengers.value.push({
            id: passengerIdCounter,
            type: roundData.isEntering ? 'entering' : 'exiting'
          });
        }, i * 300); // 每个乘客延迟300ms出现
      }
      
      // 乘客停留时间（根据难度不同）
      const stayDuration = {
        easy: 3000,
        medium: 2000,
        hard: 1500
      }[currentDifficulty.value];
      
      // 停留后显示输入框让用户输入当前人数
      setTrackedTimeout(() => {
        // 清空可见乘客
        visiblePassengers.value = [];
        
        // 显示输入框
        showRoundInput.value = true;
        roundAnswer.value = null;
        
        // 计算正确答案
        gameStore.calculateRoundAnswer(
          roundData.isEntering ? roundData.change : -roundData.change
        );
      }, stayDuration);
    };
    
    // 提交每轮答案
    const submitRoundAnswer = () => {
      if (roundAnswer.value === null) return;
      
      // 更新训练次数
      trainingNum.value++;
      
      // 检查答案是否正确
      if (roundAnswer.value === correctRoundAnswer.value) {
        successNum.value++;
      }
      
      // 更新乘客数量
      gameStore.updatePassengers(
        correctRoundAnswer.value - currentPassengers.value
      );
      
      // 隐藏输入框
      showRoundInput.value = false;
      
      // 进入下一轮
      gameStore.nextRound();
      setTrackedTimeout(nextRound, 800); // 轮次间延迟
    };
    
    // 提交最终答案
    const submitAnswer = async () => {
      if (userAnswer.value === null) return;
      
      isAnswerCorrect.value = (userAnswer.value === correctAnswer.value);
      
      // 最终答案也算一次训练
      trainingNum.value++;
      if (isAnswerCorrect.value) {
        successNum.value++;
      }
      
      try {
        // 从localStorage获取用户 token（这个我还没有搞明白）
        const token = localStorage.getItem('token');
        if (!token) {
          console.warn('用户未登录，成绩将不会被保存');
          gameState.value = 'showResult';
          return;
        }
        
        // 获取用户 ID（假设存储在 localStorage 或 Vuex/Pinia store 中）
        const user = JSON.parse(localStorage.getItem('user') || '{}');
        const userId = user.id;
        
        if (!userId) {
          console.warn('无法获取用户ID，成绩将不会被保存');
          gameState.value = 'showResult';
          return;
        }
        
        // 计算准确率
        const accuracy = trainingNum.value > 0 ? successNum.value / trainingNum.value : 0;
        
        // 发送成绩到后端
        const response = await fetch('/api/v1/bus/scores', {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json',
            'Authorization': `Bearer ${token}`
          },
          body: JSON.stringify({
            UserID: userId,
            SuccessNum: successNum.value,
            Level: currentDifficulty.value,
            Accuracy: accuracy,
            TrainingNum: trainingNum.value
          })
        });
        
        if (!response.ok) {
          throw new Error(`HTTP error! status: ${response.status}`);
        }
        
        const result = await response.json();
        console.log('成绩提交成功:', result);
      } catch (error) {
        console.error('提交成绩错误:', error);
        // 即使提交失败，也继续显示结果
      }
      
      gameState.value = 'showResult';
    };
    
    // 重置游戏
    const resetGame = () => {
      // 清理所有定时器
      timers.value.forEach(timerId => clearTimeout(timerId));
      timers.value = [];
      
      gameState.value = 'difficultySelect';
      userAnswer.value = null;
      roundAnswer.value = null;
      showRoundInput.value = false;
      visiblePassengers.value = [];
    };
    
    // 组件卸载时清理所有定时器
    onUnmounted(() => {
      timers.value.forEach(timerId => clearTimeout(timerId));
    });
    
    return {
      gameState,
      currentDifficulty,
      currentRound,
      totalRounds,
      currentPassengers,
      userAnswer,
      roundAnswer,
      showRoundInput,
      correctAnswer,
      isAnswerCorrect,
      successNum,
      trainingNum,
      visiblePassengers,
      passengerImage,
      startGame,
      submitRoundAnswer,
      submitAnswer,
      resetGame
    };
  }
};
</script>

<style scoped>
.bus-game-container {
  max-width: 800px;
  margin: 0 auto;
  padding: 20px;
}

.bus-display {
  display: flex;
  align-items: center;
  justify-content: center;
  margin-bottom: 30px;
}

.bus-icon {
  position: relative;
  margin-right: 50px;
}

.bus-icon img {
  width: 200px;
}

.passenger-count {
  position: absolute;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  font-size: 2rem;
  font-weight: bold;
  color: #333;
}

.passenger-animation-area {
  position: relative;
  width: 300px;
  height: 150px;
  border: 1px dashed #ccc;
}

.passenger {
  position: absolute;
  transition: all 0.5s ease;
}

.passenger img {
  width: 40px;
  height: 40px;
}

/* 进入动画 */
.passenger-enter-active.entering {
  animation: enter 1s forwards;
}

/* 离开动画 */
.passenger-enter-active.exiting {
  animation: exit 1s forwards;
}

/* 难度相关的动画速度 */
.speed-easy { animation-duration: 1.5s; }
.speed-medium { animation-duration: 1s; }
.speed-hard { animation-duration: 0.6s; }

@keyframes enter {
  0% { 
    transform: translateX(100px); 
    opacity: 0; 
  }
  100% { 
    transform: translateX(0); 
    opacity: 1; 
  }
}

@keyframes exit {
  0% { 
    transform: translateX(0); 
    opacity: 1; 
  }
  100% { 
    transform: translateX(-100px); 
    opacity: 0; 
  }
}

.answer-input, .round-input {
  text-align: center;
  margin-top: 30px;
}

.answer-input input, .round-input input {
  font-size: 1.2rem;
  padding: 10px;
  margin: 0 10px;
  width: 100px;
}

.answer-input button, .round-input button {
  font-size: 1.2rem;
  padding: 10px 20px;
  background-color: #4CAF50;
  color: white;
  border: none;
  border-radius: 4px;
  cursor: pointer;
}

.round-info {
  text-align: center;
  font-size: 1.2rem;
  margin-top: 20px;
}
</style>