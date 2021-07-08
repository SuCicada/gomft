package main

// http://www.360doc.com/content/19/1030/09/1367418_869900590.shtml
//https://www.zhihu.com/question/22862246
//https://github.com/t9t/gomft
// https://blog.csdn.net/weinierbian/article/details/45649729
// https://blog.csdn.net/ly510587/article/details/100370841
// https://github.com/Hilaver/NtfsResolution/
import (
	"fmt"
	"github.com/t9t/gomft/fragment"
	"github.com/t9t/gomft/mft"
	"io"
	"log"
	"os"

	"github.com/t9t/gomft/bootsect"
)

var bootSector bootsect.BootSector

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

	bootSector, err = bootsect.Parse(buf)
	if err != nil {
		log.Fatalln("Unable to parse boot sector")
	}
	log.Printf("Boot sector of C:\\:\n%+v\n", bootSector)
	const supportedOemId = "NTFS    "
	if bootSector.OemId != supportedOemId {
		log.Fatalln("Unknown OemId (file system type) %q (expected %q)\n", bootSector.OemId, supportedOemId)
	}
	fmt.Println(" BytesPerSector ", bootSector.BytesPerSector)      // 扇区
	fmt.Println(" SectorsPerCluster", bootSector.SectorsPerCluster) // 簇
	fmt.Println(" MftClusterNumber", bootSector.MftClusterNumber)
	bytesPerCluster := bootSector.BytesPerSector *
		bootSector.SectorsPerCluster
	mftPosInBytes := int64(bootSector.MftClusterNumber) *
		int64(bytesPerCluster)
	fmt.Println("==== mftPosInBytes ==== ", mftPosInBytes)

	mftSizeInBytes := bootSector.FileRecordSegmentSizeInBytes
	fmt.Println("===== mftSizeInBytes ====", mftSizeInBytes)
	mftData := make([]byte, mftSizeInBytes)
	//_, err = in.Seek(mftPosInBytes, 0)
	//valumeIndex := mftPosInBytes + int64(mftSizeInBytes*3)
	//_, _ = in.Seek(valumeIndex, 0)
	//_, _ = io.ReadFull(in, mftData)
	//record, err := mft.ParseRecord(mftData)
	//attrs := record.FindAttributes(mft.AttributeTypeVolumeName)
	//res := (attrs[0].Data)
	//fmt.Println("Value", record)
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
		for i, entry := range res.Entries {
			//fmt.Println(entry)
			fmt.Print("[", i, "] ", entry.FileName.Name, " ")
		}
		fmt.Println()
	}
	attrs = record.FindAttributes(mft.AttributeTypeIndexAllocation)
	if len(attrs) > 0 {
		log.Println("ParseDataRuns")
		//fmt.Println(attrs[0].Data)

		dataRuns, _ := mft.ParseDataRuns(attrs[0].Data)
		bytesPerCluster := int64(bootSector.SectorsPerCluster) * int64(bootSector.BytesPerSector)
		frag := mft.DataRunsToFragments(dataRuns, int(bytesPerCluster))
		fmt.Println(frag)
		data := make([]byte)
		fragment.NewReader(in, frag).Read(data)
		fmt.Println(res)
		//for _, r := range res {
		//	fmt.Println(r)
		//	bytesPerCluster := int64(bootSector.SectorsPerCluster) * int64(bootSector.BytesPerSector)
		//	index := r.OffsetCluster * bytesPerCluster
		//	size := r.LengthInClusters * uint64(bytesPerCluster)
		//	data := make([]byte, size)
		//	_, _ = in.Seek(index, 0)
		//	_, _ = io.ReadFull(in, data)
		//	res, _ := mft.ParseDataRuns(data)
		//	fmt.Println(res)
	}
}
