import { PointToastData } from "@/components/PointToast";
import { UserPoints } from "@/interface/index";
import { useState } from "react";

async function fetchUserPoints(): Promise<UserPoints | null> {
  try {
    const res = await fetch("/api/users/me/points?limit=1&offset=0");
    if (!res.ok) return null;
    return res.json();
  } catch {
    return null;
  }
}

export function usePointToast() {
  const [toastData, setToastData] = useState<PointToastData | null>(null);

  // アクション実行前に呼ぶ
  const captureBeforePoints = async (): Promise<UserPoints | null> => {
    return fetchUserPoints();
  };

  // アクション実行後に呼ぶ
  const showToastAfter = async (before: UserPoints | null) => {
    const after = await fetchUserPoints();
    if (!after || !before) return;
    if (after.totalPoints === before.totalPoints) return;

    setToastData({
      earnedPoints: after.totalPoints - before.totalPoints,
      prevTotalPoints: before.totalPoints,
      nextTotalPoints: after.totalPoints,
      prevLevel: before.level,
      nextLevel: after.level,
      prevCurrentLevelWidth: before.currentLevelWidth,
      nextCurrentLevelWidth: after.currentLevelWidth,
      prevNextLevelPoints: before.nextLevelPoints,
      nextNextLevelPoints: after.nextLevelPoints,
    });
  };

  const closeToast = () => setToastData(null);

  return { toastData, captureBeforePoints, showToastAfter, closeToast };
}
