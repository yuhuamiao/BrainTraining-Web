export const generateRoundData = (difficulty, currentCount) => {
    const config = {
      easy: { maxChange: 3, speed: 1500 },
      medium: { maxChange: 5, speed: 1000 },
      hard: { maxChange: 8, speed: 600 }
    }
    
    const { maxChange, speed } = config[difficulty]
    const change = Math.floor(Math.random() * maxChange) + 1
    const isEntering = Math.random() > 0.5
    
    return {
      change,
      isEntering,
      speed,
      newCount: isEntering 
        ? currentCount + change 
        : Math.max(0, currentCount - change)
    }
  }