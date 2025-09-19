# 🎪 Haoma QR Code Setup Guide

## Physical Carnival Setup with QR Codes

This guide explains how to set up the **physical carnival locations** with **QR codes** for the location-based cyber security quiz.

---

## 🗺️ **7 Carnival Node Locations**

Each carnival node should be set up at a **different physical location** with its **unique QR code**:

### **Node 1: Cryptography Station**
- **QR Code Content**: `NODE_CRYPTO_001`
- **Location**: Computer lab or cryptography display area
- **Theme**: The art of secret codes

### **Node 2: Network Security Station** 
- **QR Code Content**: `NODE_NETWORK_002`
- **Location**: Network equipment area or server room
- **Theme**: Protecting digital pathways

### **Node 3: Web Security Station**
- **QR Code Content**: `NODE_WEB_003` 
- **Location**: Web development area
- **Theme**: Guarding web applications

### **Node 4: Malware Analysis Station**
- **QR Code Content**: `NODE_MALWARE_004`
- **Location**: Cybersecurity lab
- **Theme**: Understanding digital threats

### **Node 5: Social Engineering Station**
- **QR Code Content**: `NODE_SOCIAL_005`
- **Location**: Communication/presentation area
- **Theme**: Human vulnerability exploitation

### **Node 6: Incident Response Station**
- **QR Code Content**: `NODE_INCIDENT_006`
- **Location**: Emergency response/operations center
- **Theme**: Handling security breaches

### **Node 7: Phishing Detection (PhDT) Station**
- **QR Code Content**: `NODE_PHDT_007`
- **Location**: Email security demonstration area
- **Theme**: Binary YES/NO phishing detection

---

## 📱 **QR Code Generation**

Generate QR codes with these exact texts:

```bash
# QR Code Contents (case-sensitive)
NODE_CRYPTO_001
NODE_NETWORK_002
NODE_WEB_003
NODE_MALWARE_004
NODE_SOCIAL_005
NODE_INCIDENT_006
NODE_PHDT_007
```

**Online QR Generator**: https://qr-code-generator.com/
**Command Line**: `qrencode -o node1.png "NODE_CRYPTO_001"`

---

## 🎯 **Game Flow Overview**

### **Student Journey:**
1. **Register** → Create account via mobile app/web
2. **Start Session** → Initialize carnival session
3. **Find Node 1** → Locate first physical location
4. **Scan QR Code** → `POST /nodes/scan` with `NODE_CRYPTO_001`
5. **Answer 5 Questions** → Complete cryptography challenges
6. **Move to Node 2** → Find next physical location  
7. **Scan Next QR** → `POST /nodes/scan` with `NODE_NETWORK_002`
8. **Repeat** → Until all 7 nodes completed
9. **View Results** → Final score and leaderboard position

### **Key Benefits:**
- ✅ **Physical Movement** - Students explore the campus/facility
- ✅ **Location-Based Learning** - Each station themed appropriately  
- ✅ **Social Interaction** - Students meet at different locations
- ✅ **Prevents Cheating** - Must physically visit each location
- ✅ **Gamification** - Treasure hunt style experience

---

## 🏗️ **Physical Setup Instructions**

### **For Each Node Location:**

1. **Print QR Code** (minimum 3x3 inches for easy scanning)
2. **Create Station Sign** with:
   - Node number and theme name
   - QR code prominently displayed  
   - Brief instructions: "Scan to start Node X"
   - Carnival/Persian themed decorations

3. **Station Layout**:
   ```
   ╔════════════════════════╗
   ║    🎪 NODE 1: CRYPTO   ║  
   ║                        ║
   ║    [QR CODE HERE]      ║
   ║   NODE_CRYPTO_001      ║
   ║                        ║
   ║ "Scan to unlock the    ║
   ║  mysteries of secret   ║
   ║  codes and ciphers!"   ║
   ╚════════════════════════╝
   ```

4. **Optional Enhancements**:
   - Theme-related props/displays
   - Persian/carnival decorations  
   - Directional signs to next location
   - Staff/volunteers for assistance

---

## 📋 **Location Planning Template**

| **Node** | **QR Code** | **Suggested Location** | **Setup Requirements** |
|----------|-------------|------------------------|----------------------|
| 1: Crypto | `NODE_CRYPTO_001` | Computer Lab | Encryption books/posters |
| 2: Network | `NODE_NETWORK_002` | Server Room/IT Area | Network equipment display |
| 3: Web Sec | `NODE_WEB_003` | Web Dev Lab | Browser security demos |
| 4: Malware | `NODE_MALWARE_004` | Security Lab | Virus/threat examples |
| 5: Social | `NODE_SOCIAL_005` | Common Area | Social engineering examples |
| 6: Incident | `NODE_INCIDENT_006` | Operations Center | Incident response procedures |
| 7: PhDT | `NODE_PHDT_007` | Email/Comm Area | Phishing email examples |

---

## 🔧 **Technical Configuration**

The QR codes map to these categories in the system:

```go
// QR Code Mapping (in carnival.go)
nodeMapping := map[string]struct {
    number   int
    category string
}{
    "NODE_CRYPTO_001":    {1, "Cryptography"},
    "NODE_NETWORK_002":   {2, "Network Security"},
    "NODE_WEB_003":       {3, "Web Security"}, 
    "NODE_MALWARE_004":   {4, "Malware Analysis"},
    "NODE_SOCIAL_005":    {5, "Social Engineering"},
    "NODE_INCIDENT_006":  {6, "Incident Response"},
    "NODE_PHDT_007":      {7, "PhDT"},
}
```

**To Add More Nodes**: Simply add new mappings and generate corresponding QR codes.

---

## 📱 **Mobile App Integration**

Students can use any **QR scanner app** or the **camera app** on their phones to scan codes, then:

1. **Copy the QR content** (e.g., `NODE_CRYPTO_001`)
2. **Open the Haoma web app** or mobile interface
3. **Paste/enter the node code** in the scan interface
4. **Start answering questions** for that node

**Alternative**: Build a **native mobile app** with integrated QR scanner that directly calls the API.

---

## 🎪 **Ready to Set Up Your Carnival!**

*Transform your learning space into Haoma's mystical carnival where each location becomes a cyber-security trial. Let the journey begin!*

**Zendeh bāsh!** ✨
