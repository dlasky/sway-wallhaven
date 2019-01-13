package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/Toyz/GoHaven"
	"github.com/urfave/cli"
	bolt "go.etcd.io/bbolt"
)

func main() {
	app := cli.NewApp()
	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:  "set",
			Usage: "",
		},
		cli.BoolFlag{
			Name:  "fetch",
			Usage: "",
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
			err = downloadWallpapers(width, height)
			if err != nil {
				log.Fatal(err)
			}
		}

		if c.Bool("res") {
			getResolution()
		}

		if c.Bool("set") {
			setWallpaper()
		}

		if c.Bool("get") {
			getWallpaper()
		}

		db, err := bolt.Open("my.db", 0600, nil)
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

	out, err := exec.Command("/usr/bin/swaymsg", "-rt", "get_outputs").Output()
	if err != nil {
		return 0, 0, err
	}
	outputs := make(SwayOutputs, 0, 0)
	json.Unmarshal(out, &outputs)
	rect := outputs[0].Rect
	return rect.Width, rect.Height, nil
}

func downloadWallpapers(width, height int) error {
	res := new(GoHaven.Resolutions)
	res.Set(fmt.Sprintf("%vx%v", width, height))

	gh := GoHaven.New()
	ghi, err := gh.Search("landscape", res)
	if err != nil {
		log.Fatal(err)
	}
	for _, res := range ghi.Results {
		detail, err := res.ImageID.Details()
		if err != nil {
			return err
		}
		p, err := detail.Download(".")
		if err != nil {
			return err
		}
		fmt.Println(p)
	}
	return nil
}

func setWallpaper() error {
	images, err := filepath.Glob("wallhaven-*")
	if err != nil {
		return err
	}
	rand.Seed(time.Now().Unix())
	img := fmt.Sprint(images[rand.Intn(len(images))])
	imgFull := fmt.Sprintf("/home/aerolith/Development/go/wallhaven/%v", img)
	cmd := exec.Command("/usr/bin/swaymsg", "output", "*", "bg", imgFull, "fill")
	out, err := cmd.CombinedOutput()
	fmt.Printf("%s", out)
	return err
}

func getWallpaper() error {
	return nil
}
