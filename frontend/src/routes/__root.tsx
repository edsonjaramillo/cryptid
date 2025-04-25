import { createRootRoute, Outlet } from '@tanstack/react-router';
import { Navigation } from '../components/navigation';
import { Toaster } from '../components/ui/toast';

export const Route = createRootRoute({
  head: () => ({ meta: [{ title: 'Hyde | Encryption Helper' }] }),
  component: Root,
});

function Root() {
  return (
    <>
      <Navigation />
      <main className="h-dvh bg-black">
        <Outlet />
      </main>
      <Toaster />
    </>
  );
}
