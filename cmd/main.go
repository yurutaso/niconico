package main

import (
	"flag"
	"fmt"
	"github.com/yurutaso/niconico"
	"log"
	"os"
	"os/user"
	"strings"
)

func main() {
	fs := flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	var (
		flagHelp      = fs.Bool("h", false, "Help")
		flagTimeshift = fs.Bool("t", false, "Timeshift (default: false)")
		argVideoOut   = fs.String("o", "", `Name of Output (default: "", which means using video title as filename )`)
		argEmail      = fs.String("e", "", `Email address`)
		argPassword   = fs.String("p", "", `Password`)
	)

	fs.Parse(os.Args[1:])
	for 0 < fs.NArg() {
		fs.Parse(fs.Args()[1:])
	}

	if *flagHelp {
		fmt.Println(`Usage: NicoNico videoid -e email -p password [-t] [-o fileout] [-h]`)
		return
	}
	if len(*argPassword) == 0 || len(*argEmail) == 0 {
		log.Fatal(fmt.Errorf(`You must set email address and password`))
	}

	id := os.Args[1]
	if *flagTimeshift {
		nc := niconico.NewNicoClient()
		nc.SetUser(*argEmail, *argPassword)
		nc.Login()

		liveVideo, err := nc.GetLiveInfo(id)
		if err != nil {
			log.Fatal(err)
		}
		fileout := ``
		if len(*argVideoOut) == 0 {
			title, err := nc.GetLiveTitle(id)
			if err != nil {
				log.Fatal(err)
			}
			fileout = title + `.mp4`
		} else {
			usr, _ := user.Current()
			fileout = strings.Replace(*argVideoOut, "~", usr.HomeDir, 1)
		}
		nc.DownloadTimeshift(liveVideo, fileout)
	} else {
		nc, err := niconico.NewNicoDriver(*argEmail, *argPassword)
		if err != nil {
			log.Fatal(err)
		}
		defer nc.Close()

		info, err := nc.GetVideoInfo(id)
		if err != nil {
			log.Fatal(err)
		}

		var fileout string
		if len(*argVideoOut) == 0 {
			fileout = info.Title + `.mp4`
		} else {
			fileout = *argVideoOut
		}
		nc.DownloadVideo(info.URL, fileout)

		/*
			nc := niconico.NewNicoClient()
			nc.SetUser(*argEmail, *argPassword)
			nc.Login()
			//videoURL, err := nc.GetVideoURL(id)
			videoURL, err := nc.GetVideoURLHTML5(id)
			if err != nil {
				log.Fatal(err)
			}
			doc, err := nc.GetVideoCookie(id)
			if err != nil {
				log.Fatal(err)
			}
			fileout := ``
			if len(*argVideoOut) == 0 {
				fileout = niconico.GetTitle(doc) + `.mp4`
			} else {
				fileout = *argVideoOut
			}
			nc.DownloadVideo(videoURL, fileout)
		*/
	}
}
