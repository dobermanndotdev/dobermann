import Link from "next/link";
import { AccountMenu } from "./AccountMenu";
import { PropsWithChildren } from "react";

export function Header() {
  return (
    <header className="px-6 flex justify-between border-b">
      <div className="font-bold py-4">
        <h1>
          <Link href="/monitors">Dobermann</Link>
        </h1>
      </div>
      <div>
        <MenuItem>
          <AccountMenu />
        </MenuItem>
      </div>
    </header>
  );
}

function MenuItem({ children }: PropsWithChildren) {
  return (
    <div className="flex justify-center h-full hover:bg-gray-300 transition-colors duration-300">
      {children}
    </div>
  );
}
