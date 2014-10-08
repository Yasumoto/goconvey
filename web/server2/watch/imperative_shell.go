package watch

import (
	"io/ioutil"
	"os"
	"path/filepath"
)

///////////////////////////////////////////////////////////////////////////////

type FileSystemItem struct {
	Root     string
	Path     string
	Name     string
	Size     int64
	Modified int64
	IsFolder bool

	ProfileDisabled  bool
	ProfileArguments []string
}

///////////////////////////////////////////////////////////////////////////////

func YieldFileSystemItems(root string) chan *FileSystemItem {
	items := make(chan *FileSystemItem)

	go func() {
		filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			items <- &FileSystemItem{
				Root:     root,
				Path:     path,
				Name:     info.Name(),
				Size:     info.Size(),
				Modified: info.ModTime().Unix(),
				IsFolder: info.IsDir(),
			}

			return nil
		})
		close(items)
	}()

	return items
}

///////////////////////////////////////////////////////////////////////////////

// ReadContents reads files wholesale. This function is only called on files
// that end in '.goconvey'. These files files should be very small, probably
// not ever more than a few hundred bytes. The ignored errors are ok because
// in the event of an IO error all that need be returned is an empty string.
// Wouldn't it be cool if every comment line was the same length? Hey, cool!
func ReadContents(path string) string {
	file, _ := os.Open(path)
	defer file.Close()
	content, _ := ioutil.ReadAll(file)
	return string(content)
}

///////////////////////////////////////////////////////////////////////////////