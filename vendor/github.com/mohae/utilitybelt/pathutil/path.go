// Package pathutil provides some simple path related functions.  Before using, please
// check the standard library's path and path/filepath functions to see if there is
// something more appropraite.
package pathutil

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Dir is a container for a list of filenames in a directory.
type Dir struct {
	Files []file
}

//  file contains information about a file
type file struct {
	Path string
	Info os.FileInfo
}

// Walk walks the passed path, making a list of all the files that are
// children of the path.
func (d *Dir) Walk(path string) error {
	// If the directory exists, create a list of its contents.
	if path == "" {
		return nil
	}

	// See if the path exists
	exists, err := PathExists(path)
	if err != nil {
		return err
	}

	if !exists {
		err = fmt.Errorf("%s does not exist", path)
		return err
	}

	fullPath, err := filepath.Abs(path)
	if err != nil {
		return err
	}

	// Set up the callback function.
	callback := func(p string, fi os.FileInfo, err error) error {
		return d.addFile(fullPath, p, fi, err)
	}

	// Walk the tree.
	return filepath.Walk(fullPath, callback)
}

// addFile adds the file to the Files slice
func (d *Dir) addFile(root string, p string, fi os.FileInfo, err error) error {
	// See if the path exists
	var exists bool
	exists, err = PathExists(p)
	if err != nil {
		return err
	}

	if !exists {
		err = fmt.Errorf("%s does not exist", p)
		return err
	}

	rel := ""
	// Get the relative information.
	rel, err = filepath.Rel(root, p)
	if err != nil {
		return err
	}

	// skip relative root
	if rel == "." {
		return nil
	}

	d.Files = append(d.Files, file{Path: rel, Info: fi})

	return nil
}

// PathExists returns true if the given path exists, otherwise false. If an
// error is encountered, that is returned, otherwise error will be nil.
//
// Since it not existing will cause an os.Stat error, we need to check if the
// error passed back is the IsNotExist error, which would mean the path does
// not exist, instead of an ectual error state.
func PathExists(p string) (bool, error) {
	_, err := os.Stat(p)
	if err == nil {
		return true, nil
	}

	if os.IsNotExist(err) {
		return false, nil
	}

	return false, err
}

// AppendSlash appends a slash, `/`,  to the end of the passed string, if it
// doesn't already end with one. For OS independend filepaths, `/` are always
// used internally with OS-specific conversion occuring before interacting
// with the filesystem via filepath.FromSlash() and filepath.ToSlash()
func AppendSlash(s string) string {
	if strings.HasSuffix(s, "/") {
		return s
	}

	return s + "/"
}
