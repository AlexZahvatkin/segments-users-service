package storage

import (
	"context"

	"github.com/AlexZahvatkin/segments-users-service/internal/models"
)

type Storage interface {
	AddUser(ctx context.Context, name string) (models.User, error)
	DeleteUser(ctx context.Context, id int64) error
	GetAllUsersId(ctx context.Context) ([]int64, error)
	GetUserById(ctx context.Context, id int64) (models.User, error)
	AddSegment(ctx context.Context, arg models.AddSegmentParams) (models.Segment, error)
	DeleteSegment(ctx context.Context, name string) error
	GetSegmentByName(ctx context.Context, name string) (models.Segment, error)
	AddUserIntoSegment(ctx context.Context, arg models.AddUserIntoSegmentParams) (models.UsersInSegment, error)
	AddUserIntoSegmentWithExpireDatetime(ctx context.Context, arg models.AddUserIntoSegmentWithExpireDatetimeParams) (models.UsersInSegment, error)
	AddUserIntoSegmentWithTTLInHours(ctx context.Context, arg models.AddUserIntoSegmentWithTTLInHoursParams) (models.UsersInSegment, error)
	GetSegmentsByUserId(ctx context.Context, userID int64) ([]string, error)
	RemoveUserFromSegment(ctx context.Context, arg models.RemoveUserFromSegmentParams) error
	GetSegmentsHistoryByUserId(ctx context.Context, arg models.GetSegmentsHistoryByUserIdParams) ([]models.UsersInSegmentsHistory, error)
}
