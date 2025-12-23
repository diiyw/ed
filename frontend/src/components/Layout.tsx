import { Link, Outlet, useLocation } from 'react-router-dom';
import { Server, FolderGit2, Home } from 'lucide-react';
import { cn } from '@/lib/utils';

export function Layout() {
  const location = useLocation();

  const navItems = [
    { to: '/', label: 'Home', icon: Home },
    { to: '/ssh', label: 'SSH Configs', icon: Server },
    { to: '/projects', label: 'Projects', icon: FolderGit2 },
  ];

  return (
    <div className="min-h-screen bg-background">
      {/* Header */}
      <header className="border-b border-border bg-card">
        <div className="container mx-auto px-4 py-4">
          <div className="flex items-center justify-between">
            <div className="flex items-center space-x-2">
              <div className="text-2xl">🚀</div>
              <h1 className="text-2xl font-bold text-foreground">Easy Deploy</h1>
            </div>
            <nav className="flex space-x-1">
              {navItems.map((item) => {
                const Icon = item.icon;
                const isActive = location.pathname === item.to;
                return (
                  <Link
                    key={item.to}
                    to={item.to}
                    className={cn(
                      'flex items-center space-x-2 rounded-md px-4 py-2 text-sm font-medium transition-colors',
                      isActive
                        ? 'bg-primary text-primary-foreground'
                        : 'text-muted-foreground hover:bg-muted hover:text-foreground'
                    )}
                  >
                    <Icon className="h-4 w-4" />
                    <span>{item.label}</span>
                  </Link>
                );
              })}
            </nav>
          </div>
        </div>
      </header>

      {/* Main Content */}
      <main className="container mx-auto px-4 py-8">
        <Outlet />
      </main>

      {/* Footer */}
      <footer className="border-t border-border bg-card mt-auto">
        <div className="container mx-auto px-4 py-4 text-center text-sm text-muted-foreground">
          Easy Deploy - SSH Deployment Manager
        </div>
      </footer>
    </div>
  );
}
