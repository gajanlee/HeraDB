package HeraDB

type Header struct  {

}

type Record struct {

}

type Block struct {
	schema *uint32
	header Header
	length uint32
	timestamp uint32
	records []Record
}