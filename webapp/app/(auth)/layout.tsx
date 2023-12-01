import { PropsWithChildren } from "react";

export default function Layout({ children }: PropsWithChildren) {
  return (
    <div className="w-full flex justify-center items-center h-screen">
      <div className="w-11/12 md:w-2/4 lg:w-2/6 border px-8 py-10">
        {children}
      </div>
    </div>
  );
}
