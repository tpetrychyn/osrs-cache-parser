package fs

type FSFile struct {
	FileId   uint16
	NameHash int32
	Contents []byte
}
