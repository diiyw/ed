# Frontend Development Proxy Configuration

## Overview

The Vite development server is configured to proxy API and WebSocket requests to the backend Go server. This allows the frontend to run on `http://localhost:5173` while the backend runs on `http://localhost:8080`, avoiding CORS issues during development.

## Configuration

### Vite Proxy Settings

The proxy is configured in `vite.config.ts`:

- **API Proxy** (`/api`): Proxies all HTTP requests starting with `/api` to the backend
- **WebSocket Proxy** (`/ws`): Proxies all WebSocket connections starting with `/ws` to the backend

### Environment Variables

The following environment variables control the proxy behavior:

#### `.env` file:

```env
# API Configuration - Used by the application code
VITE_API_URL=http://localhost:8080/api
VITE_WS_URL=ws://localhost:8080

# Development Proxy Configuration - Used by Vite dev server
VITE_API_PROXY_TARGET=http://localhost:8080
VITE_WS_PROXY_TARGET=ws://localhost:8080
```

## How It Works

### Development Mode

When running `npm run dev`:

1. Frontend runs on `http://localhost:5173`
2. Backend runs on `http://localhost:8080`
3. API requests to `/api/*` are proxied to `http://localhost:8080/api/*`
4. WebSocket connections to `/ws/*` are proxied to `ws://localhost:8080/ws/*`

Example:
- Frontend makes request to: `http://localhost:5173/api/ssh`
- Vite proxy forwards to: `http://localhost:8080/api/ssh`

### Production Mode

In production, the frontend is built as static files and served directly by the Go backend. No proxy is needed since both frontend and backend are served from the same origin.

## Customizing Backend Location

If your backend is running on a different host or port, update the `.env` file:

```env
VITE_API_PROXY_TARGET=http://localhost:3000
VITE_WS_PROXY_TARGET=ws://localhost:3000
```

## Testing the Proxy

1. Start the backend server:
   ```bash
   go run main.go
   ```

2. Start the frontend dev server:
   ```bash
   cd frontend
   npm run dev
   ```

3. Open `http://localhost:5173` in your browser
4. The frontend should successfully communicate with the backend through the proxy

## Troubleshooting

### CORS Errors

If you see CORS errors, ensure:
- The backend CORS middleware is properly configured
- The proxy `changeOrigin: true` setting is enabled in `vite.config.ts`

### WebSocket Connection Failures

If WebSocket connections fail:
- Verify the backend WebSocket handler is running
- Check that `ws: true` is set in the proxy configuration
- Ensure the WebSocket URL uses the correct protocol (`ws://` or `wss://`)

### Proxy Not Working

If the proxy isn't forwarding requests:
- Restart the Vite dev server after changing `vite.config.ts`
- Check the Vite dev server console for proxy-related errors
- Verify the backend is running and accessible at the configured target URL
