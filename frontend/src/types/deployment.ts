export interface DeploymentLog {
  type: 'log' | 'status' | 'error';
  data: string;
  timestamp: string;
}

export interface DeploymentStatus {
  id: string;
  projectName: string;
  status: 'pending' | 'running' | 'success' | 'failed';
  startedAt: string;
  completedAt?: string;
}
