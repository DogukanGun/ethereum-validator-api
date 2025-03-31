# Ethereum Validator API

A production-ready API service for querying Ethereum validator information, including sync committee duties and block rewards. This project demonstrates my ability to build scalable, maintainable blockchain infrastructure using modern development practices.

## Project Overview

You can test the application using either the frontend interface or direct API calls:

- **Frontend Interface**: [https://sf.dogukangun.de](https://sf.dogukangun.de)
- **API Endpoint**: [https://sf-api.dogukangun.de](https://sf-api.dogukangun.de)
- **API Documentation**: [https://sf-api.dogukangun.de/swagger/index.html](https://sf-api.dogukangun.de/swagger/index.html)

Test the application:
- Frontend: [https://sf.dogukangun.de](https://sf.dogukangun.de)
- API: [https://sf-api.dogukangun.de](https://sf-api.dogukangun.de)
- API Documentation: [https://sf-api.dogukangun.de/swagger/index.html](https://sf-api.dogukangun.de/swagger/index.html)

## Architecture

The project follows a clean, modular architecture:

```
ethereum-validator-api/
├── handler/                # HTTP request handlers
│   ├── blockrewardHandler.go
│   ├── syncdutiesHandler.go
│   ├── handler.go
│   └── types.go
├── service/               # Business logic layer
│   ├── ethereumService.go
│   └── ethereumService_test.go
├── utils/                 # Shared utilities
│   ├── env.go
│   └── setupEndpoints.go
├── tests/                 # Integration tests
│   ├── ethereumService_test.go
│   └── env_utils_test.go
├── platform/              # Frontend application
│   ├── app/              
│   │   ├── page.tsx      # Main interface
│   │   ├── layout.tsx    # App layout
│   │   └── globals.css
│   ├── public/           
│   └── package.json
├── main.go               # Application entry
├── go.mod               # Go dependencies
├── docker-compose.yml   # Docker config
└── Dockerfile           # API service build
```

## Design Choices

### Backend Architecture
1. **Clean Architecture**
   - Clear separation of concerns (handlers, services, utils)
   - Dependency flow from outer to inner layers
   - Testable and maintainable design

2. **Handler Layer**
   - Request/response handling
   - Input validation
   - Error management

3. **Service Layer**
   - Contains core business logic for Ethereum interactions
   - Encapsulated Ethereum node communication
   - Comprehensive test coverage

4. **Utils Layer**
   - Environment configuration management
   - Endpoint setup and routing
   - Reusable helper functions

### Frontend Architecture
- Next.js App Router for modern React patterns
- Server components for optimal performance
- TypeScript for type safety
- TailwindCSS for responsive design

## API Endpoints

### 1. Get Sync Committee Duties
```bash
curl -X GET 'http://localhost:3004/syncduties/4700000' \
  -H 'Accept: application/json'
```

Response:
```json
{
  "validators": [
    "0x1234...",
    "0x5678..."
  ],
  "sync_info": {
    "sync_period": 123,
    "committee_size": 512
  }
}
```

### 2. Get Block Rewards
```bash
curl -X GET 'http://localhost:3004/blockreward/4700000' \
  -H 'Accept: application/json'
```

Response:
```json
{
  "status": "mev",
  "reward": 123456,
  "block_info": {
    "proposer_payment": 100000,
    "is_mev_boost": true
  }
}
```

## Building and Running

### Prerequisites
- Go 1.21+
- Node.js v18+
- Docker and Docker Compose (optional)

### Method 1: Local Development

1. **Backend Setup**
```bash
cd ethereum-validator-api

# Install Go dependencies
go mod download

# Set up environment variables
cp .env.example .env
# Edit .env with your configuration

# Run the backend
go run main.go
```

2. **Frontend Setup**
```bash
# Navigate to frontend directory
cd platform

# Install dependencies
npm install

# Set up environment variables
cp .env.example .env.local
# Edit .env.local with your configuration

# Run the development server
npm run dev
```

### Method 2: Docker Deployment

```bash
# Build and run all services
docker-compose up --build

# Or run services separately
docker-compose up api
docker-compose up frontend
```

The services will be available at:
- API: http://localhost:3004
- Frontend: http://localhost:3003
- Swagger Documentation: http://localhost:3004/swagger/index.html

## Environment Variables

### Backend (.env)
```env
ETH_RPC=<ethereum-node-url>
CORS_ORIGIN=http://localhost:3003
```

### Frontend (.env.local)
```env
NEXT_PUBLIC_API_URL=http://localhost:3004
```

## Testing

You have three options to test the API:

### 1. Frontend Interface

The easiest way to test the API is through the frontend interface:

1. Visit [https://sf.dogukangun.de](https://sf.dogukangun.de)
2. Enter a slot number in the input field (e.g., 4700000)
3. Click "Fetch Data" to see both sync committee duties and block rewards
4. The interface provides example slot numbers and explanations

### 2. Postman Collection

For development and testing, you can use the provided Postman collection:

1. Download the [postman_collection.json](./postman_collection.json) file
2. Open Postman and click "Import"
3. Drag and drop the downloaded file or browse to select it
4. The collection includes both endpoints with example responses
5. The base URL is pre-configured to the live API

You can switch between environments by changing the `baseUrl` variable:
- Production: `https://sf-api.dogukangun.de`
- Local: `http://localhost:3004`

### 3. Direct API Calls

You can also make direct API calls using curl or any HTTP client:

## Frameworks and Libraries Used

### Backend
- **Gin**: High-performance HTTP web framework
- **go-ethereum**: Ethereum client implementation
- **swagger/swag**: API documentation
- **testify**: Testing framework

### Frontend
- **Next.js 15.2.4**: React framework
- **React 19**: UI library
- **TailwindCSS**: Utility-first CSS framework
- **TypeScript**: Type-safe JavaScript

## Author

Dogukan Gundogan