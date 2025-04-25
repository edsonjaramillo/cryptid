import { createFileRoute } from '@tanstack/react-router';
import { EncryptForm } from '../components/encrypt-form';
import { H1 } from '../components/ui/text';

export const Route = createFileRoute('/encrypt')({
  component: RouteComponent,
});

function RouteComponent() {
  return (
    <div className="mx-auto flex h-dvh w-responsive flex-col items-center justify-center gap-8">
      <H1 textColor="white">Encrypt</H1>
      <EncryptForm />
    </div>
  );
}
