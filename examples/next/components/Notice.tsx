export function Notice({ title, children }: { title: string; children: React.ReactNode }) {
  return (
    <aside className="tg-surface(note) tg-col tg-gap(xs) tg-pad(md)" role="status">
      <strong className="tg-text(label)">{title}</strong>
      <p>{children}</p>
    </aside>
  );
}
