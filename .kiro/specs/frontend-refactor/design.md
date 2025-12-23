# Design Document

## Overview

This design document outlines the architecture and implementation approach for refactoring the Easy Deploy application from a terminal-based UI to a modern web-based frontend. The system will be split into two main components:

1. **Frontend**: A React-based single-page application (SPA) using shadcn/ui components and Tailwind CSS
2. **Backend**: A Go-based REST API server that exposes the existing functionality through HTTP endpoints

The frontend will communicate with the backend via RESTful HTTP requests, maintaining the same data models and business logic while providing an improved user experience through a graphical web interface.

### Technology Stack

**Frontend:**
- React 18+ with TypeScript
- Vite for build tooling and development server
- Tailwind CSS for styling
- shadcn/ui for pre-built accessible components
- React Router for client-side routing
- Axios for HTTP requests
- Lucide React for icons

**Backend:**
- Go 1.21+
- Gin web framework for HTTP routing and middleware
- Existing SSH and configuration logic
- CORS middleware for cross-origin requests
- WebSocket support for streaming deployment logs

## Architecture


### System Architecture Diagram

```
┌─────────────────────────────────────────────────────────────┐
│                         Browser                              │
│  ┌───────────────────────────────────────────────────────┐  │
│  │              React Frontend (SPA)                     │  │
│  │  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐  │  │
│  │  │   Pages     │  │ Components  │  │   Services  │  │  │
│  │  │             │  │             │  │             │  │  │
│  │  │ - Home      │  │ - SSHCard   │  │ - API       │  │  │
│  │  │ - SSH List  │  │ - SSHForm   │  │ - WebSocket │  │  │
│  │  │ - Projects  │  │ - Project   │  │             │  │  │
│  │  │ - Deploy    │  │ - Deploy    │  │             │  │  │
│  │  └─────────────┘  └─────────────┘  └─────────────┘  │  │
│  └───────────────────────────────────────────────────────┘  │
└─────────────────────────────────────────────────────────────┘
                            │
                            │ HTTP/REST + WebSocket
                            ▼
┌─────────────────────────────────────────────────────────────┐
│                    Go Backend Server                         │
│  ┌───────────────────────────────────────────────────────┐  │
│  │                  REST API Layer                       │  │
│  │  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐  │  │
│  │  │   Routes    │  │  Handlers   │  │ Middleware  │  │  │
│  │  │             │  │             │  │             │  │  │
│  │  │ /api/ssh    │  │ SSH CRUD    │  │ CORS        │  │  │
│  │  │ /api/proj   │  │ Project     │  │ Logging     │  │  │
│  │  │ /api/deploy │  │ Deploy      │  │ Error       │  │  │
│  │  │ /ws/logs    │  │ Test        │  │             │  │  │
│  │  └─────────────┘  └─────────────┘  └─────────────┘  │  │
│  └───────────────────────────────────────────────────────┘  │
│  ┌───────────────────────────────────────────────────────┐  │
│  │              Business Logic Layer                     │  │
│  │  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐  │  │
│  │  │ SSH Client  │  │   Config    │  │   Deploy    │  │  │
│  │  │  (existing) │  │  Manager    │  │   Engine    │  │  │
│  │  └─────────────┘  └─────────────┘  └─────────────┘  │  │
│  └───────────────────────────────────────────────────────┘  │
│  ┌───────────────────────────────────────────────────────┐  │
│  │              Data Persistence Layer                   │  │
│  │              (config.json - existing)                 │  │
│  └───────────────────────────────────────────────────────┘  │
└─────────────────────────────────────────────────────────────┘
```

### Communication Flow

1. User interacts with React components in the browser
2. Components call service functions that make HTTP requests
3. Backend API handlers receive requests and validate input
4. Handlers call business logic functions (existing SSH/config code)
5. Business logic performs operations and returns results
6. Handlers format responses and send back to frontend
7. Frontend updates UI based on response data

For real-time deployment logs:
1. Frontend establishes WebSocket connection
2. Backend streams deployment output through WebSocket
3. Frontend displays logs in real-time as they arrive

## Components and Interfaces


### Frontend Components

#### Pages

**HomePage** (`/`)
- Main landing page with navigation to SSH configs and projects
- Displays application title and quick action buttons
- Uses shadcn Card components for visual sections

**SSHListPage** (`/ssh`)
- Displays all SSH configurations in a grid layout
- Each configuration shown as a shadcn Card with actions
- Includes "Add New" button to create configurations
- Empty state when no configurations exist

**SSHFormPage** (`/ssh/new` and `/ssh/:id/edit`)
- Form for creating/editing SSH configurations
- Uses shadcn Form components with validation
- Conditional fields based on authentication type
- Cancel and Save buttons

**ProjectListPage** (`/projects`)
- Displays all projects in a grid layout
- Each project shown as a shadcn Card with server badges
- Includes "Add New" button to create projects
- Empty state when no projects exist

**ProjectFormPage** (`/projects/new` and `/projects/:id/edit`)
- Form for creating/editing projects
- Multi-line text areas for build/deploy scripts
- Multi-select for deploy servers
- Cancel and Save buttons

**DeployPage** (`/deploy/:projectId`)
- Displays deployment progress and logs
- Real-time log streaming via WebSocket
- Shows deployment status (running, success, failed)
- Cancel deployment button

#### Reusable Components

**SSHCard**
- Props: `sshConfig`, `onEdit`, `onDelete`, `onTest`
- Displays SSH configuration details in a card
- Action buttons for edit, delete, test
- Status indicator badge

**ProjectCard**
- Props: `project`, `onEdit`, `onDelete`, `onDeploy`
- Displays project details in a card
- Shows assigned server badges
- Action buttons for edit, delete, deploy

**SSHTestDialog**
- Props: `sshConfig`, `isOpen`, `onClose`
- Modal dialog for testing SSH connections
- Shows loading spinner during test
- Displays success/error results

**DeleteConfirmDialog**
- Props: `title`, `message`, `isOpen`, `onConfirm`, `onCancel`
- Reusable confirmation dialog
- Uses shadcn AlertDialog component

**DeployLogViewer**
- Props: `logs`, `status`
- Displays deployment logs with syntax highlighting
- Auto-scrolls to bottom as new logs arrive
- Color-coded status indicators

**Layout**
- Props: `children`
- Main application layout with navigation
- Sidebar or top navigation bar
- Consistent header and footer


### Backend API Endpoints

#### SSH Configuration Endpoints

**GET /api/ssh**
- Returns list of all SSH configurations
- Response: `{ configs: SSHConfig[] }`

**GET /api/ssh/:id**
- Returns a single SSH configuration by name
- Response: `{ config: SSHConfig }`

**POST /api/ssh**
- Creates a new SSH configuration
- Request body: `SSHConfig`
- Response: `{ config: SSHConfig, message: string }`

**PUT /api/ssh/:id**
- Updates an existing SSH configuration
- Request body: `SSHConfig`
- Response: `{ config: SSHConfig, message: string }`

**DELETE /api/ssh/:id**
- Deletes an SSH configuration
- Response: `{ message: string }`

**POST /api/ssh/:id/test**
- Tests SSH connection
- Response: `{ success: boolean, message: string, output?: string }`

#### Project Endpoints

**GET /api/projects**
- Returns list of all projects
- Response: `{ projects: Project[] }`

**GET /api/projects/:id**
- Returns a single project by name
- Response: `{ project: Project }`

**POST /api/projects**
- Creates a new project
- Request body: `Project`
- Response: `{ project: Project, message: string }`

**PUT /api/projects/:id**
- Updates an existing project
- Request body: `Project`
- Response: `{ project: Project, message: string }`

**DELETE /api/projects/:id**
- Deletes a project
- Response: `{ message: string }`

**POST /api/projects/:id/deploy**
- Initiates deployment for a project
- Response: `{ deploymentId: string, message: string }`

#### WebSocket Endpoints

**WS /ws/deploy/:deploymentId**
- Streams deployment logs in real-time
- Messages: `{ type: "log" | "status" | "error", data: string }`

## Data Models


### TypeScript Interfaces (Frontend)

```typescript
interface SSHConfig {
  name: string;
  host: string;
  port: number;
  user: string;
  authType: 'password' | 'key' | 'agent';
  password?: string;
  keyFile?: string;
  keyPass?: string;
}

interface Project {
  name: string;
  buildInstructions: string;
  deployScript: string;
  deployServers: string[];
  createdAt: string;
  updatedAt: string;
}

interface APIResponse<T> {
  data?: T;
  message?: string;
  error?: string;
}

interface DeploymentLog {
  type: 'log' | 'status' | 'error';
  data: string;
  timestamp: string;
}

interface DeploymentStatus {
  id: string;
  projectName: string;
  status: 'pending' | 'running' | 'success' | 'failed';
  startedAt: string;
  completedAt?: string;
}
```

### Go Structs (Backend)

The existing structs in `types.go` will be reused:
- `SSHConfig`
- `Project`
- `Config`

New structs for API responses:

```go
type APIResponse struct {
    Data    interface{} `json:"data,omitempty"`
    Message string      `json:"message,omitempty"`
    Error   string      `json:"error,omitempty"`
}

type DeploymentLog struct {
    Type      string    `json:"type"`
    Data      string    `json:"data"`
    Timestamp time.Time `json:"timestamp"`
}

type DeploymentStatus struct {
    ID          string    `json:"id"`
    ProjectName string    `json:"projectName"`
    Status      string    `json:"status"`
    StartedAt   time.Time `json:"startedAt"`
    CompletedAt *time.Time `json:"completedAt,omitempty"`
}
```


## Correctness Properties

*A property is a characteristic or behavior that should hold true across all valid executions of a system-essentially, a formal statement about what the system should do. Properties serve as the bridge between human-readable specifications and machine-verifiable correctness guarantees.*

### Frontend Properties

**Property 1: Data fetching on mount**
*For any* page component that displays data from the backend, mounting the component should trigger an API call to fetch that data.
**Validates: Requirements 1.2, 2.1, 7.1**

**Property 2: Error handling for network failures**
*For any* API call that fails due to network error, the frontend should display an error message to the user.
**Validates: Requirements 1.4**

**Property 3: Required field display for SSH configs**
*For any* SSH configuration rendered in the UI, the display should include name, host, port, user, and authentication type fields.
**Validates: Requirements 2.2**

**Property 4: Status indicators for SSH configs**
*For any* SSH configuration displayed in a list, a status indicator should be rendered.
**Validates: Requirements 2.5**

**Property 5: Form validation before submission**
*For any* form with invalid required fields, submitting the form should not trigger an API call and should display validation errors.
**Validates: Requirements 3.2, 3.4**

**Property 6: Conditional field display by auth type**
*For any* SSH configuration form, when an authentication type is selected, only the fields relevant to that auth type should be visible.
**Validates: Requirements 3.5**

**Property 7: Form pre-population for editing**
*For any* SSH configuration being edited, the edit form should be initialized with all current values from that configuration.
**Validates: Requirements 4.1**

**Property 8: Real-time validation on field change**
*For any* form field that is modified, validation should run and display errors if the new value is invalid.
**Validates: Requirements 4.2**

**Property 9: API call on save**
*For any* valid form submission, the frontend should make an API call with the form data.
**Validates: Requirements 3.3, 4.3**

**Property 10: UI update after successful operation**
*For any* successful API response (create, update, delete), the frontend should update the relevant list view and display a success message.
**Validates: Requirements 3.3, 4.4, 5.3**

**Property 11: Confirmation dialog before deletion**
*For any* delete action initiated by the user, a confirmation dialog should be displayed before making the delete API call.
**Validates: Requirements 5.1, 7.5**

**Property 12: Warning for configs in use**
*For any* SSH configuration that is referenced by one or more projects, attempting to delete it should display a warning message.
**Validates: Requirements 5.5**

**Property 13: Loading indicator during async operations**
*For any* asynchronous API call in progress, a loading indicator should be displayed to the user.
**Validates: Requirements 6.2**

**Property 14: Test result display**
*For any* SSH connection test result (success or failure), the appropriate message and details should be displayed to the user.
**Validates: Requirements 6.3, 6.4**

**Property 15: Required field display for projects**
*For any* project rendered in the UI, the display should include name, build instructions preview, deploy script preview, and assigned servers.
**Validates: Requirements 7.2**

**Property 16: Real-time log streaming**
*For any* deployment in progress, logs received via WebSocket should be displayed in the UI as they arrive.
**Validates: Requirements 8.2**

**Property 17: Deployment status display**
*For any* completed deployment (success or failure), the appropriate status message and details should be displayed.
**Validates: Requirements 8.3, 8.4**


### Backend Properties

**Property 18: HTTP status codes for errors**
*For any* API endpoint that encounters an error, the response should include an appropriate HTTP status code (4xx for client errors, 5xx for server errors).
**Validates: Requirements 9.5**

**Property 19: CORS headers on preflight**
*For any* OPTIONS preflight request, the backend should respond with appropriate CORS headers including Access-Control-Allow-Origin, Access-Control-Allow-Methods, and Access-Control-Allow-Headers.
**Validates: Requirements 10.2, 10.4**

## Error Handling

### Frontend Error Handling

**Network Errors:**
- All API calls wrapped in try-catch blocks
- Display user-friendly error messages using shadcn Toast component
- Retry mechanism for transient failures
- Offline detection and appropriate messaging

**Validation Errors:**
- Client-side validation before API calls
- Display inline error messages on form fields
- Prevent form submission when validation fails
- Clear error messages on field correction

**API Error Responses:**
- Parse error messages from backend responses
- Display specific error details to users
- Log errors to console for debugging
- Graceful degradation when features unavailable

**WebSocket Errors:**
- Handle connection failures with reconnection logic
- Display connection status to users
- Buffer logs during disconnection
- Fallback to polling if WebSocket unavailable

### Backend Error Handling

**Request Validation:**
- Validate all incoming request data
- Return 400 Bad Request for invalid input
- Include detailed error messages in response
- Log validation failures for monitoring

**SSH Operation Errors:**
- Catch SSH connection failures
- Return appropriate error codes (500 for server errors)
- Include error details in response
- Log SSH errors with context

**File System Errors:**
- Handle config.json read/write failures
- Return 500 Internal Server Error
- Log file system errors
- Implement retry logic for transient failures

**CORS Errors:**
- Log rejected CORS requests
- Return appropriate headers
- Configure allowed origins properly
- Handle preflight requests correctly


## Testing Strategy

### Frontend Testing

**Unit Testing:**
- Test framework: Vitest
- Component testing: React Testing Library
- Test individual components in isolation
- Mock API calls using MSW (Mock Service Worker)
- Test form validation logic
- Test utility functions and hooks
- Coverage target: 80% for critical paths

**Property-Based Testing:**
- Library: fast-check (JavaScript property-based testing)
- Minimum iterations: 100 per property test
- Each property test tagged with: `Feature: frontend-refactor, Property {number}: {property_text}`
- Generate random valid/invalid form data
- Test API error handling with various error responses
- Test component rendering with various data shapes

**Integration Testing:**
- Test complete user flows (add SSH config, deploy project)
- Test API service layer with real HTTP calls to mock server
- Test WebSocket connection and message handling
- Test routing and navigation

**Example Unit Tests:**
- SSHCard renders with correct data
- SSHForm validates required fields
- DeleteConfirmDialog shows correct message
- DeployLogViewer displays logs correctly

### Backend Testing

**Unit Testing:**
- Test framework: Go testing package
- Test API handlers with mock requests
- Test business logic functions
- Test configuration loading/saving
- Test SSH client operations with mock connections
- Coverage target: 80% for critical paths

**Property-Based Testing:**
- Library: gopter (Go property-based testing)
- Minimum iterations: 100 per property test
- Each property test tagged with: `Feature: frontend-refactor, Property {number}: {property_text}`
- Generate random valid/invalid API requests
- Test error handling with various error conditions
- Test CORS headers with various origins

**Integration Testing:**
- Test complete API flows with real HTTP server
- Test WebSocket streaming
- Test file system operations with temp files
- Test SSH operations with mock SSH server

**Example Unit Tests:**
- GET /api/ssh returns all configurations
- POST /api/ssh creates new configuration
- DELETE /api/ssh/:id removes configuration
- POST /api/ssh/:id/test validates connection

### End-to-End Testing

**Tool:** Playwright
- Test complete user workflows in real browser
- Test frontend-backend integration
- Test WebSocket communication
- Test error scenarios
- Run against local development environment

**Example E2E Tests:**
- User can add, edit, and delete SSH configuration
- User can create project and deploy it
- User sees real-time deployment logs
- User sees appropriate error messages on failures

## Implementation Notes

### Frontend Setup

1. Initialize React project with Vite and TypeScript
2. Install and configure Tailwind CSS
3. Install shadcn/ui CLI and add required components
4. Set up React Router for navigation
5. Configure Axios for API calls
6. Set up WebSocket client for log streaming
7. Create folder structure: src/{pages, components, lib, hooks, services, types}

### Backend Refactoring

1. Install Gin web framework
2. Create API router and handler functions
3. Add CORS middleware
4. Implement WebSocket handler for log streaming
5. Refactor existing TUI code into reusable business logic
6. Add API response types
7. Update main.go to start HTTP server instead of TUI

### Development Workflow

1. Backend runs on http://localhost:8080
2. Frontend dev server runs on http://localhost:5173
3. Frontend proxies API requests to backend during development
4. Use environment variables for API URL configuration
5. Hot reload enabled for both frontend and backend

### Deployment Considerations

1. Build frontend as static files
2. Serve frontend files from Go backend using embed
3. Single binary deployment with embedded frontend
4. Configure production API URL
5. Enable HTTPS in production
6. Set appropriate CORS origins for production

### Migration Path

1. Keep existing TUI code in separate files
2. Add new API layer alongside TUI
3. Test API thoroughly before removing TUI
4. Provide command-line flag to choose TUI or web mode
5. Eventually deprecate TUI after web UI is stable
