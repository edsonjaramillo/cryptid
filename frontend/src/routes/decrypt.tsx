import { createFileRoute } from '@tanstack/react-router';
import { DecryptForm } from '../components/decrypt-form';
import { H1 } from '../components/ui/text';

export const Route = createFileRoute('/decrypt')({
  component: RouteComponent,
});

function RouteComponent() {
  return (
    <div className="mx-auto flex h-dvh w-responsive flex-col items-center justify-center gap-8">
      <H1 textColor="white">Decrypt</H1>
      <DecryptForm />
    </div>
  );
}
