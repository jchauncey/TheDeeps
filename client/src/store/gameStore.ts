import { create } from 'zustand'

interface Position {
  x: number
  y: number
}

interface GameState {
  playerPosition: Position
  health: number
  maxHealth: number
  level: number
  dungeonLevel: number
  movePlayer: (direction: 'up' | 'down' | 'left' | 'right') => void
}

export const useGameStore = create<GameState>((set) => ({
  playerPosition: { x: 12, y: 10 },
  health: 100,
  maxHealth: 100,
  level: 1,
  dungeonLevel: 1,
  movePlayer: (direction) => {
    set((state) => {
      const newPosition = { ...state.playerPosition }
      switch (direction) {
        case 'up':
          newPosition.y = Math.max(0, newPosition.y - 1)
          break
        case 'down':
          newPosition.y = Math.min(19, newPosition.y + 1)
          break
        case 'left':
          newPosition.x = Math.max(0, newPosition.x - 1)
          break
        case 'right':
          newPosition.x = Math.min(24, newPosition.x + 1)
          break
      }
      return { playerPosition: newPosition }
    })
  },
})) 