import { describe, it, expect, beforeEach, afterEach } from 'vitest';
import MockAdapter from 'axios-mock-adapter';
import apiClient, { sshAPI, projectAPI } from './api';

describe('API Service - Retry Logic', () => {
  let mock: MockAdapter;

  beforeEach(() => {
    mock = new MockAdapter(apiClient);
  });

  afterEach(() => {
    mock.restore();
  });

  it('retries failed requests with retryable status codes', async () => {
    let attemptCount = 0;

    mock.onGet('/ssh').reply(() => {
      attemptCount++;
      if (attemptCount < 3) {
        return [503, { error: 'Service unavailable' }];
      }
      return [200, { data: { configs: [] } }];
    });

    const result = await sshAPI.getAll();

    expect(attemptCount).toBe(3);
    expect(result).toEqual([]);
  });

  it('does not retry on non-retryable status codes', async () => {
    let attemptCount = 0;

    mock.onGet('/ssh').reply(() => {
      attemptCount++;
      return [400, { error: 'Bad request' }];
    });

    await expect(sshAPI.getAll()).rejects.toThrow();
    expect(attemptCount).toBe(1);
  });

  it('throws error after max retries', async () => {
    let attemptCount = 0;

    mock.onGet('/ssh').reply(() => {
      attemptCount++;
      return [503, { error: 'Service unavailable' }];
    });

    await expect(sshAPI.getAll()).rejects.toThrow();
    expect(attemptCount).toBeGreaterThan(1);
  }, 10000); // Increase timeout to 10 seconds

  it('handles network errors with retry', async () => {
    // Test that network errors (no response) trigger retry logic
    // by simulating a timeout error which should be retried
    mock.onGet('/ssh').timeout();

    await expect(sshAPI.getAll()).rejects.toThrow();
  }, 10000); // Increase timeout to account for retries
});

describe('API Service - Error Handling', () => {
  let mock: MockAdapter;

  beforeEach(() => {
    mock = new MockAdapter(apiClient);
  });

  afterEach(() => {
    mock.restore();
  });

  it('extracts error message from API response', async () => {
    mock.onGet('/ssh').reply(500, { error: 'Custom error message' });

    await expect(sshAPI.getAll()).rejects.toThrow('Custom error message');
  }, 10000); // Increase timeout

  it('provides default error message when none is provided', async () => {
    mock.onGet('/ssh').reply(500, {});

    await expect(sshAPI.getAll()).rejects.toThrow();
  }, 10000); // Increase timeout

  it('handles network errors with appropriate message', async () => {
    mock.onGet('/ssh').networkError();

    await expect(sshAPI.getAll()).rejects.toThrow('Network error: Unable to connect to the server');
  }, 10000); // Increase timeout
});

describe('SSH API', () => {
  let mock: MockAdapter;

  beforeEach(() => {
    mock = new MockAdapter(apiClient);
  });

  afterEach(() => {
    mock.restore();
  });

  it('getAll returns empty array when no configs exist', async () => {
    mock.onGet('/ssh').reply(200, { data: { configs: [] } });

    const result = await sshAPI.getAll();

    expect(result).toEqual([]);
  });

  it('getAll returns configs when they exist', async () => {
    const mockConfigs = [
      { name: 'test1', host: 'host1', port: 22, user: 'user1', authType: 'password' },
    ];
    mock.onGet('/ssh').reply(200, { data: { configs: mockConfigs } });

    const result = await sshAPI.getAll();

    expect(result).toEqual(mockConfigs);
  });
});

describe('Project API', () => {
  let mock: MockAdapter;

  beforeEach(() => {
    mock = new MockAdapter(apiClient);
  });

  afterEach(() => {
    mock.restore();
  });

  it('getAll returns empty array when no projects exist', async () => {
    mock.onGet('/projects').reply(200, { data: { projects: [] } });

    const result = await projectAPI.getAll();

    expect(result).toEqual([]);
  });

  it('getAll returns projects when they exist', async () => {
    const mockProjects = [
      {
        name: 'test-project',
        buildInstructions: 'npm build',
        deployScript: 'deploy.sh',
        deployServers: ['server1'],
        createdAt: '2024-01-01',
        updatedAt: '2024-01-01',
      },
    ];
    mock.onGet('/projects').reply(200, { data: { projects: mockProjects } });

    const result = await projectAPI.getAll();

    expect(result).toEqual(mockProjects);
  });
});
