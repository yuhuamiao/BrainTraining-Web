<template>
    <div class="result-display">
      <h2 v-if="isCorrect">回答正确！</h2>
      <h2 v-else>回答错误</h2>
      
      <p>正确答案: {{ correctAnswer }}</p>
      
      <button @click="$emit('restart')">再玩一次</button>
    </div>
  </template>
  
  <script>
  export default {
    name: 'ResultDisplay',
    props: {
      isCorrect: Boolean,
      correctAnswer: Number
    },
    emits: ['restart']
  }
  </script>
  
  <style scoped>
  .result-display {
    text-align: center;
    padding: 20px;
  }
  
  .result-display h2 {
    color: #4CAF50;
    font-size: 2rem;
    margin-bottom: 20px;
  }
  
  .result-display p {
    font-size: 1.2rem;
    margin-bottom: 20px;
  }
  
  .result-display button {
    padding: 10px 20px;
    background-color: #4CAF50;
    color: white;
    border: none;
    border-radius: 4px;
    cursor: pointer;
    font-size: 1rem;
  }
  </style>