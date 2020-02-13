package modules

import (
	"context"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"time"
)

type TrackService interface {
	Tracks(ctx context.Context, excludedIds []string) ([]*Track, error)
	FindTrackByIdAndCounter(ctx context.Context, track *Track) (*Track, error)
	UpdateTrackBPM(ctx context.Context, track *Track) error
}

type TrackUpdateRequest struct {
	URL              string
	ID               string            `json:"id,omitempty"`
	Counter          int64             `json:"counter,omitempty" gorm:"AUTO_INCREMENT"`
	BPM              float64           `json:"bpm,omitempty"`
	UpdatedAt        time.Time         `json:"updated_at,omitempty"`
}

// Track defines a resource that be played/streamed
type Track struct {
	ID               string            `json:"id,omitempty"`
	BPM              float64           `json:"bpm,omitempty"`
	Counter          int64             `json:"counter,omitempty" gorm:"AUTO_INCREMENT"`
	Title            string            `json:"title,omitempty"`
	AlbumID          string            `json:"album_id,omitempty"`
	SortTitle        string            `json:"sort_title,omitempty"`
	Lyrics           string            `json:"lyrics,omitempty"`
	IsExplicit       bool              `json:"is_explicit,omitempty"`
	ReleasedAt       time.Time         `json:"released_at,omitempty"`
	Published        bool              `json:"published"`
	IsDeleted        bool              `json:"is_deleted"`
	CreatedAt        time.Time         `json:"created_at,omitempty"`
	UpdatedAt        time.Time         `json:"updated_at,omitempty"`
	DeletedAt        *time.Time        `json:"deleted_at,omitempty"`
	OriginalResource *UrlResource      `json:"original_resource,omitempty" gorm:"column:original_resource"`
}

type UrlResource struct {
	URL      *url.URL             `json:"url"`
	Bucket   string               `json:"bucket"`
	Type     UploadedResourceType `json:"type"`
	FileName string               `json:"filename"`
	FilePath string               `json:"filepath"`
}

type TrackBpmResult struct {
	error []string
	success int
	failure int
}

type UploadedResourceType string

// Value get value of Jsonb
func (u UrlResource) Value() (driver.Value, error) {
	j, err := json.Marshal(u)
	return j, err
}

// Scan scan value into Hash
func (u *UrlResource) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New(fmt.Sprint("Failed to unmarshal JSONB value:", value))
	}

	return json.Unmarshal(bytes, u)
}
