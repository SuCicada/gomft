package main

// http://www.360doc.com/content/19/1030/09/1367418_869900590.shtml
//https://www.zhihu.com/question/22862246
//https://github.com/t9t/gomft

import (
	"fmt"
	"github.com/t9t/gomft/mft"
	"io"
	"log"
	"os"

	"github.com/t9t/gomft/bootsect"
)

func main() {
	log.SetFlags(log.Lshortfile | log.LstdFlags)

	in, err := os.Open(`\\.\C:`)
	if err != nil {
		log.Fatalln("Unable to open C:", err)
	}
	defer in.Close()

	buf := make([]byte, 512)
	_, err = io.ReadFull(in, buf)
	if err != nil {
		log.Fatalln("Unable to read bootsector data", err)
	}

	bootSector, err := bootsect.Parse(buf)
	if err != nil {
		log.Fatalln("Unable to parse boot sector")
	}
	log.Printf("Boot sector of C:\\:\n%+v\n", bootSector)
	const supportedOemId = "NTFS    "
	if bootSector.OemId != supportedOemId {
		log.Fatalln("Unknown OemId (file system type) %q (expected %q)\n", bootSector.OemId, supportedOemId)
	}
	bytesPerCluster := bootSector.BytesPerSector *
		bootSector.SectorsPerCluster
	mftPosInBytes := int64(bootSector.MftClusterNumber) *
		int64(bytesPerCluster)
	fmt.Println("mftPosInBytes", mftPosInBytes)

	mftSizeInBytes := bootSector.FileRecordSegmentSizeInBytes

	mftData := make([]byte, mftSizeInBytes)
	_, err = in.Seek(mftPosInBytes, 0)

	for i := 0; i < 10000; i++ {
		_, err = in.Seek(mftPosInBytes+int64(mftSizeInBytes*i), 0)
		_, err = io.ReadFull(in, mftData)
		fmt.Printf("Parsing $MFT file record\n")

		record, err := mft.ParseRecord(mftData)
		if err != nil {
			log.Println(err.Error())
		}
		attrs := record.FindAttributes(mft.AttributeTypeFileName)
		if len(attrs) > 0 {
			fileName, _ := mft.ParseFileName(attrs[0].Data)
			fmt.Println(fileName)
		}
	}
}
