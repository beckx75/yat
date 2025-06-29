package ataglib

import(
	"fmt"
	"os"
	"encoding/binary"
	"unicode/utf16"
	"github.com/rs/zerolog/log"
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
		if (b & 0x80) == 0x80 {
			amd.TagHeaderFlagUnsyncronisation = true
		}
		if (b & 0x40) == 0x40 {
			amd.TagHeaderFlagExtendedHeader = true
			if amd.TagVersion == "ID3v2.4.0" {
				log.Warn().Msg("Extended Header not supported yet, sorry...")
			} else {
				log.Warn().Msg("Extended Header not supported yet, sorry...")
			}
		}
		if (b & 0x20) == 0x20 {
			amd.TagHeaderFlagExperimentalIndicator = true
		}
		if amd.TagVersion == "ID3v2.4.0" {
			if (b & 0x10) == 0x10 {
				amd.TagHeaderFlagFooterPresent = true
			}
		}
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
	fmt.Println("amd.Tagsize befor ID3-Frame parsing:", amd.TagSize)

	return nil
}

func (amd *AudioMetadata)parseID3Frames(file *os.File, bytesread *uint32) error {
	var b byte
	var err error
	// for *bytesread < (amd.TagSize - ID3V2_HEADERSIZE) {
	for {
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
		fmt.Println("Framesize:", frameSize)
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
			te, desc, val, err := parseTXXXFrame(file, frameSize)
			if err != nil {
				return err
			}
			fmt.Printf("Tag: TXXX -> Desc: %s\t Value: %s (%d)\n", desc, val, te)
			*bytesread = *bytesread + frameSize
		}else	if frameId[0] == 'T' {
			te, val, err := parseTextFrame(file, frameSize)
			if err != nil {
				return err
			}
			fmt.Printf("Tag: %s -> Value: %s (%d)\n", frameId, val, te)
			*bytesread = *bytesread + frameSize
		} else if frameId == "APIC" {
			fmt.Println(frameId)
			fmt.Println("Framesize:", frameSize)
			fmt.Println("Bytea read:",*bytesread)
			fmt.Println("Sum:", *bytesread + frameSize)
			fmt.Println("Tagsize:", amd.TagSize)
			foundT := false
			foundTX := false
			var dummy byte
			for {
				err = binary.Read(file, binary.LittleEndian, &dummy)
				if err != nil {
					return err
				}
				*bytesread = *bytesread + 1
				if (string(dummy) == "T") && (!foundT) {
					foundT = true
				}
				if (string(dummy) == "X") && (foundT) && (!foundTX) {
					foundTX = true
				}
				if (string(dummy) == "X") && (foundT) && (foundTX) {
					fmt.Println("Bytesread so far:", *bytesread)
					panic("founc TXX")
				}
			}
			*bytesread = *bytesread + frameSize
			var x byte
			for i:=0; i<10; i++ {
				_ = binary.Read(file, binary.BigEndian, &x)
				fmt.Println("Nachlese:", string(x))
			}
		}	else {
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
		log.Error().Msgf("TagEncoding greater as defined, use 0x01: ", b)
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
			log.Warn().Msg("Big-Endian-BOM 0xFEFF not supported yet :(")
			return te, "", nil
		}
		default:
		log.Warn().Msgf("TextEncoding not supported yet:", te)
		return te, "", nil
	}
	return te, val, err
}

func parseTXXXFrame(file *os.File, frameSize uint32) (TagEnc, string, string, error) {
	var te TagEnc
	var desc string
	var val string

	var b byte
	var err error
	var bytecount uint32 = 0

	err = binary.Read(file, binary.LittleEndian, &b)
	if err != nil {
		return te, val, desc, err
	}
	bytecount++
	if b > 0x03 {
		log.Error().Msgf("TagEncoding greater as defined, use 0x01: ", b)
		te = 0x01
	} else {
		te = TagEnc(b)
	}

	bytes := []byte{}
	// read DESCRIPTION
	for {
		err = binary.Read(file, binary.LittleEndian, &b)
		if err != nil {
			return te, val, desc, err
		}
		bytecount++
		if (b == 0x00) && (te == TE_ISO8152) {
			desc = string(bytes)
			break
		}
		bytes = append(bytes, b)
		if bytecount == frameSize {
			break
		}
	}
	
	// read Value
	valsize := frameSize - bytecount
	bval := make([]byte, valsize)
	err = binary.Read(file, binary.LittleEndian, &bval)
	if err != nil {
		return te, val, desc, err
	}
	
	val = string(bval)
	
	return te, val, desc, nil
}

func parsePRIVFrame(file *os.File, frameSize uint32) ( string, []byte, error) {
	// <Header for 'Private frame', ID: "PRIV">
	// 	Owner identifier        <text string> $00
	// The private data        <binary data>
	oid := ""
	pd := []byte{}
	
	return oid, pd, nil
}
