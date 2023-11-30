import { ViewHeader } from "../(components)/ViewHeader";

export default async function MonitorsPage() {
  return (
    <>
      <ViewHeader title="Monitors">
        <span>
          Websites and applications added to Cloudflare under this account are
          listed here. Select one to update configurations or monitor security
          and performance.
        </span>
      </ViewHeader>
    </>
  );
}
