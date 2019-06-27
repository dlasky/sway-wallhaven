package main

import (
	"bytes"
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/dlasky/go-wallhaven"
	"github.com/urfave/cli"
	"go.i3wm.org/i3"
)

func main() {

	//i3 overrides to work with sway
	i3.SocketPathHook = func() (string, error) {
		out, err := exec.Command("sway", "--get-socketpath").CombinedOutput()
		if err != nil {
			return "", fmt.Errorf("getting sway socketpath: %v (output: %s)", err, out)
		}
		return string(out), nil
	}

	i3.IsRunningHook = func() bool {
		out, err := exec.Command("pgrep", "-c", "sway\\$").CombinedOutput()
		if err != nil {
			log.Printf("sway running: %v (output: %s)", err, out)
		}
		return bytes.Compare(out, []byte("1")) == 0
	}

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
				cli.StringFlag{
					Name:  "search",
					Usage: "search landscape",
					Value: "landscape",
				},
			},
			Action: func(c *cli.Context) error {
				width, height, err := getResolution()
				if err != nil {
					log.Fatal(err)
				}

				err = downloadWallpapers(c, width, height)
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
		{
			Name:    "Delete",
			Aliases: []string{"rm"},
			Usage:   "removes the currently set wallpaper and sets a new one",
			Action: func(c *cli.Context) error {
				err := removeWallpaper(c)
				if err != nil {
					log.Fatal(err)
				}
				err = setWallpaper(c)
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

func getResolution() (int64, int64, error) {

	outputs, err := i3.GetOutputs()
	if err != nil {
		return 0, 0, err
	}
	rect := outputs[0].Rect
	return rect.Width, rect.Height, nil

}

func downloadWallpapers(c *cli.Context, width, height int64) error {
	term := c.String("search")

	results, err := wallhaven.SearchWallpapers(&wallhaven.Search{
		Query: wallhaven.Q{
			Tags: []string{term},
		},
		AtLeast: wallhaven.Resolution{
			Width:  width,
			Height: height,
		},
	})
	if err != nil {
		return err
	}

	for _, wp := range results.Data {
		err := wp.Download(getCachePathFromCtx(c))
		if err != nil {
			return err
		}
		fmt.Println(wp.URL)
	}

	return nil
}

func setWallpaper(c *cli.Context) error {

	dirPath := getCachePathFromCtx(c)
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
	fmt.Println(img)

	results, err := i3.RunCommand("output * bg " + img + " fill")
	if err != nil {
		return err
	}

	var ok = true
	for _, result := range results {
		ok = ok && result.Success
	}

	if ok {
		db, err := getDbFromCtx(c)
		if err != nil {
			return err
		}
		err = db.setWallpaper(img)
		if err != nil {
			return err
		}
		err = db.close()
		if err != nil {
			return err
		}
	}

	return nil
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
	return db.close()
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
	_, err = i3.RunCommand("output * bg " + wallpaper + " fill")
	if err != nil {
		return err
	}
	return nil
}

func removeWallpaper(c *cli.Context) error {
	db, err := getDbFromCtx(c)
	if err != nil {
		return err
	}
	wallpaper, err := db.getWallpaper()
	if err != nil {
		return err
	}
	err = os.Remove(wallpaper)
	if err != nil {
		return err
	}
	return db.close()
}
