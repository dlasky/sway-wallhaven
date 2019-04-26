package main

import (
	"os"
	"path/filepath"

	"github.com/urfave/cli"

	bolt "go.etcd.io/bbolt"
)

var settings = []byte("settings")
var wallpaper = []byte("wallpaper")
var metaData = []byte("metadata")

type db struct {
	db *bolt.DB
}

func getDbFromCtx(c *cli.Context) (*db, error) {
	return getDb(getConfigPathFromCtx(c))
}

func getDb(cpath string) (*db, error) {
	path := getConfigPath(cpath)
	err := os.MkdirAll(path, os.ModePerm)
	if err != nil {
		return nil, err
	}
	dbc, err := bolt.Open(filepath.Join(path, "sway-wallhaven.db"), 0600, nil)
	if err != nil {
		return nil, err
	}
	db := db{db: dbc}
	return &db, nil
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

func (d *db) setMetaData(path string, meta []byte) error {
	return nil
}

func (d *db) close() error {
	return d.db.Close()
}
