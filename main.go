package main

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"flag"
	"fmt"
	"io"
	"io/fs"
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

	var config Config

	flag.StringVar(&config.Output, "output", "", "出力ファイルのパス")
	flag.StringVar(&config.Output, "o", "", "出力ファイルのパス（短縮）")
	flag.BoolVar(&config.Flat, "flat", false, "ファイル名のみモード")

	flag.Parse()

	config.Directories = flag.Args()

	err := lsrmd5(config)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println("complete")
	fmt.Println(time.Since(startTime))
}

func lsrmd5(config Config) error {
	type Entry struct {
		Name string
		MD5  string
	}

	if config.Output == "" {
		return errors.New("outputオプションが指定されていません。")
	}

	if len(config.Directories) == 0 {
		return errors.New("対象ディレクトリが指定されていません。")
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
			} else {
				_, err = fmt.Fprintf(resultFile, "%s  %s\n", md5String, strings.ReplaceAll(path, string(os.PathSeparator), "/"))
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
