package models

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type Session struct {
	BaseModel
	DeviceUUID string
	Name       string
	TotalTime  float64
	Completed  bool
}

func GetSession(ctx *gin.Context, db *gorm.DB, sessionID uuid.UUID) (*Session, error) {
	var session *Session
	result := db.WithContext(ctx).
		Model(&Session{}).
		Where("id = ?", sessionID).
		First(&session)

	if result.Error != nil || result.RowsAffected < 1 {
		return nil, errors.WithStack(result.Error)
	}

	return session, nil
}

func GetSessions(ctx *gin.Context, db *gorm.DB, deviceUUID string) ([]Session, error) {
	var sessions []Session
	result := db.WithContext(ctx).
		Model(&Session{}).
		Where("device_uuid = ? and completed = false", deviceUUID).
		Find(&sessions)

	if result.Error != nil {
		return nil, errors.WithStack(result.Error)
	}

	return sessions, nil
}

func UpdateSession(ctx *gin.Context, db *gorm.DB, session *Session) error {
	result := db.WithContext(ctx).
		Model(&Session{}).
		Where("id = ?", session.ID).
		Updates(map[string]interface{}{
			"name":       session.Name,
			"total_time": session.TotalTime,
			"completed":  session.Completed,
		})

	return errors.WithStack(result.Error)
}

func UpdateSessionTotalTime(ctx *gin.Context, db *gorm.DB, sessionID uuid.UUID, totalTime float64) error {
	result := db.WithContext(ctx).
		Model(&Session{}).
		Where("id = ?", sessionID).
		Updates(map[string]interface{}{
			"total_time": totalTime,
		})

	return errors.WithStack(result.Error)
}
