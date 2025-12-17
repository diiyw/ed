# easyd
A simple TUI-based deployer for personal usage built with Bubbletea.

## Features

- **Servers**: Add, edit, delete, and test SSH server configurations
- **Project Management**: Create projects with build instructions and deploy scripts
- **Deployment**: Deploy projects to multiple servers with real-time logging

## Getting Started

```bash
# Build the application
go build

# Run the application
./ed
```

## Usage

1. **Servers**:
   - Add SSH configurations for your servers
   - Test connections before deployment
   - Support for password, key, and SSH agent authentication

2. **Project Management**:
   - Create projects with build instructions and deploy scripts
   - Assign multiple deploy servers to each project
   - Edit and manage project configurations

3. **Deployment**:
   - Select a project and deploy to all configured servers
   - View real-time deployment logs
   - Automatic execution of build and deploy commands

## Configuration

All data is stored in `config.json` in the current directory. The application will create this file automatically on first run.
