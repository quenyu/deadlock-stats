import { memo, ReactNode } from 'react';
import { Link, LinkProps } from 'react-router-dom';
import { cn } from '@/shared/lib/utils';

export enum AppLinkTheme {
  PRIMARY = 'primary',
  SECONDARY = 'secondary',
}

interface AppLinkProps extends LinkProps {
  className?: string,
  theme?: AppLinkTheme,
  children?: ReactNode,
}

export const AppLink = memo(({
	children,
	className,
	to,
	theme = AppLinkTheme.PRIMARY,
	...otherProps
}: AppLinkProps) => (
	<Link
		className={cn(className, {}, [theme])}
		to={to}
		{...otherProps}
	>
		{children}
	</Link>
));