package database_test

import (
	"context"
	"database/sql"
	"os"
	"testing"
	"time"

	"github.com/AlexZahvatkin/segments-users-service/internal/database"
	"github.com/AlexZahvatkin/segments-users-service/internal/models"
	"github.com/AlexZahvatkin/segments-users-service/internal/utils/testutils"
	"github.com/stretchr/testify/assert"
)

var (
	databaseURL string
)

func TestMain(m *testing.M) {
	testutils.LoadEnv()
	databaseURL = os.Getenv("TEST_DATABASE_URL")

	os.Exit(m.Run())
}


func TestAddUser(t *testing.T) {
	query := database.TestDB(t, databaseURL)
	user := models.NewTestUser()
	res, err := query.AddUser(context.Background(), user.Name)
	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, user.Name, res.Name)
}

func TestGetAllUsers(t *testing.T) {
	query := database.TestDB(t, databaseURL)
	user := models.NewTestUser()
	query.AddUser(context.Background(), user.Name)
	query.AddUser(context.Background(), user.Name + "2")
	res, err := query.GetAllUsersId(context.Background())
	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, 2, len(res))
}

func TestDeleteUser(t *testing.T) {
	query := database.TestDB(t, databaseURL)
	user := models.NewTestUser()
	addedUser, err := query.AddUser(context.Background(), user.Name)
	assert.NotNil(t, addedUser)
	assert.NoError(t, err)
	ids, err := query.GetAllUsersId(context.Background())
	assert.NoError(t, err)
	assert.NotNil(t, ids[0])
	err = query.DeleteUser(context.Background(), ids[0])
	assert.NoError(t, err)
}

func TestAddSegment(t *testing.T) {
	query := database.TestDB(t, databaseURL)
	segment := models.NewTestSegment()
	res, err := query.AddSegment(context.Background(), database.AddSegmentParams{
		Name: segment.Name,
		Description: segment.Description,
	})
	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, segment.Description, res.Description)
	assert.Equal(t, segment.Name, res.Name)
}

func TestDeleteSegment(t *testing.T) {
	query := database.TestDB(t, databaseURL)
	segment := models.NewTestSegment()
	res, err := query.AddSegment(context.Background(), database.AddSegmentParams{
        Name: segment.Name,
        Description: segment.Description,
	})
	assert.NoError(t, err)
	assert.NotNil(t, res)
	err = query.DeleteSegment(context.Background(), segment.Name)
	assert.NoError(t, err)
	_, err = query.GetSegmentByName(context.Background(), segment.Name)
	assert.Error(t, sql.ErrNoRows)
}

func TestAddUserIntoSegment(t *testing.T) {
	query := database.TestDB(t, databaseURL)
	segment := models.NewTestSegment()
	user := models.NewTestUser()
	_ ,err := query.AddSegment(context.Background(), database.AddSegmentParams{
		Name: segment.Name,
    })
	assert.NoError(t, err)
	addedUser, err := query.AddUser(context.Background(), user.Name)
	assert.NoError(t, err)
	res, err := query.AddUserIntoSegment(context.Background(), database.AddUserIntoSegmentParams{
		UserID: addedUser.ID,
		SegmentName: segment.Name,
	})
	assert.NoError(t, err)
	assert.Equal(t, addedUser.ID, res.UserID)
	assert.Equal(t, segment.Name, res.SegmentName)
}

func TestAddUserIntoSegmentWithTTLInHours(t *testing.T) {
	query := database.TestDB(t, databaseURL)
	segment := models.NewTestSegment()
	user := models.NewTestUser()
	_ ,err := query.AddSegment(context.Background(), database.AddSegmentParams{
		Name: segment.Name,
    })
	assert.NoError(t, err)
	addedUser, err := query.AddUser(context.Background(), user.Name)
	assert.NoError(t, err)
	timeAfterHour := time.Now().Add(time.Hour)
	res, err := query.AddUserIntoSegmentWithTTLInHours(context.Background(), database.AddUserIntoSegmentWithTTLInHoursParams{
		UserID: addedUser.ID,
		SegmentName: segment.Name,
		NumberOfHours: 1,
	})
	assert.NoError(t, err)
	assert.Equal(t, addedUser.ID, res.UserID)
	assert.Equal(t, segment.Name, res.SegmentName)
	assert.NotNil(t, res.ExpireAt)
	assert.Equal(t, timeAfterHour.Year(), res.ExpireAt.Time.Year())
	assert.Equal(t, timeAfterHour.Month(), res.ExpireAt.Time.Month())
	assert.Equal(t, timeAfterHour.Day(), res.ExpireAt.Time.Day())
	assert.Equal(t, timeAfterHour.Hour(), res.ExpireAt.Time.Hour())
	assert.Equal(t, timeAfterHour.Minute(), res.ExpireAt.Time.Minute())
}

func TestGetSegmentsByUserId(t *testing.T) {
	query := database.TestDB(t, databaseURL)
	segment1 := models.NewTestSegment()
	segment2 := models.NewTestSegment()
	segment2.Name = "TestSegment2"
	user := models.NewTestUser()
	_ ,err := query.AddSegment(context.Background(), database.AddSegmentParams{
		Name: segment1.Name,
    })
	assert.NoError(t, err)
	_ , err = query.AddSegment(context.Background(), database.AddSegmentParams{
		Name: segment2.Name,
    })
	assert.NoError(t, err)
	addedUser, err := query.AddUser(context.Background(), user.Name)
	assert.NoError(t, err)
	segmentForUser1, err := query.AddUserIntoSegmentWithTTLInHours(context.Background(), database.AddUserIntoSegmentWithTTLInHoursParams{
		UserID: addedUser.ID,
		SegmentName: segment1.Name,
		NumberOfHours: 1,
	})
	assert.NoError(t, err)
	segmentForUser2, err := query.AddUserIntoSegmentWithTTLInHours(context.Background(), database.AddUserIntoSegmentWithTTLInHoursParams{
		UserID: addedUser.ID,
		SegmentName: segment2.Name,
		NumberOfHours: 1,
	})
	assert.NoError(t, err)
	res, err := query.GetSegmentsByUserId(context.Background(), addedUser.ID)
	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, 2, len(res))
	assert.Equal(t, segmentForUser1.SegmentName, res[0])
	assert.Equal(t, segmentForUser2.SegmentName, res[1])
}

func TestRemoveUserFromSegment(t *testing.T) {
	query := database.TestDB(t, databaseURL)
	segment := models.NewTestSegment()
	user := models.NewTestUser()
	_ ,err := query.AddSegment(context.Background(), database.AddSegmentParams{
		Name: segment.Name,
    })
	assert.NoError(t, err)
	addedUser, err := query.AddUser(context.Background(), user.Name)
	assert.NoError(t, err)
	addRec, err := query.AddUserIntoSegment(context.Background(), database.AddUserIntoSegmentParams{
		UserID: addedUser.ID,
		SegmentName: segment.Name,
	})
	assert.NoError(t, err)
	assert.Equal(t, addedUser.ID, addRec.UserID)
	assert.Equal(t, segment.Name, addRec.SegmentName)
	err = query.RemoveUserFromSegment(context.Background(), database.RemoveUserFromSegmentParams{
		UserID: addedUser.ID,
		SegmentName: segment.Name,
	})
	assert.NoError(t, err)
	_, err = query.GetSegmentsByUserId(context.Background(), addedUser.ID)
	assert.Error(t, sql.ErrNoRows)
}

func TestGetSegmentsHistoryByUserId(t *testing.T) {
	query := database.TestDB(t, databaseURL)
	segment := models.NewTestSegment()
	user := models.NewTestUser()
	_ ,err := query.AddSegment(context.Background(), database.AddSegmentParams{
		Name: segment.Name,
    })
	assert.NoError(t, err)
	addedUser, err := query.AddUser(context.Background(), user.Name)
	assert.NoError(t, err)
	addRec, err := query.AddUserIntoSegment(context.Background(), database.AddUserIntoSegmentParams{
		UserID: addedUser.ID,
		SegmentName: segment.Name,
	})
	assert.NoError(t, err)
	assert.Equal(t, addedUser.ID, addRec.UserID)
	assert.Equal(t, segment.Name, addRec.SegmentName)
	err = query.RemoveUserFromSegment(context.Background(), database.RemoveUserFromSegmentParams{
		UserID: addedUser.ID,
		SegmentName: segment.Name,
	})
	assert.NoError(t, err)
	res, err := query.GetSegmentsHistoryByUserId(context.Background(), database.GetSegmentsHistoryByUserIdParams{
		UserID: addRec.UserID,
		FromDate: time.Now().Add(-time.Hour),
		ToDate: time.Now().Add(time.Hour),
	}) 
	assert.NoError(t, err)
	assert.Equal(t, 2, len(res))
	assert.Equal(t, segment.Name, res[0].SegmentName)
	assert.Equal(t, segment.Name, res[1].SegmentName)
	assert.Equal(t, user.ID, res[0].UserID)
	assert.Equal(t, user.ID, res[1].UserID)
	assert.Equal(t, "inserted", res[0].ActionType)
	assert.Equal(t, "deleted", res[1].ActionType)
}