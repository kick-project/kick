package clauses

import (
	"gorm.io/gorm/clause"
)

// OrIgnore - INSERT [OR IGNORE] clause
var OrIgnore = clause.Insert{Modifier: "OR IGNORE"}

// OrReplace - INSERT [OR REPLACE] clause
var OrReplace = clause.Insert{Modifier: "OR REPLACE"}
