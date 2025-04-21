import { createFileRoute } from '@tanstack/react-router';
import { EncryptForm } from '../components/encrypt-form';
import { H1 } from '../components/ui/text';

export const Route = createFileRoute('/')({
  component: App,
});

function App() {
  return (
    <div className="h-dvh bg-black flex flex-col items-center justify-center gap-8">
      <H1 textColor="white">Hyde</H1>
      <EncryptForm />
    </div>
  );
}
