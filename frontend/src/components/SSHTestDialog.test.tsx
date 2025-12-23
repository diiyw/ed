import { describe, it, expect, vi, beforeEach } from 'vitest';
import { render, screen, waitFor, cleanup } from '@testing-library/react';
import * as fc from 'fast-check';
import { SSHTestDialog } from './SSHTestDialog';
import { sshAPI } from '@/services/api';
import type { SSHConfig, SSHTestResult } from '@/types';

// Mock the API module
vi.mock('@/services/api', () => ({
  sshAPI: {
    test: vi.fn(),
  },
}));

describe('SSHTestDialog', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  /**
   * Feature: frontend-refactor, Property 13: Loading indicator during async operations
   * 
   * Property: For any asynchronous API call in progress, a loading indicator
   * should be displayed to the user.
   * 
   * Validates: Requirements 6.2
   */
  describe('Property 13: Loading indicator during async operations', () => {
    it('should display loading indicator while testing SSH connection', async () => {
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
            // Setup: Mock API to return a delayed response
            let resolveTest: (value: SSHTestResult) => void;
            const testPromise = new Promise<SSHTestResult>((resolve) => {
              resolveTest = resolve;
            });
            vi.mocked(sshAPI.test).mockReturnValue(testPromise);

            // Action: Open the dialog
            render(
              <SSHTestDialog
                sshConfig={config}
                isOpen={true}
                onClose={() => {}}
              />
            );

            // Assert: Loading indicator should be visible while API call is in progress
            await waitFor(() => {
              const spinners = document.querySelectorAll('.animate-spin');
              expect(spinners.length).toBeGreaterThan(0);
            });

            // Resolve the promise to complete the test
            resolveTest!({ success: true, message: 'Connected' });

            // Wait for loading to complete
            await waitFor(() => {
              const spinners = document.querySelectorAll('.animate-spin');
              expect(spinners.length).toBe(0);
            });

            // Cleanup for next iteration
            cleanup();
            vi.clearAllMocks();
          }
        ),
        { numRuns: 10 }
      );
    }, 30000);
  });

  /**
   * Feature: frontend-refactor, Property 14: Test result display
   * 
   * Property: For any SSH connection test result (success or failure),
   * the appropriate message and details should be displayed to the user.
   * 
   * Validates: Requirements 6.3, 6.4
   */
  describe('Property 14: Test result display', () => {
    it('should display success message and details for successful connection tests', async () => {
      await fc.assert(
        fc.asyncProperty(
          fc.record({
            name: fc.string({ minLength: 1, maxLength: 50 }).filter(s => s.trim().length > 0),
            host: fc.string({ minLength: 1, maxLength: 100 }).filter(s => s.trim().length > 0),
            port: fc.integer({ min: 1, max: 65535 }),
            user: fc.string({ minLength: 1, maxLength: 50 }).filter(s => s.trim().length > 0),
            authType: fc.constantFrom('password', 'key', 'agent'),
          }),
          fc.record({
            success: fc.constant(true),
            message: fc.string({ minLength: 2, maxLength: 200 }).filter(s => s.trim().length >= 2),
            output: fc.option(fc.string({ minLength: 2 }).filter(s => s.trim().length >= 2), { nil: undefined }),
          }),
          async (config: SSHConfig, testResult: SSHTestResult) => {
            // Setup: Mock API to return success result
            vi.mocked(sshAPI.test).mockResolvedValue(testResult);

            // Action: Open the dialog
            render(
              <SSHTestDialog
                sshConfig={config}
                isOpen={true}
                onClose={() => {}}
              />
            );

            // Assert: Success message should be displayed
            await waitFor(() => {
              const successMessages = screen.queryAllByText(/connection successful/i);
              expect(successMessages.length).toBeGreaterThan(0);
            });

            // Assert: Result message should be displayed within the success container
            await waitFor(() => {
              const successContainers = document.querySelectorAll('[class*="bg-success"]');
              expect(successContainers.length).toBeGreaterThan(0);
              // Get the last (most recent) success container
              const lastSuccessContainer = successContainers[successContainers.length - 1];
              expect(lastSuccessContainer?.textContent).toContain(testResult.message.trim());
              
              // Assert: Output should be displayed if present
              if (testResult.output) {
                expect(lastSuccessContainer?.textContent).toContain(testResult.output.trim());
              }
            });

            // Cleanup for next iteration
            cleanup();
            vi.clearAllMocks();
          }
        ),
        { numRuns: 10 }
      );
    }, 30000);

    it('should display error message and details for failed connection tests', async () => {
      await fc.assert(
        fc.asyncProperty(
          fc.record({
            name: fc.string({ minLength: 1, maxLength: 50 }).filter(s => s.trim().length > 0),
            host: fc.string({ minLength: 1, maxLength: 100 }).filter(s => s.trim().length > 0),
            port: fc.integer({ min: 1, max: 65535 }),
            user: fc.string({ minLength: 1, maxLength: 50 }).filter(s => s.trim().length > 0),
            authType: fc.constantFrom('password', 'key', 'agent'),
          }),
          fc.record({
            success: fc.constant(false),
            message: fc.string({ minLength: 2, maxLength: 200 }).filter(s => s.trim().length >= 2),
          }),
          async (config: SSHConfig, testResult: SSHTestResult) => {
            // Setup: Mock API to return failure result
            vi.mocked(sshAPI.test).mockResolvedValue(testResult);

            // Action: Open the dialog
            render(
              <SSHTestDialog
                sshConfig={config}
                isOpen={true}
                onClose={() => {}}
              />
            );

            // Assert: Failure message should be displayed
            await waitFor(() => {
              const failureMessages = screen.queryAllByText(/connection failed/i);
              expect(failureMessages.length).toBeGreaterThan(0);
            });

            // Assert: Error message should be displayed within the error container
            await waitFor(() => {
              const errorContainer = document.querySelector('[class*="bg-error"]');
              expect(errorContainer).not.toBeNull();
              expect(errorContainer?.textContent).toContain(testResult.message.trim());
            });

            // Cleanup for next iteration
            cleanup();
            vi.clearAllMocks();
          }
        ),
        { numRuns: 10 }
      );
    }, 30000);
  });

  describe('Unit Tests', () => {
    describe('Loading State', () => {
      it('should display loading spinner during connection test', async () => {
        const mockConfig: SSHConfig = {
          name: 'test-server',
          host: 'example.com',
          port: 22,
          user: 'admin',
          authType: 'password',
        };

        // Mock API to return a delayed response
        let resolveTest: (value: SSHTestResult) => void;
        const testPromise = new Promise<SSHTestResult>((resolve) => {
          resolveTest = resolve;
        });
        vi.mocked(sshAPI.test).mockReturnValue(testPromise);

        const { unmount } = render(
          <SSHTestDialog
            sshConfig={mockConfig}
            isOpen={true}
            onClose={() => {}}
          />
        );

        // Loading indicator should be visible
        await waitFor(() => {
          const spinners = document.querySelectorAll('.animate-spin');
          expect(spinners.length).toBeGreaterThan(0);
        });

        // Resolve the promise
        resolveTest!({ success: true, message: 'Connected' });

        // Wait for loading to complete
        await waitFor(() => {
          const spinners = document.querySelectorAll('.animate-spin');
          expect(spinners.length).toBe(0);
        });

        unmount();
      });
    });

    describe('Success Result Display', () => {
      it('should display success message when connection test succeeds', async () => {
        const mockConfig: SSHConfig = {
          name: 'test-server',
          host: 'example.com',
          port: 22,
          user: 'admin',
          authType: 'password',
        };

        const mockResult: SSHTestResult = {
          success: true,
          message: 'Successfully connected to server',
          output: 'SSH connection established',
        };

        vi.mocked(sshAPI.test).mockResolvedValue(mockResult);

        render(
          <SSHTestDialog
            sshConfig={mockConfig}
            isOpen={true}
            onClose={() => {}}
          />
        );

        await waitFor(() => {
          expect(screen.getByText(/connection successful/i)).toBeInTheDocument();
        });

        expect(screen.getByText('Successfully connected to server')).toBeInTheDocument();
        expect(screen.getByText('SSH connection established')).toBeInTheDocument();
      });

      it('should display success message without output when output is not provided', async () => {
        const mockConfig: SSHConfig = {
          name: 'test-server',
          host: 'example.com',
          port: 22,
          user: 'admin',
          authType: 'key',
        };

        const mockResult: SSHTestResult = {
          success: true,
          message: 'Connection verified',
        };

        vi.mocked(sshAPI.test).mockResolvedValue(mockResult);

        render(
          <SSHTestDialog
            sshConfig={mockConfig}
            isOpen={true}
            onClose={() => {}}
          />
        );

        await waitFor(() => {
          expect(screen.getByText(/connection successful/i)).toBeInTheDocument();
        });

        expect(screen.getByText('Connection verified')).toBeInTheDocument();
      });
    });

    describe('Error Result Display', () => {
      it('should display error message when connection test fails', async () => {
        const mockConfig: SSHConfig = {
          name: 'test-server',
          host: 'example.com',
          port: 22,
          user: 'admin',
          authType: 'password',
        };

        const mockResult: SSHTestResult = {
          success: false,
          message: 'Authentication failed: Invalid credentials',
        };

        vi.mocked(sshAPI.test).mockResolvedValue(mockResult);

        render(
          <SSHTestDialog
            sshConfig={mockConfig}
            isOpen={true}
            onClose={() => {}}
          />
        );

        await waitFor(() => {
          expect(screen.getByText(/connection failed/i)).toBeInTheDocument();
        });

        expect(screen.getByText('Authentication failed: Invalid credentials')).toBeInTheDocument();
      });

      it('should display error when API call throws exception', async () => {
        const mockConfig: SSHConfig = {
          name: 'test-server',
          host: 'example.com',
          port: 22,
          user: 'admin',
          authType: 'password',
        };

        vi.mocked(sshAPI.test).mockRejectedValue(new Error('Network timeout'));

        render(
          <SSHTestDialog
            sshConfig={mockConfig}
            isOpen={true}
            onClose={() => {}}
          />
        );

        await waitFor(() => {
          expect(screen.getByText(/^error$/i)).toBeInTheDocument();
        });

        expect(screen.getByText('Network timeout')).toBeInTheDocument();
      });
    });

    describe('Dialog Behavior', () => {
      it('should display connection details in the dialog', async () => {
        const mockConfig: SSHConfig = {
          name: 'prod-server',
          host: 'production.example.com',
          port: 2222,
          user: 'deploy',
          authType: 'key',
        };

        vi.mocked(sshAPI.test).mockResolvedValue({
          success: true,
          message: 'Connected',
        });

        render(
          <SSHTestDialog
            sshConfig={mockConfig}
            isOpen={true}
            onClose={() => {}}
          />
        );

        await waitFor(() => {
          expect(screen.getByText(/testing connection: prod-server/i)).toBeInTheDocument();
        });

        expect(screen.getByText('production.example.com')).toBeInTheDocument();
        expect(screen.getByText('2222')).toBeInTheDocument();
        expect(screen.getByText('deploy')).toBeInTheDocument();
        // Check for auth type in the specific context (there are multiple "key" texts)
        const authLabels = screen.getAllByText('key');
        expect(authLabels.length).toBeGreaterThan(0);
      });

      it('should not render when sshConfig is null', () => {
        const { container } = render(
          <SSHTestDialog
            sshConfig={null}
            isOpen={true}
            onClose={() => {}}
          />
        );

        expect(container.firstChild).toBeNull();
      });
    });
  });
});
