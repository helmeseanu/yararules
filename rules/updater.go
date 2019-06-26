/*type program struct {
	wg   sync.WaitGroup
	quit chan struct{}
}

var client *http.Client*/	

/*func main() {
	u, _:= urlhelper.Parse("www.github.com/Yara-Rules/rules/Packers")

	err:=gt.Get("C:/Users/helme/Desktop/Pipelist/rules",u)
	if(err!=nil){
		fmt.Printf("eroare")
	}
	prg := &program{}

	if err := svc.Run(prg); err != nil {
		log.Fatal(err)
	}
}*/

/*func (p *program) Init(env svc.Environment) error {
	return nil
}

func (p *program) Start() error {
	p.quit = make(chan struct{})

   	proxy := goproxy.NewProxyHttpServer()
	proxy.OnResponse().DoFunc(func(resp *http.Response, ctx *goproxy.ProxyCtx) *http.Response {
		body,_ := ioutil.ReadAll(resp.Body)
		resp.Body.Close()

		s1 := rand.NewSource(time.Now().UnixNano())
    	r1 := rand.New(s1)
    	st1 := strconv.Itoa(r1.Intn(1000))

		err := ioutil.WriteFile("./cache/"+st1+".txt", body, 0644)

	if err==nil{
		conn, _ := npipe.Dial(`\\.\pipe\proxyscanner`)
		fmt.Fprintf(conn, st1+".txt")

		tmp := make([]byte, 1)
		_,_ = conn.Read(tmp)
		conn.Close()

		if bytes.Equal(tmp, []byte("0")){
			body2,_ := ioutil.ReadFile("./warning.html")
			s := string(body2)
			t := &http.Response{
  				Status:        "200 OK",
  				StatusCode:    200,
  				Proto:         "HTTP/1.1",
  				ProtoMajor:    1,
  				ProtoMinor:    1,
  				Body:          ioutil.NopCloser(bytes.NewBufferString(s)),
  				ContentLength: int64(len(body)),
  				Request:       ctx.Req,
  				Header:        make(http.Header, 0),
			}
			os.Remove("./cache/"+st1+".txt")
			return t;
		}else if bytes.Equal(tmp, []byte("1")){
			resp.Body=ioutil.NopCloser(bytes.NewBuffer(body))
			os.Remove("./cache/"+st1+".txt")
			return resp
		}
	}
		os.Remove("./cache/"+st1+".txt")
		return nil
	})

	go func() {
    	log.Fatal(http.ListenAndServe(":8080", proxy))
    }()

    p.wg.Add(1)
	go func() {
		<-p.quit
		p.wg.Done()
	}()

	return nil
}

func (p *program) Stop() error {
	close(p.quit)
	p.wg.Wait()
	return nil
}*/

package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"sync"

	getter "github.com/hashicorp/go-getter"
)

func main() {
	var mode getter.ClientMode
	mode = getter.ClientModeDir

	/*pwd, err := os.Getwd()
	if err != nil {
		log.Fatalf("Error getting wd: %s", err)
	}*/

	opts := []getter.ClientOption{}

	ctx, cancel := context.WithCancel(context.Background())

	client := &getter.Client{
		Ctx:     ctx,
		Src:     "github.com/helmeseanu/yararules",
		Dst:     "./rulesfolder",
		Pwd:     "./rulesfolder",
		Mode:    mode,
		Options: opts,
	}

	wg := sync.WaitGroup{}
	wg.Add(1)
	errChan := make(chan error, 2)
	go func() {
		defer wg.Done()
		defer cancel()
		if err := client.Get(); err != nil {
			errChan <- err
		}
	}()

	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt)

	select {
	case sig := <-c:
		signal.Reset(os.Interrupt)
		cancel()
		wg.Wait()
		log.Printf("signal %v", sig)
	case <-ctx.Done():
		wg.Wait()
		log.Printf("success!")
	case err := <-errChan:
		wg.Wait()
		log.Fatalf("Error downloading: %s", err)
	}
}