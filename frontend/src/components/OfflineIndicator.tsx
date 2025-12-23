import { useOnline } from '@/hooks/use-online';
import { WifiOff } from 'lucide-react';
import { useEffect } from 'react';
import { toast } from '@/hooks/use-toast';

export function OfflineIndicator() {
  const isOnline = useOnline();

  useEffect(() => {
    if (!isOnline) {
      toast({
        variant: 'destructive',
        title: 'No internet connection',
        description: 'You are currently offline. Some features may not work.',
      });
    } else {
      // Show a success toast when coming back online
      toast({
        variant: 'success',
        title: 'Back online',
        description: 'Your internet connection has been restored.',
      });
    }
  }, [isOnline]);

  if (isOnline) {
    return null;
  }

  return (
    <div className="fixed bottom-4 left-4 z-50 flex items-center gap-2 rounded-md bg-destructive px-4 py-2 text-sm text-destructive-foreground shadow-lg">
      <WifiOff className="h-4 w-4" />
      <span>You are offline</span>
    </div>
  );
}
