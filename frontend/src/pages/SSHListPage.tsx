import { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { Plus, Loader2 } from 'lucide-react';
import { Button } from '@/components/ui/button';
import { SSHCard } from '@/components/SSHCard';
import { SSHTestDialog } from '@/components/SSHTestDialog';
import { DeleteConfirmDialog } from '@/components/DeleteConfirmDialog';
import { sshAPI, projectAPI } from '@/services/api';
import type { SSHConfig, Project } from '@/types';

export function SSHListPage() {
  const navigate = useNavigate();
  const [configs, setConfigs] = useState<SSHConfig[]>([]);
  const [projects, setProjects] = useState<Project[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [testingConfig, setTestingConfig] = useState<SSHConfig | null>(null);
  const [deletingConfig, setDeletingConfig] = useState<SSHConfig | null>(null);

  useEffect(() => {
    fetchData();
  }, []);

  const fetchData = async () => {
    try {
      setLoading(true);
      setError(null);
      const [configsData, projectsData] = await Promise.all([
        sshAPI.getAll(),
        projectAPI.getAll(),
      ]);
      setConfigs(configsData);
      setProjects(projectsData);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to fetch data');
    } finally {
      setLoading(false);
    }
  };

  const getProjectsUsingConfig = (configName: string): Project[] => {
    return projects.filter((project) =>
      project.deployServers.includes(configName)
    );
  };

  const handleEdit = (config: SSHConfig) => {
    navigate(`/ssh/edit/${encodeURIComponent(config.name)}`);
  };

  const handleDelete = (config: SSHConfig) => {
    setDeletingConfig(config);
  };

  const confirmDelete = async () => {
    if (!deletingConfig) return;

    try {
      await sshAPI.delete(deletingConfig.name);
      setConfigs(configs.filter((c) => c.name !== deletingConfig.name));
      setDeletingConfig(null);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to delete SSH configuration');
      setDeletingConfig(null);
    }
  };

  const handleTest = (config: SSHConfig) => {
    setTestingConfig(config);
  };

  const getDeleteWarning = (): string | undefined => {
    if (!deletingConfig) return undefined;
    
    const projectsInUse = getProjectsUsingConfig(deletingConfig.name);
    if (projectsInUse.length === 0) return undefined;

    const projectNames = projectsInUse.map((p) => p.name).join(', ');
    return `This SSH configuration is currently used by ${projectsInUse.length} project${
      projectsInUse.length > 1 ? 's' : ''
    }: ${projectNames}. Deleting it may affect these projects.`;
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
          <Button onClick={fetchData} className="mt-4">
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
          <h1 className="text-3xl font-bold text-foreground">SSH Configurations</h1>
          <p className="text-muted-foreground mt-1">
            Manage your server connections and authentication
          </p>
        </div>
        <Button onClick={() => navigate('/ssh/new')} className="bg-primary hover:bg-primary/90 text-primary-foreground">
          <Plus className="h-4 w-4 mr-2" />
          Add SSH Config
        </Button>
      </div>

      {/* Empty State */}
      {configs.length === 0 ? (
        <div className="bg-card border border-border rounded-lg p-12 text-center">
          <div className="text-6xl mb-4">🗄️</div>
          <h2 className="text-2xl font-semibold text-foreground mb-2">No SSH Configurations Found</h2>
          <p className="text-muted-foreground mb-6 max-w-md mx-auto">
            You haven't added any SSH server configurations yet. Add your first configuration to get started.
          </p>
          <Button onClick={() => navigate('/ssh/new')} className="bg-primary hover:bg-primary/90 text-primary-foreground">
            <Plus className="h-4 w-4 mr-2" />
            Add Your First SSH Config
          </Button>
        </div>
      ) : (
        /* Config Grid */
        <div className="grid md:grid-cols-2 lg:grid-cols-3 gap-6">
          {configs.map((config) => (
            <SSHCard
              key={config.name}
              sshConfig={config}
              onEdit={handleEdit}
              onDelete={handleDelete}
              onTest={handleTest}
            />
          ))}
        </div>
      )}

      {/* Test Dialog */}
      <SSHTestDialog
        sshConfig={testingConfig}
        isOpen={testingConfig !== null}
        onClose={() => setTestingConfig(null)}
      />

      {/* Delete Confirmation Dialog */}
      <DeleteConfirmDialog
        title="Delete SSH Configuration"
        message={`Are you sure you want to delete "${deletingConfig?.name}"? This action cannot be undone.`}
        warning={getDeleteWarning()}
        isOpen={deletingConfig !== null}
        onConfirm={confirmDelete}
        onCancel={() => setDeletingConfig(null)}
      />
    </div>
  );
}
