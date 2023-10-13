package databases

import (
	"database/sql"
	"fmt"
	"github.com/Meraj/PoSql"
	"github.com/joho/godotenv"
	"os"
	"strconv"
	"time"
)

var PostgresDB *sql.DB

func init() {
	err := godotenv.Load("local.env")
	if err != nil {
		panic(err.Error())
	}
	ConnectToPostgresDB()
}
func ConnectToPostgresDB() {
	port, _ := strconv.ParseInt(os.Getenv("DB_PORT"), 10, 64)
	connection := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		port,
		os.Getenv("DB_USER"),
		"postgres",
		os.Getenv("DB_NAME"),
	)
	var err error
	PostgresDB, err = sql.Open("postgres", connection)
	if err != nil {
		fmt.Print(err.Error())
		time.Sleep(5000)
		ConnectToPostgresDB()
	}
	migration()
}

func migration() {
	print("asd")
	var migration PoSql.DatabaseCreator
	migration.DB = PostgresDB
	migration = migration.Table("snapp_users").
		ID().
		Column("snapp_id", "VARCHAR not null UNIQUE")

	migration = migration.Table("mentors").
		ID().
		Column("name", "VARCHAR (255) not null").
		Column("photo", "VARCHAR")

	migration = migration.Table("participants").
		ID().
		Column("name", "VARCHAR (255) not null").
		Column("code", "VARCHAR").
		Column("photo", "VARCHAR").
		Column("is_active", "BOOLEAN default true").
		Column("mentor_id", "BIGINT REFERENCES mentors (id)")

	migration = migration.Table("voting").
		ID().
		Column("name", "VARCHAR (255) not null").
		Column("description", "TEXT").
		Column("winner_id", "BIGINT REFERENCES participants (id)").
		Column("started_at", "TIMESTAMP NOT NULL").
		Column("ended_at", "TIMESTAMP NOT NULL")

	migration = migration.Table("voting_votes_cache").
		ID().
		Column("voting_id", "BIGINT REFERENCES voting (id)").
		Column("participant_id", "BIGINT REFERENCES participants (id)").
		Column("votes", "BIGINT").
		Column("is_winner", "Boolean")

	migration = migration.Table("user_voting").
		ID().
		Column("voting_id", "BIGINT REFERENCES voting (id)").
		Column("owner_id", "BIGINT REFERENCES snapp_users (id)").
		Column("vote_id", "BIGINT REFERENCES participants (id)")
	migration = migration.Table("vouchers").
		ID().
		Column("owner_id", "BIGINT REFERENCES snapp_users (id)").
		Column("name", "VARCHAR (255) not null").
		Column("description", "TEXT").
		Column("code", "VARCHAR (255) not null").
		Column("icon", "VARCHAR").
		Column("is_new", "BOOLEAN default true")

	migration = migration.Table("banners").
		ID().
		Column("is_active", "BOOLEAN default true").
		Column("image", "VARCHAR not null").
		Column("link", "VARCHAR")

	migration.Table("files").
		Column("id", "uuid primary key").
		Column("name", "VARCHAR(255)").
		Column("hash", "VARCHAR(255)").
		Column("content_type", "VARCHAR(255)").
		Column("size", "BIGINT").
		Column("created_at", "TIMESTAMP")

	migration.Table("users").
		ID().
		Column("email", "VARCHAR not null UNIQUE").
		Column("password", "VARCHAR not null").
		Column("is_superuser", "BOOLEAN default false").
		Timestamp("created_at")

	migration.Init()
}
