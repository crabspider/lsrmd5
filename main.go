package main

import (
	"crypto/md5"
	"encoding/hex"
	"flag"
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

type Config struct {
	Directories []string
	Output      string
	Flat        bool
}

func main() {
	startTime := time.Now()

	var output string
	var flat bool

	flag.StringVar(&output, "output", "", "出力ファイルのパス")
	flag.StringVar(&output, "o", "", "出力ファイルのパス（短縮）")
	flag.BoolVar(&flat, "flat", false, "ファイル名のみモード")

	flag.Parse()

	config := Config{
		Directories: flag.Args(),
		Output:      output,
		Flat:        flat,
	}

	if config.Output == "" {
		fmt.Println("outputオプションが指定されていません。")
		os.Exit(1)
	}

	if len(config.Directories) == 0 {
		fmt.Println("対象ディレクトリが指定されていません。")
		os.Exit(1)
	}

	err := lsrmd5(config)
	if err != nil {
		log.Fatalln(err)
	}

	log.Println("complete")
	log.Println(time.Since(startTime))
}

func lsrmd5(config Config) error {
	type Entry struct {
		Name string
		MD5  string
	}

	var entries []Entry

	resultFile, err := os.Create(config.Output)
	if err != nil {
		return err
	}

	for _, dir := range config.Directories {
		err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}

			if d.IsDir() {
				return nil
			}

			md5String, err := calcMD5(path)
			if err != nil {
				return err
			}

			if config.Flat {
				entries = append(entries, Entry{
					Name: d.Name(),
					MD5:  md5String,
				})
				if err != nil {
					return err
				}
			} else {
				_, err = fmt.Fprintf(resultFile, "%s  %s\n", md5String, strings.ReplaceAll(path, "\\", "/"))
				if err != nil {
					return err
				}
			}

			return nil
		})
		if err != nil {
			return err
		}
	}

	if config.Flat {
		// ソート
		sort.SliceStable(entries, func(i, j int) bool { return entries[i].Name < entries[j].Name })

		// 出力
		for _, entry := range entries {
			_, err = fmt.Fprintf(resultFile, "%s  %s\n", entry.MD5, entry.Name)
			if err != nil {
				return err
			}
		}
	}

	err = resultFile.Close()
	if err != nil {
		return err
	}

	return nil
}

func calcMD5(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}

	h := md5.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}

	err = f.Close()
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(h.Sum(nil)), nil
}
