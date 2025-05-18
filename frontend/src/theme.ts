import { extendTheme } from '@chakra-ui/react'

export const theme = extendTheme({
  config: {
    initialColorMode: 'light',
    useSystemColorMode: true,
  },
  colors: {
    brand: {
      50: '#f0f9ff',
      100: '#e0f2fe',
      200: '#bae6fd',
      300: '#7dd3fc',
      400: '#38bdf8',
      500: '#0ea5e9',
      600: '#0284c7',
      700: '#0369a1',
      800: '#075985',
      900: '#0c4a6e',
    },
  },
  styles: {
    global: (props: { colorMode: 'light' | 'dark' }) => ({
      '*': {
        margin: 0,
        padding: 0,
        boxSizing: 'border-box',
      },
      'html, body, #root': {
        bg: props.colorMode === 'light' ? 'gray.50' : 'gray.900',
        color: props.colorMode === 'light' ? 'gray.800' : 'gray.100',
        minHeight: '100vh',
        margin: 0,
        padding: 0,
        width: '100%',
      },
      'div#root': {
        display: 'flex',
        flexDirection: 'column',
        minHeight: '100vh',
        width: '100%',
        overflow: 'hidden',
      }
    }),
  },
  components: {
    Button: {
      defaultProps: {
        colorScheme: 'brand',
      },
    },
    Card: {
      baseStyle: (props: { colorMode: 'light' | 'dark' }) => ({
        container: {
          bg: props.colorMode === 'light' ? 'white' : 'gray.800',
          borderColor: props.colorMode === 'light' ? 'gray.200' : 'gray.700',
        },
      }),
    },
    Heading: {
      baseStyle: (props: { colorMode: 'light' | 'dark' }) => ({
        color: props.colorMode === 'light' ? 'gray.800' : 'white',
      }),
    },
  },
}) 