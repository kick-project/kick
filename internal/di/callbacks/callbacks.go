package callbacks

import "github.com/kick-project/kick/internal/resources/client/plumb"

// MakePlumb callback injector
type MakePlumb func(url, ref string) (*plumb.Plumb, error)
