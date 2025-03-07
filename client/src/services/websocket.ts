import { io, Socket } from 'socket.io-client'

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

    this.socket.on('gameState', (_state: any) => {
      // Update game state when receiving updates from server
      // TODO: Implement state update logic when needed
      // const store = useGameStore.getState()
      // store.updateState(_state)
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