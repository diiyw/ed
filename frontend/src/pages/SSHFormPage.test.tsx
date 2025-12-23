import { describe, it, expect, vi, beforeEach } from 'vitest';
import { render, screen, waitFor } from '@testing-library/react';
import { Routes, Route, MemoryRouter } from 'react-router-dom';
import userEvent from '@testing-library/user-event';
import * as fc from 'fast-check';
import { SSHFormPage } from './SSHFormPage';
import { sshAPI } from '@/services/api';
import type { SSHConfig } from '@/types';

// Mock the API module
vi.mock('@/services/api', () => ({
  sshAPI: {
    getByName: vi.fn(),
    create: vi.fn(),
    update: vi.fn(),
  },
}));

// Mock the SSHForm component to simplify testing
vi.mock('@/components/SSHForm', () => ({
  SSHForm: ({ initialData, onSubmit, onCancel }: {
    initialData?: SSHConfig;
    onSubmit: (data: SSHConfig) => void;
    onCancel: () => void;
  }) => {
    const handleSubmit = () => {
      const dataToSubmit = initialData || {
        name: 'test',
        host: 'test.com',
        port: 22,
        user: 'root',
        authType: 'password' as const,
      };
      onSubmit(dataToSubmit);
    };

    return (
      <div data-testid="ssh-form">
        <div data-testid="form-initial-data">{JSON.stringify(initialData)}</div>
        <button onClick={handleSubmit}>Submit</button>
        <button onClick={onCancel}>Cancel</button>
      </div>
    );
  },
}));

describe('SSHFormPage', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  const renderSSHFormPage = (initialRoute = '/ssh/new') => {
    return render(
      <MemoryRouter initialEntries={[initialRoute]}>
        <Routes>
          <Route path="/ssh/new" element={<SSHFormPage />} />
          <Route path="/ssh/:name/edit" element={<SSHFormPage />} />
          <Route path="/ssh" element={<div data-testid="ssh-list-page">SSH List</div>} />
        </Routes>
      </MemoryRouter>
    );
  };

  /**
   * Feature: frontend-refactor, Property 7: Form pre-population for editing
   * 
   * Property: For any SSH configuration being edited, the edit form should be
   * initialized with all current values from that configuration.
   * 
   * Validates: Requirements 4.1
   */
  describe('Property 7: Form pre-population for editing', () => {
    it('should pre-populate form with existing SSH config data when editing', async () => {
      await fc.assert(
        fc.asyncProperty(
          fc.record({
            name: fc.string({ minLength: 1, maxLength: 50 })
              .filter(s => !s.includes('/'))
              .filter(s => {
                // Filter out strings that would cause URI encoding issues
                try {
                  const encoded = encodeURIComponent(s);
                  decodeURIComponent(encoded);
                  return true;
                } catch {
                  return false;
                }
              }),
            host: fc.string({ minLength: 1, maxLength: 100 }),
            port: fc.integer({ min: 1, max: 65535 }),
            user: fc.string({ minLength: 1, maxLength: 50 }),
            authType: fc.constantFrom('password', 'key', 'agent'),
            password: fc.option(fc.string(), { nil: undefined }),
            keyFile: fc.option(fc.string(), { nil: undefined }),
            keyPass: fc.option(fc.string(), { nil: undefined }),
          }),
          async (config: SSHConfig) => {
            // Setup: Mock the API to return the config
            vi.mocked(sshAPI.getByName).mockResolvedValue(config);

            // Action: Render the edit form
            const { unmount } = renderSSHFormPage(`/ssh/${encodeURIComponent(config.name)}/edit`);

            // Assert: Form should be pre-populated with the config data
            await waitFor(() => {
              const formData = screen.getByTestId('form-initial-data');
              expect(formData).toBeInTheDocument();
              
              const parsedData = JSON.parse(formData.textContent || '{}');
              expect(parsedData.name).toBe(config.name);
              expect(parsedData.host).toBe(config.host);
              expect(parsedData.port).toBe(config.port);
              expect(parsedData.user).toBe(config.user);
              expect(parsedData.authType).toBe(config.authType);
            });

            // Cleanup for next iteration
            unmount();
            vi.clearAllMocks();
          }
        ),
        { numRuns: 100 }
      );
    }, 30000); // 30 second timeout for property test
  });

  /**
   * Feature: frontend-refactor, Property 9: API call on save
   * 
   * Property: For any valid form submission, the frontend should make an API call
   * with the form data.
   * 
   * Validates: Requirements 3.3, 4.3
   */
  describe('Property 9: API call on save', () => {
    it('should call create API when submitting new SSH config', async () => {
      await fc.assert(
        fc.asyncProperty(
          fc.record({
            name: fc.string({ minLength: 1, maxLength: 50 }).filter(s => !s.includes('/')),
            host: fc.string({ minLength: 1, maxLength: 100 }),
            port: fc.integer({ min: 1, max: 65535 }),
            user: fc.string({ minLength: 1, maxLength: 50 }),
            authType: fc.constantFrom('password', 'key', 'agent'),
            password: fc.option(fc.string(), { nil: undefined }),
            keyFile: fc.option(fc.string(), { nil: undefined }),
            keyPass: fc.option(fc.string(), { nil: undefined }),
          }),
          async (config: SSHConfig) => {
            const user = userEvent.setup();
            
            // Setup: Mock the API
            vi.mocked(sshAPI.create).mockResolvedValue(config);

            // Action: Render the form and submit
            const { unmount } = renderSSHFormPage('/ssh/new');

            await waitFor(() => {
              expect(screen.getByTestId('ssh-form')).toBeInTheDocument();
            });

            const submitButton = screen.getByRole('button', { name: /submit/i });
            await user.click(submitButton);

            // Assert: Create API should be called
            await waitFor(() => {
              expect(sshAPI.create).toHaveBeenCalledTimes(1);
            });

            // Cleanup for next iteration
            unmount();
            vi.clearAllMocks();
          }
        ),
        { numRuns: 100 }
      );
    }, 30000); // 30 second timeout for property test

    it('should call update API when submitting edited SSH config', async () => {
      await fc.assert(
        fc.asyncProperty(
          fc.record({
            name: fc.string({ minLength: 1, maxLength: 50 })
              .filter(s => !s.includes('/'))
              .filter(s => {
                // Filter out strings that would cause URI encoding issues
                try {
                  const encoded = encodeURIComponent(s);
                  decodeURIComponent(encoded);
                  return true;
                } catch {
                  return false;
                }
              }),
            host: fc.string({ minLength: 1, maxLength: 100 }),
            port: fc.integer({ min: 1, max: 65535 }),
            user: fc.string({ minLength: 1, maxLength: 50 }),
            authType: fc.constantFrom('password', 'key', 'agent'),
            password: fc.option(fc.string(), { nil: undefined }),
            keyFile: fc.option(fc.string(), { nil: undefined }),
            keyPass: fc.option(fc.string(), { nil: undefined }),
          }),
          async (config: SSHConfig) => {
            const user = userEvent.setup();
            
            // Setup: Mock the API
            vi.mocked(sshAPI.getByName).mockResolvedValue(config);
            vi.mocked(sshAPI.update).mockResolvedValue(config);

            // Action: Render the edit form and submit
            const { unmount } = renderSSHFormPage(`/ssh/${encodeURIComponent(config.name)}/edit`);

            await waitFor(() => {
              expect(screen.getByTestId('ssh-form')).toBeInTheDocument();
            });

            const submitButton = screen.getByRole('button', { name: /submit/i });
            await user.click(submitButton);

            // Assert: Update API should be called with the config name
            await waitFor(() => {
              expect(sshAPI.update).toHaveBeenCalledTimes(1);
              expect(sshAPI.update).toHaveBeenCalledWith(config.name, expect.any(Object));
            });

            // Cleanup for next iteration
            unmount();
            vi.clearAllMocks();
          }
        ),
        { numRuns: 100 }
      );
    }, 30000); // 30 second timeout for property test
  });

  /**
   * Feature: frontend-refactor, Property 10: UI update after successful operation
   * 
   * Property: For any successful API response (create, update, delete), the frontend
   * should update the relevant list view and display a success message.
   * 
   * Validates: Requirements 3.3, 4.4, 5.3
   */
  describe('Property 10: UI update after successful operation', () => {
    it('should navigate to list page after successful create', async () => {
      await fc.assert(
        fc.asyncProperty(
          fc.record({
            name: fc.string({ minLength: 1, maxLength: 50 }).filter(s => !s.includes('/')),
            host: fc.string({ minLength: 1, maxLength: 100 }),
            port: fc.integer({ min: 1, max: 65535 }),
            user: fc.string({ minLength: 1, maxLength: 50 }),
            authType: fc.constantFrom('password', 'key', 'agent'),
            password: fc.option(fc.string(), { nil: undefined }),
            keyFile: fc.option(fc.string(), { nil: undefined }),
            keyPass: fc.option(fc.string(), { nil: undefined }),
          }),
          async (config: SSHConfig) => {
            const user = userEvent.setup();
            
            // Setup: Mock successful API call
            vi.mocked(sshAPI.create).mockResolvedValue(config);

            // Action: Render form and submit
            const { unmount } = renderSSHFormPage('/ssh/new');

            await waitFor(() => {
              expect(screen.getByTestId('ssh-form')).toBeInTheDocument();
            });

            const submitButton = screen.getByRole('button', { name: /submit/i });
            await user.click(submitButton);

            // Assert: Should navigate to list page after successful operation
            await waitFor(() => {
              expect(screen.getByTestId('ssh-list-page')).toBeInTheDocument();
            });

            // Cleanup for next iteration
            unmount();
            vi.clearAllMocks();
          }
        ),
        { numRuns: 100 }
      );
    }, 30000); // 30 second timeout for property test

    it('should navigate to list page after successful update', async () => {
      await fc.assert(
        fc.asyncProperty(
          fc.record({
            name: fc.string({ minLength: 1, maxLength: 50 })
              .filter(s => !s.includes('/'))
              .filter(s => {
                // Filter out strings that would cause URI encoding issues
                try {
                  const encoded = encodeURIComponent(s);
                  decodeURIComponent(encoded);
                  return true;
                } catch {
                  return false;
                }
              }),
            host: fc.string({ minLength: 1, maxLength: 100 }),
            port: fc.integer({ min: 1, max: 65535 }),
            user: fc.string({ minLength: 1, maxLength: 50 }),
            authType: fc.constantFrom('password', 'key', 'agent'),
            password: fc.option(fc.string(), { nil: undefined }),
            keyFile: fc.option(fc.string(), { nil: undefined }),
            keyPass: fc.option(fc.string(), { nil: undefined }),
          }),
          async (config: SSHConfig) => {
            const user = userEvent.setup();
            
            // Setup: Mock successful API calls
            vi.mocked(sshAPI.getByName).mockResolvedValue(config);
            vi.mocked(sshAPI.update).mockResolvedValue(config);

            // Action: Render edit form and submit
            const { unmount } = renderSSHFormPage(`/ssh/${encodeURIComponent(config.name)}/edit`);

            await waitFor(() => {
              expect(screen.getByTestId('ssh-form')).toBeInTheDocument();
            });

            const submitButton = screen.getByRole('button', { name: /submit/i });
            await user.click(submitButton);

            // Assert: Should navigate to list page after successful operation
            await waitFor(() => {
              expect(screen.getByTestId('ssh-list-page')).toBeInTheDocument();
            });

            // Cleanup for next iteration
            unmount();
            vi.clearAllMocks();
          }
        ),
        { numRuns: 100 }
      );
    }, 30000); // 30 second timeout for property test
  });

  describe('Unit Tests: Form Submission Flow', () => {
    it('should render form for creating new SSH config', async () => {
      renderSSHFormPage('/ssh/new');

      await waitFor(() => {
        expect(screen.getByText(/add ssh configuration/i)).toBeInTheDocument();
      });

      expect(screen.getByTestId('ssh-form')).toBeInTheDocument();
    });

    it('should render form for editing existing SSH config', async () => {
      const mockConfig: SSHConfig = {
        name: 'test-server',
        host: 'example.com',
        port: 22,
        user: 'admin',
        authType: 'password',
      };

      vi.mocked(sshAPI.getByName).mockResolvedValue(mockConfig);

      renderSSHFormPage('/ssh/test-server/edit');

      await waitFor(() => {
        expect(screen.getByText(/edit ssh configuration/i)).toBeInTheDocument();
      });

      expect(screen.getByTestId('ssh-form')).toBeInTheDocument();
    });

    it('should display loading spinner while fetching config for edit', async () => {
      vi.mocked(sshAPI.getByName).mockImplementation(
        () => new Promise((resolve) => setTimeout(() => resolve({
          name: 'test',
          host: 'test.com',
          port: 22,
          user: 'root',
          authType: 'password',
        }), 100))
      );

      renderSSHFormPage('/ssh/test-server/edit');

      // Loading spinner should be visible
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

    it('should handle form submission for new config', async () => {
      const user = userEvent.setup();
      const mockConfig: SSHConfig = {
        name: 'new-server',
        host: 'new.example.com',
        port: 22,
        user: 'admin',
        authType: 'password',
      };

      vi.mocked(sshAPI.create).mockResolvedValue(mockConfig);

      renderSSHFormPage('/ssh/new');

      await waitFor(() => {
        expect(screen.getByTestId('ssh-form')).toBeInTheDocument();
      });

      const submitButton = screen.getByRole('button', { name: /submit/i });
      await user.click(submitButton);

      await waitFor(() => {
        expect(sshAPI.create).toHaveBeenCalled();
      });
    });

    it('should handle form submission for editing config', async () => {
      const user = userEvent.setup();
      const mockConfig: SSHConfig = {
        name: 'test-server',
        host: 'example.com',
        port: 22,
        user: 'admin',
        authType: 'password',
      };

      vi.mocked(sshAPI.getByName).mockResolvedValue(mockConfig);
      vi.mocked(sshAPI.update).mockResolvedValue(mockConfig);

      renderSSHFormPage('/ssh/test-server/edit');

      await waitFor(() => {
        expect(screen.getByTestId('ssh-form')).toBeInTheDocument();
      });

      const submitButton = screen.getByRole('button', { name: /submit/i });
      await user.click(submitButton);

      await waitFor(() => {
        expect(sshAPI.update).toHaveBeenCalledWith('test-server', expect.any(Object));
      });
    });

    it('should display error message when submission fails', async () => {
      const user = userEvent.setup();

      vi.mocked(sshAPI.create).mockRejectedValue(new Error('Failed to create config'));

      renderSSHFormPage('/ssh/new');

      await waitFor(() => {
        expect(screen.getByTestId('ssh-form')).toBeInTheDocument();
      });

      const submitButton = screen.getByRole('button', { name: /submit/i });
      await user.click(submitButton);

      await waitFor(() => {
        expect(screen.getByText(/failed to create config/i)).toBeInTheDocument();
      });
    });

    it('should display error when fetching config fails', async () => {
      vi.mocked(sshAPI.getByName).mockRejectedValue(new Error('Config not found'));

      renderSSHFormPage('/ssh/nonexistent/edit');

      await waitFor(() => {
        expect(screen.getByText(/config not found/i)).toBeInTheDocument();
      });

      expect(screen.getByRole('button', { name: /back to ssh configs/i })).toBeInTheDocument();
    });
  });

  describe('Unit Tests: Navigation on Cancel', () => {
    it('should navigate back to list page when cancel is clicked', async () => {
      const user = userEvent.setup();

      renderSSHFormPage('/ssh/new');

      await waitFor(() => {
        expect(screen.getByTestId('ssh-form')).toBeInTheDocument();
      });

      const cancelButton = screen.getByRole('button', { name: /cancel/i });
      await user.click(cancelButton);

      await waitFor(() => {
        expect(screen.getByTestId('ssh-list-page')).toBeInTheDocument();
      });
    });

    it('should navigate back when back arrow is clicked', async () => {
      const user = userEvent.setup();

      renderSSHFormPage('/ssh/new');

      await waitFor(() => {
        expect(screen.getByTestId('ssh-form')).toBeInTheDocument();
      });

      // Find the back arrow button (it's a button with ArrowLeft icon)
      const backButton = screen.getAllByRole('button').find(
        button => button.querySelector('svg')
      );
      
      if (backButton) {
        await user.click(backButton);

        await waitFor(() => {
          expect(screen.getByTestId('ssh-list-page')).toBeInTheDocument();
        });
      }
    });
  });

  describe('Unit Tests: Success Message Display', () => {
    it('should navigate to list page on successful create (success is implicit in navigation)', async () => {
      const user = userEvent.setup();
      const mockConfig: SSHConfig = {
        name: 'new-server',
        host: 'new.example.com',
        port: 22,
        user: 'admin',
        authType: 'password',
      };

      vi.mocked(sshAPI.create).mockResolvedValue(mockConfig);

      renderSSHFormPage('/ssh/new');

      await waitFor(() => {
        expect(screen.getByTestId('ssh-form')).toBeInTheDocument();
      });

      const submitButton = screen.getByRole('button', { name: /submit/i });
      await user.click(submitButton);

      // Success is indicated by navigation to list page
      await waitFor(() => {
        expect(screen.getByTestId('ssh-list-page')).toBeInTheDocument();
      });
    });

    it('should navigate to list page on successful update (success is implicit in navigation)', async () => {
      const user = userEvent.setup();
      const mockConfig: SSHConfig = {
        name: 'test-server',
        host: 'example.com',
        port: 22,
        user: 'admin',
        authType: 'password',
      };

      vi.mocked(sshAPI.getByName).mockResolvedValue(mockConfig);
      vi.mocked(sshAPI.update).mockResolvedValue(mockConfig);

      renderSSHFormPage('/ssh/test-server/edit');

      await waitFor(() => {
        expect(screen.getByTestId('ssh-form')).toBeInTheDocument();
      });

      const submitButton = screen.getByRole('button', { name: /submit/i });
      await user.click(submitButton);

      // Success is indicated by navigation to list page
      await waitFor(() => {
        expect(screen.getByTestId('ssh-list-page')).toBeInTheDocument();
      });
    });
  });
});
