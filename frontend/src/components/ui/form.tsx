import type { JSX } from 'react';
import { cn } from '../../lib/cn';

interface FormProps extends React.ComponentProps<'form'> {}

export function Form({ className, children, ...props }: FormProps): JSX.Element {
  return (
    <form className={cn('mx-auto w-full max-w-form space-y-4', className)} {...props} noValidate>
      {children}
    </form>
  );
}
