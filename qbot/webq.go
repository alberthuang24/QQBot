package qbot

import (
	"github.com/hongjianghuang/phantomjs"
	"github.com/hongjianghuang/pixterm/ansimage"
	"golang.org/x/crypto/ssh/terminal"
	"runtime"
	"time"
	"github.com/lucasb-eyer/go-colorful"
	"fmt"
	"os"
)

const (
	// LoginURL web qq 登录页url
	LoginURL = "http://web2.qq.com"
)

// Webq 用于存储phantomjs打开http://web.qq.com/的信息
type Webq struct {
	WebPage *phantomjs.WebPage
	// CurrentHTML 用来描述当前phantomjs打开的html内容.
	CurrentHTML string
}

// Make ...
func Make() *Webq {
	phantomjs.DefaultProcess.Open()
	page, err := phantomjs.DefaultProcess.CreateWebPage()
	if err != nil {
		panic(err)
	}
	return &Webq{
		WebPage: page,
		CurrentHTML: "",
	}
}

// ToURL 用来跳转到新的页面
func (b *Webq) ToURL(url string) {
	b.WebPage.Open(url)
	//等待页面js加载完毕

	//优化这里的解决方案
	time.Sleep(2 * time.Second) //2 s
}

// Login 调起QQ登录二维码渲染在shell terminal
func (b *Webq) Login() {
	b.ToURL(LoginURL)
	defer b.WebPage.Close()
	defer phantomjs.DefaultProcess.Close()

	info, err := b.WebPage.Evaluate(`function(){
		var a = document.getElementsByTagName('iframe')[0].src;
		return  { url: a};
	}`);

	if err != nil {
		panic(err)
	}

	obj := info.(map[string]interface{})

	qrCodeURL := obj["url"].(string)
	b.ToURL(qrCodeURL)

	info, err = b.WebPage.Evaluate(`function(){
		return document.getElementById('qrlogin_img').src;
	}`);

	if err != nil {
		panic(err)
	}

	qrCodeURL = info.(string)	

	// get terminal size
	tx, ty, err := terminal.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		panic(err)
	}

	tx = 90

	ty = 70

	// get scale mode from flag
	sm := ansimage.ScaleMode(0)

	// get dithering mode from flag
	dm := ansimage.DitheringMode(0)

	mc, err := colorful.Hex("#000000") // RGB color from Hex format

	if err != nil {
		panic(err)
	}

	sfy, sfx := 1, 1 // 8x4 --> with dithering

	pix, err := ansimage.NewScaledFromURL(qrCodeURL, sfy*ty, sfx*tx, mc, sm, dm)

	if err != nil {
		panic(err)
	}
	ansimage.ClearTerminal()
	pix.SetMaxProcs(runtime.NumCPU()) // maximum number of parallel goroutines!
	pix.Draw()
	fmt.Println()
}