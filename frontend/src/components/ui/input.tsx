import type { JSX } from 'react';
import { useFormContext } from 'react-hook-form';
import { cn } from '../../lib/cn';
import { Span } from './text';

interface InputRequiredProps {
  id: string;
  type: React.ComponentProps<'input'>['type'];
  field: string;
  autoComplete: React.ComponentProps<'input'>['autoComplete'];
}

type InputProps = React.ComponentProps<'input'> & InputRequiredProps;
export function Input({
  autoComplete,
  className,
  field,
  id,
  required,
  type,
  ...props
}: InputProps): JSX.Element {
  const { register } = useFormContext();
  const style = cn(
    'h-9 w-full rounded-base border border-muted bg-transparent px-3 py-2 text-sm text-white shadow-base placeholder:text-gray focus-visible:ring-2 focus-visible:ring-info focus-visible:ring-offset-2 focus-visible:outline-0 disabled:cursor-not-allowed disabled:text-muted disabled:placeholder:text-muted',
    className,
  );
  return (
    <input
      autoComplete={autoComplete}
      className={style}
      id={id}
      type={type}
      {...register(field, { required })}
      {...props}
    />
  );
}

type InputGroupProps = React.ComponentProps<'div'>;
export function InputGroup({
  children,
  className,
  ...props
}: InputGroupProps): JSX.Element {
  const style = cn('grid w-full max-w-form items-center gap-2', className);
  return (
    <div className={style} {...props}>
      {children}
    </div>
  );
}

type InputColumnsProps = React.ComponentProps<'div'>;
export function InputColumns({
  children,
  className,
  ...props
}: InputColumnsProps): JSX.Element {
  const style = cn('grid grid-cols-2 gap-4', className);
  return (
    <div className={style} {...props}>
      {children}
    </div>
  );
}

type InputDescriptionProps = React.ComponentProps<'span'>;
export function InputDescription({
  className,
  children,
  ...props
}: InputDescriptionProps): JSX.Element {
  return (
    <Span
      className={cn(className, 'line-clamp-1')}
      size="sm"
      textColor="gray"
      {...props}
    >
      {children}
    </Span>
  );
}

interface InputErrorRequiredProps {
  field: string;
}
type InputErrorProps = React.ComponentProps<'span'> & InputErrorRequiredProps;
export function InputError({
  className,
  field,
  ...props
}: InputErrorProps): JSX.Element | undefined {
  const { formState } = useFormContext();
  const { errors } = formState;
  const error = errors?.[field];

  if (!error || !error.message) {
    return undefined;
  }

  return (
    <Span
      className={cn('line-clamp-1 font-semibold', className)}
      size="sm"
      textColor="danger"
      {...props}
    >
      {error.message.toString()}
    </Span>
  );
}
