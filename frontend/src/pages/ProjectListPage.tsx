import { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { Plus, Loader2 } from 'lucide-react';
import { Button } from '@/components/ui/button';
import { ProjectCard } from '@/components/ProjectCard';
import { DeleteConfirmDialog } from '@/components/DeleteConfirmDialog';
import { projectAPI } from '@/services/api';
import type { Project } from '@/types';

export function ProjectListPage() {
  const navigate = useNavigate();
  const [projects, setProjects] = useState<Project[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [deletingProject, setDeletingProject] = useState<Project | null>(null);

  useEffect(() => {
    fetchProjects();
  }, []);

  const fetchProjects = async () => {
    try {
      setLoading(true);
      setError(null);
      const data = await projectAPI.getAll();
      setProjects(data);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to fetch projects');
    } finally {
      setLoading(false);
    }
  };

  const handleEdit = (project: Project) => {
    navigate(`/projects/edit/${encodeURIComponent(project.name)}`);
  };

  const handleDelete = (project: Project) => {
    setDeletingProject(project);
  };

  const confirmDelete = async () => {
    if (!deletingProject) return;

    try {
      await projectAPI.delete(deletingProject.name);
      setProjects(projects.filter((p) => p.name !== deletingProject.name));
      setDeletingProject(null);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to delete project');
      setDeletingProject(null);
    }
  };

  const handleDeploy = (project: Project) => {
    navigate(`/deploy/${encodeURIComponent(project.name)}`);
  };

  if (loading) {
    return (
      <div className="flex items-center justify-center min-h-[400px]">
        <Loader2 className="h-8 w-8 animate-spin text-primary" />
      </div>
    );
  }

  if (error) {
    return (
      <div className="max-w-6xl mx-auto">
        <div className="bg-destructive/10 border border-destructive/30 rounded-lg p-6 text-center">
          <p className="text-destructive font-semibold mb-2">Error</p>
          <p className="text-card-foreground">{error}</p>
          <Button onClick={fetchProjects} className="mt-4">
            Try Again
          </Button>
        </div>
      </div>
    );
  }

  return (
    <div className="max-w-6xl mx-auto space-y-6">
      {/* Header */}
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-3xl font-bold text-foreground">Projects</h1>
          <p className="text-muted-foreground mt-1">
            Manage your deployment projects and workflows
          </p>
        </div>
        <Button onClick={() => navigate('/projects/new')} className="bg-secondary hover:bg-secondary/90 text-secondary-foreground">
          <Plus className="h-4 w-4 mr-2" />
          Add Project
        </Button>
      </div>

      {/* Empty State */}
      {projects.length === 0 ? (
        <div className="bg-card border border-border rounded-lg p-12 text-center">
          <div className="text-6xl mb-4">📦</div>
          <h2 className="text-2xl font-semibold text-foreground mb-2">No Projects Found</h2>
          <p className="text-muted-foreground mb-6 max-w-md mx-auto">
            You haven't created any projects yet. Create your first project to start deploying.
          </p>
          <Button onClick={() => navigate('/projects/new')} className="bg-secondary hover:bg-secondary/90 text-secondary-foreground">
            <Plus className="h-4 w-4 mr-2" />
            Create Your First Project
          </Button>
        </div>
      ) : (
        /* Project Grid */
        <div className="grid md:grid-cols-2 lg:grid-cols-3 gap-6">
          {projects.map((project) => (
            <ProjectCard
              key={project.name}
              project={project}
              onEdit={handleEdit}
              onDelete={handleDelete}
              onDeploy={handleDeploy}
            />
          ))}
        </div>
      )}

      {/* Delete Confirmation Dialog */}
      <DeleteConfirmDialog
        title="Delete Project"
        message={`Are you sure you want to delete "${deletingProject?.name}"? This action cannot be undone.`}
        isOpen={deletingProject !== null}
        onConfirm={confirmDelete}
        onCancel={() => setDeletingProject(null)}
      />
    </div>
  );
}
