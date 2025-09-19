# 🎪 Haoma API Complete Scenario Guide

## Complete User Journey: From Registration to 7-Node Completion

This guide walks through the **complete API call scenario** for a user from account creation through completing all 7 carnival nodes.

---

## 🚀 Prerequisites Setup

### 1. Start PostgreSQL Database
```bash
make pg-start
# OR manually:
docker-compose -f docker-compose.dev.yml up -d postgres
```

### 2. Start Haoma Server
```bash
# With database running:
make dev
# OR start complete environment:
make dev-env
```

### 3. Seed Sample Data
```bash
make seed
```

---

## 📋 Complete API Flow

### **Step 1: User Registration**
Create a new player account.

**Request:**
```bash
curl -X POST http://localhost:8080/api/v1/auth/signup \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Rostam Dastan",
    "email": "rostam@haoma.dev",
    "password": "cyber_guardian_2024"
  }'
```

**Response:**
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "name": "Rostam Dastan",
  "email": "rostam@haoma.dev", 
  "message": "🎪 Welcome to Haoma's carnival! Your account has been created."
}
```

---

### **Step 2: User Login (Optional Verification)**
Authenticate the user to verify credentials.

**Request:**
```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "rostam@haoma.dev",
    "password": "cyber_guardian_2024"
  }'
```

**Response:**
```json
{
  "player": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "name": "Rostam Dastan",
    "email": "rostam@haoma.dev"
  },
  "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "token_type": "Bearer",
  "expires_in": 86400,
  "message": "🎪 Welcome back to the carnival!"
}
```

---

### **Step 3: Start Carnival Session**
Create a new session (without starting any node yet).

**Request:**
```bash
curl -X POST http://localhost:8080/api/v1/sessions/start \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
```

**Response:**
```json
{
  "session_id": "123e4567-e89b-12d3-a456-426614174001",
  "message": "🎪 Session created! Scan a node QR code at any carnival location to begin your journey."
}
```

---

### **Step 4: Scan QR Code at First Location**
Player physically visits a carnival location and scans the QR code to start Node 1.

**Request:**
```bash
curl -X POST http://localhost:8080/api/v1/nodes/scan \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..." \
  -d '{
    "node_code": "NODE_001",
    "session_id": "123e4567-e89b-12d3-a456-426614174001"
  }'
```

**Response:**
```json
{
  "session_id": "123e4567-e89b-12d3-a456-426614174001",
  "node": {
    "number": 1,
    "questions": [
      {
        "id": "q1-uuid-here",
        "text": "What is the primary purpose of AES encryption?",
        "option_a": "Data compression", 
        "option_b": "Data protection",
        "option_c": "Data transmission",
        "option_d": "Data validation"
      },
      {
        "id": "q2-uuid-here", 
        "text": "Which port is commonly used for HTTPS?",
        "option_a": "80",
        "option_b": "443", 
        "option_c": "8080",
        "option_d": "3389"
      }
      // ... 3 more questions (5 total per node)
    ]
  },
  "message": "🎪 Welcome to Node 1! Answer all 5 questions to continue your journey."
}
```

---

### **Steps 5-39: Answer All Questions & Scan QR Codes**

**For questions 1-4 in each node:**

**Request:**
```bash
curl -X POST http://localhost:8080/api/v1/sessions/123e4567-e89b-12d3-a456-426614174001/answer \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..." \
  -d '{
    "question_id": "q1-uuid-here",
    "answer": "B"
  }'
```

**Response (Questions 1-4):**
```json
{
  "is_correct": true,
  "node_completed": false,
  "message": "✅ Correct! 3 questions remaining in this node."
}
```

**Response (Question 5 of each node - Node Completion):**
```json
{
  "is_correct": true,
  "node_completed": true,
  "current_score": 580,
  "message": "🎪 Node completed! Check your updated leaderboard position. Find the next location to continue."
}
```

**Note:** Each node completion now automatically updates the leaderboard with current progress, allowing real-time competition tracking.

**After completing each node, player physically moves to next location and scans QR code:**

```bash
# Scan QR at Node 2 location
curl -X POST http://localhost:8080/api/v1/nodes/scan \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..." \
  -d '{
    "node_code": "NODE_002",
    "session_id": "123e4567-e89b-12d3-a456-426614174001"
  }'
```

**After final question of Node 7 (Game Completion):**
```json
{
  "is_correct": true,
  "node_completed": true,
  "current_score": 2847,
  "message": "🎉 Congratulations! You've completed all 7 nodes of the carnival!"
}
```

**Note:** The game automatically completes when all 7 nodes are finished. No separate "finish" endpoint is required.

---

### **Step 40: View Leaderboard**
Check the **taxteh-ye sharaf** (board of honor) to see top performers.

**Request:**
```bash
curl -X GET http://localhost:8080/api/v1/leaderboard
```

**Response:**
```json
{
  "entries": [
    {
      "rank": 1,
      "player_name": "Cyrus The Great", 
      "final_score": 3200,
      "completion_time": "35m12s",
      "achieved_at": "2025-09-18T14:15:22Z"
    },
    {
      "rank": 2,
      "player_name": "Darius CodeBreaker",
      "final_score": 3100,
      "completion_time": "38m45s", 
      "achieved_at": "2025-09-18T13:45:33Z"
    },
    {
      "rank": 3,
      "player_name": "Rostam Dastan",
      "final_score": 2847,
      "completion_time": "47m32s",
      "achieved_at": "2025-09-18T15:23:10Z"
    }
    // ... up to 10 entries
  ]
}
```

---

## 🎯 **API Summary**

| **Step** | **Endpoint** | **Purpose** |
|----------|-------------|-------------|
| 1 | `POST /auth/signup` | Create new player account |
| 2 | `POST /auth/login` | Verify credentials and get access token |
| 3 | `POST /sessions/start` | Create carnival session |
| 4, 11, 18... | `POST /nodes/scan` | Scan QR codes at physical locations |
| 5-39 | `POST /sessions/{id}/answer` | Answer 35 questions (5 per node) |
| 40 | `GET /leaderboard` | View top 10 champions |

**Key Changes:**
- **No session tokens**: Authentication uses JWT access tokens only
- **Session ID in requests**: Session management through direct session_id parameters  
- **Per-node leaderboard**: Leaderboard updates after each node completion
- **Automatic completion**: No separate finish endpoint needed
- **Duplicate prevention**: Each question can only be answered once per session

---

## ⚖️ **Game Rules Recap**

- **7 Categories**: System selects 7 unique categories from 8 available
- **35 Questions Total**: 5 questions per node (4 category + 1 fun)
- **Scoring Formula**: `(correct_answers × 100) - time_penalty`
- **Time Penalty**: 1 point deducted per 30 seconds elapsed
- **Session Limit**: 2 hours maximum
- **PhDT Questions**: Binary YES/NO only (options C/D are null)
- **Leaderboard**: Top 10, ties broken by faster completion time

---

## 🐘 **PostgreSQL Schema**

The enhanced system now uses **PostgreSQL** with these main tables:
- `players` - User accounts with encrypted passwords
- `sessions` - Game sessions with scoring
- `categories` - Question categories (8 total)
- `questions` - Question bank with PhDT support
- `attempts` - Individual answer records  
- `leaderboard_entries` - Top performer records

---

## 🎪 **Ready to Begin?**

```bash
# Start the carnival
make dev-setup

# Visit the mystical endpoints
open http://localhost:8080/docs  # Swagger UI
```

*The carnival awaits your players! Zendeh bāsh!* ✨
