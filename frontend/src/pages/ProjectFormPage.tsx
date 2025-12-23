import { useState, useEffect } from 'react';
import { useNavigate, useParams } from 'react-router-dom';
import { ArrowLeft, Loader2 } from 'lucide-react';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { ProjectForm } from '@/components/ProjectForm';
import { projectAPI, sshAPI } from '@/services/api';
import type { Project, SSHConfig } from '@/types';

export function ProjectFormPage() {
  const navigate = useNavigate();
  const { name } = useParams<{ name: string }>();
  const isEdit = Boolean(name);

  const [initialData, setInitialData] = useState<Project | undefined>();
  const [availableServers, setAvailableServers] = useState<string[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [submitting, setSubmitting] = useState(false);

  const fetchData = async () => {
    try {
      setLoading(true);
      setError(null);

      // Fetch available SSH servers
      const sshConfigs = await sshAPI.getAll();
      setAvailableServers(sshConfigs.map((cfg: SSHConfig) => cfg.name));

      // Fetch project data if editing
      if (isEdit && name) {
        const data = await projectAPI.getByName(decodeURIComponent(name));
        setInitialData(data);
      }
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to fetch data');
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchData();
  }, [isEdit, name]);

  const handleSubmit = async (data: Project) => {
    try {
      setSubmitting(true);
      setError(null);

      if (isEdit && name) {
        await projectAPI.update(decodeURIComponent(name), data);
      } else {
        await projectAPI.create(data);
      }

      navigate('/projects');
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to save project');
      setSubmitting(false);
    }
  };

  const handleCancel = () => {
    navigate('/projects');
  };

  if (loading) {
    return (
      <div className="flex items-center justify-center min-h-[400px]">
        <Loader2 className="h-8 w-8 animate-spin text-primary" />
      </div>
    );
  }

  if (error && isEdit) {
    return (
      <div className="max-w-2xl mx-auto">
        <div className="bg-destructive/10 border border-destructive/30 rounded-lg p-6 text-center">
          <p className="text-destructive font-semibold mb-2">Error</p>
          <p className="text-card-foreground">{error}</p>
          <Button onClick={() => navigate('/projects')} className="mt-4">
            Back to Projects
          </Button>
        </div>
      </div>
    );
  }

  if (availableServers.length === 0 && !loading) {
    return (
      <div className="max-w-2xl mx-auto">
        <div className="bg-accent/10 border border-accent/30 rounded-lg p-6 text-center">
          <p className="text-accent font-semibold mb-2">No SSH Configurations</p>
          <p className="text-card-foreground mb-4">
            You need to add at least one SSH configuration before creating a project.
          </p>
          <Button onClick={() => navigate('/ssh/new')} className="bg-primary hover:bg-primary/90 text-primary-foreground">
            Add SSH Configuration
          </Button>
        </div>
      </div>
    );
  }

  return (
    <div className="max-w-2xl mx-auto space-y-6">
      {/* Header */}
      <div className="flex items-center space-x-4">
        <Button
          variant="ghost"
          size="icon"
          onClick={() => navigate('/projects')}
          className="text-muted-foreground hover:text-foreground"
        >
          <ArrowLeft className="h-5 w-5" />
        </Button>
        <div>
          <h1 className="text-3xl font-bold text-foreground">
            {isEdit ? 'Edit Project' : 'Add Project'}
          </h1>
          <p className="text-muted-foreground mt-1">
            {isEdit ? 'Update your project configuration' : 'Configure a new deployment project'}
          </p>
        </div>
      </div>

      {/* Form Card */}
      <Card className="bg-card border-border">
        <CardHeader>
          <CardTitle className="text-foreground">Project Details</CardTitle>
        </CardHeader>
        <CardContent>
          {error && !isEdit && (
            <div className="bg-destructive/10 border border-destructive/30 rounded-lg p-4 mb-6">
              <p className="text-destructive text-sm">{error}</p>
            </div>
          )}
          <ProjectForm
            initialData={initialData}
            availableServers={availableServers}
            onSubmit={handleSubmit}
            onCancel={handleCancel}
          />
        </CardContent>
      </Card>

      {submitting && (
        <div className="fixed inset-0 bg-black/50 flex items-center justify-center z-50">
          <div className="bg-card border border-border rounded-lg p-6 flex items-center space-x-4">
            <Loader2 className="h-6 w-6 animate-spin text-primary" />
            <span className="text-foreground">Saving project...</span>
          </div>
        </div>
      )}
    </div>
  );
}
