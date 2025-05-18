# Greeter

A multilingual greeting application that provides personalized greetings in multiple languages using Google Cloud Translation API.

## Features

- Time-sensitive greetings (morning/afternoon/evening/night)
- Support for 28+ languages
- Translation caching to minimize API costs
- Modern web interface with dark/light mode
- Real-time translation statistics
- RESTful API server

## Prerequisites

1. Go 1.x
2. Node.js and npm
3. Google Cloud Platform account with Translation API enabled
4. Service account credentials

## Setup

1. Clone and install dependencies:
```bash
git clone <repository-url>
cd Greeter
go mod download

cd frontend
npm install
cd ..
```

2. Configure Google Cloud:
   - Place your service account key as `credentials.json` in project root
   - Create a `setup.sh` file with the following content:
     ```bash
     #!/bin/bash
     export GOOGLE_APPLICATION_CREDENTIALS="$(pwd)/credentials.json"
     export GOOGLE_CLOUD_PROJECT_ID="your-project-id"  # Replace with your Google Cloud Project ID
     ```
   - Update the `GOOGLE_CLOUD_PROJECT_ID` with your actual project ID
   - Run: `source setup.sh`

## Usage

### Web Interface

1. Build and start:
```bash
cd frontend && npm run build && cd ..
go run . -server
```

2. Visit http://localhost:8080

### CLI

```bash
go run . "John" es  # Greet John in Spanish
go run . "Maria" fr # Greet Maria in French
```

### Development Mode

```bash
# Terminal 1: Start backend
go run . -server

# Terminal 2: Start frontend
cd frontend
npm run dev
```

Visit http://localhost:5173

### Supported Languages

Primary:
- English (en), Spanish (es), French (fr), German (de)
- Chinese (zh), Japanese (ja), Korean (ko), Hindi (hi)

Additional languages available through the web interface.

## Testing

```bash
# Frontend tests
cd frontend && npm test

# Backend tests
go test -v
```

## Environment Variables

Required:
- `GOOGLE_APPLICATION_CREDENTIALS`: Path to credentials JSON
- `GOOGLE_CLOUD_PROJECT_ID`: Google Cloud Project ID

## Note

This application uses Google Cloud Translation API, which is a paid service. The program implements caching to minimize costs, but new translations will incur charges based on character count. 