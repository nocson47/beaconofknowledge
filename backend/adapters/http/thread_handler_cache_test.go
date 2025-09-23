package http

import (
	"context"
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/alicebob/miniredis/v2"
	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
	"github.com/nocson47/beaconofknowledge/internal/entities"
	"github.com/stretchr/testify/require"
)

// fakeThreadService implements usecases.ThreadService for testing
type fakeThreadService struct {
	called int
	thread *entities.Thread
}

func (f *fakeThreadService) CreateThread(ctx context.Context, t *entities.Thread) (int, error) {
	return 0, nil
}
func (f *fakeThreadService) GetThreadByID(ctx context.Context, id int) (*entities.Thread, error) {
	f.called++
	return f.thread, nil
}
func (f *fakeThreadService) GetAllThreads(ctx context.Context) ([]*entities.Thread, error) {
	return nil, nil
}
func (f *fakeThreadService) UpdateThread(ctx context.Context, t *entities.Thread) error { return nil }
func (f *fakeThreadService) DeleteThread(ctx context.Context, id int) error             { return nil }

func TestGetThreadByID_CacheAside(t *testing.T) {
	// start miniredis
	mr, err := miniredis.Run()
	require.NoError(t, err)
	defer mr.Close()

	rdb := redis.NewClient(&redis.Options{Addr: mr.Addr()})

	// prepare fake service that returns a thread
	th := &entities.Thread{ID: 1, UserID: 1, Title: "Hi", Body: "Body"}
	fake := &fakeThreadService{thread: th}

	handler := NewThreadHandler(fake, rdb)

	app := fiber.New()
	app.Get("/threads/:id", handler.GetThreadByID)

	// first request: cache miss -> service called once
	req := httptest.NewRequest("GET", "/threads/1", nil)
	resp, err := app.Test(req)
	require.NoError(t, err)
	require.Equal(t, 200, resp.StatusCode)
	require.Equal(t, 1, fake.called)

	// second request: should hit cache; service not called again
	req2 := httptest.NewRequest("GET", "/threads/1", nil)
	resp2, err := app.Test(req2)
	require.NoError(t, err)
	require.Equal(t, 200, resp2.StatusCode)
	require.Equal(t, 1, fake.called)

	// verify cached value exists in redis
	ctx := context.Background()
	key := "thread:1"
	val, err := rdb.Get(ctx, key).Result()
	require.NoError(t, err)
	var cached entities.Thread
	err = json.Unmarshal([]byte(val), &cached)
	require.NoError(t, err)
	require.Equal(t, th.Title, cached.Title)
}
