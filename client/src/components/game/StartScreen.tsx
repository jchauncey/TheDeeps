import { Box, Image, Flex, Text, Spinner } from '@chakra-ui/react'
import { useState, useEffect } from 'react'

interface StartScreenProps {
  onCharactersLoaded: () => void;
  isLoading?: boolean;
}

export const StartScreen = ({ onCharactersLoaded, isLoading = false }: StartScreenProps) => {
  const [imageLoaded, setImageLoaded] = useState(false);
  
  // Trigger the onCharactersLoaded callback after the logo has loaded
  useEffect(() => {
    if (imageLoaded && !isLoading) {
      // Add a small delay for a better visual experience
      const timer = setTimeout(() => {
        onCharactersLoaded();
      }, 1000);
      
      return () => clearTimeout(timer);
    }
  }, [imageLoaded, isLoading, onCharactersLoaded]);

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
          position="relative"
        >
          <Image
            src="/logo.png"
            alt="The Deeps Logo"
            maxW="600px"
            w="100%"
            objectFit="contain"
            onLoad={() => setImageLoaded(true)}
            onError={(e) => console.error('Logo failed to load:', e)}
            opacity={imageLoaded ? 1 : 0}
            transition="opacity 0.5s ease-in-out"
          />
          {!imageLoaded && (
            <Text color="white" fontSize="2xl" position="absolute">Loading logo...</Text>
          )}
        </Box>
        
        {isLoading && (
          <Flex direction="column" align="center" mt={8}>
            <Spinner size="xl" color="purple.500" mb={4} />
            <Text color="white" fontSize="xl">Loading game data...</Text>
          </Flex>
        )}
      </Flex>
    </Box>
  )
} 