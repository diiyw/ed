import { useState, useEffect } from 'react';
import { Loader2, CheckCircle, XCircle } from 'lucide-react';
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
} from './ui/dialog';
import { Button } from './ui/button';
import { sshAPI } from '@/services/api';
import type { SSHConfig, SSHTestResult } from '@/types';

interface SSHTestDialogProps {
  sshConfig: SSHConfig | null;
  isOpen: boolean;
  onClose: () => void;
}

export function SSHTestDialog({ sshConfig, isOpen, onClose }: SSHTestDialogProps) {
  const [testing, setTesting] = useState(false);
  const [result, setResult] = useState<SSHTestResult | null>(null);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    if (isOpen && sshConfig) {
      testConnection();
    } else {
      // Reset state when dialog closes
      setResult(null);
      setError(null);
      setTesting(false);
    }
  }, [isOpen, sshConfig]);

  const testConnection = async () => {
    if (!sshConfig) return;

    setTesting(true);
    setResult(null);
    setError(null);

    try {
      const testResult = await sshAPI.test(sshConfig.name);
      setResult(testResult);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to test connection');
    } finally {
      setTesting(false);
    }
  };

  if (!sshConfig) return null;

  return (
    <Dialog open={isOpen} onOpenChange={onClose}>
      <DialogContent className="bg-background-light border-gray-600">
        <DialogHeader>
          <DialogTitle className="text-white">
            Testing Connection: {sshConfig.name}
          </DialogTitle>
          <DialogDescription className="text-gray-400">
            {sshConfig.user}@{sshConfig.host}:{sshConfig.port}
          </DialogDescription>
        </DialogHeader>

        <div className="space-y-4">
          {/* Connection Details */}
          <div className="bg-background-dark p-4 rounded-lg space-y-2 text-sm">
            <div className="flex justify-between">
              <span className="text-gray-400">Host:</span>
              <span className="text-white font-mono">{sshConfig.host}</span>
            </div>
            <div className="flex justify-between">
              <span className="text-gray-400">Port:</span>
              <span className="text-white font-mono">{sshConfig.port}</span>
            </div>
            <div className="flex justify-between">
              <span className="text-gray-400">User:</span>
              <span className="text-white font-mono">{sshConfig.user}</span>
            </div>
            <div className="flex justify-between">
              <span className="text-gray-400">Auth:</span>
              <span className="text-white font-mono">{sshConfig.authType}</span>
            </div>
          </div>

          {/* Testing Status */}
          {testing && (
            <div className="flex items-center justify-center space-x-2 py-4">
              <Loader2 className="h-6 w-6 animate-spin text-info" />
              <span className="text-info">Testing connection...</span>
            </div>
          )}

          {/* Success Result */}
          {!testing && result && result.success && (
            <div className="bg-success/10 border border-success/30 rounded-lg p-4">
              <div className="flex items-center space-x-2 mb-2">
                <CheckCircle className="h-5 w-5 text-success" />
                <span className="text-success font-semibold">Connection Successful</span>
              </div>
              <p className="text-sm text-gray-300">{result.message}</p>
              {result.output && (
                <pre className="mt-2 text-xs text-gray-400 bg-background-dark p-2 rounded overflow-x-auto">
                  {result.output}
                </pre>
              )}
            </div>
          )}

          {/* Error Result */}
          {!testing && result && !result.success && (
            <div className="bg-error/10 border border-error/30 rounded-lg p-4">
              <div className="flex items-center space-x-2 mb-2">
                <XCircle className="h-5 w-5 text-error" />
                <span className="text-error font-semibold">Connection Failed</span>
              </div>
              <p className="text-sm text-gray-300">{result.message}</p>
            </div>
          )}

          {/* Error */}
          {!testing && error && (
            <div className="bg-error/10 border border-error/30 rounded-lg p-4">
              <div className="flex items-center space-x-2 mb-2">
                <XCircle className="h-5 w-5 text-error" />
                <span className="text-error font-semibold">Error</span>
              </div>
              <p className="text-sm text-gray-300">{error}</p>
            </div>
          )}

          {/* Actions */}
          <div className="flex justify-end space-x-2">
            {!testing && (
              <>
                <Button variant="outline" onClick={testConnection}>
                  Test Again
                </Button>
                <Button onClick={onClose}>Close</Button>
              </>
            )}
          </div>
        </div>
      </DialogContent>
    </Dialog>
  );
}
