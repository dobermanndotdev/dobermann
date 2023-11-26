import { PropsWithChildren } from "react";
import Link from "next/link";

export default async function Layout({ children }: PropsWithChildren) {
  return (
    <>
      <header className="border px-10 py-4 flex justify-between">
        <div>Dobermann</div>
        <div>
          <Link href="/api/auth/logout" className="btn btn-xs">
            Logout
          </Link>
        </div>
      </header>
      <section className="px-10 py-4">{children}</section>
    </>
  );
}
