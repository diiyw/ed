import { Server, Edit, Trash2, TestTube } from 'lucide-react';
import { Card, CardContent, CardDescription, CardFooter, CardHeader, CardTitle } from './ui/card';
import { Button } from './ui/button';
import { Badge } from './ui/badge';
import type { SSHConfig } from '@/types';

interface SSHCardProps {
  sshConfig: SSHConfig;
  onEdit: (config: SSHConfig) => void;
  onDelete: (config: SSHConfig) => void;
  onTest: (config: SSHConfig) => void;
}

export function SSHCard({ sshConfig, onEdit, onDelete, onTest }: SSHCardProps) {
  return (
    <Card className="bg-card border-border hover:border-primary transition-colors">
      <CardHeader>
        <div className="flex items-start justify-between">
          <div className="flex items-center space-x-2">
            <Server className="h-5 w-5 text-primary" />
            <CardTitle className="text-foreground">{sshConfig.name}</CardTitle>
          </div>
          <Badge variant="outline" className="text-muted-foreground">
            {sshConfig.authType}
          </Badge>
        </div>
        <CardDescription className="text-muted-foreground">
          {sshConfig.user}@{sshConfig.host}:{sshConfig.port}
        </CardDescription>
      </CardHeader>
      <CardContent>
        <div className="space-y-2 text-sm text-card-foreground">
          <div className="flex justify-between">
            <span className="text-muted-foreground">Host:</span>
            <span className="font-mono">{sshConfig.host}</span>
          </div>
          <div className="flex justify-between">
            <span className="text-muted-foreground">Port:</span>
            <span className="font-mono">{sshConfig.port}</span>
          </div>
          <div className="flex justify-between">
            <span className="text-muted-foreground">User:</span>
            <span className="font-mono">{sshConfig.user}</span>
          </div>
        </div>
      </CardContent>
      <CardFooter className="flex justify-end space-x-2">
        <Button
          variant="ghost"
          size="sm"
          onClick={() => onTest(sshConfig)}
          className="text-primary hover:text-primary/80"
        >
          <TestTube className="h-4 w-4 mr-1" />
          Test
        </Button>
        <Button
          variant="ghost"
          size="sm"
          onClick={() => onEdit(sshConfig)}
          className="text-accent hover:text-accent/80"
        >
          <Edit className="h-4 w-4 mr-1" />
          Edit
        </Button>
        <Button
          variant="ghost"
          size="sm"
          onClick={() => onDelete(sshConfig)}
          className="text-destructive hover:text-destructive/80"
        >
          <Trash2 className="h-4 w-4 mr-1" />
          Delete
        </Button>
      </CardFooter>
    </Card>
  );
}
