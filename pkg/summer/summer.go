package summer

import (
	"context"
	"errors"
	"fmt"

	"github.com/utilyre/summer/pkg/pipeline"
	"golang.org/x/sync/errgroup"
)

var (
	ErrInvalidText        = errors.New("invalid text")
	ErrNonPositiveInteger = errors.New("non-positive integer")
)

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
	algo          Algorithm
	readWorkers   int
	digestWorkers int
}

// WithAlgorithm sets the used hash function for SumTree.
func WithAlgorithm(algo Algorithm) Option {
	return func(o *options) error {
		o.algo = algo
		return nil
	}
}

// WithReadWorkers sets number of read workers for SumTree.
func WithReadWorkers(workers int) Option {
	return func(o *options) error {
		if workers <= 0 {
			return OptionError{
				Which: "read workers",
				Value: workers,
				Err:   ErrNonPositiveInteger,
			}
		}

		o.readWorkers = workers
		return nil
	}
}

// WithDigestWorkers sets number of digest workers for SumTree.
func WithDigestWorkers(workers int) Option {
	return func(o *options) error {
		if workers <= 0 {
			return OptionError{
				Which: "digest workers",
				Value: workers,
				Err:   ErrNonPositiveInteger,
			}
		}

		o.digestWorkers = workers
		return nil
	}
}

// SumTree recursively generates checksums for each file under all roots in
// parallel.
func SumTree(
	ctx context.Context,
	roots []string,
	opts ...Option,
) ([]Checksum, error) {
	o := &options{
		algo:          AlgorithmMD5,
		readWorkers:   1,
		digestWorkers: 1,
	}
	for _, opt := range opts {
		if err := opt(o); err != nil {
			return nil, err
		}
	}

	g, ctx := errgroup.WithContext(ctx)

	var pl pipeline.Pipeline
	pl.Append(o.readWorkers, readPipe{g})
	pl.Append(o.digestWorkers, digestPipe{g, o.algo})
	out := pl.Pipe(ctx, walkPipe{g, roots}.Pipe(ctx, nil))

	var checksums []Checksum
	for v := range out {
		cs := v.(Checksum)
		checksums = append(checksums, cs)
	}

	if err := g.Wait(); err != nil {
		return nil, err
	}

	return checksums, nil
}
