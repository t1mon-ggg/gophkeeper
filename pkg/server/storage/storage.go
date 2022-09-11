package storage

import (
	"context"
	"encoding/hex"
	"errors"
	"net"
	"time"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/mgutz/ansi"
	"golang.org/x/crypto/bcrypt"

	"github.com/t1mon-ggg/gophkeeper/pkg/logging"
	"github.com/t1mon-ggg/gophkeeper/pkg/models"
)

const (
	dbSchema = `CREATE TABLE IF NOT EXISTS public.users (
		id bigserial NOT NULL UNIQUE,
		username varchar NOT NULL UNIQUE,
		secret text NOT NULL
	);

	CREATE TABLE IF NOT EXISTS public.keeper (
		id bigserial NOT NULL UNIQUE,
		username varchar NOT NULL UNIQUE,
		checksum text NOT NULL UNIQUE,
		time  timestamp NOT NULL,
		content text,
		CONSTRAINT keeper_fk FOREIGN KEY (username) REFERENCES users(username) ON UPDATE CASCADE ON DELETE CASCADE
	);

	CREATE TABLE IF NOT EXISTS public.openpgp (
		id bigserial NOT NULL UNIQUE,
		username varchar NOT NULL UNIQUE,
		publickey text NOT NULL,
		time timestamp NOT NULL,
		confirmed bool NULL DEFAULT false,
		CONSTRAINT publickeys_fk FOREIGN KEY (username) REFERENCES users(username) ON UPDATE CASCADE ON DELETE CASCADE
	);

	CREATE TABLE IF NOT EXISTS public.logs (
		id bigserial NOT NULL UNIQUE,
		username varchar NOT NULL UNIQUE,
		ip varchar NOT NULL,
		action text NOT NULL,
		time timestamp NOT NULL,
		checksum text NOT NULL,
		CONSTRAINT logs_fk FOREIGN KEY (username) REFERENCES users(username) ON UPDATE NO ACTION ON DELETE NO ACTION,
		CONSTRAINT logs_checksum_fk FOREIGN KEY (checksum) REFERENCES keeper(checksum) ON UPDATE NO ACTION ON DELETE NO ACTION
	);

	CREATE UNIQUE INDEX IF NOT EXISTS users_names_idx ON public.users USING btree (username);
	CREATE UNIQUE INDEX IF NOT EXISTS keeper_content_idx ON public.keeper USING btree (username, checksum, time);
	CREATE UNIQUE INDEX IF NOT EXISTS openpgp_idx ON public.openpgp USING btree (username, publickey, time);`

	createUser = `INSERT INTO public.users (username, secret) VALUES( $1, $2 );`
	deleteUser = `DELETE FROM public.users WHERE username=$1;`
	checkUser  = `SELECT username, secret FROM public.users WHERE username=$1;`
	checkPGP   = `SELECT publickey, confirmed FROM public.openpgp WHERE username=$1;`
	logAction  = `INSERT INTO public.logs (username, ip, action, time, checksum, sign) VALUES( $1, $2, $3, $4, $5, $6 );`
	logView    = `SELECT time, ip, action, checksum, sign FROM public.logs WHERE username=$1;`
	pushData   = `INSERT INTO public.keeper (username, checksum, time, content) VALUES ( $1, $2, $3, $4);`
	pullData   = `SELECT content FROM public.keeper WHERE sername=$1 AND checksum=$2`
	addPGP     = `INSERT INTO public.openpgp (username, publickey, time) VALUES ( $1, $2, $3 );`
	delPGP     = `DELETE FROM public.openpgp WHERE publickey=$1;`
	confirmPGP = `UPDATE public.openpgp	SET confirmed=true WHERE username=$1 AND publickey=$2;`
)

type psqlStorage struct {
	logger logging.Logger
	db     *pgxpool.Pool
}

func New(dsn string, logger logging.Logger) (Storage, error) {
	s := new(psqlStorage)
	s.logger = logger.WithPrefix("storage")
	s.log().Debug(nil, "Connecting to Postgres database")

	db, err := pgxpool.Connect(context.Background(), dsn)
	if err != nil {
		s.log().Error(err, "error in database connection")
		return nil, err
	}
	s.log().Info(nil, "Database storage connected")
	s.db = db
	err = s.create()
	if err != nil {
		s.log().Fatal(err, "db schema could not applied")
		return nil, err
	}
	return s, nil
}

func (s *psqlStorage) SignUp(username, password string) error {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return err
	}
	s.log().Tracef("Password hash is %v", nil, hex.EncodeToString(hashed))
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	_, err = s.db.Exec(ctx, createUser, username, hex.EncodeToString(hashed))
	if err != nil {
		return err
	}
	s.log().Infof("user %s created", nil, ansi.Color(username, "green+b"))
	return nil
}

func (s *psqlStorage) SignIn(user models.User) error {
	usr := new(string)
	pwd := new(string)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	err := s.db.QueryRow(ctx, checkUser, user.Username).Scan(usr, pwd)
	if err != nil {
		if err == pgx.ErrNoRows {
			s.log().Trace(nil, "provided username not found")
			return errors.New("no such user or password")
		}
		s.log().Trace(nil, "error in user request query")
		return err
	}
	hashed, err := hex.DecodeString(*pwd)
	if err != nil {
		s.log().Tracef("hex decode password hash completed with error %v", err)
		return errors.New("no such user or password")
	}
	err = bcrypt.CompareHashAndPassword(hashed, []byte(user.Password))
	if err != nil {
		s.log().Tracef("error in password verification [%v]", err)
		s.log().Debugf("wrong creditentials provided [user \"%s\" with password \"%s\"]", nil, ansi.Color(user.Username, "red+b"), ansi.Color(user.Password, "red+b"))
		return errors.New("no such user or password")
	}
	return nil
}

func (s *psqlStorage) DeleteUser(username string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	tag, err := s.db.Exec(ctx, deleteUser, username)
	if err != nil {
		return err
	}

	if tag.RowsAffected() == 0 {
		s.log().Trace(nil, "no affected users in delete query")
		return errors.New("no such user")
	}
	s.log().Infof("user %s deleted", nil, ansi.Color(username, "green+b"))
	return nil
}

func (s *psqlStorage) Push(username, checksum string, data []byte) error {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	now := time.Now()
	_, err := s.db.Exec(ctx, pushData, username, checksum, now, data)
	if err != nil {
		s.log().Debug(err, "content save action failed")
		return errors.New("content push failed")
	}

	// s.SaveLog(username, "push")
	return nil
}

func (s *psqlStorage) Pull(username, checksum string) ([]byte, error) {
	data := []byte{}
	return data, nil
}

func (s *psqlStorage) SaveLog(username, action, checksum, sign string, ip *net.IP, date time.Time) error {
	return nil
}

func (s *psqlStorage) GetLog(username string) ([]models.Action, error) {
	actions := []models.Action{}
	return actions, nil
}

func (s *psqlStorage) AddPGP(username, publickey string) error {
	return nil
}

func (s *psqlStorage) DeletePGP(publickey string) error {
	return nil
}

func (s *psqlStorage) Close() {
	s.db.Close()
	s.log().Info(nil, "database connection closed")
}

func (s *psqlStorage) Ping() error {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	err := s.db.Ping(ctx)
	return err
}

func (s *psqlStorage) create() error {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	_, err := s.db.Exec(ctx, dbSchema)
	if err != nil {
		return err
	}
	s.log().Debug(nil, "database schema applied successfully")
	return nil
}

func (s *psqlStorage) log() logging.Logger {
	return s.logger
}
