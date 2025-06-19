package ataglib

import(
	"fmt"
	"os"
	"encoding/binary"
	"unicode/utf16"
)

func (amd *AudioMetadata)parseID3Header(file *os.File, bytesread *uint32) error {
	var b byte
	var err error
	var id3tagVersionMajor uint8 = 2
	var id3tagVersionMinor uint8 = 0
	var id3tagVersionResvision uint8 = 0

	err = binary.Read(file, binary.LittleEndian, &id3tagVersionMinor)
	if err != nil {
		return err
	}
	*bytesread++
	err = binary.Read(file, binary.LittleEndian, &id3tagVersionResvision)
	if err != nil {
		return err
	}
	*bytesread++
	amd.TagVersion = fmt.Sprintf("%sv%d.%d.%d",
		string(amd.TagIdentifier),
		id3tagVersionMajor, id3tagVersionMinor, id3tagVersionResvision)

	// header flags
	err = binary.Read(file, binary.LittleEndian, &b)
	if err != nil {
		return err
	}
	*bytesread++
	if b == 0x00 {
		amd.TagHeaderFlagUnsyncronisation = false
		amd.TagHeaderFlagExtendedHeader = false
		amd.TagHeaderFlagExperimentalIndicator = false
		amd.TagHeaderFlagFooterPresent = false
	} else {
		fmt.Println("WARN sorry, Tag-Header Flags for ID3v2-Tags are not supported yet...")
	}

	// ID3v2 size: 4 * %0xxxxxxx
	var rawSize [4]byte
	for i:=0;i<4;i++ {
		err = binary.Read(file, binary.BigEndian, &rawSize[i])
		if err != nil {
			return err
		}
		*bytesread++
	}
	amd.TagSize = id3v23bytesizeToUint32(rawSize) - ID3V2_HEADERSIZE

	return nil
}

func (amd *AudioMetadata)parseID3Frames(file *os.File, bytesread *uint32) error {
	var b byte
	var err error
	for *bytesread < (amd.TagSize - ID3V2_HEADERSIZE) {
		// read frameheader 10 byte
		var frameId string
		var frameSize uint32
		var frameFlags uint16

		// frameID 4byte
		for i:=0;i<4;i++ {
			err = binary.Read(file, binary.LittleEndian, &b)
			if err != nil {
				return err
			}
			*bytesread++
			if b == 0x00 {
				// fmt.Println("starting crazy 0x00 bytes... perhaps padding!?")
				for b == 0x00 {
					err = binary.Read(file, binary.LittleEndian, &b)
					if err != nil {
						return err
					}
					*bytesread++
					if *bytesread >= (amd.TagSize - ID3V2_HEADERSIZE) {
						fmt.Println("-------------------------------")
						return nil
					}
				}
			}
			frameId += string(b)
		}
		// fmt.Println("Frame-ID:", frameId)
		
		// frameSize = frameSize - ID3V2_FRAMEHEADER_SIZE
		err = binary.Read(file, binary.BigEndian, &frameSize)
		if err != nil {
			return err
		}
		*bytesread = *bytesread + 4
		// fmt.Printf("Framesize: %d\n", frameSize)

		// FRAME FLAGS 2byte
		err = binary.Read(file, binary.BigEndian, &frameFlags)
		if err != nil {
			return err
		}
		*bytesread = *bytesread + 2
		if frameFlags != 0x0000 {
			fmt.Println("FrameFlags not supported yet, sorry...")
		}
		if frameId == "TXXX" {
			fmt.Println(frameId, "not supported yet:")
			dummy := make([]byte, int(frameSize))
			err = binary.Read(file, binary.LittleEndian, &dummy)
			if err != nil {
				return err
			}
			*bytesread = *bytesread + frameSize
		}else	if frameId[0] == 'T' {
			te, val, err := parseTextFrame(file, frameSize)
			if err != nil {
				return err
			}
			fmt.Printf("Tag: %s -> Value: %s (%d)\n", frameId, val, te)
			*bytesread = *bytesread + frameSize
		} else {
			fmt.Println(frameId, "not supported yet:")
			dummy := make([]byte, int(frameSize))
			err = binary.Read(file, binary.LittleEndian, &dummy)
			if err != nil {
				return err
			}
			*bytesread = *bytesread + frameSize
			// fmt.Printf("dummy-length: %d - frameSize: %d\n", len(dummy), frameSize)
		}
	}
	fmt.Println("-------------------------------")
	return nil
}

func parseTextFrame(file *os.File, frameSize uint32) (TagEnc, string, error) {
	var te TagEnc
	var val string

	var b byte
	var err error
	var bytecount uint32 = 0
	
	err = binary.Read(file, binary.LittleEndian, &b)
	if err != nil {
		return te, val, err
	}
	bytecount++
	if b > 0x03 {
		fmt.Println("TagEncoding greater as defined, use 0x01: ", b)
		te = 0x01
	} else {
		te = TagEnc(b)
	}
	
	switch te{
		case TE_ISO8152:
		bytes := make([]byte, int(frameSize - 1))
		err = binary.Read(file, binary.LittleEndian, &bytes)
		if err != nil {
			return te, val, err
		}
		bytecount = bytecount + frameSize - 1
		val = string(bytes)
		case TE_UTF16BOM:
		var bom uint16
		err = binary.Read(file, binary.BigEndian, &bom)
		if err != nil {
			return te, val, err
		}
		frameSize--

		if bom == 0xFFFE {
			runes := []rune{}
			for {
				var char uint16 
				err := binary.Read(file, binary.LittleEndian, &char)
				if err != nil {
					return te, val, err
				}
				bytecount = bytecount + 2
				if char == 0x0000 {
					break
				}
				runes = append(runes, utf16.Decode([]uint16{char})...)
			}
			val = string(runes)
		} else if bom == 0xFEFF {
			fmt.Println("Big-Endian-BOM 0xFEFF not supported yet :(")
			return te, "", nil
		}
		default:
		fmt.Println("TextEncoding not supported yet:", te)
		return te, "", nil
	}
	return te, val, err
}

