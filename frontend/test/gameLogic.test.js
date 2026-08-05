import assert from 'node:assert/strict'
import test from 'node:test'
import { getSudoku } from 'sudoku-gen'
import { busDifficulty, generateBusSession } from '../src/utils/games/bus.js'

function seededRandom(seed) {
  let state = seed >>> 0
  return () => {
    state = (state * 1664525 + 1013904223) >>> 0
    return state / 4294967296
  }
}

test('bus sessions preserve passenger count invariants', () => {
  for (const level of Object.keys(busDifficulty)) {
    for (let seed = 1; seed <= 100; seed += 1) {
      const session = generateBusSession(level, seededRandom(seed))
      assert.equal(session.events.length, busDifficulty[level].rounds)
      let count = session.startCount
      for (const event of session.events) {
        count += event.entering ? event.amount : -event.amount
        assert.ok(count >= 0)
        assert.equal(count, event.countAfter)
      }
      assert.equal(count, session.finalCount)
    }
  }
})

test('sudoku generator returns a solved value for every given cell', () => {
  for (const level of ['easy', 'medium', 'hard']) {
    const sudoku = getSudoku(level)
    assert.equal(sudoku.puzzle.length, 81)
    assert.match(sudoku.solution, /^[1-9]{81}$/)
    for (let index = 0; index < 81; index += 1) {
      if (sudoku.puzzle[index] !== '-') assert.equal(sudoku.puzzle[index], sudoku.solution[index])
    }
  }
})
