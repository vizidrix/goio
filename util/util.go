package main

import (
	"github.com/vizidrix/goio"
	"time"
)

func main() {
	MakeCert()
}

func MakeCert() {
	var err error
	var cert *goio.Cert
	private := "private.pem"
	public := "public.pem"
	if cert, err = goio.MakeCert("intel", 2048, []string{"localhost"}, 30*time.Minute, false); err != nil {
		panic(err)
	}
	cert.WritePrivate(private)
	cert.WritePublic(public)
}
	/*
	var err error
	if len(os.Args) > 0 {
		switch os.Args[1] {
		case "pack":
			TarPackDir(os.Args[2:])
			break
		case "cert":
			MakeCert(os.Args[2:])
			return
		}
	}
	*/
//}

/*
func TarPackDir(args []string) {
	switch args[0] {
	case "toml":
		TarPackToml(args[1:])
		return
	case "web":
		TarPackWeb(args[1:])
		return
	default:
		TarPack(args)
	}
}

func tarPack(targetDir, outFile string, filters ...tp.FilePredicate) {
	var buffer []byte
	if buffer, err = tp.TarDir(targetDir, "", filters...); err != nil {
		fmt.Printf("Error zipping  TarPack files")
	}
	ioutil.WriteFile(outFile, buffer, 0666)
}

func TarPackToml(args []string) {
	targetDir := ""
	outFile := "toml.tar"
	if len(args) > 0 {
		targetDir = args[0]
		outFile = args[1]
	}
	fileFilter := tp.OnFiles(tp.PathMatchAny("*.toml"))
	var buffer []byte
	if buffer, err = tp.TarDir(targetDir, "", folderFilter, fileFilter); err != nil {
		fmt.Printf("Error zipping  TarPack files")
	}
	ioutil.WriteFile(outFile, buffer, 0666)
}

func TarPackWeb(args []string) {
	targetDir := ""
	outFile := "web.tar"
	if len(args) > 0 {
		targetDir = args[0]
		outFile = args[1]
	}
	folderFilter := tp.OnDirs(tp.Not(tp.NameContainsAny(".git")))
	fileFilter := tp.OnFiles(tp.PathMatchAny("*.js", "*.html", "*.png", "*.css", "*.ttf"))
	var buffer []byte
	if buffer, err = tp.TarDir(targetDir, "", folderFilter, fileFilter); err != nil {
		fmt.Printf("Error zipping Web files")
	}
	ioutil.WriteFile(outFile, buffer, 0666)
}
*/

	/*
	fileWhiteList := []string{
		"*.js",
		"*.html",
		"*.png",
		"*.css",
		"*.ttf",
	}
	*/
/*
	if len(os.Args) > 1 {
		var cert *goio.Cert
		if os.Args[1] == "cert" {
			if cert, err = goio.MakeCert("intel", 2048, []string{"localhost"}, 30*time.Minute, false); err != nil {
				panic(err)
			}
			cert.WritePrivate("private.pem")
			cert.WritePublic("public.pem")
			return
		}
	}
	*/
	//decodeString := "intelinde"

	//hash := sha256.Sum256([]byte(decodeString))
	//key := hash[:]

	//buffer := new(bytes.Buffer)
	//aes_w_handle, _ := aes.NewWriter(buffer, key)
	//zip_w_handle, _ := gzip.NewWriterLevel(aes_w_handle, gzip.BestCompression)

	//zip_w_handle, _ := gzip.NewWriterLevel(buffer, gzip.BestCompression)
	//tar_w_handle := tar.NewWriter(zip_w_handle)

	//tar_w_handle := tar.NewWriter(buffer)
	//writeDir(tar_w_handle, targetDirPath, "", fileWhiteList)



	/*
		files, _ := ioutil.ReadDir(targetDirPath)
		for _, file := range files {
			if !file.IsDir() {
				fmt.Printf("Writing file: %s\n", file.Name())
				if err = writeFile(tar_w_handle, targetDirPath+"/"+file.Name()); err != nil {
					fmt.Printf("Error writing file: %s\n", err)
					return
				}
			}
		}
	*/

		/*
	if err = tar_w_handle.Close(); err != nil {
		fmt.Printf("Error closing tar:\n\t- %s\n", err)
	}
	*/
	/*
		if err = zip_w_handle.Close(); err != nil {
			fmt.Printf("Error closing zip:\n\t- %s\n", err)
		}
	*/
	//aes_w_handle.Close()

	
//}

/*
func isWhitelistedFile(whitelist []string, name string) bool {
	for _, entry := range whitelist {
		if match, _ := filepath.Match(entry, name); match {
			return true
		}
	}
	return false
}

func writeDir(handle *tar.Writer, rootPath, relPath string, fileWhiteList []string) (err error) {
	fmt.Printf("Writing dir: %s => %s\n", rootPath, relPath)
	var files []os.FileInfo
	if files, err = ioutil.ReadDir(rootPath); err != nil {
		return
	}
	for _, file := range files {
		newPath := relPath
		if len(relPath) != 0 && relPath[:len(relPath)] != "/" {
			newPath += "/"
		}
		if file.IsDir() { // File is a dir and not a file
			if strings.Contains(file.Name(), ".git") {
				continue // TODO: Improve this filter
			}
			writeDir(handle, rootPath+"/"+file.Name(), newPath+file.Name(), fileWhiteList)
		} else { // File is a file and not a dir
			if !isWhitelistedFile(fileWhiteList, file.Name()) {
				continue
			}
			fmt.Printf("Writing file: %s -> %s\n", relPath, file.Name())
			if err = writeFile(handle, rootPath+"/"+file.Name(), newPath+file.Name()); err != nil {
				fmt.Printf("Error writing file: %s\n", err)
				return
			}
		}
	}
	return
}

func writeFile(handle *tar.Writer, rootPath, relPath string) error {
	var file *os.File
	var stat os.FileInfo
	var buffer []byte
	var err error
	if file, err = os.OpenFile(rootPath, os.O_RDONLY, os.ModePerm); err != nil {
		return err
	}
	if stat, err = file.Stat(); err != nil {
		return err
	}
	file.Close()
	header := &tar.Header{
		Name: relPath,
		Size: stat.Size(),
	}
	if err := handle.WriteHeader(header); err != nil {
		return err
	}
	if buffer, err = ioutil.ReadFile(rootPath); err != nil {
		return err
	}
	if _, err := handle.Write(buffer); err != nil {
		return err
	}
	return nil
}
*/