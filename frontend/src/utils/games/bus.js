export const busDifficulty = {
  easy: { rounds: 3, maxChange: 3, delay: 1500, startMin: 6, startRange: 5 },
  medium: { rounds: 5, maxChange: 5, delay: 1050, startMin: 10, startRange: 7 },
  hard: { rounds: 7, maxChange: 8, delay: 720, startMin: 15, startRange: 10 },
}

export function generateBusSession(level, random = Math.random) {
  const config = busDifficulty[level] || busDifficulty.easy
  const startCount = config.startMin + Math.floor(random() * config.startRange)
  let currentCount = startCount
  const events = []

  for (let index = 0; index < config.rounds; index += 1) {
    const entering = currentCount === 0 || random() >= 0.45
    const limit = entering ? config.maxChange : Math.min(config.maxChange, currentCount)
    const amount = Math.max(1, 1 + Math.floor(random() * limit))
    currentCount += entering ? amount : -amount
    events.push({ entering, amount, countAfter: currentCount })
  }

  return { startCount, events, finalCount: currentCount, delay: config.delay }
}
