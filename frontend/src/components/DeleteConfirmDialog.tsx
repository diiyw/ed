import { AlertTriangle } from 'lucide-react';
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from './ui/dialog';
import { Button } from './ui/button';

interface DeleteConfirmDialogProps {
  title: string;
  message: string;
  isOpen: boolean;
  onConfirm: () => void;
  onCancel: () => void;
  warning?: string;
}

export function DeleteConfirmDialog({
  title,
  message,
  isOpen,
  onConfirm,
  onCancel,
  warning,
}: DeleteConfirmDialogProps) {
  return (
    <Dialog open={isOpen} onOpenChange={onCancel}>
      <DialogContent className="bg-background-light border-gray-600">
        <DialogHeader>
          <DialogTitle className="text-white flex items-center space-x-2">
            <AlertTriangle className="h-5 w-5 text-warning" />
            <span>{title}</span>
          </DialogTitle>
          <DialogDescription className="text-gray-400">
            {message}
          </DialogDescription>
        </DialogHeader>

        {warning && (
          <div className="bg-warning/10 border border-warning/30 rounded-lg p-4">
            <div className="flex items-start space-x-2">
              <AlertTriangle className="h-5 w-5 text-warning flex-shrink-0 mt-0.5" />
              <div>
                <p className="text-sm font-semibold text-warning mb-1">Warning</p>
                <p className="text-sm text-gray-300">{warning}</p>
              </div>
            </div>
          </div>
        )}

        <DialogFooter>
          <Button variant="outline" onClick={onCancel}>
            Cancel
          </Button>
          <Button variant="destructive" onClick={onConfirm}>
            Delete
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  );
}
