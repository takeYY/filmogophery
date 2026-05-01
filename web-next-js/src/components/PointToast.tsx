"use client";

import { useEffect, useState } from "react";

export type PointToastData = {
  earnedPoints: number;
  prevTotalPoints: number;
  nextTotalPoints: number;
  prevLevel: number;
  nextLevel: number;
  prevCurrentLevelWidth: number;
  nextCurrentLevelWidth: number;
  prevNextLevelPoints: number;
  nextNextLevelPoints: number;
};

type Props = {
  data: PointToastData | null;
  onClose: () => void;
};

const DISPLAY_MS = 5000;

export function PointToast({ data, onClose }: Props) {
  // アニメーション用: マウント後に true にしてプログレスバーを伸ばす
  const [animated, setAnimated] = useState(false);
  const isLevelUp = data !== null && data.nextLevel > data.prevLevel;

  useEffect(() => {
    if (!data) {
      setAnimated(false);
      return;
    }
    // 次フレームでアニメーション開始
    const raf = requestAnimationFrame(() => setAnimated(true));
    const timer = setTimeout(onClose, DISPLAY_MS);
    return () => {
      cancelAnimationFrame(raf);
      clearTimeout(timer);
    };
  }, [data]);

  if (!data) return null;

  const prevPercent = Math.min(
    100,
    Math.max(
      0,
      Math.round(
        ((data.prevCurrentLevelWidth - data.prevNextLevelPoints) /
          data.prevCurrentLevelWidth) *
          100,
      ),
    ),
  );

  const nextPercent = Math.min(
    100,
    Math.max(
      0,
      Math.round(
        ((data.nextCurrentLevelWidth - data.nextNextLevelPoints) /
          data.nextCurrentLevelWidth) *
          100,
      ),
    ),
  );

  return (
    <div
      style={{
        position: "fixed",
        bottom: "1.5rem",
        right: "1.5rem",
        zIndex: 9999,
        minWidth: "280px",
        maxWidth: "340px",
      }}
    >
      <div
        className={`card bg-dark border-info text-light shadow`}
        style={{ borderRadius: "0.75rem", overflow: "hidden" }}
      >
        {/* レベルアップバナー */}
        {isLevelUp && (
          <div
            className="text-center py-2 fw-bold"
            style={{
              background: "linear-gradient(90deg, #0dcaf0, #6610f2)",
              letterSpacing: "0.1em",
              fontSize: "1rem",
            }}
          >
            🎉 LEVEL UP! Lv.{data.prevLevel} → Lv.{data.nextLevel}
          </div>
        )}

        <div className="card-body py-3 px-3">
          {/* 獲得ポイント */}
          <div className="d-flex justify-content-between align-items-center mb-2">
            <span className="small text-muted">ポイント獲得</span>
            <span className="fw-bold text-info fs-5">
              +{data.earnedPoints} pt
            </span>
          </div>

          {/* レベル表示 */}
          <div className="d-flex justify-content-between align-items-center mb-1">
            <span className="small">Lv.{data.nextLevel}</span>
            <span className="small text-muted">
              {data.nextTotalPoints.toLocaleString()} pt
            </span>
          </div>

          {/* プログレスバー */}
          <div
            className="progress bg-secondary"
            style={{ height: "10px", borderRadius: "5px" }}
          >
            <div
              className="progress-bar bg-info"
              role="progressbar"
              style={{
                width: `${animated ? nextPercent : prevPercent}%`,
                transition: "width 0.8s ease-in-out",
              }}
              aria-valuenow={animated ? nextPercent : prevPercent}
              aria-valuemin={0}
              aria-valuemax={100}
            />
          </div>

          {/* 次レベルまで */}
          <div className="text-end mt-1">
            <small className="text-muted">
              次のレベルまで あと{" "}
              <span className="text-info">
                {data.nextNextLevelPoints.toLocaleString()} pt
              </span>
            </small>
          </div>
        </div>

        {/* 閉じるボタン */}
        <button
          onClick={onClose}
          className="btn-close btn-close-white position-absolute"
          style={{ top: "0.5rem", right: "0.5rem" }}
          aria-label="閉じる"
        />
      </div>
    </div>
  );
}
