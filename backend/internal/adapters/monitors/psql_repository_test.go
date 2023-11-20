package monitors_test

import (
	"context"
	"database/sql"
	"log"
	"os"
	"testing"
	"time"

	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"

	"github.com/flowck/doberman/internal/adapters/monitors"
	"github.com/flowck/doberman/internal/common/psql"
	"github.com/flowck/doberman/tests"
)

var (
	db  *sql.DB
	ctx context.Context
)

func TestPsqlRepository_Insert(t *testing.T) {
	repo := monitors.NewPsqlRepository(db)
	accID := tests.FixtureAndSaveAccount(ctx, t, db)
	teamID, _ := tests.FixtureAndSaveTeam(ctx, t, db, accID)
	memberID, _ := tests.FixtureAndSaveTeamMember(ctx, t, db, accID)

	for i := 0; i < 5; i++ {
		mo := tests.FixtureMonitor(t, accID, teamID, memberID)
		assert.NoError(t, repo.Insert(ctx, mo))
	}
}

func TestMain(m *testing.M) {
	var cancel context.CancelFunc
	ctx, cancel = context.WithTimeout(context.Background(), time.Minute*1)
	defer cancel()

	var err error
	db, err = psql.Connect(os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal(err)
	}

	err = psql.ApplyMigrations(db, "../../../misc/sql/migrations")
	if err != nil {
		log.Fatal(err)
	}

	os.Exit(m.Run())
}
