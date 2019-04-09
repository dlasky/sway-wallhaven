package main

import (
	"fmt"
	"log"

	"github.com/urfave/cli"

	bolt "go.etcd.io/bbolt"
)

var settings = []byte("settings")
var wallpaper = []byte("wallpaper")

type db struct {
	db *bolt.DB
}

func getDbFromCtx(c *cli.Context) (*db, error) {
	return getDb(getConfigPathFromCtx(c))
}

func getDb(cpath string) (*db, error) {
	//TODO: (dlasky) handle folder creation in .config if nil
	path := getConfigPath(cpath)
	dbc, err := bolt.Open(fmt.Sprintf("%vsway-wallhaven.db", path), 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	db := db{db: dbc}
	return &db, err
}

func (d *db) setWallpaper(path string) error {
	return d.db.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists(settings)
		if err != nil {
			return err
		}
		return bucket.Put(wallpaper, []byte(path))
	})
}

func (d *db) getWallpaper() (string, error) {
	var out []byte
	err := d.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(settings)
		tmp := bucket.Get(wallpaper)
		out = make([]byte, len(tmp))
		copy(out, tmp)
		return nil
	})
	return string(out), err
}

func (d *db) close() error {
	return d.db.Close()
}
