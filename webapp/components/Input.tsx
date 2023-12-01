import { HtmlHTMLAttributes } from "react";

interface Props extends HtmlHTMLAttributes<HTMLInputElement> {
  label?: string;
  error?: string;
  type?: "text" | "email" | "password";
}

export function InputField({ label, error, type = "text", ...props }: Props) {
  return (
    <label className="form-control w-full">
      {!!label && (
        <div className="label">
          <span className="label-text">{label}</span>
        </div>
      )}
      <input {...props} className="input input-bordered w-full" />
      {error && (
        <div className="label">
          <span className="label-text-alt">Bottom Left label</span>
        </div>
      )}
    </label>
  );
}
