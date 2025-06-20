package ataglib

import(
	"fmt"
	"os"
	//	"log/slog"
	"encoding/binary"
)

type FileIdentifier string
const (
	FI_ID3 FileIdentifier = "ID3"
	FI_FLAC = "FLAC"
)

type TagEnc byte
const (
	TE_ISO8152 TagEnc = 0x00 // ID3v2.4 terminated with 0x00
	TE_UTF16BOM = 0x01  //with BOM -> terminated with 0x00 00
	TE_UTF16 = 0x02 // wihout BOM -> terminated with 0x00 00
	TE_UTF8 = 0x03 // ID3v2.4 terminated with 0x00
)

const(
	ID3V2_HEADERSIZE uint32 = 10
	ID3V2_FRAMEHEADER_SIZE uint32 = 10
)

type AudioMetadata struct {
	Filepath string

	TagIdentifier FileIdentifier
	// ID3 only
	TagVersion string
	TagHeaderFlagUnsyncronisation bool // Bit 7
	TagHeaderFlagExtendedHeader bool // Bit 6
	TagHeaderFlagExperimentalIndicator bool  // Bit 5
	// ID3v2.4 only
	TagHeaderFlagFooterPresent bool // Bit 4
	
	TagSize uint32
}

func NewAudioMetadata(fp string, tagHeaderOnly bool) (*AudioMetadata, error) {
	amd := AudioMetadata{Filepath: fp}
	file, err := os.OpenFile(amd.Filepath, os.O_RDONLY, 0644)
	defer file.Close()
	if err != nil {
		return nil, err
	}
	var b byte
	var bytesread uint32 = 0

	// first 3 bytes of audiofile (ID3, fLa)
	var fileIdentifier string
	for i:=0;i<3;i++ {
		err = binary.Read(file, binary.BigEndian, &b)
		if err != nil {
			return nil, err
		}
		fileIdentifier = fmt.Sprintf("%s%s", fileIdentifier, string(b))
	}
	bytesread =+ 3
	switch fileIdentifier{
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


	
	// bytesread, err = readHeader(file, &amd)
	// if bytesread != 10 {
	// 	return nil, fmt.Errorf("read bytes for header are not equal to 10: %d",
	// 		bytesread)
	// }
	// if tagHeaderOnly {
	// 	return &amd, nil
	// }
	// switch amd.TagVersion {
	// case "ID3v2.3.0":
	// 	err = parseID3v23(file)
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// default:
	// 	slog.Warn("sorry, TagVersion not supported yet :(", "Tagversion", amd.TagVersion)
	// }
	return &amd, nil
}

// func readHeader(file *os.File, amd *AudioMetadata) (int, error) {
// 	var b byte
// 	var err error
// 	bytesread := 0
// 	tagtype := ""
// 	for i:=0;i<=2;i++ {
// 		err = binary.Read(file, binary.LittleEndian, &b)
// 		if err != nil {
// 			return bytesread, err
// 		}
// 		bytesread++
// 		tagtype += string(b)
// 	}

// 	switch tagtype {
// 	case "ID3":
// 		amd.TagType = "ttID3"
// 		var id3tagVersionMajor uint8 = 2
// 		var id3tagVersionMinor uint8 = 0
//  		var id3tagVersionResvision uint8 = 0
		
// 		err = binary.Read(file, binary.LittleEndian, &id3tagVersionMinor)
// 		if err != nil {
// 			return bytesread, err
// 		}
// 		fmt.Println(id3tagVersionMinor)
// 		bytesread++
// 		err = binary.Read(file, binary.LittleEndian, &id3tagVersionResvision)
// 		if err != nil {
// 			return bytesread, err
// 		}
// 		bytesread++
// 		amd.TagVersion = fmt.Sprintf("%sv%d.%d.%d",
// 			string(amd.TagType),
// 			id3tagVersionMajor, id3tagVersionMinor, id3tagVersionResvision)
			
// 		// ID3v2 flags             %abc00000
// 		err = binary.Read(file, binary.LittleEndian, &amd.TagFlags)
// 		if err != nil {
// 			return bytesread, err
// 		}
// 		bytesread++
// 		if amd.TagFlags != 0x00 {
// 			return bytesread, fmt.Errorf("ID3v23: tags with 'FLAGS' is not supported yet, sorry...")
// 		}
// 		// ID3v2 size              4 * %0xxxxxxx
// 		var rawSize [4]byte
// 		for i:=0;i<4;i++ {
// 			err = binary.Read(file, binary.BigEndian, &rawSize[i])
// 			if err != nil {
// 				return bytesread, err
// 			}
// 			bytesread++
// 		}
// 		amd.TagSize = id3v23bytesizeToUint32(rawSize) - ID3V23_HEADERSIZE

// 	default:
// 		fmt.Println("tag-type not supported yet: ", tagtype)
// 	}
// 	return bytesread, nil
// }

func id3v23bytesizeToUint32(bsize [4]byte) uint32 {
	msb := uint32(bsize[0]) << 21
	msb2 := uint32(bsize[1]) << 14
	lsb2 := uint32(bsize[2]) << 7
	lsb := uint32(bsize[3])
	var size uint32
	size = msb | msb2 | lsb2 | lsb
	return size
}

// func bytesToUint16LE(bs [2]byte) uint16 {
// 	msb := uint16(bs[0])
// 	lsb := uint16(bs[1]) << 8
// 	var be uint16
// 	be = msb | lsb
// 	return be
// }

// func readUint16(file *os.File) (uint16, error) {
// 	bytes := [2]byte{}
// 	var b byte
// 	for i:=0;i<2;i++{
// 		err := binary.Read(file, binary.BigEndian, &b)
// 		if err != nil {
// 			return 0, err
// 		}
// 		bytes[i] = b
// 	}
// 	return bytesToUint16LE(bytes), nil
// }

