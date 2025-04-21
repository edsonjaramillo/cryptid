import type { VariantProps } from 'class-variance-authority';
import type { JSX } from 'react';
import { cva } from 'class-variance-authority';
import { cn } from '../../lib/cn';

interface TextVariants extends VariantProps<typeof textVariants> {}
export const textVariants = cva('antialiased', {
  variants: {
    size: {
      'xs': 'text-xs',
      'sm': 'text-sm',
      'base': 'text-base',
      'lg': 'text-lg',
      'xl': 'text-xl',
      '2xl': 'text-2xl',
      '3xl': 'text-3xl',
      '4xl': 'text-4xl',
    },
    textColor: {
      black: 'text-black',
      gray: 'text-gray',
      muted: 'text-muted',
      white: 'text-white',
      primary: 'text-primary',
      success: 'text-success',
      info: 'text-info',
      warning: 'text-warning',
      danger: 'text-danger',
    },
    seperator: {
      true: 'border-b border-gray pb-2',
    },
  },
  compoundVariants: [
    {
      size: ['xl', '2xl', '3xl', '4xl'],
      className: 'scroll-m-20 tracking-tight',
    },
  ],
  defaultVariants: { size: 'base', textColor: 'black', seperator: false },
});

interface Header1Props extends React.ComponentProps<'h1'>, TextVariants {}
export function H1({
  className,
  children,
  size,
  textColor,
  seperator,
  ...props
}: Header1Props): JSX.Element {
  const style = cn(
    textVariants({
      seperator,
      size: size || '4xl',
      textColor,
      className: 'font-extrabold lg:text-5xl',
    }),
    className,
  );
  return (
    <h1 className={style} {...props}>
      {children}
    </h1>
  );
}

interface Header2Props extends React.ComponentProps<'h2'>, TextVariants {}
export function H2({
  className,
  children,
  size,
  textColor,
  seperator,
  ...props
}: Header2Props): JSX.Element {
  const style = cn(
    textVariants({
      seperator,
      size: size || '3xl',
      textColor,
      className: 'font-semibold first:mt-0',
    }),
    className,
  );
  return (
    <h2 className={style} {...props}>
      {children}
    </h2>
  );
}

interface Header3Props extends React.ComponentProps<'h3'>, TextVariants {}
export function H3({
  className,
  children,
  size,
  textColor,
  seperator,
  ...props
}: Header3Props): JSX.Element {
  const style = cn(
    textVariants({
      seperator,
      size: size || '2xl',
      textColor,
      className: 'font-semibold',
    }),
    className,
  );
  return (
    <h3 className={style} {...props}>
      {children}
    </h3>
  );
}

interface Header4Props extends React.ComponentProps<'h4'>, TextVariants {}
export function H4({
  className,
  children,
  size,
  textColor,
  seperator,
  ...props
}: Header4Props): JSX.Element {
  const style = cn(
    textVariants({
      seperator,
      size: size || 'xl',
      textColor,
      className: 'font-semibold',
    }),
    className,
  );
  return (
    <h4 className={style} {...props}>
      {children}
    </h4>
  );
}

interface ParagraphProps extends React.ComponentProps<'p'>, TextVariants {}
export function Paragraph({
  className,
  children,
  size,
  textColor,
  seperator,
  ...props
}: ParagraphProps): JSX.Element {
  const style = cn(
    textVariants({
      seperator,
      size: size || 'base',
      textColor,
      className: 'leading-7 [&:not(:first-child)]:mt-6',
    }),
    className,
  );
  return (
    <p className={style} {...props}>
      {children}
    </p>
  );
}

interface LabelProps extends React.ComponentProps<'label'>, TextVariants {}
export function Label({
  className,
  children,
  size,
  textColor,
  seperator,
  ...props
}: LabelProps): JSX.Element {
  const style = cn(
    textVariants({
      seperator,
      size: size || 'sm',
      textColor,
      className: 'leading-none font-medium',
    }),
    className,
  );
  return (
    <label className={style} {...props}>
      {children}
    </label>
  );
}

interface SpanProps extends React.ComponentProps<'span'>, TextVariants {}
export function Span({
  className,
  children,
  size,
  textColor,
  seperator,
  ...props
}: SpanProps): JSX.Element {
  const style = cn(textVariants({ textColor, size, seperator }), className);
  return (
    <span className={style} {...props}>
      {children}
    </span>
  );
}
