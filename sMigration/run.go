package sMigration

import (
	"context"
	"fmt"
	"io/fs"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"reflect"
	"regexp"
	"runtime"
	"sort"
	"strings"
	"time"

	"gitlab.com/soteapps/packages/v2022/db/migrations"
	"gitlab.com/soteapps/packages/v2022/db/seeds"
	"gitlab.com/soteapps/packages/v2022/sCustom"
	"gitlab.com/soteapps/packages/v2022/sDatabase"
	"gitlab.com/soteapps/packages/v2022/sError"
	"gitlab.com/soteapps/packages/v2022/sLogger"
	"golang.org/x/exp/slices"
)

const (
	MigrationTableName             = "sote_db_log"
	MigrationType                  = "migration"
	SeedingType                    = "seeding"
	SeedingAction                  = "seed"
	goFileType                     = ".go"
	sqlFileType                    = ".sql"
	SeedsSubDir                    = "/db/seeds"
	MigrationsSubDir               = "/db/migrations"
	SeedsPackageName               = "seeds"
	MigrationsPackageName          = "migrations"
	ExternalDefaultStackTraceSkips = 3 // when you call functions/methods from outside this file
	internalDefaultStackTraceSkips = 2 // when you call functions/methods from within this file
)

type Config struct{ DBConnInfo sDatabase.ConnInfo }

type MigrationFiles struct {
	FilePath         string
	FileName         string
	MigrationName    string
	MigrationVersion string
	FileType         string
}

type MigrationQueryInfo struct {
	Params                  []interface{}
	VersionPreparedSubQuery string
	SetupType               string // sMigration.MigrationType/sMigration.SeedingType
	MigrationAction         string // sMigration.MigrationType/sMigration.SeedingAction
	StartMigrationMsg       string // seed
	EndMigrationMsg         string // seed
}

func init() {
	sLogger.SetLogMessagePrefix(LOGMESSAGEPREFIX)
}

// New sets up the migration and seeding info .
// setupType either sMigration.SeedingType or sMigration.MigrationType
func New(ctx context.Context, environment string, setupType string, stackSkips int) (mDir string, dbConnInfo sDatabase.ConnInfo,
	soteErr sError.SoteError) {
	sLogger.DebugMethod()

	var (
		mSubDir string
		err     error
		table   = sDatabase.TableInfo{
			Name: MigrationTableName,
			PrimaryKey: &sDatabase.PrimaryKeyInfo{
				Columns: []string{"version"},
			},
			Description: "Contains Migration and Seeding Logs for sDatabase package",
		}

		columns = []sDatabase.ColumnInfo{
			{
				Name:        "version",
				DataType:    sDatabase.BIGINTEGER,
				Description: "This is a 14 integer-long unique identifier for migration or seeding file",
			},
			{
				Name:        "migration_name",
				DataType:    sDatabase.STRING,
				Description: "This is the migration or seeding file name without the version",
			},
			{
				Name:        "migration_action",
				DataType:    sDatabase.STRING,
				Description: "This states which action was used either seed or migration",
			},
			{
				Name:        "migration_time",
				DataType:    sDatabase.TIMESTAMP,
				Default:     sDatabase.CURRENTTIMESTAMP,
				Description: "This is the time the migration/seeding took place",
			},
		}
	)

	switch setupType {
	case SeedingType:
		mSubDir = SeedsSubDir
	case MigrationType:
		mSubDir = MigrationsSubDir
	default:
		soteErr = sError.GetSError(sError.ErrGenericError, sError.BuildParams([]string{"invalid setup type"}), sError.EmptyMap)
		return
	}

	if dbConnInfo, soteErr = sDatabase.New(ctx, environment); soteErr.ErrCode != nil {
		return
	}

	// set up the migration table
	if soteErr = dbConnInfo.CreateTableIfNotExists(ctx, table, columns, setupType); soteErr.ErrCode != nil {
		return
	}

	// set up the migration directory
	_, file, _, _ := runtime.Caller(stackSkips)
	sLogger.Info(fmt.Sprintf("Caller File %v", file))
	basePath := filepath.Dir(path.Join(path.Dir(file)))
	mDir = basePath + mSubDir
	if err = os.MkdirAll(basePath+MigrationsSubDir, os.ModePerm); err == nil { // migration dir
		_ = os.WriteFile(basePath+MigrationsSubDir+"/init.go",
			[]byte(fmt.Sprintf("package %v\n", MigrationsPackageName)+
				"import \"gitlab.com/soteapps/packages/v2022/sDatabase\"\n"+
				"type Config struct{DBConnInfo sDatabase.ConnInfo}"),
			os.ModePerm)
		_ = exec.Command("go", "fmt", basePath+MigrationsSubDir+"/init.go").Run()
	}
	if err = os.MkdirAll(basePath+SeedsSubDir, os.ModePerm); err == nil { // seeding dir
		_ = os.WriteFile(basePath+SeedsSubDir+"/init.go",
			[]byte(fmt.Sprintf("package %v\n", SeedsPackageName)+
				"import \"gitlab.com/soteapps/packages/v2022/sDatabase\"\n"+
				"type Config struct{DBConnInfo sDatabase.ConnInfo}"), os.ModePerm)
		_ = exec.Command("go", "fmt", basePath+SeedsSubDir+"/init.go").Run()
	}

	return
}

// Run  runs either MigrationType| SeedingAction based on the action run|setup.
// This is purposely setup to be used by command line e.g. go run main.go -e development -s MigrationType -a setup
func Run(ctx context.Context, environment string, service string, action string) (soteErr sError.SoteError) {
	sLogger.DebugMethod()

	var (
		mDir     string
		isAction bool
	)
	switch service {
	case MigrationType:
		if action == "setup" {
			if mDir, _, soteErr = New(ctx, environment, MigrationType, internalDefaultStackTraceSkips); soteErr.ErrCode == nil {
				sLogger.Info(mDir)
			}

		} else if action == "run" {
			soteErr = Migrate(ctx, environment)
		} else {
			isAction = true
		}

	case SeedingAction:
		if action == "setup" {
			if mDir, _, soteErr = New(ctx, environment, SeedingType, internalDefaultStackTraceSkips); soteErr.ErrCode == nil {
				sLogger.Info(mDir)
			}

		} else if action == "run" {
			soteErr = Seed(ctx, environment)
		} else {
			isAction = true
		}
	default:
		soteErr = sError.GetSError(sError.ErrInvalidParameterValue,
			sError.BuildParams([]string{"service", service, fmt.Sprintf("%v|%v", MigrationType, SeedingAction)}),
			sError.EmptyMap)
	}

	if isAction {
		soteErr = sError.GetSError(sError.ErrInvalidParameterValue, sError.BuildParams([]string{"action", action, "setup|run"}), sError.EmptyMap)
	}

	return
}

//  By default, this function migrates|seeds all .go & .sql files withing MigrationsSubDir | SeedsSubDir folder
// setupType either sMigration.SeedingType or sMigration.MigrationType
func migrationAndSeeding(ctx context.Context, environment string, setupType string, stackSkips int) (soteErr sError.SoteError) {
	sLogger.DebugMethod()

	var (
		migrationDir   string
		migrationFiles = make([]MigrationFiles, 0)
		migrationQInfo = &MigrationQueryInfo{}
		dbConnInfo     = sDatabase.ConnInfo{}
	)

	if migrationDir, dbConnInfo, soteErr = New(ctx, environment, setupType, stackSkips); soteErr.ErrCode != nil {
		return
	}

	config := Config{DBConnInfo: dbConnInfo}
	if migrationFiles, migrationQInfo, soteErr = config.getMigrationAndSeedsFiles(migrationDir, setupType); soteErr.ErrCode != nil {
		return
	}

	soteErr = config.run(ctx, migrationFiles, migrationQInfo)

	return
}

// return migrations and seeds files from MigrationsSubDir | SeedsSubDir directory based on the setupType.
// they should be in format 14 digits, underscore,alphanumeric .go e.g. 20221122000001_test12.go.
// The part after the first underscore, represents the function name. The example will be func Test12(ctx context.Context, dbConnInfo sDatabase.ConnInfo)
func (config Config) getMigrationAndSeedsFiles(migrationDir string, setupType string) (migrationFiles []MigrationFiles,
	migrationQInfo *MigrationQueryInfo, soteErr sError.SoteError) {
	sLogger.DebugMethod()

	var (
		files     []fs.FileInfo
		err       error
		vParamStr string
		mVersions = make([]string, 0)
		mNames    = make([]string, 0)
	)
	if files, err = ioutil.ReadDir(migrationDir); err != nil {
		soteErr = sError.GetSError(sError.ErrGenericError, sError.BuildParams([]string{err.Error()}), sError.EmptyMap)
		return
	}

	migrationFiles = make([]MigrationFiles, 0)
	migrationQInfo = &MigrationQueryInfo{
		Params: make([]interface{}, 0),
	}

	i := 0
	for _, file := range files {
		re := regexp.MustCompile(`(?i)^([0-9]{14})(_*[A-Z]+[A-Z0-9_]*)(\.go|\.sql)$`)
		match := re.FindStringSubmatch(file.Name())
		if len(match) == 4 {
			mFile := MigrationFiles{
				FileName:         file.Name(),
				FilePath:         migrationDir + "/" + match[0],
				MigrationVersion: match[1],
				FileType:         strings.ToLower(match[3]),
			}

			re, _ = regexp.Compile(`_\w`)
			mFile.MigrationName = re.ReplaceAllStringFunc(regexp.MustCompile(`_+`).ReplaceAllString(match[2], "_"), func(m string) string {
				return strings.ToUpper(m[1:])
			})
			if slices.Contains(mVersions, mFile.MigrationVersion) {
				soteErr = sError.GetSError(sError.ErrDuplicateItems, sError.BuildParams([]string{fmt.Sprintf("Migration version %v in migration %v",
					mFile.MigrationVersion, mFile.FilePath)}), sError.EmptyMap)
				return
			}

			if slices.Contains(mNames, mFile.MigrationName) {
				soteErr = sError.GetSError(sError.ErrDuplicateItems, sError.BuildParams([]string{fmt.Sprintf("Migration name %v in migration %v",
					mFile.MigrationName, mFile.FilePath)}), sError.EmptyMap)
				return
			}

			switch setupType {
			case SeedingType:
				if mFile.FileType == goFileType { // check if the function for the go file exists
					soteErr = sCustom.UserFuncExists(mFile.MigrationName, seeds.Config{DBConnInfo: config.DBConnInfo})
				}
				migrationQInfo.EndMigrationMsg = "seeded"
				migrationQInfo.StartMigrationMsg = "seeding"
				migrationQInfo.SetupType = SeedingType
				migrationQInfo.MigrationAction = SeedingAction
			case MigrationType:
				if mFile.FileType == goFileType { // check if the function for the go file exists
					soteErr = sCustom.UserFuncExists(mFile.MigrationName, migrations.Config{DBConnInfo: config.DBConnInfo})
				}
				migrationQInfo.EndMigrationMsg = "migrated"
				migrationQInfo.StartMigrationMsg = "migrating"
				migrationQInfo.SetupType = MigrationType
				migrationQInfo.MigrationAction = MigrationType
			default:
				soteErr = sError.GetSError(sError.ErrGenericError, sError.BuildParams([]string{"invalid setup type"}), sError.EmptyMap)
				return
			}

			if soteErr.ErrCode != nil {
				soteErr.FmtErrMsg += fmt.Sprintf(" for migration %v", mFile.FilePath)
				log.Fatalln(soteErr.FmtErrMsg)
			}

			mVersions = append(mVersions, mFile.MigrationVersion)
			mNames = append(mNames, mFile.MigrationName)
			migrationFiles = append(migrationFiles, mFile)
			migrationQInfo.Params = append(migrationQInfo.Params, mFile.MigrationVersion)
			vParamStr += fmt.Sprintf("$%v,", i+1)
			i++
		}
	}

	migrationQInfo.VersionPreparedSubQuery = strings.TrimSuffix(vParamStr, ",")

	return
}

// runs all the migrations and seeds files
func (config Config) run(ctx context.Context, migrationFiles []MigrationFiles,
	migrationQInfo *MigrationQueryInfo) (soteErr sError.SoteError) {
	sLogger.DebugMethod()

	mStart := time.Now()
	if len(migrationFiles) > 0 {
		sort.Slice(migrationFiles, func(i, j int) bool {
			return migrationFiles[i].MigrationVersion < migrationFiles[j].MigrationVersion
		}) // sort the files in ascending order
		var (
			tRows              sDatabase.SRows
			err                error
			existingMigrations = make([]MigrationFiles, 0)
			resp               []reflect.Value
			ok                 bool
			fContent           []byte
		)

		qStmt := fmt.Sprintf("SELECT version::TEXT AS version,migration_name FROM %v WHERE version IN (%v)", MigrationTableName,
			migrationQInfo.VersionPreparedSubQuery)
		tRows, soteErr = config.DBConnInfo.QueryDBStmt(ctx, qStmt, migrationQInfo.SetupType, migrationQInfo.Params...)
		if soteErr.ErrCode != nil {
			return
		}

		defer tRows.Close()
		for tRows.Next() {
			existingMigration := new(MigrationFiles)
			if err = tRows.Scan(
				&existingMigration.MigrationVersion,
				&existingMigration.MigrationName,
			); err != nil {
				soteErr = sError.GetSError(sError.ErrGenericError, sError.BuildParams([]string{err.Error()}), sError.EmptyMap)
				return
			}

			existingMigrations = append(existingMigrations, *existingMigration)
		}

		for _, file := range migrationFiles {
			if len(existingMigrations) > 0 {
				if idSkip := slices.IndexFunc(existingMigrations, func(m MigrationFiles) bool {
					// checks if the exact migration file was already migrated in a previous migration
					return m.MigrationVersion == file.MigrationVersion && m.MigrationName == file.MigrationName
				}); idSkip != -1 {
					continue // skip this migration
				}

				if idDup := slices.IndexFunc(existingMigrations, func(m MigrationFiles) bool {
					// checks if we have a different migration with similar version was already migrated in a previous migration
					return m.MigrationVersion == file.MigrationVersion
				}); idDup != -1 {
					sLogger.Info(fmt.Sprintf("Skipping %v - a migration with a similar version(%v) was previously migrated",
						file.FilePath, file.MigrationVersion))
					continue // skip this migration
				}
			}

			// migrate the file
			sLogger.Info(fmt.Sprintf("== %v %v: %v", file.MigrationVersion, file.MigrationName, migrationQInfo.StartMigrationMsg))
			start := time.Now()
			if file.FileType == goFileType {
				switch migrationQInfo.SetupType {
				case SeedingType:
					resp, soteErr = sCustom.CallUserFunc(file.MigrationName, seeds.Config{DBConnInfo: config.DBConnInfo},
						ctx)
				case MigrationType:
					resp, soteErr = sCustom.CallUserFunc(file.MigrationName, migrations.Config{DBConnInfo: config.DBConnInfo},
						ctx)
				default:
					soteErr = sError.GetSError(sError.ErrGenericError, sError.BuildParams([]string{"invalid setup type"}), sError.EmptyMap)
					return
				}

				if soteErr.ErrCode != nil {
					log.Fatalln(soteErr.FmtErrMsg)
				}

				for _, value := range resp {
					if soteErr, ok = value.Interface().(sError.SoteError); ok && soteErr.ErrCode != nil {
						log.Fatalln(soteErr.FmtErrMsg)
					}
				}
			} else if file.FileType == sqlFileType {
				if fContent, err = ioutil.ReadFile(file.FilePath); err != nil {
					soteErr = sError.GetSError(sError.ErrGenericError, sError.BuildParams([]string{err.Error()}), sError.EmptyMap)
					log.Fatalln(soteErr.FmtErrMsg)
				}

				qStmts := strings.Split(string(fContent), ";")
				for _, qStmt = range qStmts {
					if soteErr = config.DBConnInfo.ExecDBStmt(ctx, qStmt, migrationQInfo.SetupType); soteErr.ErrCode != nil {
						log.Fatalln(soteErr.FmtErrMsg)
					}
				}
			} else {
				continue // skip the file
			}

			config.DBConnInfo.ExecDBStmt(ctx,
				fmt.Sprintf("INSERT INTO %v (version, migration_name, migration_action) VALUES ($1,$2,$3)", MigrationTableName),
				migrationQInfo.SetupType, file.MigrationVersion, file.MigrationName, migrationQInfo.SetupType)
			sLogger.Info(fmt.Sprintf("== %v %v: %v %v", file.MigrationVersion, file.MigrationName, migrationQInfo.MigrationAction, time.Since(start)))
		}
	}

	sLogger.Info(fmt.Sprintf("All %v Done. Took %v", migrationQInfo.SetupType, time.Since(mStart)))

	// print done message

	return
}
