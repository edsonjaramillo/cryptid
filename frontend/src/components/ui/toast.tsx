import type { JSX } from 'react';
import { cva } from 'class-variance-authority';
import { toast as sonnerToast, Toaster as SonnerToaster } from 'sonner';
import { cn } from '../../lib/cn';
import { Span } from './text';

export function Toaster(): JSX.Element {
  return <SonnerToaster position="bottom-right" richColors />;
}

export function toast(toast: Omit<ToastProps, 'id'>): string | number {
  return sonnerToast.custom(id => (
    <Toast description={toast.description} id={id} status={toast.status} title={toast.title} />
  ));
}

export const toastVariants = cva('', {
  variants: {
    backgroundColor: {
      success: 'bg-success-accent',
      error: 'bg-danger-accent',
      info: 'bg-info-accent',
      warning: 'bg-warning-accent',
    },
    textColor: {
      success: 'text-success-hover',
      error: 'text-danger-hover',
      info: 'text-info-hover',
      warning: 'text-warning-hover',
    },
    ring: {
      success: 'inset-ring-success/5',
      error: 'inset-ring-danger/5',
      info: 'inset-ring-info/5',
      warning: 'inset-ring-warning/5',
    },
  },
  compoundVariants: [
    {
      ring: ['success', 'error', 'info', 'warning'],
      className: 'inset-ring-2',
    },
  ],
  defaultVariants: {},
});

export function Toast(props: ToastProps): JSX.Element {
  const { id, title, description } = props;
  const backgroundStyle = toastVariants({ backgroundColor: props.status });
  const textColor = toastVariants({ textColor: props.status });
  const ring = toastVariants({ ring: props.status });

  return (
    <button
      className={cn(
        backgroundStyle,
        ring,
        'group/toast relative min-w-toast cursor-pointer appearance-none rounded-base p-4',
      )}
      onClick={() => sonnerToast.dismiss(id)}
      type="button"
    >
      <div className="flex items-center gap-2">
        {getIcon(props.status, textColor)}
        <div className="flex flex-col gap-2">
          <Span className={cn(textColor, 'text-left font-semibold')} size="xs">
            {title}
          </Span>
          {description && (
            <Span className={cn(textColor, 'text-left')} size="xs">
              {description}
            </Span>
          )}
        </div>
      </div>
      <ToastCloseButton backgroundStyle={backgroundStyle} ring={ring} textColor={textColor} />
    </button>
  );
}

interface ToastProps {
  id: string | number;
  status: 'success' | 'error' | 'info' | 'warning';
  title: string;
  description?: string;
}

interface ToastCloseButtonProps {
  backgroundStyle: string;
  textColor: string;
  ring: string;
}
function ToastCloseButton({
  backgroundStyle,
  textColor,
  ring,
}: ToastCloseButtonProps): JSX.Element {
  return (
    <div
      className={cn(
        'absolute -top-1 -left-1 rounded-full p-1 opacity-0 transition-opacity duration-base group-hover/toast:opacity-100',
        backgroundStyle,
        ring,
      )}
    >
      <svg
        className={cn('size-3 fill-white', textColor)}
        stroke="currentColor"
        strokeWidth={1.5}
        viewBox="0 0 24 24"
        xmlns="http://www.w3.org/2000/svg"
      >
        <path d="M6 18 18 6M6 6l12 12" strokeLinecap="round" strokeLinejoin="round" />
      </svg>
    </div>
  );
}

interface IconProps {
  textColor: string;
}
type ToastIconProps = React.ComponentProps<'svg'> & IconProps;

function ToastIcon({ textColor, children, ...props }: ToastIconProps): JSX.Element {
  return (
    <svg
      className={cn('size-5 shrink-0', textColor)}
      fill="currentColor"
      viewBox="0 0 24 24"
      xmlns="http://www.w3.org/2000/svg"
      {...props}
    >
      {children}
    </svg>
  );
}

function SuccessIcon({ textColor }: IconProps): JSX.Element {
  return (
    <ToastIcon textColor={textColor}>
      <path
        clipRule="evenodd"
        d="M2.25 12c0-5.385 4.365-9.75 9.75-9.75s9.75 4.365 9.75 9.75-4.365 9.75-9.75 9.75S2.25 17.385 2.25 12Zm13.36-1.814a.75.75 0 1 0-1.22-.872l-3.236 4.53L9.53 12.22a.75.75 0 0 0-1.06 1.06l2.25 2.25a.75.75 0 0 0 1.14-.094l3.75-5.25Z"
        fillRule="evenodd"
      />
    </ToastIcon>
  );
}

function InfoIcon({ textColor }: IconProps): JSX.Element {
  return (
    <ToastIcon textColor={textColor}>
      <path
        clipRule="evenodd"
        d="M2.25 12c0-5.385 4.365-9.75 9.75-9.75s9.75 4.365 9.75 9.75-4.365 9.75-9.75 9.75S2.25 17.385 2.25 12Zm8.706-1.442c1.146-.573 2.437.463 2.126 1.706l-.709 2.836.042-.02a.75.75 0 0 1 .67 1.34l-.04.022c-1.147.573-2.438-.463-2.127-1.706l.71-2.836-.042.02a.75.75 0 1 1-.671-1.34l.041-.022ZM12 9a.75.75 0 1 0 0-1.5.75.75 0 0 0 0 1.5Z"
        fillRule="evenodd"
      />
    </ToastIcon>
  );
}

function WarningIcon({ textColor }: IconProps): JSX.Element {
  return (
    <ToastIcon textColor={textColor}>
      <path
        clipRule="evenodd"
        d="M9.401 3.003c1.155-2 4.043-2 5.197 0l7.355 12.748c1.154 2-.29 4.5-2.599 4.5H4.645c-2.309 0-3.752-2.5-2.598-4.5L9.4 3.003ZM12 8.25a.75.75 0 0 1 .75.75v3.75a.75.75 0 0 1-1.5 0V9a.75.75 0 0 1 .75-.75Zm0 8.25a.75.75 0 1 0 0-1.5.75.75 0 0 0 0 1.5Z"
        fillRule="evenodd"
      />
    </ToastIcon>
  );
}

function ErrorIcon({ textColor }: IconProps): JSX.Element {
  return (
    <ToastIcon textColor={textColor}>
      <path
        clipRule="evenodd"
        d="M2.25 12c0-5.385 4.365-9.75 9.75-9.75s9.75 4.365 9.75 9.75-4.365 9.75-9.75 9.75S2.25 17.385 2.25 12ZM12 8.25a.75.75 0 0 1 .75.75v3.75a.75.75 0 0 1-1.5 0V9a.75.75 0 0 1 .75-.75Zm0 8.25a.75.75 0 1 0 0-1.5.75.75 0 0 0 0 1.5Z"
        fillRule="evenodd"
      />
    </ToastIcon>
  );
}

function getIcon(status: 'success' | 'error' | 'info' | 'warning', textColor: string): JSX.Element {
  switch (status) {
    case 'success':
      return <SuccessIcon textColor={textColor} />;
    case 'info':
      return <InfoIcon textColor={textColor} />;
    case 'warning':
      return <WarningIcon textColor={textColor} />;
    case 'error':
      return <ErrorIcon textColor={textColor} />;
  }
}
