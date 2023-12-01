import { Button } from "@@/components/Button";
import { Heading } from "@@/components/Heading";
import { InputField } from "@@/components/Input";
import Link from "next/link";

export default function CreateAccountPage() {
  return (
    <>
      <Heading>Get started with Dobbermann</Heading>
      <form className="mt-4 flex flex-col gap-2">
        <InputField label="Account Name" />
        <InputField label="E-mail" type="email" />
        <InputField label="Password" type="password" />
        <Button className="btn-primary my-4">Create account</Button>
        <span className="block text-sm">
          Do you already have an account? <Link href="/login">Log in</Link>
        </span>
      </form>
    </>
  );
}
