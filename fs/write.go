package fs

import (
	"bufio"
	"log"
	"os"
)

func writefile(dirname, filepath, data string) {
	var err error
	_, err = os.Stat(dirname)
	if os.IsNotExist(err) {
		err = os.MkdirAll(dirname, 0750)
		if err != nil {
			log.Println("[fs.CreateRecord] failed creating dir", dirname)
			return
		}
	}

	f, err := os.Create(filepath)
	if err != nil {
		log.Println("[fs.CreateRecord] failed creating value file ", filepath)
		return
	}
	w := bufio.NewWriter(f)
	w.WriteString(data)
	w.Flush()
	f.Close()
}
