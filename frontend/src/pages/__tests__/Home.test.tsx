import { render, screen, fireEvent, waitFor } from '@testing-library/react'
import '@testing-library/jest-dom'
import { ChakraProvider } from '@chakra-ui/react'
import axios from 'axios'
import Home from '../Home'
import { theme } from '../../theme'

// Mock axios
jest.mock('axios')
const mockedAxios = axios as jest.Mocked<typeof axios>

describe('Home Component', () => {
  beforeEach(() => {
    render(
      <ChakraProvider theme={theme}>
        <Home />
      </ChakraProvider>
    )
  })

  it('renders main elements', () => {
    expect(screen.getByText('Welcome to Greeter')).toBeInTheDocument()
    expect(screen.getByPlaceholderText('Enter your name')).toBeInTheDocument()
    expect(screen.getByText('Select a Language')).toBeInTheDocument()
    expect(screen.getByText('English')).toBeInTheDocument()
  })

  it('shows validation message when submitting without name', async () => {
    const englishButton = screen.getByText('English')
    fireEvent.click(englishButton)

    await waitFor(() => {
      expect(screen.getByText('Name Required')).toBeInTheDocument()
    })
  })

  it('shows more languages when clicking show more', () => {
    const showMoreButton = screen.getByText('Show More Languages')
    fireEvent.click(showMoreButton)
    
    expect(screen.getByText('Greek')).toBeInTheDocument()
    expect(screen.getByText('Show Less Languages')).toBeInTheDocument()
  })

  it('successfully submits greeting request', async () => {
    const mockResponse = {
      data: {
        greeting: 'Hello, Test!',
        stats: {
          apiCalls: 1,
          charsSent: 10,
          costEstimate: 0.00001,
          cacheHits: 0
        }
      }
    }
    mockedAxios.get.mockResolvedValueOnce(mockResponse)

    // Enter name
    const nameInput = screen.getByPlaceholderText('Enter your name')
    fireEvent.change(nameInput, { target: { value: 'Test' } })

    // Click language button
    const englishButton = screen.getByText('English')
    fireEvent.click(englishButton)

    // Wait for and verify greeting
    await waitFor(() => {
      expect(screen.getByText('Hello, Test!')).toBeInTheDocument()
      expect(screen.getByText('Translation Statistics:')).toBeInTheDocument()
      expect(screen.getByText('API Calls: 1')).toBeInTheDocument()
    })

    // Verify API call
    expect(mockedAxios.get).toHaveBeenCalledWith('/api/greet', {
      params: { name: 'Test', language: 'en' }
    })
  })

  it('handles API error gracefully', async () => {
    mockedAxios.get.mockRejectedValueOnce(new Error('API Error'))

    // Enter name
    const nameInput = screen.getByPlaceholderText('Enter your name')
    fireEvent.change(nameInput, { target: { value: 'Test' } })

    // Click language button
    const englishButton = screen.getByText('English')
    fireEvent.click(englishButton)

    // Wait for and verify error message
    await waitFor(() => {
      expect(screen.getByText('Failed to get greeting. Please try again.')).toBeInTheDocument()
    })
  })
}) 