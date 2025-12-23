import { useForm, Controller } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import * as z from 'zod';
import { Button } from './ui/button';
import { Input } from './ui/input';
import { Label } from './ui/label';
import { Textarea } from './ui/textarea';
import { MultiSelect } from './ui/multi-select';
import type { Project } from '@/types';

const projectSchema = z.object({
  name: z.string().min(1, 'Name is required'),
  buildInstructions: z.string().min(1, 'Build instructions are required'),
  deployScript: z.string().min(1, 'Deploy script is required'),
  deployServers: z.array(z.string()).min(1, 'At least one server must be selected'),
  createdAt: z.string().optional(),
  updatedAt: z.string().optional(),
});

type ProjectFormData = z.infer<typeof projectSchema>;

interface ProjectFormProps {
  initialData?: Project;
  availableServers: string[];
  onSubmit: (data: Project) => void;
  onCancel: () => void;
}

export function ProjectForm({ initialData, availableServers, onSubmit, onCancel }: ProjectFormProps) {
  const {
    register,
    handleSubmit,
    control,
    formState: { errors, isSubmitting },
  } = useForm<ProjectFormData>({
    resolver: zodResolver(projectSchema),
    defaultValues: initialData || {
      name: '',
      buildInstructions: '',
      deployScript: '',
      deployServers: [],
      createdAt: new Date().toISOString(),
      updatedAt: new Date().toISOString(),
    },
  });

  const onFormSubmit = (data: ProjectFormData) => {
    const projectData: Project = {
      ...data,
      createdAt: initialData?.createdAt || new Date().toISOString(),
      updatedAt: new Date().toISOString(),
    };
    onSubmit(projectData);
  };

  return (
    <form onSubmit={handleSubmit(onFormSubmit)} className="space-y-6">
      {/* Name */}
      <div className="space-y-2">
        <Label htmlFor="name" className="text-white">
          Project Name *
        </Label>
        <Input
          id="name"
          {...register('name')}
          placeholder="my-project"
          className={errors.name ? 'border-error' : ''}
        />
        {errors.name && (
          <p className="text-sm text-error">{errors.name.message}</p>
        )}
      </div>

      {/* Build Instructions */}
      <div className="space-y-2">
        <Label htmlFor="buildInstructions" className="text-white">
          Build Instructions *
        </Label>
        <Textarea
          id="buildInstructions"
          {...register('buildInstructions')}
          placeholder="npm install && npm run build"
          rows={4}
          className={errors.buildInstructions ? 'border-error' : ''}
        />
        {errors.buildInstructions && (
          <p className="text-sm text-error">{errors.buildInstructions.message}</p>
        )}
        <p className="text-xs text-gray-400">
          Commands to build your project before deployment
        </p>
      </div>

      {/* Deploy Script */}
      <div className="space-y-2">
        <Label htmlFor="deployScript" className="text-white">
          Deploy Script *
        </Label>
        <Textarea
          id="deployScript"
          {...register('deployScript')}
          placeholder="rsync -avz ./dist/ user@server:/var/www/html/"
          rows={4}
          className={errors.deployScript ? 'border-error' : ''}
        />
        {errors.deployScript && (
          <p className="text-sm text-error">{errors.deployScript.message}</p>
        )}
        <p className="text-xs text-gray-400">
          Commands to deploy your project to the servers
        </p>
      </div>

      {/* Deploy Servers */}
      <div className="space-y-2">
        <Label htmlFor="deployServers" className="text-white">
          Deploy Servers *
        </Label>
        <Controller
          name="deployServers"
          control={control}
          render={({ field }) => (
            <MultiSelect
              options={availableServers}
              value={field.value}
              onChange={field.onChange}
              placeholder="Select servers to deploy to"
            />
          )}
        />
        {errors.deployServers && (
          <p className="text-sm text-error">{errors.deployServers.message}</p>
        )}
        <p className="text-xs text-gray-400">
          Select one or more SSH configurations to deploy to
        </p>
      </div>

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
