package shell

import "crypto/md5"
import "encoding/hex"
import "log"
import "io"
import "os"

func hash_file_md5(filePath string) (string, error) {
	var returnMD5String string
	file, err := os.Open(filePath)
	if err != nil {
		return returnMD5String, err
	}
	defer file.Close()
	hash := md5.New()
	if _, err := io.Copy(hash, file); err != nil {
		return returnMD5String, err
	}
	hashInBytes := hash.Sum(nil)[:16]
	returnMD5String = hex.EncodeToString(hashInBytes)
	return returnMD5String, nil

}

// hash thyself:  hash_file_md5(os.Args[0])
func Md5sum(filename string) string {
	filename = Path(filename)
	log.Println("shell.Md5sum() START filename =", filename)
	hash, err := hash_file_md5(filename)
	if err == nil {
		log.Println("shell.Md5sum() hash =", hash)
		return hash
	}
	log.Println("shell.Md5sum() failed")
	return ""
}
