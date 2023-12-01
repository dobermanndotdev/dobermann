import { PropsWithChildren } from "react";

interface Props extends PropsWithChildren {
  className?: string;
  variant?: "text-lg" | "text-xl" | "text-2xl" | "text-3xl";
  as?: "h1" | "h2" | "h3" | "h4" | "h5" | "h6";
}

export function Heading({
  as: Comp = "h1",
  children,
  variant = "text-xl",
  className = "",
}: Props) {
  return (
    <Comp className={`font-bold ${variant} ${className}`}>{children}</Comp>
  );
}
