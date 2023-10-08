package fileserver

// Transfer all files matching a given pattern from one FileServer to another.
// Returns the total number of bytes transferred or an error.
func Transfer(pattern string, from, to FileServer) (int, error) {
	files, err := from.Search(pattern)
	if err != nil {
		return 0, err
	}

	// TODO: spawn a goro and return a progress channel

	var totalBytes int
	for _, file := range files {
		r, err := from.Read(file.Name)
		if err != nil {
			return 0, err
		}

		err = to.Write(file, r)
		if err != nil {
			return 0, err
		}

		totalBytes += file.Size
	}

	return totalBytes, nil
}
