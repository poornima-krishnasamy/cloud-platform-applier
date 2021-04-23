package sysutil // import "github.com/cloud-platform-applier/sysutil"

Package sysutil provides utility functions needed for the
cloud-platform-applier

FUNCTIONS

func GetEnvStringOrDefault(key, def string) string
func GetRequiredEnvString(key string) string

TYPES

type Clock struct{}
    Clock implements ClockInterface with the standard time library functions.

func (c *Clock) Now() time.Time
    Now returns current time

func (c *Clock) Since(t time.Time) time.Duration
    Since returns time since t

func (c *Clock) Sleep(d time.Duration)
    Sleep sleeps for d duration

type ClockInterface interface {
	Now() time.Time
	Since(time.Time) time.Duration
	Sleep(time.Duration)
}
    ClockInterface allows for mocking out the functionality of the standard time
    library when testing.

type FileSystem struct{}

func (f *FileSystem) ListFolders(path string) ([]string, error)
    ListFolders take the path as input, list all the folders in the give path
    and return a array of strings containing the list of folders