import type { VariantProps } from 'class-variance-authority';
import type { JSX } from 'react';
import { cva } from 'class-variance-authority';
import { cn } from '../../lib/cn';

type ButtonVariants = VariantProps<typeof buttonVariants>;

interface RequiredButtonProps {
  type: React.ComponentProps<'button'>['type'];
}

export type ButtonProps = React.ComponentProps<'button'> & RequiredButtonProps & ButtonVariants;

export const buttonVariants = cva(
  'flex h-9 cursor-pointer items-center justify-center gap-2 rounded-base px-4 py-2 text-sm font-medium whitespace-nowrap transition-colors duration-base focus-visible:ring-2 focus-visible:ring-info focus-visible:ring-offset-2 focus-visible:outline-none [&_svg]:pointer-events-none [&_svg]:size-4 [&_svg]:shrink-0',
  {
    variants: {
      color: {
        primary:
          'bg-primary text-white not-disabled:hover:bg-primary-hover not-disabled:focus:bg-primary-hover not-disabled:active:bg-primary-pressed',
        success:
          'bg-success text-white not-disabled:hover:bg-success-hover not-disabled:focus:bg-success-hover not-disabled:active:bg-success-pressed',
        info: 'bg-info text-white not-disabled:hover:bg-info-hover not-disabled:focus:bg-info-hover not-disabled:active:bg-info-pressed',
        danger:
          'bg-danger text-white not-disabled:hover:bg-danger-hover not-disabled:focus:bg-danger-hover not-disabled:active:bg-danger-pressed',
        warning:
          'bg-warning text-white not-disabled:hover:bg-warning-hover not-disabled:focus:bg-warning-hover not-disabled:active:bg-warning-pressed',
        none: 'bg-white text-black',
      },
      disabled: {
        true: 'cursor-not-allowed bg-muted text-black',
      },
      width: {
        fit: 'w-fit',
        full: 'w-full',
      },
    },
    defaultVariants: { color: 'none', width: 'fit', disabled: false },
  },
);

export function Button({
  type,
  className,
  children,
  color,
  disabled,
  width,
  ...props
}: ButtonProps): JSX.Element {
  const style = buttonVariants({ color, disabled, width });
  return (
    <button className={cn(style, className)} disabled={disabled} type={type} {...props}>
      {children}
    </button>
  );
}
