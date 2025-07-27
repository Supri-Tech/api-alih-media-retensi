
package models

import "time"

type Pemusnahan struct {
	ID         int
	Status     string
	TglLaporan time.Time
	CreatedAt  time.Time
	UpdatedAt  time.Time
}
