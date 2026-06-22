package services

import (
	"context"

	"github.com/rs/zerolog"
	"gorm.io/gorm"

	"filmogophery/internal/app/repositories"
	"filmogophery/internal/app/responses"
	"filmogophery/internal/pkg/constant"
	"filmogophery/internal/pkg/gen/model"
)

// ポイント付与量の定義
const (
	PointsForReview = 20 // レビュー投稿

	// 視聴記録: 上映時間による段階的ポイント
	PointsForWatchShort  = 10 // ~90分
	PointsForWatchMedium = 15 // 91~150分
	PointsForWatchLong   = 20 // 151分~
)

// レベル設計:
// Lv1→2: 100pt, Lv2→3: 200pt, Lv3→4: 400pt, Lv4→5: 800pt
// Lv5以降: 1,000pt固定
var levelThresholds = []int32{0, 100, 300, 700, 1500}

const fixedLevelPoints int32 = 1000

// CalcLevel は累計ポイントからレベルを計算する
func CalcLevel(totalPoints int32) int32 {
	for i := len(levelThresholds) - 1; i >= 0; i-- {
		if totalPoints >= levelThresholds[i] {
			baseLevel := int32(i + 1)
			if baseLevel < int32(len(levelThresholds)) {
				return baseLevel
			}
			// Lv5以降: 固定ポイントで加算
			extra := (totalPoints - levelThresholds[len(levelThresholds)-1]) / fixedLevelPoints
			return int32(len(levelThresholds)) + int32(extra)
		}
	}
	return 1
}

// CalcNextLevelPoints は次のレベルアップまでの残りポイントを計算する
func CalcNextLevelPoints(totalPoints, level int32) int32 {
	if int(level) < len(levelThresholds) {
		return levelThresholds[level] - totalPoints
	}
	// Lv5以降: 次の1,000pt区切りまでの残り
	pointsInCurrentLevel := (totalPoints - levelThresholds[len(levelThresholds)-1]) % fixedLevelPoints
	return fixedLevelPoints - pointsInCurrentLevel
}

// CalcCurrentLevelWidth は現在のレベル幅（レベルアップに必要なポイント数）を返す
func CalcCurrentLevelWidth(level int32) int32 {
	// levelThresholds[i] → [i+1] の差分がそのレベルの幅
	// level=1: thresholds[1]-thresholds[0] = 100
	// level=2: thresholds[2]-thresholds[1] = 200
	// ...
	// level=4: thresholds[4]-thresholds[3] = 800
	// level>=5: 固定
	if int(level) < len(levelThresholds) {
		return levelThresholds[level] - levelThresholds[level-1]
	}
	return fixedLevelPoints
}

// CalcWatchPoints は上映時間からポイントを計算する
func CalcWatchPoints(runtimeMinutes int32) int32 {
	switch {
	case runtimeMinutes <= 90:
		return PointsForWatchShort
	case runtimeMinutes <= 150:
		return PointsForWatchMedium
	default:
		return PointsForWatchLong
	}
}

type (
	IPointService interface {
		// 視聴記録ポイントを付与する
		GrantWatchHistoryPoints(ctx context.Context, tx *gorm.DB, operator *model.Users, watchHistory *model.WatchHistory, movie *model.Movies) error
		// レビュー投稿ポイントを付与する
		GrantReviewPoints(ctx context.Context, tx *gorm.DB, operator *model.Users, review *model.Reviews) error
		// ユーザーのポイント・レベルを取得する
		GetUserPoints(ctx context.Context, userID int32) (*model.UserPoints, error)
	}

	pointService struct {
		pointRepo repositories.IPointRepository
	}
)

func NewPointService(pointRepo repositories.IPointRepository) IPointService {
	return &pointService{pointRepo}
}

// 視聴記録ポイントを付与する
func (s *pointService) GrantWatchHistoryPoints(
	ctx context.Context, tx *gorm.DB, operator *model.Users, watchHistory *model.WatchHistory, movie *model.Movies,
) error {
	log := zerolog.Ctx(ctx)

	points := CalcWatchPoints(movie.RuntimeMinutes)

	// 現在のポイントを取得してレベルを計算
	current, err := s.pointRepo.FindOrCreateByUserID(ctx, tx, operator.ID)
	if err != nil {
		log.Error().Msgf("failed to get user points(userID=%d): %s", operator.ID, err.Error())
		return responses.InternalServerError()
	}
	newLevel := CalcLevel(current.TotalPoints + points)

	if _, err := s.pointRepo.AddPoints(ctx, tx, operator.ID, points, newLevel); err != nil {
		log.Error().Msgf("failed to add watch history points(userID=%d): %s", operator.ID, err.Error())
		return responses.InternalServerError()
	}

	if err := s.pointRepo.SaveHistory(ctx, tx, &model.PointHistory{
		UserID:      operator.ID,
		Points:      points,
		Action:      constant.PointActionWatchHistory,
		ReferenceID: watchHistory.ID,
	}); err != nil {
		log.Error().Msgf("failed to save point history(userID=%d): %s", operator.ID, err.Error())
		return responses.InternalServerError()
	}

	log.Info().Msgf("granted %d points for watch history(userID=%d)", points, operator.ID)
	return nil
}

// レビュー投稿ポイントを付与する
func (s *pointService) GrantReviewPoints(
	ctx context.Context, tx *gorm.DB, operator *model.Users, review *model.Reviews,
) error {
	log := zerolog.Ctx(ctx)

	// 現在のポイントを取得してレベルを計算
	current, err := s.pointRepo.FindOrCreateByUserID(ctx, tx, operator.ID)
	if err != nil {
		log.Error().Msgf("failed to get user points(userID=%d): %s", operator.ID, err.Error())
		return responses.InternalServerError()
	}
	newLevel := CalcLevel(current.TotalPoints + PointsForReview)

	if _, err := s.pointRepo.AddPoints(ctx, tx, operator.ID, PointsForReview, newLevel); err != nil {
		log.Error().Msgf("failed to add review points(userID=%d): %s", operator.ID, err.Error())
		return responses.InternalServerError()
	}

	if err := s.pointRepo.SaveHistory(ctx, tx, &model.PointHistory{
		UserID:      operator.ID,
		Points:      PointsForReview,
		Action:      constant.PointActionReview,
		ReferenceID: review.ID,
	}); err != nil {
		log.Error().Msgf("failed to save point history(userID=%d): %s", operator.ID, err.Error())
		return responses.InternalServerError()
	}

	log.Info().Msgf("granted %d points for review(userID=%d)", PointsForReview, operator.ID)
	return nil
}

// ユーザーのポイント・レベルを取得する
func (s *pointService) GetUserPoints(ctx context.Context, userID int32) (*model.UserPoints, error) {
	log := zerolog.Ctx(ctx)

	up, err := s.pointRepo.FindOrCreateByUserID(ctx, nil, userID)
	if err != nil {
		log.Error().Msgf("failed to get user points(userID=%d): %s", userID, err.Error())
		return nil, responses.InternalServerError()
	}
	return up, nil
}
