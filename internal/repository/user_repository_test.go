package repository

import (
	"testing"

	"github.com/ekideno/postly/internal/domain"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTestDB(t *testing.T) *UserRepository {
	err := godotenv.Load("../../.env")
	if err != nil {
		t.Fatalf("Error loading .env file: %v", err)
	}

	repo, err := NewGormUserRepository()
	require.NoError(t, err)

	err = repo.db.Exec("TRUNCATE TABLE users CASCADE").Error
	require.NoError(t, err)

	return repo
}

func TestUserRepository_Create(t *testing.T) {
	repo := setupTestDB(t)

	tests := []struct {
		name    string
		user    *domain.User
		wantErr bool
	}{
		{
			name: "успешное создание пользователя",
			user: &domain.User{
				ID:       "test-id-1",
				Username: "testuser1",
				Email:    "test1@example.com",
			},
			wantErr: false,
		},
		{
			name: "ошибка при создании пользователя с существующим ID",
			user: &domain.User{
				ID:       "test-id-1",
				Username: "testuser2",
				Email:    "test2@example.com",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := repo.Create(tt.user)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestUserRepository_GetByID(t *testing.T) {
	repo := setupTestDB(t)

	testUser := &domain.User{
		ID:       "test-id-2",
		Username: "testuser2",
		Email:    "test2@example.com",
	}
	err := repo.Create(testUser)
	require.NoError(t, err)

	tests := []struct {
		name    string
		id      string
		want    *domain.User
		wantErr bool
	}{
		{
			name:    "успешное получение пользователя",
			id:      "test-id-2",
			want:    testUser,
			wantErr: false,
		},
		{
			name:    "пользователь не найден",
			id:      "non-existent-id",
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := repo.GetByID(tt.id)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, got)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want.ID, got.ID)
				assert.Equal(t, tt.want.Username, got.Username)
				assert.Equal(t, tt.want.Email, got.Email)
			}
		})
	}
}

func TestUserRepository_DeleteByID(t *testing.T) {
	repo := setupTestDB(t)

	testUser := &domain.User{
		ID:       "test-id-3",
		Username: "testuser3",
		Email:    "test3@example.com",
	}
	err := repo.Create(testUser)
	require.NoError(t, err)

	tests := []struct {
		name    string
		id      string
		wantErr bool
	}{
		{
			name:    "успешное удаление пользователя",
			id:      "test-id-3",
			wantErr: false,
		},
		{
			name:    "ошибка при удалении несуществующего пользователя",
			id:      "non-existent-id",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := repo.DeleteByID(tt.id)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				_, err := repo.GetByID(tt.id)
				assert.Error(t, err)
			}
		})
	}
}
