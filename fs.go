package fs

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// EnsureFile creates a file if it doesn't exist, with default mode 0644.
func EnsureFile(path string) error {
	return EnsureFileWithMode(path, 0644)
}

// EnsureFileWithMode creates a file if it doesn't exist, with the specified mode.
func EnsureFileWithMode(path string, mode os.FileMode) error {
	info, err := os.Stat(path)
	if err == nil {
		// Check if it's a directory
		if info.IsDir() {
			return fmt.Errorf("EnsureFile failed: %s is a directory", path)
		}

		return nil // File already exists
	}
	if !os.IsNotExist(err) {
		return fmt.Errorf("EnsureFile failed to check file existence: %w", err)
	}

	// Check if the directory exists
	dir := filepath.Dir(path)
	err = EnsureDirWithMode(dir, 0755)
	if err != nil {
		return fmt.Errorf("EnsureFile failed to ensure directory: %w", err)
	}

	file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY, mode)
	if err != nil {
		return fmt.Errorf("EnsureFile failed to create file: %w", err)
	}

	defer file.Close()

	return nil
}

// EnsureDir creates a directory if it doesn't exist, with default mode 0755.
func EnsureDir(path string) error {
	return EnsureDirWithMode(path, 0755|os.ModeDir)
}

// EnsureDirWithMode creates a directory if it doesn't exist, with the specified mode.
func EnsureDirWithMode(path string, mode os.FileMode) error {
	info, err := os.Stat(path)
	if err == nil {
		// Check if it's a file
		if !info.IsDir() {
			return fmt.Errorf("EnsureDir failed: %s is a file", path)
		}

		return nil // Directory already exists
	}
	if !os.IsNotExist(err) {
		return fmt.Errorf("EnsureDir failed to check directory existence: %w", err)
	}

	// Check if the parent directory exists
	parent := filepath.Dir(path)
	err = EnsureDirWithMode(parent, 0755)
	if err != nil {
		return fmt.Errorf("EnsureDir failed to ensure parent directory: %w", err)
	}

	err = os.Mkdir(path, mode)
	if err != nil {
		return fmt.Errorf("EnsureDir failed to create directory: %w", err)
	}

	return nil
}

// Exists checks if a file or directory exists.
func Exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, fmt.Errorf("Exists failed to check if file exists: %w", err)
}

// ReadDir reads the content of a directory and returns a list of file names.
// The order of the files is not guaranteed.
func ReadDir(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("ReadDir failed to open directory: %w", err)
	}
	defer file.Close()

	return file.Readdirnames(-1)
}

// ReadDirRec reads the content of a directory recursively and returns a list of file names.
// The order of the files is not guaranteed.
func ReadDirRec(path string) ([]string, error) {
	var files []string

	err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("ReadDirRec failed to walk directory: %w", err)
		}

		if info.IsDir() {
			return nil
		}

		files = append(files, path)
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("ReadDirRec failed to walk directory: %w", err)
	}

	return files, nil
}

// ReadJson reads the content of a JSON file and unmarshals it into a struct.
//
// Example:
//
//	var v MyStruct
//	err := ReadJson("file.json", &v)
//	if err != nil {
//	    fmt.Println(err)
//	    return
//	}
func ReadJson[T any](path string, v *T) error {
	content, err := ReadBytes(path)
	if err != nil {
		return fmt.Errorf("ReadJson failed to read file: %w", err)
	}

	return json.Unmarshal(content, v)
}

// ReadText reads the content of a file and returns it as a string.
// fmt.Println(string(content))
//
// Example:
//
//	content, err := ReadText("file.txt")
//	if err != nil {
//	    fmt.Println(err)
//	    return
//	}
func ReadText(path string) (string, error) {
	content, err := ReadBytes(path)
	if err != nil {
		return "", fmt.Errorf("ReadText failed to read file: %w", err)
	}

	return string(content), nil
}

// ReadBytes reads the content of a file and returns it as a byte slice.
//
// Example:
//
//	content, err := ReadBytes("file.txt")
//	if err != nil {
//	    fmt.Println(err)
//	    return
//	}
func ReadBytes(path string) ([]byte, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("ReadBytes failed to open file: %w", err)
	}
	defer file.Close()

	fileSize, err := GetSize(path)
	if err != nil {
		return nil, fmt.Errorf("ReadBytes failed to get file size: %w", err)
	}

	reader := bufio.NewReader(file)

	content := make([]byte, fileSize)
	totalBytesRead := 0
	for totalBytesRead < fileSize {
		bytesRead, err := reader.Read(content[totalBytesRead:])
		if err != nil {
			if err == io.EOF {
				break // End of file reached
			}
			return nil, fmt.Errorf("ReadBytes failed to read file: %w", err)
		}
		totalBytesRead += bytesRead
	}

	return content, nil
}

// GetSize returns the size of a file in bytes.
// Crucially, it returns int instead of int64. This is to make `make` easier to use
// with the result of this function.
func GetSize(path string) (int, error) {
	info, err := os.Stat(path)
	if err != nil {
		return 0, fmt.Errorf("GetSize failed to get file stat: %w", err)
	}

	return int(info.Size()), nil
}

// WriteJson writes a struct to a file as JSON.
func WriteJson[T any](path string, v T) error {
	content, err := json.Marshal(v)
	if err != nil {
		return fmt.Errorf("WriteJson failed to marshal content: %w", err)
	}

	return WriteBytes(path, content)
}

// WriteJson writes a struct to a file as JSON with a specific file mode.
func WriteJsonWithMode[T any](path string, v T, mode os.FileMode) error {
	content, err := json.Marshal(v)
	if err != nil {
		return fmt.Errorf("WriteJson failed to marshal content: %w", err)
	}

	return WriteBytesWithMode(path, content, mode)
}

// WriteText writes a string to a file.
func WriteText(path, content string) error {
	err := WriteBytes(path, []byte(content))
	if err != nil {
		return fmt.Errorf("WriteText failed to write content to file: %w", err)
	}

	return nil
}

// WriteText writes a string to a file with a specific file mode.
func WriteTextWithMode(path, content string, mode os.FileMode) error {
	err := WriteBytesWithMode(path, []byte(content), mode)
	if err != nil {
		return fmt.Errorf("WriteText failed to write content to file: %w", err)
	}

	return nil
}

// WriteBytes writes a byte slice to a file.
func WriteBytes(path string, content []byte) error {
	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("WriteBytes failed to create file: %w", err)
	}
	defer file.Close()

	_, err = file.Write(content)
	if err != nil {
		return fmt.Errorf("WriteBytes failed to write content to file: %w", err)
	}

	return nil
}

// WriteBytes writes a byte slice to a file with a specific file mode.
func WriteBytesWithMode(path string, content []byte, mode os.FileMode) error {
	err := os.WriteFile(path, content, mode)
	if err != nil {
		return fmt.Errorf("WriteBytes failed to write content to file: %w", err)
	}

	return nil
}
