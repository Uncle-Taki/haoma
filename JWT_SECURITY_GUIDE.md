# 🔐 Haoma JWT Security Implementation

## Security Enhancement: JWT Token-Based Authentication

This document explains the **JWT security implementation** that protects Haoma's API from unauthorized access.

---

## 🚨 **Security Threats Addressed**

### **Before JWT Implementation:**
❌ **No Authentication** - Anyone could use any `player_id`  
❌ **Session Hijacking** - Sessions accessible without ownership verification  
❌ **Data Exposure** - Player data accessible to unauthorized users  
❌ **API Abuse** - No rate limiting or access control  

### **After JWT Implementation:**  
✅ **Player Authentication** - JWT tokens verify user identity  
✅ **Session Security** - Session tokens tied to authenticated players  
✅ **Access Control** - Middleware validates all protected endpoints  
✅ **Token Expiration** - Automatic timeout prevents long-term abuse  

---

## 🔑 **JWT Token Types**

### **Player Access Token**
- **Purpose**: Authenticates the player for all API calls
- **Duration**: 24 hours
- **Contains**: `player_id`, `player_name`, `player_email`
- **Usage**: `Authorization: Bearer <access_token>`

**Note**: Session management is now handled through direct `session_id` parameters in API calls, eliminating the need for separate session tokens.

---

## 🛡️ **Security Flow**

### **Registration & Login:**
```bash
# 1. Register (public)
POST /auth/signup
# Response: Welcome message

# 2. Login (public) 
POST /auth/login
# Response: access_token + player info
{
  "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "token_type": "Bearer",
  "expires_in": 86400,
  "player": { "id": "...", "name": "...", "email": "..." }
}
```

### **Protected Endpoints:**
```bash
# 3. Start Session (requires JWT)
POST /sessions/start
Authorization: Bearer <access_token>
# Response: session_id

# 4. Scan Node QR (requires JWT)
POST /nodes/scan  
Authorization: Bearer <access_token>
Body: {"node_code": "NODE_001", "session_id": "uuid"}  # session_id optional for first scan

# 5. Answer Questions (requires JWT)
POST /sessions/{session_id}/answer
Authorization: Bearer <access_token>
Body: {"question_id": "uuid", "answer": "B"}
```

---

## 🔧 **Implementation Details**

### **JWT Service (`internal/infrastructure/auth/jwt.go`):**
- **Token Generation**: Creates signed JWT tokens with claims
- **Token Validation**: Verifies signatures and expiration
- **Refresh Capability**: Extends token lifetime when needed

### **Middleware (`internal/infrastructure/auth/middleware.go`):**
- **JWTMiddleware**: Validates player access tokens and extracts player information

### **Route Protection (`internal/adapters/http/handlers.go`):**
```go
// Public routes (no authentication)
authPublic.POST("/signup", handler.Signup)
authPublic.POST("/login", handler.Login) 
api.GET("/leaderboard", handler.GetLeaderboard)

// Protected routes (JWT required)
authProtected.GET("/profile", handler.GetProfile)
sessions.POST("/start", handler.StartSession)
nodes.POST("/scan", handler.ScanNodeQR)
sessions.POST("/:id/answer", handler.SubmitAnswer)
```

---

## 🔒 **Environment Security**

### **JWT Secret Configuration:**
```bash
# .env file (NEVER commit this!)
JWT_SECRET=your_super_secret_jwt_key_change_in_production_at_least_32_characters

# Production requirements:
# - Minimum 32 characters
# - Random, cryptographically secure
# - Different for each environment
# - Stored in secure key management
```

### **Environment Variables:**
```bash
# Database (keep secure)
DB_PASSWORD=strong_database_password

# JWT (critical security)
JWT_SECRET=use_a_proper_secret_key_manager_in_production

# Server
GIN_MODE=release  # In production
```

---

## 📱 **Client Integration**

### **Mobile/Web App Flow:**
```javascript
// 1. Login and store tokens
const loginResponse = await fetch('/api/v1/auth/login', {
  method: 'POST',
  body: JSON.stringify({ email, password }),
  headers: { 'Content-Type': 'application/json' }
});

const { access_token } = await loginResponse.json();
localStorage.setItem('access_token', access_token);

// 2. Use access token for all requests
const headers = {
  'Authorization': `Bearer ${localStorage.getItem('access_token')}`,
  'Content-Type': 'application/json'
};

// 3. Start session and get session_id
const sessionResponse = await fetch('/api/v1/sessions/start', {
  method: 'POST',
  headers
});

const { session_id } = await sessionResponse.json();
localStorage.setItem('session_id', session_id);

// 4. Scan QR code with session_id
await fetch('/api/v1/nodes/scan', {
  method: 'POST',
  headers,
  body: JSON.stringify({ 
    node_code: 'NODE_CRYPTO_001',
    session_id: localStorage.getItem('session_id')
  })
});

// 5. Answer questions with session_id in URL
await fetch(`/api/v1/sessions/${localStorage.getItem('session_id')}/answer`, {
  method: 'POST',
  headers,
  body: JSON.stringify({ question_id: 'uuid', answer: 'B' })
});
```

---

## 🛠️ **Testing with cURL**

### **Complete Secure API Flow:**
```bash
# 1. Register
curl -X POST http://localhost:8080/api/v1/auth/signup \
  -H "Content-Type: application/json" \
  -d '{"name":"Rostam","email":"rostam@haoma.dev","password":"test123"}'

# 2. Login and extract token
RESPONSE=$(curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"rostam@haoma.dev","password":"test123"}')

ACCESS_TOKEN=$(echo $RESPONSE | jq -r '.access_token')

# 3. Start session with JWT
SESSION_RESPONSE=$(curl -X POST http://localhost:8080/api/v1/sessions/start \
  -H "Authorization: Bearer $ACCESS_TOKEN" \
  -H "Content-Type: application/json")

SESSION_ID=$(echo $SESSION_RESPONSE | jq -r '.session_id')

# 4. Scan QR with session_id
curl -X POST http://localhost:8080/api/v1/nodes/scan \
  -H "Authorization: Bearer $ACCESS_TOKEN" \
  -H "Content-Type: application/json" \
  -d "{\"node_code\":\"NODE_001\",\"session_id\":\"$SESSION_ID\"}"

# 5. Answer question with session_id in URL
curl -X POST http://localhost:8080/api/v1/sessions/$SESSION_ID/answer \
  -H "Authorization: Bearer $ACCESS_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"question_id":"question_uuid","answer":"B"}'
```

---

## 🔍 **Security Best Practices**

### **Production Deployment:**
1. **Use HTTPS only** - Never transmit JWT over HTTP
2. **Secure secret storage** - Use AWS Secrets Manager, HashiCorp Vault, etc.
3. **Token rotation** - Implement refresh token mechanism
4. **Rate limiting** - Prevent brute force attacks
5. **Audit logging** - Log authentication events
6. **CORS configuration** - Restrict origins in production

### **Token Management:**
1. **Short expiration** - Access tokens expire in 24h
2. **Secure storage** - Use httpOnly cookies or secure local storage
3. **Automatic logout** - Clear tokens on expiration
4. **Logout endpoint** - Invalidate tokens server-side
5. **Session ID handling** - Session IDs are managed in API request parameters

---

## 🎪 **Security Benefits for ELECOMP 1404**

✅ **Student Data Protection** - Each student can only access their own data  
✅ **Session Isolation** - Students cannot interfere with others' sessions  
✅ **Fair Competition** - Prevents cheating through API manipulation  
✅ **Academic Integrity** - Ensures legitimate participation only  
✅ **Professional Standards** - Teaches real-world security practices  

**Your carnival is now enterprise-secure!** 🔐✨
