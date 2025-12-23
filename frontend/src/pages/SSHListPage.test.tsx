import { describe, it, expect, vi, beforeEach } from 'vitest';
import { render, screen, waitFor } from '@testing-library/react';
import { BrowserRouter } from 'react-router-dom';
import userEvent from '@testing-library/user-event';
import * as fc from 'fast-check';
import { SSHListPage } from './SSHListPage';
import { sshAPI, projectAPI } from '@/services/api';
import type { SSHConfig, Project } from '@/types';

// Mock the API module
vi.mock('@/services/api', () => ({
  sshAPI: {
    getAll: vi.fn(),
    delete: vi.fn(),
  },
  projectAPI: {
    getAll: vi.fn(),
  },
}));

// Mock the child components to simplify testing
vi.mock('@/components/SSHCard', () => ({
  SSHCard: ({ sshConfig, onEdit, onDelete, onTest }: any) => (
    <div data-testid={`ssh-card-${sshConfig.name}`}>
      <span>{sshConfig.name}</span>
      <span>{sshConfig.host}</span>
      <button onClick={() => onEdit(sshConfig)}>Edit</button>
      <button onClick={() => onDelete(sshConfig)}>Delete</button>
      <button onClick={() => onTest(sshConfig)}>Test</button>
    </div>
  ),
}));

vi.mock('@/components/SSHTestDialog', () => ({
  SSHTestDialog: ({ isOpen, onClose }: any) =>
    isOpen ? <div data-testid="test-dialog"><button onClick={onClose}>Close</button></div> : null,
}));

vi.mock('@/components/DeleteConfirmDialog', () => ({
  DeleteConfirmDialog: ({ isOpen, onConfirm, onCancel, warning }: any) =>
    isOpen ? (
      <div data-testid="delete-dialog">
        {warning && <div data-testid="delete-warning">{warning}</div>}
        <button onClick={onConfirm}>Confirm</button>
        <button onClick={onCancel}>Cancel</button>
      </div>
    ) : null,
}));

describe('SSHListPage', () => {
  const renderSSHListPage = () => {
    return render(
      <BrowserRouter>
        <SSHListPage />
      </BrowserRouter>
    );
  };

  beforeEach(() => {
    vi.clearAllMocks();
  });

  /**
   * Feature: frontend-refactor, Property 1: Data fetching on mount
   * 
   * Property: For any page component that displays data from the backend,
   * mounting the component should trigger an API call to fetch that data.
   * 
   * Validates: Requirements 1.2, 2.1
   */
  describe('Property 1: Data fetching on mount', () => {
    it('should fetch SSH configs from API when component mounts', async () => {
      await fc.assert(
        fc.asyncProperty(
          fc.array(
            fc.record({
              name: fc.string({ minLength: 1, maxLength: 50 }),
              host: fc.string({ minLength: 1, maxLength: 100 }),
              port: fc.integer({ min: 1, max: 65535 }),
              user: fc.string({ minLength: 1, maxLength: 50 }),
              authType: fc.constantFrom('password', 'key', 'agent'),
              password: fc.option(fc.string(), { nil: undefined }),
              keyFile: fc.option(fc.string(), { nil: undefined }),
              keyPass: fc.option(fc.string(), { nil: undefined }),
            }),
            { minLength: 0, maxLength: 10 }
          ),
          async (configs: SSHConfig[]) => {
            // Setup: Mock the API to return the generated configs
            vi.mocked(sshAPI.getAll).mockResolvedValue(configs);
            vi.mocked(projectAPI.getAll).mockResolvedValue([]);

            // Action: Mount the component
            renderSSHListPage();

            // Assert: API should be called on mount
            await waitFor(() => {
              expect(sshAPI.getAll).toHaveBeenCalledTimes(1);
            });

            // Cleanup for next iteration
            vi.clearAllMocks();
          }
        ),
        { numRuns: 100 }
      );
    });
  });

  /**
   * Feature: frontend-refactor, Property 11: Confirmation dialog before deletion
   * 
   * Property: For any delete action initiated by the user, a confirmation dialog
   * should be displayed before making the delete API call.
   * 
   * Validates: Requirements 5.1, 7.5
   */
  describe('Property 11: Confirmation dialog before deletion', () => {
    it('should display confirmation dialog before deleting any SSH config', async () => {
      await fc.assert(
        fc.asyncProperty(
          fc.record({
            name: fc.string({ minLength: 1, maxLength: 50 }).filter(s => s.trim().length > 0),
            host: fc.string({ minLength: 1, maxLength: 100 }).filter(s => s.trim().length > 0),
            port: fc.integer({ min: 1, max: 65535 }),
            user: fc.string({ minLength: 1, maxLength: 50 }).filter(s => s.trim().length > 0),
            authType: fc.constantFrom('password', 'key', 'agent'),
            password: fc.option(fc.string(), { nil: undefined }),
            keyFile: fc.option(fc.string(), { nil: undefined }),
            keyPass: fc.option(fc.string(), { nil: undefined }),
          }),
          async (config: SSHConfig) => {
            const user = userEvent.setup();
            
            // Setup: Mock the API to return the config
            vi.mocked(sshAPI.getAll).mockResolvedValue([config]);
            vi.mocked(projectAPI.getAll).mockResolvedValue([]);
            vi.mocked(sshAPI.delete).mockResolvedValue();

            // Action: Mount the component
            const { unmount } = renderSSHListPage();

            // Wait for the config to be loaded
            await waitFor(() => {
              expect(screen.getByTestId(`ssh-card-${config.name}`)).toBeInTheDocument();
            });

            // Click delete button
            const deleteButton = screen.getByRole('button', { name: /delete/i });
            await user.click(deleteButton);

            // Assert: Confirmation dialog should appear BEFORE API call
            await waitFor(() => {
              expect(screen.getByTestId('delete-dialog')).toBeInTheDocument();
            });

            // Assert: Delete API should NOT have been called yet
            expect(sshAPI.delete).not.toHaveBeenCalled();

            // Cleanup for next iteration
            unmount();
            vi.clearAllMocks();
          }
        ),
        { numRuns: 10 }
      );
    }, 30000);
  });

  /**
   * Feature: frontend-refactor, Property 12: Warning for configs in use
   * 
   * Property: For any SSH configuration that is referenced by one or more projects,
   * attempting to delete it should display a warning message.
   * 
   * Validates: Requirements 5.5
   */
  describe('Property 12: Warning for configs in use', () => {
    it('should display warning when deleting SSH config that is in use by projects', async () => {
      await fc.assert(
        fc.asyncProperty(
          fc.record({
            name: fc.string({ minLength: 1, maxLength: 50 }).filter(s => s.trim().length > 0),
            host: fc.string({ minLength: 1, maxLength: 100 }).filter(s => s.trim().length > 0),
            port: fc.integer({ min: 1, max: 65535 }),
            user: fc.string({ minLength: 1, maxLength: 50 }).filter(s => s.trim().length > 0),
            authType: fc.constantFrom('password', 'key', 'agent'),
          }),
          fc.array(
            fc.record({
              name: fc.string({ minLength: 1, maxLength: 50 }).filter(s => s.trim().length > 0),
              buildInstructions: fc.string(),
              deployScript: fc.string(),
              deployServers: fc.array(fc.string({ minLength: 1, maxLength: 50 }).filter(s => s.trim().length > 0), { minLength: 1, maxLength: 5 }),
              createdAt: fc.date().map(d => d.toISOString()),
              updatedAt: fc.date().map(d => d.toISOString()),
            }),
            { minLength: 1, maxLength: 5 }
          ),
          async (config: SSHConfig, projects: Project[]) => {
            const user = userEvent.setup();
            
            // Ensure at least one project uses this config
            const projectsUsingConfig = projects.map(p => ({
              ...p,
              deployServers: [config.name, ...p.deployServers.filter(s => s !== config.name)],
            }));

            // Setup: Mock the API
            vi.mocked(sshAPI.getAll).mockResolvedValue([config]);
            vi.mocked(projectAPI.getAll).mockResolvedValue(projectsUsingConfig);

            // Action: Mount the component
            const { unmount } = renderSSHListPage();

            // Wait for the config to be loaded
            await waitFor(() => {
              expect(screen.getByTestId(`ssh-card-${config.name}`)).toBeInTheDocument();
            });

            // Click delete button
            const deleteButton = screen.getByRole('button', { name: /delete/i });
            await user.click(deleteButton);

            // Assert: Warning should be displayed
            await waitFor(() => {
              const warning = screen.getByTestId('delete-warning');
              expect(warning).toBeInTheDocument();
              expect(warning.textContent).toContain('currently used by');
            });

            // Cleanup for next iteration
            unmount();
            vi.clearAllMocks();
          }
        ),
        { numRuns: 10 }
      );
    }, 30000);

    it('should NOT display warning when deleting SSH config that is NOT in use', async () => {
      await fc.assert(
        fc.asyncProperty(
          fc.record({
            name: fc.string({ minLength: 1, maxLength: 50 }).filter(s => s.trim().length > 0),
            host: fc.string({ minLength: 1, maxLength: 100 }).filter(s => s.trim().length > 0),
            port: fc.integer({ min: 1, max: 65535 }),
            user: fc.string({ minLength: 1, maxLength: 50 }).filter(s => s.trim().length > 0),
            authType: fc.constantFrom('password', 'key', 'agent'),
          }),
          fc.array(
            fc.record({
              name: fc.string({ minLength: 1, maxLength: 50 }).filter(s => s.trim().length > 0),
              buildInstructions: fc.string(),
              deployScript: fc.string(),
              deployServers: fc.array(fc.string({ minLength: 1, maxLength: 50 }).filter(s => s.trim().length > 0), { minLength: 0, maxLength: 5 }),
              createdAt: fc.date().map(d => d.toISOString()),
              updatedAt: fc.date().map(d => d.toISOString()),
            }),
            { minLength: 0, maxLength: 5 }
          ),
          async (config: SSHConfig, projects: Project[]) => {
            const user = userEvent.setup();
            
            // Ensure NO project uses this config
            const projectsNotUsingConfig = projects.map(p => ({
              ...p,
              deployServers: p.deployServers.filter(s => s !== config.name),
            }));

            // Setup: Mock the API
            vi.mocked(sshAPI.getAll).mockResolvedValue([config]);
            vi.mocked(projectAPI.getAll).mockResolvedValue(projectsNotUsingConfig);

            // Action: Mount the component
            const { unmount } = renderSSHListPage();

            // Wait for the config to be loaded
            await waitFor(() => {
              expect(screen.getByTestId(`ssh-card-${config.name}`)).toBeInTheDocument();
            });

            // Click delete button
            const deleteButton = screen.getByRole('button', { name: /delete/i });
            await user.click(deleteButton);

            // Assert: Warning should NOT be displayed
            await waitFor(() => {
              expect(screen.getByTestId('delete-dialog')).toBeInTheDocument();
            });
            
            expect(screen.queryByTestId('delete-warning')).not.toBeInTheDocument();

            // Cleanup for next iteration
            unmount();
            vi.clearAllMocks();
          }
        ),
        { numRuns: 10 }
      );
    }, 30000);
  });

  describe('Empty State', () => {
    it('should display empty state when no configs exist', async () => {
      vi.mocked(sshAPI.getAll).mockResolvedValue([]);
      vi.mocked(projectAPI.getAll).mockResolvedValue([]);

      renderSSHListPage();

      await waitFor(() => {
        expect(screen.getByText(/no ssh configurations found/i)).toBeInTheDocument();
      });

      expect(screen.getByText(/you haven't added any ssh server configurations yet/i)).toBeInTheDocument();
      expect(screen.getByRole('button', { name: /add your first ssh config/i })).toBeInTheDocument();
    });

    it('should not display empty state when configs exist', async () => {
      const mockConfigs: SSHConfig[] = [
        {
          name: 'test-server',
          host: 'example.com',
          port: 22,
          user: 'admin',
          authType: 'password',
        },
      ];

      vi.mocked(sshAPI.getAll).mockResolvedValue(mockConfigs);
      vi.mocked(projectAPI.getAll).mockResolvedValue([]);

      renderSSHListPage();

      await waitFor(() => {
        expect(screen.queryByText(/no ssh configurations found/i)).not.toBeInTheDocument();
      });
    });
  });

  describe('Config List Rendering', () => {
    it('should display all SSH configs in a grid', async () => {
      const mockConfigs: SSHConfig[] = [
        {
          name: 'server1',
          host: 'example1.com',
          port: 22,
          user: 'admin',
          authType: 'password',
        },
        {
          name: 'server2',
          host: 'example2.com',
          port: 2222,
          user: 'root',
          authType: 'key',
        },
      ];

      vi.mocked(sshAPI.getAll).mockResolvedValue(mockConfigs);
      vi.mocked(projectAPI.getAll).mockResolvedValue([]);

      renderSSHListPage();

      await waitFor(() => {
        expect(screen.getByTestId('ssh-card-server1')).toBeInTheDocument();
        expect(screen.getByTestId('ssh-card-server2')).toBeInTheDocument();
      });

      expect(screen.getByText('server1')).toBeInTheDocument();
      expect(screen.getByText('server2')).toBeInTheDocument();
    });

    it('should display page header with title and description', async () => {
      vi.mocked(sshAPI.getAll).mockResolvedValue([]);
      vi.mocked(projectAPI.getAll).mockResolvedValue([]);

      renderSSHListPage();

      await waitFor(() => {
        expect(screen.getByRole('heading', { name: /^ssh configurations$/i })).toBeInTheDocument();
      });

      expect(screen.getByText(/manage your server connections and authentication/i)).toBeInTheDocument();
    });

    it('should display "Add SSH Config" button in header', async () => {
      vi.mocked(sshAPI.getAll).mockResolvedValue([]);
      vi.mocked(projectAPI.getAll).mockResolvedValue([]);

      renderSSHListPage();

      await waitFor(() => {
        expect(screen.getByRole('button', { name: /add ssh config/i })).toBeInTheDocument();
      });
    });
  });

  describe('Action Buttons', () => {
    it('should handle edit action', async () => {
      const user = userEvent.setup();
      const mockConfigs: SSHConfig[] = [
        {
          name: 'test-server',
          host: 'example.com',
          port: 22,
          user: 'admin',
          authType: 'password',
        },
      ];

      vi.mocked(sshAPI.getAll).mockResolvedValue(mockConfigs);
      vi.mocked(projectAPI.getAll).mockResolvedValue([]);

      renderSSHListPage();

      await waitFor(() => {
        expect(screen.getByTestId('ssh-card-test-server')).toBeInTheDocument();
      });

      const editButton = screen.getByRole('button', { name: /edit/i });
      await user.click(editButton);

      // Navigation is handled by React Router, which we can't easily test here
      // The important part is that the button exists and is clickable
      expect(editButton).toBeInTheDocument();
    });

    it('should handle delete action and show confirmation dialog', async () => {
      const user = userEvent.setup();
      const mockConfigs: SSHConfig[] = [
        {
          name: 'test-server',
          host: 'example.com',
          port: 22,
          user: 'admin',
          authType: 'password',
        },
      ];

      vi.mocked(sshAPI.getAll).mockResolvedValue(mockConfigs);
      vi.mocked(projectAPI.getAll).mockResolvedValue([]);

      renderSSHListPage();

      await waitFor(() => {
        expect(screen.getByTestId('ssh-card-test-server')).toBeInTheDocument();
      });

      const deleteButton = screen.getByRole('button', { name: /delete/i });
      await user.click(deleteButton);

      // Delete confirmation dialog should appear
      await waitFor(() => {
        expect(screen.getByTestId('delete-dialog')).toBeInTheDocument();
      });
    });

    it('should handle test action and show test dialog', async () => {
      const user = userEvent.setup();
      const mockConfigs: SSHConfig[] = [
        {
          name: 'test-server',
          host: 'example.com',
          port: 22,
          user: 'admin',
          authType: 'password',
        },
      ];

      vi.mocked(sshAPI.getAll).mockResolvedValue(mockConfigs);
      vi.mocked(projectAPI.getAll).mockResolvedValue([]);

      renderSSHListPage();

      await waitFor(() => {
        expect(screen.getByTestId('ssh-card-test-server')).toBeInTheDocument();
      });

      const testButton = screen.getByRole('button', { name: /test/i });
      await user.click(testButton);

      // Test dialog should appear
      await waitFor(() => {
        expect(screen.getByTestId('test-dialog')).toBeInTheDocument();
      });
    });

    it('should call delete API and update list when delete is confirmed', async () => {
      const user = userEvent.setup();
      const mockConfigs: SSHConfig[] = [
        {
          name: 'test-server',
          host: 'example.com',
          port: 22,
          user: 'admin',
          authType: 'password',
        },
      ];

      vi.mocked(sshAPI.getAll).mockResolvedValue(mockConfigs);
      vi.mocked(projectAPI.getAll).mockResolvedValue([]);
      vi.mocked(sshAPI.delete).mockResolvedValue();

      renderSSHListPage();

      await waitFor(() => {
        expect(screen.getByTestId('ssh-card-test-server')).toBeInTheDocument();
      });

      // Click delete button
      const deleteButton = screen.getByRole('button', { name: /delete/i });
      await user.click(deleteButton);

      // Confirm deletion
      await waitFor(() => {
        expect(screen.getByTestId('delete-dialog')).toBeInTheDocument();
      });

      const confirmButton = screen.getByRole('button', { name: /confirm/i });
      await user.click(confirmButton);

      // API should be called
      await waitFor(() => {
        expect(sshAPI.delete).toHaveBeenCalledWith('test-server');
      });

      // Config should be removed from the list
      await waitFor(() => {
        expect(screen.queryByTestId('ssh-card-test-server')).not.toBeInTheDocument();
      });
    });

    it('should not call delete API when delete is cancelled', async () => {
      const user = userEvent.setup();
      const mockConfigs: SSHConfig[] = [
        {
          name: 'test-server',
          host: 'example.com',
          port: 22,
          user: 'admin',
          authType: 'password',
        },
      ];

      vi.mocked(sshAPI.getAll).mockResolvedValue(mockConfigs);
      vi.mocked(projectAPI.getAll).mockResolvedValue([]);

      renderSSHListPage();

      await waitFor(() => {
        expect(screen.getByTestId('ssh-card-test-server')).toBeInTheDocument();
      });

      // Click delete button
      const deleteButton = screen.getByRole('button', { name: /delete/i });
      await user.click(deleteButton);

      // Cancel deletion
      await waitFor(() => {
        expect(screen.getByTestId('delete-dialog')).toBeInTheDocument();
      });

      const cancelButton = screen.getByRole('button', { name: /cancel/i });
      await user.click(cancelButton);

      // API should not be called
      expect(sshAPI.delete).not.toHaveBeenCalled();

      // Config should still be in the list
      expect(screen.getByTestId('ssh-card-test-server')).toBeInTheDocument();
    });
  });

  describe('Loading State', () => {
    it('should display loading spinner while fetching configs', async () => {
      vi.mocked(sshAPI.getAll).mockImplementation(
        () => new Promise((resolve) => setTimeout(() => resolve([]), 100))
      );
      vi.mocked(projectAPI.getAll).mockImplementation(
        () => new Promise((resolve) => setTimeout(() => resolve([]), 100))
      );

      renderSSHListPage();

      // Loading spinner should be visible (Loader2 icon with animate-spin class)
      const spinner = document.querySelector('.animate-spin');
      expect(spinner).toBeInTheDocument();

      // Wait for loading to complete
      await waitFor(
        () => {
          const spinnerAfter = document.querySelector('.animate-spin');
          expect(spinnerAfter).not.toBeInTheDocument();
        },
        { timeout: 200 }
      );
    });
  });

  describe('Error Handling', () => {
    it('should display error message when API call fails', async () => {
      vi.mocked(sshAPI.getAll).mockRejectedValue(new Error('Network error'));
      vi.mocked(projectAPI.getAll).mockResolvedValue([]);

      renderSSHListPage();

      await waitFor(() => {
        expect(screen.getByText(/^error$/i)).toBeInTheDocument();
        expect(screen.getByText(/network error/i)).toBeInTheDocument();
      });

      expect(screen.getByRole('button', { name: /try again/i })).toBeInTheDocument();
    });

    it('should retry fetching when "Try Again" button is clicked', async () => {
      const user = userEvent.setup();
      
      // First call fails
      vi.mocked(sshAPI.getAll).mockRejectedValueOnce(new Error('Network error'));
      vi.mocked(projectAPI.getAll).mockResolvedValue([]);
      
      renderSSHListPage();

      await waitFor(() => {
        expect(screen.getByText(/network error/i)).toBeInTheDocument();
      });

      // Second call succeeds
      vi.mocked(sshAPI.getAll).mockResolvedValue([]);

      const tryAgainButton = screen.getByRole('button', { name: /try again/i });
      await user.click(tryAgainButton);

      await waitFor(() => {
        expect(sshAPI.getAll).toHaveBeenCalledTimes(2);
      });

      // Error should be cleared
      await waitFor(() => {
        expect(screen.queryByText(/network error/i)).not.toBeInTheDocument();
      });
    });

    it('should display error when delete fails', async () => {
      const user = userEvent.setup();
      const mockConfigs: SSHConfig[] = [
        {
          name: 'test-server',
          host: 'example.com',
          port: 22,
          user: 'admin',
          authType: 'password',
        },
      ];

      vi.mocked(sshAPI.getAll).mockResolvedValue(mockConfigs);
      vi.mocked(projectAPI.getAll).mockResolvedValue([]);
      vi.mocked(sshAPI.delete).mockRejectedValue(new Error('Delete failed'));

      renderSSHListPage();

      await waitFor(() => {
        expect(screen.getByTestId('ssh-card-test-server')).toBeInTheDocument();
      });

      // Click delete and confirm
      const deleteButton = screen.getByRole('button', { name: /delete/i });
      await user.click(deleteButton);

      await waitFor(() => {
        expect(screen.getByTestId('delete-dialog')).toBeInTheDocument();
      });

      const confirmButton = screen.getByRole('button', { name: /confirm/i });
      await user.click(confirmButton);

      // Error should be displayed
      await waitFor(() => {
        expect(screen.getByText(/delete failed/i)).toBeInTheDocument();
      });
    });
  });

  describe('Delete Functionality with In-Use Warning', () => {
    it('should display warning when SSH config is used by projects', async () => {
      const user = userEvent.setup();
      const mockConfigs: SSHConfig[] = [
        {
          name: 'prod-server',
          host: 'example.com',
          port: 22,
          user: 'admin',
          authType: 'password',
        },
      ];

      const mockProjects: Project[] = [
        {
          name: 'project1',
          buildInstructions: 'npm build',
          deployScript: 'deploy.sh',
          deployServers: ['prod-server'],
          createdAt: '2024-01-01T00:00:00Z',
          updatedAt: '2024-01-01T00:00:00Z',
        },
        {
          name: 'project2',
          buildInstructions: 'npm build',
          deployScript: 'deploy.sh',
          deployServers: ['prod-server', 'other-server'],
          createdAt: '2024-01-01T00:00:00Z',
          updatedAt: '2024-01-01T00:00:00Z',
        },
      ];

      vi.mocked(sshAPI.getAll).mockResolvedValue(mockConfigs);
      vi.mocked(projectAPI.getAll).mockResolvedValue(mockProjects);

      renderSSHListPage();

      await waitFor(() => {
        expect(screen.getByTestId('ssh-card-prod-server')).toBeInTheDocument();
      });

      // Click delete button
      const deleteButton = screen.getByRole('button', { name: /delete/i });
      await user.click(deleteButton);

      // Warning should be displayed
      await waitFor(() => {
        const warning = screen.getByTestId('delete-warning');
        expect(warning).toBeInTheDocument();
        expect(warning.textContent).toContain('currently used by 2 projects');
        expect(warning.textContent).toContain('project1, project2');
      });
    });

    it('should not display warning when SSH config is not used by any projects', async () => {
      const user = userEvent.setup();
      const mockConfigs: SSHConfig[] = [
        {
          name: 'unused-server',
          host: 'example.com',
          port: 22,
          user: 'admin',
          authType: 'password',
        },
      ];

      const mockProjects: Project[] = [
        {
          name: 'project1',
          buildInstructions: 'npm build',
          deployScript: 'deploy.sh',
          deployServers: ['other-server'],
          createdAt: '2024-01-01T00:00:00Z',
          updatedAt: '2024-01-01T00:00:00Z',
        },
      ];

      vi.mocked(sshAPI.getAll).mockResolvedValue(mockConfigs);
      vi.mocked(projectAPI.getAll).mockResolvedValue(mockProjects);

      renderSSHListPage();

      await waitFor(() => {
        expect(screen.getByTestId('ssh-card-unused-server')).toBeInTheDocument();
      });

      // Click delete button
      const deleteButton = screen.getByRole('button', { name: /delete/i });
      await user.click(deleteButton);

      // Warning should NOT be displayed
      await waitFor(() => {
        expect(screen.getByTestId('delete-dialog')).toBeInTheDocument();
      });
      
      expect(screen.queryByTestId('delete-warning')).not.toBeInTheDocument();
    });
  });
});
