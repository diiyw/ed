import { useForm } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import * as z from 'zod';
import { Button } from './ui/button';
import { Input } from './ui/input';
import { Label } from './ui/label';
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from './ui/select';
import type { SSHConfig } from '@/types';

const sshConfigSchema = z.object({
  name: z.string().min(1, 'Name is required'),
  host: z.string().min(1, 'Host is required'),
  port: z.number().min(1).max(65535, 'Port must be between 1 and 65535'),
  user: z.string().min(1, 'User is required'),
  authType: z.enum(['password', 'key', 'agent']),
  password: z.string().optional(),
  keyFile: z.string().optional(),
  keyPass: z.string().optional(),
});

type SSHConfigFormData = z.infer<typeof sshConfigSchema>;

interface SSHFormProps {
  initialData?: SSHConfig;
  onSubmit: (data: SSHConfig) => void;
  onCancel: () => void;
}

export function SSHForm({ initialData, onSubmit, onCancel }: SSHFormProps) {
  const {
    register,
    handleSubmit,
    watch,
    setValue,
    formState: { errors, isSubmitting },
  } = useForm<SSHConfigFormData>({
    resolver: zodResolver(sshConfigSchema),
    defaultValues: initialData || {
      name: '',
      host: '',
      port: 22,
      user: '',
      authType: 'password',
      password: '',
      keyFile: '',
      keyPass: '',
    },
  });

  const authType = watch('authType');

  const onFormSubmit = (data: SSHConfigFormData) => {
    onSubmit(data as SSHConfig);
  };

  return (
    <form onSubmit={handleSubmit(onFormSubmit)} className="space-y-6">
      {/* Name */}
      <div className="space-y-2">
        <Label htmlFor="name" className="text-white">
          Name *
        </Label>
        <Input
          id="name"
          {...register('name')}
          placeholder="my-server"
          className={errors.name ? 'border-error' : ''}
        />
        {errors.name && (
          <p className="text-sm text-error">{errors.name.message}</p>
        )}
      </div>

      {/* Host */}
      <div className="space-y-2">
        <Label htmlFor="host" className="text-white">
          Host *
        </Label>
        <Input
          id="host"
          {...register('host')}
          placeholder="example.com"
          className={errors.host ? 'border-error' : ''}
        />
        {errors.host && (
          <p className="text-sm text-error">{errors.host.message}</p>
        )}
      </div>

      {/* Port */}
      <div className="space-y-2">
        <Label htmlFor="port" className="text-white">
          Port *
        </Label>
        <Input
          id="port"
          type="number"
          {...register('port', { valueAsNumber: true })}
          placeholder="22"
          className={errors.port ? 'border-error' : ''}
        />
        {errors.port && (
          <p className="text-sm text-error">{errors.port.message}</p>
        )}
      </div>

      {/* User */}
      <div className="space-y-2">
        <Label htmlFor="user" className="text-white">
          User *
        </Label>
        <Input
          id="user"
          {...register('user')}
          placeholder="root"
          className={errors.user ? 'border-error' : ''}
        />
        {errors.user && (
          <p className="text-sm text-error">{errors.user.message}</p>
        )}
      </div>

      {/* Auth Type */}
      <div className="space-y-2">
        <Label htmlFor="authType" className="text-white">
          Authentication Type *
        </Label>
        <Select
          value={authType}
          onValueChange={(value) => setValue('authType', value as 'password' | 'key' | 'agent')}
        >
          <SelectTrigger>
            <SelectValue placeholder="Select auth type" />
          </SelectTrigger>
          <SelectContent>
            <SelectItem value="password">Password</SelectItem>
            <SelectItem value="key">SSH Key</SelectItem>
            <SelectItem value="agent">SSH Agent</SelectItem>
          </SelectContent>
        </Select>
      </div>

      {/* Conditional Fields Based on Auth Type */}
      {authType === 'password' && (
        <div className="space-y-2">
          <Label htmlFor="password" className="text-white">
            Password
          </Label>
          <Input
            id="password"
            type="password"
            {...register('password')}
            placeholder="••••••••"
          />
        </div>
      )}

      {authType === 'key' && (
        <>
          <div className="space-y-2">
            <Label htmlFor="keyFile" className="text-white">
              Key File Path
            </Label>
            <Input
              id="keyFile"
              {...register('keyFile')}
              placeholder="/path/to/private/key"
            />
          </div>
          <div className="space-y-2">
            <Label htmlFor="keyPass" className="text-white">
              Key Password (if encrypted)
            </Label>
            <Input
              id="keyPass"
              type="password"
              {...register('keyPass')}
              placeholder="••••••••"
            />
          </div>
        </>
      )}

      {/* Form Actions */}
      <div className="flex justify-end space-x-2 pt-4">
        <Button type="button" variant="outline" onClick={onCancel}>
          Cancel
        </Button>
        <Button type="submit" disabled={isSubmitting}>
          {isSubmitting ? 'Saving...' : 'Save'}
        </Button>
      </div>
    </form>
  );
}
