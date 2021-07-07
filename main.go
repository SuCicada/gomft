package main

// http://www.360doc.com/content/19/1030/09/1367418_869900590.shtml
//https://www.zhihu.com/question/22862246
//https://github.com/t9t/gomft
//
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
	fmt.Println(" BytesPerSector", bootSector.BytesPerSector)
	fmt.Println(" SectorsPerCluster", bootSector.SectorsPerCluster)
	fmt.Println(" MftClusterNumber", bootSector.MftClusterNumber)
	bytesPerCluster := bootSector.BytesPerSector *
		bootSector.SectorsPerCluster
	mftPosInBytes := int64(bootSector.MftClusterNumber) *
		int64(bytesPerCluster)
	fmt.Println("mftPosInBytes", mftPosInBytes)

	mftSizeInBytes := bootSector.FileRecordSegmentSizeInBytes

	mftData := make([]byte, mftSizeInBytes)
	_, err = in.Seek(mftPosInBytes, 0)
	rootIndex := mftPosInBytes + int64(mftSizeInBytes*5)
	//for i := 0; i < 10; i++ {
	//	fmt.Println("----------", i, "-------------")
	//	index := mftPosInBytes + int64(mftSizeInBytes*i)
	showFileRecord(in, rootIndex, mftData)
}
func showFileRecord(in *os.File, index int64, mftData []byte) {
	_, _ = in.Seek(index, 0)
	_, _ = io.ReadFull(in, mftData)
	fmt.Println("Parsing $MFT file record", index)

	record, err := mft.ParseRecord(mftData)
	if err != nil {
		log.Println(err.Error())
	}
	//fmt.Println("record", record)
	attrs := record.FindAttributes(mft.AttributeTypeFileName)
	if len(attrs) > 0 {
		res, _ := mft.ParseFileName(attrs[0].Data)
		fmt.Println(res.Name)
	}
	attrs = record.FindAttributes(mft.AttributeTypeIndexRoot)
	if len(attrs) > 0 {
		log.Println("ParseIndexRoot")
		//res, _ := mft.ParseFileReference(attrs[0].Data)
		res, _ := mft.ParseIndexRoot(attrs[0].Data)
		for _, entry := range res.Entries {
			fmt.Println(entry)
			fmt.Print(entry.FileName.Name, " ")
		}
		fmt.Println()
	}
	attrs = record.FindAttributes(mft.AttributeTypeIndexAllocation)
	if len(attrs) > 0 {
		log.Println("ParseIndexRoot")
		res, _ := mft.ParseDataRuns(attrs[0].Data)
		for _, r := range res {
			fmt.Println(r)
			newIndex :=
				showFileRecord(in, newIndex, mftData)

		}
	}
}
