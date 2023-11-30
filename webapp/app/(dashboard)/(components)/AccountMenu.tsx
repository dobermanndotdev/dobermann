import { Icon } from "./Icon";

export function AccountMenu() {
  return (
    <details className="dropdown dropdown-bottom dropdown-end h-full w-full">
      <summary
        role="button"
        className="list-none h-full flex gap-1 items-center justify-center px-2 cursor-pointer"
      >
        <Icon name="ri-user-fill" />
        <Icon name="ri-arrow-down-s-fill" />
      </summary>
      <div className="p-2 shadow menu dropdown-content z-[1] bg-base-100 rounded-box w-52">
        <div className="border-b flex flex-col gap-1 pb-3 px-4">
          <span className="py-0">My profile</span>
          <span className="text-xs py-0">user@email.com</span>
        </div>

        <ul className="mt-1">
          <li>
            <a>Account</a>
          </li>
          <li>
            <a>Log out</a>
          </li>
        </ul>
      </div>
    </details>
  );
}
