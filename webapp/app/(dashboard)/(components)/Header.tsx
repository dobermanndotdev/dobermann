import Link from "next/link";
import { AccountMenu } from "./AccountMenu";

export function Header() {
  return (
    <header className="px-6 py-4 flex justify-between border-b">
      <div className="font-bold">
        <h1>
          <Link href="/monitors">Dobermann</Link>
        </h1>
      </div>
      <div>
        <AccountMenu />
      </div>
    </header>
  );
}
