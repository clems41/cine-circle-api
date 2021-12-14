package migration

import (
	"gorm.io/gorm"
)

type Migration interface {
	JobName() string            // Name your migration job
	Job(DB *gorm.DB) error      // Here is where you put your migration operations, consider using a transaction
	IsJobDone(DB *gorm.DB) bool // Write a method that checks if your job worked
	PreventJobRerun() bool      // Should we test that IsJobDone() sends false before running your job?
}

/*
	88  88b           d88  88888888ba     ,ad8888ba,    88888888ba  888888888888  888888888888    db         888b      88  888888888888
	88  888b         d888  88      "8b   d8"'    `"8b   88      "8b      88            88        d88b        8888b     88       88
	88  88`8b       d8'88  88      ,8P  d8'        `8b  88      ,8P      88            88       d8'`8b       88 `8b    88       88
	88  88 `8b     d8' 88  888888888P'  88          88  888888888P'      88            88      d8'  `8b      88  `8b   88       88
	88  88  `8b   d8'  88  88""""""'    88          88  88""""88'        88            88     d8Y8888Y8b     88   `8b  88       88
	88  88   `8b d8'   88  88           Y8,        ,8P  88    `8b        88            88    d8""""""""8b    88    `8b 88       88
	88  88    `888'    88  88            Y8a.    .a8P   88     `8b       88            88   d8'        `8b   88     `8888       88
	88  88     `8'     88  88             `"Y8888Y"'    88      `8b      88            88  d8'          `8b  88      `888       88


	Here is where you defined new action to do AT THE END of the 'migrations' list
	Modifying existing jobs should be punished by castration, create a new job to update your data like a good boy
	Each job is responsible of managing (or not) a transaction. Here is a potato for the long explanation

	The job should NEVER use packages from the "internal" directory, because it would always use the last
	version, which is not always compatible with old jobs. That is why you should ALWAYS write your migrations jobs
	without importing anything from the project itself : Each job should be independent. For instance :
	it should define its own structs for gorm, or using plain SQL etc.

*/

func GetMigrationJobs() []Migration {
	var migrations []Migration
	//migrations = append(migrations, jobs.ExampleJob{})
	// ^ INSERT NEW JOBS HERE, ALWAYS AT THE BOTTOM, AND BE CAREFUL WHILE MERGING
	return migrations
}
