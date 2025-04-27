import { Link } from '@tanstack/react-router';
import { buttonVariants } from './ui/button';
import { H1, H2 } from './ui/text';

export function Hero() {
  return (
    <div className="mx-auto flex h-dvh w-responsive flex-col items-center justify-center gap-8">
      <H1 textColor="white" className="text-center text-pretty">
        Secure Your Data
      </H1>
      <H2 textColor="gray" size="2xl" className="text-center text-pretty">
        Start using Hyde by choosing an action
      </H2>
      <div className="flex gap-4">
        <Link to="/encrypt" className={buttonVariants({ color: 'primary' })}>
          Encrypt
        </Link>
        <Link to="/decrypt" className={buttonVariants({ color: 'none' })}>
          Decrypt
        </Link>
      </div>
    </div>
  );
}
