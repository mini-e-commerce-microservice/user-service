package users_test

import (
	"context"
	"github.com/DATA-DOG/go-sqlmock"
	wsqlx "github.com/SyaibanAhmadRamadhan/sqlx-wrapper"
	"github.com/go-faker/faker/v4"
	"github.com/guregu/null/v5"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/require"
	"math/rand"
	"regexp"
	"testing"
	"user-service/internal/model"
	"user-service/internal/repositories"
	"user-service/internal/repositories/users"
)

func Test_repository_FindOneUser(t *testing.T) {
	dbMock, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer dbMock.Close()

	ctx := context.TODO()
	sqlxDB := sqlx.NewDb(dbMock, "sqlmock")

	sqlxx := wsqlx.NewRdbms(sqlxDB)

	r := users.NewRepository(sqlxx)

	t.Run("should be return correct", func(t *testing.T) {
		expectedInput := users.FindOneUserInput{
			ID:    null.IntFrom(rand.Int63()),
			Email: null.StringFrom(faker.Email()),
		}
		expectedOutput := users.FindOneUserOutput{
			Data: model.User{
				ID:    expectedInput.ID.Int64,
				Email: expectedInput.Email.String,
			},
		}

		mock.ExpectQuery(regexp.QuoteMeta(
			`SELECT id, email FROM users WHERE id = $1 AND email = $2`,
		)).WithArgs(expectedInput.ID.Int64, expectedInput.Email.String).WillReturnRows(sqlmock.NewRows([]string{
			"id", "email",
		}).AddRow(
			expectedOutput.Data.ID,
			expectedOutput.Data.Email,
		))

		output, err := r.FindOneUser(ctx, expectedInput)
		require.NoError(t, err)
		require.Equal(t, expectedOutput, output)
	})

	t.Run("should be return error no record", func(t *testing.T) {
		expectedInput := users.FindOneUserInput{
			ID:    null.IntFrom(rand.Int63()),
			Email: null.StringFrom(faker.Email()),
		}
		mock.ExpectQuery(regexp.QuoteMeta(
			`SELECT id, email FROM users WHERE id = $1 AND email = $2`,
		)).WithArgs(expectedInput.ID.Int64, expectedInput.Email.String).
			WillReturnError(wsqlx.ErrRecordNoRows)

		_, err = r.FindOneUser(ctx, expectedInput)
		require.ErrorIs(t, err, repositories.ErrRecordNotFound)
	})
}
