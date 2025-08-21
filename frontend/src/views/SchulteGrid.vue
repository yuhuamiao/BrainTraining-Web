<template>
    <div class="schulte-container">
      <h2>舒尔特表格训练</h2>
      
      <!-- 游戏说明 -->
      <div v-if="!gameStarted" class="instructions">
        <p>请在30秒内按顺序点击1到25的数字</p>
        <button @click="startGame">开始游戏</button>
      </div>
      
      <!-- 游戏界面 -->
      <div v-else class="game-area">
        <!-- 计时器 -->
        <div class="timer">剩余时间: {{ timeLeft }}秒</div>
        
        <!-- 网格 -->
        <div class="grid" :style="gridStyle">
          <div 
            v-for="(number, index) in flattenedMatrix" 
            :key="index"
            class="grid-cell"
            :class="{ 
              'selected': selectedCells.includes(number),
              'correct': correctOrder.includes(number)
            }"
            @click="handleCellClick(number)"
          >
            {{ number }}
          </div>
        </div>
        
        <!-- 游戏结果 -->
        <div v-if="gameEnded" class="result">
          <h3 v-if="isCompleted">恭喜！您在{{ timeUsed }}秒内完成了挑战</h3>
          <h3 v-else>时间到！您完成了{{ selectedCells.length }}个数字</h3>
          <button @click="restartGame">再试一次</button>
        </div>
      </div>
    </div>
  </template>
  
  <script>
  import { ref, computed, onMounted, onUnmounted } from 'vue'
  import axios from 'axios'
  
  export default {
    name: 'SchulteGrid',
    setup() {
      // 响应式数据
      const matrix = ref([])
      const gameStarted = ref(false)
      const gameEnded = ref(false)
      const timeLeft = ref(30)
      const selectedCells = ref([])
      const correctOrder = ref([])
      const timer = ref(null)
      const startTime = ref(null)
      const timeUsed = ref(0)
      
      // 计算属性
      const flattenedMatrix = computed(() => {
        return matrix.value.flat()
      })
      
      const gridStyle = computed(() => {
        const size = matrix.value.length
        return {
          gridTemplateColumns: `repeat(${size}, 1fr)`,
          gridTemplateRows: `repeat(${size}, 1fr)`
        }
      })
      
      const isCompleted = computed(() => {
        return selectedCells.value.length === 25
      })
      
      // 方法
      const fetchMatrix = async () => {
        try {
          const response = await axios.get('/api/v1/schulte/matrix?size=5')
          matrix.value = response.data.matrix
        } catch (error) {
          console.error('获取矩阵失败:', error)
          // 可以在这里添加错误处理逻辑
        }
      }
      
      const startGame = () => {
        gameStarted.value = true
        gameEnded.value = false
        timeLeft.value = 30
        selectedCells.value = []
        correctOrder.value = Array.from({length: 25}, (_, i) => i + 1)
        startTime.value = Date.now()
        
        // 启动计时器
        timer.value = setInterval(() => {
          timeLeft.value--
          
          if (timeLeft.value <= 0) {
            endGame()
          }
        }, 1000)
      }
      
      const handleCellClick = (number) => {
        if (gameEnded.value) return
        
        const expectedNumber = selectedCells.value.length + 1
        
        if (number === expectedNumber) {
          selectedCells.value.push(number)
          
          // 检查是否完成所有数字
          if (selectedCells.value.length === 25) {
            endGame()
          }
        }
      }
      
      const endGame = () => {
        clearInterval(timer.value)
        gameEnded.value = true
        timeUsed.value = ((Date.now() - startTime.value) / 1000).toFixed(2)
        
        // 提交成绩
        submitScore()
      }
      
      const submitScore = async () => {
        try {
          const token = localStorage.getItem('token')
          const userId = localStorage.getItem('userId') // 假设用户ID已存储
          
          await axios.post('/api/v1/schulte/scores', {
            isPassed: isCompleted.value,
            userId: userId,
            trainingNum: 1
          }, {
            headers: {
              'Authorization': `Bearer ${token}`
            }
          })
        } catch (error) {
          console.error('提交成绩失败:', error)
        }
      }
      
      const restartGame = () => {
        fetchMatrix()
        startGame()
      }
      
      // 生命周期钩子
      onMounted(() => {
        fetchMatrix()
      })
      
      onUnmounted(() => {
        if (timer.value) {
          clearInterval(timer.value)
        }
      })
      
      return {
        matrix,
        gameStarted,
        gameEnded,
        timeLeft,
        selectedCells,
        correctOrder,
        timeUsed,
        flattenedMatrix,
        gridStyle,
        isCompleted,
        startGame,
        handleCellClick,
        restartGame
      }
    }
  }
  </script>
  
  <style scoped>
  .schulte-container {
    display: flex;
    flex-direction: column;
    align-items: center;
    padding: 20px;
  }
  
  .instructions {
    text-align: center;
    margin-bottom: 20px;
  }
  
  .game-area {
    display: flex;
    flex-direction: column;
    align-items: center;
  }
  
  .timer {
    font-size: 24px;
    margin-bottom: 20px;
    font-weight: bold;
  }
  
  .grid {
    display: grid;
    gap: 5px;
    margin-bottom: 20px;
  }
  
  .grid-cell {
    width: 50px;
    height: 50px;
    display: flex;
    justify-content: center;
    align-items: center;
    border: 1px solid #ccc;
    border-radius: 5px;
    font-size: 18px;
    font-weight: bold;
    cursor: pointer;
    transition: all 0.2s;
  }
  
  .grid-cell:hover {
    background-color: #f0f0f0;
  }
  
  .grid-cell.selected {
    background-color: #4caf50;
    color: white;
  }
  
  .grid-cell.correct {
    background-color: #e8f5e9;
  }
  
  .result {
    text-align: center;
  }
  
  button {
    padding: 10px 20px;
    background-color: #4caf50;
    color: white;
    border: none;
    border-radius: 5px;
    cursor: pointer;
    font-size: 16px;
    margin-top: 10px;
  }
  
  button:hover {
    background-color: #45a049;
  }
  </style>