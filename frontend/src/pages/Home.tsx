import { useState } from 'react'
import {
  Box,
  Button,
  Card,
  CardBody,
  FormControl,
  FormLabel,
  Grid,
  Heading,
  Input,
  Stack,
  Text,
  useToast,
  VStack,
  useColorMode,
  Icon,
  Center,
} from '@chakra-ui/react'
import { ChevronDownIcon, ChevronUpIcon } from '@chakra-ui/icons'
import axios from 'axios'

interface GreetingResponse {
  greeting: string
  stats?: {
    apiCalls: number
    charsSent: number
    costEstimate: number
    cacheHits: number
  }
}

// Primary languages shown by default
const PRIMARY_LANGUAGES = [
  { code: 'en', name: 'English' },
  { code: 'es', name: 'Spanish' },
  { code: 'fr', name: 'French' },
  { code: 'de', name: 'German' },
  { code: 'zh', name: 'Chinese' },
  { code: 'ja', name: 'Japanese' },
  { code: 'ko', name: 'Korean' },
  { code: 'hi', name: 'Hindi' },
]

// Additional languages shown when expanded
const SECONDARY_LANGUAGES = [
  { code: 'it', name: 'Italian' },
  { code: 'pt', name: 'Portuguese' },
  { code: 'nl', name: 'Dutch' },
  { code: 'pl', name: 'Polish' },
  { code: 'ru', name: 'Russian' },
  { code: 'ar', name: 'Arabic' },
  { code: 'tr', name: 'Turkish' },
  { code: 'vi', name: 'Vietnamese' },
  { code: 'th', name: 'Thai' },
  { code: 'sv', name: 'Swedish' },
  { code: 'da', name: 'Danish' },
  { code: 'fi', name: 'Finnish' },
  { code: 'el', name: 'Greek' },
  { code: 'he', name: 'Hebrew' },
  { code: 'id', name: 'Indonesian' },
  { code: 'ms', name: 'Malay' },
  { code: 'fa', name: 'Persian' },
  { code: 'bn', name: 'Bengali' },
  { code: 'ta', name: 'Tamil' },
  { code: 'uk', name: 'Ukrainian' },
]

const Home = () => {
  const { colorMode } = useColorMode()
  const [name, setName] = useState('')
  const [loading, setLoading] = useState('')
  const [result, setResult] = useState<GreetingResponse | null>(null)
  const [showMore, setShowMore] = useState(false)
  const toast = useToast()

  const handleSubmit = async (selectedLanguage: string) => {
    if (!name.trim()) {
      toast({
        title: 'Name Required',
        description: 'Please enter your name first.',
        status: 'warning',
        duration: 3000,
        isClosable: true,
      })
      return
    }

    setLoading(selectedLanguage)
    try {
      const response = await axios.get<GreetingResponse>('/api/greet', {
        params: { name, language: selectedLanguage },
      })
      setResult(response.data)
    } catch (error) {
      toast({
        title: 'Error',
        description: 'Failed to get greeting. Please try again.',
        status: 'error',
        duration: 5000,
        isClosable: true,
      })
    } finally {
      setLoading('')
    }
  }

  const renderLanguageButtons = (languages: typeof PRIMARY_LANGUAGES) => (
    <Grid
      templateColumns={{
        base: 'repeat(2, 1fr)',
        sm: 'repeat(3, 1fr)',
        md: 'repeat(4, 1fr)',
        lg: 'repeat(4, 1fr)',
      }}
      gap={3}
    >
      {languages.map((lang) => (
        <Button
          key={lang.code}
          onClick={() => handleSubmit(lang.code)}
          isLoading={loading === lang.code}
          variant="outline"
          size="lg"
          height="auto"
          py={3}
          whiteSpace="normal"
          borderColor={colorMode === 'light' ? 'gray.200' : 'gray.600'}
          _hover={{
            bg: colorMode === 'light' ? 'brand.50' : 'brand.900',
            borderColor: colorMode === 'light' ? 'brand.500' : 'brand.400',
          }}
        >
          {lang.name}
        </Button>
      ))}
    </Grid>
  )

  return (
    <VStack 
      spacing={8} 
      w="100%" 
      maxW="container.xl"
      mx="auto"
      px={4}
      py={8}
    >
      <Box textAlign="center" w="100%" px={4}>
        <Heading 
          size="2xl" 
          mb={4} 
          bgGradient={colorMode === 'light' 
            ? 'linear(to-r, brand.500, brand.700)' 
            : 'linear(to-r, brand.400, brand.600)'
          }
          bgClip="text"
        >
          Welcome to Greeter
        </Heading>
        <Text fontSize="lg" color={colorMode === 'light' ? 'gray.600' : 'gray.400'}>
          Get personalized greetings in multiple languages
        </Text>
      </Box>

      <Grid
        templateColumns={{ base: '1fr', lg: '1fr 1fr' }}
        gap={8}
        w="100%"
        alignItems="start"
      >
        <Card w="100%">
          <CardBody>
            <VStack spacing={6}>
              <FormControl isRequired>
                <FormLabel>Your Name</FormLabel>
                <Input
                  value={name}
                  onChange={(e) => setName(e.target.value)}
                  placeholder="Enter your name"
                  size="lg"
                  bg={colorMode === 'light' ? 'white' : 'gray.700'}
                  borderColor={colorMode === 'light' ? 'gray.200' : 'gray.600'}
                  _hover={{
                    borderColor: colorMode === 'light' ? 'brand.500' : 'brand.400',
                  }}
                />
              </FormControl>

              <Box w="100%">
                <Text 
                  fontSize="lg" 
                  fontWeight="medium" 
                  mb={4}
                  color={colorMode === 'light' ? 'gray.700' : 'gray.300'}
                >
                  Select a Language
                </Text>
                
                <VStack spacing={3} align="stretch" w="100%">
                  {renderLanguageButtons(PRIMARY_LANGUAGES)}
                  
                  <Box
                    style={{
                      maxHeight: showMore ? '1000px' : '0',
                      overflow: 'hidden',
                      transition: 'max-height 0.3s ease-in-out'
                    }}
                  >
                    {renderLanguageButtons(SECONDARY_LANGUAGES)}
                  </Box>
                </VStack>

                <Button
                  onClick={() => setShowMore(!showMore)}
                  variant="ghost"
                  size="md"
                  mt={4}
                  w="full"
                  rightIcon={<Icon as={showMore ? ChevronUpIcon : ChevronDownIcon} />}
                  color={colorMode === 'light' ? 'gray.600' : 'gray.400'}
                  _hover={{
                    bg: colorMode === 'light' ? 'gray.100' : 'gray.700',
                  }}
                >
                  {showMore ? 'Show Less Languages' : 'Show More Languages'}
                </Button>
              </Box>
            </VStack>
          </CardBody>
        </Card>

        <Card 
          w="100%" 
          bg={result ? (colorMode === 'light' ? 'brand.50' : 'brand.900') : (colorMode === 'light' ? 'gray.50' : 'gray.800')}
          height="100%"
          minH="400px"
          opacity={result ? 1 : 0.8}
        >
          <CardBody>
            {result ? (
              <>
                <Text 
                  fontSize="2xl" 
                  textAlign="center" 
                  color={colorMode === 'light' ? 'brand.700' : 'brand.200'} 
                  fontWeight="bold"
                  mb={4}
                >
                  {result.greeting}
                </Text>

                {result.stats && (
                  <Box 
                    pt={4} 
                    borderTop="1px" 
                    borderColor={colorMode === 'light' ? 'brand.100' : 'brand.800'}
                  >
                    <Text 
                      fontSize="sm" 
                      color={colorMode === 'light' ? 'gray.600' : 'gray.400'} 
                      mb={2}
                    >
                      Translation Statistics:
                    </Text>
                    <Stack 
                      spacing={1} 
                      fontSize="sm" 
                      color={colorMode === 'light' ? 'gray.600' : 'gray.400'}
                    >
                      <Text>API Calls: {result.stats?.apiCalls ?? 0}</Text>
                      <Text>Characters Translated: {result.stats?.charsSent ?? 0}</Text>
                      <Text>
                        Estimated Cost: ${result.stats ? result.stats.costEstimate.toFixed(5) : '0.00000'}
                      </Text>
                      <Text>Cache Hits: {result.stats?.cacheHits ?? 0}</Text>
                    </Stack>
                  </Box>
                )}
              </>
            ) : (
              <Center h="100%" flexDirection="column" textAlign="center">
                <Text
                  fontSize="xl"
                  color={colorMode === 'light' ? 'gray.600' : 'gray.400'}
                  mb={4}
                >
                  Your greeting will appear here
                </Text>
                <Text
                  fontSize="md"
                  color={colorMode === 'light' ? 'gray.500' : 'gray.500'}
                >
                  Enter your name and select a language to get started
                </Text>
              </Center>
            )}
          </CardBody>
        </Card>
      </Grid>
    </VStack>
  )
}

export default Home 