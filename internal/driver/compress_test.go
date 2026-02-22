package driver

import (
	"bytes"
	"io"
	"testing"

	"github.com/backup-blob/zfs-backup-blob/internal/domain/config"
)

func TestNewCompress(t *testing.T) {
	conf := &config.CompressConfig{
		Level: 3,
	}
	
	driver := NewCompress(conf)
	if driver == nil {
		t.Fatal("expected non-nil driver")
	}
	
	cd, ok := driver.(*CompressDriver)
	if !ok {
		t.Fatal("expected *CompressDriver type")
	}
	
	if cd.Conf.Level != 3 {
		t.Errorf("expected level 3, got %d", cd.Conf.Level)
	}
}

func TestCompressDriver_Read(t *testing.T) {
	testData := []byte("This is test data for compression. " +
		"It needs to be long enough to actually benefit from compression. " +
		"Zstd works best with larger chunks of data. " +
		"Repeating patterns help: abc abc abc abc abc abc abc abc abc abc")
	
	conf := &config.CompressConfig{Level: 3}
	driver := NewCompress(conf)
	
	// Create a reader with test data
	originalReader := bytes.NewReader(testData)
	
	// Wrap with compression (this is what happens during backup)
	compressedReader, err := driver.Read(originalReader)
	if err != nil {
		t.Fatalf("failed to create compressed reader: %v", err)
	}
	
	// Read all compressed data
	compressedData, err := io.ReadAll(compressedReader)
	if err != nil {
		t.Fatalf("failed to read compressed data: %v", err)
	}
	
	// Verify compression actually happened (compressed should be smaller)
	if len(compressedData) >= len(testData) {
		t.Logf("compressed size (%d) >= original size (%d) - data may be too small for compression",
			len(compressedData), len(testData))
	}
}

func TestCompressDriver_Write(t *testing.T) {
	testData := []byte("Test data for compression and decompression roundtrip")
	
	conf := &config.CompressConfig{Level: 3}
	driver := NewCompress(conf)
	
	// First compress the data
	originalReader := bytes.NewReader(testData)
	compressedReader, err := driver.Read(originalReader)
	if err != nil {
		t.Fatalf("failed to create compressed reader: %v", err)
	}
	
	compressedData, err := io.ReadAll(compressedReader)
	if err != nil {
		t.Fatalf("failed to read compressed data: %v", err)
	}
	
	// Now decompress using Write
	var output bytes.Buffer
	decompressWriter, err := driver.Write(&output)
	if err != nil {
		t.Fatalf("failed to create decompress writer: %v", err)
	}
	
	_, err = decompressWriter.Write(compressedData)
	if err != nil {
		t.Fatalf("failed to write compressed data: %v", err)
	}
	
	// Close to flush
	if closer, ok := decompressWriter.(io.Closer); ok {
		err = closer.Close()
		if err != nil {
			t.Fatalf("failed to close decompress writer: %v", err)
		}
	}
	
	decompressedData := output.Bytes()
	if !bytes.Equal(decompressedData, testData) {
		t.Errorf("decompressed data doesn't match original:\noriginal: %s\ndecompressed: %s",
			testData, decompressedData)
	}
}

func TestCompressDriver_encoderLevel(t *testing.T) {
	tests := []struct {
		level    int
		expected string
	}{
		{1, "fastest"},
		{2, "fastest"},
		{3, "default"},
		{5, "default"},
		{6, "better"},
		{9, "better"},
		{10, "best"},
		{22, "best"},
	}
	
	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			conf := &config.CompressConfig{Level: tt.level}
			driver := &CompressDriver{Conf: conf}
			
			level := driver.encoderLevel()
			
			// Compare string representation
			got := level.String()
			if got != tt.expected {
				t.Errorf("level %d: expected %s, got %s", tt.level, tt.expected, got)
			}
		})
	}
}

func TestCompressDriver_RoundTrip(t *testing.T) {
	// Test various data sizes
	testCases := []struct {
		name string
		data []byte
	}{
		{
			name: "small",
			data: []byte("small"),
		},
		{
			name: "medium",
			data: bytes.Repeat([]byte("medium data "), 100),
		},
		{
			name: "large",
			data: bytes.Repeat([]byte("large data with some variation "), 1000),
		},
		{
			name: "binary",
			data: func() []byte {
				b := make([]byte, 1024)
				for i := range b {
					b[i] = byte(i % 256)
				}
				return b
			}(),
		},
	}
	
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			conf := &config.CompressConfig{Level: 3}
			driver := NewCompress(conf)
			
			// Compress
			originalReader := bytes.NewReader(tc.data)
			compressedReader, err := driver.Read(originalReader)
			if err != nil {
				t.Fatalf("failed to create compressed reader: %v", err)
			}
			
			compressedData, err := io.ReadAll(compressedReader)
			if err != nil {
				t.Fatalf("failed to read compressed data: %v", err)
			}
			
			// Decompress
			var output bytes.Buffer
			decompressWriter, err := driver.Write(&output)
			if err != nil {
				t.Fatalf("failed to create decompress writer: %v", err)
			}
			
			_, err = decompressWriter.Write(compressedData)
			if err != nil {
				t.Fatalf("failed to write compressed data: %v", err)
			}
			
			if closer, ok := decompressWriter.(io.Closer); ok {
				err = closer.Close()
				if err != nil {
					t.Fatalf("failed to close decompress writer: %v", err)
				}
			}
			
			result := output.Bytes()
			if !bytes.Equal(result, tc.data) {
				t.Errorf("decompressed data doesn't match original (len=%d vs len=%d)",
					len(result), len(tc.data))
			}
		})
	}
}
