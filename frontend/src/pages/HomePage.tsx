import { Link } from 'react-router-dom';
import { Server, FolderGit2, ArrowRight } from 'lucide-react';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Button } from '@/components/ui/button';

export function HomePage() {
  return (
    <div className="max-w-6xl mx-auto space-y-8">
      {/* Welcome Section */}
      <div className="text-center space-y-4">
        <div className="text-6xl">🚀</div>
        <h1 className="text-4xl font-bold text-foreground">Welcome to Easy Deploy</h1>
        <p className="text-xl text-muted-foreground max-w-2xl mx-auto">
          Manage your SSH configurations and deploy projects to multiple servers with ease
        </p>
      </div>

      {/* Quick Actions */}
      <div className="grid md:grid-cols-2 gap-6 mt-12">
        {/* SSH Configurations Card */}
        <Card className="bg-card border-border hover:border-primary transition-all hover:shadow-lg hover:shadow-primary/20">
          <CardHeader>
            <div className="flex items-center space-x-3">
              <div className="p-3 bg-primary/10 rounded-lg">
                <Server className="h-8 w-8 text-primary" />
              </div>
              <div>
                <CardTitle className="text-foreground text-2xl">SSH Configurations</CardTitle>
                <CardDescription className="text-muted-foreground">
                  Manage your server connections
                </CardDescription>
              </div>
            </div>
          </CardHeader>
          <CardContent className="space-y-4">
            <p className="text-card-foreground">
              Add, edit, delete, and test SSH server configurations. Configure authentication
              using passwords, SSH keys, or SSH agent.
            </p>
            <Link to="/ssh">
              <Button className="w-full bg-primary hover:bg-primary/90 text-primary-foreground">
                Manage SSH Configs
                <ArrowRight className="ml-2 h-4 w-4" />
              </Button>
            </Link>
          </CardContent>
        </Card>

        {/* Projects Card */}
        <Card className="bg-card border-border hover:border-secondary transition-all hover:shadow-lg hover:shadow-secondary/20">
          <CardHeader>
            <div className="flex items-center space-x-3">
              <div className="p-3 bg-secondary/10 rounded-lg">
                <FolderGit2 className="h-8 w-8 text-secondary" />
              </div>
              <div>
                <CardTitle className="text-foreground text-2xl">Projects</CardTitle>
                <CardDescription className="text-muted-foreground">
                  Deploy your applications
                </CardDescription>
              </div>
            </div>
          </CardHeader>
          <CardContent className="space-y-4">
            <p className="text-card-foreground">
              Create projects with build instructions and deploy scripts. Deploy to multiple
              servers simultaneously with real-time logs.
            </p>
            <Link to="/projects">
              <Button className="w-full bg-secondary hover:bg-secondary/90 text-secondary-foreground">
                Manage Projects
                <ArrowRight className="ml-2 h-4 w-4" />
              </Button>
            </Link>
          </CardContent>
        </Card>
      </div>

      {/* Features Section */}
      <div className="mt-16">
        <h2 className="text-2xl font-bold text-foreground text-center mb-8">Features</h2>
        <div className="grid md:grid-cols-3 gap-6">
          <div className="text-center space-y-2">
            <div className="text-4xl">🔐</div>
            <h3 className="text-lg font-semibold text-foreground">Secure Authentication</h3>
            <p className="text-sm text-muted-foreground">
              Support for password, SSH key, and SSH agent authentication
            </p>
          </div>
          <div className="text-center space-y-2">
            <div className="text-4xl">⚡</div>
            <h3 className="text-lg font-semibold text-foreground">Fast Deployment</h3>
            <p className="text-sm text-muted-foreground">
              Deploy to multiple servers simultaneously with real-time feedback
            </p>
          </div>
          <div className="text-center space-y-2">
            <div className="text-4xl">📊</div>
            <h3 className="text-lg font-semibold text-foreground">Live Logs</h3>
            <p className="text-sm text-muted-foreground">
              Monitor deployment progress with streaming logs in real-time
            </p>
          </div>
        </div>
      </div>
    </div>
  );
}
