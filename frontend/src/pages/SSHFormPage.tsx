import { useState, useEffect } from 'react';
import { useNavigate, useParams } from 'react-router-dom';
import { ArrowLeft, Loader2 } from 'lucide-react';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { SSHForm } from '@/components/SSHForm';
import { sshAPI } from '@/services/api';
import type { SSHConfig } from '@/types';

export function SSHFormPage() {
  const navigate = useNavigate();
  const { name } = useParams<{ name: string }>();
  const isEdit = Boolean(name);

  const [initialData, setInitialData] = useState<SSHConfig | undefined>();
  const [loading, setLoading] = useState(isEdit);
  const [error, setError] = useState<string | null>(null);
  const [submitting, setSubmitting] = useState(false);

  useEffect(() => {
    if (isEdit && name) {
      try {
        fetchConfig(decodeURIComponent(name));
      } catch (error) {
        // If decoding fails, use the name as-is
        fetchConfig(name);
      }
    }
  }, [isEdit, name]);

  const fetchConfig = async (configName: string) => {
    try {
      setLoading(true);
      setError(null);
      const data = await sshAPI.getByName(configName);
      setInitialData(data);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to fetch SSH configuration');
    } finally {
      setLoading(false);
    }
  };

  const handleSubmit = async (data: SSHConfig) => {
    try {
      setSubmitting(true);
      setError(null);

      if (isEdit && name) {
        try {
          await sshAPI.update(decodeURIComponent(name), data);
        } catch (decodeError) {
          // If decoding fails, use the name as-is
          await sshAPI.update(name, data);
        }
      } else {
        await sshAPI.create(data);
      }

      navigate('/ssh');
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to save SSH configuration');
      setSubmitting(false);
    }
  };

  const handleCancel = () => {
    navigate('/ssh');
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
          <Button onClick={() => navigate('/ssh')} className="mt-4">
            Back to SSH Configs
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
          onClick={() => navigate('/ssh')}
          className="text-muted-foreground hover:text-foreground"
        >
          <ArrowLeft className="h-5 w-5" />
        </Button>
        <div>
          <h1 className="text-3xl font-bold text-foreground">
            {isEdit ? 'Edit SSH Configuration' : 'Add SSH Configuration'}
          </h1>
          <p className="text-muted-foreground mt-1">
            {isEdit ? 'Update your server connection details' : 'Configure a new server connection'}
          </p>
        </div>
      </div>

      {/* Form Card */}
      <Card className="bg-card border-border">
        <CardHeader>
          <CardTitle className="text-foreground">Connection Details</CardTitle>
        </CardHeader>
        <CardContent>
          {error && !isEdit && (
            <div className="bg-destructive/10 border border-destructive/30 rounded-lg p-4 mb-6">
              <p className="text-destructive text-sm">{error}</p>
            </div>
          )}
          <SSHForm
            initialData={initialData}
            onSubmit={handleSubmit}
            onCancel={handleCancel}
          />
        </CardContent>
      </Card>

      {submitting && (
        <div className="fixed inset-0 bg-black/50 flex items-center justify-center z-50">
          <div className="bg-card border border-border rounded-lg p-6 flex items-center space-x-4">
            <Loader2 className="h-6 w-6 animate-spin text-primary" />
            <span className="text-foreground">Saving configuration...</span>
          </div>
        </div>
      )}
    </div>
  );
}
