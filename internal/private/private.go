package private

import (
	"context"
	"io"

	"github.com/containers/image/v5/docker/reference"
	"github.com/containers/image/v5/types"
)

// ImageDestination is an internal extension to the types.ImageDestination
// interface.
type ImageDestination interface {
	types.ImageDestination

	// PutBlobWithOptions is a wrapper around PutBlob.  If
	// options.LayerIndex is set, the blob will be committed directly.
	// Either by the calling goroutine or by another goroutine already
	// committing layers.
	//
	// Please note that TryReusingBlobWithOptions and PutBlobWithOptions
	// *must* be used the together.  Mixing the two with non "WithOptions"
	// functions is not supported.
	PutBlobWithOptions(ctx context.Context, stream io.Reader, blobinfo types.BlobInfo, options PutBlobOptions) (types.BlobInfo, error)

	// TryReusingBlobWithOptions is a wrapper around TryReusingBlob.  If
	// options.LayerIndex is set, the reused blob will be recoreded as
	// already pulled.
	//
	// Please note that TryReusingBlobWithOptions and PutBlobWithOptions
	// *must* be used the together.  Mixing the two with non "WithOptions"
	// functions is not supported.
	TryReusingBlobWithOptions(ctx context.Context, blobinfo types.BlobInfo, options TryReusingBlobOptions) (bool, types.BlobInfo, error)
}

// PutBlobOptions are used in PutBlobWithOptions.
type PutBlobOptions struct {
	// Cache to look up blob infos.
	Cache types.BlobInfoCache
	// Denotes whether the blob is a config or not.
	IsConfig bool
	// Indicates an empty layer.
	EmptyLayer bool
	// The corresponding index in the layer slice.
	LayerIndex *int
}

// TryReusingBlobOptions are used in TryReusingBlobWithOptions.
type TryReusingBlobOptions struct {
	// Cache to look up blob infos.
	Cache types.BlobInfoCache
	// Use an equivalent of the desired blob.
	CanSubstitute bool
	// Indicates an empty layer.
	EmptyLayer bool
	// The corresponding index in the layer slice.
	LayerIndex *int
	// The reference of the image that contains the target blob.
	SrcRef reference.Named
}

// ImageSourceChunk is a portion of a blob.
// This API is experimental and can be changed without bumping the major version number.
type ImageSourceChunk struct {
	Offset uint64
	Length uint64
}

// ImageSourceSeekable is an image source that permits to fetch chunks of the entire blob.
// This API is experimental and can be changed without bumping the major version number.
type ImageSourceSeekable interface {
	// GetBlobAt returns a stream for the specified blob.
	// The specified chunks must be not overlapping and sorted by their offset.
	GetBlobAt(context.Context, types.BlobInfo, []ImageSourceChunk) (chan io.ReadCloser, chan error, error)
}

// ImageDestinationPartial is a service to store a blob by requesting the missing chunks to a ImageSourceSeekable.
// This API is experimental and can be changed without bumping the major version number.
type ImageDestinationPartial interface {
	// PutBlobPartial writes contents of stream and returns data representing the result.
	PutBlobPartial(ctx context.Context, stream ImageSourceSeekable, srcInfo types.BlobInfo, cache types.BlobInfoCache) (types.BlobInfo, error)
}

// BadPartialRequestError is returned by ImageSourceSeekable.GetBlobAt on an invalid request.
type BadPartialRequestError struct {
	Status string
}

func (e BadPartialRequestError) Error() string {
	return e.Status
}
