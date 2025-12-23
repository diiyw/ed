# Implementation Plan

- [x] 1. Set up frontend project structure


  - Initialize Vite + React + TypeScript project in `frontend/` directory
  - Install and configure Tailwind CSS
  - Install shadcn/ui CLI and initialize
  - Set up folder structure: src/{pages, components, lib, hooks, services, types}
  - Configure path aliases in tsconfig.json
  - _Requirements: 12.1, 12.2, 12.3_

- [ ] 1.1 Write unit tests for project setup
  - Test that Vite config is valid
  - Test that Tailwind is properly configured
  - _Requirements: 12.1_

- [x] 2. Install frontend dependencies


  - Install React Router for routing
  - Install Axios for HTTP requests
  - Install Lucide React for icons
  - Install date-fns for date formatting
  - Install necessary shadcn/ui components (Button, Card, Form, Dialog, Toast, etc.)
  - _Requirements: 1.1, 11.3_

- [x] 3. Create TypeScript type definitions


  - Create types/api.ts with SSHConfig, Project, APIResponse interfaces
  - Create types/deployment.ts with DeploymentLog, DeploymentStatus interfaces
  - Export all types from types/index.ts
  - _Requirements: 12.3_

- [-] 4. Implement API service layer

  - Create services/api.ts with Axios instance configuration
  - Implement SSH CRUD functions (getSSHConfigs, createSSHConfig, updateSSHConfig, deleteSSHConfig, testSSHConnection)
  - Implement Project CRUD functions (getProjects, createProject, updateProject, deleteProject, deployProject)
  - Add error handling and response parsing
  - _Requirements: 1.2, 9.1, 9.2, 12.2_

- [ ] 4.1 Write property test for API error handling
  - **Property 2: Error handling for network failures**
  - **Validates: Requirements 1.4**

- [x] 5. Implement WebSocket service


  - Create services/websocket.ts for deployment log streaming
  - Implement connection, message handling, and reconnection logic
  - Add event emitter for log messages
  - _Requirements: 8.2, 9.4_

- [ ] 5.1 Write unit tests for WebSocket service
  - Test connection establishment
  - Test message parsing
  - Test reconnection logic
  - _Requirements: 8.2_

- [-] 6. Create reusable UI components

  - Create components/Layout.tsx with navigation
  - Create components/SSHCard.tsx for displaying SSH configs
  - Create components/ProjectCard.tsx for displaying projects
  - Create components/SSHTestDialog.tsx for connection testing
  - Create components/DeleteConfirmDialog.tsx for confirmations
  - Create components/DeployLogViewer.tsx for deployment logs
  - _Requirements: 2.1, 7.1, 11.3_

- [ ] 6.1 Write property test for required field display
  - **Property 3: Required field display for SSH configs**
  - **Validates: Requirements 2.2**

- [ ] 6.2 Write property test for status indicators
  - **Property 4: Status indicators for SSH configs**
  - **Validates: Requirements 2.5**

- [ ] 6.3 Write property test for project field display
  - **Property 15: Required field display for projects**
  - **Validates: Requirements 7.2**

- [ ] 6.4 Write unit tests for UI components
  - Test SSHCard rendering
  - Test ProjectCard rendering
  - Test dialog components
  - _Requirements: 2.1, 7.1_

- [x] 7. Create SSH configuration form component


  - Create components/SSHForm.tsx with all fields
  - Implement form validation using React Hook Form
  - Add conditional field display based on auth type
  - Implement real-time validation
  - _Requirements: 3.1, 3.2, 3.4, 3.5, 4.2_

- [ ] 7.1 Write property test for form validation
  - **Property 5: Form validation before submission**
  - **Validates: Requirements 3.2, 3.4**

- [ ] 7.2 Write property test for conditional fields
  - **Property 6: Conditional field display by auth type**
  - **Validates: Requirements 3.5**

- [ ] 7.3 Write property test for real-time validation
  - **Property 8: Real-time validation on field change**
  - **Validates: Requirements 4.2**

- [ ] 7.4 Write unit tests for SSH form
  - Test form rendering
  - Test validation logic
  - Test conditional field display
  - _Requirements: 3.1, 3.5_

- [x] 8. Create project form component





  - Create components/ProjectForm.tsx with all fields
  - Implement form validation
  - Add multi-select for deploy servers
  - Add text areas for build/deploy scripts
  - _Requirements: 7.4_

- [x] 8.1 Write unit tests for project form

  - Test form rendering
  - Test validation logic
  - Test server selection
  - _Requirements: 7.4_


- [x] 9. Implement HomePage





  - Create pages/HomePage.tsx with main navigation
  - Add welcome message and quick action buttons
  - Use shadcn Card components for layout
  - _Requirements: 1.1_

- [x] 9.1 Write unit test for HomePage


  - Test that HomePage renders without errors
  - Test navigation links are present
  - _Requirements: 1.1_

- [x] 10. Implement SSHListPage






  - Create pages/SSHListPage.tsx
  - Fetch SSH configs on mount
  - Display configs in grid using SSHCard components
  - Add "Add New" button
  - Implement empty state
  - Handle edit, delete, and test actions
  - _Requirements: 2.1, 2.2, 2.3, 2.5_

- [x] 10.1 Write property test for data fetching on mount


  - **Property 1: Data fetching on mount**
  - **Validates: Requirements 1.2, 2.1**

- [x] 10.2 Write unit tests for SSHListPage

  - Test empty state display
  - Test config list rendering
  - Test action buttons
  - _Requirements: 2.1, 2.3_


- [x] 11. Implement SSHFormPage





  - Create pages/SSHFormPage.tsx for add/edit
  - Use SSHForm component
  - Handle form submission
  - Navigate back on success or cancel
  - Display success/error messages
  - _Requirements: 3.1, 3.3, 4.1, 4.3, 4.4, 4.5_

- [x] 11.1 Write property test for form pre-population

  - **Property 7: Form pre-population for editing**
  - **Validates: Requirements 4.1**

- [x] 11.2 Write property test for API call on save

  - **Property 9: API call on save**
  - **Validates: Requirements 3.3, 4.3**

- [x] 11.3 Write property test for UI update after success

  - **Property 10: UI update after successful operation**
  - **Validates: Requirements 3.3, 4.4, 5.3**

- [x] 11.4 Write unit tests for SSHFormPage

  - Test form submission flow
  - Test navigation on cancel
  - Test success message display
  - _Requirements: 3.3, 4.4_


- [x] 12. Implement delete confirmation functionality



  - Add delete confirmation dialog to SSHListPage
  - Check if SSH config is in use by projects
  - Display warning if in use
  - Handle delete API call on confirmation
  - _Requirements: 5.1, 5.2, 5.3, 5.4, 5.5_

- [x] 12.1 Write property test for confirmation dialog


  - **Property 11: Confirmation dialog before deletion**
  - **Validates: Requirements 5.1, 7.5**

- [x] 12.2 Write property test for in-use warning



  - **Property 12: Warning for configs in use**
  - **Validates: Requirements 5.5**

- [x] 12.3 Write unit tests for delete functionality

  - Test confirmation dialog display
  - Test delete API call
  - Test cancel behavior
  - _Requirements: 5.1, 5.4_

- [x] 13. Implement SSH connection testing




  - Add test connection functionality to SSHListPage
  - Display SSHTestDialog with loading state
  - Show test results (success/failure)
  - Allow dismissing results
  - _Requirements: 6.1, 6.2, 6.3, 6.4, 6.5_

- [x] 13.1 Write property test for loading indicator


  - **Property 13: Loading indicator during async operations**
  - **Validates: Requirements 6.2**

- [x] 13.2 Write property test for test result display


  - **Property 14: Test result display**
  - **Validates: Requirements 6.3, 6.4**

- [x] 13.3 Write unit tests for connection testing

  - Test loading state display
  - Test success result display
  - Test error result display
  - _Requirements: 6.2, 6.3, 6.4_

- [x] 14. Implement ProjectListPage






  - Create pages/ProjectListPage.tsx
  - Fetch projects on mount
  - Display projects in grid using ProjectCard components
  - Add "Add New" button
  - Implement empty state
  - Handle edit, delete, and deploy actions
  - _Requirements: 7.1, 7.2, 7.3, 7.5_

- [x] 14.1 Write unit tests for ProjectListPage


  - Test empty state display
  - Test project list rendering
  - Test action buttons
  - _Requirements: 7.1, 7.3_

- [-] 15. Implement ProjectFormPage

  - Create pages/ProjectFormPage.tsx for add/edit
  - Use ProjectForm component
  - Handle form submission
  - Navigate back on success or cancel
  - Display success/error messages
  - _Requirements: 7.4_

- [ ] 15.1 Write unit tests for ProjectFormPage
  - Test form submission flow
  - Test navigation on cancel
  - Test success message display
  - _Requirements: 7.4_

- [x] 16. Implement DeployPage


  - Create pages/DeployPage.tsx
  - Establish WebSocket connection for logs
  - Display real-time deployment logs using DeployLogViewer
  - Show deployment status (running, success, failed)
  - Add cancel deployment button
  - Handle WebSocket errors and reconnection
  - _Requirements: 8.1, 8.2, 8.3, 8.4, 8.5_

- [ ] 16.1 Write property test for log streaming
  - **Property 16: Real-time log streaming**
  - **Validates: Requirements 8.2**

- [ ] 16.2 Write property test for deployment status
  - **Property 17: Deployment status display**
  - **Validates: Requirements 8.3, 8.4**

- [ ] 16.3 Write unit tests for DeployPage
  - Test WebSocket connection
  - Test log display
  - Test status display
  - Test cancel button
  - _Requirements: 8.2, 8.3, 8.4, 8.5_

- [x] 17. Set up React Router


  - Create App.tsx with router configuration
  - Define routes for all pages
  - Add 404 page
  - Implement navigation in Layout component
  - _Requirements: 1.1_

- [ ] 17.1 Write unit tests for routing
  - Test all routes render correctly
  - Test navigation between pages
  - Test 404 page
  - _Requirements: 1.1_


- [x] 18. Implement global error handling




  - Create error boundary component
  - Add toast notifications for API errors
  - Implement retry logic for failed requests
  - Add offline detection
  - _Requirements: 1.4_

- [x] 18.1 Write unit tests for error handling


  - Test error boundary
  - Test toast notifications
  - Test retry logic
  - _Requirements: 1.4_


- [x] 19. Apply Tailwind styling and theme



  - Configure Tailwind with custom color palette (purple, emerald, amber)
  - Create theme configuration matching TUI colors
  - Apply consistent spacing and typography
  - Ensure responsive design
  - _Requirements: 11.1, 11.2, 1.3_

- [ ] 20. Checkpoint - Ensure frontend tests pass
  - Ensure all tests pass, ask the user if questions arise.

- [x] 21. Set up Go backend API structure






  - Create api/ directory for API handlers
  - Create api/router.go for route definitions
  - Create api/handlers/ for handler functions
  - Create api/middleware/ for middleware
  - _Requirements: 9.1, 9.2, 12.1_


- [x] 22. Install Go backend dependencies





  - Install Gin web framework: `go get github.com/gin-gonic/gin`
  - Install CORS middleware: `go get github.com/gin-contrib/cors`
  - Install WebSocket library: `go get github.com/gorilla/websocket`
  - Update go.mod and go.sum
  - _Requirements: 9.1, 10.1_


- [x] 23. Implement SSH configuration API handlers




  - Create api/handlers/ssh.go
  - Implement GetSSHConfigs handler (GET /api/ssh)
  - Implement GetSSHConfig handler (GET /api/ssh/:id)
  - Implement CreateSSHConfig handler (POST /api/ssh)
  - Implement UpdateSSHConfig handler (PUT /api/ssh/:id)
  - Implement DeleteSSHConfig handler (DELETE /api/ssh/:id)
  - Implement TestSSHConnection handler (POST /api/ssh/:id/test)
  - _Requirements: 9.1, 9.3_

- [x] 23.1 Write property test for HTTP status codes


  - **Property 18: HTTP status codes for errors**
  - **Validates: Requirements 9.5**

- [x] 23.2 Write unit tests for SSH API handlers


  - Test GET /api/ssh returns all configs
  - Test POST /api/ssh creates config
  - Test PUT /api/ssh/:id updates config
  - Test DELETE /api/ssh/:id removes config
  - Test POST /api/ssh/:id/test validates connection



  - _Requirements: 9.1, 9.3_


- [ ] 24. Implement project API handlers

  - Create api/handlers/projects.go
  - Implement GetProjects handler (GET /api/projects)
  - Implement GetProject handler (GET /api/projects/:id)
  - Implement CreateProject handler (POST /api/projects)
  - Implement UpdateProject handler (PUT /api/projects/:id)


  - Implement DeleteProject handler (DELETE /api/projects/:id)
  - Implement DeployProject handler (POST /api/projects/:id/deploy)
  - _Requirements: 9.2, 9.4_

- [ ] 24.1 Write unit tests for project API handlers
  - Test GET /api/projects returns all projects
  - Test POST /api/projects creates project
  - Test PUT /api/projects/:id updates project
  - Test DELETE /api/projects/:id removes project
  - Test POST /api/projects/:id/deploy initiates deployment
  - _Requirements: 9.2, 9.4_

- [ ] 25. Implement WebSocket handler for deployment logs
  - Create api/handlers/websocket.go
  - Implement WebSocket upgrade handler
  - Implement log streaming logic
  - Handle client disconnections
  - Add deployment status updates
  - _Requirements: 9.4, 8.2_

- [ ] 25.1 Write unit tests for WebSocket handler
  - Test WebSocket connection upgrade
  - Test log message streaming
  - Test connection cleanup
  - _Requirements: 9.4_

- [x] 26. Implement CORS middleware


  - Create api/middleware/cors.go
  - Configure allowed origins
  - Handle preflight requests
  - Add CORS headers to responses
  - Log CORS errors
  - _Requirements: 10.1, 10.2, 10.4_

- [ ] 26.1 Write property test for CORS headers
  - **Property 19: CORS headers on preflight**



  - **Validates: Requirements 10.2, 10.4**


- [ ] 26.2 Write unit tests for CORS middleware
  - Test OPTIONS preflight handling
  - Test CORS headers on responses
  - Test allowed origins


  - _Requirements: 10.1, 10.2_





- [ ] 27. Implement error handling middleware

  - Create api/middleware/error.go
  - Add panic recovery


  - Format error responses consistently

  - Log errors with context
  - _Requirements: 9.5_

- [ ] 27.1 Write unit tests for error middleware
  - Test panic recovery
  - Test error response formatting


  - Test error logging
  - _Requirements: 9.5_

- [ ] 28. Implement logging middleware

  - Create api/middleware/logging.go
  - Log all incoming requests
  - Log response status and duration
  - Add request ID for tracing
  - _Requirements: 10.5_

- [ ] 29. Set up API router

  - Create api/router.go
  - Register all routes



  - Apply middleware
  - Configure route groups
  - _Requirements: 9.1, 9.2_

- [ ] 29.1 Write integration tests for API routes
  - Test all routes are registered




  - Test middleware is applied
  - Test route groups work correctly
  - _Requirements: 9.1, 9.2_

- [ ] 30. Refactor main.go for API mode
  - Add command-line flag for TUI vs API mode
  - Start HTTP server in API mode


  - Keep TUI code for backward compatibility
  - Configure server port and host
  - _Requirements: 9.1_

- [ ] 30.1 Write unit tests for main.go
  - Test server starts correctly



  - Test command-line flags

  - Test graceful shutdown
  - _Requirements: 9.1_

- [ ] 31. Create API response types


  - Add APIResponse struct to types.go
  - Add DeploymentLog struct
  - Add DeploymentStatus struct
  - Add JSON tags for serialization
  - _Requirements: 9.5_

- [ ] 32. Implement deployment engine

  - Create api/deploy/engine.go
  - Implement deployment orchestration logic
  - Handle build and deploy steps




  - Stream logs to WebSocket clients
  - Track deployment status
  - _Requirements: 8.1, 8.2, 8.3, 8.4_





- [ ] 32.1 Write unit tests for deployment engine
  - Test deployment orchestration
  - Test log streaming
  - Test status tracking
  - Test error handling


  - _Requirements: 8.1, 8.2, 8.3, 8.4_

- [ ] 33. Checkpoint - Ensure backend tests pass
  - Ensure all tests pass, ask the user if questions arise.

- [ ] 34. Configure frontend development proxy

  - Update vite.config.ts to proxy API requests to backend
  - Configure proxy for WebSocket connections
  - Set up environment variables for API URL
  - _Requirements: 1.2_

- [ ] 35. Test frontend-backend integration locally
  - Start backend server on port 8080
  - Start frontend dev server on port 5173
  - Test all API endpoints from frontend
  - Test WebSocket connection
  - Verify CORS is working
  - _Requirements: 1.2, 10.1_

- [ ] 35.1 Write end-to-end tests
  - Test complete SSH config workflow
  - Test complete project workflow
  - Test deployment workflow
  - Test error scenarios
  - _Requirements: 1.1, 2.1, 7.1, 8.1_

- [ ] 36. Build frontend for production

  - Run `npm run build` to create production bundle
  - Verify build output in dist/ directory
  - Test production build locally
  - _Requirements: 1.1_

- [ ] 37. Embed frontend in Go binary

  - Use Go embed to include frontend dist/ files
  - Serve static files from embedded filesystem
  - Configure fallback to index.html for SPA routing
  - Update main.go to serve embedded files
  - _Requirements: 1.1_

- [ ] 37.1 Write unit tests for static file serving
  - Test static files are served correctly
  - Test SPA fallback routing
  - Test embedded files are accessible
  - _Requirements: 1.1_

- [ ] 38. Final integration testing
  - Build complete application with embedded frontend
  - Test all features in production mode
  - Verify no console errors
  - Test on different browsers
  - _Requirements: 1.1, 1.5_

- [ ] 39. Update documentation
  - Update README.md with new web UI instructions
  - Document API endpoints
  - Add development setup instructions
  - Add deployment instructions
  - _Requirements: 1.1_

- [ ] 40. Final checkpoint - Ensure all tests pass
  - Ensure all tests pass, ask the user if questions arise.
