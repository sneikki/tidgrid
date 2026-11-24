// Joins class names. tidgrid check reads every string literal passed to cn().
export function cn(...names: (string | false | null | undefined)[]): string {
  return names.filter(Boolean).join(" ");
}
