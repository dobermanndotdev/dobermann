import { PropsWithChildren } from "react";

interface Props extends PropsWithChildren {
  title: string;
}

export function ViewHeader({ children, title }: Props) {
  return (
    <section>
      <h2 className="font-bold text-2xl mb-1">{title}</h2>
      <div className="text-sm">{children}</div>
    </section>
  );
}
