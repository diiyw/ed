import { BrowserRouter, Routes, Route } from 'react-router-dom';
import { Layout } from './components/Layout';
import { HomePage } from './pages/HomePage';
import { SSHListPage } from './pages/SSHListPage';
import { SSHFormPage } from './pages/SSHFormPage';
import { ProjectListPage } from './pages/ProjectListPage';
import { ProjectFormPage } from './pages/ProjectFormPage';
import { DeployPage } from './pages/DeployPage';
import { ErrorBoundary } from './components/ErrorBoundary';
import { Toaster } from './components/ui/toaster';
import { OfflineIndicator } from './components/OfflineIndicator';

function App() {
  return (
    <ErrorBoundary>
      <BrowserRouter>
        <Routes>
          <Route path="/" element={<Layout />}>
            <Route index element={<HomePage />} />
            <Route path="ssh" element={<SSHListPage />} />
            <Route path="ssh/new" element={<SSHFormPage />} />
            <Route path="ssh/edit/:name" element={<SSHFormPage />} />
            <Route path="projects" element={<ProjectListPage />} />
            <Route path="projects/new" element={<ProjectFormPage />} />
            <Route path="projects/edit/:name" element={<ProjectFormPage />} />
            <Route path="deploy/:name" element={<DeployPage />} />
            <Route path="*" element={<NotFoundPage />} />
          </Route>
        </Routes>
        <Toaster />
        <OfflineIndicator />
      </BrowserRouter>
    </ErrorBoundary>
  );
}

function NotFoundPage() {
  return (
    <div className="flex flex-col items-center justify-center min-h-[400px] space-y-4">
      <div className="text-6xl">🔍</div>
      <h1 className="text-4xl font-bold text-white">404 - Page Not Found</h1>
      <p className="text-gray-400">The page you're looking for doesn't exist.</p>
      <a
        href="/"
        className="inline-flex items-center justify-center rounded-md bg-primary px-4 py-2 text-sm font-medium text-white hover:bg-primary/90 transition-colors"
      >
        Go Home
      </a>
    </div>
  );
}

export default App;
