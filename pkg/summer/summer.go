package summer

import (
	"context"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"iter"
	"os"

	"github.com/utilyre/summer/pkg/pipeline"
)

var (
	ErrNilValue           = errors.New("nil value")
	ErrInvalidText        = errors.New("invalid text")
	ErrNonPositiveInteger = errors.New("non-positive integer")
)

type Summer struct {
	options
}

func New(opts ...Option) (*Summer, error) {
	o := options{
		fsys:       os.DirFS("."),
		algo:       AlgorithmMD5,
		readJobs:   1,
		digestJobs: 1,
		recursive:  false,
	}
	for _, opt := range opts {
		if err := opt(&o); err != nil {
			return nil, err
		}
	}

	return &Summer{options: o}, nil
}

func (s *Summer) Sum(ctx context.Context, names []string) (iter.Seq[Checksum], error) {
	var pl pipeline.Pipeline[Checksum]
	pl.Append(s.readJobs, readPipe{s.fsys})
	pl.Append(s.digestJobs, digestPipe{s.algo})

	var namesCh <-chan Checksum
	if s.recursive {
		namesCh = walkDirs(ctx, s.fsys, names)
	} else {
		namesCh = walkFiles(ctx, s.fsys, names)
	}

	out := pl.Pipe(ctx, namesCh)

	return func(yield func(Checksum) bool) {
		for cs := range out {
			if !yield(cs) {
				return
			}
		}
	}, nil
}

// An Algorithm represents a supported hash function.
type Algorithm int

const (
	AlgorithmMD5 Algorithm = iota + 1
	AlgorithmSha1
	AlgorithmSha256
	AlgorithmSha512
)

func (algo Algorithm) String() string {
	switch algo {
	case AlgorithmMD5:
		return "md5"
	case AlgorithmSha1:
		return "sha1"
	case AlgorithmSha256:
		return "sha256"
	case AlgorithmSha512:
		return "sha512"
	default:
		return ""
	}
}

func (Algorithm) Type() string {
	return "md5|sha1|sha256|sha512"
}

func (algo *Algorithm) Set(value string) error {
	return algo.UnmarshalText([]byte(value))
}

func (algo Algorithm) MarshalText() ([]byte, error) {
	return []byte(algo.String()), nil
}

func (algo *Algorithm) UnmarshalText(text []byte) error {
	s := string(text)
	switch s {
	case "md5":
		*algo = AlgorithmMD5
		return nil
	case "sha1":
		*algo = AlgorithmSha1
		return nil
	case "sha256":
		*algo = AlgorithmSha256
		return nil
	case "sha512":
		*algo = AlgorithmSha512
		return nil
	default:
		return fmt.Errorf("algorithm: string '%s': %w", s, ErrInvalidText)
	}
}

// A Checksum represents name and hash of a particular file.
type Checksum struct {
	Name string
	Hash []byte
	Err  error

	body io.ReadCloser
}

// An Option represents an optional parameter given to SumTree.
type Option func(o *options) error

// An OptionError represents a validation error of Option.
type OptionError struct {
	Which string
	Value any
	Err   error
}

func (e OptionError) Error() string {
	return fmt.Sprintf("option '%s': value %v: %v", e.Which, e.Value, e.Err)
}

func (e OptionError) Unwrap() error {
	return e.Err
}

type options struct {
	fsys       fs.FS
	algo       Algorithm
	readJobs   int
	digestJobs int
	recursive  bool
}

// WithFS determines what filesystem to walk.
func WithFS(fsys fs.FS) Option {
	return func(o *options) error {
		if fsys == nil {
			return OptionError{
				Which: "fs",
				Value: fsys,
				Err:   ErrNilValue,
			}
		}

		o.fsys = fsys
		return nil
	}
}

// WithAlgorithm determines what hash function to use.
func WithAlgorithm(algo Algorithm) Option {
	return func(o *options) error {
		o.algo = algo
		return nil
	}
}

// WithReadJobs determines how many jobs to spin up for reading.
func WithReadJobs(n int) Option {
	return func(o *options) error {
		if n <= 0 {
			return OptionError{
				Which: "read jobs",
				Value: n,
				Err:   ErrNonPositiveInteger,
			}
		}

		o.readJobs = n
		return nil
	}
}

// WithReadJobs determines how many jobs to spin up for digesting.
func WithDigestJobs(n int) Option {
	return func(o *options) error {
		if n <= 0 {
			return OptionError{
				Which: "digest jobs",
				Value: n,
				Err:   ErrNonPositiveInteger,
			}
		}

		o.digestJobs = n
		return nil
	}
}

// WithRecursive determines whether to walk recursively.
func WithRecursive(v bool) Option {
	return func(o *options) error {
		o.recursive = v
		return nil
	}
}
