import { Icon } from "./Icon";

export function AccountMenu() {
  return (
    <details className="dropdown dropdown-bottom dropdown-end">
      <summary className="m-1 btn btn-xs">
        <Icon name="ri-user-fill" />
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
