import { FolderGit2, Edit, Trash2, Rocket, Server } from 'lucide-react';
import { Card, CardContent, CardDescription, CardFooter, CardHeader, CardTitle } from './ui/card';
import { Button } from './ui/button';
import { Badge } from './ui/badge';
import type { Project } from '@/types';

interface ProjectCardProps {
  project: Project;
  onEdit: (project: Project) => void;
  onDelete: (project: Project) => void;
  onDeploy: (project: Project) => void;
}

export function ProjectCard({ project, onEdit, onDelete, onDeploy }: ProjectCardProps) {
  const truncate = (text: string, maxLength: number) => {
    if (text.length <= maxLength) return text;
    return text.substring(0, maxLength) + '...';
  };

  return (
    <Card className="bg-card border-border hover:border-secondary transition-colors">
      <CardHeader>
        <div className="flex items-start justify-between">
          <div className="flex items-center space-x-2">
            <FolderGit2 className="h-5 w-5 text-secondary" />
            <CardTitle className="text-foreground">{project.name}</CardTitle>
          </div>
          <Badge variant="success">Active</Badge>
        </div>
        <CardDescription className="text-muted-foreground">
          {project.deployServers.length} server{project.deployServers.length !== 1 ? 's' : ''} configured
        </CardDescription>
      </CardHeader>
      <CardContent className="space-y-3">
        {/* Deploy Servers */}
        <div>
          <div className="flex items-center space-x-2 mb-2">
            <Server className="h-4 w-4 text-primary" />
            <span className="text-sm text-muted-foreground">Servers:</span>
          </div>
          <div className="flex flex-wrap gap-2">
            {project.deployServers.length > 0 ? (
              project.deployServers.map((server) => (
                <Badge key={server} variant="outline" className="text-xs">
                  {server}
                </Badge>
              ))
            ) : (
              <span className="text-sm text-muted">No servers configured</span>
            )}
          </div>
        </div>

        {/* Build Instructions Preview */}
        {project.buildInstructions && (
          <div>
            <span className="text-sm text-muted-foreground">Build:</span>
            <p className="text-sm text-card-foreground font-mono mt-1 bg-background p-2 rounded">
              {truncate(project.buildInstructions, 60)}
            </p>
          </div>
        )}

        {/* Deploy Script Preview */}
        {project.deployScript && (
          <div>
            <span className="text-sm text-muted-foreground">Deploy:</span>
            <p className="text-sm text-card-foreground font-mono mt-1 bg-background p-2 rounded">
              {truncate(project.deployScript, 60)}
            </p>
          </div>
        )}
      </CardContent>
      <CardFooter className="flex justify-end space-x-2">
        <Button
          variant="default"
          size="sm"
          onClick={() => onDeploy(project)}
          className="bg-secondary hover:bg-secondary/80 text-secondary-foreground"
        >
          <Rocket className="h-4 w-4 mr-1" />
          Deploy
        </Button>
        <Button
          variant="ghost"
          size="sm"
          onClick={() => onEdit(project)}
          className="text-accent hover:text-accent/80"
        >
          <Edit className="h-4 w-4 mr-1" />
          Edit
        </Button>
        <Button
          variant="ghost"
          size="sm"
          onClick={() => onDelete(project)}
          className="text-destructive hover:text-destructive/80"
        >
          <Trash2 className="h-4 w-4 mr-1" />
          Delete
        </Button>
      </CardFooter>
    </Card>
  );
}
