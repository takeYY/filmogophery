export type PointHistoryItem = {
  id: number;
  points: number;
  action: "watch_history" | "review";
  referenceId: number;
  createdAt: string | null;
};

export type UserPoints = {
  totalPoints: number;
  level: number;
  nextLevelPoints: number;
  currentLevelWidth: number;
  pointHistory: PointHistoryItem[];
};
