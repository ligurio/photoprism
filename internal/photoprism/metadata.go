package photoprism

import (
	"path/filepath"

	"github.com/photoprism/photoprism/internal/meta"
	"github.com/photoprism/photoprism/pkg/fs"
	"github.com/photoprism/photoprism/pkg/txt"
)

// MetaData returns exif meta data of a media file.
func (m *MediaFile) MetaData() (result meta.Data, err error) {
	m.metaDataOnce.Do(func() {
		err = m.metaData.Exif(m.FileName())

		if jsonFile := fs.TypeJson.FindSub(m.FileName(), fs.HiddenPath, false); jsonFile == "" {
			log.Debugf("mediafile: no json sidecar file found for %s", txt.Quote(filepath.Base(m.FileName())))
		} else if jsonErr := m.metaData.JSON(jsonFile); jsonErr != nil {
			log.Warn(jsonErr)
		} else {
			err = nil
		}

		if err != nil {
			log.Debugf("mediafile: %s", err.Error())
		}
	})

	return m.metaData, err
}
