import type { FieldErrors } from 'react-hook-form';

import { standardSchemaResolver } from '@hookform/resolvers/standard-schema';
import { FormProvider, useForm } from 'react-hook-form';
import * as v from 'valibot';
import { Form } from '../components/ui/form';
import { Input, InputError, InputGroup } from '../components/ui/input';
import { Label } from '../components/ui/text';
import { Button } from './ui/button';
import { toast } from './ui/toast';

const encryptSchema = v.object({
  password: v.pipe(v.string(), v.minLength(1, 'Password is required')),
  file: v.pipe(v.instance(FileList), v.minLength(1, 'File is required')),
});

type EncryptFormData = v.InferOutput<typeof encryptSchema>;

export function EncryptForm() {
  const form = useForm<EncryptFormData>({
    resolver: standardSchemaResolver(encryptSchema),
    defaultValues: {
      password: 'asdf',
    },
  });

  async function onSubmit(formData: EncryptFormData) {
    console.warn(formData);
    toast({ status: 'success', title: 'Encrypted successfully' });
  }

  async function onError(errors: FieldErrors<EncryptFormData>) {
    console.error(errors);
    toast({ status: 'error', title: 'Check form requirements.' });
  }

  return (
    <FormProvider {...form}>
      <Form onSubmit={form.handleSubmit(onSubmit, onError)}>
        <InputGroup>
          <Label textColor="white" htmlFor="file">
            File(s)
          </Label>
          <Input
            type="file"
            id="file"
            field="file"
            placeholder="File"
            autoComplete="file"
            multiple
            required
          />
          <InputError field="file" />
        </InputGroup>
        <InputGroup>
          <Label textColor="white" htmlFor="password">
            Password
          </Label>
          <Input
            type="password"
            id="password"
            field="password"
            autoComplete="new-password"
            placeholder="Password"
          />
          <InputError field="password" />
        </InputGroup>
        <Button
          type="submit"
          color="primary"
          className="ml-auto block"
          disabled={form.formState.isSubmitting}
        >
          {form.formState.isSubmitting ? 'Encrypting...' : 'Encrypt'}
        </Button>
      </Form>
    </FormProvider>
  );
}
