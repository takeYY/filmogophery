// components/MovieCard.tsx
/**
 * 映画カード共通コンポーネント
 */

import { posterUrlPrefix } from "@/constants/poster";
import { Genre, Movie } from "@/interface/index";
import Image from "next/image";

type Props = {
  movie: Movie;
  onClick?: () => void;
  /** ポスター画像のwidth（デフォルト: 200） */
  imageWidth?: number;
  /** ポスター画像のheight（デフォルト: 200） */
  imageHeight?: number;
  /** カード本体の追加クラス（例: "position-relative"） */
  className?: string;
  /** カード左上に表示するオーバーレイ要素（ウォッチリストボタン等） */
  overlay?: React.ReactNode;
  /** カード下部に表示するページ固有のアクション（Reviewボタン等） */
  actions?: React.ReactNode;
};

export function MovieCard({
  movie,
  onClick,
  imageWidth = 200,
  imageHeight = 200,
  className = "",
  overlay,
  actions,
}: Props) {
  const posterSrc =
    posterUrlPrefix + (movie.posterURL ?? "/Agz71U0wcesx87micVn731Z1vPu.jpg");

  const cardContent = (
    <div className="row g-0">
      <div className="col-md-4">
        <Image
          src={posterSrc}
          className="card-img-top w-100 h-auto"
          alt={`${movie.title}のポスター`}
          width={imageWidth}
          height={imageHeight}
          style={{ objectFit: "cover" }}
        />
      </div>
      <div className="col-md-8">
        <div className="card-body text-light">
          {/* タイトル */}
          <h5 className="card-title">{movie.title}</h5>
          {/* ジャンル */}
          {movie.genres.length !== 0 && (
            <div className="card-text d-grid gap-2 d-md-block">
              {movie.genres.map((g: Genre) => (
                <button
                  key={g.code}
                  type="button"
                  className="btn btn-outline-info btn-sm"
                >
                  {g.name}
                </button>
              ))}
            </div>
          )}
          {/* 公開日 */}
          <p className="card-text">
            公開日：{movie.releaseDate.substring(0, 10)}
          </p>
          {/* 概要 */}
          <p className="card-text">
            {movie.overview.length > 40
              ? movie.overview.substring(0, 37) + "..."
              : movie.overview}
          </p>
          {/* ページ固有のアクション */}
          {actions && (
            <div className="border-top border-success mt-2 pt-2">
              <div className="row">
                <div className="col text-center">{actions}</div>
              </div>
            </div>
          )}
        </div>
      </div>
    </div>
  );

  if (onClick) {
    return (
      <button
        className={`card mb-2 bg-dark border-info ${className}`}
        onClick={onClick}
        style={{ cursor: "pointer", textAlign: "left", width: "100%" }}
      >
        {cardContent}
      </button>
    );
  }

  return (
    <div className={`card mb-2 bg-dark border-info ${className}`}>
      {overlay}
      {cardContent}
    </div>
  );
}
