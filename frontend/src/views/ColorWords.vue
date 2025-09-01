<template>
    <div class="color-words-training">
      <div class="header">
        <h2>多色文字训练</h2>
        <div class="stats">
          <div class="timer">时间: {{ formatTime(remainingTime) }}</div>
          <div class="score">得分: {{ score }}</div>
          <div class="accuracy">准确率: {{ accuracy }}%</div>
        </div>
      </div>
  
      <div v-if="!gameStarted" class="setup-screen">
        <div class="difficulty-selector">
          <h3>选择难度</h3>
          <div class="difficulty-options">
            <button 
              v-for="level in difficultyLevels" 
              :key="level.value" 
              @click="selectDifficulty(level.value)"
              :class="{ active: selectedDifficulty === level.value }"
            >
              {{ level.label }}
            </button>
          </div>
        </div>
        <button class="start-button" @click="startGame">开始训练</button>
      </div>
  
      <div v-else class="game-screen">
        <div class="color-matrix">
          <div 
            v-for="(row, rowIndex) in matrix" 
            :key="rowIndex" 
            class="matrix-row"
          >
            <div
              v-for="(cell, colIndex) in row"
              :key="colIndex"
              class="color-cell"
              :class="{ highlighted: isHighlighted(rowIndex, colIndex) }"
              :style="getCellStyle(rowIndex, colIndex)"
            >
              {{ cell.word }}
            </div>
          </div>
        </div>
  
        <div class="color-buttons">
          <button
            v-for="color in colorOptions"
            :key="color.value"
            @click="checkAnswer(color.value)"
            :disabled="!isAnswering || gameEnded"
            class="color-button"
            :style="{ backgroundColor: getButtonColor(color.value) }"
          >
            {{ color.label }}
          </button>
        </div>
  
        <div v-if="gameEnded" class="results">
          <h3>训练结束!</h3>
          <p>最终得分: {{ score }}</p>
          <p>准确率: {{ accuracy }}%</p>
          <button @click="restartGame">再试一次</button>
          <button @click="submitScore">保存成绩</button>
        </div>
      </div>
  
      <div v-if="loading" class="loading-overlay">
        <div class="spinner"></div>
        <p>加载中...</p>
      </div>
    </div>
  </template>
  
  <script>
  import { ref, computed, onMounted, onUnmounted } from 'vue'
  import axios from 'axios'
  
  export default {
    name: 'ColorWords',
    setup() {
      // 游戏状态
      const gameStarted = ref(false)
      const gameEnded = ref(false)
      const isAnswering = ref(false)
      const loading = ref(false)
      
      // 游戏数据
      const matrix = ref([])
      const selectedDifficulty = ref('easy')
      const currentColor = ref('')
      const highlightedCell = ref({ row: -1, col: -1 })
      
      // 统计数据
      const score = ref(0)
      const totalQuestions = ref(0)
      const correctAnswers = ref(0)
      const remainingTime = ref(0)
      const timerInterval = ref(null)
      const gameDuration = 60 // 60秒游戏时间
      
      // WebSocket连接
      const ws = ref(null)
      const isConnected = ref(false)
      
      // 难度选项
      const difficultyLevels = [
        { value: 'easy', label: '简单' },
        { value: 'medium', label: '中等' },
        { value: 'hard', label: '困难' }
      ]
      
      // 颜色选项
      const colorOptions = [
        { value: 'red', label: '红' },
        { value: 'blue', label: '蓝' },
        { value: 'green', label: '绿' },
        { value: 'yellow', label: '黄' },
        { value: 'black', label: '黑' }
      ]
      
      // 计算准确率
      const accuracy = computed(() => {
        if (totalQuestions.value === 0) return 0
        return Math.round((correctAnswers.value / totalQuestions.value) * 100)
      })
      
      // 获取矩阵数据
      const fetchMatrix = async () => {
        try {
          loading.value = true
          const response = await axios.get('/api/v1/color_words/matrix', {
            params: { difficulty: selectedDifficulty.value }
          })
          matrix.value = response.data.matrix
        } catch (error) {
          console.error('获取矩阵失败:', error)
          alert('获取训练数据失败，请重试')
        } finally {
          loading.value = false
        }
      }
      
      // 建立WebSocket连接
      const connectWebSocket = () => {
        const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:'
        const wsUrl = `${protocol}//${window.location.host}/api/v1/color_words/color`
        
        ws.value = new WebSocket(wsUrl)
        
        ws.value.onopen = () => {
          console.log('WebSocket连接已建立')
          isConnected.value = true
        }
        
        ws.value.onmessage = (event) => {
          const data = JSON.parse(event.data)
          currentColor.value = data.color
          highlightRandomCell()
          isAnswering.value = true
          
          // 5秒后自动隐藏高亮，如果用户没有回答则记为错误
          setTimeout(() => {
            if (isAnswering.value) {
              isAnswering.value = false
              totalQuestions.value++
            }
          }, 5000)
        }
        
        ws.value.onclose = () => {
          console.log('WebSocket连接已关闭')
          isConnected.value = false
        }
        
        ws.value.onerror = (error) => {
          console.error('WebSocket错误:', error)
          isConnected.value = false
        }
      }
      
      // 随机高亮一个单元格
      const highlightRandomCell = () => {
        const rows = matrix.value.length
        const cols = matrix.value[0].length
        const row = Math.floor(Math.random() * rows)
        const col = Math.floor(Math.random() * cols)
        
        highlightedCell.value = { row, col }
      }
      
      // 检查单元格是否高亮
      const isHighlighted = (row, col) => {
        return highlightedCell.value.row === row && highlightedCell.value.col === col && isAnswering.value
      }
      
      // 获取单元格样式
      const getCellStyle = (row, col) => {
        if (isHighlighted(row, col)) {
          return { color: currentColor.value }
        }
        return { color: 'gray' }
      }
      
      // 获取按钮颜色
      const getButtonColor = (colorValue) => {
        const colorMap = {
          red: '#ff4757',
          blue: '#3742fa',
          green: '#2ed573',
          yellow: '#ffa502',
          black: '#2f3542'
        }
        return colorMap[colorValue] || 'gray'
      }
      
      // 检查答案
      const checkAnswer = (selectedColor) => {
        if (!isAnswering.value) return
        
        totalQuestions.value++
        if (selectedColor === currentColor.value) {
          score.value += 10
          correctAnswers.value++
        }
        
        isAnswering.value = false
      }
      
      // 开始游戏
      const startGame = async () => {
        await fetchMatrix()
        connectWebSocket()
        gameStarted.value = true
        gameEnded.value = false
        score.value = 0
        totalQuestions.value = 0
        correctAnswers.value = 0
        remainingTime.value = gameDuration
        
        // 启动计时器
        timerInterval.value = setInterval(() => {
          remainingTime.value--
          
          if (remainingTime.value <= 0) {
            endGame()
          }
        }, 1000)
      }
      
      // 结束游戏
      const endGame = () => {
        clearInterval(timerInterval.value)
        gameEnded.value = true
        
        // 关闭WebSocket连接
        if (ws.value) {
          ws.value.close()
        }
      }
      
      // 重新开始游戏
      const restartGame = () => {
        gameStarted.value = false
        gameEnded.value = false
        isAnswering.value = false
      }
      
      // 提交成绩
      const submitScore = async () => {
        try {
          loading.value = true
          
          // 假设我们有用户ID，实际应用中应从认证信息获取
          const userId = localStorage.getItem('userId') || 'anonymous'
          
          const scoreData = {
            userId: userId,
            success_num: correctAnswers.value,
            level: selectedDifficulty.value,
            accuracy: accuracy.value / 100,
            training_num: totalQuestions.value
          }
          
          await axios.post('/api/v1/color_words/scores', scoreData)
          alert('成绩保存成功!')
        } catch (error) {
          console.error('提交成绩失败:', error)
          alert('保存成绩失败，请重试')
        } finally {
          loading.value = false
        }
      }
      
      // 选择难度
      const selectDifficulty = (difficulty) => {
        selectedDifficulty.value = difficulty
      }
      
      // 格式化时间显示
      const formatTime = (seconds) => {
        const mins = Math.floor(seconds / 60)
        const secs = seconds % 60
        return `${mins}:${secs < 10 ? '0' : ''}${secs}`
      }
      
      // 组件卸载时清理
      onUnmounted(() => {
        if (timerInterval.value) {
          clearInterval(timerInterval.value)
        }
        
        if (ws.value) {
          ws.value.close()
        }
      })
      
      return {
        gameStarted,
        gameEnded,
        isAnswering,
        loading,
        matrix,
        selectedDifficulty,
        score,
        remainingTime,
        accuracy,
        difficultyLevels,
        colorOptions,
        isHighlighted,
        getCellStyle,
        getButtonColor,
        checkAnswer,
        startGame,
        restartGame,
        submitScore,
        selectDifficulty,
        formatTime
      }
    }
  }
  </script>
  
  <style scoped>
  .color-words-training {
    max-width: 800px;
    margin: 0 auto;
    padding: 20px;
    font-family: 'Arial', sans-serif;
  }
  
  .header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 20px;
  }
  
  .stats {
    display: flex;
    gap: 15px;
  }
  
  .setup-screen {
    text-align: center;
    padding: 40px 0;
  }
  
  .difficulty-selector {
    margin-bottom: 30px;
  }
  
  .difficulty-options {
    display: flex;
    justify-content: center;
    gap: 10px;
    margin-top: 15px;
  }
  
  .difficulty-options button {
    padding: 10px 20px;
    border: 2px solid #ddd;
    background: white;
    border-radius: 5px;
    cursor: pointer;
    transition: all 0.3s;
  }
  
  .difficulty-options button.active {
    border-color: #3498db;
    background: #3498db;
    color: white;
  }
  
  .start-button {
    padding: 12px 30px;
    background: #2ecc71;
    color: white;
    border: none;
    border-radius: 5px;
    font-size: 18px;
    cursor: pointer;
    transition: background 0.3s;
  }
  
  .start-button:hover {
    background: #27ae60;
  }
  
  .color-matrix {
    margin: 20px 0;
    display: flex;
    flex-direction: column;
    align-items: center;
  }
  
  .matrix-row {
    display: flex;
  }
  
  .color-cell {
    width: 60px;
    height: 60px;
    display: flex;
    align-items: center;
    justify-content: center;
    margin: 4px;
    border: 2px solid #ddd;
    border-radius: 8px;
    font-size: 20px;
    font-weight: bold;
    transition: all 0.3s;
  }
  
  .color-cell.highlighted {
    border-color: #3498db;
    box-shadow: 0 0 10px rgba(52, 152, 219, 0.5);
    transform: scale(1.05);
  }
  
  .color-buttons {
    display: flex;
    justify-content: center;
    gap: 15px;
    margin: 30px 0;
  }
  
  .color-button {
    padding: 12px 25px;
    border: none;
    border-radius: 5px;
    color: white;
    font-size: 16px;
    font-weight: bold;
    cursor: pointer;
    transition: transform 0.2s, opacity 0.2s;
  }
  
  .color-button:hover:not(:disabled) {
    transform: translateY(-2px);
  }
  
  .color-button:disabled {
    opacity: 0.6;
    cursor: not-allowed;
  }
  
  .results {
    text-align: center;
    padding: 20px;
    border: 2px solid #2ecc71;
    border-radius: 10px;
    margin-top: 20px;
  }
  
  .results button {
    margin: 10px;
    padding: 10px 20px;
    border: none;
    border-radius: 5px;
    cursor: pointer;
  }
  
  .results button:first-child {
    background: #3498db;
    color: white;
  }
  
  .results button:last-child {
    background: #2ecc71;
    color: white;
  }
  
  .loading-overlay {
    position: fixed;
    top: 0;
    left: 0;
    right: 0;
    bottom: 0;
    background: rgba(255, 255, 255, 0.8);
    display: flex;
    flex-direction: column;
    justify-content: center;
    align-items: center;
    z-index: 1000;
  }
  
  .spinner {
    width: 40px;
    height: 40px;
    border: 4px solid #f3f3f3;
    border-top: 4px solid #3498db;
    border-radius: 50%;
    animation: spin 1s linear infinite;
    margin-bottom: 10px;
  }
  
  @keyframes spin {
    0% { transform: rotate(0deg); }
    100% { transform: rotate(360deg); }
  }
  
  @media (max-width: 600px) {
    .color-cell {
      width: 50px;
      height: 50px;
      font-size: 16px;
    }
    
    .color-buttons {
      flex-wrap: wrap;
    }
    
    .header {
      flex-direction: column;
      gap: 10px;
    }
    
    .stats {
      justify-content: center;
    }
  }
  </style>