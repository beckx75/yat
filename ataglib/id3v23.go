package ataglib

import(
	"fmt"
	"os"
	"encoding/binary"
	"unicode/utf16"
)

func parseID3v23(file *os.File) error {
	var b byte
	var err error

	fmt.Println("ich bin hiiiiier :)")
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
