import { ChakraProvider, Box, useColorMode } from '@chakra-ui/react'
import { BrowserRouter as Router, Routes, Route } from 'react-router-dom'
import { theme } from './theme.ts'
import Navbar from './components/Navbar'
import Home from './pages/Home'

function Layout() {
  const { colorMode } = useColorMode()
  return (
    <Box 
      minH="100vh" 
      w="100%" 
      bg={colorMode === 'light' ? 'gray.50' : 'gray.900'}
      display="flex"
      flexDirection="column"
    >
      <Navbar />
      <Box 
        flex="1"
        w="100%"
        bg={colorMode === 'light' ? 'gray.50' : 'gray.900'}
      >
        <Routes>
          <Route path="/" element={<Home />} />
        </Routes>
      </Box>
    </Box>
  )
}

function App() {
  return (
    <ChakraProvider theme={theme}>
      <Router>
        <Layout />
      </Router>
    </ChakraProvider>
  )
}

export default App
