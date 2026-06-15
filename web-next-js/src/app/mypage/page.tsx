// mypage/page.tsx
/**
 * マイページ（ポイント・レベル）
 * パス: /mypage
 */
"use client";

import { useAuth } from "@/hooks/useAuth";
import { UserPoints } from "@/interface/index";
import { useEffect, useState } from "react";

const ACTION_LABELS: Record<string, string> = {
  watch_history: "視聴記録",
  review: "レビュー投稿",
};

export default function MyPage() {
  const [loading, setLoading] = useState(true);
  const [userPoints, setUserPoints] = useState<UserPoints | null>(null);
  const [guideOpen, setGuideOpen] = useState(false);

  useEffect(() => {
    const saved = localStorage.getItem("pointGuideOpen");
    // 初回（未設定）はデフォルトで開く、一度閉じたら以降は閉じた状態を維持
    setGuideOpen(saved === null ? true : saved === "true");
  }, []);

  const toggleGuide = (next: boolean) => {
    setGuideOpen(next);
    localStorage.setItem("pointGuideOpen", String(next));
  };

  const { checked } = useAuth();

  useEffect(() => {
    if (!checked) return;
    const fetchPoints = async () => {
      setLoading(true);
      try {
        const res = await fetch("/api/users/me/points");
        if (res.ok) {
          const data: UserPoints = await res.json();
          setUserPoints(data);
        }
      } catch {
        console.error("ポイント情報の取得に失敗しました");
      } finally {
        setLoading(false);
      }
    };
    fetchPoints();
  }, [checked]);

  if (loading) {
    return (
      <div className="container py-5 text-center">
        <div className="spinner-border text-info" role="status">
          <span className="visually-hidden">Loading...</span>
        </div>
      </div>
    );
  }

  if (!userPoints) {
    return (
      <div className="container py-5 text-center text-light">
        <p>ポイント情報を取得できませんでした。</p>
      </div>
    );
  }

  const earnedInLevel =
    userPoints.currentLevelWidth - userPoints.nextLevelPoints;
  const progressPercent = Math.min(
    100,
    Math.max(
      0,
      Math.round((earnedInLevel / userPoints.currentLevelWidth) * 100),
    ),
  );

  return (
    <div className="container py-4">
      <h3 className="text-center mb-4">マイページ</h3>

      {/* ポイント・レベルカード */}
      <div className="row justify-content-center mb-4">
        <div className="col-md-6">
          <div className="card bg-dark border-info text-light">
            <div className="card-body text-center">
              <div className="mb-3">
                <span
                  className="badge bg-info text-dark fs-5 px-4 py-2"
                  style={{ borderRadius: "2rem" }}
                >
                  Lv. {userPoints.level}
                </span>
              </div>
              <h5 className="card-title mb-1">
                {userPoints.totalPoints.toLocaleString()} pt
              </h5>
              <p className="text-muted small mb-3">累計ポイント</p>

              {/* レベルアップまでのプログレスバー */}
              <div className="mb-1">
                <div
                  className="progress bg-secondary"
                  style={{ height: "10px" }}
                >
                  <div
                    className="progress-bar bg-info"
                    role="progressbar"
                    style={{ width: `${progressPercent}%` }}
                    aria-valuenow={progressPercent}
                    aria-valuemin={0}
                    aria-valuemax={100}
                  />
                </div>
              </div>
              <p className="text-light small">
                次のレベルまで あと{" "}
                <span className="text-info fw-bold">
                  {userPoints.nextLevelPoints.toLocaleString()} pt
                </span>
              </p>
            </div>
          </div>
        </div>
      </div>

      {/* ポイント獲得方法（アコーディオン） */}
      <div className="row justify-content-center mb-4">
        <div className="col-md-6">
          <div className="border border-secondary rounded">
            <button
              className="w-100 d-flex justify-content-between align-items-center px-3 py-2 bg-dark text-light border-0"
              style={{ fontSize: "0.875rem" }}
              onClick={() => toggleGuide(!guideOpen)}
              aria-expanded={guideOpen}
            >
              <span>ポイント獲得方法</span>
              <i className={`bi bi-chevron-${guideOpen ? "up" : "down"}`} />
            </button>
            {guideOpen && (
              <ul className="list-group list-group-flush">
                <li className="list-group-item bg-dark text-light d-flex justify-content-between border-secondary">
                  <span>視聴記録（〜90分）</span>
                  <span className="text-info">+10 pt</span>
                </li>
                <li className="list-group-item bg-dark text-light d-flex justify-content-between border-secondary">
                  <span>視聴記録（91〜150分）</span>
                  <span className="text-info">+15 pt</span>
                </li>
                <li className="list-group-item bg-dark text-light d-flex justify-content-between border-secondary">
                  <span>視聴記録（151分〜）</span>
                  <span className="text-info">+20 pt</span>
                </li>
                <li className="list-group-item bg-dark text-light d-flex justify-content-between border-secondary">
                  <span>レビュー投稿</span>
                  <span className="text-info">+20 pt</span>
                </li>
              </ul>
            )}
          </div>
        </div>
      </div>

      {/* ポイント履歴 */}
      <div className="row justify-content-center">
        <div className="col-md-6">
          <h5 className="text-light mb-3">ポイント履歴</h5>
          {userPoints.pointHistory.length === 0 ? (
            <p className="text-muted">まだポイント履歴がありません。</p>
          ) : (
            <ul className="list-group">
              {userPoints.pointHistory.map((item) => (
                <li
                  key={item.id}
                  className="list-group-item bg-dark text-light border-secondary d-flex justify-content-between align-items-center"
                >
                  <div>
                    <span className="me-2">
                      {ACTION_LABELS[item.action] ?? item.action}
                    </span>
                    {item.createdAt && (
                      <small className="text-muted">
                        {new Date(item.createdAt).toLocaleDateString("ja-JP")}
                      </small>
                    )}
                  </div>
                  <span className="badge bg-info text-dark">
                    +{item.points} pt
                  </span>
                </li>
              ))}
            </ul>
          )}
        </div>
      </div>
    </div>
  );
}
