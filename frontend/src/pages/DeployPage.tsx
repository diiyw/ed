import { useState, useEffect } from 'react';
import { useNavigate, useParams } from 'react-router-dom';
import { ArrowLeft, Loader2, XCircle } from 'lucide-react';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { DeployLogViewer } from '@/components/DeployLogViewer';
import { projectAPI } from '@/services/api';
import { DeploymentWebSocket } from '@/services/websocket';
import type { Project, DeploymentLog } from '@/types';

export function DeployPage() {
  const navigate = useNavigate();
  const { name } = useParams<{ name: string }>();

  const [project, setProject] = useState<Project | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [deploymentId, setDeploymentId] = useState<string | null>(null);
  const [logs, setLogs] = useState<DeploymentLog[]>([]);
  const [status, setStatus] = useState<'pending' | 'running' | 'success' | 'failed'>('pending');
  const [ws, setWs] = useState<DeploymentWebSocket | null>(null);

  useEffect(() => {
    if (name) {
      fetchProjectAndDeploy(decodeURIComponent(name));
    }
  }, [name]);

  useEffect(() => {
    // Cleanup WebSocket on unmount
    return () => {
      if (ws) {
        ws.disconnect();
      }
    };
  }, [ws]);

  const fetchProjectAndDeploy = async (projectName: string) => {
    try {
      setLoading(true);
      setError(null);

      // Fetch project details
      const projectData = await projectAPI.getByName(projectName);
      setProject(projectData);

      // Start deployment
      const result = await projectAPI.deploy(projectName);
      setDeploymentId(result.deploymentId);
      setStatus('running');

      // Connect to WebSocket for logs
      const websocket = new DeploymentWebSocket(result.deploymentId);
      
      websocket.onMessage((log) => {
        setLogs((prev) => [...prev, log]);
        
        // Update status based on log type
        if (log.type === 'status') {
          if (log.data.includes('success') || log.data.includes('completed')) {
            setStatus('success');
          } else if (log.data.includes('failed') || log.data.includes('error')) {
            setStatus('failed');
          }
        } else if (log.type === 'error') {
          setStatus('failed');
        }
      });

      websocket.onError((err) => {
        console.error('WebSocket error:', err);
        setError('Connection to deployment logs failed');
      });

      websocket.onClose(() => {
        console.log('WebSocket connection closed');
      });

      websocket.connect();
      setWs(websocket);

    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to start deployment');
      setStatus('failed');
    } finally {
      setLoading(false);
    }
  };

  const handleCancel = () => {
    if (ws) {
      ws.disconnect();
    }
    navigate('/projects');
  };

  if (loading) {
    return (
      <div className="flex items-center justify-center min-h-[400px]">
        <Loader2 className="h-8 w-8 animate-spin text-primary" />
      </div>
    );
  }

  if (error && !project) {
    return (
      <div className="max-w-4xl mx-auto">
        <div className="bg-destructive/10 border border-destructive/30 rounded-lg p-6 text-center">
          <XCircle className="h-12 w-12 text-destructive mx-auto mb-4" />
          <p className="text-destructive font-semibold mb-2">Deployment Failed</p>
          <p className="text-card-foreground">{error}</p>
          <Button onClick={() => navigate('/projects')} className="mt-4">
            Back to Projects
          </Button>
        </div>
      </div>
    );
  }

  return (
    <div className="max-w-4xl mx-auto space-y-6">
      {/* Header */}
      <div className="flex items-center justify-between">
        <div className="flex items-center space-x-4">
          <Button
            variant="ghost"
            size="icon"
            onClick={() => navigate('/projects')}
            className="text-muted-foreground hover:text-foreground"
          >
            <ArrowLeft className="h-5 w-5" />
          </Button>
          <div>
            <h1 className="text-3xl font-bold text-foreground">
              Deploying: {project?.name}
            </h1>
            <p className="text-muted-foreground mt-1">
              {project?.deployServers.length} server(s) • Deployment ID: {deploymentId}
            </p>
          </div>
        </div>
        {status === 'running' && (
          <Button variant="destructive" onClick={handleCancel}>
            Cancel Deployment
          </Button>
        )}
      </div>

      {/* Project Info */}
      <Card className="bg-card border-border">
        <CardHeader>
          <CardTitle className="text-foreground">Deployment Configuration</CardTitle>
        </CardHeader>
        <CardContent className="space-y-4">
          <div>
            <span className="text-sm text-muted-foreground">Build Instructions:</span>
            <pre className="mt-1 text-sm text-card-foreground bg-background p-3 rounded overflow-x-auto">
              {project?.buildInstructions}
            </pre>
          </div>
          <div>
            <span className="text-sm text-muted-foreground">Deploy Script:</span>
            <pre className="mt-1 text-sm text-card-foreground bg-background p-3 rounded overflow-x-auto">
              {project?.deployScript}
            </pre>
          </div>
          <div>
            <span className="text-sm text-muted-foreground">Target Servers:</span>
            <div className="mt-1 flex flex-wrap gap-2">
              {project?.deployServers.map((server) => (
                <span
                  key={server}
                  className="px-2 py-1 bg-primary/20 text-primary rounded text-sm"
                >
                  {server}
                </span>
              ))}
            </div>
          </div>
        </CardContent>
      </Card>

      {/* Deployment Logs */}
      <DeployLogViewer logs={logs} status={status} />

      {/* Actions */}
      {(status === 'success' || status === 'failed') && (
        <div className="flex justify-end space-x-2">
          <Button onClick={() => navigate('/projects')}>
            Back to Projects
          </Button>
          {status === 'failed' && (
            <Button
              onClick={() => {
                setLogs([]);
                setStatus('pending');
                if (name) {
                  fetchProjectAndDeploy(decodeURIComponent(name));
                }
              }}
              className="bg-secondary hover:bg-secondary/90"
            >
              Retry Deployment
            </Button>
          )}
        </div>
      )}
    </div>
  );
}
