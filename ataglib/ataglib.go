package ataglib

import(
	"fmt"
	"os"
	"encoding/binary"
	"unicode/utf16"
)

type TagType string
const (
	ttID3 TagType = "ID3"
	ttVorbisComment = "VorbisComment"
)

type TagEnc byte
const (
	teISO8152 TagEnc = 0x00
	teUTF16 = 0x01
	teUTF8 = 0x02
)

const(
	ID3V23_HEADERSIZE uint32 = 10
	ID3V23_FRAEMHEADER_SIZE uint32 = 10
)

type AudioMetadata struct {
	Filepath string
	
	TagType TagType
	TagVersion string
	TagFlags byte
	TagSize uint32

	TextTags []TextTag
}

type TextTag struct {
	OrgValue string
	MapValue string
	Enc TagEnc
}

func NewAudioMetadata(fp string, tagHeaderOnly bool) (*AudioMetadata, error) {
	amd := AudioMetadata{Filepath: fp}
	file, err := os.OpenFile(amd.Filepath, os.O_RDONLY, 0644)
	defer file.Close()
	if err != nil {
		return nil, err
	}

	var bytesread int
	bytesread, err = readHeader(file, &amd)
	fmt.Printf("read %d bytes\n", bytesread)
	if tagHeaderOnly {
		return &amd, nil
	}
	
	return &amd, nil
}

func parseID3v23(file *os.File) error {
	var b byte
	var err error

	// read the first frame...
	frameId := "" // 4byte
	for i:=0;i<4;i++ {
		err = binary.Read(file, binary.LittleEndian, &b)
		if err != nil {
			return err
		}
		frameId += string(b)
	}
	fmt.Println(frameId)
	var frameSize uint32
	err = binary.Read(file, binary.BigEndian, &frameSize)
	if err != nil {
		return err
	}
	// frameSize = frameSize - ID3V23_FRAEMHEADER_SIZE
	fmt.Printf("Framesize: %d\n", frameSize)
	
	var frameFlags uint16
	err = binary.Read(file, binary.BigEndian, &frameFlags)
	if err != nil {
		return err
	}
	// var frameDatabytes []byte
	// var i uint32
	// for i=0;i<frameSize;i++ {
	// 	err = binary.Read(file, binary.BigEndian, &b)
	// 	if err != nil {
	// 		return err
	// 	}
	// 	frameDatabytes = append(frameDatabytes, b)
	// }
	switch frameId {
	case "TIT2":
		var encoding byte
		err = binary.Read(file, binary.BigEndian, &encoding)
		if err != nil {
			return err
		}
		fmt.Printf("encoding: 0x%02x\n", encoding)
		if encoding == teUTF16 {
			var bom uint16
			err = binary.Read(file, binary.BigEndian, &bom)
			if err != nil {
				return err
			}
			fmt.Println(bom)
			bytecount := 0
			if bom == 0xFFFE {
				fmt.Println("Little-Endian-Bom 0xFFFE...")
				runes := []rune{}
				for {
					var char uint16 
					err := binary.Read(file, binary.LittleEndian, &char)
					if err != nil {
						return err
					}
					bytecount = bytecount + 2
					if char == 0 {
						fmt. Printf("Bytecount: %d, Framesize: %d\n", bytecount, frameSize)
						break
					}
					runes = append(runes, utf16.Decode([]uint16{char})...)
				}
				fmt.Println(string(runes))
			} else {
				fmt.Println("Big-Endian-Bom 0xFEFF...")
			}
		}
	default:
		fmt.Printf("yet unsupported frameId %s\n", frameId)
	}
	
	
	return nil
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

func bytesToUint16LE(bs [2]byte) uint16 {
	msb := uint16(bs[0])
	lsb := uint16(bs[1]) << 8
	var be uint16
	be = msb | lsb
	return be
}

func readUint16(file *os.File) (uint16, error) {
	bytes := [2]byte{}
	var b byte
	for i:=0;i<2;i++{
		err := binary.Read(file, binary.BigEndian, &b)
		if err != nil {
			return 0, err
		}
		bytes[i] = b
	}
	return bytesToUint16LE(bytes), nil
}

func readHeader(file *os.File, amd *AudioMetadata) (int, error) {
	var b byte
	var err error
	bytesread := 0
	tagtype := ""
	for i:=0;i<=2;i++ {
		err = binary.Read(file, binary.LittleEndian, &b)
		if err != nil {
			return bytesread, err
		}
		bytesread++
		fmt.Printf("%d (0x%02x) - %s\n", b, b, string(b))
		tagtype += string(b)
	}

	switch tagtype {
	case "ID3":
		amd.TagType = ttID3
		var id3tagVersionMajor uint8 = 2
		var id3tagVersionMinor uint8 = 0
 		var id3tagVersionResvision uint8 = 0
		
		fmt.Println("found ID3-Tag")
		err = binary.Read(file, binary.LittleEndian, &id3tagVersionMinor)
		if err != nil {
			return bytesread, err
		}
		bytesread++
		err = binary.Read(file, binary.LittleEndian, &id3tagVersionResvision)
		if err != nil {
			return bytesread, err
		}
		bytesread++
			fmt.Printf("found ID3v%d.%d.%d Audio-Metadata...\n",
				id3tagVersionMajor, id3tagVersionMinor, id3tagVersionResvision)
		amd.TagVersion = fmt.Sprintf("%sv%d.%d.%d",
			string(amd.TagType),
			id3tagVersionMajor, id3tagVersionMinor, id3tagVersionResvision)
			
		// ID3v2 flags             %abc00000
		err = binary.Read(file, binary.LittleEndian, &amd.TagFlags)
		if err != nil {
			return bytesread, err
		}
		bytesread++
		if amd.TagFlags != 0x00 {
			return bytesread, fmt.Errorf("ID3v23: tags with 'FLAGS' is not supported yet, sorry...")
		}
		// ID3v2 size              4 * %0xxxxxxx
		var rawSize [4]byte
		for i:=0;i<4;i++ {
			err = binary.Read(file, binary.BigEndian, &rawSize[i])
			if err != nil {
				return bytesread, err
			}
		}
		amd.TagSize = id3v23bytesizeToUint32(rawSize) - ID3V23_HEADERSIZE

	default:
		fmt.Println("tag-type not supported yet: ", tagtype)
	}
	return bytesread, nil
}
