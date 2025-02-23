package fs

import (
	"os"
	"sort"
	"testing"
)

// Prefer standard library functions internally in tests
// as using fs to test fs is a bit circular

func TestEnsureFile(t *testing.T) {
	// Expect to create a file if it does not exist already
	t.Run("file does not exist", func(t *testing.T) {
		path := "testfile_ensure_file_1.txt"
		defer os.Remove(path)

		err := EnsureFile(path)
		if err != nil {
			t.Errorf("EnsureFile failed: %v", err)
		}
	})

	// Expect to not create a file if it already exists
	t.Run("file exists", func(t *testing.T) {
		path := "testfile_ensure_file_2.txt"
		defer os.Remove(path)

		err := WriteText(path, "test content")
		if err != nil {
			t.Fatalf("WriteText failed: %v", err)
		}

		err = EnsureFile(path)
		if err != nil {
			t.Errorf("EnsureFile failed: %v", err)
		}

		content, err := ReadText(path)
		if err != nil {
			t.Errorf("ReadText failed: %v", err)
		}

		if content != "test content" {
			t.Errorf("Expected content to be 'test content', got '%s'", content)
		}
	})

	// Expect default mode to be 0644
	t.Run("file has mode 0644 by default", func(t *testing.T) {
		path := "testfile_ensure_file_3.txt"
		defer os.Remove(path)

		err := EnsureFile(path)
		if err != nil {
			t.Errorf("EnsureFile failed: %v", err)
		}

		info, err := os.Stat(path)
		if err != nil {
			t.Errorf("os.Stat failed: %v", err)
		}

		if info.Mode().Perm() != 0644 {
			t.Errorf("Expected file mode to be 0644, got %#o", info.Mode().Perm())
		}

		if info.IsDir() {
			t.Errorf("Expected file to not be a directory")
		}
	})
}

func TestEnsureDir(t *testing.T) {
	// Expect to create a directory if it does not exist already
	t.Run("directory does not exist", func(t *testing.T) {
		path := "ensure_dir_1"
		defer os.RemoveAll(path)

		err := EnsureDir(path)
		if err != nil {
			t.Errorf("EnsureDir failed: %v", err)
		}
	})

	// Expect to not create a directory if it already exists
	t.Run("directory exists", func(t *testing.T) {
		path := "ensure_dir_2"
		defer os.RemoveAll(path)

		err := os.Mkdir(path, 0755)
		if err != nil {
			t.Fatalf("os.Mkdir failed: %v", err)
		}

		err = EnsureDir(path)
		if err != nil {
			t.Errorf("EnsureDir failed: %v", err)
		}
	})

	// Expect default mode to be 0755
	t.Run("directory has mode 0755 by default", func(t *testing.T) {
		path := "ensure_dir_4"
		defer os.RemoveAll(path)

		err := EnsureDir(path)
		if err != nil {
			t.Errorf("EnsureDir failed: %v", err)
		}

		info, err := os.Stat(path)
		if err != nil {
			t.Errorf("os.Stat failed: %v", err)
		}

		if info.Mode().Perm() != 0755 {
			t.Errorf("Expected directory mode to be 0755, got %#o", info.Mode().Perm())
		}

		if !info.IsDir() {
			t.Errorf("Expected directory to be a directory")
		}
	})
}

func TestExists(t *testing.T) {
	// Expect to return true if a file exists
	t.Run("file exists", func(t *testing.T) {
		path := "exists_1.txt"
		defer os.Remove(path)

		err := WriteText(path, "test content")
		if err != nil {
			t.Fatalf("WriteText failed: %v", err)
		}

		exists, err := Exists(path)
		if err != nil {
			t.Errorf("Exists failed: %v", err)
		}

		if !exists {
			t.Errorf("Expected file to exist")
		}
	})

	// Expect to return true if a directory exists
	t.Run("directory exists", func(t *testing.T) {
		path := "exists_2"
		defer os.RemoveAll(path)

		err := os.Mkdir(path, 0755)
		if err != nil {
			t.Fatalf("os.Mkdir failed: %v", err)
		}

		exists, err := Exists(path)
		if err != nil {
			t.Errorf("Exists failed: %v", err)
		}

		if !exists {
			t.Errorf("Expected directory to exist")
		}
	})

	// Expect to return false if a file does not exist
	t.Run("file does not exist", func(t *testing.T) {
		path := "exists_3.txt"
		defer os.Remove(path)

		exists, err := Exists(path)
		if err != nil {
			t.Errorf("Exists failed: %v", err)
		}

		if exists {
			t.Errorf("Expected file to not exist")
		}
	})
}

func TestReadDir(t *testing.T) {
	// Expect to return a list of file names
	t.Run("read directory", func(t *testing.T) {
		path := "read_dir"
		defer os.RemoveAll(path) // Ensure cleanup

		err := os.Mkdir(path, 0755)
		if err != nil {
			t.Fatalf("os.Mkdir failed: %v", err)
		}

		files := []string{"read_dir_1.txt", "read_dir_2.txt", "read_dir_3.txt"}
		for _, file := range files {
			err := WriteText(path+"/"+file, "test content")
			if err != nil {
				t.Fatalf("WriteText failed: %v", err)
			}
		}

		names, err := ReadDir(path)
		if err != nil {
			t.Errorf("ReadDir failed: %v", err)
		}

		if len(names) != len(files) {
			t.Errorf("Expected %d files, got %d", len(files), len(names))
		}

		// The order of the files is not guaranteed
		sort.Strings(names)
		sort.Strings(files)

		for i, name := range names {
			if name != files[i] {
				t.Errorf("Expected file name to be %s, got %s", files[i], name)
			}
		}
	})
}

func TestReadDirRec(t *testing.T) {
	// Expect to return a list of file names
	t.Run("read directory recursively", func(t *testing.T) {
		path := "read_dir_rec"
		defer os.RemoveAll(path)

		err := os.Mkdir(path, 0755)
		if err != nil {
			t.Fatalf("os.Mkdir failed: %v", err)
		}

		files := []string{"read_dir_rec_1.txt", "read_dir_rec_2.txt", "read_dir_rec_3.txt"}
		for _, file := range files {
			err := WriteText(path+"/"+file, "test content")
			if err != nil {
				t.Fatalf("WriteText failed: %v", err)
			}
		}

		names, err := ReadDirRec(path)
		if err != nil {
			t.Errorf("ReadDirRec failed: %v", err)
		}

		if len(names) != len(files) {
			t.Errorf("Expected %d files, got %d", len(files), len(names))
		}

		// The order of the files is not guaranteed
		sort.Strings(names)
		sort.Strings(files)

		for i, name := range names {
			if name != path+"/"+files[i] {
				t.Errorf("Expected file name to be %s, got %s", path+"/"+files[i], name)
			}
		}
	})
}

func TestReadJson(t *testing.T) {
	// Expect to read and unmarshal a JSON file
	t.Run("read JSON file", func(t *testing.T) {
		path := "read_json.json"
		defer os.Remove(path)

		err := WriteText(path, `{"key": "value"}`)
		if err != nil {
			t.Fatalf("WriteText failed: %v", err)
		}

		var v struct {
			Key string `json:"key"`
		}

		err = ReadJson(path, &v)
		if err != nil {
			t.Errorf("ReadJson failed: %v", err)
		}

		if v.Key != "value" {
			t.Errorf("Expected key to be 'value', got '%s'", v.Key)
		}
	})
}

func TestReadText(t *testing.T) {
	// Expect to read the content of a file
	t.Run("read text file", func(t *testing.T) {
		path := "read_text.txt"
		defer os.Remove(path)

		err := WriteText(path, "test content")
		if err != nil {
			t.Fatalf("WriteText failed: %v", err)
		}

		content, err := ReadText(path)
		if err != nil {
			t.Errorf("ReadText failed: %v", err)
		}

		if content != "test content" {
			t.Errorf("Expected content to be 'test content', got '%s'", content)
		}
	})
}

func TestReadBytes(t *testing.T) {
	// Expect to read the content of a file
	t.Run("read bytes file", func(t *testing.T) {
		path := "read_bytes.txt"
		defer os.Remove(path)

		err := WriteText(path, "test content")
		if err != nil {
			t.Fatalf("WriteText failed: %v", err)
		}

		content, err := ReadBytes(path)
		if err != nil {
			t.Errorf("ReadBytes failed: %v", err)
		}

		if string(content) != "test content" {
			t.Errorf("Expected content to be 'test content', got '%s'", content)
		}
	})
}

func TestGetSize(t *testing.T) {
	// Expect to return the size of a file
	t.Run("get file size", func(t *testing.T) {
		path := "get_size.txt"
		defer os.Remove(path)

		err := WriteText(path, "test content")
		if err != nil {
			t.Fatalf("WriteText failed: %v", err)
		}

		size, err := GetSize(path)
		if err != nil {
			t.Errorf("GetSize failed: %v", err)
		}

		if size != 12 {
			t.Errorf("Expected size to be 12, got %d", size)
		}
	})
}

func TestWriteJson(t *testing.T) {
	// Expect to marshal and write a struct to a JSON file
	t.Run("write JSON file", func(t *testing.T) {
		path := "write_json.json"
		defer os.Remove(path)

		var v struct {
			Key string `json:"key"`
		}
		v.Key = "value"

		err := WriteJson(path, v)
		if err != nil {
			t.Errorf("WriteJson failed: %v", err)
		}

		content, err := ReadText(path)
		if err != nil {
			t.Errorf("ReadText failed: %v", err)
		}

		if content != `{"key":"value"}` {
			t.Errorf("Expected content to be '{\"key\":\"value\"}', got '%s'", content)
		}
	})
}

func TestWriteJsonWithMode(t *testing.T) {
	// Expect to write a JSON file with a specific mode
	t.Run("write JSON file with mode", func(t *testing.T) {
		path := "write_json_mode.json"
		defer os.Remove(path)

		var v struct {
			Key string `json:"key"`
		}
		v.Key = "value"

		err := WriteJsonWithMode(path, v, 0644)
		if err != nil {
			t.Errorf("WriteJsonWithMode failed: %v", err)
		}

		info, err := os.Stat(path)
		if err != nil {
			t.Errorf("os.Stat failed: %v", err)
		}

		if info.Mode().Perm() != 0644 {
			t.Errorf("Expected file mode to be 0644, got %#o", info.Mode().Perm())
		}
	})
}

func TestWriteText(t *testing.T) {
	// Expect to write content to a file
	t.Run("write text file", func(t *testing.T) {
		path := "write_text.txt"
		defer os.Remove(path)

		err := WriteText(path, "test content")
		if err != nil {
			t.Errorf("WriteText failed: %v", err)
		}

		content, err := ReadText(path)
		if err != nil {
			t.Errorf("ReadText failed: %v", err)
		}

		if content != "test content" {
			t.Errorf("Expected content to be 'test content', got '%s'", content)
		}
	})
}

func TestWriteTextWithMode(t *testing.T) {
	// Expect to write a text file with a specific mode
	t.Run("write text file with mode", func(t *testing.T) {
		path := "write_text_mode.txt"
		defer os.Remove(path)

		err := WriteTextWithMode(path, "test content", 0644)
		if err != nil {
			t.Errorf("WriteTextWithMode failed: %v", err)
		}

		info, err := os.Stat(path)
		if err != nil {
			t.Errorf("os.Stat failed: %v", err)
		}

		if info.Mode().Perm() != 0644 {
			t.Errorf("Expected file mode to be 0644, got %#o", info.Mode().Perm())
		}
	})
}

func TestWriteBytes(t *testing.T) {
	// Expect to write content to a file
	t.Run("write bytes file", func(t *testing.T) {
		path := "write_bytes.txt"
		defer os.Remove(path)

		err := WriteBytes(path, []byte("test content"))
		if err != nil {
			t.Errorf("WriteBytes failed: %v", err)
		}

		content, err := ReadText(path)
		if err != nil {
			t.Errorf("ReadText failed: %v", err)
		}

		if content != "test content" {
			t.Errorf("Expected content to be 'test content', got '%s'", content)
		}
	})
}

func TestWriteBytesWithMode(t *testing.T) {
	// Expect to write a byte slice to a file with a specific mode
	t.Run("write bytes file with mode", func(t *testing.T) {
		path := "write_bytes_mode.txt"
		defer os.Remove(path)

		err := WriteBytesWithMode(path, []byte("test content"), 0644)
		if err != nil {
			t.Errorf("WriteBytesWithMode failed: %v", err)
		}

		info, err := os.Stat(path)
		if err != nil {
			t.Errorf("os.Stat failed: %v", err)
		}

		if info.Mode().Perm() != 0644 {
			t.Errorf("Expected file mode to be 0644, got %#o", info.Mode().Perm())
		}
	})
}
