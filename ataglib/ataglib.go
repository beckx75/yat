package ataglib

import(
	"fmt"
	"os"
	"encoding/binary"
	//	"golang.org/x/text/encoding/unicode"
)

type TagEnc byte
const (
	ISO8152 TagEnc = 0x00
	UTF16 = 0x01
	UTF8 = 0x02
)

const(
	ID3V23_HEADERSIZE uint32 = 10
	ID3V23_FRAEMHEADER_SIZE uint32 = 10
)

type AudioMetadata struct {
	Filepath string
	TextTags []TextTag
}

type TextTag struct {
	OrgValue string
	MapValue string
	Enc TagEnc
}

func NewAudioAudiometadata(fp string) (*AudioMetadata, error) {
	amd := AudioMetadata{Filepath: fp}
	file, err := os.OpenFile(amd.Filepath, os.O_RDONLY, 0644)
	defer file.Close()
	if err != nil {
		return nil, err
	}

	var b byte
	tagtype := ""
	for i:=0;i<=2;i++ {
		err = binary.Read(file, binary.LittleEndian, &b)
		if err != nil {
			return nil, err
		}
		fmt.Printf("%d (0x%02x) - %s\n", b, b, string(b))
		tagtype += string(b)
	}

	switch tagtype {
	case "ID3":
		fmt.Println("found ID3-Tag")
		var id3v2Version uint8
		err = binary.Read(file, binary.LittleEndian, &id3v2Version)
		if err != nil {
			return nil, err
		}
		if id3v2Version == 2 {
			fmt.Println("not implemented yet...")
		} else if id3v2Version == 3 {
			err = parseID3v23(file)
			if err != nil {
				return nil, err
			}
		} else if id3v2Version == 4 {
			fmt.Println("not implemented yet...")
		} else {
			return nil, fmt.Errorf("malformed ID3v2-Tag Entry... unknown Version %02d", id3v2Version)
		}
	default:
		fmt.Printf("unsupported Tag-Type %s\n", tagtype)
	}
	
	return &amd, nil
}

func parseID3v23(file *os.File) error {
	var b byte
	var err error
	err = binary.Read(file, binary.LittleEndian, &b)
	if err != nil {
		return err
	}
	if b != 0 {
		return fmt.Errorf("malformated ID3v2.3-Tag Entry...")
	}
	// ID3v2 flags             %abc00000
	var flagsByte byte
	err = binary.Read(file, binary.LittleEndian, &flagsByte)
	if err != nil {
		return err
	}
	if flagsByte != 0x00 {
		return fmt.Errorf("ID3v23: tags with 'FLAGS' is not supported yet, sorry...")
	}
	// ID3v2 size              4 * %0xxxxxxx
	var rawSize [4]byte
	for i:=0;i<4;i++ {
		err = binary.Read(file, binary.BigEndian, &rawSize[i])
		if err != nil {
			return err
		}
	}
	id3v2tagsize := id3v23bytesizeToUint32(rawSize) - ID3V23_HEADERSIZE
	fmt.Println(id3v2tagsize)

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
	frameSize = frameSize - ID3V23_FRAEMHEADER_SIZE
	fmt.Println(frameSize)
	
	var frameFlags uint16
	err = binary.Read(file, binary.BigEndian, &frameFlags)
	if err != nil {
		return err
	}
	var frameDatabytes []byte
	var i uint32
	for i=0;i<frameSize;i++ {
		err = binary.Read(file, binary.BigEndian, &b)
		if err != nil {
			return err
		}
		frameDatabytes = append(frameDatabytes, b)
	}
	switch frameId {
	case "TIT2":
		encoding := frameDatabytes[0]
		fmt.Printf("encoding: 0x%02x\n", encoding)
		if encoding == UTF16 {
			bom := bytesToUint16BE(frameDatabytes[1:3])
			fmt.Println(bom)
			if bom == 0xFFFE {
				fmt.Println("Little-Endian-Bom 0xFFFE...")
			} else {
				fmt.Println("Big-Endian-Bom 0xFEFF...")
			}
		}
		fmt.Printf("Title:", string(frameDatabytes[1:]))
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

func bytesToUint16BE(bs []byte) uint16 {
	msb := uint16(bs[0]) << 8
	lsb := uint16(bs[1])
	var be uint16
	be = msb | lsb
	return be
}
