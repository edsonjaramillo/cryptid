import { createRootRoute, Outlet } from '@tanstack/react-router';
import { Toaster } from '../components/ui/toast';

export const Route = createRootRoute({
  head: () => ({
    meta: [{ title: 'Hyde | Encryption Helper' }],
  }),
  component: () => (
    <>
      <Outlet />
      <Toaster />
    </>
  ),
});
