import React, { useEffect, useState } from 'react';
import { Box, Center, Image, Flex, Text, Spinner } from '@chakra-ui/react';
import { keyframes } from '@emotion/react';

interface SplashScreenProps {
  onInitializationComplete: () => void;
  minDisplayTime?: number;
}

const fadeIn = keyframes`
  from { opacity: 0; }
  to { opacity: 1; }
`;

const fadeOut = keyframes`
  from { opacity: 1; }
  to { opacity: 0; }
`;

const pulse = keyframes`
  0% { transform: scale(1); }
  50% { transform: scale(1.05); }
  100% { transform: scale(1); }
`;

const SplashScreen: React.FC<SplashScreenProps> = ({ 
  onInitializationComplete, 
  minDisplayTime = 2000 
}) => {
  const [isVisible, setIsVisible] = useState(true);
  const [isAnimatingOut, setIsAnimatingOut] = useState(false);
  const [loadingText, setLoadingText] = useState('Initializing');

  // Cycle through loading text states
  useEffect(() => {
    const loadingStates = ['Initializing', 'Initializing.', 'Initializing..', 'Initializing...'];
    let currentIndex = 0;
    
    const interval = setInterval(() => {
      currentIndex = (currentIndex + 1) % loadingStates.length;
      setLoadingText(loadingStates[currentIndex]);
    }, 500);
    
    return () => clearInterval(interval);
  }, []);

  useEffect(() => {
    const startTime = Date.now();
    
    // This timeout ensures the splash screen is shown for at least minDisplayTime
    const timer = setTimeout(() => {
      const elapsedTime = Date.now() - startTime;
      if (elapsedTime >= minDisplayTime) {
        setIsAnimatingOut(true);
        setTimeout(() => {
          setIsVisible(false);
          onInitializationComplete();
        }, 1000); // Duration of fade out animation
      }
    }, minDisplayTime);

    return () => clearTimeout(timer);
  }, [minDisplayTime, onInitializationComplete]);

  if (!isVisible) return null;

  return (
    <Box
      position="fixed"
      top="0"
      left="0"
      width="100vw"
      height="100vh"
      bg="#121212"
      zIndex="9999"
      animation={isAnimatingOut ? `${fadeOut} 1s ease-out forwards` : `${fadeIn} 1s ease-in`}
    >
      <Center height="100%">
        <Flex direction="column" align="center">
          <Box
            animation={`${pulse} 2s infinite ease-in-out`}
            mb={6}
          >
            <Image 
              src="/thedeeps.png" 
              alt="The Deeps Logo" 
              maxWidth="500px"
            />
          </Box>
          <Flex align="center" mt={6}>
            <Spinner 
              thickness="4px"
              speed="0.65s"
              emptyColor="gray.700"
              color="cyan.300"
              size="md"
              mr={3}
            />
            <Text fontSize="xl" color="cyan.300">{loadingText}</Text>
          </Flex>
        </Flex>
      </Center>
    </Box>
  );
};

export default SplashScreen; 