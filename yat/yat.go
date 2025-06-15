package yat

import (
	"io/fs"
	"strings"
	"slices"
	"path/filepath"
	"beckx.online/yat/ataglib"
)

type YatData struct {
	Files []string
	AudioMetadatas []*ataglib.AudioMetadata
}

func NewYatData(args []string) (*YatData, error) {
	yd := new(YatData)
	var err error
	yd.Files, err = getAudiofiles(args, []string{".mp3"})
	if err != nil {
		return nil, err
	} else {
		return yd, nil
	}
}

func (yd *YatData) ReadAudioMetadata(tagHeaderOnly bool) error {
	for _, file := range yd.Files{
		amd, err := ataglib.NewAudioMetadata(file, tagHeaderOnly)
		if err != nil {
			return err
		}
		yd.AudioMetadatas = append(yd.AudioMetadatas, amd)
	}
	return nil
}

func getAudiofiles(args []string, pattern []string) ([]string, error) {
	files := []string{}
	for _, arg := range args{
		err := filepath.WalkDir(arg, func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if d.IsDir(){
				return nil
			}
			ext := strings.ToLower(filepath.Ext(path))
			if slices.Contains(pattern, ext) {
				abspath, err := filepath.Abs(path)
				if err != nil {
					return err
				}
				if !slices.Contains(files, abspath) {
					files = append(files, abspath)
				}
			}
			return nil
		})
		if err != nil {
			return nil, err
		}
	}
	return files, nil
}
