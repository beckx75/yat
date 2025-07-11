package ataglib

import (
	"encoding/binary"
	"fmt"
	"os"

	"github.com/rs/zerolog/log"
)

type FileIdentifier string

const (
	FI_ID3  FileIdentifier = "ID3"
	FI_FLAC                = "FLAC"
)

type TagEnc byte

const (
	TE_ISO8152  TagEnc = 0x00 // ID3v2.4 terminated with 0x00
	TE_UTF16BOM        = 0x01 //with BOM -> terminated with 0x00 00
	TE_UTF16           = 0x02 // wihout BOM -> terminated with 0x00 00
	TE_UTF8            = 0x03 // ID3v2.4 terminated with 0x00
)

const (
	ID3V2_HEADERSIZE       uint32 = 10
	ID3V2_FRAMEHEADER_SIZE uint32 = 10
)

type AudioMetadata struct {
	Filepath string

	TagIdentifier FileIdentifier
	// ID3 only
	TagVersion                         string
	TagHeaderFlagUnsyncronisation      bool // Bit 7
	TagHeaderFlagExtendedHeader        bool // Bit 6
	TagHeaderFlagExperimentalIndicator bool // Bit 5
	// ID3v2.4 only
	TagHeaderFlagFooterPresent bool // Bit 4

	TagSize    uint32
	TextFrames map[string]*TextFrame
}

type TextFrame struct {
	OrigFramId string
	YatFrameId string
	FrameEnc   TagEnc
	Value      string
	Desc       string // for ID3 TXXX-Frames
}

func NewAudioMetadata(fp string, tagHeaderOnly bool) (*AudioMetadata, error) {
	log.Info().Msg("new AudioMetadata")
	amd := AudioMetadata{
		Filepath:   fp,
		TextFrames: make(map[string]*TextFrame),
	}
	file, err := os.OpenFile(amd.Filepath, os.O_RDONLY, 0644)
	defer file.Close()
	if err != nil {
		return nil, err
	}
	var b byte
	var bytesread uint32 = 0

	// first 3 bytes of audiofile (ID3, fLa)
	var fileIdentifier string
	for i := 0; i < 3; i++ {
		err = binary.Read(file, binary.BigEndian, &b)
		if err != nil {
			return nil, err
		}
		fileIdentifier = fmt.Sprintf("%s%s", fileIdentifier, string(b))
	}
	bytesread = +3
	switch fileIdentifier {
	case "ID3":
		amd.TagIdentifier = FI_ID3
		err = amd.parseID3Header(file, &bytesread)
		if err != nil {
			return nil, err
		}
		if bytesread != 10 {
			return nil, fmt.Errorf("read bytes for header are not equal to 10: %d",
				bytesread)
		}
		err = amd.parseID3Frames(file, &bytesread)
		if err != nil {
			return nil, err
		}
	case "fLa":
		fmt.Println("open flac stream")
	default:
		fmt.Println("currently not supported other audio-frame-tag:", fileIdentifier)
	}

	return &amd, nil
}

func (amd *AudioMetadata) String() string {
	s := fmt.Sprintf("AudioMetadata-Content for '%s'\n", amd.Filepath)
	s = fmt.Sprintf("%s\tTag-Identifier: %s\n", s, amd.TagIdentifier)
	for _, tf := range amd.TextFrames {
		s = fmt.Sprintf("%s\tFrame-Id: '%s'\t'%s'\n", s, tf.OrigFramId, tf.Value)
	}
	return s
}

func id3v23bytesizeToUint32(bsize [4]byte) uint32 {
	msb := uint32(bsize[0]) << 21
	msb2 := uint32(bsize[1]) << 14
	lsb2 := uint32(bsize[2]) << 7
	lsb := uint32(bsize[3])
	var size uint32
	size = msb | msb2 | lsb2 | lsb
	return size
}
