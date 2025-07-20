import { memo, forwardRef, ReactNode } from 'react';
import { Link, LinkProps } from 'react-router-dom';
import { cn } from '@/shared/lib/utils';
import { isExternalRoute } from '@/shared/constants/routes';

export enum AppLinkTheme {
  PRIMARY = 'primary',
  SECONDARY = 'secondary',
}

interface AppLinkProps extends LinkProps {
  className?: string;
  theme?: AppLinkTheme;
  children?: ReactNode;
}

export const AppLink = memo(forwardRef<HTMLAnchorElement, AppLinkProps>(({
	children,
	className,
	to,
	theme = AppLinkTheme.PRIMARY,
	...otherProps
}, ref) => {
  if (isExternalRoute(to.toString())) {
    return (
      <a
        ref={ref}
        href={to.toString()}
        className={cn(className, {}, [theme])}
        {...otherProps}
      >
        {children}
      </a>
    );
  }

  return (
    <Link
      ref={ref}
      to={to}
      className={cn(className, {}, [theme])}
      {...otherProps}
    >
      {children}
    </Link>
  );
}));