# Uriel - Virtual Office Platform

> Transform remote work with immersive virtual offices that bring teams together

Uriel is a virtual office collaboration platform similar to Gather.town, designed to help distributed teams feel connected through persistent digital spaces. It combines the spontaneity of in-person interaction with the flexibility of remote work.

## 🎯 Vision

Create virtual workspaces where teams can:
- **Naturally interact** through proximity-based communication
- **Collaborate seamlessly** with integrated tools and real-time features  
- **Build culture** through spontaneous conversations and shared spaces
- **Stay productive** with meeting rooms, screen sharing, and workflows

## ✨ Key Features

### 🏢 Virtual Workspaces
- **Custom office layouts** with rooms, meeting areas, and social spaces
- **Proximity-based audio** - walk up and start talking naturally
- **Visual presence** - see who's around and what they're working on
- **Interactive objects** - whiteboards, meeting tables, quiet zones

### 🗣️ Natural Communication  
- **Spatial audio** that gets louder as you get closer
- **Quick conversations** without scheduling meetings
- **Screen sharing** for instant collaboration
- **Meeting rooms** for focused discussions

### 🔗 Seamless Integrations
- **Slack** - sync status and map channels to rooms
- **Google Calendar/Outlook** - auto-create meeting rooms
- **File sharing** - drag and drop documents
- **Productivity tools** - bring your existing workflow

### 📊 Team Insights
- **Engagement analytics** - understand team collaboration patterns
- **Room utilization** - optimize office layouts
- **Meeting efficiency** - track and improve meeting culture

## 🏗️ Architecture

### Backend Stack
- **Go** with Gin framework for high-performance APIs
- **MongoDB** for flexible, scalable data storage
- **WebSockets** for real-time presence and communication
- **Redis** for caching and session management
- **JWT** authentication with role-based authorization

### Key Components
- **User Management** - profiles, presence, authentication
- **Workspace Management** - multi-tenant office spaces
- **Room System** - virtual spaces with objects and interactions
- **Real-time Engine** - WebSocket-based position and voice updates
- **Meeting System** - scheduled and instant meetings
- **Integration Layer** - third-party service connections

## 📁 Project Structure

```
uriel/
├── cmd/uriel/           # Application entry point
├── internal/            # Core application code
│   ├── auth/           # Authentication & authorization
│   ├── config/         # Configuration management
│   ├── database/       # Database operations
│   ├── models/         # Data models
│   └── user/         # User management (legacy name)
├── docs/               # API documentation & design
│   ├── API.md         # Comprehensive API specification
│   ├── SCHEMA.md      # Database schema design
│   └── Discussion.md  # Routes & architecture details
└── README.md          # This file
```

## 🚀 Getting Started

### Prerequisites
- Go 1.24+
- MongoDB 4.4+
- Redis 6.0+

### Quick Start

1. **Clone the repository**
   ```bash
   git clone https://github.com/palSagnik/uriel.git
   cd uriel
   ```

2. **Install dependencies**
   ```bash
   go mod download
   ```

3. **Set up environment variables**
   ```bash
   cp .env.example .env
   # Edit .env with your MongoDB and Redis configurations
   ```

4. **Run the application**
   ```bash
   go run cmd/uriel/main.go
   ```

5. **Access the API**
   ```
   http://localhost:8080/api/v1
   ```

## 📖 Documentation

- **[API Reference](docs/API.md)** - Complete API endpoint documentation
- **[Database Schema](docs/SCHEMA.md)** - MongoDB collection design
- **[Routes & Architecture](docs/Discussion.md)** - Detailed implementation specs

## 🛠️ Development Status

**Current Phase:** API Design & Database Schema ✅

**Completed:**
- ✅ Comprehensive API design
- ✅ Database schema planning
- ✅ Authentication system architecture
- ✅ WebSocket event specifications

**In Progress:**
- 🔄 Backend implementation
- 🔄 Real-time WebSocket layer

**Planned:**
- ⏳ Test-driven development (TDD)
- ⏳ Integration implementations
- ⏳ Frontend client
- ⏳ Mobile applications

## 🤝 Contributing

We welcome contributions! Please read our contributing guidelines and feel free to submit issues and pull requests.

### Development Workflow
1. Follow the existing API design in `docs/API.md`
2. Implement with test-driven development
3. Ensure database operations align with `docs/SCHEMA.md`
4. Add comprehensive tests
5. Update documentation as needed

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- Inspired by [Gather.town](https://gather.town) and the vision of more human remote work
- Built with modern Go practices and scalable architecture patterns
- Designed for the future of distributed team collaboration

---

**Ready to transform your remote work experience?** 🚀

[Get Started](#getting-started) | [View API Docs](docs/API.md) | [Contribute](#contributing)