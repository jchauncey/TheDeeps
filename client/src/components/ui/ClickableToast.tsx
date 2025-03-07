import { Box, useToast, UseToastOptions } from '@chakra-ui/react';

interface ClickableToastProps extends UseToastOptions {
  onClose?: () => void;
}

export const useClickableToast = () => {
  const toast = useToast();

  const showClickableToast = (props: ClickableToastProps) => {
    const { title, description, status, duration = 5000, position = 'top', onClose } = props;

    return toast({
      position,
      duration,
      render: ({ onClose: closeToast }) => (
        <Box
          color='white'
          p={3}
          bg={status === 'success' ? 'green.500' : status === 'error' ? 'red.500' : status === 'info' ? 'blue.500' : 'gray.500'}
          borderRadius='md'
          onClick={() => {
            closeToast();
            if (onClose) onClose();
          }}
          cursor="pointer"
          boxShadow="md"
        >
          {title && <Box fontWeight="bold" mb={1}>{title}</Box>}
          {description && <Box>{description}</Box>}
        </Box>
      ),
    });
  };

  return showClickableToast;
}; 