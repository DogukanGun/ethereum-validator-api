# Ethereum Validator API Platform

This is the frontend platform for the Ethereum Validator API, built with Next.js and React.

## 🚀 Getting Started

### Prerequisites

- Node.js (v18 or higher)
- npm or yarn
- Docker (optional, for containerized deployment)

### Installation

1. Clone the repository:
```bash
git clone <your-repository-url>
cd platform
```

2. Install dependencies:
```bash
npm install
# or
yarn install
```

3. Set up environment variables:
Create a `.env.local` file in the root directory and add necessary environment variables:
```env
NEXT_PUBLIC_API_URL=http://localhost:3003
```

### Development

To run the development server:

```bash
npm run dev
# or
yarn dev
```

The application will be available at `http://localhost:3000`.

### Production Build

To create a production build:

```bash
npm run build
npm start
# or
yarn build
yarn start
```

### Docker Deployment

To build and run using Docker:

```bash
docker build -t ethereum-validator-platform .
docker run -p 3003:3003 ethereum-validator-platform
```

## 🛠 Technology Stack

- **Frontend Framework**: Next.js 15.2.4
- **UI Library**: React 19
- **Styling**: TailwindCSS
- **Language**: TypeScript
- **Development Tools**: ESLint, PostCSS

## 📡 API Endpoints

### Validator Information

#### Get Validator Status
```bash
curl -X GET http://localhost:3003/api/validator/status/{validatorId}
```

Response:
```json
{
  "validatorId": "123",
  "status": "active",
  "balance": "32.5",
  "effectiveness": "98.5"
}
```

#### Get Validator Performance
```bash
curl -X GET http://localhost:3003/api/validator/performance/{validatorId}
```

Response:
```json
{
  "validatorId": "123",
  "proposedBlocks": 50,
  "missedAttestations": 2,
  "rewards24h": "0.01"
}
```

### Staking Operations

#### Initiate Staking
```bash
curl -X POST http://localhost:3003/api/staking/initiate \
  -H "Content-Type: application/json" \
  -d '{
    "amount": "32",
    "withdrawalCredentials": "0x1234...",
    "pubkey": "0x5678..."
  }'
```

Response:
```json
{
  "transactionHash": "0xabcd...",
  "status": "pending"
}
```

#### Get Staking Rewards
```bash
curl -X GET http://localhost:3003/api/staking/rewards/{address}
```

Response:
```json
{
  "totalRewards": "1.5",
  "lastReward": "0.01",
  "apr": "4.5"
}
```

## 📦 Postman Collection

You can import the following Postman collection to test the API endpoints:

```json
{
  "info": {
    "name": "Ethereum Validator API",
    "schema": "https://schema.getpostman.com/json/collection/v2.1.0/collection.json"
  },
  "item": [
    {
      "name": "Get Validator Status",
      "request": {
        "method": "GET",
        "url": "{{baseUrl}}/api/validator/status/:validatorId",
        "path": {
          "validatorId": "123"
        }
      }
    },
    {
      "name": "Get Validator Performance",
      "request": {
        "method": "GET",
        "url": "{{baseUrl}}/api/validator/performance/:validatorId",
        "path": {
          "validatorId": "123"
        }
      }
    },
    {
      "name": "Initiate Staking",
      "request": {
        "method": "POST",
        "url": "{{baseUrl}}/api/staking/initiate",
        "header": {
          "Content-Type": "application/json"
        },
        "body": {
          "mode": "raw",
          "raw": {
            "amount": "32",
            "withdrawalCredentials": "0x1234...",
            "pubkey": "0x5678..."
          }
        }
      }
    },
    {
      "name": "Get Staking Rewards",
      "request": {
        "method": "GET",
        "url": "{{baseUrl}}/api/staking/rewards/:address",
        "path": {
          "address": "0x1234..."
        }
      }
    }
  ],
  "variable": [
    {
      "key": "baseUrl",
      "value": "http://localhost:3003"
    }
  ]
}
```

To use this collection:
1. Copy the JSON above
2. Open Postman
3. Click "Import" -> "Raw text"
4. Paste the JSON and click "Import"
5. Set up an environment variable `baseUrl` with your API's base URL

## 🤝 Contributing

Please read our contributing guidelines before submitting pull requests.

## 📝 License

This project is licensed under the MIT License.
