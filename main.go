package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"time"

	"github.com/Toyz/GoHaven"
	"github.com/urfave/cli"
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
	}
	app.Name = "wallhaven swaywm"
	app.Usage = "download and set wallpapers"
	app.Commands = []cli.Command{
		{
			Name:    "fetch",
			Aliases: []string{"f"},
			Usage:   "fetch new wallpapers from wallhaven",
			Flags: []cli.Flag{
				cli.StringFlag{Name: "search"},
			},
			Action: func(c *cli.Context) error {
				width, height, err := getResolution()
				if err != nil {
					log.Fatal(err)
				}
				err = downloadWallpapers(c.String("search"), c.String("cache"), width, height)
				if err != nil {
					log.Fatal(err)
				}

				return nil
			},
		},
		{
			Name:    "resolution",
			Aliases: []string{"r"},
			Usage:   "print resolution and exit (useful for debugging)",
			Action: func(c *cli.Context) error {
				w, h, err := getResolution()
				if err != nil {
					log.Fatal(err)
				}
				fmt.Printf("%vx%v", w, h)
				return nil
			},
		},
		{
			Name:    "set",
			Aliases: []string{"s"},
			Usage:   "set sets a new randomized wallpaper and exits",
			Action: func(c *cli.Context) error {
				err := setWallpaper(c)
				if err != nil {
					log.Fatal(err)
				}
				return nil
			},
		},
		{
			Name:    "get",
			Aliases: []string{"g"},
			Usage:   "returns the currently set wallpaper and exits",
			Action: func(c *cli.Context) error {
				err := getWallpaper(c)
				if err != nil {
					log.Fatal(err)
				}
				return nil
			},
		},
		{
			Name:    "restore",
			Aliases: []string{"r"},
			Usage:   "restores the previously set wallpaper and exits",
			Action: func(c *cli.Context) error {
				err := restoreWallpaper(c)
				if err != nil {
					log.Fatal(err)
				}
				return nil
			},
		},
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

func setWallpaper(c *cli.Context) error {

	dirPath := getCachePathFromCtx(c)
	fmt.Println(dirPath)
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
	conn.Close()
	if bytes.Compare(msg.Payload, []byte(`[ { "success": true } ]`)) == 0 {
		db, err := getDbFromCtx(c)
		if err != nil {
			return err
		}
		err = db.setWallpaper(img)
		if err != nil {
			return err
		}
		db.close()
	} else {
		return fmt.Errorf("%s", msg.Payload)
	}

	return err
}

func getWallpaper(c *cli.Context) error {
	db, err := getDbFromCtx(c)
	if err != nil {
		return err
	}
	wallpaper, err := db.getWallpaper()
	if err != nil {
		return err
	}
	fmt.Printf("%v", wallpaper)
	return nil
}

func restoreWallpaper(c *cli.Context) error {
	db, err := getDbFromCtx(c)
	if err != nil {
		return err
	}
	wallpaper, err := db.getWallpaper()
	if err != nil {
		return err
	}
	conn, err := getSocket()
	if err != nil {
		return err
	}
	msg, err := trip(conn, message{Type: messageTypeRunCommand, Payload: []byte("output * bg " + wallpaper + " fill")})
	if err != nil {
		return err
	}
	fmt.Printf("%s", msg.Payload)
	conn.Close()
	// if bytes.Compare(msg.Payload, []byte(`[ { "success": true } ]`)) == 0 {
	return nil
}
