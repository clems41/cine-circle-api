package migration

import (
	"cine-circle-api/internal/constant/settingsConst"
	"cine-circle-api/pkg/test/setupTestCase"
	"cine-circle-api/pkg/utils/pathUtils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/tools/go/packages"
	"io/ioutil"
	"reflect"
	"strings"
	"testing"
)

var (
	// PackagesThatMustBeNotImported define package list that should not be imported by migration jobs.
	PackagesThatMustBeNotImported = []string{
		"internal/repository/model",
		"internal/domain",
	}
)

// TestMigrationJobs check that migrations jobs are still working from the beginning with a new clean database
func TestMigrationJobs(t *testing.T) {
	DB, clean := setupTestCase.OpenCleanDatabaseFromTemplate(t)
	defer clean()

	// Run the main process
	currentVersion, upToDateVersion, err := RunMigrations(DB)
	require.NoError(t, err)
	require.Equal(t, upToDateVersion, currentVersion)
}

// TestNewJobsDoNotUseTheInternalLibrary will check that any internal packages are imported by migration jobs
func TestNewJobsDoNotUseTheInternalLibrary(t *testing.T) {
	jobList := GetMigrationJobs()

	for _, job := range jobList {
		assert.NotEqual(t, job, 2)

		ty := reflect.TypeOf(job)

		pkgPath := ty.PkgPath()

		config := packages.Config{
			Mode: packages.NeedImports | packages.NeedDeps,
		}
		pcks, err := packages.Load(&config, pkgPath)
		assert.NoError(t, err, "Package was not found")
		assert.Len(t, pcks, 1)
		pck := pcks[0]

		for subPkgPath, subPkg := range pck.Imports {
			for _, packageToAvoid := range PackagesThatMustBeNotImported {
				assert.False(t, strings.Contains(subPkgPath, packageToAvoid),
					"The job %s imports the internal package %s, it should not in order to\n"+
						"ensure the job is runnable when the internal package changes.\n"+
						"Here is the import stack :\n%s", job.JobName(), packageToAvoid, subPkg)
			}
		}
	}
}

// TestGetMigrationJobs check that all new migration file created in jobs directory are added into list from migration directory. It is useful to avoid missing jobs
func TestGetMigrationJobs(t *testing.T) {
	jobList := GetMigrationJobs()
	nbJobs := len(jobList)

	// Get number of files from jobs directory
	rootPath, err := pathUtils.GetRootProjectPath()
	require.NoError(t, err)
	jobPath := rootPath + settingsConst.RelativeJobsPathFromRootDir
	files, _ := ioutil.ReadDir(jobPath)
	nbFiles := len(files)

	require.Equal(t, nbFiles, nbJobs, "Number of jobs defined in migration list is not the same that files in jobs directory. "+
		"You probably miss to ass your new job."+
		"NbJobs : %d \t NbFiles : %d", nbJobs, nbFiles)
}
