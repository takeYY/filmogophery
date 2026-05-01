import { ok } from "neverthrow";
import { Logger } from "pino";
import {
  findPointHistory as _findPointHistory,
  findUserPoints as _findUserPoints,
} from "../../repositories/points/points.repository";

// レベル設計:
// Lv1→2: 100pt, Lv2→3: 200pt, Lv3→4: 400pt, Lv4→5: 800pt
// Lv5以降: 1,000pt固定
const LEVEL_THRESHOLDS = [0, 100, 300, 700, 1500] as const;
const FIXED_LEVEL_POINTS = 1000;

// calcLevel は累計ポイントからレベルを計算する
export function calcLevel(totalPoints: number): number {
  for (let i = LEVEL_THRESHOLDS.length - 1; i >= 0; i--) {
    if (totalPoints >= LEVEL_THRESHOLDS[i]) {
      const baseLevel = i + 1;
      if (baseLevel < LEVEL_THRESHOLDS.length) {
        return baseLevel;
      }
      // Lv5以降: 固定ポイントで加算
      const extra = Math.floor(
        (totalPoints - LEVEL_THRESHOLDS[LEVEL_THRESHOLDS.length - 1]) /
          FIXED_LEVEL_POINTS,
      );
      return LEVEL_THRESHOLDS.length + extra;
    }
  }
  return 1;
}

// calcNextLevelPoints は次のレベルアップまでの残りポイントを計算する
export function calcNextLevelPoints(
  totalPoints: number,
  level: number,
): number {
  if (level < LEVEL_THRESHOLDS.length) {
    return LEVEL_THRESHOLDS[level] - totalPoints;
  }
  // Lv5以降: 次の1,000pt区切りまでの残り
  const pointsInCurrentLevel =
    (totalPoints - LEVEL_THRESHOLDS[LEVEL_THRESHOLDS.length - 1]) %
    FIXED_LEVEL_POINTS;
  return FIXED_LEVEL_POINTS - pointsInCurrentLevel;
}

// calcCurrentLevelWidth は現在のレベル幅を返す
export function calcCurrentLevelWidth(level: number): number {
  if (level < LEVEL_THRESHOLDS.length) {
    return LEVEL_THRESHOLDS[level] - LEVEL_THRESHOLDS[level - 1];
  }
  return FIXED_LEVEL_POINTS;
}

export type UserPointsResult = {
  totalPoints: number;
  level: number;
  nextLevelPoints: number;
  currentLevelWidth: number;
  pointHistory: {
    id: number;
    points: number;
    action: string;
    referenceId: number;
    createdAt: string | null;
  }[];
};

type Deps = {
  findUserPoints?: typeof _findUserPoints;
  findPointHistory?: typeof _findPointHistory;
};

export async function getUserPoints(
  logger: Logger,
  userId: number,
  limit: number,
  offset: number,
  {
    findUserPoints = _findUserPoints,
    findPointHistory = _findPointHistory,
  }: Deps = {},
) {
  logger.info({ userId }, "getUserPoints called");

  const [up, history] = await Promise.all([
    findUserPoints(userId),
    findPointHistory(userId, limit, offset),
  ]);

  const totalPoints = up?.totalPoints ?? 0;
  const level = calcLevel(totalPoints);
  const nextLevelPoints = calcNextLevelPoints(totalPoints, level);
  const currentLevelWidth = calcCurrentLevelWidth(level);

  return ok({
    totalPoints,
    level,
    nextLevelPoints,
    currentLevelWidth,
    pointHistory: history.map((h) => ({
      id: h.id,
      points: h.points,
      action: h.action,
      referenceId: h.referenceId,
      createdAt: h.createdAt,
    })),
  } satisfies UserPointsResult);
}
