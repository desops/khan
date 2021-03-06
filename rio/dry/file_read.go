package dry

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"syscall"

	"khan.rip/rio/util"
)

type File struct {
	info    *util.FileInfo // nil info means file not present (deleted)
	content []byte         // nil content means content not cached. (zero length slice means empty file.)
}

func (f *File) String() string {
	return fmt.Sprintf("info %s content %#v\n", f.info, string(f.content))
}

type Reader struct {
	r *bytes.Reader
}

func (r *Reader) Read(p []byte) (int, error) {
	return r.r.Read(p)
}
func (r *Reader) Close() error {
	return nil
}

func (host *Host) Open(fpath string) (io.ReadCloser, error) {
	host.fsmu.Lock()
	file := host.fs[fpath]
	if file == nil && host.cascade != nil {
		// we don't want to hold this lock while SSH does its thing if we don't have to
		host.fsmu.Unlock()
		return host.cascade.Open(fpath)
	}
	defer host.fsmu.Unlock()

	if file == nil || file.info == nil {
		return nil, &os.PathError{Op: "open", Path: fpath, Err: syscall.ENOENT}
	}

	// TODO: Someday, be super cool and emulate a bunch of common permission errors.
	if file.info.IsDir() {
		return nil, &os.PathError{Op: "open", Path: fpath, Err: syscall.EISDIR}
	}

	reader := &Reader{
		r: bytes.NewReader(file.content),
	}
	return reader, nil
}

func (host *Host) ReadFile(fpath string) ([]byte, error) {
	buf := &bytes.Buffer{}
	fh, err := host.Open(fpath)
	if err != nil {
		return nil, err
	}
	defer fh.Close()
	_, err = io.Copy(buf, fh)
	if err != nil {
		return nil, err
	}
	err = fh.Close()
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
