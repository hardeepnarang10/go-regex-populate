package goregexpopulate

type Populate interface {
	Populate(any) error
}
