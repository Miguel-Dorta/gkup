package hasher

import (
	"github.com/Miguel-Dorta/gkup/pkg"
	"github.com/Miguel-Dorta/gkup/pkg/files"
	"github.com/Miguel-Dorta/gkup/pkg/threadSafe"
	"golang.org/x/sync/errgroup"
	"sync"
)

// MultiHasher is a type for making hashing operations concurrently
type MultiHasher struct {
	workers []*Hasher
}

// NewMultiHasher creates a new MultiHasher object
func NewMultiHasher(algorithm string) (*MultiHasher, error) {
	var err error

	workers := make([]*Hasher, pkg.NumberOfThreads)
	for i := range workers {
		workers[i], err = New(algorithm)
		if err != nil {
			return nil, err
		}
	}
	return &MultiHasher{workers: workers}, nil
}

// CheckFiles checks concurrently whether the files listed in the path slice provided match with the info contained in their names.
// That means that they follow the specification from files.GetFileFromName() and their information is correct.
// This process is aimed to detect file corruption or filename defects.
func (mh *MultiHasher) CheckFiles(paths []string) bool {
	var wg sync.WaitGroup
	var errsFound fuse
	pathsSafe := threadSafe.NewStringList(paths)

	for _, w := range mh.workers {
		wg.Add(1)
		go w.fileChecker(pathsSafe, &errsFound, &wg)
	}

	wg.Wait()
	return bool(errsFound)
}

// GetFiles creates concurrently a list of files.File from the path list provided
func (mh *MultiHasher) GetFiles(paths []string) ([]*files.File, error) {
	var eg errgroup.Group
	pathsSafe := threadSafe.NewStringList(paths)
	filesSafe := threadSafe.NewFileList(make([]*files.File, 0, len(paths)))

	for _, w := range mh.workers {
		eg.Go(func() error {
			return w.fileGetter(pathsSafe, filesSafe)
		})
	}

	if err := eg.Wait(); err != nil {
		return nil, err
	}
	return filesSafe.GetList(), nil
}

// HashFiles gets concurrently the hash from the list of files.File provided
func (mh *MultiHasher) HashFiles(files []*files.File) error {
	var eg errgroup.Group
	filesSafe := threadSafe.NewFileList(files)

	for _, w := range mh.workers {
		eg.Go(func() error {
			return w.fileHasher(filesSafe)
		})
	}

	if err := eg.Wait(); err != nil {
		return err
	}
	return nil
}
