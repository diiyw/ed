import { describe, it, expect, vi, beforeEach } from 'vitest';
import { renderHook, act } from '@testing-library/react';
import { useToast, toast } from './use-toast';

describe('useToast', () => {
  beforeEach(() => {
    // Clear any existing toasts before each test
    const { result } = renderHook(() => useToast());
    act(() => {
      result.current.toasts.forEach((t) => result.current.dismiss(t.id));
    });
  });

  it('adds a toast when toast() is called', () => {
    const { result } = renderHook(() => useToast());

    act(() => {
      toast({
        title: 'Test toast',
        description: 'Test description',
      });
    });

    expect(result.current.toasts).toHaveLength(1);
    expect(result.current.toasts[0].title).toBe('Test toast');
    expect(result.current.toasts[0].description).toBe('Test description');
  });

  it('dismisses a toast by id', async () => {
    const { result } = renderHook(() => useToast());

    let toastId: string;
    act(() => {
      const t = toast({
        title: 'Test toast',
      });
      toastId = t.id;
    });

    // Wait for state to update
    await vi.waitFor(() => {
      expect(result.current.toasts.length).toBeGreaterThan(0);
    });

    act(() => {
      result.current.dismiss(toastId!);
    });

    // Toast should be marked as closed (open: false)
    const dismissedToast = result.current.toasts.find(t => t.id === toastId);
    expect(dismissedToast?.open).toBe(false);
  });

  it('limits the number of toasts to TOAST_LIMIT', () => {
    const { result } = renderHook(() => useToast());

    act(() => {
      // Add 10 toasts
      for (let i = 0; i < 10; i++) {
        toast({
          title: `Toast ${i}`,
        });
      }
    });

    // Should only keep the last 5 toasts (TOAST_LIMIT = 5)
    expect(result.current.toasts.length).toBeLessThanOrEqual(5);
  });

  it('supports different toast variants', () => {
    const { result } = renderHook(() => useToast());

    act(() => {
      toast({
        title: 'Error toast',
        variant: 'destructive',
      });
    });

    expect(result.current.toasts[0].variant).toBe('destructive');
  });

  it('allows updating a toast', () => {
    const { result } = renderHook(() => useToast());

    let toastInstance: ReturnType<typeof toast>;
    act(() => {
      toastInstance = toast({
        title: 'Original title',
      });
    });

    act(() => {
      toastInstance.update({
        id: toastInstance.id,
        title: 'Updated title',
      });
    });

    expect(result.current.toasts[0].title).toBe('Updated title');
  });
});
