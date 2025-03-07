import { Box, Image, HStack, Button, Text, Flex } from '@chakra-ui/react'
import { useState } from 'react'

interface StartScreenProps {
  onNewGame: () => void;
  onLoadGame: () => void;
}

export const StartScreen = ({ onNewGame, onLoadGame }: StartScreenProps) => {
  const [imageLoaded, setImageLoaded] = useState(false);

  return (
    <Box
      position="fixed"
      top={0}
      left={0}
      right={0}
      bottom={0}
      bg="#291326"
      display="flex"
      alignItems="center"
      justifyContent="center"
    >
      <Flex 
        direction="column" 
        align="center" 
        justify="center" 
        width="100%" 
        height="100%"
        px={4}
      >
        <Box 
          width="100%" 
          display="flex" 
          justifyContent="center" 
          alignItems="center"
          mb={12}
        >
          <Image
            src="/logo.png"
            alt="The Deeps Logo"
            maxW="600px"
            w="100%"
            objectFit="contain"
            onLoad={() => setImageLoaded(true)}
            onError={(e) => console.error('Logo failed to load:', e)}
          />
          {!imageLoaded && (
            <Text color="white" fontSize="2xl" position="absolute">Loading logo...</Text>
          )}
        </Box>
        <HStack spacing={6}>
          <Button
            size="md"
            bg="#6B46C1"
            color="white"
            width="200px"
            height="50px"
            fontSize="xl"
            onClick={onNewGame}
            _hover={{
              transform: 'scale(1.05)',
              bg: '#805AD5'
            }}
            _active={{
              bg: '#553C9A'
            }}
            border="2px solid"
            borderColor="purple.200"
          >
            New Game
          </Button>
          <Button
            size="md"
            bg="transparent"
            color="white"
            width="200px"
            height="50px"
            fontSize="xl"
            onClick={onLoadGame}
            border="2px solid"
            borderColor="purple.200"
            _hover={{
              transform: 'scale(1.05)',
              bg: 'rgba(107, 70, 193, 0.2)'
            }}
            _active={{
              bg: 'rgba(107, 70, 193, 0.4)'
            }}
          >
            Load Game
          </Button>
        </HStack>
      </Flex>
    </Box>
  )
} 