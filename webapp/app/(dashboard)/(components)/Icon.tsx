interface Props {
  name: string;
}

export function Icon({ name }: Props) {
  return <i className={`${name}`}></i>;
}
