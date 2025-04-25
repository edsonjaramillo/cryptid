import type { FieldErrors } from 'react-hook-form';
import type { UniversalFormData } from '../utils/validation';
import { standardSchemaResolver } from '@hookform/resolvers/standard-schema';
import { FormProvider, useForm } from 'react-hook-form';
import { downloadFile } from '../lib/download';
import { universalSchema } from '../utils/validation';
import { Button } from './ui/button';
import { Form } from './ui/form';
import { Input, InputError, InputGroup } from './ui/input';
import { Label } from './ui/text';
import { toast } from './ui/toast';

export function DecryptForm() {
  const form = useForm<UniversalFormData>({ resolver: standardSchemaResolver(universalSchema) });

  async function onSubmit(data: UniversalFormData) {
    const formData = new FormData();
    formData.append('password', data.password);
    formData.append('file', data.file[0]);

    const response = await fetch('http://localhost:8080/decrypt', {
      method: 'POST',
      body: formData,
    });

    if (!response.ok) {
      toast({ status: 'error', title: response.statusText });
      return;
    }

    const blob = await response.blob();
    downloadFile('decrypt', blob, data.file[0].name);

    form.resetField('file');
    toast({ status: 'success', title: 'Encrypted file ready to download.' });
  }

  async function onError(_: FieldErrors<UniversalFormData>) {
    toast({ status: 'error', title: 'Check form requirements.' });
  }

  return (
    <FormProvider {...form}>
      <Form onSubmit={form.handleSubmit(onSubmit, onError)}>
        <InputGroup>
          <Label textColor="white" htmlFor="file">
            File
          </Label>
          <Input
            type="file"
            id="file"
            field="file"
            placeholder="File"
            autoComplete="off"
            required
          />
          <InputError field="file" />
        </InputGroup>
        <InputGroup>
          <Label textColor="white" htmlFor="password">
            Password
          </Label>
          <Input type="password" id="password" field="password" autoComplete="current-password" />
          <InputError field="password" />
        </InputGroup>
        <Button
          type="submit"
          color="primary"
          className="ml-auto block"
          disabled={form.formState.isSubmitting}
        >
          {form.formState.isSubmitting ? 'Decrypting...' : 'Decrypt'}
        </Button>
      </Form>
    </FormProvider>
  );
}
