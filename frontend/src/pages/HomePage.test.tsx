import { describe, it, expect } from 'vitest';
import { render, screen } from '@testing-library/react';
import { BrowserRouter } from 'react-router-dom';
import { HomePage } from './HomePage';

describe('HomePage', () => {
  const renderHomePage = () => {
    return render(
      <BrowserRouter>
        <HomePage />
      </BrowserRouter>
    );
  };

  describe('Rendering', () => {
    it('should render without errors', () => {
      renderHomePage();
      
      // Check that the main heading is present
      expect(screen.getByText(/welcome to easy deploy/i)).toBeInTheDocument();
    });

    it('should display welcome message', () => {
      renderHomePage();
      
      expect(screen.getByText(/welcome to easy deploy/i)).toBeInTheDocument();
      expect(screen.getByText(/manage your ssh configurations and deploy projects/i)).toBeInTheDocument();
    });

    it('should display SSH configurations card', () => {
      renderHomePage();
      
      expect(screen.getByRole('heading', { name: /ssh configurations/i })).toBeInTheDocument();
      expect(screen.getByText(/manage your server connections/i)).toBeInTheDocument();
    });

    it('should display projects card', () => {
      renderHomePage();
      
      expect(screen.getByText(/^projects$/i)).toBeInTheDocument();
      expect(screen.getByText(/deploy your applications/i)).toBeInTheDocument();
    });

    it('should display features section', () => {
      renderHomePage();
      
      expect(screen.getByText(/^features$/i)).toBeInTheDocument();
      expect(screen.getByText(/secure authentication/i)).toBeInTheDocument();
      expect(screen.getByText(/fast deployment/i)).toBeInTheDocument();
      expect(screen.getByText(/live logs/i)).toBeInTheDocument();
    });
  });

  describe('Navigation Links', () => {
    it('should have navigation link to SSH configurations page', () => {
      renderHomePage();
      
      const sshLink = screen.getByRole('link', { name: /manage ssh configs/i });
      expect(sshLink).toBeInTheDocument();
      expect(sshLink).toHaveAttribute('href', '/ssh');
    });

    it('should have navigation link to projects page', () => {
      renderHomePage();
      
      const projectsLink = screen.getByRole('link', { name: /manage projects/i });
      expect(projectsLink).toBeInTheDocument();
      expect(projectsLink).toHaveAttribute('href', '/projects');
    });

    it('should have quick action buttons for both SSH and Projects', () => {
      renderHomePage();
      
      const sshButton = screen.getByRole('button', { name: /manage ssh configs/i });
      const projectsButton = screen.getByRole('button', { name: /manage projects/i });
      
      expect(sshButton).toBeInTheDocument();
      expect(projectsButton).toBeInTheDocument();
    });
  });

  describe('Layout and Structure', () => {
    it('should use card components for layout', () => {
      const { container } = renderHomePage();
      
      // Check that card elements are present (shadcn cards have specific structure)
      const cards = container.querySelectorAll('[class*="card"]');
      expect(cards.length).toBeGreaterThan(0);
    });

    it('should display icons for SSH and Projects sections', () => {
      renderHomePage();
      
      // Check for Lucide icons by their SVG elements
      const { container } = renderHomePage();
      const svgElements = container.querySelectorAll('svg');
      
      // Should have multiple SVG icons (Server, FolderGit2, ArrowRight icons)
      expect(svgElements.length).toBeGreaterThan(0);
    });
  });
});
