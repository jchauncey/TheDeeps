import { io, Socket } from 'socket.io-client'
import { useGameStore } from '../store/gameStore'

class WebSocketService {
  private socket: Socket | null = null

  connect() {
    if (this.socket) return

    this.socket = io('http://localhost:8080', {
      reconnection: true,
      reconnectionAttempts: 5,
      reconnectionDelay: 1000,
    })

    this.socket.on('connect', () => {
      console.log('Connected to game server')
    })

    this.socket.on('gameState', (state: any) => {
      // Update game state when receiving updates from server
      const store = useGameStore.getState()
      // TODO: Update store with received state
    })

    this.socket.on('disconnect', () => {
      console.log('Disconnected from game server')
    })
  }

  disconnect() {
    if (this.socket) {
      this.socket.disconnect()
      this.socket = null
    }
  }

  sendMove(direction: 'up' | 'down' | 'left' | 'right') {
    if (this.socket) {
      this.socket.emit('move', { direction })
    }
  }
}

export const wsService = new WebSocketService() 