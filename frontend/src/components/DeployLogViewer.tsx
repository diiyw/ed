import { useEffect, useRef } from 'react';
import { Loader2, CheckCircle, XCircle, Clock } from 'lucide-react';
import { Card, CardContent, CardHeader, CardTitle } from './ui/card';
import { Badge } from './ui/badge';
import type { DeploymentLog } from '@/types';

interface DeployLogViewerProps {
  logs: DeploymentLog[];
  status: 'pending' | 'running' | 'success' | 'failed';
}

export function DeployLogViewer({ logs, status }: DeployLogViewerProps) {
  const logEndRef = useRef<HTMLDivElement>(null);

  // Auto-scroll to bottom when new logs arrive
  useEffect(() => {
    logEndRef.current?.scrollIntoView({ behavior: 'smooth' });
  }, [logs]);

  const getStatusIcon = () => {
    switch (status) {
      case 'pending':
        return <Clock className="h-5 w-5 text-muted" />;
      case 'running':
        return <Loader2 className="h-5 w-5 animate-spin text-info" />;
      case 'success':
        return <CheckCircle className="h-5 w-5 text-success" />;
      case 'failed':
        return <XCircle className="h-5 w-5 text-error" />;
    }
  };

  const getStatusBadge = () => {
    switch (status) {
      case 'pending':
        return <Badge variant="outline">Pending</Badge>;
      case 'running':
        return <Badge variant="info">Running</Badge>;
      case 'success':
        return <Badge variant="success">Success</Badge>;
      case 'failed':
        return <Badge variant="destructive">Failed</Badge>;
    }
  };

  const getLogColor = (type: DeploymentLog['type']) => {
    switch (type) {
      case 'log':
        return 'text-gray-300';
      case 'status':
        return 'text-info';
      case 'error':
        return 'text-error';
      default:
        return 'text-gray-300';
    }
  };

  return (
    <Card className="bg-background-light border-gray-600">
      <CardHeader>
        <div className="flex items-center justify-between">
          <div className="flex items-center space-x-2">
            {getStatusIcon()}
            <CardTitle className="text-white">Deployment Logs</CardTitle>
          </div>
          {getStatusBadge()}
        </div>
      </CardHeader>
      <CardContent>
        <div className="bg-background-dark rounded-lg p-4 h-96 overflow-y-auto font-mono text-sm">
          {logs.length === 0 ? (
            <div className="flex items-center justify-center h-full text-gray-500">
              Waiting for logs...
            </div>
          ) : (
            <div className="space-y-1">
              {logs.map((log, index) => (
                <div key={index} className={`${getLogColor(log.type)} whitespace-pre-wrap`}>
                  <span className="text-gray-500 mr-2">
                    [{new Date(log.timestamp).toLocaleTimeString()}]
                  </span>
                  {log.data}
                </div>
              ))}
              <div ref={logEndRef} />
            </div>
          )}
        </div>
      </CardContent>
    </Card>
  );
}
