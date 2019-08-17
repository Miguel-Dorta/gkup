package hasher

// fuse will have a race condition, but it doesn't matter :)
type fuse bool

func (f *fuse) trigger() {
	if !*f {
		*f = true
	}
}
