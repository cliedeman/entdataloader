package integration

import (
	"context"
	"entgo.io/ent/dialect"
	"github.com/cliedeman/entdataloader/internal/integration/ent"
	"github.com/cliedeman/entdataloader/internal/integration/ent/enttest"
	"github.com/cliedeman/entdataloader/internal/integration/ent/migrate"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/require"
	"reflect"
	"runtime"
	"strings"
	"sync"
	"testing"
	"time"
)

func TestSQLite(t *testing.T) {
	client := enttest.Open(t, dialect.SQLite, "file:ent?mode=memory&cache=shared&_fk=1", opts)
	defer client.Close()
	for _, tt := range tests {
		name := runtime.FuncForPC(reflect.ValueOf(tt).Pointer()).Name()
		t.Run(name[strings.LastIndex(name, ".")+1:], func(t *testing.T) {
			drop(t, client)
			tt(t, client)
		})
	}
}

var (
	opts = enttest.WithMigrateOptions(
		migrate.WithDropIndex(true),
		migrate.WithDropColumn(true),
	)
	tests = [...]func(*testing.T, *ent.Client){
		O2MSingleBatch,
		O2MTwoBatches,
		M2O,
	}
)

func O2MSingleBatch(t *testing.T, client *ent.Client) {
	require := require.New(t)
	ctx := context.Background()
	usr1, usr2 := setup(client, ctx)

	dataloader := ent.NewUserPostsDataloader(ctx, client.Debug(), ent.UserPostsDataloaderWithFilter(func(q *ent.PostQuery) *ent.PostQuery {
		// example filter
		return q.Where()
	}))
	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		posts, err := dataloader.Load(usr1.ID)
		require.NoError(err)
		require.Equal(len(posts), 5)
	}()

	go func() {
		defer wg.Done()
		posts, err := dataloader.Load(usr2.ID)
		require.NoError(err)
		require.Equal(len(posts), 10)
	}()

	wg.Wait()
}

func setup(client *ent.Client, ctx context.Context) (*ent.User, *ent.User) {
	usr1 := client.User.Create().SetUsername("user1").SaveX(ctx)
	usr2 := client.User.Create().SetUsername("user2").SaveX(ctx)
	for i := 0; i < 5; i++ {
		client.Post.Create().SetAuthor(usr1).SaveX(ctx)
	}
	for i := 0; i < 10; i++ {
		client.Post.Create().SetAuthor(usr2).SaveX(ctx)
	}
	return usr1, usr2
}

func O2MTwoBatches(t *testing.T, client *ent.Client) {
	require := require.New(t)
	ctx := context.Background()
	usr1, usr2 := setup(client, ctx)

	dataloader := ent.NewUserPostsDataloader(ctx, client.Debug())
	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		posts, err := dataloader.Load(usr1.ID)
		require.NoError(err)
		require.Equal(len(posts), 5)
	}()

	go func() {
		// Delay second batch for 2 calls
		time.Sleep(time.Second)
		defer wg.Done()
		posts, err := dataloader.Load(usr2.ID)
		require.NoError(err)
		require.Equal(len(posts), 10)
	}()

	wg.Wait()
}

func M2O(t *testing.T, client *ent.Client) {
	require := require.New(t)
	ctx := context.Background()
	setup(client, ctx)

	dataloader := ent.NewPostAuthorDataloader(ctx, client.Debug())
	var wg sync.WaitGroup

	posts := client.Post.Query().AllX(ctx)

	for _, p := range posts {
		go func(p *ent.Post) {
			defer wg.Done()
			wg.Add(1)
			u, err := dataloader.Load(p.AuthorID)
			require.NoError(err)
			require.NotNil(u)
		}(p)
	}

	wg.Wait()
}

func drop(t *testing.T, client *ent.Client) {
	t.Log("drop data from database")
	ctx := context.Background()
	client.Post.Delete().ExecX(ctx)
	client.User.Delete().ExecX(ctx)
}
