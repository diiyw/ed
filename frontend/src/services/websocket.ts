import type { DeploymentLog } from '@/types';

type MessageHandler = (log: DeploymentLog) => void;
type ErrorHandler = (error: Error) => void;
type CloseHandler = () => void;

export class DeploymentWebSocket {
  private ws: WebSocket | null = null;
  private url: string;
  private messageHandlers: Set<MessageHandler> = new Set();
  private errorHandlers: Set<ErrorHandler> = new Set();
  private closeHandlers: Set<CloseHandler> = new Set();
  private reconnectAttempts = 0;
  private maxReconnectAttempts = 5;
  private reconnectDelay = 1000;
  private shouldReconnect = true;

  constructor(deploymentId: string) {
    const wsProtocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
    const wsBaseUrl = import.meta.env.VITE_WS_URL || `${wsProtocol}//${window.location.host}`;
    this.url = `${wsBaseUrl}/ws/deploy/${deploymentId}`;
  }

  connect(): void {
    try {
      this.ws = new WebSocket(this.url);

      this.ws.onopen = () => {
        console.log('WebSocket connected');
        this.reconnectAttempts = 0;
      };

      this.ws.onmessage = (event) => {
        try {
          const log: DeploymentLog = JSON.parse(event.data);
          this.messageHandlers.forEach((handler) => handler(log));
        } catch (error) {
          console.error('Failed to parse WebSocket message:', error);
        }
      };

      this.ws.onerror = (event) => {
        console.error('WebSocket error:', event);
        const error = new Error('WebSocket connection error');
        this.errorHandlers.forEach((handler) => handler(error));
      };

      this.ws.onclose = () => {
        console.log('WebSocket closed');
        this.closeHandlers.forEach((handler) => handler());

        // Attempt to reconnect if enabled
        if (this.shouldReconnect && this.reconnectAttempts < this.maxReconnectAttempts) {
          this.reconnectAttempts++;
          console.log(`Reconnecting... Attempt ${this.reconnectAttempts}/${this.maxReconnectAttempts}`);
          setTimeout(() => this.connect(), this.reconnectDelay * this.reconnectAttempts);
        }
      };
    } catch (error) {
      console.error('Failed to create WebSocket:', error);
      const err = error instanceof Error ? error : new Error('Failed to create WebSocket');
      this.errorHandlers.forEach((handler) => handler(err));
    }
  }

  disconnect(): void {
    this.shouldReconnect = false;
    if (this.ws) {
      this.ws.close();
      this.ws = null;
    }
  }

  onMessage(handler: MessageHandler): () => void {
    this.messageHandlers.add(handler);
    return () => this.messageHandlers.delete(handler);
  }

  onError(handler: ErrorHandler): () => void {
    this.errorHandlers.add(handler);
    return () => this.errorHandlers.delete(handler);
  }

  onClose(handler: CloseHandler): () => void {
    this.closeHandlers.add(handler);
    return () => this.closeHandlers.delete(handler);
  }

  isConnected(): boolean {
    return this.ws !== null && this.ws.readyState === WebSocket.OPEN;
  }
}

export default DeploymentWebSocket;
