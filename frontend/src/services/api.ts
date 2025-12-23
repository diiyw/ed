import axios, { type AxiosInstance, type AxiosError, type InternalAxiosRequestConfig } from 'axios';
import type { SSHConfig, Project, APIResponse, SSHTestResult } from '@/types';

// Retry configuration
const MAX_RETRIES = 3;
const RETRY_DELAY = 1000; // 1 second
const RETRYABLE_STATUS_CODES = [408, 429, 500, 502, 503, 504];

// Create axios instance with default config
const apiClient: AxiosInstance = axios.create({
  baseURL: import.meta.env.VITE_API_URL || '/api',
  timeout: 30000,
  headers: {
    'Content-Type': 'application/json',
  },
});

// Request interceptor to add retry count
apiClient.interceptors.request.use(
  (config: InternalAxiosRequestConfig) => {
    // Initialize retry count if not present
    if (!config.headers) {
      config.headers = {} as unknown;
    }
    return config;
  },
  (error) => Promise.reject(error)
);

// Response interceptor for error handling and retry logic
apiClient.interceptors.response.use(
  (response) => response,
  async (error: AxiosError<APIResponse<unknown>>) => {
    const config = error.config as InternalAxiosRequestConfig & { _retryCount?: number };
    
    if (!config) {
      return Promise.reject(error);
    }

    // Initialize retry count
    config._retryCount = config._retryCount || 0;

    // Check if we should retry
    const shouldRetry = 
      config._retryCount < MAX_RETRIES &&
      (!error.response || RETRYABLE_STATUS_CODES.includes(error.response.status));

    if (shouldRetry) {
      config._retryCount += 1;
      
      // Wait before retrying with exponential backoff
      const delay = RETRY_DELAY * Math.pow(2, config._retryCount - 1);
      await new Promise(resolve => setTimeout(resolve, delay));
      
      return apiClient(config);
    }

    // Handle network errors
    if (!error.response) {
      throw new Error('Network error: Unable to connect to the server');
    }

    // Extract error message from response
    const errorMessage = error.response.data?.error || error.message || 'An unknown error occurred';
    throw new Error(errorMessage);
  }
);

// SSH Configuration API functions
export const sshAPI = {
  // Get all SSH configurations
  async getAll(): Promise<SSHConfig[]> {
    const response = await apiClient.get<APIResponse<{ configs: SSHConfig[] }>>('/ssh');
    return response.data.data?.configs || [];
  },

  // Get a single SSH configuration by name
  async getByName(name: string): Promise<SSHConfig> {
    const response = await apiClient.get<APIResponse<{ config: SSHConfig }>>(`/ssh/${encodeURIComponent(name)}`);
    if (!response.data.data?.config) {
      throw new Error('SSH configuration not found');
    }
    return response.data.data.config;
  },

  // Create a new SSH configuration
  async create(config: SSHConfig): Promise<SSHConfig> {
    const response = await apiClient.post<APIResponse<{ config: SSHConfig }>>('/ssh', config);
    if (!response.data.data?.config) {
      throw new Error('Failed to create SSH configuration');
    }
    return response.data.data.config;
  },

  // Update an existing SSH configuration
  async update(name: string, config: SSHConfig): Promise<SSHConfig> {
    const response = await apiClient.put<APIResponse<{ config: SSHConfig }>>(
      `/ssh/${encodeURIComponent(name)}`,
      config
    );
    if (!response.data.data?.config) {
      throw new Error('Failed to update SSH configuration');
    }
    return response.data.data.config;
  },

  // Delete an SSH configuration
  async delete(name: string): Promise<void> {
    await apiClient.delete(`/ssh/${encodeURIComponent(name)}`);
  },

  // Test SSH connection
  async test(name: string): Promise<SSHTestResult> {
    const response = await apiClient.post<SSHTestResult>(`/ssh/${encodeURIComponent(name)}/test`);
    return response.data;
  },
};

// Project API functions
export const projectAPI = {
  // Get all projects
  async getAll(): Promise<Project[]> {
    const response = await apiClient.get<APIResponse<{ projects: Project[] }>>('/projects');
    return response.data.data?.projects || [];
  },

  // Get a single project by name
  async getByName(name: string): Promise<Project> {
    const response = await apiClient.get<APIResponse<{ project: Project }>>(`/projects/${encodeURIComponent(name)}`);
    if (!response.data.data?.project) {
      throw new Error('Project not found');
    }
    return response.data.data.project;
  },

  // Create a new project
  async create(project: Project): Promise<Project> {
    const response = await apiClient.post<APIResponse<{ project: Project }>>('/projects', project);
    if (!response.data.data?.project) {
      throw new Error('Failed to create project');
    }
    return response.data.data.project;
  },

  // Update an existing project
  async update(name: string, project: Project): Promise<Project> {
    const response = await apiClient.put<APIResponse<{ project: Project }>>(
      `/projects/${encodeURIComponent(name)}`,
      project
    );
    if (!response.data.data?.project) {
      throw new Error('Failed to update project');
    }
    return response.data.data.project;
  },

  // Delete a project
  async delete(name: string): Promise<void> {
    await apiClient.delete(`/projects/${encodeURIComponent(name)}`);
  },

  // Deploy a project
  async deploy(name: string): Promise<{ deploymentId: string; message: string }> {
    const response = await apiClient.post<APIResponse<{ deploymentId: string; message: string }>>(
      `/projects/${encodeURIComponent(name)}/deploy`
    );
    if (!response.data.data) {
      throw new Error('Failed to start deployment');
    }
    return response.data.data;
  },
};

export default apiClient;
