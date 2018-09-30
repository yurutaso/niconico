package niconico

import (
	"fmt"
	"github.com/sclevine/agouti"
	"log"
	"os/exec"
)

const (
	LOGIN_URL = `https://account.nicovideo.jp/login`
)

type NicoDriver struct {
	driver   *agouti.WebDriver
	page     *agouti.Page
	password string
	email    string
}

type VideoInfo struct {
	URL   string
	Title string
}

func NewNicoDriver(email, password string) (*NicoDriver, error) {
	nc := &NicoDriver{email: email, password: password}
	// Launch chromedriver
	/*
		options := agouti.ChromeOptions("prefs", map[string]interface{}{
			"download.default_directory": savedir,
			"download.prompt_for_download": false,
			"download.directory_upgrade": true
		})
		nc.driver = agouti.ChromeDriver(options)
	*/
	nc.driver = agouti.ChromeDriver()
	if err := nc.driver.Start(); err != nil {
		return nil, err
	}

	// Open a page
	var err error
	nc.page, err = nc.driver.NewPage()
	if err != nil {
		return nil, err
	}

	// Login
	if err := nc.Login(); err != nil {
		return nil, err
	}

	return nc, nil
}

func (nc *NicoDriver) Login() error {
	if err := nc.page.Navigate("https://account.nicovideo.jp/login"); err != nil {
		return err
	}
	if err := nc.page.FindByID("input__mailtel").Fill(nc.email); err != nil {
		return err
	}
	if err := nc.page.FindByID("input__password").Fill(nc.password); err != nil {
		return err
	}
	if err := nc.page.FindByID("login__submit").Click(); err != nil {
		return err
	}
	return nil
}

func (nc *NicoDriver) GetVideoInfo(videoID string) (*VideoInfo, error) {
	if err := nc.page.Navigate("http://www.nicovideo.jp/watch/" + videoID); err != nil {
		return nil, err
	}
	var err error
	info := &VideoInfo{}

	// Get url of the video
	info.URL, err = nc.page.FindByClass("MainVideoPlayer").Find("video").Attribute("src")
	if err != nil {
		return nil, err
	}

	// Get title of the video
	info.Title, err = nc.page.FindByClass(`VideoTitle`).Text()
	if err != nil {
		return nil, err
	}

	return info, nil
}

func (nc *NicoDriver) DownloadVideo(videoURL, fileout string) {
	fmt.Printf("Saving video as %s from %s.", fileout, videoURL)
	// Using http.Get() and io.Copy() oftern fails with error "unexpected EOF"
	cmd := exec.Command("ffmpeg", "-y", "-i", videoURL, "-vcodec", "copy", "-acodec", "copy", fileout)
	result, err := cmd.CombinedOutput()
	if err != nil {
		log.Println(err)
	}
	fmt.Println(string(result))
	fmt.Println(`Complete!`)
}

func (nc *NicoDriver) Close() {
	if nc.page != nil {
		if err := nc.page.CloseWindow(); err != nil {
			log.Printf("%s", err)
		}
	}
	if nc.driver != nil {
		nc.driver.Stop()
	}
}
