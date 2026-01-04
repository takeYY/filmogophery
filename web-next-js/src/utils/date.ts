export function formatRelativeTime(target: Date): string {
  const now = new Date();
  const diffMs = now.getTime() - target.getTime();
  const diffDays = Math.floor(diffMs / (1000 * 60 * 60 * 24));

  if (diffDays < 30) {
    return `${diffDays}日前`;
  } else if (diffDays < 365) {
    const diffMonths = Math.floor(diffDays / 30);
    return `${diffMonths}ヶ月前`;
  } else {
    const diffYears = Math.floor(diffDays / 365);
    return `${diffYears}年前`;
  }
}
