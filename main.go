package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"time"

	"github.com/Toyz/GoHaven"
	"github.com/urfave/cli"
	bolt "go.etcd.io/bbolt"
)

func main() {
	app := cli.NewApp()
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "cache",
			Usage: "--cache ~/.wallhaven sets a working directory to store files",
		},
		cli.StringFlag{
			Name:  "config",
			Usage: "--config ~/.config/sway-wallhaven/config",
		},
		cli.BoolFlag{
			Name:  "set",
			Usage: "--set sets a random wallpaper from the currently available",
		},
		cli.BoolFlag{
			Name:  "fetch",
			Usage: "--fetch attempts to fetch new wallpapers for ",
		},
		cli.StringFlag{
			Name:  "search",
			Usage: "--search the term to search for on wallhaven",
			Value: "landscape",
		},
		cli.BoolFlag{
			Name:  "res",
			Usage: "--show current resolution",
		},
	}
	app.Name = "wallhaven swaywm"
	app.Usage = "download and set wallpapers"
	app.Action = func(c *cli.Context) error {

		if c.Bool("fetch") {
			width, height, err := getResolution()
			if err != nil {
				log.Fatal(err)
			}
			err = downloadWallpapers(c.String("search"), c.String("cache"), width, height)
			if err != nil {
				log.Fatal(err)
			}
		}

		if c.Bool("res") {
			w, h, err := getResolution()
			if err != nil {
				log.Fatal(err)
			}
			fmt.Printf("%vx%v", w, h)
		}

		if c.Bool("set") {
			err := setWallpaper(getCachePath(c.String("cache")))
			if err != nil {
				log.Fatal(err)
			}
		}

		if c.Bool("get") {
			getWallpaper()
		}

		db, err := bolt.Open("wallhaven.db", 0600, nil)
		if err != nil {
			log.Fatal(err)
		}
		defer db.Close()

		return nil
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}

}

func getResolution() (int, int, error) {

	conn, err := getSocket()
	defer conn.Close()

	msg, err := trip(conn, message{Type: messageTypeGetOutputs})
	outputs := make(SwayOutputs, 0, 0)
	err = json.Unmarshal(msg.Payload, &outputs)
	if err != nil {
		return 0, 0, err
	}
	rect := outputs[0].Rect
	return rect.Width, rect.Height, nil

}

func downloadWallpapers(term string, path string, width, height int) error {
	res := new(GoHaven.Resolutions)
	res.Set(fmt.Sprintf("%vx%v", width, height))

	gh := GoHaven.New()
	ghi, err := gh.Search(term, res)
	if err != nil {
		return err
	}
	for _, res := range ghi.Results {
		detail, err := res.ImageID.Details()
		if err != nil {
			return err
		}
		fmt.Println(getCachePath(path))
		p, err := detail.Download(getCachePath(path))
		if err != nil {
			return err
		}
		fmt.Println(p)
	}
	return nil
}

func setWallpaper(dirPath string) error {
	images, err := filepath.Glob(fmt.Sprintf("%v/wallhaven-*", dirPath))
	if err != nil {
		return err
	}
	count := len(images)
	if count == 0 {
		return fmt.Errorf("No images found at path %v", dirPath)
	}
	rand.Seed(time.Now().Unix())
	img := fmt.Sprint(images[rand.Intn(len(images))])

	conn, err := getSocket()
	if err != nil {
		return err
	}
	msg, err := trip(conn, message{Type: messageTypeRunCommand, Payload: []byte("output * bg " + img + " fill")})
	if err != nil {
		return err
	}
	fmt.Printf("%s", msg.Payload)
	conn.Close()

	return err
}

func getWallpaper() error {
	return nil
}
