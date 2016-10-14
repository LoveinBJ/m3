	"github.com/m3db/m3db/ts"
	Open(namespace ts.ID, shard uint32, start time.Time) error
	// Write will write the id and data pair and returns an error on a write error
	Write(id ts.ID, data []byte) error
	// WriteAll will write the id and all byte slices and returns an error on a write error
	WriteAll(id ts.ID, data [][]byte) error
	Open(namespace ts.ID, shard uint32, start time.Time) error
	// Read returns the next id and data pair or error, will return io.EOF at end of volume
	Read() (id ts.ID, data []byte, err error)
	// SetInstrumentOptions sets the instrumentation options
	SetInstrumentOptions(value instrument.Options) Options
	// InstrumentOptions returns the instrumentation options
	InstrumentOptions() instrument.Options
	// SetRetentionOptions sets the retention options
	SetRetentionOptions(value retention.Options) Options
	// RetentionOptions returns the retention options
	RetentionOptions() retention.Options
	// SetFilePathPrefix sets the file path prefix for sharded TSDB files
	SetFilePathPrefix(value string) Options
	// FilePathPrefix returns the file path prefix for sharded TSDB files
	FilePathPrefix() string
	// SetNewFileMode sets the new file mode
	SetNewFileMode(value os.FileMode) Options
	// NewFileMode returns the new file mode
	NewFileMode() os.FileMode
	// SetNewDirectoryMode sets the new directory mode
	SetNewDirectoryMode(value os.FileMode) Options
	// NewDirectoryMode returns the new directory mode
	NewDirectoryMode() os.FileMode
	// SetWriterBufferSize sets the buffer size for writing TSDB files
	SetWriterBufferSize(value int) Options
	// WriterBufferSize returns the buffer size for writing TSDB files
	WriterBufferSize() int
	// SetReaderBufferSize sets the buffer size for reading TSDB files
	SetReaderBufferSize(value int) Options
	// ReaderBufferSize returns the buffer size for reading TSDB files
	ReaderBufferSize() int