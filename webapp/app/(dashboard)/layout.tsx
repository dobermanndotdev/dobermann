import { PropsWithChildren } from "react";

import { Drawer } from "./(components)/Drawer";
import { Header } from "./(components)/Header";

export default async function Layout({ children }: PropsWithChildren) {
  return (
    <>
      <Header />
      <Drawer>{children}</Drawer>
    </>
  );
}
