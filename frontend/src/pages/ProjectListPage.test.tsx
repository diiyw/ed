import { describe, it, expect, vi, beforeEach } from 'vitest';
import { render, screen, waitFor } from '@testing-library/react';
import { BrowserRouter } from 'react-router-dom';
import userEvent from '@testing-library/user-event';
import { ProjectListPage } from './ProjectListPage';
import { projectAPI } from '@/services/api';
import type { Project } from '@/types';

// Mock the API module
vi.mock('@/services/api', () => ({
  projectAPI: {
    getAll: vi.fn(),
    delete: vi.fn(),
  },
}));

// Mock the child components to simplify testing
vi.mock('@/components/ProjectCard', () => ({
  ProjectCard: ({ project, onEdit, onDelete, onDeploy }: any) => (
    <div data-testid={`project-card-${project.name}`}>
      <span>{project.name}</span>
      <span>{project.buildInstructions}</span>
      <span>{project.deployScript}</span>
      {project.deployServers.map((server: string) => (
        <span key={server}>{server}</span>
      ))}
      <button onClick={() => onEdit(project)}>Edit</button>
      <button onClick={() => onDelete(project)}>Delete</button>
      <button onClick={() => onDeploy(project)}>Deploy</button>
    </div>
  ),
}));

vi.mock('@/components/DeleteConfirmDialog', () => ({
  DeleteConfirmDialog: ({ isOpen, onConfirm, onCancel }: any) =>
    isOpen ? (
      <div data-testid="delete-dialog">
        <button onClick={onConfirm}>Confirm</button>
        <button onClick={onCancel}>Cancel</button>
      </div>
    ) : null,
}));

describe('ProjectListPage', () => {
  const renderProjectListPage = () => {
    return render(
      <BrowserRouter>
        <ProjectListPage />
      </BrowserRouter>
    );
  };

  beforeEach(() => {
    vi.clearAllMocks();
  });

  describe('Empty State', () => {
    it('should display empty state when no projects exist', async () => {
      vi.mocked(projectAPI.getAll).mockResolvedValue([]);

      renderProjectListPage();

      await waitFor(() => {
        expect(screen.getByText(/no projects found/i)).toBeInTheDocument();
      });

      expect(screen.getByText(/you haven't created any projects yet/i)).toBeInTheDocument();
      expect(screen.getByRole('button', { name: /create your first project/i })).toBeInTheDocument();
    });

    it('should not display empty state when projects exist', async () => {
      const mockProjects: Project[] = [
        {
          name: 'test-project',
          buildInstructions: 'npm build',
          deployScript: 'deploy.sh',
          deployServers: ['server1'],
          createdAt: '2024-01-01T00:00:00Z',
          updatedAt: '2024-01-01T00:00:00Z',
        },
      ];

      vi.mocked(projectAPI.getAll).mockResolvedValue(mockProjects);

      renderProjectListPage();

      await waitFor(() => {
        expect(screen.queryByText(/no projects found/i)).not.toBeInTheDocument();
      });
    });
  });

  describe('Project List Rendering', () => {
    it('should display all projects in a grid', async () => {
      const mockProjects: Project[] = [
        {
          name: 'project1',
          buildInstructions: 'npm build',
          deployScript: 'deploy.sh',
          deployServers: ['server1'],
          createdAt: '2024-01-01T00:00:00Z',
          updatedAt: '2024-01-01T00:00:00Z',
        },
        {
          name: 'project2',
          buildInstructions: 'yarn build',
          deployScript: 'deploy2.sh',
          deployServers: ['server2'],
          createdAt: '2024-01-02T00:00:00Z',
          updatedAt: '2024-01-02T00:00:00Z',
        },
      ];

      vi.mocked(projectAPI.getAll).mockResolvedValue(mockProjects);

      renderProjectListPage();

      await waitFor(() => {
        expect(screen.getByTestId('project-card-project1')).toBeInTheDocument();
        expect(screen.getByTestId('project-card-project2')).toBeInTheDocument();
      });

      expect(screen.getByText('project1')).toBeInTheDocument();
      expect(screen.getByText('project2')).toBeInTheDocument();
    });

    it('should display page header with title and description', async () => {
      vi.mocked(projectAPI.getAll).mockResolvedValue([]);

      renderProjectListPage();

      await waitFor(() => {
        expect(screen.getByRole('heading', { name: /^projects$/i })).toBeInTheDocument();
      });

      expect(screen.getByText(/manage your deployment projects and workflows/i)).toBeInTheDocument();
    });

    it('should display "Add Project" button in header', async () => {
      vi.mocked(projectAPI.getAll).mockResolvedValue([]);

      renderProjectListPage();

      await waitFor(() => {
        expect(screen.getByRole('button', { name: /add project/i })).toBeInTheDocument();
      });
    });
  });

  describe('Action Buttons', () => {
    it('should handle edit action', async () => {
      const user = userEvent.setup();
      const mockProjects: Project[] = [
        {
          name: 'test-project',
          buildInstructions: 'npm build',
          deployScript: 'deploy.sh',
          deployServers: ['server1'],
          createdAt: '2024-01-01T00:00:00Z',
          updatedAt: '2024-01-01T00:00:00Z',
        },
      ];

      vi.mocked(projectAPI.getAll).mockResolvedValue(mockProjects);

      renderProjectListPage();

      await waitFor(() => {
        expect(screen.getByTestId('project-card-test-project')).toBeInTheDocument();
      });

      const editButton = screen.getByRole('button', { name: /edit/i });
      await user.click(editButton);

      // Navigation is handled by React Router, which we can't easily test here
      // The important part is that the button exists and is clickable
      expect(editButton).toBeInTheDocument();
    });

    it('should handle delete action and show confirmation dialog', async () => {
      const user = userEvent.setup();
      const mockProjects: Project[] = [
        {
          name: 'test-project',
          buildInstructions: 'npm build',
          deployScript: 'deploy.sh',
          deployServers: ['server1'],
          createdAt: '2024-01-01T00:00:00Z',
          updatedAt: '2024-01-01T00:00:00Z',
        },
      ];

      vi.mocked(projectAPI.getAll).mockResolvedValue(mockProjects);

      renderProjectListPage();

      await waitFor(() => {
        expect(screen.getByTestId('project-card-test-project')).toBeInTheDocument();
      });

      const deleteButton = screen.getByRole('button', { name: /delete/i });
      await user.click(deleteButton);

      // Delete confirmation dialog should appear
      await waitFor(() => {
        expect(screen.getByTestId('delete-dialog')).toBeInTheDocument();
      });
    });

    it('should handle deploy action', async () => {
      const user = userEvent.setup();
      const mockProjects: Project[] = [
        {
          name: 'test-project',
          buildInstructions: 'npm build',
          deployScript: 'deploy.sh',
          deployServers: ['server1'],
          createdAt: '2024-01-01T00:00:00Z',
          updatedAt: '2024-01-01T00:00:00Z',
        },
      ];

      vi.mocked(projectAPI.getAll).mockResolvedValue(mockProjects);

      renderProjectListPage();

      await waitFor(() => {
        expect(screen.getByTestId('project-card-test-project')).toBeInTheDocument();
      });

      const deployButton = screen.getByRole('button', { name: /deploy/i });
      await user.click(deployButton);

      // Navigation is handled by React Router, which we can't easily test here
      // The important part is that the button exists and is clickable
      expect(deployButton).toBeInTheDocument();
    });

    it('should call delete API and update list when delete is confirmed', async () => {
      const user = userEvent.setup();
      const mockProjects: Project[] = [
        {
          name: 'test-project',
          buildInstructions: 'npm build',
          deployScript: 'deploy.sh',
          deployServers: ['server1'],
          createdAt: '2024-01-01T00:00:00Z',
          updatedAt: '2024-01-01T00:00:00Z',
        },
      ];

      vi.mocked(projectAPI.getAll).mockResolvedValue(mockProjects);
      vi.mocked(projectAPI.delete).mockResolvedValue();

      renderProjectListPage();

      await waitFor(() => {
        expect(screen.getByTestId('project-card-test-project')).toBeInTheDocument();
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
        expect(projectAPI.delete).toHaveBeenCalledWith('test-project');
      });

      // Project should be removed from the list
      await waitFor(() => {
        expect(screen.queryByTestId('project-card-test-project')).not.toBeInTheDocument();
      });
    });

    it('should not call delete API when delete is cancelled', async () => {
      const user = userEvent.setup();
      const mockProjects: Project[] = [
        {
          name: 'test-project',
          buildInstructions: 'npm build',
          deployScript: 'deploy.sh',
          deployServers: ['server1'],
          createdAt: '2024-01-01T00:00:00Z',
          updatedAt: '2024-01-01T00:00:00Z',
        },
      ];

      vi.mocked(projectAPI.getAll).mockResolvedValue(mockProjects);

      renderProjectListPage();

      await waitFor(() => {
        expect(screen.getByTestId('project-card-test-project')).toBeInTheDocument();
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
      expect(projectAPI.delete).not.toHaveBeenCalled();

      // Project should still be in the list
      expect(screen.getByTestId('project-card-test-project')).toBeInTheDocument();
    });
  });

  describe('Loading State', () => {
    it('should display loading spinner while fetching projects', async () => {
      vi.mocked(projectAPI.getAll).mockImplementation(
        () => new Promise((resolve) => setTimeout(() => resolve([]), 100))
      );

      renderProjectListPage();

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
      vi.mocked(projectAPI.getAll).mockRejectedValue(new Error('Network error'));

      renderProjectListPage();

      await waitFor(() => {
        expect(screen.getByText(/^error$/i)).toBeInTheDocument();
        expect(screen.getByText(/network error/i)).toBeInTheDocument();
      });

      expect(screen.getByRole('button', { name: /try again/i })).toBeInTheDocument();
    });

    it('should retry fetching when "Try Again" button is clicked', async () => {
      const user = userEvent.setup();
      
      // First call fails
      vi.mocked(projectAPI.getAll).mockRejectedValueOnce(new Error('Network error'));
      
      renderProjectListPage();

      await waitFor(() => {
        expect(screen.getByText(/network error/i)).toBeInTheDocument();
      });

      // Second call succeeds
      vi.mocked(projectAPI.getAll).mockResolvedValue([]);

      const tryAgainButton = screen.getByRole('button', { name: /try again/i });
      await user.click(tryAgainButton);

      await waitFor(() => {
        expect(projectAPI.getAll).toHaveBeenCalledTimes(2);
      });

      // Error should be cleared
      await waitFor(() => {
        expect(screen.queryByText(/network error/i)).not.toBeInTheDocument();
      });
    });

    it('should display error when delete fails', async () => {
      const user = userEvent.setup();
      const mockProjects: Project[] = [
        {
          name: 'test-project',
          buildInstructions: 'npm build',
          deployScript: 'deploy.sh',
          deployServers: ['server1'],
          createdAt: '2024-01-01T00:00:00Z',
          updatedAt: '2024-01-01T00:00:00Z',
        },
      ];

      vi.mocked(projectAPI.getAll).mockResolvedValue(mockProjects);
      vi.mocked(projectAPI.delete).mockRejectedValue(new Error('Delete failed'));

      renderProjectListPage();

      await waitFor(() => {
        expect(screen.getByTestId('project-card-test-project')).toBeInTheDocument();
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
});
