import { Box, Flex, Heading, IconButton, useColorMode } from '@chakra-ui/react'
import { FaMoon, FaSun } from 'react-icons/fa'

const Navbar = () => {
  const { colorMode, toggleColorMode } = useColorMode()

  return (
    <Box 
      bg={colorMode === 'light' ? 'white' : 'gray.800'} 
      shadow="sm" 
      position="sticky" 
      top={0} 
      zIndex={1}
      borderBottom="1px"
      borderColor={colorMode === 'light' ? 'gray.200' : 'gray.700'}
    >
      <Flex
        maxW="container.xl"
        mx="auto"
        px={4}
        h={16}
        alignItems="center"
        justifyContent="space-between"
      >
        <Heading 
          size="lg" 
          bgGradient={colorMode === 'light' 
            ? 'linear(to-r, brand.500, brand.700)' 
            : 'linear(to-r, brand.400, brand.600)'
          } 
          bgClip="text"
        >
          Greeter
        </Heading>
        <IconButton
          aria-label="Toggle color mode"
          icon={colorMode === 'light' ? <FaMoon /> : <FaSun />}
          onClick={toggleColorMode}
          variant="ghost"
          colorScheme="brand"
          fontSize="20px"
          size="lg"
          _hover={{
            bg: colorMode === 'light' ? 'brand.50' : 'brand.900',
          }}
        />
      </Flex>
    </Box>
  )
}

export default Navbar 