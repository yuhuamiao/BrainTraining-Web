import { defineStore } from 'pinia';

export const useGameStore = defineStore('busGame', {
  state: () => ({
    difficulty: 'easy', // 当前难度
    currentRound: 0,    // 当前轮次
    totalRounds: 3,     // 总轮次
    currentPassengers: 10, // 当前乘客数
    correctAnswer: 0,   // 正确答案
    gameHistory: []     // 游戏历史记录
  }),
  actions: {
    // 开始新游戏
    startNewGame(difficulty) {
      this.difficulty = difficulty;
      this.currentRound = 0;
      this.currentPassengers = 10; // 初始乘客数
      
      // 根据难度设置总轮次
      this.totalRounds = {
        easy: 3,
        medium: 5,
        hard: 7
      }[difficulty];
    },
    
    // 生成一轮变化数据
    generateRound() {
      const config = {
        easy: { maxChange: 3 },
        medium: { maxChange: 5 },
        hard: { maxChange: 8 }
      }[this.difficulty];
      
      const change = Math.floor(Math.random() * config.maxChange) + 1;
      const isEntering = Math.random() > 0.5;
      
      return {
        change,
        isEntering,
        newCount: isEntering 
          ? this.currentPassengers + change 
          : Math.max(0, this.currentPassengers - change)
      };
    },
    
    // 更新乘客数量
    updatePassengers(change) {
      this.currentPassengers = Math.max(0, this.currentPassengers + change);
    },
    
    // 进入下一轮
    nextRound() {
      if (this.currentRound < this.totalRounds) {
        this.currentRound += 1;
      }
      
      // 在所有轮次完成后记录正确答案
      if (this.currentRound === this.totalRounds) {
        this.correctAnswer = this.currentPassengers;
      }
    }
  }
});