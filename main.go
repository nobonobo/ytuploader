package main

import (
	"bytes"
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/fsnotify/fsnotify"
)

var (
	ptn = regexp.MustCompile(`.*[0-9]{2}\.mp4$`)
)

type output struct {
	timer *time.Timer
}

func (o *output) Write(b []byte) (int, error) {
	msg := strings.TrimSpace(string(b))
	if len(msg) > 0 {
		log.Print(msg)
	}
	o.timer.Reset(30 * time.Second)
	return len(b), nil
}

func proc(ctx context.Context, file string) error {
	log.Println("upload:", file)
	args := []string{"./youtubeuploader"}
	args = append(args, "-filename", file)
	args = append(args, "-notify", "false")
	args = append(args, "-language", "ja")
	args = append(args, "-chunksize", "268435456") // 256MiB
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	log.Print(strings.Join(args, " "))
	cmd := exec.CommandContext(ctx, args[0], args[1:]...)
	cmd.Stdout = &output{timer: time.AfterFunc(30*time.Second, cancel)}
	buff := bytes.NewBuffer(nil)
	cmd.Stderr = buff
	if err := cmd.Run(); err != nil {
		fmt.Println(strings.TrimSpace(buff.String()))
		return err
	}
	base := filepath.Base(file)
	parent := filepath.Dir(file)
	doneDir := filepath.Join(parent, "done")
	err := os.MkdirAll(doneDir, 0755)
	if err != nil {
		return err
	}
	if err := os.Rename(file, filepath.Join(doneDir, base)); err != nil {
		return err
	}
	log.Println("completed:", file)
	return nil
}

func check(ctx context.Context, src string) error {
	files, err := filepath.Glob(src + "/*.mp4")
	if err != nil {
		return err
	}
	srcs := []string{}
	for _, file := range files {
		if !ptn.MatchString(file) {
			continue
		}
		srcs = append(srcs, file)
	}
	sort.Slice(srcs, func(i, j int) bool {
		return srcs[i] < srcs[j]
	})
	for _, file := range srcs {
		if err := proc(ctx, file); err != nil {
			return err
		}
	}
	return nil
}

func main() {
	secrets := "client_secrets.json"
	flag.StringVar(&secrets, "secrets", secrets, "path to json file for client secrets")
	src := "."
	flag.StringVar(&src, "src", src, "src directory")
	flag.Parse()
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
	defer cancel()
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	absSrc, err := filepath.Abs(src)
	if err != nil {
		log.Fatal(err)
	}
	if err := watcher.Add(absSrc); err != nil {
		log.Fatal(err)
	}
	for {
		select {
		case <-ctx.Done():
			return
		case event := <-watcher.Events:
			log.Println("event:", event)
			if event.Op.String() == "CREATE" {
				log.Println("checking...")
				if err := check(ctx, src); err != nil {
					log.Fatalln(err)
				}
			}
		}
	}
}
