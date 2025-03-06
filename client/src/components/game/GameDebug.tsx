import { useEffect, useState } from 'react'
import { Box, Text, Button, VStack, HStack, Badge, IconButton } from '@chakra-ui/react'

interface DebugMessage {
  type: string
  timestamp: string
  status?: string
  level?: 'info' | 'warn' | 'error' | 'debug'
  message?: string
}

export const GameDebug = () => {
  const [socket, setSocket] = useState<WebSocket | null>(null)
  const [connected, setConnected] = useState(false)
  const [messages, setMessages] = useState<DebugMessage[]>([])
  const [pingTime, setPingTime] = useState<number | null>(null)
  const [isMinimized, setIsMinimized] = useState(false)

  useEffect(() => {
    console.log('GameDebug mounted')
    const ws = new WebSocket('ws://localhost:8080/ws')

    ws.onopen = () => {
      setConnected(true)
      console.log('Connected to server')
      addDebugMessage('Connected to server', 'info')
    }

    ws.onclose = () => {
      setConnected(false)
      console.log('Disconnected from server')
      addDebugMessage('Disconnected from server', 'error')
    }

    ws.onerror = (error) => {
      console.error('WebSocket error:', error)
      addDebugMessage('WebSocket error occurred', 'error')
    }

    ws.onmessage = (event) => {
      console.log('Received message:', event.data)
      const data = JSON.parse(event.data)
      if (data.type === 'pong') {
        const endTime = Date.now()
        setPingTime(endTime - startTime)
        addDebugMessage(`Ping response: ${endTime - startTime}ms`, 'debug')
      } else if (data.type === 'debug') {
        addDebugMessage(data.message, data.level || 'debug')
      }
      setMessages(prev => [...prev, data].slice(-5))
    }

    setSocket(ws)

    return () => {
      ws.close()
    }
  }, [])

  let startTime = 0

  const addDebugMessage = (message: string, level: DebugMessage['level'] = 'debug') => {
    const newMessage: DebugMessage = {
      type: 'debug',
      timestamp: new Date().toISOString(),
      message,
      level
    }
    setMessages(prev => [...prev, newMessage].slice(-5))
  }

  const sendPing = () => {
    if (socket && socket.readyState === WebSocket.OPEN) {
      startTime = Date.now()
      socket.send(JSON.stringify({ type: 'ping' }))
      addDebugMessage('Sending ping request', 'debug')
    }
  }

  const getLevelColor = (level?: string) => {
    switch (level) {
      case 'error': return 'red.500'
      case 'warn': return 'yellow.500'
      case 'info': return 'blue.500'
      default: return 'gray.500'
    }
  }

  if (isMinimized) {
    return (
      <Box
        position="fixed"
        bottom={4}
        right={4}
        bg="blackAlpha.800"
        p={2}
        borderRadius="md"
        color="white"
        zIndex={9999}
        boxShadow="lg"
        border="1px solid"
        borderColor="whiteAlpha.200"
        cursor="pointer"
        onClick={() => setIsMinimized(false)}
      >
        <HStack spacing={2}>
          <Text fontSize="sm">Debug</Text>
          <Badge colorScheme={connected ? 'green' : 'red'} variant="solid" size="sm">
            {connected ? 'Connected' : 'Disconnected'}
          </Badge>
        </HStack>
      </Box>
    )
  }

  return (
    <Box
      position="fixed"
      bottom={4}
      right={4}
      bg="blackAlpha.800"
      p={4}
      borderRadius="md"
      color="white"
      maxW="400px"
      zIndex={9999}
      boxShadow="lg"
      border="1px solid"
      borderColor="whiteAlpha.200"
    >
      <VStack gap={3} alignItems="stretch">
        <HStack justify="space-between">
          <Text fontWeight="bold">Debug Panel</Text>
          <HStack>
            <Badge colorScheme={connected ? 'green' : 'red'} px={2} py={1}>
              {connected ? 'Connected' : 'Disconnected'}
            </Badge>
            <IconButton
              aria-label="Minimize"
              icon={<Text fontSize="lg" color="white">−</Text>}
              size="sm"
              onClick={() => setIsMinimized(true)}
              bg="whiteAlpha.200"
              _hover={{ bg: 'whiteAlpha.400' }}
              _active={{ bg: 'whiteAlpha.500' }}
              borderRadius="md"
            />
          </HStack>
        </HStack>

        <Button 
          size="sm" 
          onClick={sendPing} 
          disabled={!connected}
          colorScheme="blue"
          _hover={{ bg: 'blue.600' }}
        >
          Send Ping
        </Button>

        {pingTime && (
          <Text fontSize="sm">
            Last ping: {pingTime}ms
          </Text>
        )}

        <Box>
          <Text fontSize="sm" mb={2}>Last 5 messages:</Text>
          <VStack gap={1} alignItems="stretch">
            {messages.map((msg, i) => (
              <Box
                key={i}
                fontSize="xs"
                bg="whiteAlpha.100"
                p={2}
                borderRadius="sm"
                borderLeft="3px solid"
                borderLeftColor={getLevelColor(msg.level)}
              >
                <HStack justify="space-between" mb={1}>
                  <Text color={getLevelColor(msg.level)}>
                    {msg.type} {msg.level && `[${msg.level}]`}
                  </Text>
                  <Text fontSize="xx-small" color="gray.400">
                    {new Date(msg.timestamp).toLocaleTimeString()}
                  </Text>
                </HStack>
                {msg.message && (
                  <Text color="gray.300">{msg.message}</Text>
                )}
              </Box>
            ))}
          </VStack>
        </Box>
      </VStack>
    </Box>
  )
} 