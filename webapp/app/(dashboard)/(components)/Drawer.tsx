import Link from "next/link";
import { PropsWithChildren } from "react";
import { Icon } from "./Icon";

export function Drawer({ children }: PropsWithChildren) {
  return (
    <div className="drawer lg:drawer-open">
      <input id="my-drawer-2" type="checkbox" className="drawer-toggle" />
      <section className="drawer-content px-10 py-10">{children}</section>
      <div className="drawer-side min-h-full border-r">
        <ul className="menu  text-base-content w-60">
          <li className="text-sm">
            <Link href="/monitors">
              <Icon name="ri-window-fill" />
              Monitors
            </Link>
          </li>
          <li>
            <Link href="/incidents">
              <Icon name="ri-error-warning-line" />
              Incidents
            </Link>
          </li>
        </ul>
      </div>
    </div>
  );
}
