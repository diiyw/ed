export interface SSHConfig {
  name: string;
  host: string;
  port: number;
  user: string;
  authType: 'password' | 'key' | 'agent';
  password?: string;
  keyFile?: string;
  keyPass?: string;
}

export interface Project {
  name: string;
  buildInstructions: string;
  deployScript: string;
  deployServers: string[];
  createdAt: string;
  updatedAt: string;
}

export interface APIResponse<T> {
  data?: T;
  message?: string;
  error?: string;
}

export interface SSHTestResult {
  success: boolean;
  message: string;
  output?: string;
}
