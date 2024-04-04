package meta

type File struct {
	Sha1     string
	Name     string
	Size     int64
	Location string
	UploadAt string
}

var fileMetas map[string]File

func init() {
	fileMetas = make(map[string]File)
}
func UpdateFileMeta(fmeta File) {
	fileMetas[fmeta.Sha1] = fmeta
}

func GetFileMeta(fileSha1 string) File {
	return fileMetas[fileSha1]
}

func DeleteFileMeta(fileSha1 string) {
	delete(fileMetas, fileSha1)
}
