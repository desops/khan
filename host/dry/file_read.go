package dry

import (
	"bytes"
	"io"
	"os"
	"syscall"
)

type File struct {
	info    *FileInfo // nil info means file not present (deleted)
	content []byte
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
	defer host.fsmu.Unlock()

	file := host.fs[fpath]

	if file == nil && host.cascade != nil {
		return host.cascade.Open(fpath)
	}

	if file == nil {
		return nil, &os.PathError{Op: "open", Path: fpath, Err: syscall.ENOENT}
	}

	if file != nil && file.info != nil {
		// TODO: Someday, be super cool and emulate a bunch of common permission errors.

		if file.info.isdir {
			return nil, &os.PathError{Op: "open", Path: fpath, Err: syscall.EISDIR}
		}
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