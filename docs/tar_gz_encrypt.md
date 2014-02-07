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
